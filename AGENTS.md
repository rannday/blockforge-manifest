# AGENTS.md

This file applies to this repository.

## Project

- This repo owns `manifest.schema.json`, the canonical Blockforge/Varda server manifest contract.
- Keep `docs/manifest.schema.json` synced from the root schema with `go run ./cmd/sync-docs`.
- Do not edit `docs/index.html` unless the task explicitly asks for Pages UI changes.
- Do not replace the JSON Schema contract with Go structs. Go code is tooling only.

## Go

- Run `gofmt` on changed Go files.
- Run `go test ./...` after Go or schema changes.
- Do not add dependencies unless needed for the task.
- Keep CLI output stable and readable; README examples may depend on it.

## Schema

- Add or update fixtures when changing validation behavior.
- Put valid examples under `examples/` or `testdata/valid/`.
- Put rejected examples under `testdata/invalid/`.
- Only bump `schema_version` when old consumers must reject the manifest.
- Keep root schema and docs schema byte-for-byte equivalent after sync.

## GitHub Pages

- Pages deploys through the manual `Deploy Pages` GitHub Actions workflow.
- Do not add push-triggered Pages deploys unless explicitly requested.
- Do not introduce a separate Pages branch.

## Style

- Prefer small, focused diffs.
- Preserve existing names and formatting unless the task asks for a rename.
- Do not touch generated/synced files directly when a tool owns them.
