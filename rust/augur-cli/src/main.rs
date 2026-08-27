//! Lint OpenTelemetry Collector configs for best practices.

mod output;

use std::collections::BTreeMap;
use std::io::{IsTerminal, Read, Write};
use std::path::PathBuf;
use std::process::ExitCode;

use augur_core::{config, Linter, Severity};
use clap::Parser;
use output::Format;

#[derive(Parser)]
#[command(
    name = "augur",
    version,
    about = "Lint OpenTelemetry Collector configs for best practices",
    long_about = "Lint OpenTelemetry Collector configs for best practices.\n\n\
When multiple files are provided, they are deep-merged into a single effective \
config before linting, matching the collector's own --config behavior (maps \
merge recursively; slices and scalars are replaced by the later file).\n\n\
Use \"-\" to read a config from stdin (e.g. cat config.yaml | augur -)."
)]
struct Cli {
    /// Config files to lint; "-" reads from stdin
    #[arg(required = true, value_name = "config.yaml")]
    files: Vec<String>,

    /// Output format: text, json, github
    #[arg(short, long, default_value = "text")]
    output: String,

    /// Treat warnings as errors
    #[arg(short, long)]
    strict: bool,

    /// Only show failures, suppress warnings
    #[arg(short, long)]
    quiet: bool,

    /// Comma-separated rule IDs to skip
    #[arg(short = 'k', long, default_value = "")]
    skip: String,

    /// Disable colored output
    #[arg(long)]
    no_color: bool,

    /// Additional policy directory (merged with built-in rules)
    #[arg(short, long)]
    policy: Option<PathBuf>,

    /// Path to a .env-style file used to resolve ${env:VAR} placeholders
    /// (repeatable; later files override)
    #[arg(long = "env-file")]
    env_files: Vec<PathBuf>,
}

fn main() -> ExitCode {
    let cli = Cli::parse();
    match run(&cli) {
        Ok(true) => ExitCode::SUCCESS,
        Ok(false) => {
            eprintln!("Error: lint failures detected");
            ExitCode::FAILURE
        }
        Err(e) => {
            eprintln!("Error: {e}");
            ExitCode::FAILURE
        }
    }
}

fn run(cli: &Cli) -> Result<bool, Box<dyn std::error::Error>> {
    let format: Format = cli.output.parse()?;

    let mut builder = Linter::builder().skip_rules(
        cli.skip
            .split(',')
            .map(str::trim)
            .filter(|s| !s.is_empty())
            .map(str::to_string),
    );
    if cli.quiet {
        builder = builder.severities([Severity::Deny]);
    }
    if let Some(dir) = &cli.policy {
        for (name, src) in read_policy_dir(dir)? {
            builder = builder.policy(name, src);
        }
    }
    if !cli.env_files.is_empty() {
        let mut env = BTreeMap::new();
        for path in &cli.env_files {
            let raw = std::fs::read_to_string(path)
                .map_err(|e| format!("opening env file {}: {e}", path.display()))?;
            env.extend(config::parse_env_file(&path.to_string_lossy(), &raw)?);
        }
        builder = builder.env(env);
    }

    let linter = builder.build()?;

    let mut stdin_seen = false;
    let mut sources = Vec::with_capacity(cli.files.len());
    for path in &cli.files {
        if path == "-" {
            if stdin_seen {
                return Err("stdin (-) can only be specified once".into());
            }
            stdin_seen = true;
            let mut buf = String::new();
            std::io::stdin().read_to_string(&mut buf)?;
            sources.push(("<stdin>".to_string(), buf));
        } else {
            let raw = std::fs::read_to_string(path)
                .map_err(|e| format!("opening {path:?}: {e}"))?;
            sources.push((path.clone(), raw));
        }
    }

    let result = linter.lint_merged(&sources)?;
    let failed = output::has_failures(&result.findings, cli.strict);

    let no_color = cli.no_color || !std::io::stdout().is_terminal();
    let mut stdout = std::io::stdout().lock();
    output::write(&mut stdout, format, std::slice::from_ref(&result), no_color)?;
    stdout.flush()?;

    Ok(!failed)
}

/// Loads every non-test .rego file under `dir`, recursively.
fn read_policy_dir(dir: &PathBuf) -> Result<Vec<(String, String)>, Box<dyn std::error::Error>> {
    let mut out = Vec::new();
    let mut stack = vec![dir.clone()];
    while let Some(current) = stack.pop() {
        for entry in std::fs::read_dir(&current)
            .map_err(|e| format!("reading policy dir {}: {e}", current.display()))?
        {
            let path = entry?.path();
            if path.is_dir() {
                stack.push(path);
            } else if path.extension().is_some_and(|e| e == "rego")
                && !path.file_name().unwrap().to_string_lossy().ends_with("_test.rego")
            {
                let src = std::fs::read_to_string(&path)?;
                out.push((path.to_string_lossy().to_string(), src));
            }
        }
    }
    out.sort();
    Ok(out)
}
