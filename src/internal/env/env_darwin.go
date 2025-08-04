//go:build darwin
// +build darwin

package env

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SetSystemPath 永久设置系统 PATH 环境变量 (macOS 平台实现)
func SetSystemPath(path string) error {
	// 获取当前的 PATH
	currentPath := os.Getenv("PATH")

	// 如果新路径不在当前 PATH 中，则添加它
	if !strings.Contains(currentPath, path) {
		// 将新路径添加到 PATH 的开头
		newPath := path
		if currentPath != "" {
			newPath = path + string(filepath.ListSeparator) + currentPath
		}

		// 更新系统环境变量
		return UpdateEnvironmentVariable("PATH", newPath)
	}

	return nil
}

// UpdateEnvironmentVariable 更新 macOS 系统环境变量
func UpdateEnvironmentVariable(key, value string) error {
	// 当前进程立即生效
	err := os.Setenv(key, value)
	if err != nil {
		return err
	}

	// 尝试更新用户级环境变量
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("无法获取用户主目录: %v", err)
	}

	// 尝试更新用户的 shell 配置文件 (优先.zshrc，其次.bash_profile)
	zshrcPath := filepath.Join(homeDir, ".zshrc")
	if err := updateShellProfile(zshrcPath, key, value); err != nil {
		bashProfilePath := filepath.Join(homeDir, ".bash_profile")
		if err := updateShellProfile(bashProfilePath, key, value); err != nil {
			return fmt.Errorf("无法更新环境变量: %v", err)
		}
	}

	return nil
}

// updateShellProfile 更新用户的 shell 配置文件
func updateShellProfile(filePath, key, value string) error {
	// 检查文件是否存在，如果不存在则创建
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		file.Close()
	}

	// 读取现有文件内容
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 读取所有行
	var lines []string
	scanner := bufio.NewScanner(file)
	keyFound := false
	exportStr := fmt.Sprintf("export %s=\"%s\"", key, value)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, fmt.Sprintf("export %s=", key)) {
			lines = append(lines, exportStr)
			keyFound = true
		} else {
			lines = append(lines, line)
		}
	}

	// 如果没有找到该环境变量，则添加它
	if !keyFound {
		lines = append(lines, exportStr)
	}

	// 写回文件
	return os.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644)
}

// SetPersistentEnvVar 在 macOS 上永久设置用户级环境变量
func SetPersistentEnvVar(name, value string) error {
	// 获取用户主目录
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// 目标文件为 ~/.zshrc (macOS Catalina及以后版本的默认shell)
	profile := filepath.Join(home, ".zshrc")

	// 以追加模式打开文件，若不存在则创建
	f, err := os.OpenFile(profile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// 写入 export 语句
	_, err = fmt.Fprintf(f, "\nexport %s=\"%s\"\n", name, value)
	return err
}
