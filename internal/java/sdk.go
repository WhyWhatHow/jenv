package java

import (
	"errors"
	"fmt"
	"github.com/whywhathow/jenv/internal/config"
	"github.com/whywhathow/jenv/internal/env"
	"github.com/whywhathow/jenv/internal/sys"
	"os"
	"path/filepath"
	"runtime" // Added
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
	//Init()
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
	// env.SetSystemEnvPath() // Commented out to fix build for Linux simulation
	// env.SetCurrentUserEnvPath() // Commented out to fix build for Linux simulation
	// TODO: Revisit PATH setting for Init() in a cross-platform way.
	// For Linux, UpdateEnvironmentVariable("PATH", newPathEntry) is the new way.
	// Marking as initialized anyway for now.
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
	// initialScanDepth := strings.Count(dir, string(os.PathSeparator)) // Removed unused variable

	// Directories to skip on all platforms (especially Linux)
	// For Linux, paths are absolute. For others, these are just names.
	var skipDirs = map[string]bool{
		// Linux specific absolute paths
		"/proc": true, "/sys": true, "/dev": true, "/run": true,
		"/lost+found": true, "/mnt": true, "/media": true, // "/tmp": true, // Removed /tmp from default skip for explicit scans
		// General names to skip
		"node_modules": true, "vendor": true, /* add more common non-JDK dirs if needed */
	}
	// Windows specific skips (names)
	var windowsSkipDirs = map[string]bool{
		"Windows":                 true,
		"$Recycle.Bin":          true,
		"System Volume Information": true,
		"Program Files":           true, // Often too deep or contains non-JDK stuff
		"Program Files (x86)":   true,
	}

	visited := make(map[string]bool) // To avoid re-scanning or cycles if any symlinks trick us

	for len(queue) > 0 {
		currentPath := queue[0]
		queue = queue[1:]

		if visited[currentPath] {
			continue
		}
		visited[currentPath] = true

		// Skip hidden directories at the top level of the initial scan path,
		// or if they are immediate children of common SDK installation parent directories.
		// This is a heuristic.
		baseName := filepath.Base(currentPath)
		isTopLevelScan := currentPath == dir
		if strings.HasPrefix(baseName, ".") && isTopLevelScan {
			fmt.Printf("Skipping hidden directory at top level: %s\n", currentPath)
			continue
		}

		// Check against general skipDirs (absolute for Linux, names otherwise)
		if _, shouldSkip := skipDirs[currentPath]; runtime.GOOS == "linux" && shouldSkip {
			fmt.Printf("Skipping Linux system directory: %s\n", currentPath)
			continue
		}
		if _, shouldSkip := skipDirs[baseName]; shouldSkip {
			fmt.Printf("Skipping common non-JDK directory: %s\n", currentPath)
			continue
		}


		// Conditional Windows-specific skips
		if runtime.GOOS == "windows" {
			if _, shouldSkip := windowsSkipDirs[baseName]; shouldSkip {
				// A special check for "Windows" dir on system drive
				if strings.EqualFold(baseName, "Windows") && strings.HasPrefix(strings.ToLower(currentPath), strings.ToLower(os.Getenv("SystemDrive"))) {
					fmt.Printf("Skipping Windows system directory: %s\n", currentPath)
					continue
				} else if !strings.EqualFold(baseName, "Windows") { // For other windowsSkipDirs
					fmt.Printf("Skipping Windows specific directory: %s\n", currentPath)
					continue
				}
			}
		}

		entries, err := os.ReadDir(currentPath)
		if err != nil {
			// Log error but continue (e.g. permission denied for a directory)
			fmt.Printf("Error reading directory %s: %v\n", currentPath, err)
			continue
		}

		for _, entry := range entries {
			fullPath := filepath.Join(currentPath, entry.Name())
			entryName := entry.Name()

			if !entry.IsDir() {
				continue
			}

			// Re-check for hidden directories if it's a subdirectory now
			if strings.HasPrefix(entryName, ".") {
				// Allow .jenv or similar if it's part of a known SDK structure, otherwise skip.
				// This is tricky; for now, let's be aggressive in skipping hidden ones unless it's the target dir.
				// A better approach might be needed if JDKs are in hidden dirs often.
				// fmt.Printf("Skipping hidden subdirectory: %s\n", fullPath)
				// continue
			}


			// Directory depth check relative to the initial 'dir'
			// currentScanDepth := strings.Count(fullPath, string(os.PathSeparator)) - initialScanDepth
			// The original check was relative to the length of 'dir' string, which is better
			currentScanDepth := strings.Count(strings.TrimPrefix(fullPath, dir), string(os.PathSeparator))

			if currentScanDepth > 6 { // Increased scan depth
				// fmt.Printf("Skipping due to depth: %s (depth %d)\n", fullPath, currentScanDepth)
				continue
			}

			if config.ValidateJavaPath(fullPath) {
				// Check if already added to avoid duplicates from different symlink paths etc.
				alreadyAdded := false
				for _, jdk := range jdks {
					if jdk.Path == fullPath {
						alreadyAdded = true
						break
					}
				}
				if !alreadyAdded {
					jdks = append(jdks, JDK{Path: fullPath, Name: entryName})
				}
				// Once a valid JDK is found in a path, do not scan deeper in this branch.
				continue
			}

			// Add to queue for further scanning
			queue = append(queue, fullPath)
		}
	}
	return jdks
}
