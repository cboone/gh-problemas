# Phase 1 Implementation Plan: Scaffold through Read-Only MVP

## Context

The gh-problemas repository currently has no code — only a high-level plan, LICENSE, gitignore, and a stub README. This plan covers Phase 1A (scaffold + read-only MVP dashboard) and Phase 1B (pagination, comments, config baseline), broken into 8 sequential milestones that each produce testable, working code.

The high-level plan lives at `docs/plans/todo/2026-02-10-initial-high-level-plan.md` and is the source of truth for scope and architecture decisions.

**Goal:** A user can `gh problemas` in a repo, see open issues, drill into an issue detail with rendered markdown, scroll comments, and return to the dashboard — all keyboard-driven, read-only, zero-config.

---

## Milestone 1 — Project Scaffold and CLI Shell (Phase 1A)

**Goal:** Runnable `gh problemas` binary with `--help` and `--version`. No API calls yet.

### Files to create

**`go.mod`** — `go mod init github.com/hpg/gh-problemas` then `go get`:
- `github.com/cli/go-gh/v2`
- `github.com/charmbracelet/bubbletea`
- `github.com/charmbracelet/bubbles`
- `github.com/charmbracelet/lipgloss`
- `github.com/charmbracelet/glamour`
- `github.com/spf13/cobra`
- `github.com/spf13/viper`

**`main.go`** — Minimal entry point calling `cmd.Execute()`.

**`cmd/root.go`** — Cobra root command:
- `Use: "gh-problemas"`, version set via `ldflags`
- `RunE` will eventually launch Bubble Tea; for M1, prints a placeholder
- Handles `--help`, `--version` automatically via Cobra

**`Makefile`** — Build, test, lint, install targets with `VERSION` and `LDFLAGS`.

### Verification
- `go build ./...` succeeds
- `go run . --help` and `go run . --version` produce expected output

---

## Milestone 2 — Data Models and Issue Listing (Phase 1A)

**Goal:** `internal/data/` can fetch issues from GitHub GraphQL API for the current repo.

### Files to create

**`internal/data/models.go`** — Core value types:
- `Issue` — Number, Title, State, CreatedAt, UpdatedAt, Author, Labels, Assignees, Milestone, CommentCount, ReactionCount, Body
- `Label` — Name, Color (hex)
- `PageInfo` — HasNextPage, EndCursor
- `IssueListResult` — Issues, PageInfo
- `IssueListOptions` — States, Labels, OrderBy, First, After

**`internal/data/issues.go`** — Issue data client:
- `Querier` interface: `Do(query string, variables map[string]interface{}, response interface{}) error` — abstracts `*api.GraphQLClient` for testability
- `IssueClient` struct accepting `Querier`, owner, repo
- `List(opts IssueListOptions) (IssueListResult, error)` — uses the GraphQL `ListIssues` query from the high-level plan
- `Get(number int) (Issue, error)` — fetches single issue with body via separate GraphQL query
- Internal response structs mirroring GraphQL JSON shape, with `toIssue()` converter

### Tests (`internal/data/issues_test.go`)
- Mock `Querier` returning canned JSON; test List with 3 issues, empty response, and GraphQL error
- Test Get with full body, and not-found error

---

## Milestone 3 — Utility Helpers (Phase 1A)

**Goal:** Reusable formatting utilities needed by the UI layer.

### Files to create

**`internal/utils/time.go`**
- `RelativeTime(t time.Time) string` — "2m ago", "3h ago", "5d ago", "2mo ago", "1y ago"

**`internal/utils/color.go`**
- `HexToColor(hex string) lipgloss.Color` — handles with/without `#`
- `ContrastColor(backgroundHex string) lipgloss.Color` — W3C luminance algorithm, returns black or white

**`internal/utils/markdown.go`**
- `RenderMarkdown(content string, width int) (string, error)` — Glamour with `WithAutoStyle()` and `WithWordWrap(width)`

