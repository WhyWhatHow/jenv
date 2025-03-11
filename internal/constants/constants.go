package constants

const (
	// Windows 系统环境变量路径
	ENV_SYSTEM_PATH     = `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`
	ENV_USER_PATH       = `Environment`
	ENV_WIN_JAVA_HOME   = "%JAVA_HOME%\\bin;"
	ENV_LINUX_JAVA_HOME = "$JAVA_HOME/bin:"

	// 默认配置路径
	DEFAULT_CONFIG_FILE = "config.json"
	DEFAULT_FOLDER      = ".jdks"
	DEFAULT_BACKUP_FILE = "backup.json"

	// 默认符号链接路径
	DEFAULT_SYMLINK_PATH = "C:\\Java\\JAVA_HOME"
	DEFAULT_SYMLINK_NAME = "JAVA_HOME"
)
