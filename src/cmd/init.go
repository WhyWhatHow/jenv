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

Platform-specific behavior:
  Windows: Uses Windows registry for environment variable management
           Requires administrator privileges for system-wide installation

  Linux/macOS: Uses shell configuration files for environment variables
               Can run with or without root privileges
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

	// Check platform-specific requirements and provide clear guidance
	if runtime.GOOS == "windows" {
		fmt.Printf("%s: %s\n",
			style.Info.Render("Platform"),
			style.Info.Render("Windows - Using registry-based environment variable management"))

		if !sys.IsAdmin() {
			fmt.Printf("%s: %s\n",
				style.Error.Render("Error"),
				style.Error.Render("Administrator privileges required on Windows"))
			fmt.Printf("%s: %s\n",
				style.Info.Render("Solution"),
				style.Info.Render("Please run PowerShell as Administrator and try again"))
			return
		}

		fmt.Printf("%s: %s\n",
			style.Success.Render("Privileges"),
			style.Success.Render("Administrator privileges detected - proceeding with system-wide setup"))
	} else {
		fmt.Printf("%s: %s\n",
			style.Info.Render("Platform"),
			style.Info.Render("Unix-like - Using shell configuration files for environment variables"))

		if !sys.IsAdmin() {
			fmt.Printf("%s: %s\n",
				style.Warning.Render("Privileges"),
				style.Warning.Render("Running without root privileges"))
			fmt.Printf("%s: %s\n",
				style.Info.Render("Configuration"),
				style.Info.Render("Will create user-level configuration in your home directory"))
		} else {
			fmt.Printf("%s: %s\n",
				style.Success.Render("Privileges"),
				style.Success.Render("Root privileges detected - will create system-wide configuration"))
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

	// Provide platform-specific next steps
	fmt.Println()
	fmt.Println(style.Header.Render("📋 Next Steps:"))

	if runtime.GOOS == "windows" {
		fmt.Printf("  %s %s\n",
			style.Info.Render("1."),
			style.Input.Render("Add jenv to your PATH: jenv add-to-path"))
		fmt.Printf("  %s %s\n",
			style.Info.Render("2."),
			style.Input.Render("Scan for Java installations: jenv scan c:\\"))
		fmt.Printf("  %s %s\n",
			style.Info.Render("3."),
			style.Input.Render("Add a Java version: jenv add <name> <path>"))
		fmt.Printf("  %s %s\n",
			style.Info.Render("4."),
			style.Input.Render("Switch to a Java version: jenv use <name>"))

		fmt.Println()
		fmt.Printf("%s: %s\n",
			style.Info.Render("Windows Note"),
			style.Info.Render("Environment variables are managed through Windows registry"))
		fmt.Printf("%s: %s\n",
			style.Info.Render("Alternative"),
			style.Info.Render("You can also manage environment variables via Control Panel → System → Advanced"))
	} else {
		fmt.Printf("  %s %s\n",
			style.Info.Render("1."),
			style.Input.Render("Add jenv to your PATH: jenv add-to-path"))
		fmt.Printf("  %s %s\n",
			style.Info.Render("2."),
			style.Input.Render("Reload your shell configuration:"))
		fmt.Printf("     %s\n", style.Input.Render("source ~/.bashrc    # for bash"))
		fmt.Printf("     %s\n", style.Input.Render("source ~/.zshrc     # for zsh"))
		fmt.Printf("     %s\n", style.Input.Render("source ~/.config/fish/config.fish  # for fish"))
		fmt.Printf("  %s %s\n",
			style.Info.Render("3."),
			style.Input.Render("Scan for Java installations: jenv scan /usr/lib/jvm"))
		fmt.Printf("  %s %s\n",
			style.Info.Render("4."),
			style.Input.Render("Add a Java version: jenv add <name> <path>"))
		fmt.Printf("  %s %s\n",
			style.Info.Render("5."),
			style.Input.Render("Switch to a Java version: jenv use <name>"))

		fmt.Println()
		fmt.Printf("%s: %s\n",
			style.Info.Render("Unix Note"),
			style.Info.Render("Environment variables are managed through shell configuration files"))
	}
}
