//! Linting engine for OpenTelemetry Collector configs.
//!
//! Evaluates the same Rego policies as the Go implementation — they are shared
//! verbatim from `rules/policy/` — against a parsed config, and returns the
//! findings. The crate performs no I/O of its own: callers hand it strings, so
//! it works unchanged in a CLI, a server, or a wasm module in a browser.

pub mod config;

use std::collections::{BTreeMap, BTreeSet};

use serde::Serialize;
use serde_json::{Map, Value};

include!(concat!(env!("OUT_DIR"), "/policies.rs"));

/// Severity of a policy finding.
#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord, Serialize)]
#[serde(rename_all = "lowercase")]
pub enum Severity {
    /// A blocking violation that must be fixed.
    Deny,
    /// An advisory finding for best practices.
    Warn,
}

impl Severity {
    fn rule(self) -> &'static str {
        match self {
            Self::Deny => "data.main.deny",
            Self::Warn => "data.main.warn",
        }
    }
}

/// A single policy violation.
#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
pub struct Finding {
    pub rule_id: String,
    pub severity: Severity,
    pub message: String,
    pub file: String,
}

/// All findings for one evaluated config.
#[derive(Debug, Clone, Serialize)]
pub struct LintResult {
    pub file: String,
    #[serde(serialize_with = "null_if_empty")]
    pub findings: Vec<Finding>,
}

/// Go marshals a nil slice as `null`, and a clean config produces exactly that.
/// Emitting `[]` instead would silently change the shape consumers parse, so
/// the wart is preserved deliberately.
fn null_if_empty<S: serde::Serializer>(v: &[Finding], s: S) -> Result<S::Ok, S::Error> {
    if v.is_empty() {
        s.serialize_none()
    } else {
        s.collect_seq(v)
    }
}

#[derive(Debug)]
pub enum Error {
    Policy(String),
    Eval(String),
    Config(config::ConfigError),
}

impl std::fmt::Display for Error {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Self::Policy(e) => write!(f, "compiling policies: {e}"),
            Self::Eval(e) => write!(f, "evaluating policies: {e}"),
            Self::Config(e) => write!(f, "{e}"),
        }
    }
}

impl std::error::Error for Error {}

impl From<config::ConfigError> for Error {
    fn from(e: config::ConfigError) -> Self {
        Self::Config(e)
    }
}

/// Builder for a [`Linter`].
#[derive(Default)]
pub struct Builder {
    extra_policies: Vec<(String, String)>,
    disable_builtins: bool,
    skip_rules: BTreeSet<String>,
    severities: BTreeSet<Severity>,
    env: BTreeMap<String, String>,
}

impl Builder {
    /// Adds a custom policy module. `name` is used in compile errors.
    pub fn policy(mut self, name: impl Into<String>, source: impl Into<String>) -> Self {
        self.extra_policies.push((name.into(), source.into()));
        self
    }

    /// Excludes the bundled OTEL-* rules. At least one [`Builder::policy`] must
    /// then be supplied.
    pub fn without_builtins(mut self) -> Self {
        self.disable_builtins = true;
        self
    }

    /// Drops findings whose rule ID matches any of `ids`.
    pub fn skip_rules<I, S>(mut self, ids: I) -> Self
    where
        I: IntoIterator<Item = S>,
        S: Into<String>,
    {
        self.skip_rules
            .extend(ids.into_iter().map(Into::into).filter(|s| !s.is_empty()));
        self
    }

    /// Restricts findings to the given severities. Empty means all.
    pub fn severities<I: IntoIterator<Item = Severity>>(mut self, s: I) -> Self {
        self.severities.extend(s);
        self
    }

    /// Supplies values for `${env:VAR}` substitution.
    pub fn env(mut self, env: BTreeMap<String, String>) -> Self {
        self.env.extend(env);
        self
    }

    /// Compiles the policy set.
    pub fn build(self) -> Result<Linter, Error> {
        let mut engine = regorus::Engine::new();

        // OPA leaves an expression undefined when a builtin hits a runtime type
        // error and carries on; regorus aborts the whole evaluation by default.
        // Configs with an empty section (`grpc:` with no body) are routine, so
        // matching OPA here is what keeps findings from vanishing entirely.
        engine.set_strict_builtin_errors(false);

        let mut any = false;
        if !self.disable_builtins {
            for (name, src) in POLICIES {
                engine
                    .add_policy((*name).to_string(), (*src).to_string())
                    .map_err(|e| Error::Policy(e.to_string()))?;
                any = true;
            }
        }
        for (name, src) in &self.extra_policies {
            engine
                .add_policy(name.clone(), src.clone())
                .map_err(|e| Error::Policy(e.to_string()))?;
            any = true;
        }
        if !any {
            return Err(Error::Policy(
                "no policies configured: without_builtins requires at least one policy".into(),
            ));
        }

        Ok(Linter {
            engine,
            skip_rules: self.skip_rules,
            severities: self.severities,
            env: self.env,
        })
    }
}

