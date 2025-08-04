# JEnv 用户指南

## 安装指南

### 方式一：从发布版本安装
1. 访问 [Releases页面](https://github.com/WhyWhatHow/jenv/releases)
2. 下载最新版本的JEnv

### 方式二：从源码编译
前置要求：
- Go 1.21或更高版本
- Git
- Windows系统需要管理员权限

编译步骤：
```bash
git clone https://github.com/WhyWhatHow/jenv.git
cd jenv/src
go build -o jenv
```

## 基本命令

### 初始化
```bash
jenv init                # 初始化JEnv配置
```

### JDK管理
```bash
jenv add jdk8 "路径"     # 添加JDK并指定别名
jenv list               # 查看已安装的JDK列表
jenv use jdk8           # 切换到指定JDK版本
jenv remove jdk8        # 移除指定JDK
jenv current            # 显示当前使用的JDK
```

### 系统配置
```bash
jenv add-to-path        # 将JEnv添加到系统PATH
jenv update             # 更新环境变量
```

## 权限要求（Windows）

### 管理员权限
Windows系统下需要管理员权限来创建系统级符号链接：
1. 以管理员身份运行PowerShell
2. 开启开发者模式（Windows 10+）
3. 首次运行时接受UAC提示

## 常见问题

### Q: 为什么需要管理员权限？
A: Windows系统限制要求创建系统级符号链接必须具有管理员权限。JEnv使用符号链接来管理Java版本，这种方式比反复修改系统PATH更高效。

### Q: 如何验证安装是否成功？
A: 运行 `jenv --version` 命令检查版本信息。

### Q: 命令执行时提示权限不足怎么办？
A: 确保使用管理员权限运行终端，并在UAC提示时点击"是"。