package sys

import (
	"os"
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
