# gh-problemas

## Overview

A terminal UI for triaging and managing GitHub issues, built with
[Bubble Tea](https://github.com/charmbracelet/bubbletea) and
[Cobra](https://github.com/spf13/cobra).

## Structure

```text
cmd/            Cobra command definitions (root command, CLI entry point)
internal/
  config/       User configuration loading
  data/         GitHub GraphQL API clients (issues, comments, users, pagination)
  ui/           Bubble Tea application shell, styles, and key bindings
    components/ Reusable TUI components
    views/      Top-level views (dashboard, detail)
  utils/        Helpers (color, markdown rendering, time formatting)
tests/scrut/    End-to-end CLI tests (scrut snapshot tests)
docs/plans/     Implementation plans
```

## Development

```sh
make build      # Build the binary to bin/gh-problemas
make test       # Run unit tests
make lint       # Run golangci-lint
make vet        # Run go vet
make fmt        # Check gofmt compliance
make cover      # Run tests with coverage report
make tidy       # Tidy go.mod
make install    # Install as a gh extension
make help       # Show all available targets
```

Version is injected at build time via `-ldflags "-X main.version=..."`.

## Dependencies

- Go (version pinned in `go.mod`)
- `gh` CLI (for authentication context at runtime)
- `golangci-lint` (for `make lint`)
