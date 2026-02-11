# Phase 1B Network and API Error Tests

## Status bar classifies network errors

```scrut
$ go test ./internal/ui/components -run TestStatusBar_SetError_Network -v
=== RUN   TestStatusBar_SetError_Network
--- PASS: TestStatusBar_SetError_Network* (glob)
PASS
ok  *internal/ui/components* (glob)
```

## Status bar provides API guidance for auth and missing repos

```scrut
$ go test ./internal/ui/components -run 'TestStatusBar_SetError_API401|TestStatusBar_SetError_API404' -v
=== RUN   TestStatusBar_SetError_API401
--- PASS: TestStatusBar_SetError_API401* (glob)
=== RUN   TestStatusBar_SetError_API404
--- PASS: TestStatusBar_SetError_API404* (glob)
PASS
ok  *internal/ui/components* (glob)
```
