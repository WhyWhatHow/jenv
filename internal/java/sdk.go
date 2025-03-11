package java

import (
	"fmt"
	"jenv-go/internal/config"
	"jenv-go/internal/env"
	"jenv-go/internal/sys"
)

/**
 *1.  init config.json, backup.json
 *2.  set env : JAVA_HOME and Path
 */
func Init() error {
	//1. generate config.json
	// 获取配置实例
	cfg, err := config.GetInstance()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	// 如果已经初始化过，则返回
	if cfg.Initialized {
		return fmt.Errorf("Already init jenv")
	}
	//2. generate backup.json
	config.BackupEnvPath() // 备份环境变量

	//3. set JAVA_HOME
	defaultSymlinkPath := config.GetDefaultSymlinkPath()
	cfg.SetSymlinkPath(defaultSymlinkPath)

	if err := env.SetEnv("JAVA_HOME", defaultSymlinkPath); err != nil {
		return fmt.Errorf("设置环境变量失败: %v", err)
	}
	// 4. set PATH
	env.SetSystemEnvPath()
	env.SetCurrentUserEnvPath()
	// 标记为已初始化
	cfg.Initialized = true

	// 保存配置
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("保存配置失败: %v", err)
	}

	return nil
}

// ListJdks 返回所有已注册的 Java JDK
func ListJdks() (map[string]config.JDK, error) {
	// 获取配置实例
	cfg, err := config.GetInstance()
	if err != nil {
		return nil, err
	}
	return cfg.Jdks, nil
}

// AddJDK 添加新的 JDK
func AddJDK(name, path string) error {
	// 获取配置实例
	cfg, err := config.GetInstance()
	if err != nil {
		return err
	}

	// 添加 JDK
	if err := cfg.AddJDK(name, path); err != nil {
		return err
	}

	// 保存配置
	return cfg.Save()
}

// RemoveJDK 移除 JDK
func RemoveJDK(name string) error {
	// 获取配置实例
	cfg, err := config.GetInstance()
	if err != nil {
		return err
	}

	// 移除 JDK
	if err := cfg.RemoveJDK(name); err != nil {
		return err
	}

	// 保存配置
	return cfg.Save()
}

// UseJDK 设置当前使用的 JDK
func UseJDK(name string) error {
	// 获取配置实例
	cfg, err := config.GetInstance()
	if err != nil {
		return err
	}

	// 检查 JDK 是否存在
	jdk, exists := cfg.Jdks[name]

	if !exists {
		return config.ErrJDKNotFound
	}

	// 创建符号链接
	if err := sys.CreateSymlink(jdk.Path, cfg.SymlinkPath); err != nil {
		return fmt.Errorf("创建符号链接失败: %v", err)
	}

	// 设置当前 JDK
	if err := cfg.SetCurrentJDK(name); err != nil {
		return err
	}

	// 保存配置
	return cfg.Save()
}

// GetCurrentJDK 获取当前使用的 JDK
func GetCurrentJDK() (config.JDK, error) {
	// 获取配置实例
	cfg, err := config.GetInstance()
	if err != nil {
		return config.JDK{}, err
	}

	// 如果没有设置当前 JDK，返回错误
	if cfg.Current == "" {
		return config.JDK{}, fmt.Errorf("当前未设置 JDK")
	}

	// 查找当前 JDK
	for _, jdk := range cfg.Jdks {
		if jdk.Name == cfg.Current {
			return jdk, nil
		}
	}

	return config.JDK{}, config.ErrJDKNotFound
}
