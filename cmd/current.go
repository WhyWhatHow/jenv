package cmd

import (
	"fmt"
	"github.com/whywhathow/jenv/internal/java"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show the current JDK",
	Long: `Show the current JDK configuration.

This command displays the name and path of the currently active JDK.
If no JDK is currently set, it will prompt you to configure one first.`,
	Example: "  jenv current",
	Run:     runCurrent,
}

func init() {
	rootCmd.AddCommand(currentCmd)
}

func runCurrent(cmd *cobra.Command, args []string) {
	// Get current JDK
	currentJDK, err := java.GetCurrentJDK()
	if err != nil {
		fmt.Println("No JDK is currently configured.")
		fmt.Println("Please use 'jenv use <name>' to set a JDK first.")
		return
	}

	//TODO [whywhathow] [2025/3/12]  [opt] 找第三方库 更好的输出到Terminal

	// Use tabwriter for formatted output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tPATH")
	fmt.Fprintf(w, "%s\t%s\n", currentJDK.Name, currentJDK.Path)
	w.Flush()
}
