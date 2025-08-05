package java

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/whywhathow/jenv/internal/config"
	"github.com/whywhathow/jenv/internal/env"
	"github.com/whywhathow/jenv/internal/sys"
)

type JDK struct {
	Name string
	Path string
}

var ErrNoJDKConfigured = errors.New("no JDK configured")
var cfg *config.Config

// Simple cache for directory scan results to improve performance and avoid redundant scans
var (
	scanCache      = make(map[string]time.Time)
	scanCacheMutex sync.RWMutex
	cacheTimeout   = 5 * time.Minute // Cache results for 5 minutes
)

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

// ScanResult 包含扫描操作的详细结果
type ScanResult struct {
	JDKs     []JDK
	Duration time.Duration
	Scanned  int // 实际扫描的目录数
	Skipped  int // 因权限等问题跳过的目录数
	Excluded int // 因已存在或规则被排除的目录数
}

// Task 定义了需要扫描的目录任务
type WorkerTask struct {
	Path  string
	Depth int
}

// WorkerResult 是工人完成任务后返回的结果
type WorkerResult struct {
	FoundJDK    *JDK
	SubDirTasks []WorkerTask
	IsSkipped   bool
	IsExcluded  bool
}

// maxDepth 定义了最大扫描深度
const maxDepth = 5

// numWorkers 定义了并发的工人数量
var numWorkers = runtime.NumCPU() * 2 // 保持动态，但可以根据需要调整

// ScanJDK 是一个简单的包装器，只返回找到的JDK列表
func ScanJDK(dir string) []JDK {
	result := ScanJDKWithStats(dir)
	return result.JDKs
}

// ScanJDKWithStats 使用健壮的并发模型执行JDK扫描并返回详细统计信息
func ScanJDKWithStats(dir string) ScanResult {
	start := time.Now()
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("Directory does not exist: %s\n", dir)
		return ScanResult{Duration: time.Since(start)}
	}

	existingPaths := getExistingJDKPaths()

	// --- 调度中心-工人 并发模型 ---
	tasksChan := make(chan WorkerTask, numWorkers*2)
	resultsChan := make(chan WorkerResult, numWorkers*2)
	var workerWg sync.WaitGroup

	// 1. 启动固定数量的工人
	for i := 0; i < numWorkers; i++ {
		workerWg.Add(1)
		go jdkScannerWorker(tasksChan, resultsChan, &workerWg, existingPaths)
	}

	// 2. 启动调度中心 goroutine
	var finalJDKs []JDK
	var stats ScanResult
	var dispatcherWg sync.WaitGroup
	dispatcherWg.Add(1)

	go func() {
		defer dispatcherWg.Done()
		taskQueue := []WorkerTask{{Path: dir, Depth: 1}}
		pendingTasks := 1

		for pendingTasks > 0 {
			var currentTask WorkerTask
			var sendChan chan WorkerTask

			if len(taskQueue) > 0 {
				currentTask = taskQueue[0]
				sendChan = tasksChan // 只有队列中有任务时，才准备发送
			}

			select {
			case sendChan <- currentTask:
				taskQueue = taskQueue[1:] // 任务已发送，从队列移除

			case result := <-resultsChan:
				pendingTasks-- // 一个任务完成了

				// 更新统计数据
				stats.Scanned++
				if result.IsSkipped {
					stats.Skipped++
				}
				if result.IsExcluded {
					stats.Excluded++
				}

				// 如果找到了 JDK，收集它
				if result.FoundJDK != nil {
					finalJDKs = append(finalJDKs, *result.FoundJDK)
				}

				// 将新的子任务加入队列
				if len(result.SubDirTasks) > 0 {
					taskQueue = append(taskQueue, result.SubDirTasks...)
					pendingTasks += len(result.SubDirTasks)
				}
			}
		}

		// 所有任务处理完毕，关闭任务通道，让工人们退出
		close(tasksChan)
	}()

	// 3. 等待调度中心和所有工人都完成
	dispatcherWg.Wait()
	workerWg.Wait()

	stats.JDKs = finalJDKs
	stats.Duration = time.Since(start)

	if len(stats.JDKs) == 0 {
		fmt.Println("No valid JDKs found in the specified directory")
	}

	return stats
}

