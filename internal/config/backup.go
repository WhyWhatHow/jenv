package config

import (
	"encoding/json"
	"github.com/whywhathow/jenv/internal/constants"
	"github.com/whywhathow/jenv/internal/env"
	"os"
	"path/filepath"
	"runtime"
)

// PathBackup represents the structure for backing up PATH environment variables
type PathBackup struct {
	UserPath   string `json:"user_path,omitempty"`
	SystemPath string `json:"system_path,omitempty"`
}

// BackupEnvPath creates a backup of the PATH environment variables
// On Windows, it backs up both user and system PATH variables
func BackupEnvPath() error {
	// Get user home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Create .jdks directory if it doesn't exist
	jdksDir := filepath.Join(home, constants.DEFAULT_FOLDER)
	if err := os.MkdirAll(jdksDir, 0755); err != nil {
		return err
	}

	// Check if backup file already exists
	backupPath := filepath.Join(jdksDir, constants.DEFAULT_BACKUP_FILE)
	if _, err := os.Stat(backupPath); err == nil {
		// Backup file already exists, no need to create it again
		return nil
	}

	// Create backup structure
	backup := PathBackup{}

	// Get PATH information based on OS
	if runtime.GOOS == "windows" {
		// Get user PATH from registry
		userPath, err := env.QueryUserEnvironmentVariable("PATH")
		if err == nil {
			backup.UserPath = userPath
		}
		// Get system PATH from registry
		sysPath, err := env.QuerySystemEnvironmentVariable("PATH")
		if err == nil {
			backup.SystemPath = sysPath
		}
	} else {
		// For non-Windows systems, just get the PATH environment variable
		backup.UserPath = os.Getenv("PATH")
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(backup, "", "  ")
	if err != nil {
		return err
	}

	// Write to backup.json file
	return os.WriteFile(backupPath, data, 0644)
}

// RestorePathFromBackup restores PATH environment variables from backup
func RestorePathFromBackup() error {
	// Get user home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Check if backup file exists
	backupPath := filepath.Join(home, ".jdks", constants.DEFAULT_BACKUP_FILE)
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return err
	}

	// Unmarshal backup data
	var backup PathBackup
	if err := json.Unmarshal(data, &backup); err != nil {
		return err
	}

	// Restore PATH based on OS
	if runtime.GOOS == "windows" {
		// Only attempt to restore if we have admin privileges
		// Restore user PATH
		if backup.UserPath != "" {
			env.UpdateUserEnvironmentVariable("PATH", backup.UserPath)
		}

		// Restore system PATH
		if backup.SystemPath != "" {
			env.UpdateSystemEnvironmentVariable("PATH", backup.SystemPath)
		}
	} else {
		// For non-Windows systems, just set the PATH environment variable
		if backup.UserPath != "" {
			os.Setenv("PATH", backup.UserPath)
		}
	}

	return nil
}

// userPath, or systemPath update
//
//	func UpdateBackUp(key string, value string) {
//		backup, err := getBackupPath()
//		if err == null {
//			 if(key=="userpath")
//		}
//
// }
func GetDefaultBackupFilePath() string {
	// Get user home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	// Check if backup file exists
	backupPath := filepath.Join(home, ".jdks", constants.DEFAULT_BACKUP_FILE)
	return backupPath
}
