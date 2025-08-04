package sys

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
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
	// 在Linux/Unix系统中，只有在创建系统级符号链接时才需要root权限
	// 用户级符号链接不需要特殊权限
	if runtime.GOOS != "windows" {
		// 检查是否为系统级路径（如/opt, /usr等）
		if isSystemPath(newPath) && os.Geteuid() != 0 {
			return fmt.Errorf("creating symlink in system directory requires root privileges: %s", newPath)
		}
	} else {
		// Windows系统仍然需要管理员权限
		if err := RequireAdmin(); err != nil {
			return err
		}
	}

	// 确保父目录存在
	if err := os.MkdirAll(filepath.Dir(newPath), 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %v", err)
	}

	// If the target path exists, remove it first
	if _, err := os.Lstat(newPath); err == nil {
		if err := os.Remove(newPath); err != nil {
			return fmt.Errorf("failed to remove existing path: %v", err)
		}
	}

	// Create the symbolic link
	if err := os.Symlink(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to create symbolic link: %v", err)
	}

	return nil
}

// isSystemPath checks if a path is in a system directory that requires root privileges
func isSystemPath(path string) bool {
	if runtime.GOOS == "windows" {
		return false // Windows handles this differently
	}

	// Common system directories that require root privileges
	systemPaths := []string{
		"/opt/",
		"/usr/",
		"/etc/",
		"/var/",
		"/bin/",
		"/sbin/",
		"/lib/",
		"/lib64/",
	}

	for _, sysPath := range systemPaths {
		if strings.HasPrefix(path, sysPath) {
			return true
		}
	}

	return false
}
