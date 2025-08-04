//go:build windows

package env

import (
	"fmt"
	"github.com/whywhathow/jenv/internal/constants"
	"github.com/whywhathow/jenv/internal/sys"
	"golang.org/x/sys/windows/registry"
	"strings"
	"syscall"
	"unsafe"
)

const (
	HWND_BROADCAST   = 0xFFFF
	WM_SETTINGCHANGE = 0x001A
)

//const ENV_SYSTEM_PATH = `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`

// doSetEnv sets both system and user environment variables in Windows
func doSetEnv(key, value string) error {
	if err := UpdateSystemEnvironmentVariable(key, value); err != nil {
		return err
	}
	return UpdateUserEnvironmentVariable(key, value)
}

// SetSystemPath 永久设置系统 PATH 环境变量 (Windows 平台实现)
func SetSystemEnvPath() error {
	// 获取当前的 PATH
	currentPath, err := QuerySystemEnvironmentVariable("PATH")
	if err != nil {
		return err
	}

	// 1. checking JAVA_HOME exist in Path or not
	currentPath = cleanPath(currentPath, getDefaultJavaHome())
	newPath := getDefaultJavaHome() + ";" + currentPath
	return UpdateSystemEnvironmentVariable("PATH", newPath)
}

// 删掉 path 中 binPath 路径。
func cleanPath(path string, binPath string) string {
	newPath := ``
	paths := strings.Split(path, ";")
	for _, p := range paths {
		if p == binPath {
			continue
		}
		if len(newPath) == 0 {
			newPath = p
		} else {
			newPath = newPath + ";" + p
		}
	}
	return newPath

}

// SetCurrentUserEnvPath sets the PATH environment variable for the current user
func SetCurrentUserEnvPath() error {
	// 获取当前的 PATH
	currentPath, err := QueryUserEnvironmentVariable("PATH")
	if err != nil {
		return err
	}
	path := cleanPath(currentPath, constants.ENV_WIN_JAVA_HOME)
	var newPath = constants.ENV_WIN_JAVA_HOME + ";" + path
	// 更新系统环境变量
	return UpdateUserEnvironmentVariable("PATH", newPath)

}

// broadcastEnvironmentChange notifies the system about environment variable changes
func broadcastEnvironmentChange() error {
	user32 := syscall.NewLazyDLL("user32.dll")
	sendMessageTimeout := user32.NewProc("SendMessageTimeoutW")
	envStr, err := syscall.UTF16PtrFromString("Environment")
	if err != nil {
		return fmt.Errorf("failed to convert environment string: %v", err)
	}

	ret, _, err := sendMessageTimeout.Call(
		HWND_BROADCAST,
		WM_SETTINGCHANGE,
		0,
		uintptr(unsafe.Pointer(envStr)),
		0,
		200,
		0,
	)

	if ret == 0 {
		return fmt.Errorf("failed to broadcast environment change: %v", err)
	}
	return nil
}

// UpdateSystemEnvironmentVariable updates a Windows system environment variable
func UpdateSystemEnvironmentVariable(key, value string) error {
	if !sys.IsAdmin() {
		return fmt.Errorf("please run this command with admin privileges")
	}

	k, err := registry.OpenKey(registry.LOCAL_MACHINE, constants.ENV_SYSTEM_PATH, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %v", err)
	}
	defer k.Close()

	if err := k.SetStringValue(key, value); err != nil {
		return fmt.Errorf("failed to set registry value for %v: %v", key, err)
	}

	return broadcastEnvironmentChange()
}

// UpdateUserEnvironmentVariable updates a Windows user environment variable
func UpdateUserEnvironmentVariable(key, value string) error {
	if !sys.IsAdmin() {
		return fmt.Errorf("please run this command with admin privileges")
	}

	k, err := registry.OpenKey(registry.CURRENT_USER, constants.ENV_USER_PATH, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("Setting current Registry key:  %v failed: %v.", key, err)
	}
	defer k.Close()

	err = k.SetStringValue(key, value)
	if err != nil {
		return fmt.Errorf("Setting current Registry key:  %v failed: %v.", key, err)
	}

	return broadcastEnvironmentChange()
}

// QueryUserEnvironmentVariable retrieves a Windows user environment variable
func QueryUserEnvironmentVariable(key string) (string, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, constants.ENV_USER_PATH, registry.QUERY_VALUE)
	if err != nil {
		return "", fmt.Errorf("failed to open registry key: %v", err)
	}
	defer k.Close()
	value, _, err := k.GetStringValue(key)
	if err != nil {
		return "", fmt.Errorf("failed to get registry value for %v: %v", key, err)
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
