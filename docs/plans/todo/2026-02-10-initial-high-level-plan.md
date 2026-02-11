# gh-problemas — High-Level Plan

A `gh` CLI extension that provides a rich terminal user interface (TUI) for managing GitHub issues. Written in Go.

```
gh problemas
```

---

## Vision

**gh-problemas** fills a gap in the `gh` extension ecosystem: a dedicated, keyboard-driven TUI for *issue management*. While `gh-dash` is excellent for PR dashboards with some issue visibility, there is no focused tool for issue triage, exploration, and lifecycle management from the terminal. gh-problemas is that tool.

---

## Tech Stack

| Layer | Choice | Rationale |
|-------|--------|-----------|
| Language | **Go** | Static binaries, no runtime deps, gh ecosystem standard |
| TUI framework | **Bubble Tea** (charmbracelet/bubbletea) | Elm architecture, inline + fullscreen modes, dominant in gh extensions |
| TUI components | **Bubbles** (charmbracelet/bubbles) | List, table, viewport, text input, spinner — battle-tested |
| Styling | **Lip Gloss** (charmbracelet/lipgloss) | CSS-like terminal styling, composable |
| Markdown | **Glamour** (charmbracelet/glamour) | Render issue bodies as styled markdown in the terminal |
| GitHub API | **go-gh v2** (cli/go-gh) | Inherits `gh` auth, REST + GraphQL clients, repo context |
| CLI framework | **Cobra** | Subcommands, flags, help generation |
| Config | **Viper** + YAML | User-configurable dashboards, keybindings, themes |
| CI/CD | **gh-extension-precompile** action | Automated cross-platform releases, purpose-built for gh extensions |

---

## Core Features

### 1. Issue Dashboard (default view)

The main screen users see when running `gh problemas`.

- **Configurable sections** — each section is a filtered view of issues (e.g. "My Issues", "Needs Triage", "High Priority", "Stale")
- **Column display** — number, title, author, labels, milestone, assignees, age, comment count, reactions
- **Real-time sorting** — by created, updated, comments, reactions, or priority
- **Section switching** — tab/shift-tab between sections
- **Auto-refresh** — configurable polling interval
- **Multi-repo support** — view issues across multiple repositories in one dashboard

### 2. Issue Detail View

Drill into a single issue from the dashboard.

- **Rendered markdown body** via Glamour
- **Full comment thread** with author, timestamp, and reactions
- **Metadata sidebar** — labels, milestone, assignees, project, linked PRs
- **Timeline events** — label changes, assignments, references, cross-links
- **Inline actions** — comment, close/reopen, assign, label, edit (opens `$EDITOR`)

### 3. Issue Creation

Create new issues without leaving the terminal.

- **Template selection** — list and pick from repo issue templates
- **Form mode** — fill in title, body, labels, milestone, assignees interactively
- **Editor handoff** — open `$EDITOR` for composing longer bodies
- **Preview** — render the markdown before submitting

### 4. Bulk Operations

Act on multiple issues at once.

- **Multi-select** — spacebar to toggle, visual indicators
- **Bulk label** — add/remove labels to selection
- **Bulk assign** — assign/unassign users
- **Bulk close/reopen** — with optional comment
- **Bulk milestone** — move issues to a milestone

### 5. Search & Filter

Powerful filtering beyond what sections provide.

- **Fuzzy search** — across title, body, labels, author
- **Filter bar** — composable filters: `is:open label:bug assignee:@me milestone:v2`
- **Saved filters** — persist named filters in config
- **Full-text search** — delegates to GitHub's search API for server-side search

### 6. Labels & Milestones Management

First-class views for organizing primitives.

- **Label browser** — view all labels with colors, create/edit/delete
- **Milestone browser** — view milestones with progress bars (open/closed issue counts)
- **Filter-from-label** — jump to a filtered issue list for any label or milestone

### 7. Keyboard-Driven Workflow

- **Vim-style navigation** — `j`/`k`, `g`/`G`, `/` for search
- **Command palette** — `:` or `ctrl+p` to invoke any action by name
- **Customizable keybindings** — override any binding in config
- **Context-sensitive help** — `?` shows available actions for the current view

---

## Architecture

