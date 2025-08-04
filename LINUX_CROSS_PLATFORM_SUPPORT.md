# JEnv Linux跨平台支持功能总结

## 概述

本次更新为jenv项目添加了完整的Linux平台支持，实现了真正的跨平台Java版本管理工具。在保持Windows功能不变的基础上，新增了Linux特有的功能和优化。

## 新增和修改的文件

### 新增文件
- `src/internal/shell/shell.go` - Unix-like系统的Shell环境管理
- `src/internal/shell/shell_test.go` - Shell模块单元测试
- `src/internal/env/env_unix_test.go` - Unix系统环境变量管理测试
- `src/cmd/init.go` - 新的初始化命令
- `LINUX_CROSS_PLATFORM_SUPPORT.md` - 本文档

### 修改文件
- `src/internal/env/env_unix.go` - 完善Unix系统环境变量管理（重命名自env_linux.go）
- `src/internal/constants/constants.go` - 添加跨平台常量定义
- `src/internal/config/config.go` - 优化跨平台配置管理
- `src/internal/sys/system.go` - 增强系统工具类
- `src/internal/style/styles.go` - 添加Warning样式
- `src/internal/style/theme.go` - 添加Warning颜色主题
- `src/cmd/root.go` - 修改初始化逻辑，支持跨平台权限处理
- `README.md` - 更新文档，添加Linux使用说明

### 移除文件
- `src/internal/shell/shell_windows.go` - 移除不必要的Windows shell支持

## 主要功能实现

### 1. Linux环境变量管理 (`env_linux.go`)

#### 核心功能
- **完整的`doSetEnv`实现**：支持当前进程、系统级和用户级环境变量设置
- **多Shell环境支持**：自动检测并配置bash、zsh、fish、profile
- **智能权限处理**：根据用户权限选择系统级或用户级配置
- **环境变量持久化**：确保重启后配置仍然有效

#### 关键函数
```go
func doSetEnv(key, value string) error
func updateUserEnvironmentVariables(key, value string) error
func detectUserShells(homeDir string) []string
func updateShellEnvironment(homeDir, shell, key, value string) error
func SetSystemEnvPath() error
func SetCurrentUserEnvPath() error
```

### 2. 多Shell环境支持 (`shell/`)

#### 支持的Shell
- **Bash** (`.bashrc`): `export VAR="value"`
- **Zsh** (`.zshrc`): `export VAR="value"`
- **Fish** (`.config/fish/config.fish`): `set -gx VAR "value"`
- **Profile** (`.profile`): `export VAR="value"` (通用fallback)

#### 核心功能
- **自动检测**：扫描用户目录中的shell配置文件
- **智能配置**：为每种shell生成正确的语法
- **配置管理**：支持添加、更新、删除环境变量

### 3. 跨平台配置管理

#### 符号链接路径策略
- **Windows**: `C:\Java\JAVA_HOME`
- **Linux (root)**: `/opt/jenv/java_home`
- **Linux (user)**: `~/.jenv/java_home`
- **macOS (root)**: `/opt/jenv/java_home`
- **macOS (user)**: `~/.jenv/java_home`

#### 配置兼容性
- 统一的配置文件格式
- 平台特定的默认值
- 向后兼容性保证

### 4. 增强的系统工具类

#### 权限管理
```go
func IsAdmin() bool                    // 跨平台管理员检查
func RequireAdmin() error              // 智能权限要求
func isSystemPath(path string) bool    // 系统路径检测
```

#### 符号链接管理
- **智能权限检查**：只在需要时要求root权限
- **自动目录创建**：确保父目录存在
- **错误处理优化**：提供清晰的错误信息

### 5. 新的初始化流程

#### `jenv init` 命令
- **平台检测**：自动识别运行平台
- **权限提示**：清晰的权限要求说明
- **配置选择**：根据权限选择配置级别
- **用户指导**：提供下一步操作建议

#### 初始化逻辑
- Windows：要求管理员权限
- Linux：可选root权限，自动降级到用户级
- 配置验证和错误处理

## Linux平台使用说明

### 安装和构建

```bash
# 克隆项目
git clone https://github.com/WhyWhatHow/jenv.git
cd jenv/src

# 构建Linux版本
go build -o jenv

# 添加执行权限
chmod +x jenv
```

### 首次设置

```bash
# 1. 初始化jenv
./jenv init

# 2. 添加到PATH
./jenv add-to-path

# 3. 重新加载shell配置
source ~/.bashrc  # 或 ~/.zshrc, ~/.config/fish/config.fish
```

### 常用操作

```bash
# 扫描Java安装
./jenv scan /usr/lib/jvm
./jenv scan /opt

# 添加Java版本
./jenv add jdk8 /usr/lib/jvm/java-8-openjdk
./jenv add jdk11 /usr/lib/jvm/java-11-openjdk

# 列出可用版本
./jenv list

# 切换Java版本
./jenv use jdk11

# 查看当前版本
./jenv current
java -version
```

### 权限级别说明

#### Root用户运行
- 创建系统级配置：`/opt/jenv/java_home`
- 修改系统环境文件：`/etc/environment`
- 所有用户共享配置

#### 普通用户运行
- 创建用户级配置：`~/.jenv/java_home`
- 修改用户shell配置文件
- 仅当前用户可用

## 与Windows版本的差异

### 相同功能
- 基本的Java版本管理（add, list, use, remove, current）
- 符号链接机制
- 配置文件格式
- 命令行界面

### 平台差异

| 功能 | Windows | Linux |
|------|---------|-------|
| 权限要求 | 必须管理员 | 可选root |
| 符号链接位置 | `C:\Java\JAVA_HOME` | `/opt/jenv/` 或 `~/.jenv/` |
| 环境变量存储 | 注册表 | Shell配置文件 |
| Shell支持 | PowerShell, CMD | bash, zsh, fish |
| 配置持久化 | 系统注册表 | 配置文件 |

### 兼容性保证
- 配置文件格式完全兼容
- 命令行接口保持一致
- 可在不同平台间共享JDK配置

## 测试覆盖

### 单元测试
- `env_linux_test.go`: Linux环境变量管理测试
- `shell_test.go`: Shell环境支持测试
- `system_test.go`: 系统工具类测试（现有）

### 测试场景
- 多Shell环境检测和配置
- 权限级别处理
- 环境变量设置和清理
- 配置文件更新和回滚
- 符号链接创建和管理

## 已知限制和注意事项

### 当前限制
1. macOS支持仍在开发中
2. 某些Linux发行版可能需要额外配置
3. Fish shell需要手动安装才能被检测

### 注意事项
1. 首次使用需要运行`jenv init`
2. 环境变量更改需要重新加载shell
3. 系统级安装需要root权限
4. 建议定期备份配置文件

## 后续开发计划

### 短期目标
- [ ] 完善macOS支持
- [ ] 添加更多Linux发行版测试
- [ ] 优化错误处理和用户提示

### 长期目标
- [ ] 支持更多Shell环境
- [ ] 添加自动更新功能
- [ ] 集成包管理器安装
- [ ] Web界面管理工具

## 贡献指南

### 开发环境
- Go 1.21+
- 支持的操作系统：Windows 10+, Linux (主流发行版)
- 推荐IDE：VS Code with Go extension

### 测试要求
- 所有新功能必须包含单元测试
- 跨平台功能需要在多个平台测试
- 保持测试覆盖率 > 80%

### 代码规范
- 遵循Go官方代码规范
- 使用有意义的变量和函数名
- 添加适当的注释和文档
- 保持向后兼容性
