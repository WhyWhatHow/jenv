//go:build windows

package env

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"os"
	"path/filepath"
	"strings"
)

const SYSTEM_PATH = `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`

// SetSystemPath 永久设置系统 PATH 环境变量 (Windows 平台实现)
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
		return UpdateSystemEnvironmentVariable("PATH", newPath)
	}

	return nil
}

// UpdateSystemEnvironmentVariable 更新 Windows 系统环境变量

func UpdateSystemEnvironmentVariable(key, value string) error {

	if !IsAdmin() {
		return fmt.Errorf("please run this command with admin privileges")
	}

	// 需要管理员权限设置系统环境变量 system
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, SYSTEM_PATH, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("Setting current Registry key:  %v failed: %v.", key, err)
	}
	defer k.Close()

	return k.SetStringValue(key, value)
}

const USER_PATH = `Environment`

// UpdateUserEnvironmentVariable hint : se
/**
 * @Description: update user level environment variable
 * @param key
 * @param value
 */
func UpdateUserEnvironmentVariable(key, value string) error {
	if !IsAdmin() {
		return fmt.Errorf("please run this command with admin privileges")
	}
	// 需要管理员权限设置系统环境变量 system
	k, err := registry.OpenKey(registry.CURRENT_USER, USER_PATH, registry.SET_VALUE)
	//k, err := registry.OpenKey(registry.CURRENT_USER, USER_PATH, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("Setting current Registry key:  %v failed: %v.", key, err)
	}
	defer k.Close()

	return k.SetStringValue(key, value)
}
