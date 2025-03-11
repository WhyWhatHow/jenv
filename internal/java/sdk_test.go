package java

import (
	"jenv-go/internal/config"

	"os"
	"path/filepath"
	"testing"
)

func TestInit(t *testing.T) {
	// 设置临时测试目录

	_, _ = os.UserHomeDir()

	testDir := filepath.Join(os.TempDir(), "jenv-test")
	defer os.RemoveAll(testDir)

	// 设置配置文件路径
	config.SetConfigPath(filepath.Join(testDir, "config.json"))

	t.Run("成功初始化", func(t *testing.T) {
		// 清理环境
		os.RemoveAll(testDir)
		os.MkdirAll(testDir, 0755)

		// 执行初始化
		err := Init()
		if err != nil {
			t.Fatalf("初始化失败: %v", err)
		}

		// 验证配置
		cfg, err := config.GetInstance()
		if err != nil {
			t.Fatalf("加载配置失败: %v", err)
		}

		if !cfg.Initialized {
			t.Error("配置未标记为已初始化")
		}

		if cfg.SymlinkPath != "C:\\Java\\JAVA_HOME" {
			t.Errorf("符号链接路径不正确，期望 C:\\Java\\JAVA_HOME，实际 %s", cfg.SymlinkPath)
		}

		// 验证环境变量
		if os.Getenv("JAVA_HOME") != "C:\\Java\\JAVA_HOME" {
			t.Error("JAVA_HOME 环境变量设置不正确")
		}
	})

	t.Run("重复初始化", func(t *testing.T) {
		// 清理环境
		os.RemoveAll(testDir)
		os.MkdirAll(testDir, 0755)

		// 先执行一次初始化
		if err := Init(); err != nil {
			t.Fatalf("第一次初始化失败: %v", err)
		}

		// 再次执行初始化
		err := Init()
		if err == nil {
			t.Error("期望重复初始化返回错误，但没有")
		}
	})

	t.Run("创建目录失败", func(t *testing.T) {
		// 清理环境
		os.RemoveAll(testDir)
		os.MkdirAll(testDir, 0755)

		// 创建一个只读文件，阻止创建目录
		if err := os.WriteFile(filepath.Join(testDir, "Java"), []byte{}, 0444); err != nil {
			t.Fatalf("创建测试文件失败: %v", err)
		}

		// 执行初始化
		err := Init()
		if err == nil {
			t.Error("期望创建目录失败返回错误，但没有")
		}
	})
}
