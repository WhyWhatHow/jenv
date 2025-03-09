package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBackupPath(t *testing.T) {
	// 测试 BackupEnvPath 函数
	err := BackupEnvPath()
	if err != nil {
		t.Fatalf("BackupEnvPath failed: %v", err)
	}

	// 检查备份文件是否存在
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}
	backupPath := filepath.Join(home, ".jdks", "backup.json")
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		t.Fatalf("Backup file does not exist: %v", err)
	}
}

func TestRestorePathFromBackup(t *testing.T) {
	// 先备份当前 PATH 环境变量
	err := BackupEnvPath()
	if err != nil {
		t.Fatalf("BackupEnvPath failed: %v", err)
	}

	// 测试 RestorePathFromBackup 函数
	err = RestorePathFromBackup()
	if err != nil {
		t.Fatalf("RestorePathFromBackup failed: %v", err)
	}

	// 检查 PATH 是否被正确恢复
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}
	backupPath := filepath.Join(home, ".jdks", "backup.json")
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		t.Fatalf("Backup file does not exist: %v", err)
	}
}
