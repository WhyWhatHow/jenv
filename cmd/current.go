package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/whywhathow/jenv/internal/java"
	"github.com/whywhathow/jenv/internal/style"
)

var currentCmd = &cobra.Command{
	Aliases: []string{"cur", "now"},
	Use:     "current",
	Short:   "Show the current JDK",
	Long: `Show the current JDK configuration.

This command displays the name and path of the currently active JDK.
If no JDK is currently set, it will prompt you to configure one first.`,
	Example: `jenv current
jenv cur
jenv now`,
	Run: runCurrent,
}

func init() {
	rootCmd.AddCommand(currentCmd)
}

func runCurrent(cmd *cobra.Command, args []string) {
	// Get current JDK
	currentJDK, err := java.GetCurrentJDK()
	if err != nil {
		fmt.Println(style.Error.Render("No JDK is currently configured."))
		fmt.Println(style.Input.Render("Please use 'jenv use <name>' to set a JDK first."))
		return
	}

	// Display header
	fmt.Println(style.Header.Render("Current JDK Configuration"))

	// Display current JDK info
	fmt.Printf("%s: %s\n",
		style.Name.Render("Name"),
		style.Current.Render(currentJDK.Name))
	fmt.Printf("%s: %s\n",
		style.Name.Render("Path"),
		style.Current.Render(currentJDK.Path))
}
