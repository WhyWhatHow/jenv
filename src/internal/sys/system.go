package sys

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
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

//// IsSymlink checks if the specified path is a symbolic link
//func IsSymlink(path string) bool {
//	info, err := os.Lstat(path)
//	if err != nil {
//		return false
//	}
//	return info.Mode()&os.ModeSymlink != 0
//}

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
	if err := RequireAdmin(); err != nil {
		return err
	}

	// If the target path exists , remove it first
	if _, err := os.Stat(newPath); err == nil {
		if err := os.Remove(newPath); err != nil {
			return fmt.Errorf("Failed to remove existing symbolic link: %v", err)
		}
	}

	// Create the symbolic link
	if err := os.Symlink(oldPath, newPath); err != nil {
		return fmt.Errorf("Failed to create symbolic link: %v", err)
	}

	return nil
}
