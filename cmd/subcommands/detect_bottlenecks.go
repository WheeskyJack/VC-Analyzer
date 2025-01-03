package subcommands

import (
    "errors"
    "fmt"
    "github.com/MakeNowJust/heredoc/v2"
    "github.com/adigulalkari/VC-Analyzer/pkg/analyzer"
    "github.com/spf13/cobra"
    "os"
    "os/exec"
    "strings"
)

// Function to gather commit history data
func getCommitHistory(repoPath string) ([]analyzer.CommitInfo, error) {
    // Navigate to the repository path
    cmd := exec.Command("git", "-C", repoPath, "log", "--pretty=format:%H,%an,%ad,%s", "--numstat")
    output, err := cmd.CombinedOutput()
    if err != nil {
        return nil, fmt.Errorf("failed to run git command: %v", err)
    }

    // Parse the output
    lines := strings.Split(string(output), "\n")
    var commits []analyzer.CommitInfo
    var currentCommit *analyzer.CommitInfo
    for _, line := range lines {
        if strings.Contains(line, ",") {
            // New commit entry
            fields := strings.SplitN(line, ",", 4)
            currentCommit = &analyzer.CommitInfo{
                Hash:    fields[0],
                Author:  fields[1],
                Date:    fields[2],
                Message: fields[3],
                Files:   make(map[string]int),
            }
            commits = append(commits, *currentCommit)
        } else if currentCommit != nil {
            // File changes for the current commit
            fileFields := strings.Fields(line)
            if len(fileFields) == 3 {
                // Update file change details
                currentCommit.Files[fileFields[2]] = 1 // Mark file as changed
            }
        }
    }
    return commits, nil
}

var DetectBottlenecksCmd = &cobra.Command{
    Use:     "detect-bottlenecks <repository-path>",
    Short:   "Find bottlenecks in the commit history of a local repository",
    Aliases: []string{"d"},
    Example: heredoc.Doc(`
        $ vc-analyze detect-bottlenecks path/to/local/repo
    `),
    Args: func(cmd *cobra.Command, args []string) error {
        if len(args) < 1 {
            return errors.New("requires a repository path argument")
        }
        return nil
    },
    RunE: func(cmd *cobra.Command, args []string) error {
        repoPath := args[0] // Get the repository path from the arguments

        // Check if the repository exists
        if _, err := os.Stat(repoPath); os.IsNotExist(err) {
            return fmt.Errorf("repository path does not exist: %s", repoPath)
        }

        // Get the commit history for the repository
        commits, err := getCommitHistory(repoPath)
        if err != nil {
            return fmt.Errorf("error fetching commit history: %v", err)
        }

        // Detect bottlenecks based on the commit history
        // detectBottlenecks(commits)
        analyzer.DetectBottlenecks(commits)
        return nil
    },
}
