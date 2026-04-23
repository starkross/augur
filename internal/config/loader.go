package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"sigs.k8s.io/yaml"
)

const maxConfigSize = 10 << 20 // 10 MB

// Stdin is the reader used when a config path is "-".
// It defaults to os.Stdin and can be overridden in tests.
var Stdin io.Reader = os.Stdin

// LoadYAML reads and parses a YAML file into a generic map.
func LoadYAML(path string) (map[string]any, error) {
	if path == "-" {
		return LoadReader(Stdin, "<stdin>")
	}

	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("opening %q: %w", path, err)
	}
	defer func() { _ = f.Close() }()

	return LoadReader(f, fmt.Sprintf("%q", path))
}

// LoadReader parses YAML from r into a generic map.
// The label is used in error messages to identify the source (e.g. a file path or "<stdin>").
func LoadReader(r io.Reader, label string) (map[string]any, error) {
	data, err := io.ReadAll(io.LimitReader(r, int64(maxConfigSize)+1))
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", label, err)
	}
	if len(data) > maxConfigSize {
		return nil, fmt.Errorf("%s: input exceeds %d MB limit", label, maxConfigSize>>20)
	}

	var out map[string]any
	if err := yaml.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parsing YAML in %s: %w", label, err)
	}

	if out == nil {
		return nil, fmt.Errorf("empty or invalid YAML in %s", label)
	}

	return out, nil
}

// LoadMerged loads each path and deep-merges them in order, matching the
// OpenTelemetry Collector's confmap semantics: maps merge recursively, scalars
// and slices are replaced by the later file.
func LoadMerged(paths []string) (map[string]any, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("no config files provided")
	}

	stdinCount := 0
	for _, p := range paths {
		if p == "-" {
			stdinCount++
		}
	}
	if stdinCount > 1 {
		return nil, fmt.Errorf("stdin (-) can only be specified once")
	}

	merged, err := LoadYAML(paths[0])
	if err != nil {
		return nil, err
	}

	for _, p := range paths[1:] {
		next, err := LoadYAML(p)
		if err != nil {
			return nil, err
		}
		merged = mergeMap(merged, next)
	}

	return merged, nil
}

func mergeMap(dst, src map[string]any) map[string]any {
	if dst == nil {
		dst = map[string]any{}
	}
	for k, sv := range src {
		dv, ok := dst[k]
		if !ok {
			dst[k] = sv
			continue
		}
		if dm, dstIsMap := dv.(map[string]any); dstIsMap {
			if sm, srcIsMap := sv.(map[string]any); srcIsMap {
				dst[k] = mergeMap(dm, sm)
				continue
			}
		}
		dst[k] = sv
	}
	return dst
}
