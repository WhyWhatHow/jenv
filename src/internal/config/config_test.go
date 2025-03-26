package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInstance(t *testing.T) {
	// 清理单例实例
	instance = nil
	configPath = ""

	// 获取实例
	cfg, err := GetInstance()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// 再次获取实例，确保单例
	cfg2, err := GetInstance()
	assert.NoError(t, err)
	assert.Equal(t, cfg, cfg2)
}

func TestAddJDK(t *testing.T) {
	// 清理单例实例
	instance = nil
	configPath = ""

	// 获取实例
	cfg, err := GetInstance()
	assert.NoError(t, err)

	jdkName := "jdk21"
	jdkPath := "C:\\Users\\JhonNash\\.jdks\\azul-21.0.5"
	// 添加 JDK
	err = cfg.AddJDK(jdkName, jdkPath)
	assert.NoError(t, err)

	// 验证 JDK 是否添加成功
	assert.Len(t, cfg.Jdks, 1)
	assert.Equal(t, jdkName, cfg.Jdks["jdk21"].Name)
	assert.Equal(t, jdkPath, cfg.Jdks["jdk21"].Path)
}

func TestRemoveJDK(t *testing.T) {
	// 清理单例实例
	instance = nil
	configPath = ""

	// 获取实例
	cfg, err := GetInstance()
	assert.NoError(t, err)
	jdkName := "jdk21"
	//jdkPath := "C:\\Users\\JhonNash\\.jdks\\azul-21.0.5"
	// 添加 JDK
	//err = cfg.AddJDK(jdkName, jdkPath)
	//assert.NoError(t, err)

	// 移除 JDK
	err = cfg.RemoveJDK(jdkName)
	assert.NoError(t, err)

	// 验证 JDK 是否移除成功
	assert.Len(t, cfg.Jdks, 0)
}

func TestSaveAndLoadConfig(t *testing.T) {
	// 清理单例实例
	instance = nil
	configPath = ""

	// 获取实例
	cfg, err := GetInstance()
	assert.NoError(t, err)
	jdkName := "jdk21"
	jdkPath := "C:\\Users\\JhonNash\\.jdks\\azul-21.0.5"
	// 添加 JDK
	err = cfg.AddJDK(jdkName, jdkPath)
	assert.NoError(t, err)

	// 保存配置
	err = cfg.Save()
	assert.NoError(t, err)

	// 重新加载配置
	cfg2, err := GetInstance()
	assert.NoError(t, err)

	// 验证配置是否正确加载
	assert.Len(t, cfg2.Jdks, 1)
	assert.Equal(t, jdkName, cfg2.Jdks["jdk21"].Name)
	assert.Equal(t, jdkPath, cfg2.Jdks["jdk21"].Path)
}

func TestSetCurrentJDK(t *testing.T) {
	// 清理单例实例
	instance = nil
	configPath = ""

	// 获取实例
	cfg, err := GetInstance()
	assert.NoError(t, err)
	jdkName := "jdk21"
	jdkPath := "C:\\Users\\JhonNash\\.jdks\\azul-21.0.5"
	// 添加 JDK
	err = cfg.AddJDK(jdkName, jdkPath)
	assert.NoError(t, err)

	// 设置当前 JDK
	err = cfg.SetCurrentJDK(jdkName)
	assert.NoError(t, err)

	// 验证当前 JDK 是否设置成功
	assert.Equal(t, jdkName, cfg.Current)
}
