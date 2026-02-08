# gh-problemas

A GitHub CLI extension for managing problemas.

## Installation

### From source

```bash
gh extension install cboone/gh-problemas
```

### Local development

```bash
# Clone the repository
git clone https://github.com/cboone/gh-problemas.git
cd gh-problemas

# Build the binary
make build

# Install locally
make install
```

## Usage

```bash
# Show help
gh problemas help

# Show version
gh problemas version
```

## Building

This extension is written in Go and compiled to a binary.

```bash
# Build the binary
make build

# Run tests
make test

# Format code
make fmt

# Run linter
make lint
```

## Requirements

- Go 1.24 or higher
- GitHub CLI (`gh`)

## License

See [LICENSE](LICENSE) file.