### Tests
- `time_test.go`: table-driven with 30s/90m/25h/60d/400d cases, edge cases for future time and zero
- `color_test.go`: hex parsing, contrast against white/black/saturated red backgrounds
- `markdown_test.go`: simple render, empty string, width wrapping

---

## Milestone 4 — App Shell, Styles, Keys, Components (Phase 1A)

**Goal:** Bubble Tea app skeleton with view stack routing, styles, key bindings, status bar, and spinner.

### Files to create

**`internal/ui/styles.go`** — `Styles` struct with `DefaultStyles()`: App, Header, StatusBar, SelectedRow, NormalRow, IssueNumber, IssueTitle, LabelStyle, Spinner, ErrorText, HelpKey, HelpDesc.

**`internal/ui/keys.go`** — `KeyMap` struct with `DefaultKeyMap()`: Up (k/up), Down (j/down), Open (enter), Back (esc/backspace), Quit (q), ForceQuit (ctrl+c), Refresh (R), Help (?), PageUp, PageDown, GoToTop (g), GoToBottom (G), NextPage (L).

**`internal/ui/components/statusbar.go`** — Status bar showing repo name, key hints, and transient messages (error/success). Methods: `SetMessage()`, `SetKeyHints()`.

**`internal/ui/components/spinner.go`** — Thin wrapper around `bubbles/spinner` with `Start(label)`, `Stop()`, `IsActive()`.

**`internal/ui/app.go`** — Top-level Bubble Tea model. This is the core architectural piece:

```
View interface {
    Init() tea.Cmd
    Update(msg tea.Msg) (View, tea.Cmd)
    View() string
    KeyHints() []string
}
```

- `App` struct holds: `viewStack []View`, `statusBar`, `keys`, `styles`, `issueClient`, dimensions
- Navigation messages: `NavigateToDetailMsg{IssueNumber}`, `NavigateBackMsg{}`
- Data messages: `IssuesLoadedMsg{Result, Err}`, `IssueDetailLoadedMsg{Issue, Err}`
- `Update` handles: `tea.WindowSizeMsg` (propagate), `ForceQuit` (always quit), `Quit` on last view, `NavigateToDetailMsg` (push), `NavigateBackMsg` (pop), then delegates to current view
- `View` renders current view + status bar, with status bar taking 1 line at bottom
- Views use pointer receivers; stack stores `*DashboardView`, `*DetailView`

**Update `cmd/root.go`** — Wire up the real app:
1. `gh.CurrentRepository()` for owner/name
2. `api.DefaultGraphQLClient()` for API client
3. Create `IssueClient`, then `App`
4. `tea.NewProgram(app, tea.WithAltScreen()).Run()`

### Tests (`internal/ui/app_test.go`)
- Push/pop view stack operations
- `NavigateToDetailMsg` pushes, `NavigateBackMsg` pops (when stack > 1)
- `ForceQuit` produces `tea.Quit`; `q` on last view quits; `q` on detail delegates (no quit)

---

## Milestone 5 — Dashboard View (Phase 1A)

**Goal:** Main screen showing open issues with keyboard navigation, loading state, and error state.

### Files to create

**`internal/ui/views/dashboard.go`**:
- `issueItem` wrapping `data.Issue`, implementing `list.Item` with `FilterValue()`, `Title()`, `Description()`
- Custom `issueDelegate` for rendering: issue number, title, labels (colored), author, age, comment count
- `DashboardView` struct with `bubbles/list.Model`, `IssueClient`, spinner, loading/error state
- `Init()` returns a `tea.Cmd` that calls `issueClient.List()` and returns `IssuesLoadedMsg`
- `Update()`: on `IssuesLoadedMsg` populates list; on `Enter` emits `NavigateToDetailMsg`; on `R` triggers refresh; delegates to `list.Model` for j/k navigation
- `View()`: centered spinner during initial load, error message on failure, list view otherwise

