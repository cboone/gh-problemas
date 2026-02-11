package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hpg/gh-problemas/internal/data"
	"github.com/hpg/gh-problemas/internal/ui"
	"github.com/hpg/gh-problemas/internal/ui/components"
	"github.com/hpg/gh-problemas/internal/utils"
)

// DetailView shows a single issue with its rendered markdown body and comments.
type DetailView struct {
	viewport        viewport.Model
	issueClient     *data.IssueClient
	commentClient   *data.CommentClient
	spinner         *components.Spinner
	styles          ui.Styles
	keys            ui.KeyMap
	dateFormat      string
	issueNumber     int
	issue           *data.Issue
	comments        []data.Comment
	loading         bool
	loadingComments bool
	errMsg          string
	width           int
	height          int
}

// NewDetailView creates a new detail view for the given issue number.
func NewDetailView(client *data.IssueClient, styles ui.Styles, keys ui.KeyMap, issueNumber, width, height int) *DetailView {
	return NewDetailViewWithCommentsAndDateFormat(client, nil, styles, keys, issueNumber, width, height, "relative")
}

// NewDetailViewWithComments creates a detail view with a comment client.
func NewDetailViewWithComments(client *data.IssueClient, commentClient *data.CommentClient, styles ui.Styles, keys ui.KeyMap, issueNumber, width, height int) *DetailView {
	return NewDetailViewWithCommentsAndDateFormat(client, commentClient, styles, keys, issueNumber, width, height, "relative")
}

// NewDetailViewWithCommentsAndDateFormat creates a detail view with comments and configurable timestamp formatting.
func NewDetailViewWithCommentsAndDateFormat(client *data.IssueClient, commentClient *data.CommentClient, styles ui.Styles, keys ui.KeyMap, issueNumber, width, height int, dateFormat string) *DetailView {
	vp := viewport.New(width, height)
	vp.SetContent("Loading...")
	spinner := components.NewSpinner(styles.Spinner)
	if dateFormat == "" {
		dateFormat = "relative"
	}

	return &DetailView{
		viewport:      vp,
		issueClient:   client,
		commentClient: commentClient,
		spinner:       spinner,
		styles:        styles,
		keys:          keys,
		dateFormat:    dateFormat,
		issueNumber:   issueNumber,
		loading:       true,
		width:         width,
		height:        height,
	}
}

// Init implements ui.View.
func (d *DetailView) Init() tea.Cmd {
	client := d.issueClient
	number := d.issueNumber
	spinCmd := d.spinner.Start("Loading issue...")
	statusCmd := ui.StatusLoading("Loading issue...")
	fetchCmd := func() tea.Msg {
		issue, err := client.Get(number)
		return ui.IssueDetailLoadedMsg{Issue: issue, Err: err}
	}
	return tea.Batch(spinCmd, statusCmd, fetchCmd)
}

// Update implements ui.View.
func (d *DetailView) Update(msg tea.Msg) (ui.View, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		d.width = msg.Width
		d.height = msg.Height - 1
		d.viewport.Width = msg.Width
		d.viewport.Height = d.height
		return d, nil

	case ui.IssueDetailLoadedMsg:
		d.loading = false
		d.spinner.Stop()
		if msg.Err != nil {
			d.errMsg = fmt.Sprintf("Error loading issue: %v", msg.Err)
			return d, nil
		}
		d.errMsg = ""
		d.issue = &msg.Issue
		d.renderContent()
		// Fetch comments if we have a comment client
		if d.commentClient != nil {
			d.loadingComments = true
			cc := d.commentClient
			number := d.issueNumber
			statusCmd := ui.StatusLoading("Loading comments...")
			fetchCmd := func() tea.Msg {
				result, err := cc.List(number, 25, "")
				return ui.CommentsLoadedMsg{
					Comments: result.Comments,
					PageInfo: result.PageInfo,
					Err:      err,
				}
			}
			return d, tea.Batch(statusCmd, fetchCmd)
		}
		return d, ui.StatusInfo(fmt.Sprintf("Loaded issue #%d", msg.Issue.Number))

	case ui.CommentsLoadedMsg:
		d.loadingComments = false
		if msg.Err != nil {
			d.errMsg = fmt.Sprintf("Error loading comments: %v", msg.Err)
			return d, nil
		}
		d.comments = msg.Comments
		d.renderContent()
		if len(msg.Comments) == 0 {
			return d, ui.StatusInfo("No comments")
		}
		return d, ui.StatusInfo(fmt.Sprintf("Loaded %d comments", len(msg.Comments)))

	case tea.KeyMsg:
		if key.Matches(msg, d.keys.Back) {
			return d, func() tea.Msg { return ui.NavigateBackMsg{} }
		}
		if key.Matches(msg, d.keys.Quit) {
			return d, func() tea.Msg { return ui.NavigateBackMsg{} }
		}
	}

	// Update spinner
	spinCmd := d.spinner.Update(msg)
	if spinCmd != nil {
		cmds = append(cmds, spinCmd)
	}

	// Delegate to viewport
	var vpCmd tea.Cmd
	d.viewport, vpCmd = d.viewport.Update(msg)
	if vpCmd != nil {
		cmds = append(cmds, vpCmd)
	}

	return d, tea.Batch(cmds...)
}

