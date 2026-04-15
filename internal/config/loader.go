package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"sigs.k8s.io/yaml"
)

const maxConfigSize = 10 << 20 // 10 MB

// LoadYAML reads and parses a YAML file into a generic map.
func LoadYAML(path string) (out map[string]any, err error) {
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("opening %q: %w", path, err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("closing %q: %w", path, cerr)
		}
	}()

	data, err := io.ReadAll(io.LimitReader(f, maxConfigSize+1))
	if err != nil {
		return nil, fmt.Errorf("reading %q: %w", path, err)
	}
	if len(data) > maxConfigSize {
		return nil, fmt.Errorf("%q: file exceeds %d MB limit", path, maxConfigSize>>20)
	}

	if err := yaml.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parsing YAML in %q: %w", path, err)
	}

	if out == nil {
		return nil, fmt.Errorf("empty or invalid YAML in %q", path)
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
