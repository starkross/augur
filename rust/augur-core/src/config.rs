//! YAML loading, env substitution and deep merge.
//!
//! Fidelity note: the Go implementation feeds policies through
//! `sigs.k8s.io/yaml`, which converts YAML to JSON using go-yaml v2 and so
//! applies YAML 1.1 scalar resolution. Rust YAML crates follow the YAML 1.2
//! core schema. The differences that change lint results are handled in
//! [`coerce_yaml11`]; the remainder are documented in rust/README.md.

use std::collections::BTreeMap;

use serde_json::{Map, Value};

/// Matches `${env:VAR}` / `${ENV:VAR}` with an optional `:-default` fallback,
/// mirroring `internal/config/env.go`.
const MAX_CONFIG_SIZE: usize = 10 << 20;

#[derive(Debug)]
pub enum ConfigError {
    TooLarge,
    Empty,
    Yaml(String),
    Io(String),
}

impl std::fmt::Display for ConfigError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Self::TooLarge => write!(f, "YAML input exceeds {} MB limit", MAX_CONFIG_SIZE >> 20),
            Self::Empty => write!(f, "empty or invalid YAML"),
            Self::Yaml(e) => write!(f, "parsing YAML: {e}"),
            Self::Io(e) => write!(f, "{e}"),
        }
    }
}

impl std::error::Error for ConfigError {}

/// go-yaml v2 resolves these bare scalars as booleans (YAML 1.1). serde_norway
/// follows YAML 1.2 core and leaves them as strings, which silently breaks
/// every rule comparing `== true` / `== false` — including OTEL-032, where a
/// missed match means an unreported TLS bypass.
fn coerce_yaml11(v: serde_norway::Value) -> serde_norway::Value {
    use serde_norway::Value as V;
    match v {
        V::String(s) => match s.as_str() {
            "y" | "Y" | "yes" | "Yes" | "YES" | "on" | "On" | "ON" => V::Bool(true),
            "n" | "N" | "no" | "No" | "NO" | "off" | "Off" | "OFF" => V::Bool(false),
            _ => V::String(s),
        },
        V::Sequence(xs) => V::Sequence(xs.into_iter().map(coerce_yaml11).collect()),
        V::Mapping(m) => V::Mapping(
            m.into_iter()
                .map(|(k, val)| (coerce_yaml11(k), coerce_yaml11(val)))
                .collect(),
        ),
        other => other,
    }
}

/// Parses YAML into a JSON object, matching what the Go loader hands the
/// policies.
pub fn parse_yaml(data: &str) -> Result<Map<String, Value>, ConfigError> {
    if data.len() > MAX_CONFIG_SIZE {
        return Err(ConfigError::TooLarge);
    }
    let parsed: serde_norway::Value =
        serde_norway::from_str(data).map_err(|e| ConfigError::Yaml(e.to_string()))?;
    let parsed = coerce_yaml11(parsed);
    let json: Value =
        serde_json::to_value(&parsed).map_err(|e| ConfigError::Yaml(e.to_string()))?;
    match json {
        Value::Object(m) => Ok(m),
        Value::Null => Err(ConfigError::Empty),
        _ => Err(ConfigError::Empty),
    }
}

/// Deep-merges `src` into `dst`: maps merge recursively, scalars and sequences
/// are replaced by the later document. Matches the collector's confmap
/// semantics and `internal/config/loader.go`.
pub fn merge_map(dst: &mut Map<String, Value>, src: Map<String, Value>) {
    for (k, sv) in src {
        match (dst.get_mut(&k), sv) {
            (Some(Value::Object(dm)), Value::Object(sm)) => {
                let mut inner = std::mem::take(dm);
                merge_map(&mut inner, sm);
                *dm = inner;
            }
            (_, sv) => {
                dst.insert(k, sv);
            }
        }
    }
}

