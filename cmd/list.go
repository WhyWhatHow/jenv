package cmd

import (
	"fmt"
	"github.com/whywhathow/jenv/internal/java"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Java JDKs",
	Long: `List all Java JDKs registered in the environment.

This command displays all registered JDK installations,
showing their names, paths, and which one is currently active.`,
	Example: "  jenv list",
	Run:     RunList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func RunList(cmd *cobra.Command, args []string) {
	// Get all JDKs
	jdks, err := java.ListJdks()

	if err != nil {
		fmt.Printf("failed to get JDK list: %v\n", err)
		return
	}

	if len(jdks) == 0 {
		fmt.Println("No JDKs are currently registered")
		return
	}

	// Get current JDK
	currentJDK, err := java.GetCurrentJDK()
	currentName := ""
	if err == nil {
		currentName = currentJDK.Name
	}

	// Use tabwriter for formatted output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tPATH\tCURRENT")

	for _, jdk := range jdks {
		current := ""
		if jdk.Name == currentName {
			current = "*"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", jdk.Name, jdk.Path, current)
	}

	w.Flush()
}
