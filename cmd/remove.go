package cmd

import (
	"fmt"
	"github.com/whywhathow/jenv/internal/java"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a Java JDK",
	Long: `Remove a Java JDK from the environment.

This command will remove the specified JDK from jenv-go's management.
It will not delete the actual JDK files from your system.`,
	Example: "  jenv remove jdk8",
	Args:    cobra.ExactArgs(1),
	Run:     runRemove,
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

func runRemove(cmd *cobra.Command, args []string) {
	name := args[0]

	// Remove JDK
	if err := java.RemoveJDK(name); err != nil {
		fmt.Printf("failed to remove JDK: %v\n", err)
		return
	}

	fmt.Printf("Successfully removed JDK: %s\n", name)
	return
}
