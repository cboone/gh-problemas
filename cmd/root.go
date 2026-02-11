package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/hpg/gh-problemas/internal/config"
	"github.com/hpg/gh-problemas/internal/data"
	"github.com/hpg/gh-problemas/internal/ui"
	"github.com/hpg/gh-problemas/internal/ui/views"
	"github.com/spf13/cobra"
)

// version is set via ldflags at build time.
var version = "dev"

var rootCmd = &cobra.Command{
	Use:     "gh-problemas",
	Short:   "A TUI for GitHub issue management",
	Long:    "gh-problemas is a terminal user interface for triaging and managing GitHub issues.",
	Version: version,
	RunE:    runApp,
}

func runApp(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	repo, err := repository.Current()
	if err != nil {
		return fmt.Errorf("could not determine repository: %w\nRun this command from inside a git repository with a GitHub remote.", err)
	}

	gqlClient, err := api.DefaultGraphQLClient()
	if err != nil {
		return fmt.Errorf("could not create GitHub API client: %w\nTry running: gh auth login", err)
	}

	owner := repo.Owner
	name := repo.Name
	repoName := owner + "/" + name
	pageSize := cfg.Defaults.PageSize

	issueClient := data.NewIssueClient(gqlClient, owner, name)
	commentClient := data.NewCommentClient(gqlClient, owner, name)
	app := ui.NewApp(
		issueClient,
		repoName,
		func(a *ui.App) ui.View {
			return views.NewDashboardViewWithPageSize(a.IssueClient(), a.Styles(), a.Keys(), a.Width(), a.Height(), pageSize)
		},
		func(a *ui.App, issueNumber int) ui.View {
			return views.NewDetailViewWithComments(a.IssueClient(), commentClient, a.Styles(), a.Keys(), issueNumber, a.Width(), a.Height())
		},
	)

	p := tea.NewProgram(app, tea.WithAltScreen())
	_, err = p.Run()
	return err
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
