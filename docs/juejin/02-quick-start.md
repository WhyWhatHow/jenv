# 5 分钟从零到 Hello World:jenv 极速上手

## 你的时间很宝贵

在上一篇文章《曾经淋过雨,现在想为你撑伞》中,我讲述了创建 jenv 的初心——让 Java 新手不再被环境配置折磨。

你的痛点,我都懂。

现在,我要给你解决方案。

**承诺:5 分钟完成配置,马上开始写代码。**

不需要手动改环境变量,不需要重启电脑,不需要担心把系统搞坏。

只需要跟着这篇教程,一步一步来,你就能拥有一个**专业、高效、好用**的 Java 开发环境。

让我们开始吧!

---

## 准备工作:访问 jenv Landing Page

jenv 提供了一个专门的 Landing Page,让你可以**一键下载 JDK + jenv**,省去到处找 JDK 安装包的麻烦。

🌐 **Landing Page 地址**: [jenv-win.vercel.app](https://jenv-win.vercel.app)

### Landing Page 的特点:

✅ **自动检测平台**:自动识别你是 Windows、Linux 还是 macOS
✅ **多版本选择**:Java 8、11、17、21 任你挑
✅ **一站式下载**:JDK + jenv 一起下载,不用分开找
✅ **多语言支持**:中文/英文界面切换
✅ **定期更新**:每两周自动更新 JDK 下载链接

---

## Step 1:下载 jenv 和 JDK

### 1.1 访问 Landing Page

打开浏览器,访问 [jenv-win.vercel.app](https://jenv-win.vercel.app)

页面会自动检测你的操作系统,并显示对应的下载选项。

### 1.2 选择 JDK 版本

根据你的需求选择 JDK 版本:

- **Java 8**:大学课程、老项目常用
- **Java 11**:LTS 长期支持版本
- **Java 17**:最新 LTS 版本,推荐新手学习
- **Java 21**:最新特性,尝鲜用

💡 **新手建议**:如果不确定,选 Java 17 就对了!

### 1.3 下载文件

点击下载按钮,下载:
1. **jenv** 可执行文件
2. **JDK** 安装包

**Windows 用户**:下载 `jenv-windows-amd64.exe`
**Linux 用户**:下载 `jenv-linux-amd64`

### 1.4 安装 JDK(可选)

如果你还没有安装任何 JDK,先安装下载好的 JDK:

**Windows**:双击 `.exe` 或 `.msi` 文件,按提示安装
**Linux**:解压 `.tar.gz` 文件到 `/usr/lib/jvm` 或 `/opt`

💡 **提示**:如果你电脑上已经有 JDK 了,可以跳过这步,jenv 会自动扫描找到它们!

---

## Step 2:初始化 jenv

### 2.1 移动 jenv 到合适的位置

**Windows**:
```bash
# 建议放到 C:\Program Files\ 或 C:\tools\
move jenv-windows-amd64.exe C:\tools\jenv.exe
```

**Linux**:
```bash
# 移动到 /usr/local/bin 并赋予执行权限
sudo mv jenv-linux-amd64 /usr/local/bin/jenv
sudo chmod +x /usr/local/bin/jenv
```

### 2.2 运行初始化命令

打开终端(CMD、PowerShell 或 Bash),运行:

```bash
jenv init
```

**背后发生了什么**:
- 创建配置文件 `~/.jenv/config.json`
- 创建符号链接目录:
  - Windows: `C:\java\JAVA_HOME`
  - Linux: `/opt/jenv/java_home` (系统级) 或 `~/.jenv/java_home` (用户级)
- 准备环境变量配置

⚠️ **Windows 用户注意**:你会看到一个 UAC(用户账户控制)提示,询问是否允许程序修改系统。**点击"是"**,这是因为创建符号链接需要管理员权限。

⚠️ **Linux 用户注意**:如果你想要系统级安装,需要使用 `sudo jenv init`。如果不用 sudo,jenv 会自动使用用户级配置(推荐新手使用用户级配置,不需要 root 权限)。

### 2.3 添加 jenv 到 PATH

```bash
jenv add-to-path
```

**Windows**:这个命令会自动将 `C:\java\JAVA_HOME\bin` 添加到系统 PATH
**Linux**:这个命令会自动更新你的 shell 配置文件(.bashrc/.zshrc/.config/fish/config.fish)

💡 **Linux 用户**:执行完后,运行以下命令让配置生效:
```bash
# bash 用户
source ~/.bashrc

# zsh 用户
source ~/.zshrc

# fish 用户
source ~/.config/fish/config.fish
```

---

## Step 3:扫描并添加 JDK

### 3.1 自动扫描 JDK

jenv 可以自动扫描你电脑上所有已安装的 JDK,**只需 300ms**!

**Windows**:
```bash
jenv scan C:\
```

**Linux**:
```bash
jenv scan /usr/lib/jvm
jenv scan /opt
```

**输出示例**:
```
🔍 Scanning for JDK installations...
✓ Found: Java 8 at C:\Program Files\Java\jdk1.8.0_291
✓ Found: Java 11 at C:\Program Files\Java\jdk-11.0.12
✓ Found: Java 17 at C:\Program Files\Java\jdk-17.0.2

📊 Scan completed in 320ms
🎉 Found 3 JDK installations
```

⚡ **性能亮点**:jenv 使用 **Go 的并发处理**(goroutines),像派 10 个人同时搜索,所以速度超快!

传统工具需要 3+ 秒,jenv 只需 300ms。这就是 **10 倍速度提升**的秘密。

![jenv-add.gif](../../assets/jenv-add.gif)

### 3.2 手动添加 JDK

如果 scan 没找到,或者你想添加特定的 JDK,可以手动添加:

```bash
jenv add java8 "C:\Program Files\Java\jdk1.8.0_291"
jenv add java11 "C:\Program Files\Java\jdk-11.0.12"
jenv add java17 "C:\Program Files\Java\jdk-17.0.2"
```

💡 **命名建议**:使用简单的名称,比如 `java8`、`java11`、`java17`,方便记忆和切换。

---

## Step 4:切换 Java 版本

现在,魔法来了!切换 Java 版本,只需一条命令:

```bash
jenv use java17
```

**输出**:
```
✅ Switched to Java 17
📁 JAVA_HOME: C:\java\JAVA_HOME -> C:\Program Files\Java\jdk-17.0.2
🎯 Current version: java17
```

### 验证切换成功:

```bash
java -version
```

**输出**:
```
java version "17.0.2" 2022-01-18 LTS
Java(TM) SE Runtime Environment (build 17.0.2+8-LTS-86)
Java HotSpot(TM) 64-Bit Server VM (build 17.0.2+8-LTS-86, mixed mode, sharing)
```

🎉 **就是这么简单!** 无需重启终端,无需重启电脑,**立即生效**!

### 切换到其他版本:

```bash
jenv use java8
java -version
# 输出: java version "1.8.0_291"

jenv use java11
java -version
# 输出: java version "11.0.12"
```

---

## Step 5:查看所有已添加的 JDK

想看看 jenv 管理了哪些 JDK?

```bash
jenv list
```

**输出**:
```
📋 Installed Java versions:
  java8    -> C:\Program Files\Java\jdk1.8.0_291
  java11   -> C:\Program Files\Java\jdk-11.0.12
* java17   -> C:\Program Files\Java\jdk-17.0.2
           (* = current)

💡 Use 'jenv use <name>' to switch versions
```

`*` 表示当前正在使用的版本。

---

## 实战场景:多项目管理

### 场景 1:大学课程 + 个人学习

你的大学 Java 课程要求用 Java 8,但你在 YouTube 上学习的 Spring Boot 教程用的是 Java 17。

**传统做法**:手动改环境变量,重启 IDE,痛苦不堪。

**jenv 做法**:
```bash
# 上午:做大学作业
cd ~/school/java-homework
jenv use java8
javac HelloWorld.java

# 下午:学习 Spring Boot
cd ~/learning/spring-boot-tutorial
jenv use java17
./mvnw spring-boot:run
```

**零痛苦,无缝切换!**

### 场景 2:测试 Java 21 新特性

你听说 Java 21 有虚拟线程(Virtual Threads),想试试,但又不想破坏主项目环境。

```bash
# 添加 Java 21
jenv add java21 "C:\Program Files\Java\jdk-21"

# 切换到 Java 21 试一试
jenv use java21
java --enable-preview VirtualThreadTest.java

# 试完了,切回 Java 17
jenv use java17
```

**零风险尝鲜!**

### 场景 3:复现用户 Bug

用户报告说你的应用在 Java 8 上有 Bug,但在 Java 11 上正常。

```bash
# 在 Java 8 上测试
jenv use java8
./run-tests.sh

# 在 Java 11 上测试
jenv use java11
./run-tests.sh
```

**快速定位问题!**

---

## 进阶技巧

### 查看当前使用的 JDK

```bash
jenv current
```

**输出**:
```
Current Java version: java17
Path: C:\Program Files\Java\jdk-17.0.2
```

### 切换主题(护眼深色模式)

```bash
jenv theme dark   # 深色主题
jenv theme light  # 浅色主题
```

### 删除不需要的 JDK

```bash
jenv remove java8
jenv list  # 验证已删除
```

💡 **安全提示**:jenv 只删除管理引用,**不会删除实际的 Java 文件**,所以完全安全!

---

## FAQ:常见问题

### Q1:需要卸载现有的 Java 吗?

**A**:不需要!jenv 与现有 Java 安装完全兼容,只是帮你管理和切换,不会破坏任何东西。

### Q2:会破坏我的系统吗?

**A**:不会。jenv 只修改环境变量(有备份),不触碰 Java 安装文件。即使出问题,重启电脑就能恢复。

### Q3:必须用命令行吗?我不习惯...

**A**:目前是的,但命令超级简单:
- `jenv add` - 添加
- `jenv use` - 切换
- `jenv list` - 查看
- `jenv remove` - 删除
- `jenv current` - 当前版本

就这 5 个核心命令!GUI 版本在路线图中,敬请期待。

### Q4:为什么 Windows 需要管理员权限?

**A**:只有在 `jenv init` 时需要创建符号链接,Windows 系统要求管理员权限。**之后所有操作都不需要管理员权限了**!

这是一次性的设置,之后你就自由了。

### Q5:支持哪些 Java 版本?

**A**:全部!从 Java 6 到 Java 21,任何供应商的 JDK 都支持:
- Oracle JDK
- OpenJDK
- Amazon Corretto
- Azul Zulu
- GraalVM
- ...

只要是标准的 JDK,jenv 都能管理!

### Q6:Linux 一定要用 sudo 吗?

**A**:不一定!
- **用 sudo**:系统级安装,所有用户共享
- **不用 sudo**:用户级安装,只影响当前用户

**新手推荐**:不用 sudo,用户级安装更安全,不需要 root 权限。

---

## 总结:你已经学会了什么

🎓 **5 分钟内,你学会了**:

1. ✅ 从 Landing Page 下载 jenv + JDK
2. ✅ 初始化 jenv,配置环境
3. ✅ 扫描并添加 JDK(300ms 极速)
4. ✅ 一键切换 Java 版本
5. ✅ 管理多个 JDK,应对不同场景
6. ✅ 使用进阶技巧和主题

**现在,你可以专注于学习 Java 本身,而不是被环境配置折磨了!**

---

## 开始你的 Java 之旅

🚀 **立即开始**:
- Landing Page: [jenv-win.vercel.app](https://jenv-win.vercel.app)
- GitHub: [github.com/WhyWhatHow/jenv](https://github.com/WhyWhatHow/jenv)
- 完整文档: [Usage Guide](https://github.com/WhyWhatHow/jenv#usage)

⭐ **觉得有帮助?**
- 给个 Star,让更多人看到
- 分享给你的同学和朋友
- 在评论区分享你的使用体验

💬 **遇到问题?**
- 提 Issue: [GitHub Issues](https://github.com/WhyWhatHow/jenv/issues)
- 我会尽力帮助每一位遇到困难的新手

---

**Happy coding, 愿你的 Java 开发之路一帆风顺!✨**

---

> **下一篇预告**:《为什么 jenv 比 jenv-for-windows 快 10 倍?》
>
> 揭秘 jenv 的技术实现:符号链接、Go 并发、Dispatcher-Worker 模型...
