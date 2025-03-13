package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/whywhathow/jenv/internal/java"
)

var addCmd = &cobra.Command{
	Use:   "add [flags] <name> <path>",
	Short: "Add a new Java JDK",
	Long: `Add a new Java JDK to the environment.

This command allows you to register a new Java JDK installation
by providing a name and the path to the JDK installation directory.`,
	Example: `  jenvadd jdk8 "C:\Program Files\Java\jdk1.8.0_291"
  jenvadd -f jdk11 "C:\Program Files\Java\jdk-11.0.12"`,
	Args: cobra.ExactArgs(2),
	Run:  runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) {
	name := args[0]
	path := args[1]

	// Add JDK
	if err := java.AddJDK(name, path); err != nil {
		fmt.Printf("failed to add JDK: %v\n", err)
		return
	}

	fmt.Printf("Successfully added JDK: %s -> %s\n", name, path)
}