// jdkScannerWorker 是并发模型中的“工人”，负责处理单个目录的扫描
func jdkScannerWorker(tasks <-chan WorkerTask, results chan<- WorkerResult, wg *sync.WaitGroup, existingPaths map[string]bool) {
	defer wg.Done()
	for task := range tasks {
		res := WorkerResult{}

		// 预过滤：在读取目录之前进行检查
		if !shouldScanDirectory(task.Path) || task.Depth > maxDepth {
			// if task.Depth > maxDepth {
			res.IsSkipped = true
			results <- res
			continue
		}

		if isPathExcluded(task.Path, existingPaths) {
			res.IsExcluded = true
			results <- res
			continue
		}

		// 核心逻辑：检查当前目录是否为JDK
		if config.ValidateJavaPath(task.Path) {
			res.FoundJDK = &JDK{Path: task.Path, Name: filepath.Base(task.Path)}
			results <- res
			continue // 找到JDK后，不再扫描其子目录
		}

		// 读取子目录
		entries, err := os.ReadDir(task.Path)
		if err != nil {
			res.IsSkipped = true
			results <- res
			continue
		}

		// 准备子任务
		for _, entry := range entries {
			if entry.IsDir() {
				fullPath := filepath.Join(task.Path, entry.Name())
				res.SubDirTasks = append(res.SubDirTasks, WorkerTask{Path: fullPath, Depth: task.Depth + 1})
			}
		}
		results <- res
	}
}

// --- 以下是您原有的辅助函数，我进行了一些微调和简化，以更好地服务于新的并发模型 ---

// getExistingJDKPaths 保持不变，功能明确且高效
func getExistingJDKPaths() map[string]bool {
	existingPaths := make(map[string]bool)
	jdks, err := ListJdks()
	if err != nil {
		return existingPaths
	}
	for _, jdk := range jdks {
		normalizedPath, _ := filepath.Abs(jdk.Path)
		existingPaths[strings.ToLower(normalizedPath)] = true
	}
	return existingPaths
}

// isPathExcluded 简化和优化
func isPathExcluded(path string, existingPaths map[string]bool) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return true // 无法解析路径，跳过
	}
	normalizedPath := strings.ToLower(absPath)

	// 精确匹配或子目录匹配
	for existing := range existingPaths {
		if strings.HasPrefix(normalizedPath, existing) {
			// 确保是完全匹配或子目录，而不是 "c:\java\jdk1" 匹配 "c:\java\jdk11"
			if len(normalizedPath) == len(existing) || normalizedPath[len(existing)] == os.PathSeparator {
				return true
			}
		}
	}
	return false
}

// shouldScanDirectory performs aggressive pre-filtering to skip directories that are unlikely to contain JDKs
func shouldScanDirectory(path string) bool {

	// Get directory name for filtering
	dirName := strings.ToLower(filepath.Base(path))

	// Skip common non-JDK directories aggressively
	skipDirs := []string{
		// System directories
		"windows", "system32", "syswow64", "drivers", "winsxs",
		"$recycle.bin", "system volume information", "recovery",

		// Common application directories that won't have JDKs
		"node_modules", ".git", ".svn", ".hg", "bin", "obj", "debug", "release",
		"temp", "tmp", "cache", "logs", "log", "backup", "backups",
		"downloads", "documents", "pictures", "music", "videos", "desktop",

		// Development tools (but not JDK locations)
		"visual studio", "microsoft visual studio", "jetbrains", "intellij",
		"eclipse", "netbeans", "android studio", "xamarin",

		// Package managers and build tools
		"npm", "yarn", "gradle", "maven", ".m2", "nuget", "pip", "conda",

		// Version control and IDE files
		".vscode", ".idea", ".vs", "target", "build", "dist", "out",

		// Common non-JDK subdirectories
		"src", "source", "sources", "test", "tests", "doc", "docs", "documentation",
		"examples", "samples", "demo", "demos", "tutorial", "tutorials",
	}

	for _, skip := range skipDirs {
		if dirName == skip {
			return false
		}
	}

	// Skip directories with certain patterns
	if strings.Contains(dirName, "temp") ||
		strings.Contains(dirName, "cache") ||
		strings.Contains(dirName, "backup") ||
		strings.HasPrefix(dirName, "~") {
		return false
	}
	return true
}
