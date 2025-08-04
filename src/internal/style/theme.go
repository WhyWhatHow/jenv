package style

import "github.com/charmbracelet/lipgloss"

// Theme represents a collection of color settings for the application
type Theme struct {
	Name         string
	HeaderColor  string
	NameColor    string
	PathColor    string
	CurrentColor string
	ErrorColor   string
	WarningColor string
	SuccessColor string
	InputColor   string
}

// Available themes
var (
	DefaultTheme = Theme{
		Name:         "default",
		HeaderColor:  "87",  // Bright cyan for headers
		NameColor:    "147", // Light purple for names
		PathColor:    "246", // Light gray for paths
		CurrentColor: "159", // Soft cyan for current items
		ErrorColor:   "203", // Soft red for errors
		WarningColor: "214", // Orange for warnings
		SuccessColor: "150", // Soft green for success
		InputColor:   "153", // Soft blue for input
	}

	DarkTheme = Theme{
		Name:         "dark",
		HeaderColor:  "81",  // Deep cyan for headers
		NameColor:    "141", // Deep purple for names
		PathColor:    "244", // Darker gray for paths
		CurrentColor: "123", // Deep cyan for current items
		ErrorColor:   "196", // Deep red for errors
		WarningColor: "208", // Dark orange for warnings
		SuccessColor: "114", // Deep green for success
		InputColor:   "147", // Deep blue for input
	}

	LightTheme = Theme{
		Name:         "light",
		HeaderColor:  "31",  // Light cyan for headers
		NameColor:    "98",  // Light purple for names
		PathColor:    "242", // Medium gray for paths
		CurrentColor: "37",  // Bright cyan for current items
		ErrorColor:   "167", // Muted red for errors
		WarningColor: "172", // Light orange for warnings
		SuccessColor: "71",  // Muted green for success
		InputColor:   "67",  // Muted blue for input
	}
)

// CurrentTheme holds the currently active theme
var CurrentTheme = DefaultTheme

// ApplyTheme applies the specified theme to all styles
func ApplyTheme(theme Theme) {
	Header = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(theme.HeaderColor))

	Name = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.NameColor))

	Path = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.PathColor))

	Current = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.CurrentColor)).
		Bold(true)

	Error = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.ErrorColor)).
		Bold(true)

	Warning = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.WarningColor)).
		Bold(true)

	Success = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.SuccessColor)).
		Bold(true)

	Input = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.InputColor)).
		Italic(true)

	CurrentTheme = theme
}

// GetAvailableThemes returns a list of all available themes
func GetAvailableThemes() []Theme {
	return []Theme{DefaultTheme, DarkTheme, LightTheme}
}

// GetThemeByName returns a theme by its name
func GetThemeByName(name string) (Theme, bool) {
	for _, theme := range GetAvailableThemes() {
		if theme.Name == name {
			return theme, true
		}
	}
	return Theme{}, false
}
