# Gh Problemas

A terminal UI for triaging and managing GitHub issues.

## Installation

### Homebrew

```sh
brew install cboone/tap/gh-problemas
```

### From source

```sh
go install github.com/cboone/gh-problemas@latest
```

### From release

Download a binary from the [releases page](https://github.com/cboone/gh-problemas/releases).

### Build locally

```sh
git clone https://github.com/cboone/gh-problemas.git
cd gh-problemas
make build
./bin/gh-problemas
```

## Usage

```sh
gh-problemas
```

The app uses your `gh` authentication context. Run `gh auth login` if needed.

## License

[MIT License](./LICENSE). TL;DR: Do whatever you want with this software, just keep the copyright notice included. The authors aren't liable if something goes wrong.
