# Repository Guidelines

## Project Structure & Module Organization

This repository contains two independent Go modules:

- `app/` is the working `postcli` reference. `cmd/postcli/main.go` wires dependencies, `internal/cli/` handles commands and output, and `internal/api/` performs HTTP and JSON operations. Tests live beside their packages as `*_test.go`.
- `myapp/` is the evolving application. Keep CLI routing in `controller/`, business workflows in `service/`, and external API access in `repository/`.
- `docker/` defines the optional Go/Codex development container.

Run Go commands from the module being changed; the root is not a Go module.

## Build, Test, and Development Commands

```bash
cd app && go run ./cmd/postcli help     # Run the reference CLI
cd app && go build ./cmd/postcli        # Compile the reference CLI
cd app && go test ./...                 # Run all app tests
cd app && go vet ./...                  # Run static checks
cd myapp && go test ./...               # Compile and test the evolving app
docker compose -f docker/docker-compose.yml up --build
```

Before submitting changes, run `go test ./...` and `go vet ./...` in every affected module.

## Coding Style & Naming Conventions

Format Go files with `gofmt`. Follow standard Go naming: exported identifiers use `PascalCase`, internal identifiers use `camelCase`, and packages use short lowercase names. Prefer constructors such as `NewClient`. Pass `context.Context` first to cancellable operations. Keep normal output on `stdout`, diagnostics on `stderr`, and wrap errors with `%w` when adding context.

Maintain the dependency direction `controller -> service -> repository`. Controllers select exit codes, services implement workflows, and repositories handle transport details.

## Testing Guidelines

Use Go's `testing` package and name tests `TestBehavior` in `_test.go` files. Prefer table-driven tests for command variants and `httptest.Server` for HTTP behavior. Tests must not depend on live JSONPlaceholder. There is no fixed coverage threshold; cover success, validation, API error, malformed response, and timeout paths when relevant.

## Commit & Pull Request Guidelines

History uses short, focused summaries in English or Japanese. Prefer an imperative subject such as `Add login command parsing`, and keep unrelated changes in separate commits. Pull requests should describe behavior changes, affected module(s), verification commands, and any new exit codes or CLI output. Link relevant issues; include terminal examples instead of screenshots for CLI changes.

## Security & Configuration

Never commit API tokens, passwords, or `.env` files. Supply environment-specific URLs and credentials through environment variables. Do not print credentials or tokens to either output stream.
