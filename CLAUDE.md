# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Building
```bash
cd src
go build -ldflags "-X github.com/whywhathow/jenv/cmd.Version=1.0.0" -o jenv
```

For development builds with debug information:
```bash
cd src
go build -o jenv
```

### Testing
```bash
cd src
go test ./...
```

Run specific test:
```bash
cd src
go test ./internal/java -v
```

### Dependencies
```bash
cd src
go mod download
go mod tidy
```

### Release Build (using JReleaser)
```bash
# For tagged releases - handled by GitHub Actions
# Local assembly for testing:
jreleaser assemble
```

## Architecture Overview

### Core Structure
- **Entry Point**: `src/jenv.go` - Main entry point that calls `cmd.Execute()`
- **Commands**: `src/cmd/` - All CLI commands using Cobra framework
- **Internal Packages**: `src/internal/` - Core functionality modules

### Key Components

#### Configuration Management (`internal/config/`)
- Singleton pattern with thread-safe access using sync.RWMutex
- JSON-based configuration stored in user's home directory (`~/.jenv/config.json`)
- Cross-platform symlink management for Windows/Linux/macOS
- JDK registry using map-based storage for O(1) lookups

#### Java SDK Management (`internal/java/`)
- Concurrent JDK scanning using Dispatcher-Worker pattern with goroutines
- Platform-specific path validation (javac.exe on Windows, javac on Unix)
- Ultra-fast scanning optimized from 3s to 300ms through aggressive pre-filtering

#### Environment Management (`internal/env/`)
- Platform-specific environment variable handling
- Windows: Registry-based environment management
- Linux/Unix: Shell configuration file updates (.bashrc, .zshrc, config.fish)
- Automatic privilege escalation detection and handling

#### Cross-Platform System Integration (`internal/sys/`)
- Administrator privilege detection for Windows (UAC) and Unix (root)
- Platform-specific permission handling with graceful degradation

#### UI and Styling (`internal/style/`)
- Theme system with light/dark mode support
- Consistent styling using Charm's lipgloss library
- Terminal-friendly output formatting

### Symlink-Based Architecture
The application uses symbolic links for efficient Java version switching:
- **Windows**: `C:\java\JAVA_HOME` points to active JDK
- **Linux System**: `/opt/jenv/java_home` (root) or `~/.jenv/java_home` (user)
- **macOS**: `/opt/jenv/java_home` (root) or `~/.jenv/java_home` (user)

### Platform Support Matrix
- **Windows**: Complete support with administrator privileges required
- **Linux**: Complete support with optional root privileges
- **macOS**: Infrastructure ready, implementation in progress

### Concurrency Patterns
- Singleton configuration with double-checked locking
- Read-write locks for configuration access
- Worker pool pattern for JDK scanning performance
- Thread-safe JDK registry operations

### Error Handling
- Custom error types: `ErrJDKExists`, `ErrJDKNotFound`, `ErrInvalidPath`, `ErrNotInitialized`
- Graceful degradation when admin privileges unavailable on Linux
- Platform-specific error messages and remediation suggestions

## Development Notes

### Go Module Structure
- Module: `github.com/whywhathow/jenv`
- Go version: 1.23.4
- Main dependencies: Cobra CLI, Charm Lipgloss, testify

### Version Injection
Version is injected at build time using ldflags:
```bash
-ldflags "-X github.com/whywhathow/jenv/cmd.Version=${VERSION}"
```

### Testing Pattern
Tests use standard Go testing with testify assertions. Test files follow `*_test.go` naming convention.

### Release Process
- Uses JReleaser for cross-platform builds and GitHub releases
- Automated through GitHub Actions on tag push
- Supports multiple architectures: amd64, arm64 for Windows, Linux, macOS