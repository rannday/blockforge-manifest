package manifest

import (
	"path/filepath"
	"testing"
)

func TestValidFixtures(t *testing.T) {
	validator := newTestValidator(t)

	for _, pattern := range []string{
		filepath.Join("..", "..", "examples", "*.json"),
		filepath.Join("..", "..", "testdata", "valid", "*.json"),
	} {
		for _, path := range globTestFiles(t, pattern) {
			if err := validator.ValidateFile(path); err != nil {
				t.Errorf("%s should be valid: %v", path, err)
			}
		}
	}
}

func TestInvalidFixtures(t *testing.T) {
	validator := newTestValidator(t)

	for _, path := range globTestFiles(t, filepath.Join("..", "..", "testdata", "invalid", "*.json")) {
		if err := validator.ValidateFile(path); err == nil {
			t.Errorf("%s should be invalid", path)
		}
	}
}

func newTestValidator(t *testing.T) *Validator {
	t.Helper()

	validator, err := NewValidator(filepath.Join("..", "..", "manifest.schema.json"))
	if err != nil {
		t.Fatal(err)
	}
	return validator
}

func globTestFiles(t *testing.T, pattern string) []string {
	t.Helper()

	matches, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) == 0 {
		t.Fatalf("no files match %s", pattern)
	}
	return matches
}
