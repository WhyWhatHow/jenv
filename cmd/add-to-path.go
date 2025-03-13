package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/whywhathow/jenv/internal/env"
	"github.com/whywhathow/jenv/internal/style"
	"github.com/whywhathow/jenv/internal/sys"
	"os"
	"path/filepath"
)

var pathCmd = &cobra.Command{
	Aliases: []string{"path"},
	Use:     "add-to-path",
	Short:   "Add jenv to system PATH",
	Long: `Add jenv executable to system PATH environment variable.

This command will add the directory containing the jenv executable
to your system's PATH environment variable, allowing you to run
jenv from any location.

Requires administrator/root privileges to modify system environment variables.`,
	Example: `  jenv path
	jenv add-to-path`,
	Run: runPath,
}

func init() {
	rootCmd.AddCommand(pathCmd)
}

func runPath(cmd *cobra.Command, args []string) {
	fmt.Println(style.Header.Render("Adding jenv to System PATH"))
	if !sys.IsAdmin() {
		fmt.Println(style.Error.Render("Error: You must run this command as administrator/root"))
		os.Exit(1)
	}

	// Get the directory containing the jenv executable
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("%s: Unable to locate jenv executable\n%s\n",
			style.Error.Render("Error"),
			style.Error.Render(err.Error()))
		return
	}
	exeDir := filepath.Dir(exePath)

	// Check if path already exists
	if env.IsInPath(exeDir) {
		fmt.Printf("%s: %s\n",
			style.Info.Render("Notice"),
			style.Info.Render("jenv is already in your PATH"))
		return
	}

	// Add jenv directory to PATH
	if err := env.AddToPath(exeDir); err != nil {
		fmt.Printf("%s: Unable to modify PATH environment variable\n%s\n",
			style.Error.Render("Error"),
			style.Error.Render(err.Error()))
		return
	}

	fmt.Printf("%s\n", style.Success.Render("✓ Successfully added jenv to PATH"))
	fmt.Printf("%s: %s\n", style.Name.Render("Location"), style.Path.Render(exeDir))
	fmt.Printf("\n%s\n", style.Info.Render("➜ Please restart your terminal to apply the changes"))
	fmt.Printf("%s\n", style.Info.Render("➜ Then run 'jenv' to verify the installation"))
}