```
gh-problemas/
├── main.go                  # Entry point, Cobra root command
├── cmd/                     # CLI command definitions
│   └── root.go              # Root command, flag parsing, config loading
├── internal/
│   ├── config/              # YAML config parsing, defaults, validation
│   │   ├── config.go
│   │   └── keys.go          # Keybinding definitions
│   ├── data/                # GitHub API data layer
│   │   ├── issues.go        # Issue queries (list, get, search)
│   │   ├── labels.go        # Label CRUD
│   │   ├── milestones.go    # Milestone queries
│   │   ├── comments.go      # Comment queries and mutations
│   │   ├── pagination.go    # Cursor-based pagination helpers
│   │   └── models.go        # Data models (split by domain as they grow)
│   ├── ui/                  # Bubble Tea TUI layer
│   │   ├── app.go           # Top-level app model, view stack routing
│   │   ├── keys.go          # Keymap definitions
│   │   ├── styles.go        # Lip Gloss style definitions
│   │   ├── theme.go         # Theme support (light/dark/custom)
│   │   ├── components/      # Reusable UI components
│   │   │   ├── section.go   # A single dashboard section (list of issues)
│   │   │   ├── table.go     # Issue table component
│   │   │   ├── filter.go    # Filter bar component
│   │   │   ├── statusbar.go # Bottom status bar
│   │   │   ├── header.go    # Top header bar
│   │   │   ├── spinner.go   # Loading indicators for async operations
│   │   │   ├── help.go      # Help overlay
│   │   │   └── prompt.go    # Confirmation / input prompts
│   │   └── views/           # Full-screen view models
│   │       ├── dashboard.go # Main dashboard (sections + issue lists)
│   │       ├── detail.go    # Issue detail view
│   │       ├── create.go    # Issue creation form
│   │       ├── labels.go    # Label management view
│   │       └── milestones.go# Milestone management view
│   └── utils/               # Shared utilities
│       ├── markdown.go      # Glamour markdown rendering helpers
│       ├── time.go          # Relative time formatting
│       └── color.go         # Label color conversion (hex → ANSI)
├── config.example.yml       # Example configuration file
├── .github/
│   └── workflows/
│       └── release.yml      # Cross-compile + release on tag push
├── go.mod
├── go.sum
├── LICENSE
├── README.md
└── PLAN.md
```

### Key architectural decisions

1. **`internal/` package** — all non-`main` code is unexported. This is idiomatic Go for applications (vs libraries).

2. **Data layer separated from UI** — `internal/data/` handles all GitHub API calls and returns plain Go structs. The UI layer never makes API calls directly. This enables testing the data layer independently and swapping API strategies (REST vs GraphQL) without touching UI code.

3. **GraphQL-first for reads, REST for mutations** — GraphQL lets us fetch exactly the fields we need in a single request (issues + labels + milestones + comments). REST is simpler for mutations (close issue, add label) where we're sending a known payload.

4. **View stack routing in `app.go`** — the top-level Bubble Tea model maintains a view stack for navigation history. It delegates `Update` and `View` calls to the top of the stack. Views communicate upward via messages (e.g. `NavigateToDetailMsg{IssueNumber: 42}` to push, `NavigateBackMsg` to pop). This enables natural back-navigation with `esc`/`backspace`.

5. **Config-driven sections** — dashboard sections are defined in YAML, not hardcoded. Each section specifies a title, query filters, sort order, and display limit. This follows the pattern established by `gh-dash`.

6. **Pagination abstraction** — cursor-based pagination is handled by a shared helper in the data layer, wrapping GraphQL `pageInfo`/`endCursor` patterns. This is needed from Phase 1 since issue lists can be large.

7. **`@me` alias resolution** — the config uses `@me` as a shorthand for the authenticated user. The data layer resolves this to the actual GitHub username via `go-gh` auth context before making API calls.

8. **Loading states** — all async operations (API fetches, mutations) display loading indicators. The dashboard shows a spinner during initial fetch and inline loading states during refresh. This prevents the TUI from appearing frozen.

---

## Configuration

```yaml
# ~/.config/gh-problemas/config.yml

version: 1                    # config schema version (for future migrations)

defaults:
  repo: ""                    # default repo (empty = current repo from git context)
  refresh_interval: 300       # seconds between auto-refresh (0 = disabled)
  page_size: 50               # issues per API request
  date_format: relative       # "relative" or "2006-01-02"

theme: dark                   # "dark", "light", or path to custom theme file

sections:
  - title: My Issues
    filters:
      assignee: "@me"
      state: open
    sort: updated
    limit: 20

  - title: Needs Triage
    filters:
      "no:label": true
      state: open
    sort: created
    limit: 20

  - title: Bugs
    filters:
      label: bug
      state: open
    sort: reactions
    limit: 20

  - title: Recently Closed
    filters:
      state: closed
    sort: updated
    limit: 10

keybindings:
  global:
    quit: ["q", "ctrl+c"]
    help: ["?"]
    refresh: ["R"]
    search: ["/"]
    command_palette: ["ctrl+p"]
  dashboard:
    next_section: ["tab"]
    prev_section: ["shift+tab"]
    open_issue: ["enter"]
    new_issue: ["c"]
    toggle_select: ["space"]
    bulk_label: ["L"]
    bulk_assign: ["A"]
    bulk_close: ["X"]
  detail:
    back: ["esc", "backspace"]
    comment: ["c"]
    close_reopen: ["x"]
    assign: ["a"]
    label: ["l"]
    edit: ["e"]
    open_browser: ["o"]
```

---

## Implementation Phases

### Phase 1 — Scaffold & Read-Only Dashboard

