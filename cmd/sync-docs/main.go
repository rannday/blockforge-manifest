package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rannday/blockforge-manifest/internal/manifest"
)

func main() {
	os.Exit(run())
}

func run() int {
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

	source := filepath.Join(repoRoot, "manifest.schema.json")
	docsDir := filepath.Join(repoRoot, "docs")
	target := filepath.Join(docsDir, "manifest.schema.json")
	nojekyll := filepath.Join(docsDir, ".nojekyll")

	if err := os.MkdirAll(docsDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: create docs dir: %v\n", err)
		return 1
	}

	if err := copyFile(source, target); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: sync schema: %v\n", err)
		return 1
	}

	fh, err := os.OpenFile(nojekyll, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: refresh .nojekyll: %v\n", err)
		return 1
	}
	if err := fh.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: close .nojekyll: %v\n", err)
		return 1
	}

	fmt.Println("Synced manifest.schema.json -> docs/manifest.schema.json")
	return 0
}

func copyFile(source, target string) error {
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(target)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}
	return dst.Sync()
}
