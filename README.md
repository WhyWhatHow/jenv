<div align="center">
<img src="assets/jenv-logo.png" width="200" height="200" alt="JEnv Logo">

# Jenv: Java Environment Manager

![GitHub release](https://img.shields.io/github/v/release/WhyWhatHow/jenv)
![Build Status](https://img.shields.io/github/actions/workflow/status/WhyWhatHow/jenv/release.yml?branch=main)
</div>

## Overview

`Jenv` is a command-line tool for managing multiple Java versions on your system. It allows you to easily switch between
different Java versions, add new Java installations, and manage your Java environment.

## Features

### Efficient Java Version Management

- **Symlink-Based Architecture**
    - Fast version switching through symbolic links
    - One-time system PATH configuration
    - Changes persist across system reboots
    - Instant effect in all console windows

### Windows-First Design

- **Optimized for Windows**
    - Automatic administrator privilege handling
    - Minimized UAC prompts with least privilege principle
    - Superior performance on Windows 10/11 systems

### Modern CLI Experience

- **User-Friendly Interface**
    - Intuitive command structure
    - Light/Dark theme support
    - Colorful output for improved readability
    - Detailed help documentation

### Advanced Features

- **Smart JDK Management**
    - System-wide JDK scanning
    - Alias-based JDK management
    - Current JDK status tracking
    - Easy JDK addition and removal

### Future-Ready

- **Cross-Platform Support**
    - Windows support (Current)
    - Linux support (Implemented)
    - macOS support (Coming soon)
## Project Structure

```
.
├── cmd/            # Command implementations
│   ├── add.go      # Add JDK command
│   ├── list.go     # List JDKs command
│   ├── ...
│   └── root.go     # Root command and flags
├── internal/       # Internal packages
│   ├── config/     # Configuration management
│   ├── ...
│   └── sys/        # System utilities
├── main.go         # Main entry point
├── go.mod          # Go module definition
├── go.sum          # Go module checksums
├── doc/            # Documentation
├── assets/         # Assets like images, gifs
├── index.html      # Landing page (if applicable)
└── .github/        # GitHub configurations
    └── workflows/      # CI/CD workflows
```

## Installation

### From Release
Download the latest release from the [Releases page](https://github.com/WhyWhatHow/jenv/releases).

### Build from Source

#### Prerequisites

- Go 1.21 or higher
- Git
- Windows systems require [Administrator privileges](#symbolic-link-permissions) (for creating system symbolic links)

#### Build Steps

1. Clone the repository:
```bash
git clone https://github.com/WhyWhatHow/jenv.git
cd jenv
```

2. Build the project:

```bash
# For Windows (PowerShell)
go build -ldflags "-X github.com/whywhathow/jenv/cmd.Version=1.0.0" -o jenv.exe main.go

# For Linux/macOS
go build -ldflags "-X github.com/whywhathow/jenv/cmd.Version=1.0.0" -o jenv main.go

# For development build (with debug information)
go build -o jenv main.go
```

## Usage

![jenv.gif](assets/jenv.gif)

### Installation Verification

```bash
# Verify jenv installation
jenv --version
```

## Linux Specifics

On Linux, `jenv` provides a user-friendly way to manage Java versions. Here are some key aspects:

- **Symlink for `JAVA_HOME`**: `jenv` manages the active Java version by pointing a symbolic link to the chosen JDK installation. By default, this symlink is located at `~/.jenv/java_home`. This path is then typically used to set your `JAVA_HOME` environment variable.

- **Environment Variable Setup**:
    - When `jenv add-to-path` is run (for adding `jenv` itself to the PATH) or during the initial setup triggered by the first `jenv` command, `jenv` attempts to make the `JAVA_HOME/bin` directory (via the symlink) available in your shell's `PATH`.
    - If `jenv` commands that modify environment variables (like `add-to-path` or the initial setup) are run with `sudo`, `jenv` will attempt to create or update a script in `/etc/profile.d/jenv.sh`. This makes the settings available system-wide. For this to take effect, your shell needs to be configured to source scripts from this directory (most default shell configurations do).
    - Without `sudo`, or if system-wide modification fails (e.g., due to permissions), `jenv` will fall back to updating user-specific shell configuration files. It currently supports `~/.bashrc` (for bash shells) and `~/.zshrc` (for zsh shells), falling back to `~/.profile` for other shells.
    - After these files are modified, you may need to source the relevant file (e.g., `source ~/.bashrc`) or open a new terminal session for the changes to take effect.

- **No Root Needed (Usually)**: For most common operations like adding, removing, listing, or switching between JDKs (`jenv add`, `jenv remove`, `jenv list`, `jenv use`), `sudo` is **not** required, provided `jenv` is installed in a user-writable location and uses the default symlink path (`~/.jenv/java_home`).

### Add and remove JDK

![jenv-add.gif](assets/jenv-add.gif)

```bash
# Add a new JDK with an alias name
#jenv add <alias> <jdk_path>
jenv add jdk8 "C:\Program Files\Java\jdk1.8.0_291"
#jenv remove <alias>
jenv remove jdk8
```

### List all installed JDKs

```bash
jenv list
```

### Switch to a specific JDK version

```bash
#jenv use <alias>
jenv use jdk8
```

### Remove a JDK from jenv

```bash
#jenv remove <alias>
jenv remove jdk8
```

### Show current JDK in use

```bash
jenv current
```

### Scan system for installed JDKs
```bash
#jenv scan <path>
jenv scan c:\
```

### Add jenv to system PATH

```bash
jenv add-to-path
```

### Change UI theme (light/dark)

```bash
#jenv theme <theme_name>
jenv theme dark
```

### help & version
```bash
#jenv help [command]
jenv --version
```

## Q&A

### Administrative Privileges: Windows vs. Linux

**Windows:**
Due to Windows system restrictions, creating system-level symbolic links (which `jenv` uses by default at `C:\Java\JAVA_HOME`) requires administrator privileges. `jenv` attempts to handle UAC prompts automatically when needed. You'll typically need to run `jenv` in an administrator PowerShell or Command Prompt for operations that modify this symlink or system PATH.

**Linux:**
On Linux, `sudo` (root privileges) is generally **not** required for typical user-level management if:
- `jenv` itself is installed in a user-writable location (e.g., within your home directory).
- The default `JAVA_HOME` symlink path (`~/.jenv/java_home`) is used.
- You are managing JDKs installed in user-writable locations.

`sudo` would only be needed if you wish to:
- Install `jenv` to a system-wide directory (e.g., `/usr/local/bin`).
- Use commands like `jenv add-to-path` to create system-wide configurations (e.g., in `/etc/profile.d/jenv.sh`).
- Configure `jenv` to manage a `JAVA_HOME` symlink located in a system-protected directory.

For most users, installing and running `jenv` without `sudo` is the recommended approach on Linux.

### Why was this project created?

While Linux and macOS users have mature tools like `sdkman` and `jenv` for Java version management, Windows users have
limited options. The existing [Jenv-forWindows](https://github.com/FelixSelter/JEnv-for-Windows) solution, while
functional, faces performance issues on Windows 10 systems.

This project was born out of two motivations:

1. To create a fast, efficient Java version manager specifically optimized for Windows
2. To explore AI-assisted development using tools like `cursor` and `Trae` while learning Go programming from scratch

The goal is to provide Windows developers with a robust, performant solution for managing multiple Java environments,
similar to what Linux and macOS users already enjoy.

### How it works?

Inspired by nvm-windows, JEnv uses symlinks for Java version management, which offers several advantages:

1. **Symlink-Based Architecture**
    - Creates a single symlink (default `C:\java\JAVA_HOME` on Windows, `~/.jenv/java_home` on Linux/macOS) during initial setup. This path can be configured.
    - Switching Java versions only requires updating the symlink target.
    - No need to modify system PATH repeatedly for `JAVA_HOME` itself.
    - Changes persist across system reboots and apply to all console windows

2. **Implementation Details**
    - During initialization (first run or relevant setup commands):
        - Establishes the `JAVA_HOME` symlink (e.g., at `C:\java\JAVA_HOME` or `~/.jenv/java_home`).
        - Adds the corresponding `bin` directory (e.g., `JAVA_HOME\bin` on Windows, `JAVA_HOME/bin` on Linux/macOS) to the system PATH. This is typically a one-time setup for persistence across sessions.
        - Creates an initial symlink if a default JDK is set up or added.
    - When switching versions (`jenv use <alias>`):
        - Simply updates the `JAVA_HOME` symlink to point to the desired JDK's installation path.
        - No further PATH modifications are generally needed for `JAVA_HOME` itself, as the PATH already points to the symlinked `JAVA_HOME/bin`.
        - Changes to `JAVA_HOME` (via the symlink) are effectively immediate for new processes. Existing shells might need to be restarted or have their `JAVA_HOME` variable refreshed if they cached its value at startup.

This approach is more efficient than constantly modifying system PATH variables directly for each Java version, providing a cleaner and more reliable solution for Java version management.

## Acknowledgments

- [cobra](https://github.com/spf13/cobra) - A powerful CLI framework for Go
- [jreleaser](https://jreleaser.org/) - A release automation tool
- [nvm-windows](https://github.com/coreybutler/nvm-windows) - Inspired our symlink-based approach
- [Jenv-for-Windows](https://github.com/FelixSelter/JEnv-for-Windows) - A predecessor project for Java version
  management on Windows

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

