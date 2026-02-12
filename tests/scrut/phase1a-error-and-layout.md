# Phase 1A Error and Layout Tests

## Unknown flags produce helpful error

```scrut
$ gh-problemas --invalid-flag 2>&1 || true
Error: unknown flag: --invalid-flag
Usage:
  gh-problemas [flags]

Flags:
  -h, --help      help for gh-problemas
  -v, --version   version for gh-problemas

```

## Binary exists and is executable

```scrut
$ gh-problemas --version
gh-problemas version * (glob)
```
