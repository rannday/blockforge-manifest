package manifest

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

// Validator validates manifest JSON documents against the repository schema.
type Validator struct {
	schema *jsonschema.Schema
}

// NewValidator compiles the schema at schemaPath.
func NewValidator(schemaPath string) (*Validator, error) {
	abs, err := filepath.Abs(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("resolve schema path: %w", err)
	}

	compiler := jsonschema.NewCompiler()
	compiler.AssertFormat()

	schema, err := compiler.Compile(abs)
	if err != nil {
		return nil, fmt.Errorf("compile schema %s: %w", schemaPath, err)
	}

	return &Validator{schema: schema}, nil
}

// ValidateFile validates one manifest JSON file.
func (v *Validator) ValidateFile(path string) error {
	fh, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fh.Close()

	return v.Validate(fh)
}

// Validate validates one manifest JSON stream.
func (v *Validator) Validate(r io.Reader) error {
	instance, err := jsonschema.UnmarshalJSON(r)
	if err != nil {
		return err
	}

	return v.schema.Validate(instance)
}

// FindRepoRoot walks up from start until manifest.schema.json is found.
func FindRepoRoot(start string) (string, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "manifest.schema.json")); err == nil {
			return dir, nil
		} else if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("manifest.schema.json not found from %s", start)
		}
		dir = parent
	}
}

// ExpandInputs expands shell-style globs. Unmatched patterns remain unchanged.
func ExpandInputs(args []string) []string {
	var inputs []string
	for _, arg := range args {
		matches, err := filepath.Glob(arg)
		if err == nil && len(matches) > 0 {
			inputs = append(inputs, matches...)
			continue
		}
		inputs = append(inputs, arg)
	}
	return inputs
}

// FormatValidationError formats jsonschema validation errors for CLI output.
func FormatValidationError(err error) []string {
	var validationErr *jsonschema.ValidationError
	if !errors.As(err, &validationErr) {
		return []string{err.Error()}
	}

	var lines []string
	collectValidationErrors(validationErr, &lines)
	return lines
}

func collectValidationErrors(err *jsonschema.ValidationError, lines *[]string) {
	if len(err.Causes) == 0 {
		path := formatInstancePath(err.InstanceLocation)
		*lines = append(*lines, fmt.Sprintf("%s: %s", path, validationMessage(err)))
		return
	}

	for _, cause := range err.Causes {
		collectValidationErrors(cause, lines)
	}
}

func validationMessage(err *jsonschema.ValidationError) string {
	output := err.BasicOutput()
	if output.Error != nil {
		return output.Error.String()
	}
	return err.Error()
}

func formatInstancePath(path []string) string {
	if len(path) == 0 {
		return "<root>"
	}
	return strings.Join(path, ".")
}
