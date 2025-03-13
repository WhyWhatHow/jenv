package cmd

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/whywhathow/jenv/internal/config"
	"github.com/whywhathow/jenv/internal/java"
	"github.com/whywhathow/jenv/internal/style"
	"os"
	"sort"
)

var listCmd = &cobra.Command{
	Aliases: []string{"ls"},
	Use:     "list",
	Short:   "List all Java JDKs",
	Long: `List all Java JDKs registered in the environment.

This command displays all registered JDK installations,
showing their names, paths, and which one is currently active.`,
	Example: `  jenv list
jenv ls`,
	Run: RunList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func RunList(cmd *cobra.Command, args []string) {
	jdks, err := java.ListJdks()
	if err != nil {
		fmt.Println(style.Error.Render("Error:"), style.Error.Render(err.Error()))
		return
	}

	if len(jdks) == 0 {
		fmt.Println(style.Path.Render("＞ No JDKs registered"))
		return
	}

	// Get current JDK and sort
	currentJDK, err := java.GetCurrentJDK()
	// Display current JDK info if available
	if err == nil {
		fmt.Printf("\n%s\n", style.Header.Render("Current JDK"))
		fmt.Printf("%s: %s\n", style.Name.Render("Name"), style.Current.Render(currentJDK.Name))
		fmt.Printf("%s: %s\n\n", style.Name.Render("Path"), style.Path.Render(currentJDK.Path))
	}

	var sorted []config.JDK
	for _, jdk := range jdks {
		sorted = append(sorted, jdk)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Name < sorted[j].Name
	})

	// Create and configure table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Path", "Current"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetTablePadding(" ")
	table.SetNoWhiteSpace(true)

	//var[]string data;
	// Add data rows
	for _, jdk := range sorted {
		currentMark := "  "
		name := style.Name.Render(jdk.Name)
		path := style.Path.Render(jdk.Path)
		if jdk.Name == currentJDK.Name {
			currentMark = style.Current.Render("✓")
			name = style.Current.Render(jdk.Name)
			path = style.Current.Render(jdk.Path)
		}

		table.Append([]string{
			name,
			path,
			currentMark,
		})
	}

	// Render the table
	table.Render()
}
