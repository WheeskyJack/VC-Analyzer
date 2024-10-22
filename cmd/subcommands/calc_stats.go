package subcommands

import (
	"errors"
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/adigulalkari/VC-Analyzer/pkg/analyzer"
)

var (
	authorStats  bool
	commitSize   bool
	activeBranch bool
)

var CalcStatsCmd = &cobra.Command{
	Use:   "calc-stats <path/to/repo>",
	Short: "Calculate statistics for the local repo",
	Long:  `This command allows you to calculate various statistics for a local Git repository, including author statistics and commit size statistics.`,
	Aliases: []string{"c"},
	Example: heredoc.Doc(`
        # Calculate author statistics
        $ vc-analyze calc-stats --author-stats path/to/local/repo

        # Calculate commit size statistics
        $ vc-analyze calc-stats --commit-size path/to/local/repo

        # Calculate branch statistics
        $ vc-analyze calc-stats --active-branch path/to/local/repo
    `),
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a path to the repository")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		repoPath := args[0] // Get the repository path from the arguments

		// Check if the repository exists (optional)
		if _, err := os.Stat(repoPath); os.IsNotExist(err) {
			return fmt.Errorf("repository path does not exist: %s", repoPath)
		}

		// Check which flag is set and call the appropriate function
		if authorStats {
			// Call a function to calculate author statistics
			analyzer.AnalyzeCommitHistory(repoPath)
		} else if commitSize {
			// Call a function to calculate commit size statistics
			analyzer.AnalyzeCommitSize(repoPath)
		} else if activeBranch {
			// Call function to show branch statistics
			analyzer.AnalyzeBranchStats(repoPath)
		} else {
			return errors.New("no valid flag provided, use --author-stats , --commit-size or --active-branch")
		}

		return nil
	},
}

func init() {
	CalcStatsCmd.Flags().BoolVarP(&authorStats, "author-stats", "a", false, "Calculate statistics for each author")
	CalcStatsCmd.Flags().BoolVarP(&commitSize, "commit-size", "s", false, "Calculate the size of commits")
	CalcStatsCmd.Flags().BoolVarP(&activeBranch, "active-branch", "b", false, "Show branch statistics")

	// Custom help function to display help message clearly
	CalcStatsCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		helpMessage := `
Usage:
  vc-analyze calc-stats <path/to/repo> [flags]

This command allows you to calculate various statistics for a local Git repository.

Flags:
  -a, --author-stats    Calculate statistics for each author
  -s, --commit-size     Calculate the size of commits
  -b, --active-branch   Show branch statistics
  -h, --help            help for calc-stats

Examples:
# Calculate author statistics
$ vc-analyze calc-stats --author-stats path/to/local/repo

# Calculate commit size statistics
$ vc-analyze calc-stats --commit-size path/to/local/repo

# Calculate branch statistics
$ vc-analyze calc-stats --active-branch path/to/local/repo
`
		fmt.Fprint(cmd.OutOrStdout(), helpMessage)
	})
}