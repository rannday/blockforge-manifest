package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rannday/blockforge-manifest/internal/manifest"
)

func main() {
	os.Exit(run())
}

func run() int {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "--" {
		args = args[1:]
	}

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: go run ./cmd/validate-manifest -- <manifest.json> [more.json...]")
		return 2
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		return 1
	}

	repoRoot, err := manifest.FindRepoRoot(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		return 1
	}

	validator, err := manifest.NewValidator(filepath.Join(repoRoot, "manifest.schema.json"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		return 1
	}

	ok := true
	for _, path := range manifest.ExpandInputs(args) {
		err := validator.ValidateFile(path)
		if err == nil {
			fmt.Printf("OK: %s\n", path)
			continue
		}

		ok = false
		fmt.Printf("INVALID: %s\n", path)
		for _, line := range manifest.FormatValidationError(err) {
			fmt.Println(line)
		}
	}

	if ok {
		return 0
	}
	return 1
}
