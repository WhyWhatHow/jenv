<div align="center">
<img src="assets/jenv-logo.png" width="200" height="200" alt="JEnv Logo">

# Jenv: Java 版本管理工具

![GitHub release](https://img.shields.io/github/v/release/WhyWhatHow/jenv)
![Build Status](https://img.shields.io/github/actions/workflow/status/WhyWhatHow/jenv/release.yml?branch=main)
![Version](https://img.shields.io/badge/version-v0.6.9-blue)

[English](README.md) | 中文 | [日本語](README_jp.md)

**🚀 [通过落地页快速开始](https://jenv-win.vercel.app)** - 一键下载，自动检测平台（Windows/Linux/macOS） • 多语言支持（EN/中文/日本語） • 每两周更新

</div>

## 最新更新 (v0.6.9)

### 🚀 性能提升
- **超快 JDK 扫描**：扫描时间从 3 秒缩短至 300ms（提升 90%）。[阅读更多技术细节](doc/PERFORMANCE_zh.md)
- **并发处理**：使用 goroutines 实现了 Dispatcher-Worker 模型
- **智能过滤**：积极的预过滤，跳过不必要的目录
- **进度追踪**：实时扫描进度和详细统计数据

### ✅ 跨平台支持
- **macOS 支持已完成**：完全兼容 macOS（Intel/Apple Silicon）
- **Linux 支持已完成**：支持多种 Shell 的完整 Linux 兼容性
- **Windows 优化**：增强了路径验证和兼容性修复

### 🔧 技术增强
- **Java 路径验证**：提高了 Windows JDK 检测的可靠性
- **环境管理**：优化了跨平台环境变量处理
- **配置清理**：移除了未使用选项，提高了代码可维护性

---

## 概述

`Jenv` 是一个用于管理系统中多个 Java 版本的命令行工具。它允许你轻松地在不同的 Java 版本之间切换，添加新的 Java 安装，并管理你的 Java 环境。

## 特性

### 高效的 Java 版本管理

- **基于符号链接的架构**
    - 通过符号链接快速切换版本
    - 一次性系统 PATH 配置
    - 更改在系统重启后依然有效
    - 在所有控制台窗口中立即生效

### 跨平台支持

- **Windows 支持**
    - 基于注册表的环境变量管理（Windows 标准）
    - 自动处理管理员权限
    - 遵循最小权限原则，减少 UAC 提示
    - 在 Windows 10/11 系统上表现优异

- **Linux/Unix 支持**
    - 基于 Shell 配置文件比例的环境管理
    - 提供用户级和系统级配置选项
    - 支持多种 Shell 环境（bash, zsh, fish）
    - 智能权限处理

- **macOS 支持**
    - 支持 Intel 和 Apple Silicon 架构
    - 支持 macOS 标准 JDK 路径

### 现代 CLI 体验

- **用户友好界面**
    - 直观的命令结构
    - 支持浅色/深色主题
    - 彩色输出，提高可读性
    - 详细的帮助文档

### 高级功能

- **智能 JDK 管理**
    - 全系统 JDK 扫描
    - **超快扫描性能（3s → 300ms）**，使用并发 Dispatcher-Worker 模型
    - 基于别名的 JDK 管理
    - 当前 JDK 状态追踪
    - 轻松添加和删除 JDK

---

## 安装

### 从 Release 下载
从 [Releases 页面](https://github.com/WhyWhatHow/jenv/releases)下载最新版本。

### 从源码构建

#### 前置条件

- Go 1.21 或更高版本
- Git
- **Windows**：创建系统符号链接需要管理员权限

#### 构建步骤

1. 克隆仓库：
```bash
git clone https://github.com/WhyWhatHow/jenv.git
cd jenv
```

2. 构建项目：

```bash
cd src
# Windows (PowerShell)
go build -o jenv.exe
# Linux/macOS
go build -o jenv
```

## 使用方法

### 首次设置

```bash
# 初始化 jenv（首次使用必选）
jenv init

# 将 jenv 添加到系统 PATH
jenv add-to-path
```

### 添加和删除 JDK

```bash
# 使用别名添加新的 JDK
jenv add jdk8 "C:\Program Files\Java\jdk1.8.0_291"
# 移除 JDK
jenv remove jdk8
```

### 切换 JDK 版本

```bash
# 列出所有已安装的 JDK
jenv list

# 切换到指定的 JDK 版本
jenv use jdk8

# 显示当前使用的 JDK
jenv current
```

### 扫描系统中的 JDK
```bash
# 自动检测已安装的 JDK
jenv scan c:\  # Windows
jenv scan /usr/lib/jvm  # Linux
```

---

## 为什么创建这个项目？

虽然 Linux 和 macOS 拥有像 `sdkman` 和原始 `jenv`（基于 bash）这样成熟的 Java 版本管理工具，但 Windows 用户一直以来只有有限且往往不够理想的选择。现有的 [JEnv-for-Windows](https://github.com/FelixSelter/JEnv-for-Windows) 等方案在现代 Windows 系统上可能会遇到明显的性能瓶颈。

本项目诞生于两个核心动力：

1.  **弥补 Windows 空白**：为 Windows 开发者提供一个稳健、高性能且与 Unix 类系统体验一致（甚至更好）的 Java 版本管理工具。
2.  **性能优先**：无论系统大小或复杂度如何，都能实现近乎瞬时的 JDK 扫描和切换。

我们的目标是成为 **Windows 上事实上的标准 Java 环境管理器**，同时为跨 Windows、Linux 和 macOS 工作的开发者提供无缝、统一的体验。

## 致谢

- [cobra](https://github.com/spf13/cobra) - 强大的 Go CLI 框架
- [jreleaser](https://jreleaser.org/) - 发布自动化工具
- [nvm-windows](https://github.com/coreybutler/nvm-windows) - 启发了我们的符号链接方案
- [Jenv-for-Windows](https://github.com/FelixSelter/JEnv-for-Windows) - Windows Java 版本管理的前辈项目

## 开源协议

本项目采用 Apache License 2.0 协议开源 - 详情请参阅 [LICENSE](LICENSE) 文件。
