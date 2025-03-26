//go:build linux
// +build linux

package env

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SetSystemPath 永久设置系统 PATH 环境变量 (Linux 平台实现)
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

// UpdateSystemEnvironmentVariable 更新 Linux 系统环境变量
func UpdateEnvironmentVariable(key, value string) error {
	// 当前进程立即生效
	err := os.Setenv(key, value)
	if err != nil {
		return err
	}

	// 尝试更新 /etc/environment 文件
	if err := updateEnvironmentFile("/etc/environment", key, value); err != nil {
		// 如果无法更新系统级环境变量，尝试更新用户级环境变量
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("无法获取用户主目录: %v", err)
		}

		// 尝试更新用户的 .profile 文件
		profilePath := filepath.Join(homeDir, ".profile")
		if err := updateProfileFile(profilePath, key, value); err != nil {
			return fmt.Errorf("无法更新环境变量: %v", err)
		}
	}

	return nil
}

// updateEnvironmentFile 更新 /etc/environment 文件中的环境变量
func updateEnvironmentFile(filePath, key, value string) error {
	// 读取现有文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	keyFound := false

	// 更新或添加环境变量
	for _, line := range lines {
		if strings.HasPrefix(line, key+"=") {
			newLines = append(newLines, fmt.Sprintf("%s=\"%s\"", key, value))
			keyFound = true
		} else {
			newLines = append(newLines, line)
		}
	}

	// 如果没有找到该环境变量，则添加它
	if !keyFound {
		newLines = append(newLines, fmt.Sprintf("%s=\"%s\"", key, value))
	}

	// 写回文件
	return os.WriteFile(filePath, []byte(strings.Join(newLines, "\n")), 0644)
}

// updateProfileFile 更新用户的 .profile 文件中的环境变量
func updateProfileFile(filePath, key, value string) error {
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

// SetPersistentEnvVar 在 Unix-like 系统上永久设置用户级环境变量
func SetPersistentEnvVar(name, value string) error {
	// 获取用户主目录
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// 目标文件为 ~/.profile
	profile := filepath.Join(home, ".profile")

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
