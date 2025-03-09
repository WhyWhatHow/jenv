package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSetSystemPath(t *testing.T) {
	// 保存当前的 PATH 环境变量
	originalPath := os.Getenv("PATH")
	defer func() {
		// 测试结束后恢复原始 PATH
		err := os.Setenv("PATH", originalPath)
		if err != nil {
			t.Errorf("恢复原始 PATH 失败: %v", err)
		}
	}()

	// 测试路径
	testPath := "C:\\Test\\Path"

	// 获取当前的 PATH
	currentPath := os.Getenv("PATH")

	// 添加新路径到 PATH
	newPath := testPath
	if currentPath != "" {
		newPath = testPath + string(filepath.ListSeparator) + currentPath
	}

	// 设置新的 PATH
	err := os.Setenv("PATH", newPath)
	if err != nil {
		t.Fatalf("设置 PATH 环境变量失败: %v", err)
	}

	// 验证 PATH 是否被正确设置
	updatedPath := os.Getenv("PATH")
	if !strings.Contains(updatedPath, testPath) {
		t.Errorf("PATH 环境变量未包含测试路径\n期望包含: %s\n实际值: %s", testPath, updatedPath)
	}

	// 验证分隔符是否正确
	paths := strings.Split(updatedPath, string(filepath.ListSeparator))
	if paths[0] != testPath {
		t.Errorf("新路径未被添加到 PATH 的开头\n期望: %s\n实际值: %s", testPath, paths[0])
	}
}

func TestUpdateEnvironmentVariable(t *testing.T) {
	// 保存原始环境变量
	testKey := "JENV_TEST_VAR"
	defer func() {
		// 测试结束后清理环境变量
		err := os.Unsetenv(testKey)
		if err != nil {
			t.Errorf("清理测试环境变量失败: %v", err)
		}
	}()

	// 测试设置环境变量
	testValue := "test_value"
	err := UpdateSystemEnvironmentVariable(testKey, testValue)

	// 由于需要管理员权限，这里我们期望在非管理员权限下运行时返回错误
	if err == nil {
		// 如果没有错误，说明可能是在管理员权限下运行
		// 验证环境变量是否被正确设置
		actualValue := os.Getenv(testKey)
		if actualValue != testValue {
			t.Errorf("环境变量值设置错误，期望: %s, 实际: %s", testValue, actualValue)
		}
	} else {
		// 在非管理员权限下，应该返回权限相关的错误
		if !strings.Contains(err.Error(), "Access is denied") &&
			!strings.Contains(err.Error(), "permission denied") {
			t.Errorf("期望权限相关错误，实际错误: %v", err)
		}
	}
}
