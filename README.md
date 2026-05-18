# Blockforge Manifest

`blockforge-manifest` defines a production JSON Schema for Blockforge-compatible Minecraft server manifests.

The schema is modpack-agnostic. It can describe any modpack that publishes a supported loader installer, server configuration ZIP, and downloadable mod JAR list.

## Contract

- Root `manifest.schema.json` is canonical.
- `docs/manifest.schema.json` is published GitHub Pages copy.
- Public schema URL: `https://rannday.github.io/blockforge-manifest/manifest.schema.json`
- Current schema version: `1`
- `pack` is optional metadata. It is not required and no pack name is reserved.
- `loader.type` supports `neoforge`, `forge`, `fabric`, and `quilt`.

This schema describes a full production manifest.

Required production fields:

- `minecraft`
- `server_config`
- `loader.type`
- `loader.sha1`

All downloadable artifacts use HTTP(S) URLs.

- `loader.installer_url` ends in `.jar`
- `server_config.url` ends in `.zip`
- `mods[].url` ends in `.jar`
- `file:` URLs rejected on purpose

`server_config` is required. `mods[]` entries are full production objects with `name`, `url`, `website_url`, `sha1`, and `size`.

## Workflow

Before committing schema changes, run:

```powershell
go run ./cmd/sync-docs
go test ./...
```

`go test ./...` is the normal pass/fail check. It asserts valid fixtures pass and invalid fixtures fail. It also verifies Go tooling still compiles.

Go module dependencies are managed with `go.mod` and `go.sum`.

Manual validator checks are optional. `go test ./...` already covers these cases.

Useful CLI examples:

```powershell
go run ./cmd/validate-manifest -- examples/manifest.example.json testdata/valid/manifest.v1.json
go run ./cmd/validate-manifest -- testdata/invalid/*.json
```

The invalid fixture command exits with status `1`.

## Scripts

### `cmd/sync-docs`

Copies canonical root schema into Pages directory and ensures `docs/.nojekyll` exists as an empty GitHub Pages marker file.

It updates:

```text
docs/manifest.schema.json
docs/.nojekyll
```

Run:

```powershell
go run ./cmd/sync-docs
```

### `cmd/validate-manifest`

Validates one or more manifest JSON files against root schema. Globs are expanded by the tool, so the command works in PowerShell and shells that do not expand globs.

Run:

```powershell
go run ./cmd/validate-manifest -- examples/manifest.example.json testdata/valid/manifest.v1.json
go run ./cmd/validate-manifest -- testdata/invalid/*.json
```

Valid files print `OK: <path>`. Invalid files print `INVALID: <path>` plus readable errors.

## GitHub Pages

This repo publishes Pages from GitHub Actions instead of deploying from a branch.

One-time repository setup:

1. Open repository Settings.
2. Open Pages.
3. Set Source to `GitHub Actions`.
4. Save.

Publishing flow:

1. Commit schema and docs changes to the normal branch.
2. Open Actions.
3. Select `Deploy Pages`.
4. Click `Run workflow`.

The workflow is manual-only through `workflow_dispatch`, so Pages does not update on every push. It runs `go test ./...`, runs `go run ./cmd/sync-docs`, uploads the `docs` directory as a Pages artifact, then deploys it with `actions/deploy-pages`.

No separate deployment branch is required. The published content comes from `docs/` in the workflow checkout.

Pushes and pull requests still run CI through `.github/workflows/ci.yml`; only the manual Pages workflow publishes the site.

## Consumer Guidance

### Manifest Producers

- Validate generated manifest against `manifest.schema.json` before publish.
- Fetch published schema URL during producer or release validation if latest published schema needed.
- Do not publish if validation fails.

### Manifest Consumers

- Do not fetch schema at runtime.
- Keep installer self-contained.
- Installer-specific checks still live in Go code:
  - safe inferred `.jar` filenames
  - ZIP traversal protection
  - applying manifest to filesystem

## Compatibility Policy

- Only bump `schema_version` when old consumers must reject manifest.
- Additive optional fields need care and tests.
- Do not change meaning of required fields without schema version bump.
- Tag this repo when releasing schema changes.