/// Replaces `${env:VAR}` references in every string value. A reference with a
/// `:-default` fallback resolves to the default when the variable is missing;
/// references that resolve to nothing are left intact so the policy layer can
/// still detect them via `lib.is_env_var`.
pub fn substitute_env(v: &mut Value, env: &BTreeMap<String, String>) {
    match v {
        Value::String(s) => {
            if s.contains("${") {
                *s = substitute_string(s, env);
            }
        }
        Value::Object(m) => {
            for (_, val) in m.iter_mut() {
                substitute_env(val, env);
            }
        }
        Value::Array(xs) => {
            for val in xs.iter_mut() {
                substitute_env(val, env);
            }
        }
        _ => {}
    }
}

fn substitute_string(s: &str, env: &BTreeMap<String, String>) -> String {
    let bytes = s.as_bytes();
    let mut out = String::with_capacity(s.len());
    let mut i = 0;
    while i < bytes.len() {
        match parse_ref(s, i) {
            Some((name, default, end)) => {
                if let Some(val) = env.get(name) {
                    out.push_str(val);
                } else if let Some(d) = default {
                    out.push_str(d);
                } else {
                    out.push_str(&s[i..end]); // leave the placeholder intact
                }
                i = end;
            }
            None => {
                let ch = s[i..].chars().next().unwrap();
                out.push(ch);
                i += ch.len_utf8();
            }
        }
    }
    out
}

/// Parses `${env:NAME}` or `${env:NAME:-default}` at `start`.
/// Returns (name, default, end-offset).
fn parse_ref(s: &str, start: usize) -> Option<(&str, Option<&str>, usize)> {
    let rest = s.get(start..)?;
    let body = rest
        .strip_prefix("${env:")
        .or_else(|| rest.strip_prefix("${ENV:"))?;
    let prefix_len = rest.len() - body.len();
    let close = body.find('}')?;
    let inner = &body[..close];
    let end = start + prefix_len + close + 1;

    let (name, default) = match inner.find(":-") {
        Some(idx) => (&inner[..idx], Some(&inner[idx + 2..])),
        None => (inner, None),
    };
    if !valid_env_key(name) {
        return None;
    }
    Some((name, default, end))
}

fn valid_env_key(k: &str) -> bool {
    let mut chars = k.chars();
    match chars.next() {
        Some(c) if c == '_' || c.is_ascii_alphabetic() => {}
        _ => return false,
    }
    chars.all(|c| c == '_' || c.is_ascii_alphanumeric())
}

/// Parses a .env-style file: blank lines and `#` comments are ignored, every
/// other line must be KEY=VALUE. Quoted values are unquoted. Last wins.
pub fn parse_env_file(path: &str, contents: &str) -> Result<BTreeMap<String, String>, ConfigError> {
    let mut out = BTreeMap::new();
    for (idx, raw) in contents.lines().enumerate() {
        let line = raw.trim();
        if line.is_empty() || line.starts_with('#') {
            continue;
        }
        let line = line.strip_prefix("export ").unwrap_or(line);
        let eq = match line.find('=') {
            Some(0) | None => {
                return Err(ConfigError::Io(format!(
                    "{path}:{}: expected KEY=VALUE",
                    idx + 1
                )))
            }
            Some(i) => i,
        };
        let key = line[..eq].trim();
        if !valid_env_key(key) {
            return Err(ConfigError::Io(format!(
                "{path}:{}: invalid key {key:?}",
                idx + 1
            )));
        }
        out.insert(key.to_string(), unquote(line[eq + 1..].trim()).to_string());
    }
    Ok(out)
}

fn unquote(s: &str) -> &str {
    let b = s.as_bytes();
    if b.len() >= 2 {
        let (f, l) = (b[0], b[b.len() - 1]);
        if (f == b'"' && l == b'"') || (f == b'\'' && l == b'\'') {
            return &s[1..s.len() - 1];
        }
    }
    s
}
