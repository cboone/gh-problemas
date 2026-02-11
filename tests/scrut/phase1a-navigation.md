# Phase 1A Navigation Tests

These tests verify the application handles missing context gracefully.

## Running outside a git repo shows helpful error

```scrut
$ cd /tmp && gh-problemas 2>&1 || true
Error: could not determine repository: * (glob)
*error*exit status 128* (glob)
*Run this command from inside a git repository with a GitHub remote.* (glob)
Usage:
  gh-problemas [flags]

Flags:
  -h, --help      help for gh-problemas
  -v, --version   version for gh-problemas

```
