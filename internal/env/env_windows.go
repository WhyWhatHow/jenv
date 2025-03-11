//go:build windows

package env

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"jenv-go/internal/constants"
	"jenv-go/internal/sys"
	"runtime"
)

//const ENV_SYSTEM_PATH = `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`

func getDefaultJAVAHOME() string {
	if runtime.GOOS == "windows" {
		return constants.ENV_WIN_JAVA_HOME
	} else {
		return constants.ENV_LINUX_JAVA_HOME

	}
}
func SetEnvInWin(key, value string) error {
	err := UpdateSystemEnvironmentVariable(key, value)
	if err != nil {
		return err
	}
	err1 := UpdateUserEnvironmentVariable(key, value)
	if err1 != nil {
		return err1
	}
	return nil
}

// SetSystemPath 永久设置系统 PATH 环境变量 (Windows 平台实现)
func SetSystemEnvPath() error {
	// 获取当前的 PATH
	currentPath, err := QuerySystemEnvironmentVariable("PATH")
	if err != nil {
		return err
	}
	defaultJavaHome := getDefaultJAVAHOME()
	var newPath = defaultJavaHome + currentPath
	return UpdateSystemEnvironmentVariable("PATH", newPath)
}

func SetCurrentUserEnvPath() error {
	// 获取当前的 PATH
	currentPath, err := QueryUserEnvironmentVariable("PATH")
	if err != nil {
		return err
	}
	defaultJavaHome := getDefaultJAVAHOME()
	var newPath = defaultJavaHome + currentPath
	// 更新系统环境变量
	return UpdateUserEnvironmentVariable("PATH", newPath)

}

// UpdateSystemEnvironmentVariable 更新 Windows 系统环境变量

func UpdateSystemEnvironmentVariable(key, value string) error {

	if !sys.IsAdmin() {
		return fmt.Errorf("please run this command with admin privileges")
	}

	// 需要管理员权限设置系统环境变量 system
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, constants.ENV_SYSTEM_PATH, registry.SET_VALUE)
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
	if !sys.IsAdmin() {
		return fmt.Errorf("please run this command with admin privileges")
	}
	// 需要管理员权限设置系统环境变量 system
	k, err := registry.OpenKey(registry.CURRENT_USER, constants.ENV_USER_PATH, registry.SET_VALUE)
	//k, err := registry.OpenKey(registry.CURRENT_USER, constants.ENV_USER_PATH, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("Setting current Registry key:  %v failed: %v.", key, err)
	}
	defer k.Close()

	return k.SetStringValue(key, value)
}

// query user level environment variable
func QueryUserEnvironmentVariable(key string) (string, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, USER_PATH, registry.QUERY_VALUE)
	if err != nil {
		return "", fmt.Errorf("Setting current Registry key:  %v failed: %v.", key, err)
	}
	defer k.Close()

	value, _, err := k.GetStringValue(key)
	if err != nil {
		return "", fmt.Errorf("Setting current Registry key:  %v failed: %v.", key, err)
	}

	return value, nil
}

// query system level environment variable
func QuerySystemEnvironmentVariable(key string) (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, constants.ENV_SYSTEM_PATH, registry.QUERY_VALUE)
	if err != nil {
		return "", fmt.Errorf("Setting current Registry key:  %v failed: %v.", key, err)
	}
	defer k.Close()

	value, _, err := k.GetStringValue(key)
	if err != nil {
		return "", fmt.Errorf("Setting current Registry key:  %v failed: %v.", key, err)
	}
	return value, nil
}