Get a working TUI that displays issues from the current repo.

- [ ] Initialize Go module, install dependencies
- [ ] Cobra root command with `--help` and `--version` flags
- [ ] Implement `internal/data/` — fetch issues via GraphQL, parse into models
- [ ] Cursor-based pagination helper for GraphQL queries
- [ ] Resolve `@me` alias to authenticated username via `go-gh`
- [ ] Implement basic Bubble Tea app shell with dashboard view
- [ ] Render issues in a table with key columns (number, title, labels, author, age)
- [ ] Loading spinner during initial data fetch
- [ ] Navigation: j/k scrolling, enter to open issue detail (read-only)
- [ ] Issue detail view: rendered markdown body, comments, metadata
- [ ] View stack for back-navigation (detail -> dashboard)
- [ ] Status bar with repo name, section info, keybinding hints
- [ ] Load config from `~/.config/gh-problemas/config.yml` with sensible defaults
- [ ] Basic error display in status bar (API failures, network errors)

### Phase 2 — Sections, Filtering & Search

Make the dashboard powerful and configurable.

- [ ] Config-driven sections with independent filters and sort orders
- [ ] Tab/shift-tab section switching
- [ ] Filter bar component with live filtering
- [ ] Fuzzy search across visible issues (client-side)
- [ ] Server-side search via GitHub search API
- [ ] Saved filters in config

### Phase 3 — Issue Mutations

Enable writing, not just reading.

- [ ] Close / reopen issues
- [ ] Add / remove labels (interactive picker)
- [ ] Assign / unassign users (interactive picker)
- [ ] Set milestone
- [ ] Add comments (inline text area or `$EDITOR` handoff)
- [ ] Create new issues (template selection, form mode, editor handoff, preview)

### Phase 4 — Bulk Operations & Multi-Repo

Scale up to power-user workflows.

- [ ] Multi-select with spacebar
- [ ] Bulk label, assign, close, milestone (concurrent worker pool with rate limit awareness)
- [ ] Progress indicator for bulk operations
- [ ] Multi-repo support — configure multiple repos, unified view
- [ ] Cross-repo search

### Phase 5 — Polish & Release

Production readiness.

- [ ] Customizable keybindings (loaded from config)
- [ ] Theme support (dark/light/custom)
- [ ] Command palette
- [ ] Help overlay (`?`)
- [ ] Auto-refresh with configurable interval and ETag-based conditional requests
- [ ] Graceful rate limit handling (backoff, status bar indicator)
- [ ] `NO_COLOR` and limited terminal support
- [ ] GitHub Actions release workflow with `gh-extension-precompile`
- [ ] README with screenshots, installation instructions, configuration docs
- [ ] Add `gh-extension` topic to repository

---

## API Strategy

### Primary queries (GraphQL)

GraphQL is ideal for the dashboard because we can fetch exactly what we need:

```graphql
query ListIssues($owner: String!, $name: String!, $first: Int!, $after: String, $states: [IssueState!], $labels: [String!], $orderBy: IssueOrder!) {
  repository(owner: $owner, name: $name) {
    issues(first: $first, after: $after, states: $states, labels: $labels, orderBy: $orderBy) {
      pageInfo { hasNextPage endCursor }
      nodes {
        number
        title
        state
        createdAt
        updatedAt
        author { login }
        labels(first: 10) { nodes { name color } }
        assignees(first: 5) { nodes { login } }
        milestone { title }
        comments { totalCount }
        reactions { totalCount }
      }
    }
  }
}
```

### Mutations (REST)

REST is simpler for single-resource mutations:

```
PATCH /repos/{owner}/{repo}/issues/{number}          — close/reopen, edit
POST  /repos/{owner}/{repo}/issues                   — create
POST  /repos/{owner}/{repo}/issues/{number}/labels    — add labels
POST  /repos/{owner}/{repo}/issues/{number}/assignees  — assign
POST  /repos/{owner}/{repo}/issues/{number}/comments   — comment
```

### Rate limit awareness

- Track `X-RateLimit-Remaining` and `X-RateLimit-Reset` headers
- Display rate limit info in status bar
- Back off gracefully when approaching limits
- Use conditional requests (`If-None-Match` / ETags) for refresh polling

---

## Design Principles

1. **Speed over completeness** — show something fast, load details lazily. Paginate aggressively. Cache locally.

2. **Keyboard-first, discoverable** — every action is reachable by keyboard. Status bar always shows relevant bindings. `?` reveals all options.

3. **Config-driven, zero-config** — works out of the box with sensible defaults for the current repo. Power users customize everything via YAML.

4. **Respect the terminal** — honor `NO_COLOR`, `TERM`, dark/light backgrounds. Degrade gracefully in limited terminals. Don't break pipes.

5. **gh-native** — use `go-gh` for auth, repo context, and API access. Never ask users to configure tokens. Behave like a natural extension of `gh`.
