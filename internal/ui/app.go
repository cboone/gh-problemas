package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hpg/gh-problemas/internal/data"
	"github.com/hpg/gh-problemas/internal/ui/components"
)

// View is the interface for a screen in the view stack.
type View interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (View, tea.Cmd)
	View() string
	KeyHints() []string
}

// Navigation messages

// NavigateToDetailMsg requests navigation to an issue detail view.
type NavigateToDetailMsg struct {
	IssueNumber int
}

// NavigateBackMsg requests navigation back to the previous view.
type NavigateBackMsg struct{}

// Data messages

// IssuesLoadedMsg carries the result of loading issues.
type IssuesLoadedMsg struct {
	Result data.IssueListResult
	Err    error
}

// IssuesPageLoadedMsg carries the result of loading an additional page of issues.
type IssuesPageLoadedMsg struct {
	Result data.IssueListResult
	Append bool
	Err    error
}

// IssueDetailLoadedMsg carries the result of loading a single issue.
type IssueDetailLoadedMsg struct {
	Issue data.Issue
	Err   error
}

// CommentsLoadedMsg carries the result of loading comments.
type CommentsLoadedMsg struct {
	Comments []data.Comment
	PageInfo data.PageInfo
	Err      error
}

// ViewFactory creates the initial view to push onto the stack.
type ViewFactory func(app *App) View

// DetailViewFactory creates a detail view for a given issue number.
type DetailViewFactory func(app *App, issueNumber int) View

// App is the top-level Bubble Tea model.
type App struct {
	viewStack       []View
	statusBar       *components.StatusBar
	keys            KeyMap
	styles          Styles
	issueClient     *data.IssueClient
	width           int
	height          int
	repoName        string
	initView        ViewFactory
	detailViewFn    DetailViewFactory
}

// NewApp creates a new App with the given issue client, repo name, and view factories.
func NewApp(client *data.IssueClient, repoName string, initView ViewFactory, detailView ...DetailViewFactory) *App {
	styles := DefaultStyles()
	keys := DefaultKeyMap()
	sb := components.NewStatusBar(styles.StatusBar)
	sb.SetRepoName(repoName)

	var dvFn DetailViewFactory
	if len(detailView) > 0 {
		dvFn = detailView[0]
	}

	return &App{
		viewStack:    nil,
		statusBar:    sb,
		keys:         keys,
		styles:       styles,
		issueClient:  client,
		repoName:     repoName,
		initView:     initView,
		detailViewFn: dvFn,
	}
}

// PushView pushes a view onto the stack and returns its Init command.
func (a *App) PushView(v View) tea.Cmd {
	a.viewStack = append(a.viewStack, v)
	a.updateKeyHints()
	return v.Init()
}

// PopView removes the top view from the stack if more than one view remains.
func (a *App) PopView() {
	if len(a.viewStack) > 1 {
		a.viewStack = a.viewStack[:len(a.viewStack)-1]
		a.updateKeyHints()
	}
}

// CurrentView returns the top view on the stack, or nil.
func (a *App) CurrentView() View {
	if len(a.viewStack) == 0 {
		return nil
	}
	return a.viewStack[len(a.viewStack)-1]
}

// ViewStackLen returns the number of views on the stack.
func (a *App) ViewStackLen() int {
	return len(a.viewStack)
}

func (a *App) updateKeyHints() {
	if v := a.CurrentView(); v != nil {
		a.statusBar.SetKeyHints(v.KeyHints())
	}
}

// IssueClient returns the issue client.
func (a *App) IssueClient() *data.IssueClient {
	return a.issueClient
}

// Styles returns the app styles.
func (a *App) Styles() Styles {
	return a.styles
}

// Width returns the current terminal width.
func (a *App) Width() int {
	return a.width
}

// Height returns the current terminal height (minus status bar).
func (a *App) Height() int {
	return a.height
}

// StatusBar returns the status bar component.
func (a *App) StatusBar() *components.StatusBar {
	return a.statusBar
}

// Keys returns the app key map.
func (a *App) Keys() KeyMap {
	return a.keys
}

// Init implements tea.Model.
func (a *App) Init() tea.Cmd {
	if a.initView != nil {
		v := a.initView(a)
		return a.PushView(v)
	}
	return nil
}

// Update implements tea.Model.
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height - 1 // Reserve 1 line for status bar
		a.statusBar.SetWidth(msg.Width)
		// Propagate resize to current view
		if v := a.CurrentView(); v != nil {
			updated, cmd := v.Update(msg)
			a.viewStack[len(a.viewStack)-1] = updated
			return a, cmd
		}
		return a, nil

	case tea.KeyMsg:
		if key.Matches(msg, a.keys.ForceQuit) {
			return a, tea.Quit
		}
		if key.Matches(msg, a.keys.Quit) && len(a.viewStack) <= 1 {
			return a, tea.Quit
		}

	case NavigateToDetailMsg:
		a.statusBar.SetMessage("")
		if a.detailViewFn != nil {
			v := a.detailViewFn(a, msg.IssueNumber)
			cmd := a.PushView(v)
			return a, cmd
		}
		return a, nil

	case NavigateBackMsg:
		a.PopView()
		a.statusBar.SetMessage("")
		return a, nil
	}

	// Delegate to current view
	if v := a.CurrentView(); v != nil {
		updated, cmd := v.Update(msg)
		a.viewStack[len(a.viewStack)-1] = updated
		return a, cmd
	}

	return a, nil
}

// View implements tea.Model.
func (a *App) View() string {
	if a.CurrentView() == nil {
		return "No view loaded"
	}

	viewContent := a.CurrentView().View()
	viewHeight := a.height
	content := lipgloss.NewStyle().Height(viewHeight).Render(viewContent)
	return content + "\n" + a.statusBar.View()
}
