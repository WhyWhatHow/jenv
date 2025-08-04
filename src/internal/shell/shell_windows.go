//go:build windows

package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ShellType represents different shell types on Windows
type ShellType string

const (
	PowerShell ShellType = "powershell"
	Cmd        ShellType = "cmd"
	GitBash    ShellType = "gitbash"
)

// ShellConfig represents shell-specific configuration for Windows
type ShellConfig struct {
	Type        ShellType
	ConfigFile  string
	ExportCmd   string
	SourceCmd   string
	CommentChar string
}

// GetShellConfigs returns configurations for different shells on Windows
func GetShellConfigs() map[ShellType]ShellConfig {
	return map[ShellType]ShellConfig{
		PowerShell: {
			Type:        PowerShell,
			ConfigFile:  "Documents/PowerShell/Microsoft.PowerShell_profile.ps1",
			ExportCmd:   "$env:",
			SourceCmd:   ".",
			CommentChar: "#",
		},
		GitBash: {
			Type:        GitBash,
			ConfigFile:  ".bashrc",
			ExportCmd:   "export",
			SourceCmd:   "source",
			CommentChar: "#",
		},
	}
}

// DetectUserShells detects which shells the user has configured on Windows
func DetectUserShells() ([]ShellType, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %v", err)
	}

	var shells []ShellType
	configs := GetShellConfigs()

	for shellType, config := range configs {
		configPath := filepath.Join(homeDir, config.ConfigFile)
		
		// Check if config file exists or if we can create it
		if shellType == PowerShell {
			configDir := filepath.Dir(configPath)
			if _, err := os.Stat(configDir); err == nil {
				shells = append(shells, shellType)
			}
		} else {
			if _, err := os.Stat(configPath); err == nil {
				shells = append(shells, shellType)
			}
		}
	}

	// On Windows, we primarily use registry for environment variables
	// Shell configs are secondary
	return shells, nil
}

// GetCurrentShell attempts to detect the current shell on Windows
func GetCurrentShell() ShellType {
	// On Windows, we can't easily detect the current shell
	// Default to PowerShell as it's most common
	return PowerShell
}

// SetEnvironmentVariable sets an environment variable for Windows shells
// Note: On Windows, environment variables are primarily managed through registry
// This function is mainly for shell-specific configurations
func SetEnvironmentVariable(key, value string) error {
	shells, err := DetectUserShells()
	if err != nil {
		return fmt.Errorf("failed to detect user shells: %v", err)
	}

	// On Windows, we don't need to update shell configs for environment variables
	// as they are handled by the registry. This is mainly for compatibility.
	var errors []string
	for _, shell := range shells {
		if err := SetEnvironmentVariableForShell(shell, key, value); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", shell, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to update some shell environments: %s", strings.Join(errors, "; "))
	}

	return nil
}

// SetEnvironmentVariableForShell sets an environment variable for a specific shell on Windows
func SetEnvironmentVariableForShell(shellType ShellType, key, value string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %v", err)
	}

	configs := GetShellConfigs()
	config, exists := configs[shellType]
	if !exists {
		return fmt.Errorf("unsupported shell type: %s", shellType)
	}

	configPath := filepath.Join(homeDir, config.ConfigFile)

	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	return updateWindowsShellConfigFile(configPath, config, key, value)
}

// updateWindowsShellConfigFile updates a shell configuration file on Windows
func updateWindowsShellConfigFile(filePath string, config ShellConfig, key, value string) error {
	// Ensure file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if file, err := os.Create(filePath); err != nil {
			return fmt.Errorf("failed to create config file %s: %v", filePath, err)
		} else {
			file.Close()
		}
	}

	// Read existing content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %v", filePath, err)
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	keyFound := false

	// Generate the export line based on shell type
	var exportLine string
	var exportPrefix string

	if config.Type == PowerShell {
		exportLine = fmt.Sprintf("%s%s = \"%s\"", config.ExportCmd, key, value)
		exportPrefix = fmt.Sprintf("%s%s = ", config.ExportCmd, key)
	} else {
		exportLine = fmt.Sprintf("%s %s=\"%s\"", config.ExportCmd, key, value)
		exportPrefix = fmt.Sprintf("%s %s=", config.ExportCmd, key)
	}

	// Find and update existing environment variable setting
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, exportPrefix) {
			newLines = append(newLines, exportLine)
			keyFound = true
		} else {
			newLines = append(newLines, line)
		}
	}

	// If not found, add new setting
	if !keyFound {
		// Add comment and new export line
		newLines = append(newLines, "")
		newLines = append(newLines, fmt.Sprintf("%s Added by jenv", config.CommentChar))
		newLines = append(newLines, exportLine)
	}

	// Write back to file
	return os.WriteFile(filePath, []byte(strings.Join(newLines, "\n")), 0644)
}

// RemoveEnvironmentVariable removes an environment variable from all shell configs on Windows
func RemoveEnvironmentVariable(key string) error {
	shells, err := DetectUserShells()
	if err != nil {
		return fmt.Errorf("failed to detect user shells: %v", err)
	}

	var errors []string
	for _, shell := range shells {
		if err := RemoveEnvironmentVariableFromShell(shell, key); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", shell, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to remove from some shell environments: %s", strings.Join(errors, "; "))
	}

	return nil
}

// RemoveEnvironmentVariableFromShell removes an environment variable from a specific shell on Windows
func RemoveEnvironmentVariableFromShell(shellType ShellType, key string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %v", err)
	}

	configs := GetShellConfigs()
	config, exists := configs[shellType]
	if !exists {
		return fmt.Errorf("unsupported shell type: %s", shellType)
	}

	configPath := filepath.Join(homeDir, config.ConfigFile)

	// If file doesn't exist, nothing to remove
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil
	}

	// Read existing content
	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %v", configPath, err)
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string

	// Generate the export prefix based on shell type
	var exportPrefix string
	if config.Type == PowerShell {
		exportPrefix = fmt.Sprintf("%s%s = ", config.ExportCmd, key)
	} else {
		exportPrefix = fmt.Sprintf("%s %s=", config.ExportCmd, key)
	}

	// Filter out lines that set this environment variable
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmedLine, exportPrefix) {
			newLines = append(newLines, line)
		}
	}

	// Write back to file
	return os.WriteFile(configPath, []byte(strings.Join(newLines, "\n")), 0644)
}
