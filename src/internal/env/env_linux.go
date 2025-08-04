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

func doSetEnv(key, value string) error {
	// 当前进程立即生效
	if err := os.Setenv(key, value); err != nil {
		return fmt.Errorf("failed to set environment variable in current process: %v", err)
	}

	// 尝试更新系统级环境变量（需要root权限）
	if os.Geteuid() == 0 {
		if err := updateSystemEnvironmentVariable(key, value); err != nil {
			// 系统级更新失败，继续尝试用户级更新
			fmt.Printf("Warning: Failed to update system environment variable: %v\n", err)
		}
	}

	// 更新用户级环境变量（支持多种shell）
	return updateUserEnvironmentVariables(key, value)
}

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

// updateSystemEnvironmentVariable 更新系统级环境变量（需要root权限）
func updateSystemEnvironmentVariable(key, value string) error {
	// 尝试更新 /etc/environment 文件
	if err := updateEnvironmentFile("/etc/environment", key, value); err != nil {
		return fmt.Errorf("failed to update /etc/environment: %v", err)
	}
	return nil
}

// updateUserEnvironmentVariables 更新用户级环境变量，支持多种shell
func updateUserEnvironmentVariables(key, value string) error {
	// 使用新的shell模块来处理多shell环境
	// 由于import问题，暂时使用原有的实现方式
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %v", err)
	}

	// 检测当前使用的shell
	shells := detectUserShells(homeDir)

	var errors []string
	for _, shell := range shells {
		if err := updateShellEnvironment(homeDir, shell, key, value); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", shell, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to update some shell environments: %s", strings.Join(errors, "; "))
	}

	return nil
}

// detectUserShells 检测用户使用的shell环境
func detectUserShells(homeDir string) []string {
	var shells []string

	// 检查常见的shell配置文件
	shellConfigs := map[string]string{
		"bash":    ".bashrc",
		"zsh":     ".zshrc",
		"fish":    ".config/fish/config.fish",
		"profile": ".profile", // 通用配置文件
	}

	for shell, configFile := range shellConfigs {
		configPath := filepath.Join(homeDir, configFile)
		if shell == "fish" {
			// fish的配置文件在子目录中
			if _, err := os.Stat(filepath.Dir(configPath)); err == nil {
				shells = append(shells, shell)
			}
		} else {
			if _, err := os.Stat(configPath); err == nil {
				shells = append(shells, shell)
			}
		}
	}

	// 如果没有找到任何shell配置文件，默认使用.profile
	if len(shells) == 0 {
		shells = append(shells, "profile")
	}

	return shells
}

// updateShellEnvironment 更新特定shell的环境变量
func updateShellEnvironment(homeDir, shell, key, value string) error {
	switch shell {
	case "bash":
		return updateBashEnvironment(homeDir, key, value)
	case "zsh":
		return updateZshEnvironment(homeDir, key, value)
	case "fish":
		return updateFishEnvironment(homeDir, key, value)
	case "profile":
		return updateProfileEnvironment(homeDir, key, value)
	default:
		return fmt.Errorf("unsupported shell: %s", shell)
	}
}

// updateBashEnvironment 更新bash环境变量
func updateBashEnvironment(homeDir, key, value string) error {
	bashrcPath := filepath.Join(homeDir, ".bashrc")
	return updateShellConfigFile(bashrcPath, key, value, "export")
}

// updateZshEnvironment 更新zsh环境变量
func updateZshEnvironment(homeDir, key, value string) error {
	zshrcPath := filepath.Join(homeDir, ".zshrc")
	return updateShellConfigFile(zshrcPath, key, value, "export")
}

// updateFishEnvironment 更新fish环境变量
func updateFishEnvironment(homeDir, key, value string) error {
	configDir := filepath.Join(homeDir, ".config", "fish")
	configPath := filepath.Join(configDir, "config.fish")

	// 确保配置目录存在
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create fish config directory: %v", err)
	}

	return updateFishConfigFile(configPath, key, value)
}

// updateProfileEnvironment 更新.profile环境变量
func updateProfileEnvironment(homeDir, key, value string) error {
	profilePath := filepath.Join(homeDir, ".profile")
	return updateShellConfigFile(profilePath, key, value, "export")
}

// updateShellConfigFile 更新shell配置文件（bash/zsh/.profile通用）
func updateShellConfigFile(filePath, key, value, exportCmd string) error {
	// 确保文件存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if file, err := os.Create(filePath); err != nil {
			return fmt.Errorf("failed to create config file %s: %v", filePath, err)
		} else {
			file.Close()
		}
	}

	// 读取现有内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %v", filePath, err)
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	keyFound := false
	exportLine := fmt.Sprintf("%s %s=\"%s\"", exportCmd, key, value)
	exportPrefix := fmt.Sprintf("%s %s=", exportCmd, key)

	// 查找并更新现有的环境变量设置
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, exportPrefix) {
			newLines = append(newLines, exportLine)
			keyFound = true
		} else {
			newLines = append(newLines, line)
		}
	}

	// 如果没有找到现有设置，添加新的
	if !keyFound {
		// 添加注释说明这是jenv设置的
		newLines = append(newLines, "")
		newLines = append(newLines, "# Added by jenv")
		newLines = append(newLines, exportLine)
	}

	// 写回文件
	return os.WriteFile(filePath, []byte(strings.Join(newLines, "\n")), 0644)
}

