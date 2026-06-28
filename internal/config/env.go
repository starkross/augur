package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// envRefRe matches ${env:VAR} / ${ENV:VAR} placeholders, with an optional
// ":-default" fallback (e.g. ${env:HOST:-localhost}). Mirrors the syntax the
// OpenTelemetry Collector accepts and the patterns recognised by
// lib.is_env_var in the policy package.
var envRefRe = regexp.MustCompile(`\$\{(?:env|ENV):([A-Za-z_][A-Za-z0-9_]*)(:-([^}]*))?\}`)

// LoadEnvFile parses a .env-style file: blank lines and lines starting with
// '#' are ignored; every other line must be KEY=VALUE. Values may be wrapped
// in single or double quotes (which are stripped). The last definition of a
// key wins.
func LoadEnvFile(path string) (map[string]string, error) {
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("opening env file %q: %w", path, err)
	}
	defer func() { _ = f.Close() }()

	out := map[string]string{}
	sc := bufio.NewScanner(f)
	lineNo := 0
	for sc.Scan() {
		lineNo++
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		eq := strings.IndexByte(line, '=')
		if eq <= 0 {
			return nil, fmt.Errorf("%s:%d: expected KEY=VALUE", path, lineNo)
		}
		key := strings.TrimSpace(line[:eq])
		val := strings.TrimSpace(line[eq+1:])
		if !validEnvKey(key) {
			return nil, fmt.Errorf("%s:%d: invalid key %q", path, lineNo, key)
		}
		val = unquote(val)
		out[key] = val
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("reading env file %q: %w", path, err)
	}
	return out, nil
}

// LoadEnvFiles loads each path in order. Later files override earlier ones.
func LoadEnvFiles(paths []string) (map[string]string, error) {
	if len(paths) == 0 {
		return nil, nil
	}
	merged := map[string]string{}
	for _, p := range paths {
		m, err := LoadEnvFile(p)
		if err != nil {
			return nil, err
		}
		for k, v := range m {
			merged[k] = v
		}
	}
	return merged, nil
}

// SubstituteEnv walks v in place, replacing every ${env:VAR} / ${ENV:VAR}
// reference inside string values with the corresponding entry from env. A
// reference with a "${env:VAR:-default}" fallback resolves to the default
// when VAR is missing. References that resolve to nothing are left intact so
// that the policy layer can still detect them via lib.is_env_var.
func SubstituteEnv(v any, env map[string]string) any {
	if len(env) == 0 && !hasDefaults(v) {
		return v
	}
	switch x := v.(type) {
	case string:
		return substituteString(x, env)
	case map[string]any:
		for k, val := range x {
			x[k] = SubstituteEnv(val, env)
		}
		return x
	case []any:
		for i, val := range x {
			x[i] = SubstituteEnv(val, env)
		}
		return x
	default:
		return v
	}
}

func substituteString(s string, env map[string]string) string {
	if !strings.Contains(s, "${") {
		return s
	}
	return envRefRe.ReplaceAllStringFunc(s, func(match string) string {
		sub := envRefRe.FindStringSubmatch(match)
		name := sub[1]
		hasDefault := sub[2] != ""
		def := sub[3]
		if val, ok := env[name]; ok {
			return val
		}
		if hasDefault {
			return def
		}
		return match
	})
}

// hasDefaults is a cheap pre-check so that even with an empty env map we
// still process strings that carry "${env:VAR:-default}" fallbacks.
func hasDefaults(v any) bool {
	switch x := v.(type) {
	case string:
		return strings.Contains(x, ":-")
	case map[string]any:
		for _, val := range x {
			if hasDefaults(val) {
				return true
			}
		}
	case []any:
		for _, val := range x {
			if hasDefaults(val) {
				return true
			}
		}
	}
	return false
}

func validEnvKey(k string) bool {
	if k == "" {
		return false
	}
	for i, r := range k {
		switch {
		case r == '_':
		case r >= 'A' && r <= 'Z':
		case r >= 'a' && r <= 'z':
		case i > 0 && r >= '0' && r <= '9':
		default:
			return false
		}
	}
	return true
}

func unquote(s string) string {
	if len(s) >= 2 {
		first, last := s[0], s[len(s)-1]
		if (first == '"' && last == '"') || (first == '\'' && last == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