// View implements ui.View.
func (d *DetailView) View() string {
	if d.loading {
		return lipgloss.Place(d.width, d.height, lipgloss.Center, lipgloss.Center, d.spinner.View())
	}

	if d.errMsg != "" {
		errView := d.styles.ErrorText.Render(d.errMsg)
		return lipgloss.Place(d.width, d.height, lipgloss.Center, lipgloss.Center, errView)
	}

	return d.viewport.View()
}

// KeyHints implements ui.View.
func (d *DetailView) KeyHints() []string {
	return []string{"j/k: scroll", "esc: back", "q: back"}
}

func (d *DetailView) renderContent() {
	if d.issue == nil {
		return
	}

	var sb strings.Builder
	issue := d.issue

	// Header
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	numberStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	sb.WriteString(titleStyle.Render(issue.Title))
	sb.WriteString(" ")
	sb.WriteString(numberStyle.Render(fmt.Sprintf("#%d", issue.Number)))
	sb.WriteString("\n")

	// Metadata line
	metaParts := []string{
		fmt.Sprintf("State: %s", issue.State),
		fmt.Sprintf("Author: %s", issue.Author),
		fmt.Sprintf("Created: %s", utils.FormatTime(issue.CreatedAt, d.dateFormat)),
		fmt.Sprintf("Updated: %s", utils.FormatTime(issue.UpdatedAt, d.dateFormat)),
	}
	if issue.Milestone != "" {
		metaParts = append(metaParts, fmt.Sprintf("Milestone: %s", issue.Milestone))
	}
	if len(issue.Assignees) > 0 {
		metaParts = append(metaParts, fmt.Sprintf("Assignees: %s", strings.Join(issue.Assignees, ", ")))
	}

	metaStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	sb.WriteString(metaStyle.Render(strings.Join(metaParts, "  ")))
	sb.WriteString("\n")

	// Labels
	if len(issue.Labels) > 0 {
		var labelParts []string
		for _, l := range issue.Labels {
			bg := utils.HexToColor(l.Color)
			fg := utils.ContrastColor(l.Color)
			style := lipgloss.NewStyle().Background(bg).Foreground(fg).Padding(0, 1)
			labelParts = append(labelParts, style.Render(l.Name))
		}
		sb.WriteString(strings.Join(labelParts, " "))
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	divider := lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Render(strings.Repeat("─", d.width))
	sb.WriteString(divider)
	sb.WriteString("\n\n")

	// Body
	if issue.Body != "" {
		rendered, err := utils.RenderMarkdown(issue.Body, d.width-4)
		if err != nil {
			sb.WriteString(issue.Body)
		} else {
			sb.WriteString(rendered)
		}
	} else {
		sb.WriteString(metaStyle.Render("No description provided."))
	}

	// Comments
	if len(d.comments) > 0 {
		sb.WriteString("\n")
		sb.WriteString(divider)
		sb.WriteString("\n")
		commentHeaderStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
		sb.WriteString(commentHeaderStyle.Render(fmt.Sprintf("Comments (%d)", len(d.comments))))
		sb.WriteString("\n\n")

		authorStyle := lipgloss.NewStyle().Bold(true)
		timeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

		for i, c := range d.comments {
			sb.WriteString(authorStyle.Render(c.Author))
			sb.WriteString(" ")
			sb.WriteString(timeStyle.Render(utils.FormatTime(c.CreatedAt, d.dateFormat)))
			if c.Reactions > 0 {
				sb.WriteString(timeStyle.Render(fmt.Sprintf("  %d reactions", c.Reactions)))
			}
			sb.WriteString("\n")

			if c.Body != "" {
				rendered, err := utils.RenderMarkdown(c.Body, d.width-4)
				if err != nil {
					sb.WriteString(c.Body)
				} else {
					sb.WriteString(rendered)
				}
			}

			if i < len(d.comments)-1 {
				sb.WriteString("\n")
				thinDivider := lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Render(strings.Repeat("- ", d.width/2))
				sb.WriteString(thinDivider)
				sb.WriteString("\n\n")
			}
		}
	} else if d.loadingComments {
		sb.WriteString("\n")
		sb.WriteString(metaStyle.Render("Loading comments..."))
	}

	d.viewport.SetContent(sb.String())
}