// updateFishConfigFile 更新fish配置文件
func updateFishConfigFile(filePath, key, value string) error {
	// 确保文件存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if file, err := os.Create(filePath); err != nil {
			return fmt.Errorf("failed to create fish config file: %v", err)
		} else {
			file.Close()
		}
	}

	// 读取现有内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read fish config file: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	keyFound := false
	setLine := fmt.Sprintf("set -gx %s \"%s\"", key, value)
	setPrefix := fmt.Sprintf("set -gx %s ", key)

	// 查找并更新现有的环境变量设置
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, setPrefix) {
			newLines = append(newLines, setLine)
			keyFound = true
		} else {
			newLines = append(newLines, line)
		}
	}

	// 如果没有找到现有设置，添加新的
	if !keyFound {
		// 添加注释说明这是jenv设置的
		newLines = append(newLines, "")
		newLines = append(newLines, "# Added by jenv")
		newLines = append(newLines, setLine)
	}

	// 写回文件
	return os.WriteFile(filePath, []byte(strings.Join(newLines, "\n")), 0644)
}

// QuerySystemEnvironmentVariable 查询系统环境变量（Linux实现）
func QuerySystemEnvironmentVariable(key string) (string, error) {
	// 在Linux中，首先尝试从当前环境获取
	if value := os.Getenv(key); value != "" {
		return value, nil
	}

	// 尝试从/etc/environment读取
	if value, err := readFromEnvironmentFile("/etc/environment", key); err == nil && value != "" {
		return value, nil
	}

	// 如果都没有找到，返回空字符串
	return "", nil
}

// UpdateSystemEnvironmentVariable 更新系统环境变量（Linux实现）
func UpdateSystemEnvironmentVariable(key, value string) error {
	return updateSystemEnvironmentVariable(key, value)
}

// readFromEnvironmentFile 从环境文件中读取特定的环境变量
func readFromEnvironmentFile(filePath, key string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, key+"=") {
			// 提取值，去掉引号
			value := strings.TrimPrefix(line, key+"=")
			value = strings.Trim(value, "\"'")
			return value, nil
		}
	}

	return "", fmt.Errorf("environment variable %s not found", key)
}

// SetSystemEnvPath 永久设置系统 PATH 环境变量 (Linux 平台实现)
func SetSystemEnvPath() error {
	// 获取当前的 PATH
	currentPath := os.Getenv("PATH")

	// 获取默认的JAVA_HOME路径
	javaHomeBin := getDefaultJavaHome()

	// 清理PATH中已存在的JAVA_HOME路径
	currentPath = cleanPathLinux(currentPath, javaHomeBin)

	// 将JAVA_HOME/bin添加到PATH开头
	newPath := javaHomeBin + ":" + currentPath

	// 更新系统环境变量
	return updateSystemEnvironmentVariable("PATH", newPath)
}

// SetCurrentUserEnvPath 设置当前用户的 PATH 环境变量 (Linux 平台实现)
func SetCurrentUserEnvPath() error {
	// 获取当前的 PATH
	currentPath := os.Getenv("PATH")

	// 获取默认的JAVA_HOME路径
	javaHomeBin := getDefaultJavaHome()

	// 清理PATH中已存在的JAVA_HOME路径
	currentPath = cleanPathLinux(currentPath, javaHomeBin)

	// 将JAVA_HOME/bin添加到PATH开头
	newPath := javaHomeBin + ":" + currentPath

	// 更新用户环境变量
	return updateUserEnvironmentVariables("PATH", newPath)
}

// cleanPathLinux 清理Linux PATH中的指定路径
func cleanPathLinux(path string, binPath string) string {
	if path == "" {
		return ""
	}

	paths := strings.Split(path, ":")
	var cleanedPaths []string

	for _, p := range paths {
		if p != binPath && strings.TrimSpace(p) != "" {
			cleanedPaths = append(cleanedPaths, p)
		}
	}

	return strings.Join(cleanedPaths, ":")
}

// QueryUserEnvironmentVariable 查询用户环境变量（Linux实现）
func QueryUserEnvironmentVariable(key string) (string, error) {
	// 在Linux中，用户环境变量通常存储在shell配置文件中
	// 首先尝试从当前环境获取
	if value := os.Getenv(key); value != "" {
		return value, nil
	}

	// 尝试从用户的shell配置文件中读取
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %v", err)
	}

	// 检查常见的配置文件
	configFiles := []string{".bashrc", ".zshrc", ".profile"}
	for _, configFile := range configFiles {
		configPath := filepath.Join(homeDir, configFile)
		if value, err := readExportFromFile(configPath, key); err == nil && value != "" {
			return value, nil
		}
	}

	return "", fmt.Errorf("user environment variable %s not found", key)
}

// readExportFromFile 从shell配置文件中读取export语句的值
func readExportFromFile(filePath, key string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(content), "\n")
	exportPrefix := fmt.Sprintf("export %s=", key)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, exportPrefix) {
			// 提取值，去掉引号
			value := strings.TrimPrefix(line, exportPrefix)
			value = strings.Trim(value, "\"'")
			return value, nil
		}
	}

	return "", fmt.Errorf("export %s not found in %s", key, filePath)
}

// UpdateUserEnvironmentVariable 更新用户环境变量（Linux实现）
func UpdateUserEnvironmentVariable(key, value string) error {
	return updateUserEnvironmentVariables(key, value)
}
