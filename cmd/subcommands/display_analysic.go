package subcommands

import (
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var CmdCmd = &cobra.Command{
	Use:   "cmd",
	Short: "Display help commands in table format",
	Long:  `Displays all available commands along with their flags and shorthand in table.`,
	Example: heredoc.Doc(`#Display available commands in table format
				$vc-analyze cmd

				# Display specific help for a command, e.g., calc-stats
				$vc-analyze c -h

				#Display specific help command, e.g., check-anti-patterns
				$vc-analyze p -h`),
	Run: func(cmd *cobra.Command, args []string) {
		printCmdHelpTable()
	},
}

func printCmdHelpTable() {
	fmt.Println("Usage: vc-analyze [command/short-hand] [-h/--help]")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"command", "short-hand", "Description"})
	//Set minimum column width for 2nd column
	table.SetColMinWidth(2, 35)
	table.SetAutoWrapText(false)

	table.Append([]string{"check-anti-patterns", "p", "Find out the anti-patterns present in your repository"})
	table.Append([]string{"calc-stats", "c", "Calculate statistics for the local repository"})
	table.Append([]string{"detect-bottlenecks", "d", "Find bottlenecks in the commit history of a local repository"})
	table.Append([]string{"get", "g", "Display one or many repositories details"})

	table.SetBorder(true)
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor}, // Color for "Flag" header
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor}, // Color for "Short-hand" header
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor}, //  Color for Description header
	)
	table.SetColumnColor(
		tablewriter.Colors{tablewriter.FgHiBlackColor}, // Color for "Flag" column
		tablewriter.Colors{tablewriter.FgHiBlackColor}, // Color for "Description" column
		tablewriter.Colors{tablewriter.FgHiBlackColor}, // Color for Short-hand column
	)

	// Render the table to the console
	table.Render()
}

func init() {
	CmdCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		printCmdHelpTable() // Custom help function in table format
	})
}
