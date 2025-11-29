# Zero to Hello World in 5 Minutes with JEnv

## Your Time is Precious

In my previous article, "Why I Built JEnv: An Umbrella for Java Beginners," I shared the story of creating jenv—to spare Java beginners from environment configuration nightmares.

I understand your pain points.

Now, I'm here with the solution.

**Promise: 5 minutes to complete setup, then start coding immediately.**

No manual environment variable editing. No computer reboots. No fear of breaking your system.

Just follow this tutorial step by step, and you'll have a **professional, efficient, and user-friendly** Java development environment.

Let's get started!

---

## Preparation: Visit the JEnv Landing Page

JEnv provides a dedicated Landing Page where you can **download JDK + jenv in one click**, eliminating the hassle of hunting for JDK installers.

🌐 **Landing Page**: [jenv-win.vercel.app](https://jenv-win.vercel.app)

### Landing Page Features:

✅ **Auto-detect platform**: Automatically recognizes Windows, Linux, or macOS
✅ **Multiple versions**: Choose from Java 8, 11, 17, or 21
✅ **One-stop download**: Get JDK + jenv together—no separate searches needed
✅ **Multi-language support**: Switch between English and Chinese
✅ **Regular updates**: JDK download links updated bi-weekly

---

## Step 1: Download JEnv and JDK

### 1.1 Visit the Landing Page

Open your browser and go to [jenv-win.vercel.app](https://jenv-win.vercel.app)

The page will auto-detect your operating system and display the appropriate download options.

### 1.2 Choose JDK Version

Select a JDK version based on your needs:

- **Java 8**: Common for university courses and legacy projects
- **Java 11**: LTS (Long Term Support) version
- **Java 17**: Latest LTS version, recommended for beginners
- **Java 21**: Latest features for early adopters

💡 **Beginner Tip**: If unsure, go with Java 17!

### 1.3 Download Files

Click the download button to get:
1. **jenv** executable
2. **JDK** installer

**Windows users**: Download `jenv-windows-amd64.exe`
**Linux users**: Download `jenv-linux-amd64`

### 1.4 Install JDK (Optional)

If you don't have any JDK installed yet, install the downloaded one:

**Windows**: Double-click the `.exe` or `.msi` file and follow the installer prompts
**Linux**: Extract the `.tar.gz` file to `/usr/lib/jvm` or `/opt`

💡 **Tip**: If you already have JDK installed, skip this step—jenv will auto-detect it!

---

## Step 2: Initialize JEnv

### 2.1 Move JEnv to a Suitable Location

**Windows**:
```bash
# Recommended locations: C:\Program Files\ or C:\tools\
move jenv-windows-amd64.exe C:\tools\jenv.exe
```

**Linux**:
```bash
# Move to /usr/local/bin and grant execute permission
sudo mv jenv-linux-amd64 /usr/local/bin/jenv
sudo chmod +x /usr/local/bin/jenv
```

### 2.2 Run Initialization Command

Open a terminal (CMD, PowerShell, or Bash) and run:

```bash
jenv init
```

**What happens behind the scenes**:
- Creates config file `~/.jenv/config.json`
- Creates symlink directory:
  - Windows: `C:\java\JAVA_HOME`
  - Linux: `/opt/jenv/java_home` (system-level) or `~/.jenv/java_home` (user-level)
- Prepares environment variable configuration

⚠️ **Windows users**: You'll see a UAC (User Account Control) prompt asking for permission to modify the system. **Click "Yes"**—this is required to create symbolic links.

⚠️ **Linux users**: For system-level installation, use `sudo jenv init`. Without sudo, jenv will automatically use user-level configuration (recommended for beginners—no root privileges needed).

### 2.3 Add JEnv to PATH

```bash
jenv add-to-path
```

**Windows**: This command automatically adds `C:\java\JAVA_HOME\bin` to the system PATH
**Linux**: This command updates your shell config file (.bashrc/.zshrc/.config/fish/config.fish)

💡 **Linux users**: After execution, run the following to apply changes:
```bash
# bash users
source ~/.bashrc

# zsh users
source ~/.zshrc

# fish users
source ~/.config/fish/config.fish
```

---

## Step 3: Scan and Add JDKs

### 3.1 Auto-scan for JDKs

JEnv can automatically scan your computer for all installed JDKs in **just 300ms**!

**Windows**:
```bash
jenv scan C:\
```

**Linux**:
```bash
jenv scan /usr/lib/jvm
jenv scan /opt
```

**Sample output**:
```
🔍 Scanning for JDK installations...
✓ Found: Java 8 at C:\Program Files\Java\jdk1.8.0_291
✓ Found: Java 11 at C:\Program Files\Java\jdk-11.0.12
✓ Found: Java 17 at C:\Program Files\Java\jdk-17.0.2

📊 Scan completed in 320ms
🎉 Found 3 JDK installations
```

⚡ **Performance highlight**: JEnv uses **Go's concurrency** (goroutines), like dispatching 10 people to search simultaneously—that's why it's blazing fast!

Traditional tools take 3+ seconds; jenv only needs 300ms. That's the secret to **10x speed improvement**.

![jenv-add.gif](../../assets/jenv-add.gif)

### 3.2 Manually Add JDKs

If the scan didn't find a JDK, or you want to add a specific one, manually add it:

```bash
jenv add java8 "C:\Program Files\Java\jdk1.8.0_291"
jenv add java11 "C:\Program Files\Java\jdk-11.0.12"
jenv add java17 "C:\Program Files\Java\jdk-17.0.2"
```

💡 **Naming tip**: Use simple names like `java8`, `java11`, `java17` for easy recall and switching.

---

## Step 4: Switch Java Versions

Now, the magic happens! Switch Java versions with a single command:

```bash
jenv use java17
```

**Output**:
```
✅ Switched to Java 17
📁 JAVA_HOME: C:\java\JAVA_HOME -> C:\Program Files\Java\jdk-17.0.2
🎯 Current version: java17
```

### Verify the switch:

```bash
java -version
```

**Output**:
```
java version "17.0.2" 2022-01-18 LTS
Java(TM) SE Runtime Environment (build 17.0.2+8-LTS-86)
Java HotSpot(TM) 64-Bit Server VM (build 17.0.2+8-LTS-86, mixed mode, sharing)
```

🎉 **That's it!** No terminal restart. No computer reboot. **Instant effect**!

### Switch to other versions:

```bash
jenv use java8
java -version
# Output: java version "1.8.0_291"

jenv use java11
java -version
# Output: java version "11.0.12"
```

---

## Step 5: List All Added JDKs

Want to see which JDKs jenv is managing?

```bash
jenv list
```

**Output**:
```
📋 Installed Java versions:
  java8    -> C:\Program Files\Java\jdk1.8.0_291
  java11   -> C:\Program Files\Java\jdk-11.0.12
* java17   -> C:\Program Files\Java\jdk-17.0.2
           (* = current)

💡 Use 'jenv use <name>' to switch versions
```

The `*` indicates the currently active version.

---

## Real-World Scenarios: Multi-Project Management

### Scenario 1: University Course + Personal Learning

Your university Java course requires Java 8, but the Spring Boot tutorial you're following on YouTube uses Java 17.

**Traditional approach**: Manually edit environment variables, restart IDE—pure torture.

**JEnv approach**:
```bash
# Morning: University assignment
cd ~/school/java-homework
jenv use java8
javac HelloWorld.java

# Afternoon: Learning Spring Boot
cd ~/learning/spring-boot-tutorial
jenv use java17
./mvnw spring-boot:run
```

**Zero pain, seamless switching!**

### Scenario 2: Testing Java 21 Features

You heard Java 21 has Virtual Threads and want to try them, but you don't want to mess up your main project environment.

```bash
# Add Java 21
jenv add java21 "C:\Program Files\Java\jdk-21"

# Switch to Java 21 for experimentation
jenv use java21
java --enable-preview VirtualThreadTest.java

# Done experimenting? Switch back to Java 17
jenv use java17
```

**Zero-risk exploration!**

### Scenario 3: Reproducing User Bugs

A user reports a bug on Java 8, but it works fine on Java 11.

```bash
# Test on Java 8
jenv use java8
./run-tests.sh

# Test on Java 11
jenv use java11
./run-tests.sh
```

**Quick problem diagnosis!**

---

## Advanced Tips

### Check Current JDK

```bash
jenv current
```

**Output**:
```
Current Java version: java17
Path: C:\Program Files\Java\jdk-17.0.2
```

### Switch Theme (Dark Mode for Eye Comfort)

```bash
jenv theme dark   # Dark theme
jenv theme light  # Light theme
```

### Remove Unwanted JDKs

```bash
jenv remove java8
jenv list  # Verify removal
```

💡 **Safety note**: JEnv only removes the management reference—**it doesn't delete actual Java files**, so it's completely safe!

---

## FAQ: Common Questions

### Q1: Do I need to uninstall existing Java?

**A**: No! JEnv is fully compatible with existing Java installations. It only helps you manage and switch—it doesn't break anything.

### Q2: Will it break my system?

**A**: No. JEnv only modifies environment variables (with backup) and doesn't touch Java installation files. Even if something goes wrong, a simple reboot will restore everything.

### Q3: Must I use the command line? I'm not used to it...

**A**: Currently, yes, but the commands are super simple:
- `jenv add` - Add
- `jenv use` - Switch
- `jenv list` - View
- `jenv remove` - Remove
- `jenv current` - Current version

Just 5 core commands! A GUI version is on the roadmap—stay tuned.

### Q4: Why does Windows need admin privileges?

**A**: Only during `jenv init` to create symbolic links—Windows requires admin privileges for this. **After that, all operations are privilege-free**!

It's a one-time setup, then you're free.

### Q5: Which Java versions are supported?

**A**: All of them! From Java 6 to Java 21, any vendor's JDK is supported:
- Oracle JDK
- OpenJDK
- Amazon Corretto
- Azul Zulu
- GraalVM
- ...

As long as it's a standard JDK, jenv can manage it!

### Q6: Must Linux users use sudo?

**A**: Not necessarily!
- **With sudo**: System-level installation, shared by all users
- **Without sudo**: User-level installation, affects only current user

**Beginner recommendation**: No sudo—user-level installation is safer and doesn't require root privileges.

---

## Summary: What You've Learned

🎓 **In 5 minutes, you learned**:

1. ✅ Download jenv + JDK from Landing Page
2. ✅ Initialize jenv, configure environment
3. ✅ Scan and add JDKs (300ms ultra-fast)
4. ✅ Switch Java versions with one command
5. ✅ Manage multiple JDKs for different scenarios
6. ✅ Use advanced features and themes

**Now you can focus on learning Java itself, instead of fighting with environment configuration!**

---

## Start Your Java Journey

🚀 **Get started now**:
- Landing Page: [jenv-win.vercel.app](https://jenv-win.vercel.app)
- GitHub: [github.com/WhyWhatHow/jenv](https://github.com/WhyWhatHow/jenv)
- Full docs: [Usage Guide](https://github.com/WhyWhatHow/jenv#usage)

⭐ **Found this helpful?**
- Star on GitHub—let more people discover jenv
- Share with classmates and friends
- Comment with your experience

💬 **Questions?**
- Open an Issue: [GitHub Issues](https://github.com/WhyWhatHow/jenv/issues)
- I'll do my best to help every beginner in need

---

**Happy coding—may your Java development journey be smooth sailing! ✨**

---

> **Coming Next**: _Why JEnv is 10x Faster: A Technical Deep Dive_
>
> Unveiling jenv's technical implementation: symbolic links, Go concurrency, Dispatcher-Worker model...
