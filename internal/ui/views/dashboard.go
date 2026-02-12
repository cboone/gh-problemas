package views

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hpg/gh-problemas/internal/data"
	"github.com/hpg/gh-problemas/internal/ui"
	"github.com/hpg/gh-problemas/internal/ui/components"
	"github.com/hpg/gh-problemas/internal/utils"
)

// issueItem wraps a data.Issue for the list component.
type issueItem struct {
	issue data.Issue
}

func (i issueItem) FilterValue() string { return i.issue.Title }
func (i issueItem) Title() string       { return i.issue.Title }
func (i issueItem) Description() string {
	parts := []string{
		fmt.Sprintf("#%d", i.issue.Number),
		i.issue.Author,
		utils.RelativeTime(i.issue.CreatedAt),
	}
	if i.issue.CommentCount > 0 {
		parts = append(parts, fmt.Sprintf("%d comments", i.issue.CommentCount))
	}
	return strings.Join(parts, " | ")
}

// issueDelegate renders issue items in the list.
type issueDelegate struct {
	styles ui.Styles
}

func (d issueDelegate) Height() int                         { return 2 }
func (d issueDelegate) Spacing() int                        { return 0 }
func (d issueDelegate) Update(tea.Msg, *list.Model) tea.Cmd { return nil }

func (d issueDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(issueItem)
	if !ok {
		return
	}

	isSelected := index == m.Index()

	// Build label string
	var labelParts []string
	for _, l := range i.issue.Labels {
		bg := utils.HexToColor(l.Color)
		fg := utils.ContrastColor(l.Color)
		style := lipgloss.NewStyle().Background(bg).Foreground(fg).Padding(0, 1)
		labelParts = append(labelParts, style.Render(l.Name))
	}
	labels := strings.Join(labelParts, " ")

	// Title line
	numberStyle := d.styles.IssueNumber
	titleStyle := d.styles.IssueTitle
	if isSelected {
		numberStyle = numberStyle.Foreground(lipgloss.Color("12"))
		titleStyle = titleStyle.Foreground(lipgloss.Color("12"))
	}

	titleLine := numberStyle.Render(fmt.Sprintf("#%-5d", i.issue.Number)) + " " + titleStyle.Render(i.issue.Title)
	if labels != "" {
		titleLine += " " + labels
	}

	// Meta line
	meta := fmt.Sprintf("       %s  %s", i.issue.Author, utils.RelativeTime(i.issue.CreatedAt))
	if i.issue.CommentCount > 0 {
		meta += fmt.Sprintf("  %d comments", i.issue.CommentCount)
	}
	if i.issue.ReactionCount > 0 {
		meta += fmt.Sprintf("  %d reactions", i.issue.ReactionCount)
	}

	metaStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	if isSelected {
		metaStyle = metaStyle.Foreground(lipgloss.Color("244"))
	}

	cursor := "  "
	if isSelected {
		cursor = "> "
	}

	fmt.Fprintf(w, "%s%s\n%s%s", cursor, titleLine, "  ", metaStyle.Render(meta))
}

// DashboardView is the main view showing open issues.
type DashboardView struct {
	list        list.Model
	issueClient *data.IssueClient
	paginator   *data.Paginator
	spinner     *components.Spinner
	styles      ui.Styles
	keys        ui.KeyMap
	loading     bool
	loadingMore bool
	errMsg      string
	width       int
	height      int
	pageSize    int
}

// NewDashboardView creates a new dashboard view.
func NewDashboardView(client *data.IssueClient, styles ui.Styles, keys ui.KeyMap, width, height int) *DashboardView {
	return NewDashboardViewWithPageSize(client, styles, keys, width, height, 50)
}

// NewDashboardViewWithPageSize creates a new dashboard view with a custom page size.
func NewDashboardViewWithPageSize(client *data.IssueClient, styles ui.Styles, keys ui.KeyMap, width, height, pageSize int) *DashboardView {
	delegate := issueDelegate{styles: styles}
	l := list.New(nil, delegate, width, height)
	l.SetShowTitle(true)
	l.Title = "Open Issues"
	l.SetShowStatusBar(true)
	l.SetShowFilter(false)
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()

	spinner := components.NewSpinner(styles.Spinner)
	paginator := data.NewPaginator(pageSize)

	return &DashboardView{
		list:        l,
		issueClient: client,
		paginator:   paginator,
		spinner:     spinner,
		styles:      styles,
		keys:        keys,
		loading:     true,
		width:       width,
		height:      height,
		pageSize:    pageSize,
	}
}

// Init implements ui.View.
func (d *DashboardView) Init() tea.Cmd {
	client := d.issueClient
	pageSize := d.pageSize
	spinCmd := d.spinner.Start("Loading issues...")
	statusCmd := ui.StatusLoading("Loading issues...")
	fetchCmd := func() tea.Msg {
		result, err := client.List(data.IssueListOptions{
			States: []string{"OPEN"},
			First:  pageSize,
		})
		return ui.IssuesLoadedMsg{Result: result, Err: err}
	}
	return tea.Batch(spinCmd, statusCmd, fetchCmd)
}

