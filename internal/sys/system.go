package sys

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath" // Added
	"runtime"
	"strings"       // Added
)

// IsAdmin checks if the current process has administrator privileges
func IsAdmin() bool {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("net", "session")
		err := cmd.Run()
		return err == nil
	}

	// For Unix/Linux systems, check if the user is root
	return os.Geteuid() == 0
}

// RequireAdmin checks for administrator privileges, returns an error if not present
func RequireAdmin() error {
	if !IsAdmin() {
		var msg string
		if runtime.GOOS == "windows" {
			msg = "Please run this command as an administrator"
		} else {
			msg = "Please run this command with sudo"
		}
		return fmt.Errorf("Need Admin privilege: %s", msg)
	}
	return nil
}

// IsSymlink checks if the specified path is a symbolic link
func IsSymlink(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeSymlink != 0
}

// GetSymlinkTarget retrieves the target path of a symbolic link
func GetSymlinkTarget(path string) (string, error) {
	target, err := os.Readlink(path)
	if err != nil {
		return "", fmt.Errorf("Failed to read symbolic link: %v", err)
	}
	return target, nil
}

/*
*

	CreateSymlink creates a symbolic link, requires administrator privileges

@param oldPath The path to the existing file or directory
@param newPath The path where the symbolic link should be created
*/
func CreateSymlink(oldPath string, newPath string) error {
	// Ensure the parent directory for newPath exists.
	parentDir := filepath.Dir(newPath)
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		if mkdirErr := os.MkdirAll(parentDir, 0755); mkdirErr != nil {
			return fmt.Errorf("failed to create parent directory %s: %w", parentDir, mkdirErr)
		}
	}

	// If newPath (the symlink itself) already exists, remove it first.
	// os.Lstat is used to check the symlink itself, not its target.
	if _, err := os.Lstat(newPath); err == nil {
		if removeErr := os.Remove(newPath); removeErr != nil {
			return fmt.Errorf("failed to remove existing file/symlink at %s: %w", newPath, removeErr)
		}
	} else if !os.IsNotExist(err) {
		// For errors other than "not exist", return them.
		return fmt.Errorf("failed to check status of %s: %w", newPath, err)
	}

	// Attempt to create the symbolic link.
	err := os.Symlink(oldPath, newPath)
	if err != nil {
		if os.IsPermission(err) && runtime.GOOS == "linux" {
			// Check if newPath is in a system-protected area vs user's home.
			// This is a heuristic. A more robust check might involve comparing newPath against known system paths.
			homeDir, homeErr := os.UserHomeDir()
			if homeErr == nil && strings.HasPrefix(newPath, homeDir) {
				return fmt.Errorf("failed to create symlink due to permissions at %s (even in home directory): %w. Please check directory permissions", newPath, err)
			}
			return fmt.Errorf("failed to create symlink at %s due to permissions: %w. If this path is outside your home directory, you might need to use sudo", newPath, err)
		}
		return fmt.Errorf("failed to create symbolic link from %s to %s: %w", oldPath, newPath, err)
	}

	return nil
}
