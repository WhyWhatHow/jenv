package env

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// oldPath 保存原始的 PATH 环境变量值
var oldPath string

// SetEnv 临时设置环境变量
func SetEnv(key string, value string) error {
	if runtime.GOOS == "windows" {
		return SetEnvInWin(key, value)
	} else {
		return os.Setenv(key, value)
	}
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
