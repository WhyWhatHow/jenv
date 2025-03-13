package java

import (
	"errors"
	"fmt"
	"github.com/whywhathow/jenv/internal/config"
	"github.com/whywhathow/jenv/internal/env"
	"github.com/whywhathow/jenv/internal/sys"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type JDK struct {
	Name string
	Path string
}

var ErrNoJDKConfigured = errors.New("no JDK configured")
var cfg *config.Config

/**
 *1.  init config.json, backup.json
 *2.  set env : JAVA_HOME and Path
 */
func init() {
	cfg, _ = config.GetInstance()
	//if err != nil {
	//	return fmt.Errorf("加载配置失败: %v", err)
	//}
	Init()
}
func Init() error {
	//cfg, err := config.GetInstance()
	//if err != nil {
	//	return nil, err
	//}

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
	//cfg, err := config.GetInstance()
	//if err != nil {
	//	return nil, err
	//}
	return cfg.Jdks, nil
}

// AddJDK 添加新的 JDK
func AddJDK(name, path string) error {
	// 获取配置实例
	//cfg, err := config.GetInstance()
	//if err != nil {
	//	return err
	//}

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
	//cfg, err := config.GetInstance()
	//if err != nil {
	//	return err
	//}

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
	//cfg, err := config.GetInstance()
	//if err != nil {
	//	return err
	//}

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
	//cfg, err := config.GetInstance()
	//if err != nil {
	//	return config.JDK{}, err
	//}

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
func ScanJDK(dir string) []JDK {
	// Validate directory
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("Directory does not exist: %s\n", dir)
		return nil
	}

	// Scan for JDKs
	jdks := scanForJDKs(dir)
	if len(jdks) == 0 {
		fmt.Println("No valid JDKs found in the specified directory")
		return nil
	}

	return jdks
}

/*
*
1 . 扫描每一个目录
2. 对于是 jdk 所在的目录， ,生成后退出当前目录，不再搜索当前目录以及其 child
3. 继续遍历。
*/
func scanForJDKs(dir string) []JDK {
	var jdks []JDK
	queue := []string{dir}

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		entries, _ := os.ReadDir(path)
		for _, entry := range entries {
			fullPath := filepath.Join(path, entry.Name())

			// 预检2：通过目录特征快速判断
			if info, _ := entry.Info(); !info.IsDir() {
				continue
			}

			// 跳过 系统目录  res :  7s-2s
			if (strings.EqualFold(entry.Name(), "Windows") && strings.HasPrefix(strings.ToLower(path), strings.ToLower(os.Getenv("SystemDrive")))) ||
				strings.EqualFold(entry.Name(), "$Recycle.Bin") ||
				strings.EqualFold(entry.Name(), "System Volume Information") {
				continue
			}

			// 目录深度检查
			if strings.Count(fullPath[len(dir):], string(os.PathSeparator)) > 3 {
				log.Print("Skipping directory:", fullPath, "\n")
				continue
			}

			// 完整校验
			if config.ValidateJavaPath(fullPath) {
				jdks = append(jdks, JDK{Path: fullPath, Name: entry.Name()})
				continue // 发现有效目录后跳过子目录
			}

			// 添加到队列继续扫描
			if entry.IsDir() {
				queue = append(queue, fullPath)
			}
		}
	}
	return jdks
}
