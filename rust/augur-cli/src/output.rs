//! Output formatters. Byte-for-byte compatible with internal/output/output.go.

use augur_core::{Finding, LintResult, Severity};
use serde::Serialize;
use std::io::{self, Write};

#[derive(Clone, Copy, PartialEq, Eq)]
pub enum Format {
    Text,
    Json,
    GitHub,
}

impl std::str::FromStr for Format {
    type Err = String;
    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s {
            "text" => Ok(Self::Text),
            "json" => Ok(Self::Json),
            "github" => Ok(Self::GitHub),
            other => Err(format!("unknown output format: {other:?}")),
        }
    }
}

pub fn write(w: &mut impl Write, format: Format, results: &[LintResult], no_color: bool) -> io::Result<()> {
    match format {
        Format::Text => text(w, results, no_color),
        Format::Json => json(w, results),
        Format::GitHub => github(w, results),
    }
}

fn paint(s: &str, code: &str, no_color: bool) -> String {
    if no_color {
        s.to_string()
    } else {
        format!("\x1b[{code}m{s}\x1b[0m")
    }
}

fn text(w: &mut impl Write, results: &[LintResult], no_color: bool) -> io::Result<()> {
    let (mut denies, mut warns) = (0usize, 0usize);
    for r in results {
        if r.findings.is_empty() {
            continue;
        }
        writeln!(w, "{}", paint(&r.file, "1", no_color))?;
        for f in &r.findings {
            match f.severity {
                Severity::Deny => {
                    writeln!(w, "  {} {}", paint("FAIL", "0;31", no_color), f.message)?;
                    denies += 1;
                }
                Severity::Warn => {
                    writeln!(w, "  {} {}", paint("WARN", "0;33", no_color), f.message)?;
                    warns += 1;
                }
            }
        }
        writeln!(w)?;
    }
    let summary = if denies + warns == 0 {
        paint("✓ All checks passed", "0;32", no_color)
    } else if denies == 0 {
        paint(&format!("⚠ {warns} warning(s), 0 failure(s)"), "0;33", no_color)
    } else {
        paint(&format!("✗ {denies} failure(s), {warns} warning(s)"), "0;31", no_color)
    };
    writeln!(w, "{summary}")
}

#[derive(Serialize)]
struct JsonOutput<'a> {
    files: &'a [LintResult],
    summary: JsonSummary,
}

#[derive(Serialize, Default)]
struct JsonSummary {
    total_files: usize,
    failures: usize,
    warnings: usize,
    passed: usize,
}

fn json(w: &mut impl Write, results: &[LintResult]) -> io::Result<()> {
    let mut summary = JsonSummary::default();
    for r in results {
        let mut has_issue = false;
        for f in &r.findings {
            match f.severity {
                Severity::Deny => summary.failures += 1,
                Severity::Warn => summary.warnings += 1,
            }
            has_issue = true;
        }
        summary.total_files += 1;
        if !has_issue {
            summary.passed += 1;
        }
    }
    let out = JsonOutput { files: results, summary };
    // Go's json.Encoder uses two-space indent and appends a newline.
    let mut buf = Vec::new();
    let formatter = serde_json::ser::PrettyFormatter::with_indent(b"  ");
    let mut ser = serde_json::Serializer::with_formatter(&mut buf, formatter);
    serde::Serialize::serialize(&out, &mut ser).map_err(io::Error::other)?;
    w.write_all(&buf)?;
    writeln!(w)
}

fn github(w: &mut impl Write, results: &[LintResult]) -> io::Result<()> {
    for r in results {
        for f in &r.findings {
            let level = match f.severity {
                Severity::Deny => "error",
                Severity::Warn => "warning",
            };
            writeln!(
                w,
                "::{level} file={},title={}::{}",
                r.file,
                f.rule_id,
                escape(&f.message)
            )?;
        }
    }
    Ok(())
}

fn escape(s: &str) -> String {
    s.replace('%', "%25").replace('\n', "%0A").replace('\r', "%0D")
}

/// True when any finding should fail the run.
pub fn has_failures(findings: &[Finding], strict: bool) -> bool {
    findings
        .iter()
        .any(|f| f.severity == Severity::Deny || (strict && f.severity == Severity::Warn))
}
