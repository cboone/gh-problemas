package ui

import "github.com/charmbracelet/lipgloss"

// Styles holds all application styles.
type Styles struct {
	App         lipgloss.Style
	Header      lipgloss.Style
	StatusBar   lipgloss.Style
	SelectedRow lipgloss.Style
	NormalRow   lipgloss.Style
	IssueNumber lipgloss.Style
	IssueTitle  lipgloss.Style
	LabelStyle  lipgloss.Style
	Spinner     lipgloss.Style
	ErrorText   lipgloss.Style
	HelpKey     lipgloss.Style
	HelpDesc    lipgloss.Style
}

// DefaultStyles returns the default application styles.
func DefaultStyles() Styles {
	return Styles{
		App:         lipgloss.NewStyle().Padding(0, 1),
		Header:      lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")),
		StatusBar:   lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Background(lipgloss.Color("236")).Padding(0, 1),
		SelectedRow: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")),
		NormalRow:   lipgloss.NewStyle(),
		IssueNumber: lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Width(6),
		IssueTitle:  lipgloss.NewStyle().Bold(true),
		LabelStyle:  lipgloss.NewStyle().Padding(0, 1),
		Spinner:     lipgloss.NewStyle().Foreground(lipgloss.Color("205")),
		ErrorText:   lipgloss.NewStyle().Foreground(lipgloss.Color("9")),
		HelpKey:     lipgloss.NewStyle().Foreground(lipgloss.Color("241")),
		HelpDesc:    lipgloss.NewStyle().Foreground(lipgloss.Color("245")),
	}
}
