package sys

import (
	"os"
	"path/filepath" // Ensure filepath is imported
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsAdmin(t *testing.T) {
	// 测试Windows系统
	if runtime.GOOS == "windows" {
		// 假设当前用户不是管理员
		assert.False(t, IsAdmin())
	} else {
		// 测试Unix/Linux系统
		assert.Equal(t, os.Geteuid() == 0, IsAdmin())
	}
}

func TestRequireAdmin(t *testing.T) {
	// 测试Windows系统
	if runtime.GOOS == "windows" {
		err := RequireAdmin()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Please run this command as an administrator")
	} else {
		// 测试Unix/Linux系统
		err := RequireAdmin()
		if os.Geteuid() == 0 {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "Please run this command with sudo")
		}
	}
}

func TestIsSymlink(t *testing.T) {
	// 创建一个临时文件
	tmpFile, err := os.CreateTemp("", "test")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// 测试普通文件
	assert.False(t, IsSymlink(tmpFile.Name()))

	// 创建一个符号链接
	tmpLink := tmpFile.Name() + ".link"
	err = os.Symlink(tmpFile.Name(), tmpLink)
	assert.NoError(t, err)
	defer os.Remove(tmpLink)

	// 测试符号链接
	assert.True(t, IsSymlink(tmpLink))
}

func TestGetSymlinkTarget(t *testing.T) {
	// 创建一个临时文件
	tmpFile, err := os.CreateTemp("", "test")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// 创建一个符号链接
	tmpLink := tmpFile.Name() + ".link"
	err = os.Symlink(tmpFile.Name(), tmpLink)
	assert.NoError(t, err)
	defer os.Remove(tmpLink)

	// 测试获取符号链接目标
	target, err := GetSymlinkTarget(tmpLink)
	assert.NoError(t, err)
	assert.Equal(t, tmpFile.Name(), target)

	// 测试获取非符号链接目标
	_, err = GetSymlinkTarget(tmpFile.Name())
	assert.Error(t, err)
}

func TestCreateSymlink(t *testing.T) {
	// 创建一个临时文件
	tmpFile, err := os.CreateTemp("", "test")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// 创建一个符号链接
	tmpLink := tmpFile.Name() + ".link"
	err = CreateSymlink(tmpFile.Name(), tmpLink)
	assert.NoError(t, err)
	defer os.Remove(tmpLink)

	// 测试符号链接是否创建成功
	assert.True(t, IsSymlink(tmpLink))

	// 测试创建符号链接到已存在的符号链接
	err = CreateSymlink(tmpFile.Name(), tmpLink)
	assert.NoError(t, err)

	// 创建一个临时文件夹
	tmpDir, err := os.MkdirTemp("", "testdir")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// 创建一个指向文件夹的符号链接
	tmpDirLink := "E:\\mytest\\"
	err = CreateSymlink(tmpDir, tmpDirLink)
	assert.NoError(t, err)
	defer os.Remove(tmpDirLink)

	// 测试文件夹符号链接是否创建成功
	assert.True(t, IsSymlink(tmpDirLink))

	// 测试获取文件夹符号链接目标
	target, err := GetSymlinkTarget(tmpDirLink)
	assert.NoError(t, err)
	assert.Equal(t, tmpDir, target)
}

func TestCreateSymlink_SuccessInUserDir_Linux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Skipping Linux-specific test on non-Linux OS")
	}

	testDir, err := os.MkdirTemp("", "test_symlink_user_dir_")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(testDir)

	targetFilePath := filepath.Join(testDir, "target.txt")
	err = os.WriteFile(targetFilePath, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create target file: %v", err)
	}

	symlinkPath := filepath.Join(testDir, "symlink_to_target")

	err = CreateSymlink(targetFilePath, symlinkPath)
	if err != nil {
		t.Fatalf("CreateSymlink failed: %v", err)
	}

	// Verify symlink existence
	if _, err := os.Lstat(symlinkPath); err != nil {
		t.Fatalf("Symlink not created or not accessible: %v", err)
	}

	// Verify it's a symlink and points to the correct target
	readlinkTarget, err := os.Readlink(symlinkPath)
	if err != nil {
		t.Fatalf("Failed to read symlink (is it a symlink?): %v", err)
	}

	if readlinkTarget != targetFilePath {
		t.Errorf("Symlink points to %s, expected %s", readlinkTarget, targetFilePath)
	}
}

func TestCreateSymlink_OverwriteExistingSymlink_Linux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Skipping Linux-specific test on non-Linux OS")
	}

	testDir, err := os.MkdirTemp("", "test_symlink_overwrite_")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(testDir)

	target1Path := filepath.Join(testDir, "target1.txt")
	err = os.WriteFile(target1Path, []byte("target 1 content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create target1 file: %v", err)
	}

	target2Path := filepath.Join(testDir, "target2.txt")
	err = os.WriteFile(target2Path, []byte("target 2 content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create target2 file: %v", err)
	}

	symlinkPath := filepath.Join(testDir, "symlink_path")

	// Create initial symlink to target1
	err = os.Symlink(target1Path, symlinkPath)
	if err != nil {
		t.Fatalf("Failed to create initial symlink: %v", err)
	}

	// Call CreateSymlink to overwrite with target2
	err = CreateSymlink(target2Path, symlinkPath)
	if err != nil {
		t.Fatalf("CreateSymlink failed to overwrite: %v", err)
	}

	// Verify symlink existence
	if _, err := os.Lstat(symlinkPath); err != nil {
		t.Fatalf("Symlink not accessible after overwrite: %v", err)
	}

	// Verify it's a symlink and now points to target2
	readlinkTarget, err := os.Readlink(symlinkPath)
	if err != nil {
		t.Fatalf("Failed to read symlink after overwrite: %v", err)
	}

	if readlinkTarget != target2Path {
		t.Errorf("Symlink points to %s, expected %s after overwrite", readlinkTarget, target2Path)
	}
}

func TestCreateSymlink_CreatesParentDir_Linux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Skipping Linux-specific test on non-Linux OS")
	}

	baseTestDir, err := os.MkdirTemp("", "test_symlink_parent_")
	if err != nil {
		t.Fatalf("Failed to create base temp dir: %v", err)
	}
	defer os.RemoveAll(baseTestDir)

	targetFilePath := filepath.Join(baseTestDir, "target.txt")
	err = os.WriteFile(targetFilePath, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create target file: %v", err)
	}

	// Symlink path where 'new_parent_dir' does not exist yet
	parentDirForSymlink := filepath.Join(baseTestDir, "new_parent_dir")
	symlinkPath := filepath.Join(parentDirForSymlink, "symlink_to_target")

	err = CreateSymlink(targetFilePath, symlinkPath)
	if err != nil {
		t.Fatalf("CreateSymlink failed: %v", err)
	}

	// Verify parent directory was created
	if _, err := os.Stat(parentDirForSymlink); err != nil {
		t.Fatalf("Parent directory for symlink was not created: %v", err)
	}

	// Verify symlink existence
	if _, err := os.Lstat(symlinkPath); err != nil {
		t.Fatalf("Symlink not created or not accessible: %v", err)
	}

	// Verify it's a symlink and points to the correct target
	readlinkTarget, err := os.Readlink(symlinkPath)
	if err != nil {
		t.Fatalf("Failed to read symlink: %v", err)
	}

	if readlinkTarget != targetFilePath {
		t.Errorf("Symlink points to %s, expected %s", readlinkTarget, targetFilePath)
	}
}
