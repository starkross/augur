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