### Tests
- Create with mock client, send `IssuesLoadedMsg`, verify list item count
- Enter on selected item produces `NavigateToDetailMsg`
- `R` triggers fetch command

---

## Milestone 6 — Detail View + Scrut Baseline (Phase 1A)

**Goal:** Issue detail with rendered markdown body in scrollable viewport. Phase 1A acceptance criteria met.

### Files to create

**`internal/ui/views/detail.go`**:
- `DetailView` with `bubbles/viewport.Model`, issue data, spinner, loading/error state
- `Init()` returns command to call `issueClient.Get(number)`, returns `IssueDetailLoadedMsg`
- On load: render header (title, number, metadata line with state/author/dates/milestone/assignees/labels) + markdown body via `utils.RenderMarkdown`, set as viewport content
- `esc`/`backspace` emits `NavigateBackMsg`; `q` in detail also emits `NavigateBackMsg` (not app quit)
- Viewport handles j/k/pgup/pgdown scrolling natively

**`tests/scrut/phase1a-cli.md`** (Scrut):
- `gh-problemas --help` produces expected help text
- `gh-problemas --version` outputs version string

**`tests/scrut/phase1a-navigation.md`** (Scrut):
- Boot application in a fixture repo with deterministic mocked issue data
- Verify dashboard list renders with loading -> populated state
- Open selected issue, verify markdown body render, navigate back, then quit
- Confirm returning to dashboard preserves prior list position/selection

**`tests/scrut/phase1a-error-and-layout.md`** (Scrut):
- Unauthenticated (`401`) path shows `gh auth login` guidance and non-crashing UI
- Missing repo (`404`) path shows user-friendly error in main view + status bar
- Resize / narrow terminal snapshot coverage for dashboard and detail layout
- Empty issue list snapshot coverage (no-data state)

### Phase 1A Definition of Done
- [ ] User runs `gh problemas`, sees open issues in current repo
- [ ] User opens issue detail, returns to dashboard without losing list state (view stack)
- [ ] `go test ./...` passes for all packages
- [ ] Scrut covers help/version, boot -> list -> detail -> back -> quit, error states, and layout snapshots

---

## Milestone 7 — Pagination, Comments, User Resolution (Phase 1B)

**Goal:** Handle large repos, show comments in detail, resolve `@me` alias.

### Files to create

**`internal/data/pagination.go`**:
- `Paginator` struct: manages cursor state, page size, total loaded count
- `NextPageRequest() *PageRequest` — returns nil when no more pages
- `Update(pageInfo PageInfo, count int)` — records page result
- `HasNextPage() bool`, `TotalLoaded() int`, `Reset()`

**`internal/data/comments.go`**:
- `Comment` struct: Author, Body, CreatedAt, UpdatedAt, Reactions
- `CommentClient` accepting `Querier`, owner, repo
- `List(issueNumber, first int, after string) (CommentListResult, error)` — GraphQL query for issue comments

**`internal/data/user.go`**:
- `UserClient` with `WhoAmI() (string, error)` — `query { viewer { login } }` via GraphQL

### Files to update

**`internal/ui/views/dashboard.go`**:
- Add `Paginator` to `DashboardView`
- `L` key loads next page (appends to list); status shows "Showing X issues"
- New message type `IssuesPageLoadedMsg` with `Append bool`

**`internal/ui/views/detail.go`**:
- After issue loads, fetch comments via `CommentClient`
- New `CommentsLoadedMsg`; render comments below body with author, timestamp, rendered markdown
- Optional "Load more comments" for long threads

**`internal/ui/components/statusbar.go`**:
- Repo context (`owner/repo`) on left
- Key hints center
- Error/loading on right; truncate long errors
- Distinct network failure vs API error display

