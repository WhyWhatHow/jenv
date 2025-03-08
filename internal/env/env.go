package env

import (
	"os"
	"path/filepath"
	"strings"
)

// oldPath 保存原始的 PATH 环境变量值
var oldPath string

// SetEnv 临时设置环境变量
func SetEnv(key string, value string) error {
	return os.Setenv(key, value)
}

// GetEnv 获取环境变量
func GetEnv(key string) string {
	return os.Getenv(key)
}

// CleanJDKPath 清理 PATH 中的 JDK 相关设置
func CleanJDKPath() error {
	// 获取当前的 PATH
	path := os.Getenv("PATH")

	// 将 PATH 分割成数组
	paths := strings.Split(path, string(filepath.ListSeparator))

	// 过滤掉包含 jdk、jre 的路径
	var cleanedPaths []string
	for _, p := range paths {
		if !strings.Contains(strings.ToLower(p), "jdk") && !strings.Contains(strings.ToLower(p), "jre") {
			cleanedPaths = append(cleanedPaths, p)
		}
	}

	// 重新组合 PATH
	newPath := strings.Join(cleanedPaths, string(filepath.ListSeparator))

	// 设置新的 PATH
	return os.Setenv("PATH", newPath)
}

// RestoreOldPath 恢复原始的 PATH 环境变量
func RestoreOldPath() error {
	if oldPath != "" {
		return os.Setenv("PATH", oldPath)
	}
	return nil
}

// SetEnvWithPath 设置环境变量，并将其 bin 目录添加到 PATH 中
func SetEnvWithPath(key string, value string) error {
	// 设置环境变量
	if err := os.Setenv(key, value); err != nil {
		return err
	}

	// 获取当前的 PATH
	path := os.Getenv("PATH")

	// 添加环境变量的 bin 目录到 PATH，避免重复添加
	binPath := filepath.Join(value, "bin")
	if path == "" {
		path = binPath
	} else if !strings.Contains(path, binPath) {
		newPath := binPath + string(filepath.ListSeparator) + path
		path = newPath
	}

	// 设置新的 PATH
	if err := os.Setenv("PATH", path); err != nil {
		return err
	}

	return nil
}
