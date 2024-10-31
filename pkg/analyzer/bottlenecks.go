package analyzer

import (
    "bytes"
    "fmt"
    "os/exec"
    "regexp"
    "strings"
)

// CommitInfo represents the details of a commit, including the commit hash, author, date,
// commit message, and a map of files changed in the commit along with their respective change counts.
type CommitInfo struct {
    Hash    string
    Author  string
    Date    string
    Message string
    Files   map[string]int 
}

// MergeCommit holds the information of a merge commit, including the commit hash,
// the commit message, and a list of files affected by the merge. This structure is used
// to track merge conflicts and the files involved.
type MergeCommit struct {
    CommitHash string
    Message    string
    Files      []string
}

func DetectBottlenecks(commits []CommitInfo){
    fileChangeCounts := make(map[string]int)
    authorChangeCounts := make(map[string]int)

    for _, commit := range commits {
        for file := range commit.Files {
            fileChangeCounts[file]++
        }
        authorChangeCounts[commit.Author]++
    }

    // Arbitarily choosing fileChange threshold as 5
    fmt.Println("Potential bottleneck files (most frequently changed):")
    for file, count := range fileChangeCounts {
        if count > 2 { 
            fmt.Printf("%s: %d changes\n", file, count)
        }
    }

    // Arbitarily choosing authorChange threshold as 5
    fmt.Println("\nPotential bottlenecks caused by contributors:")
    for author, count := range authorChangeCounts {
        if count > 5 { 
            fmt.Printf("%s: %d changes\n", author, count)
        }
    }

    fmt.Println("\nPotential bottlenecks caused by frequent rollbacks:")
    mergeCommits := findMergeConflicts()
    for _, merge := range mergeCommits {
        fmt.Printf("Merge Commit: %s, Message: %s, Affected Files: %v\n", merge.CommitHash, merge.Message, merge.Files)
    }
}

// findMergeConflicts retrieves the list of merge commits from the Git history,
// checks their commit messages for conflict indicators, and collects details
// about the files affected by these merges. This function aids in identifying 
// problematic merges that could slow down development.
func findMergeConflicts() []MergeCommit {
    cmd := exec.Command("git", "log", "--merges", "--pretty=format:%H %s")
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        fmt.Println("Error executing git command:", err)
        return nil
    }

    commitLines := strings.Split(out.String(), "\n")
    var mergeCommits []MergeCommit

    for _, line := range commitLines {
        parts := strings.SplitN(line, " ", 2)
        if len(parts) < 2 {
            continue
        }
        commitHash := parts[0]
        commitMessage := parts[1]

        if isMergeConflict(commitMessage) {
            files := getChangedFiles(commitHash)
            mergeCommits = append(mergeCommits, MergeCommit{
                CommitHash: commitHash,
                Message:    commitMessage,
                Files:      files,
            })
        }
    }

    return mergeCommits
}

// isMergeConflict checks if a given commit message indicates a merge conflict
// by searching for the keyword "conflict". This helps to filter out merge commits 
// that are likely to have encountered issues, facilitating focused analysis on 
// problematic areas during conflict resolution.
func isMergeConflict(message string) bool {
    conflictPattern := `(?i)conflict`
    return regexp.MustCompile(conflictPattern).FindString(message) != ""
}

// getChangedFiles retrieves a list of files that were changed in a specific
// merge commit identified by its hash. This is crucial for understanding 
// the impact of a merge and identifying which files may require attention 
// during conflict resolution.
func getChangedFiles(commitHash string) []string {
    cmd := exec.Command("git", "diff-tree", "--no-commit-id", "--name-only", "-r", commitHash)
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        fmt.Println("Error executing git command:", err)
        return nil
    }

    files := strings.Split(out.String(), "\n")
    return files[:len(files)-1] 
}

