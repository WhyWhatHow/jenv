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

- **Cross-Platform Support (Planned)**
    - Windows support (Current)
    - Linux support (Coming soon)
    - macOS support (Coming soon)
## Project Structure

```
.
├── src/                # Source code directory
│   ├── cmd/            # Command implementations
│   │   ├── add.go      # Add JDK command
│   │   ├── list.go     # List JDKs command
│   │   ├── remove.go   # Remove JDK command
│   │   ├── use.go      # Switch JDK command
│   │   └── root.go     # Root command and flags
│   ├── internal/       # Internal packages
│   │   ├── config/     # Configuration management
│   │   ├── constants/  # Constants definitions
│   │   ├── env/        # Environment handling
│   │   ├── java/       # Java SDK management
│   │   ├── logging/    # Logging utilities
│   │   ├── style/      # UI styling
│   │   └── sys/        # System utilities
│   └── jenv.go         # Main entry point
├── doc/                # Documentation
└── .github/            # GitHub configurations
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

cd src 

# For Windows (PowerShell)
go build -ldflags "-X github.com/whywhathow/jenv/cmd.Version=1.0.0" -o jenv

# For Linux/macOS
go build -ldflags "-X github.com/whywhathow/jenv/cmd.Version=1.0.0" -o jenv 

# For development build (with debug information)
go build -o jenv 
```

## Usage

![jenv.gif](assets/jenv.gif)

### Installation Verification

```bash
# Verify jenv installation
jenv --version

```

### Add and remove JDK

![jenv-add.gif](assets/jenv-add.gif)

```bash
# Add a new JDK with an alias name
jenv add <alias> <jdk_path>
jenv add jdk8 "C:\Program Files\Java\jdk1.8.0_291"
jenv remove <alias>
jenv remove jdk8
```

### List all installed JDKs

```bash
jenv list
```

### Switch to a specific JDK version

```bash
jenv use <alias>
jenv use jdk8
```

### Remove a JDK from jenv

```bash
jenv remove <alias>
jenv remove jdk8
```

### Show current JDK in use

```bash
jenv current
```

### Scan system for installed JDKs
```bash
jenv scan c:\
```

### Add jenv to system PATH

```bash
jenv add-to-path
```

### Change UI theme (light/dark)

```bash
jenv theme <theme_name>
jenv theme dark
```

### help & version
```bash
jenv help [command]
jenv --version
```

## Q&A

### Why are administrator privileges needed?

Due to Windows system restrictions, creating system-level symbolic links requires:
`Running PowerShell as Administrator`

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
    - Creates a single symlink at `C:\java\JAVA_HOME` during installation
    - Switching Java versions only requires updating the symlink target
    - No need to modify system PATH repeatedly
    - Changes persist across system reboots and apply to all console windows

2. **Implementation Details**
    - During initialization:
        - Creates `JAVA_HOME` directory at `C:\java\JAVA_HOME`
        - Adds `JAVA_HOME\bin` to system PATH (one-time setup)
        - Creates initial symlink to default JDK
    - When switching versions:
        - Simply updates symlink target to desired JDK
        - No PATH modifications needed
        - Changes take effect immediately in all console windows

3. **Administrative Privileges**
    - Administrator privileges are only required when creating/modifying symbolic links
    - UAC prompts are handled automatically with minimal privilege scope
    - Follows the principle of least privilege, requesting only necessary permissions
    - Permission requests only occur during initialization (jenv init) and version switching (jenv use)

This approach is more efficient than constantly modifying system PATH variables, providing a cleaner and more reliable
solution for Java version management on Windows.

## Acknowledgments

- [cobra](https://github.com/spf13/cobra) - A powerful CLI framework for Go
- [jreleaser](https://jreleaser.org/) - A release automation tool
- [nvm-windows](https://github.com/coreybutler/nvm-windows) - Inspired our symlink-based approach
- [Jenv-for-Windows](https://github.com/FelixSelter/JEnv-for-Windows) - A predecessor project for Java version
  management on Windows

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