/// Evaluates configs against a compiled policy set.
pub struct Linter {
    engine: regorus::Engine,
    skip_rules: BTreeSet<String>,
    severities: BTreeSet<Severity>,
    env: BTreeMap<String, String>,
}

impl Linter {
    /// Starts building a linter with the bundled rules.
    pub fn builder() -> Builder {
        Builder::default()
    }

    /// Evaluates an already-parsed config. `label` identifies the source in the
    /// returned findings.
    pub fn lint(&self, label: &str, mut input: Map<String, Value>) -> Result<LintResult, Error> {
        let mut as_value = Value::Object(std::mem::take(&mut input));
        config::substitute_env(&mut as_value, &self.env);

        // Clone so concurrent callers don't share evaluation state.
        let mut engine = self.engine.clone();
        engine.set_input(
            regorus::Value::from_json_str(&as_value.to_string())
                .map_err(|e| Error::Eval(e.to_string()))?,
        );

        let mut findings = Vec::new();
        for severity in [Severity::Deny, Severity::Warn] {
            let value = engine
                .eval_rule(severity.rule().to_string())
                .map_err(|e| Error::Eval(e.to_string()))?;
            collect(&value, severity, label, &mut findings);
        }

        // Go sorts by rule ID alone with an unstable sort, leaving ties in
        // arbitrary order. Sorting by message too makes the output stable.
        findings.sort_by(|a, b| a.rule_id.cmp(&b.rule_id).then_with(|| a.message.cmp(&b.message)));

        Ok(LintResult {
            file: label.to_string(),
            findings: self.filter(findings),
        })
    }

    /// Parses YAML and evaluates it.
    pub fn lint_yaml(&self, label: &str, data: &str) -> Result<LintResult, Error> {
        let parsed = config::parse_yaml(data)
            .map_err(|e| config::ConfigError::Io(format!("{label:?}: {e}")))?;
        self.lint(label, parsed)
    }

    /// Deep-merges several documents in order, then evaluates the result —
    /// matching the collector's own `--config` behavior.
    pub fn lint_merged(&self, sources: &[(String, String)]) -> Result<LintResult, Error> {
        let (first_label, first) = sources
            .first()
            .ok_or_else(|| config::ConfigError::Io("no config files provided".into()))?;

        let mut merged = config::parse_yaml(first)
            .map_err(|e| config::ConfigError::Io(format!("{first_label:?}: {e}")))?;
        for (label, raw) in &sources[1..] {
            let next = config::parse_yaml(raw)
                .map_err(|e| config::ConfigError::Io(format!("{label:?}: {e}")))?;
            config::merge_map(&mut merged, next);
        }

        let label = if sources.len() > 1 {
            let names: Vec<&str> = sources.iter().map(|(n, _)| n.as_str()).collect();
            format!("merged: {}", names.join(", "))
        } else {
            first_label.clone()
        };
        self.lint(&label, merged)
    }

    fn filter(&self, findings: Vec<Finding>) -> Vec<Finding> {
        if self.skip_rules.is_empty() && self.severities.is_empty() {
            return findings;
        }
        findings
            .into_iter()
            .filter(|f| !self.skip_rules.contains(&f.rule_id))
            .filter(|f| self.severities.is_empty() || self.severities.contains(&f.severity))
            .collect()
    }
}

fn collect(value: &regorus::Value, severity: Severity, file: &str, out: &mut Vec<Finding>) {
    let Ok(set) = value.as_set() else { return };
    for item in set.iter() {
        if let Ok(message) = item.as_string() {
            out.push(Finding {
                rule_id: extract_rule_id(message),
                severity,
                message: message.to_string(),
                file: file.to_string(),
            });
        }
    }
}

/// Rule IDs are the `OTEL-001` prefix before the first colon. Mirrors
/// `extractRuleID` in internal/engine/engine.go, including its shape check.
fn extract_rule_id(message: &str) -> String {
    if let Some(idx) = message.find(':') {
        if idx > 0 {
            let candidate = message[..idx].trim();
            if is_rule_id(candidate) {
                return candidate.to_string();
            }
        }
    }
    "UNKNOWN".to_string()
}

/// `^[A-Za-z][A-Za-z0-9_]*-[A-Za-z0-9_]+$`, without pulling in a regex crate.
fn is_rule_id(s: &str) -> bool {
    let Some((head, tail)) = s.split_once('-') else {
        return false;
    };
    let mut head_chars = head.chars();
    match head_chars.next() {
        Some(c) if c.is_ascii_alphabetic() => {}
        _ => return false,
    }
    head.len() >= 1
        && head_chars.all(|c| c.is_ascii_alphanumeric() || c == '_')
        && !tail.is_empty()
        && tail.chars().all(|c| c.is_ascii_alphanumeric() || c == '_')
}
