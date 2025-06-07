package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/whywhathow/jenv/internal/config"
	"github.com/whywhathow/jenv/internal/java"
	"github.com/whywhathow/jenv/internal/style"
	// "github.com/whywhathow/jenv/internal/sys" // sys.IsAdmin check removed
	"os"
)

var Version = "dev"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jenv",
	Short: "A Java version manager",
	Long: `Jenv is a command-line tool for managing multiple Java versions.

It allows you to easily switch between different Java versions,
add new Java installations, and manage your Java environment.

First time setup:
➜ Run 'jenv add-to-path' to add jenv to your system PATH
➜ Run 'jenv scan <dir>' to find and add Java installations
➜ Run 'jenv use <name>' to select a Java version`,
	Version: Version,
}

func init() {
	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "version %s" .Version}}
Author: WhyWhatHow (https://github.com/WhyWhatHow)
Email: whywhathow.fun@gmail.com
License: Apache License 2.0
`)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

/*
*
1. Initialize configuration
2. Backup environment variables if not already done
3. Add global flags
4.
*/
func init() {
	// Check admin privileges - This global check is removed.
	// Specific commands or functions should handle permissions as needed.
	// if !sys.IsAdmin() {
	// 	fmt.Printf("%s: %s\n",
	// 		style.Error.Render("Error"),
	// 		style.Error.Render("Administrator/root privileges required"))
	// 	os.Exit(1)
	// }

	// Initialize configuration system
	cfg, err := config.GetInstance()
	if err != nil {
		fmt.Printf("%s: %s\n%s\n",
			style.Error.Render("Error"),
			style.Error.Render("Failed to initialize configuration"),
			style.Error.Render(err.Error()))
		os.Exit(1)
	}

	// Apply saved theme if exists
	if cfg.Theme != "" {
		if theme, ok := style.GetThemeByName(cfg.Theme); ok {
			style.ApplyTheme(theme)
		}
	}

	// First-time initialization
	if !cfg.Initialized {
		fmt.Println(style.Info.Render("➜ Performing first-time setup..."))

		// Backup environment variables
		fmt.Print(style.Info.Render("➜ Backing up environment variables... "))
		if err := config.BackupEnvPath(); err != nil {
			fmt.Printf("\n%s: %s\n%s\n",
				style.Error.Render("Error"),
				style.Error.Render("Failed to backup environment variables"),
				style.Error.Render(err.Error()))
			os.Exit(1)
		}
		fmt.Println(style.Success.Render("✓"))

		// Initialize Java environment  && add JAVA_HOME to PATH
		fmt.Print(style.Info.Render("➜ Initializing Java environment... "))
		if err := java.Init(); err != nil {
			fmt.Printf("\n%s: %s\n%s\n",
				style.Error.Render("Error"),
				style.Error.Render("Failed to initialize Java environment"),
				style.Error.Render(err.Error()))
			os.Exit(1)
		}
		fmt.Println(style.Success.Render("✓"))

		// Save configuration
		cfg.Initialized = true
		if err := cfg.Save(); err != nil {
			fmt.Printf("%s: %s\n%s\n",
				style.Error.Render("Error"),
				style.Error.Render("Failed to save configuration"),
				style.Error.Render(err.Error()))
			os.Exit(1)
		}

		fmt.Printf("\n%s\n", style.Success.Render("✓ Initialization complete!"))
		fmt.Printf("%s\n", style.Info.Render("➜ Run 'jenv add-to-path' to complete the setup"))
	}
}
