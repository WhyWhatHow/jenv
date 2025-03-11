package cmd

import (
	"fmt"
	"jenv-go/internal/java"

	"github.com/spf13/cobra"
)

var (
	useCmd = &cobra.Command{
		Use:   "use <name>",
		Short: "Switch to a different Java JDK",
		Long: `Switch to a different Java JDK version.

This command will set the specified JDK as the current Java version
for your environment.`,
		Example: "  jenv use jdk8",
		Args:    cobra.ExactArgs(1),
		Run:     runUse,
	}
)

func init() {
	rootCmd.AddCommand(useCmd)
}

func runUse(cmd *cobra.Command, args []string) {
	name := args[0]

	// Switch JDK
	if err := java.UseJDK(name); err != nil {
		fmt.Printf("failed to switch JDK: %v\n", err)
		return
	}

	fmt.Printf("Successfully switched to JDK: %s\n", name)
}
