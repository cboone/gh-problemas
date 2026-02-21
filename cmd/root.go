package cmd

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/cboone/gh-problemas/internal/config"
	"github.com/cboone/gh-problemas/internal/data"
	"github.com/cboone/gh-problemas/internal/ui"
	"github.com/cboone/gh-problemas/internal/ui/views"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "gh-problemas",
	Short:         "A terminal UI for triaging and managing GitHub issues",
	Long:          "gh-problemas is a terminal user interface for triaging and managing GitHub issues.",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          runApp,
}

func runApp(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	gqlClient, err := api.DefaultGraphQLClient()
	if err != nil {
		return fmt.Errorf("could not create GitHub API client: %w\nTry running: gh auth login", err)
	}

	owner, name, err := resolveRepository(cfg.Defaults.Repo, gqlClient)
	if err != nil {
		return err
	}

	repoName := owner + "/" + name
	pageSize := cfg.Defaults.PageSize
	dateFormat := cfg.Defaults.DateFormat

	issueClient := data.NewIssueClient(gqlClient, owner, name)
	commentClient := data.NewCommentClient(gqlClient, owner, name)
	app := ui.NewApp(
		issueClient,
		repoName,
		func(a *ui.App) ui.View {
			return views.NewDashboardViewWithPageSize(a.IssueClient(), a.Styles(), a.Keys(), a.Width(), a.Height(), pageSize)
		},
		func(a *ui.App, issueNumber int) ui.View {
			return views.NewDetailViewWithCommentsAndDateFormat(a.IssueClient(), commentClient, a.Styles(), a.Keys(), issueNumber, a.Width(), a.Height(), dateFormat)
		},
	)

	p := tea.NewProgram(app, tea.WithAltScreen())
	_, err = p.Run()
	return err
}

func resolveRepository(configRepo string, gqlClient data.Querier) (string, string, error) {
	owner := ""
	name := ""

	if configRepo != "" {
		parts := strings.Split(configRepo, "/")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return "", "", fmt.Errorf("invalid defaults.repo %q: expected owner/repo", configRepo)
		}
		owner = parts[0]
		name = parts[1]
	} else {
		repo, err := repository.Current()
		if err != nil {
			return "", "", fmt.Errorf("could not determine repository: %w\nRun this command from inside a git repository with a GitHub remote.", err)
		}
		owner = repo.Owner
		name = repo.Name
	}

	if owner == "@me" {
		login, err := data.NewUserClient(gqlClient).WhoAmI()
		if err != nil {
			return "", "", fmt.Errorf("resolving @me for defaults.repo: %w", err)
		}
		owner = login
	}

	return owner, name, nil
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

// SetVersion sets the version string on the root command.
func SetVersion(v string) {
	rootCmd.Version = v
}
