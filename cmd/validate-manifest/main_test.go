package main

import (
	"os"
	"testing"
)

func TestRunValidFixtures(t *testing.T) {
	t.Chdir("../..")
	withArgs(t, "validate-manifest", "--", "examples/manifest.example.json", "testdata/valid/manifest.v1.json")

	if code := run(); code != 0 {
		t.Fatalf("run() = %d, want 0", code)
	}
}

func TestRunInvalidFixtures(t *testing.T) {
	t.Chdir("../..")
	withArgs(t, "validate-manifest", "--", "testdata/invalid/*.json")

	if code := run(); code != 1 {
		t.Fatalf("run() = %d, want 1", code)
	}
}

func TestRunNoArgs(t *testing.T) {
	t.Chdir("../..")
	withArgs(t, "validate-manifest")

	if code := run(); code != 2 {
		t.Fatalf("run() = %d, want 2", code)
	}
}

func withArgs(t *testing.T, args ...string) {
	t.Helper()

	oldArgs := os.Args
	os.Args = args
	t.Cleanup(func() {
		os.Args = oldArgs
	})
}
