package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/whywhathow/jenv/internal/java"
	"github.com/whywhathow/jenv/internal/style"
	"github.com/whywhathow/jenv/internal/sys"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize jenv for first-time use",
	Long: `Initialize jenv for first-time use.

This command sets up the necessary environment variables and symbolic links
for jenv to manage Java versions on your system.

On Windows: Requires administrator privileges
On Linux/macOS: Can run with or without root privileges
  - With root: Creates system-wide configuration
  - Without root: Creates user-level configuration`,
	Example: `  jenv init`,
	Run:     runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) {
	fmt.Println(style.Header.Render("🚀 Initializing jenv..."))

	// Check platform-specific requirements
	if runtime.GOOS == "windows" {
		if !sys.IsAdmin() {
			fmt.Printf("%s: %s\n",
				style.Error.Render("Error"),
				style.Error.Render("Administrator privileges required on Windows"))
			return
		}
	} else {
		if !sys.IsAdmin() {
			fmt.Printf("%s: %s\n",
				style.Warning.Render("Warning"),
				style.Warning.Render("Running without root privileges"))
			fmt.Printf("%s: %s\n",
				style.Info.Render("Info"),
				style.Info.Render("Will create user-level configuration"))
		} else {
			fmt.Printf("%s: %s\n",
				style.Success.Render("Info"),
				style.Success.Render("Running with root privileges - will create system-wide configuration"))
		}
	}

	// Perform initialization
	if err := java.Init(); err != nil {
		fmt.Printf("%s: %s\n",
			style.Error.Render("Error"),
			style.Error.Render(fmt.Sprintf("Initialization failed: %v", err)))
		return
	}

	fmt.Printf("%s: %s\n",
		style.Success.Render("Success"),
		style.Success.Render("jenv has been initialized successfully!"))

	// Provide next steps
	fmt.Println()
	fmt.Println(style.Header.Render("📋 Next Steps:"))
	fmt.Printf("  %s %s\n",
		style.Info.Render("1."),
		style.Input.Render("Add jenv to your PATH: jenv add-to-path"))
	fmt.Printf("  %s %s\n",
		style.Info.Render("2."),
		style.Input.Render("Scan for Java installations: jenv scan /usr/lib/jvm"))
	fmt.Printf("  %s %s\n",
		style.Info.Render("3."),
		style.Input.Render("Add a Java version: jenv add <name> <path>"))
	fmt.Printf("  %s %s\n",
		style.Info.Render("4."),
		style.Input.Render("Switch to a Java version: jenv use <name>"))

	if runtime.GOOS != "windows" {
		fmt.Println()
		fmt.Printf("%s: %s\n",
			style.Warning.Render("Note"),
			style.Warning.Render("You may need to restart your shell or run 'source ~/.bashrc' for changes to take effect"))
	}
}
