# Phase 1B Comments and Detail Tests

## Comments client returns comment list

```scrut
$ go test ./internal/data -run TestCommentList_ThreeComments -v
=== RUN   TestCommentList_ThreeComments
--- PASS: TestCommentList_ThreeComments* (glob)
PASS
ok  *internal/data* (glob)
```

## Detail view uses configured date formatting

```scrut
$ go test ./internal/ui/views -run TestDetailView_UsesConfiguredDateFormat -v
=== RUN   TestDetailView_UsesConfiguredDateFormat
--- PASS: TestDetailView_UsesConfiguredDateFormat* (glob)
PASS
ok  *internal/ui/views* (glob)
```
