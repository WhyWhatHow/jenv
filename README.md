# Jenv: Java Environment Manager



Jenv-Go is a command-line tool for managing multiple Java versions on your system. It allows you to easily switch between different Java versions, add new Java installations, and manage your Java environment.

## Features

- Add and manage multiple JDK installations
- Switch between Java versions with a single command
- Automatically update environment variables (JAVA_HOME, PATH)
- Cross-platform support (Windows, Linux)
  - [x] `Windows`
  - [ ] `Linux`
  - [ ] `MacOS`

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
├── templates/          # Template files
└── .github/            # GitHub configurations
    └── workflows/      # CI/CD workflows
```

## Installation

### From Release

Download the latest release from the [Releases page](https://github.com/WhyWhatHow/jenv/releases).

### Build from Source

#### Prerequisites

- Go 1.16 or higher
- Git

#### Build Steps

1. Clone the repository:

```bash
git clone https://github.com/WhyWhatHow/jenv.git
cd jenv
```

2. Build the project:

```bash

cd src 

go mod download
# For Windows (PowerShell)
go build -ldflags "-X github.com/whywhathow/jenv/cmd.Version=1.0.0" -o jenv.exe 

# For Linux/macOS
go build -ldflags "-X github.com/whywhathow/jenv/cmd.Version=1.0.0" -o jenv 

# For development build (with debug information)
go build -o jenv.exe 
```

## Usage

### Basic Commands

```bash
# Initialize or reinitialize jenv configuration
jenv init

# Add a new JDK with an alias name
jenv add <alias> <jdk_path>
jenv add jdk8 "C:\Program Files\Java\jdk1.8.0_291"

# List all installed JDKs
jenv list

# Switch to a specific JDK version
jenv use <alias>
jenv use jdk8

# Remove a JDK from jenv
jenv remove <alias>
jenv remove jdk8

# Update environment variables (JAVA_HOME, PATH)
jenv update
```

### Additional Features

```bash
# Scan system for installed JDKs
jenv scan c:\

# Show current JDK in use
jenv current

# Add jenv to system PATH
jenv add-to-path

# Change UI theme (light/dark)
jenv theme <theme_name>
jenv theme dark
```

For detailed information about each command and its options, use:
```bash
jenv help [command]
```

## Acknowledgments

Special thanks to [Trae](https://trae.ai) for providing the fantastic agentic IDE and AI Flow paradigm that greatly
enhanced our development experience.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