### Tests
- `pagination_test.go`: initial state, sequential page requests, exhausted pages, reset
- `comments_test.go`: mock 3 comments, empty list, error propagation
- `user_test.go`: mock viewer response, error propagation
- Scrut `tests/scrut/phase1b-pagination.md`: load-more flow, exhausted cursor behavior, and status text updates
- Scrut `tests/scrut/phase1b-comments.md`: comment thread rendering with markdown + long thread viewport behavior
- Scrut `tests/scrut/phase1b-network-errors.md`: network/API failure during pagination or comments fetch with retry messaging

---

## Milestone 8 — Configuration Layer (Phase 1B)

**Goal:** Config loading from `~/.config/gh-problemas/config.yml` with sensible defaults. Phase 1B acceptance complete.

### Files to create

**`internal/config/config.go`**:
- `Config` struct: Version, Defaults (Repo, RefreshInterval, PageSize, DateFormat), Theme
- `Load() (*Config, error)` using Viper:
  - Set defaults (version=1, refresh=300, page size=50, date format=relative, theme=dark)
  - Config path: `XDG_CONFIG_HOME/gh-problemas` or `~/.config/gh-problemas`
  - Read YAML; ignore `ConfigFileNotFoundError` (zero-config works)
  - Unmarshal into `Config`

### Files to update

**`cmd/root.go`** — Load config before creating App; pass to `NewApp`.

**`internal/ui/views/dashboard.go`** — Use `cfg.Defaults.PageSize` instead of hardcoded 50.

**`internal/ui/views/detail.go`** — Use `cfg.Defaults.DateFormat` for timestamp display.

### Tests (`internal/config/config_test.go`)
- No config file: all defaults applied
- Partial override: specified values used, others default
- Invalid YAML: error returned
- XDG_CONFIG_HOME respected; fallback to `~/.config/gh-problemas`
- Scrut `tests/scrut/phase1b-config.md`: startup with no config, then startup with custom `page_size` and `date_format` to verify visible behavior

### Phase 1B Definition of Done
- [ ] Large issue lists handled via cursor-based pagination with "load more"
- [ ] Detail view shows comments and metadata
- [ ] Zero-config works (no config file needed, all defaults applied)
- [ ] `go test ./...` passes for all packages
- [ ] Scrut covers pagination, comments rendering, config defaults/overrides, and network error states

---

## Cross-Cutting: Error Handling

All API errors flow through message types (`Err` field). Pattern:
1. `tea.Cmd` catches errors, wraps in result message
2. View `Update` checks `Err`, sets view-level error state
3. View `View` renders error with retry hint
4. Status bar displays most recent error

HTTP status codes get user-friendly messages:
- 401: "Run `gh auth login` to re-authenticate"
- 403: "Check your permissions for this repository"
- 404: "Repository not found"

## Cross-Cutting: Testability

Every data client accepts the `Querier` interface, not `*api.GraphQLClient` directly. Tests use a `mockQuerier` that returns canned responses via JSON marshal/unmarshal round-trip.

## Cross-Cutting: Scrut Coverage Policy

Phase 1 must include both smoke and regression Scrut scenarios; do not treat Scrut as CLI-only checks.

- Smoke: CLI help/version and dashboard -> detail -> back -> quit loop
- Error: 401/403/404 and network failure UX for list/detail/pagination paths
- Layout: narrow-width and resize snapshots for dashboard/detail
- Data shape: empty list, long markdown body, and long comment thread rendering
- Config: no-config defaults and visible behavior from config overrides

## Verification

After all milestones:
1. `go build ./...` — compiles cleanly
2. `go test ./...` — all tests pass
3. `go vet ./...` — no issues
4. Manual: `go run . --help`, `go run . --version`
5. Manual in a real repo with `gh` auth: `go run .` shows issues, enter opens detail, esc returns, q quits
6. Scrut smoke: `scrut tests/scrut/phase1a-*.md` passes
7. Scrut Phase 1B regression: `scrut tests/scrut/phase1b-*.md` passes
8. Full Scrut suite: `scrut tests/scrut/` passes
