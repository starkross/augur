use augur_core::{config, Linter, Severity};
use std::collections::BTreeMap;

fn env(pairs: &[(&str, &str)]) -> BTreeMap<String, String> {
    pairs.iter().map(|(k, v)| (k.to_string(), v.to_string())).collect()
}

#[test]
fn yaml11_booleans_match_go() {
    // go-yaml v2 resolves these as booleans; YAML 1.2 would leave them strings,
    // which silently breaks rules comparing `== true`.
    let m = config::parse_yaml("a: yes\nb: no\nc: on\nd: off\ne: true\nf: False\ng: maybe").unwrap();
    assert_eq!(m["a"], serde_json::json!(true));
    assert_eq!(m["b"], serde_json::json!(false));
    assert_eq!(m["c"], serde_json::json!(true));
    assert_eq!(m["d"], serde_json::json!(false));
    assert_eq!(m["e"], serde_json::json!(true));
    assert_eq!(m["f"], serde_json::json!(false));
    assert_eq!(m["g"], serde_json::json!("maybe"));
}

#[test]
fn empty_document_is_rejected() {
    assert!(config::parse_yaml("").is_err());
}

#[test]
fn deep_merge_replaces_scalars_and_sequences() {
    let mut base = config::parse_yaml("a:\n  x: 1\n  keep: 2\nlist: [1, 2]\ns: old").unwrap();
    let next = config::parse_yaml("a:\n  keep: 99\n  z: 3\nlist: [9]\ns: new").unwrap();
    config::merge_map(&mut base, next);
    assert_eq!(base["a"]["x"], serde_json::json!(1));
    assert_eq!(base["a"]["keep"], serde_json::json!(99));
    assert_eq!(base["a"]["z"], serde_json::json!(3));
    // Sequences and scalars are replaced wholesale, not merged.
    assert_eq!(base["list"], serde_json::json!([9]));
    assert_eq!(base["s"], serde_json::json!("new"));
}

#[test]
fn yaml11_coerces_mapping_keys_too() {
    // Surprising but deliberate: go-yaml v2 resolves a bare `y` key as the
    // boolean true, which sigs.k8s.io/yaml then renders as the string "true".
    // Diverging here would change which keys the policies can see.
    let m = config::parse_yaml("y: 2\nn: 3\nkeep: 4").unwrap();
    assert_eq!(m["true"], serde_json::json!(2));
    assert_eq!(m["false"], serde_json::json!(3));
    assert_eq!(m["keep"], serde_json::json!(4));
    assert!(m.get("y").is_none());
}

#[test]
fn env_substitution() {
    let e = env(&[("HOST", "example.com")]);
    let mut v = serde_json::json!({
        "resolved": "${env:HOST}:4317",
        "upper":    "${ENV:HOST}",
        "default":  "${env:MISSING:-fallback}",
        "unresolved": "${env:NOPE}",
        "untouched": "plain"
    });
    config::substitute_env(&mut v, &e);
    assert_eq!(v["resolved"], serde_json::json!("example.com:4317"));
    assert_eq!(v["upper"], serde_json::json!("example.com"));
    assert_eq!(v["default"], serde_json::json!("fallback"));
    // Left intact so lib.is_env_var can still see it.
    assert_eq!(v["unresolved"], serde_json::json!("${env:NOPE}"));
    assert_eq!(v["untouched"], serde_json::json!("plain"));
}

#[test]
fn env_file_parsing() {
    let src = "# comment\n\nexport A=1\nB = \"quoted\"\nC='single'\nA=2\n";
    let got = config::parse_env_file("t.env", src).unwrap();
    assert_eq!(got["A"], "2"); // last wins
    assert_eq!(got["B"], "quoted");
    assert_eq!(got["C"], "single");
}

#[test]
fn env_file_rejects_bad_lines() {
    assert!(config::parse_env_file("t.env", "not-a-pair\n").is_err());
    assert!(config::parse_env_file("t.env", "1BAD=x\n").is_err());
}

const BAD: &str = r#"
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
exporters:
  debug: {}
processors: {}
service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: []
      exporters: [debug]
"#;

#[test]
fn lints_and_extracts_rule_ids() {
    let l = Linter::builder().build().unwrap();
    let r = l.lint_yaml("bad.yaml", BAD).unwrap();
    assert!(!r.findings.is_empty());
    assert_eq!(r.file, "bad.yaml");
    for f in &r.findings {
        assert!(f.rule_id.starts_with("OTEL-"), "unexpected id {}", f.rule_id);
        assert_eq!(f.file, "bad.yaml");
    }
    // Findings are sorted and stable.
    let mut sorted = r.findings.clone();
    sorted.sort_by(|a, b| a.rule_id.cmp(&b.rule_id).then(a.message.cmp(&b.message)));
    assert_eq!(r.findings, sorted);
}

#[test]
fn empty_section_does_not_abort_evaluation() {
    // `grpc:` with no body parses as null. OPA leaves the failing expression
    // undefined and carries on; regorus aborts unless strict errors are off.
    let cfg = "receivers:\n  otlp:\n    protocols:\n      grpc:\nservice:\n  pipelines:\n    traces:\n      receivers: [otlp]\n      exporters: [debug]\nexporters:\n  debug: {}\n";
    let l = Linter::builder().build().unwrap();
    let r = l.lint_yaml("null.yaml", cfg).expect("must not abort");
    assert!(!r.findings.is_empty());
}

#[test]
fn skip_and_severity_filters() {
    let l = Linter::builder().build().unwrap();
    let all = l.lint_yaml("c.yaml", BAD).unwrap();
    let skipped_id = all.findings[0].rule_id.clone();

    let l2 = Linter::builder().skip_rules([skipped_id.clone()]).build().unwrap();
    let r2 = l2.lint_yaml("c.yaml", BAD).unwrap();
    assert!(r2.findings.iter().all(|f| f.rule_id != skipped_id));

    let l3 = Linter::builder().severities([Severity::Deny]).build().unwrap();
    let r3 = l3.lint_yaml("c.yaml", BAD).unwrap();
    assert!(r3.findings.iter().all(|f| f.severity == Severity::Deny));
}

#[test]
fn merged_configs_use_a_joined_label() {
    let l = Linter::builder().build().unwrap();
    let sources = vec![
        ("a.yaml".to_string(), BAD.to_string()),
        ("b.yaml".to_string(), "processors:\n  batch: {}\n".to_string()),
    ];
    let r = l.lint_merged(&sources).unwrap();
    assert_eq!(r.file, "merged: a.yaml, b.yaml");
}

#[test]
fn without_builtins_requires_a_policy() {
    assert!(Linter::builder().without_builtins().build().is_err());
}

#[test]
fn custom_policy_is_evaluated() {
    let l = Linter::builder()
        .policy(
            "custom.rego",
            "package main\nimport future.keywords.if\nimport future.keywords.contains\n\
             deny contains msg if { input.nope; msg := \"CUSTOM-001: triggered\" }\n",
        )
        .build()
        .unwrap();
    let r = l.lint_yaml("c.yaml", "nope: true\nservice: {}\n").unwrap();
    assert!(r.findings.iter().any(|f| f.rule_id == "CUSTOM-001"));
}
