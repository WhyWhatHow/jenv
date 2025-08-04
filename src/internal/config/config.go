package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/whywhathow/jenv/internal/constants"
)

var (
	ErrJDKExists            = errors.New("JDK already exists")
	ErrJDKNotFound          = errors.New("JDK not found")
	ErrInvalidPath    error = errors.New("Invalid Java installation path")
	ErrNotInitialized       = errors.New("jenvnot initialized")
)

var (
	instance     *Config
	instanceLock sync.RWMutex
	configPath   string
)

type Config struct {
	Current       string         `json:"current"`
	SymlinkPath   string         `json:"symlink_path"`
	Initialized   bool           `json:"initialized"`
	EnvBackUpPath string         `json:"env_backup_path"`
	Jdks          map[string]JDK `json:"jdks"`  // 将数组改为 map
	Theme         string         `json:"theme"` // Current theme name
	// 添加互斥锁保护并发访问
	lock sync.RWMutex
}

type JDK struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// GetInstance 返回配置的单例实例
func GetInstance() (*Config, error) {
	instanceLock.RLock()
	if instance != nil {
		defer instanceLock.RUnlock()
		return instance, nil
	}
	instanceLock.RUnlock()

	instanceLock.Lock()
	defer instanceLock.Unlock()

	// 双重检查
	if instance != nil {
		return instance, nil
	}

	// 加载配置
	cfg, err := loadConfigFromFile()
	if err != nil {
		return nil, err
	}

	instance = cfg
	return instance, nil
}

// 从文件加载配置
func loadConfigFromFile() (*Config, error) {
	// 1. 获取config.json path
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}
	//2. 判断config.json 存在与否
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 如果不存在，则创建默认配置
		// 默认backupFilepath 与configFilePath 在同一个folder
		backupFilePath := GetDefaultBackupFilePath()
		defaultSymlinkPath := GetDefaultSymlinkPath()
		cfg := &Config{
			Current:       "",
			SymlinkPath:   defaultSymlinkPath,
			Initialized:   false,
			EnvBackUpPath: backupFilePath,
			Jdks:          make(map[string]JDK), // 初始化 map
		}

		// 保存配置到文件
		if err := cfg.doSave(); err != nil {
			return nil, fmt.Errorf("保存配置文件失败: %w", err)
		}

		return cfg, nil
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	// 解析配置文件
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
func (c *Config) Save() error {
	// 获取锁，防止并发修改
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.doSave()
}

// doSave 保存配置到文件
func (c *Config) doSave() error {

	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	// 新增：确保配置目录存在
	configDir := filepath.Dir(path)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// SetSymlinkPath 设置符号链接路径
func (c *Config) SetSymlinkPath(symlinkPath string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.SymlinkPath = symlinkPath
	c.doSave()
}

// SetCurrentJDK 设置当前JDK
func (c *Config) SetCurrentJDK(jdkName string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	// 检查JDK是否存在
	if _, exists := c.Jdks[jdkName]; !exists {
		return ErrJDKNotFound
	}

	c.Current = jdkName
	return c.doSave()
}

// AddJDK 添加新的JDK
func (c *Config) AddJDK(name, path string) error {
	// 验证Java路径
	if !ValidateJavaPath(path) {
		return ErrInvalidPath
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	// 检查是否已存在同名JDK
	if _, exists := c.Jdks[name]; exists {
		return ErrJDKExists
	}

	// 添加新JDK
	c.Jdks[name] = JDK{Name: name, Path: path}

	// 保存更新后的配置到文件
	return c.doSave()
}

// RemoveJDK 移除JDK
func (c *Config) RemoveJDK(name string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	// 检查JDK是否存在
	if _, exists := c.Jdks[name]; !exists {
		return ErrJDKNotFound
	}

	// 删除JDK
	delete(c.Jdks, name)

	// 如果移除的是当前JDK，则取消设置
	if c.Current == name {
		c.Current = ""
	}

	// 保存更新后的配置到文件
	return c.doSave()
}

// 保留原有的工具函数
func ValidateJavaPath(path string) bool {
	// 是否需要根据runtime不同进行不同的判断呢
	if runtime.GOOS == "windows" {
		javaExe := filepath.Join(path, "bin", "java.exe")
		_, err := os.Stat(javaExe)
		return err == nil
	} else {
		// 其他os判断
		java := filepath.Join(path, "bin", "java")
		_, err := os.Stat(java)
		return err == nil
	}
}

func GetDefaultSymlinkPath() string {
	switch runtime.GOOS {
	case "windows":
		return constants.DEFAULT_SYMLINK_PATH_WINDOWS
	case "linux":
		// 尝试使用系统级路径，如果没有权限则使用用户级路径
		if os.Geteuid() == 0 {
			return constants.DEFAULT_SYMLINK_PATH_LINUX
		}
		// 非root用户使用用户目录
		if dir, err := os.UserHomeDir(); err == nil {
			return filepath.Join(dir, constants.USER_SYMLINK_PATH_LINUX)
		}
		return constants.DEFAULT_SYMLINK_PATH_LINUX
	case "darwin":
		// 尝试使用系统级路径，如果没有权限则使用用户级路径
		if os.Geteuid() == 0 {
			return constants.DEFAULT_SYMLINK_PATH_DARWIN
		}
		// 非root用户使用用户目录
		if dir, err := os.UserHomeDir(); err == nil {
			return filepath.Join(dir, constants.USER_SYMLINK_PATH_DARWIN)
		}
		return constants.DEFAULT_SYMLINK_PATH_DARWIN
	default:
		// 默认情况下使用用户目录
		if dir, err := os.UserHomeDir(); err == nil {
			return filepath.Join(dir, constants.DEFAULT_SYMLINK_NAME)
		}
		return "/tmp/jenv_java_home"
	}
}

// CreateSymlink creates a symbolic link from source to target
func CreateSymlink(source, target string) error {
	// 如果目标已存在，先删除
	_ = os.Remove(target)

	// 创建符号链接的父目录
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return err
	}

	// 创建符号链接
	return os.Symlink(source, target)
}

// SetConfigPath sets a custom path for the config file
func SetConfigPath(path string) {
	instanceLock.Lock()
	defer instanceLock.Unlock()

	configPath = path
	// 重置实例，强制下次重新加载
	instance = nil
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
	// 如果已设置自定义配置路径，直接返回
	if configPath != "" {
		return configPath, nil
	}

	// 否则使用默认路径
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(home, constants.DEFAULT_FOLDER)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}
	var configFile = filepath.Join(configDir, constants.DEFAULT_CONFIG_FILE)

	return configFile, nil
}

// InitializeConfig 初始化配置文件，如果配置文件不存在则创建
func InitializeConfig(configFilePath string) error {

	// 检查配置文件是否存在
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// 创建默认配置
		cfg := &Config{
			Jdks:          make(map[string]JDK),
			Current:       "",
			SymlinkPath:   GetDefaultSymlinkPath(), // 使用跨平台的默认路径
			Initialized:   false,
			EnvBackUpPath: "",
		}

		// 保存配置到文件
		if err := cfg.doSave(); err != nil {
			return fmt.Errorf("保存配置文件失败: %v", err)
		}
	}

	return nil
}
