package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/whywhathow/jenv/internal/config"
	"github.com/whywhathow/jenv/internal/style"
	"github.com/whywhathow/jenv/internal/sys"
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

	// Check admin privileges - different logic for different platforms
	if runtime.GOOS == "windows" {
		// Windows always requires administrator privileges
		if !sys.IsAdmin() {
			fmt.Printf("%s: %s\n",
				style.Error.Render("Error"),
				style.Error.Render("Administrator privileges required"))
			os.Exit(1)
		}
	} else {
		// Linux/Unix: Only warn if not root, don't exit
		if !sys.IsAdmin() {
			fmt.Printf("%s: %s\n",
				style.Warning.Render("Warning"),
				style.Warning.Render("Running without root privileges. Some features may be limited."))
			fmt.Printf("%s: %s\n",
				style.Info.Render("Info"),
				style.Info.Render("jenv will use user-level configuration and symlinks."))
		}
	}

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

	// Check if jenv has been initialized
	if !cfg.Initialized {
		fmt.Printf("%s: %s\n",
			style.Warning.Render("Warning"),
			style.Warning.Render("jenv has not been initialized yet"))
		fmt.Printf("%s: %s\n",
			style.Info.Render("Info"),
			style.Info.Render("Run 'jenv init' to set up jenv for first-time use"))
		fmt.Println()
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
