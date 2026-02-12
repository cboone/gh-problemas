# Phase 1B Pagination Tests

## Paginator sequential loading behavior

```scrut
$ go test ./internal/data -run TestPaginator_SequentialPages -v
=== RUN   TestPaginator_SequentialPages
--- PASS: TestPaginator_SequentialPages* (glob)
PASS
ok  *internal/data* (glob)
```

## Paginator exhaustion behavior

```scrut
$ go test ./internal/data -run TestPaginator_Exhausted -v
=== RUN   TestPaginator_Exhausted
--- PASS: TestPaginator_Exhausted* (glob)
PASS
ok  *internal/data* (glob)
```