// Update implements ui.View.
func (d *DashboardView) Update(msg tea.Msg) (ui.View, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		d.width = msg.Width
		d.height = msg.Height - 1
		d.list.SetSize(msg.Width, d.height)
		return d, nil

	case ui.IssuesLoadedMsg:
		d.loading = false
		d.spinner.Stop()
		if msg.Err != nil {
			d.errMsg = fmt.Sprintf("Error loading issues: %v", msg.Err)
			return d, nil
		}
		d.errMsg = ""
		d.paginator.Reset()
		d.paginator.Update(msg.Result.PageInfo, len(msg.Result.Issues))
		items := make([]list.Item, len(msg.Result.Issues))
		for i, issue := range msg.Result.Issues {
			items[i] = issueItem{issue: issue}
		}
		cmd := d.list.SetItems(items)
		d.updateTitle()
		statusCmd := ui.StatusInfo(fmt.Sprintf("Showing %d issues", d.paginator.TotalLoaded()))
		return d, tea.Batch(cmd, statusCmd)

	case ui.IssuesPageLoadedMsg:
		d.loadingMore = false
		d.spinner.Stop()
		if msg.Err != nil {
			d.errMsg = fmt.Sprintf("Error loading more issues: %v", msg.Err)
			return d, nil
		}
		d.errMsg = ""
		d.paginator.Update(msg.Result.PageInfo, len(msg.Result.Issues))
		// Append new items to existing list
		existing := d.list.Items()
		for _, issue := range msg.Result.Issues {
			existing = append(existing, issueItem{issue: issue})
		}
		cmd := d.list.SetItems(existing)
		d.updateTitle()
		statusCmd := ui.StatusInfo(fmt.Sprintf("Showing %d issues", d.paginator.TotalLoaded()))
		return d, tea.Batch(cmd, statusCmd)

	case tea.KeyMsg:
		if key.Matches(msg, d.keys.Open) {
			item, ok := d.list.SelectedItem().(issueItem)
			if ok {
				return d, func() tea.Msg {
					return ui.NavigateToDetailMsg{IssueNumber: item.issue.Number}
				}
			}
		}
		if key.Matches(msg, d.keys.Refresh) {
			d.loading = true
			d.errMsg = ""
			client := d.issueClient
			pageSize := d.pageSize
			spinCmd := d.spinner.Start("Refreshing...")
			statusCmd := ui.StatusLoading("Refreshing issues...")
			fetchCmd := func() tea.Msg {
				result, err := client.List(data.IssueListOptions{
					States: []string{"OPEN"},
					First:  pageSize,
				})
				return ui.IssuesLoadedMsg{Result: result, Err: err}
			}
			return d, tea.Batch(spinCmd, statusCmd, fetchCmd)
		}
		if key.Matches(msg, d.keys.NextPage) && !d.loading && !d.loadingMore {
			req := d.paginator.NextPageRequest()
			if req != nil {
				d.loadingMore = true
				client := d.issueClient
				after := req.After
				first := req.First
				spinCmd := d.spinner.Start("Loading more...")
				statusCmd := ui.StatusLoading("Loading more issues...")
				fetchCmd := func() tea.Msg {
					result, err := client.List(data.IssueListOptions{
						States: []string{"OPEN"},
						First:  first,
						After:  after,
					})
					return ui.IssuesPageLoadedMsg{Result: result, Err: err}
				}
				return d, tea.Batch(spinCmd, statusCmd, fetchCmd)
			}
			return d, ui.StatusInfo(fmt.Sprintf("Showing %d issues", d.paginator.TotalLoaded()))
		}
	}

	// Update spinner
	spinCmd := d.spinner.Update(msg)
	if spinCmd != nil {
		cmds = append(cmds, spinCmd)
	}

	// Delegate to list
	var listCmd tea.Cmd
	d.list, listCmd = d.list.Update(msg)
	if listCmd != nil {
		cmds = append(cmds, listCmd)
	}

	return d, tea.Batch(cmds...)
}

// View implements ui.View.
func (d *DashboardView) View() string {
	if d.loading {
		return lipgloss.Place(d.width, d.height, lipgloss.Center, lipgloss.Center, d.spinner.View())
	}

	if d.errMsg != "" {
		errView := d.styles.ErrorText.Render(d.errMsg)
		return lipgloss.Place(d.width, d.height, lipgloss.Center, lipgloss.Center, errView)
	}

	return d.list.View()
}

// KeyHints implements ui.View.
func (d *DashboardView) KeyHints() []string {
	hints := []string{"j/k: navigate", "enter: open", "R: refresh"}
	if d.paginator.HasNextPage() {
		hints = append(hints, "L: load more")
	}
	hints = append(hints, "q: quit")
	return hints
}

func (d *DashboardView) updateTitle() {
	total := d.paginator.TotalLoaded()
	if d.paginator.HasNextPage() {
		d.list.Title = fmt.Sprintf("Open Issues (showing %d+)", total)
	} else {
		d.list.Title = fmt.Sprintf("Open Issues (%d)", total)
	}
}
