package style

import "github.com/charmbracelet/lipgloss"

// Style variables for consistent UI rendering
var (
	Header  = lipgloss.NewStyle()
	Name    = lipgloss.NewStyle()
	Path    = lipgloss.NewStyle()
	Current = lipgloss.NewStyle()
	Error   = lipgloss.NewStyle()
	Warning = lipgloss.NewStyle()
	Success = lipgloss.NewStyle()
	Input   = lipgloss.NewStyle()
	Info    = lipgloss.NewStyle()
)

func init() {
	// Apply default theme on package initialization
	ApplyTheme(DefaultTheme)
}
