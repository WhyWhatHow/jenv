package cmd

import (
	"fmt"
	"github.com/whywhathow/jenv/internal/java"
	"github.com/whywhathow/jenv/internal/style"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Aliases: []string{"rm"},
	Use:     "remove <name>",
	Short:   "Remove a Java JDK",
	Long: `Remove a Java JDK from the environment.

This command will remove the specified JDK from jenv-go's management.
It will not delete the actual JDK files from your system.

Use -f or --force flag to skip confirmation prompt.`,
	Example: ` jenv remove jdk8
jenv remove -f jdk11
jenv rm jdk8
jenv rm -f jdk11 `,
	Args: cobra.ExactArgs(1),
	Run:  runRemove,
}

var force bool

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation prompt")
}

func runRemove(cmd *cobra.Command, args []string) {
	name := args[0]

	// Get JDK info for confirmation
	jdks, err := java.ListJdks()
	if err != nil {
		fmt.Printf("%s: %v\n", style.Error.Render("Error"), style.Error.Render(err.Error()))
		return
	}

	jdk, exists := jdks[name]
	if !exists {
		fmt.Printf("%s: %s\n", style.Error.Render("Error"), style.Error.Render("JDK not found"))
		return
	}

	// Show JDK info and confirm removal
	if !force {
		fmt.Println(style.Header.Render("\nRemoving JDK"))
		fmt.Printf("%s: %s\n", style.Name.Render("Name"), style.Current.Render(jdk.Name))
		fmt.Printf("%s: %s\n\n", style.Name.Render("Path"), style.Path.Render(jdk.Path))

		fmt.Print(style.Input.Render("Are you sure you want to remove this JDK? [y/N] "))
		var confirm string
		fmt.Scanln(&confirm)

		if confirm != "y" && confirm != "Y" {
			fmt.Println(style.Input.Render("\nOperation cancelled"))
			return
		}
	}

	// Remove JDK
	if err := java.RemoveJDK(name); err != nil {
		fmt.Printf("%s: %v\n", style.Error.Render("Failed to remove JDK"), style.Error.Render(err.Error()))
		return
	}

	fmt.Printf("%s: %s\n", style.Success.Render("Successfully removed JDK"), style.Name.Render(name))
	return
}
