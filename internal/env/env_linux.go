//go:build linux
// +build linux

package env

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/unix" // For checking write access
)

// Helper function to check write access to a directory
func canWriteToDir(dirPath string) bool {
	return unix.Access(dirPath, unix.W_OK) == nil
}

// Helper function to update or append a key-value pair in a script file
// Ensures the line is in the format `export KEY="VALUE"`
func updateScriptFile(filePath, key, value string) error {
	// Ensure directory exists if it's a user file path
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Attempt to create directory if it doesn't exist (e.g. for ~/.config/jenv/)
		// For /etc/profile.d, this won't create /etc or /etc/profile.d
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}


	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open or create file %s: %w", filePath, err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	keyFound := false
	// Lines to search for: `export KEY=`, `KEY=` (at the beginning of a line or after spaces)
	// We will replace them with `export KEY="VALUE"`
	prefixDashExport := fmt.Sprintf("export %s=", key) // e.g. export JAVA_HOME=
	prefixKeyOnly := fmt.Sprintf("%s=", key)           // e.g. JAVA_HOME=
	targetLine := fmt.Sprintf("export %s=\"%s\"", key, value)

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, prefixDashExport) || strings.HasPrefix(trimmedLine, prefixKeyOnly) {
			// Check if it's part of a comment
			if strings.HasPrefix(trimmedLine, "#") {
				lines = append(lines, line) // Keep commented line as is
				continue
			}
			if !keyFound { // Replace only the first occurrence if multiple are somehow present
				lines = append(lines, targetLine)
				keyFound = true
			} else {
				// If key was already found and replaced, and we find another declaration,
				// it's safer to keep it to avoid accidental data loss from complex files.
				// Or, we could choose to remove subsequent declarations. For now, keep.
				lines = append(lines, line)
			}
		} else {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	if !keyFound {
		// If the file is not empty and doesn't end with a newline, add one.
		if len(lines) > 0 && lines[len(lines)-1] != "" {
			lines = append(lines, "")
		}
		lines = append(lines, targetLine)
	}

	// Truncate file before writing new content
	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate file %s: %w", filePath, err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek file %s: %w", filePath, err)
	}

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		if _, err := fmt.Fprintln(writer, line); err != nil {
			return fmt.Errorf("failed to write to file %s: %w", filePath, err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer for %s: %w", filePath, err)
	}
	return nil
}

// Stubs for Windows-specific functions to allow Linux build.
// These should not be called on Linux.
// They are defined here because env.go calls them, and env_windows.go
// might not be providing them correctly for all build/link scenarios.

func SetEnvInWin(key string, value string) error {
	return fmt.Errorf("SetEnvInWin called on non-Windows OS")
}

func QuerySystemEnvironmentVariable(variable string) (string, error) {
	return "", fmt.Errorf("QuerySystemEnvironmentVariable called on non-Windows OS (" + variable + ")")
}

func UpdateSystemEnvironmentVariable(key string, value string) error {
	return fmt.Errorf("UpdateSystemEnvironmentVariable called on non-Windows OS")
}

func QueryUserEnvironmentVariable(variable string) (string, error) {
	return "", fmt.Errorf("QueryUserEnvironmentVariable called on non-Windows OS (%s)", variable)
}

func UpdateUserEnvironmentVariable(key string, value string) error {
	return fmt.Errorf("UpdateUserEnvironmentVariable called on non-Windows OS")
}

// SetSystemPath permanently sets the system PATH environment variable (Linux platform implementation)
func SetSystemPath(newPathEntry string) error {
	currentPath := os.Getenv("PATH")
	// Prepend the new path entry
	var newPath string
	if currentPath == "" {
		newPath = newPathEntry
	} else {
		// Avoid adding duplicate paths
		paths := strings.Split(currentPath, string(filepath.ListSeparator))
		for _, p := range paths {
			if p == newPathEntry {
				// Path already exists, do nothing
				return nil
			}
		}
		newPath = newPathEntry + string(filepath.ListSeparator) + currentPath
	}
	return UpdateEnvironmentVariable("PATH", newPath)
}

// UpdateEnvironmentVariable updates Linux system or user environment variables.
// It tries to set the variable system-wide in /etc/profile.d/jenv.sh first.
// If that fails, it falls back to user-specific shell configuration files.
func UpdateEnvironmentVariable(key, value string) error {
	// Set for current process immediately
	if err := os.Setenv(key, value); err != nil {
		return fmt.Errorf("failed to set environment variable for current process: %w", err)
	}

	jenvShPath := "/etc/profile.d/jenv.sh"
	if canWriteToDir("/etc/profile.d") {
		err := updateScriptFile(jenvShPath, key, value)
		if err == nil {
			// Successfully updated system-wide script
			return nil
		}
		// If updateScriptFile failed despite directory being writable, log and fallback
		fmt.Fprintf(os.Stderr, "Warning: Failed to write to %s: %v. Attempting user-level configuration.\n", jenvShPath, err)
	} else {
		fmt.Fprintf(os.Stderr, "Warning: No write access to /etc/profile.d. Attempting user-level configuration for %s.\n", key)
	}

	// Fallback to user-specific configuration
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	shell := os.Getenv("SHELL")
	var configFile string

	if strings.Contains(shell, "bash") {
		configFile = filepath.Join(homeDir, ".bashrc")
	} else if strings.Contains(shell, "zsh") {
		configFile = filepath.Join(homeDir, ".zshrc")
	} else {
		// General fallback for other shells (sh, dash, etc.)
		configFile = filepath.Join(homeDir, ".profile")
	}

	fmt.Fprintf(os.Stdout, "Attempting to update user configuration: %s\n", configFile)
	err = updateScriptFile(configFile, key, value)
	if err != nil {
		return fmt.Errorf("failed to update user environment file %s: %w", configFile, err)
	}

	return nil
}

// updateEnvironmentFile is deprecated. Use UpdateEnvironmentVariable.
// This function's logic for /etc/environment is being replaced by /etc/profile.d/jenv.sh usage.
func updateEnvironmentFile(filePath, key, value string) error {
	return fmt.Errorf("updateEnvironmentFile is deprecated; /etc/environment is no longer managed by this tool directly")
}

// updateProfileFile is deprecated. Use UpdateEnvironmentVariable.
// This function's logic for user .profile is being replaced by more specific shell config files.
func updateProfileFile(filePath, key, value string) error {
	return fmt.Errorf("updateProfileFile is deprecated; user profile updates are handled by UpdateEnvironmentVariable")
}


// SetPersistentEnvVar 在 Unix-like 系统上永久设置用户级环境变量
// This function might be redundant given the new UpdateEnvironmentVariable.
// For now, it's kept as per instructions, but its direct Fprintf usage should be reviewed.
func SetPersistentEnvVar(name, value string) error {
	// 获取用户主目录
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory for SetPersistentEnvVar: %w", err)
	}

	// 目标文件为 ~/.profile - this is a specific choice of this function
	profilePath := filepath.Join(home, ".profile")

	// The new updateScriptFile function handles creation and update more robustly.
	// However, SetPersistentEnvVar's original behavior was simple append.
	// To maintain its distinct behavior for now (as per instructions to keep it):
	f, err := os.OpenFile(profilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open/create %s for SetPersistentEnvVar: %w", profilePath, err)
	}
	defer f.Close()

	// Ensure newline before and after the export statement
	if _, err := fmt.Fprintf(f, "\nexport %s=\"%s\"\n", name, value); err != nil {
		return fmt.Errorf("failed to write to %s for SetPersistentEnvVar: %w", profilePath, err)
	}
	return nil
}
