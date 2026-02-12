# Phase 1A CLI Tests

## Help flag produces usage information

```scrut
$ gh-problemas --help
gh-problemas is a terminal user interface for triaging and managing GitHub issues.

Usage:
  gh-problemas [flags]

Flags:
  -h, --help      help for gh-problemas
  -v, --version   version for gh-problemas
```

## Version flag produces version string

```scrut
$ gh-problemas --version
gh-problemas version * (glob)
```

## Short help flag works

```scrut
$ gh-problemas -h
gh-problemas is a terminal user interface for triaging and managing GitHub issues.

Usage:
  gh-problemas [flags]

Flags:
  -h, --help      help for gh-problemas
  -v, --version   version for gh-problemas
```

## Short version flag works

```scrut
$ gh-problemas -v
gh-problemas version * (glob)
```
