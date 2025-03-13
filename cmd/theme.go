package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/whywhathow/jenv/internal/config"
	"github.com/whywhathow/jenv/internal/style"
	"strings"
)

var themeCmd = &cobra.Command{
	Aliases: []string{"t"},
	Use:     "theme [name]",
	Short:   "Manage color themes",
	Long: `Manage color themes for the CLI interface.

Without arguments, this command lists all available themes.
With a theme name argument, it switches to the specified theme.`,
	Example: "  jenv theme\n  jenv theme dark",
	Run:     runTheme,
}

func init() {
	rootCmd.AddCommand(themeCmd)
}

func runTheme(cmd *cobra.Command, args []string) {
	// If no arguments provided, list available themes
	if len(args) == 0 {
		listThemes()
		return
	}

	// Switch to the specified theme
	themeName := args[0]
	if theme, ok := style.GetThemeByName(themeName); ok {
		style.ApplyTheme(theme)
		// Save theme to configuration
		cfg, err := config.GetInstance()
		if err != nil {
			fmt.Printf("%s: %s\n",
				style.Error.Render("Error"),
				"Failed to get configuration instance")
			return
		}
		cfg.Theme = themeName
		if err := cfg.Save(); err != nil {
			fmt.Printf("%s: %s\n",
				style.Error.Render("Error"),
				"Failed to save theme configuration")
			return
		}
		fmt.Printf("%s: %s\n",
			style.Success.Render("Successfully switched to theme"),
			style.Name.Render(themeName))
	} else {
		fmt.Printf("%s: Theme '%s' not found\n",
			style.Error.Render("Error"),
			themeName)
		listThemes()
	}
}

func listThemes() {
	fmt.Println(style.Header.Render("Available Themes"))
	fmt.Println(strings.Repeat("─", 30))

	for _, theme := range style.GetAvailableThemes() {
		if theme.Name == style.CurrentTheme.Name {
			fmt.Printf("%s %s\n",
				style.Current.Render("✓"),
				style.Name.Render(theme.Name))
		} else {
			fmt.Printf("  %s\n",
				style.Name.Render(theme.Name))
		}
	}
}
