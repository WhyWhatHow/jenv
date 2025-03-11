# Jenv-Go: Java Environment Manager

[![CI/CD Pipeline](https://github.com/WhyWhatHow/jenv-go-v2/actions/workflows/ci.yml/badge.svg)](https://github.com/WhyWhatHow/jenv-go-v2/actions/workflows/ci.yml)

Jenv-Go is a command-line tool for managing multiple Java versions on your system. It allows you to easily switch between different Java versions, add new Java installations, and manage your Java environment.

## Features

- Add and manage multiple JDK installations
- Switch between Java versions with a single command
- Automatically update environment variables (JAVA_HOME, PATH)
- Cross-platform support (Windows, Linux)

## Installation

Download the latest release from the [Releases page](https://github.com/WhyWhatHow/jenv-go-v2/releases).

## Usage

```bash
# Add a new JDK
jenv add jdk8 "C:\Program Files\Java\jdk1.8.0_291"

# List all installed JDKs
jenv list

# Switch to a specific JDK
jenv use jdk8

# Remove a JDK
jenv remove jdk8

# Update environment variables
jenv update

# Initialize or reinitialize jenv
jenv init
```

## Development Guide

### Project Structure

```
.
├── cmd/            # Command implementations
├── internal/       # Internal packages
│   ├── config/     # Configuration management
│   ├── constants/  # Constants definitions
│   ├── env/        # Environment variable handling
│   ├── java/       # Java SDK management
│   └── sys/        # System utilities
├── .github/        # GitHub configuration
│   └── workflows/  # CI/CD workflows
└── jenv.go         # Main entry point
```

### CI/CD Pipeline

This project uses GitHub Actions for continuous integration and delivery. The workflow is defined in `.github/workflows/ci.yml`.

#### Workflow Features

- **Automated Testing**: All tests are run on every push and pull request
- **Linting**: Code quality checks with golangci-lint
- **Cross-platform Builds**: Builds for Windows and Linux
- **Automated Releases**: Creates GitHub releases when tags are pushed

### How to Use the CI/CD Pipeline

#### For Development

1. Fork and clone the repository
2. Create a new branch for your feature or bugfix
3. Make your changes and commit them
4. Push your branch and create a pull request
5. The CI pipeline will automatically run tests and linting

#### For Releases

1. Update the version in `cmd/root.go` (in the `Version` field)
2. Commit your changes and push to the main branch
3. Create and push a new tag with the version number:

```bash
git tag v1.0.0
git push origin v1.0.0
```

4. The CI/CD pipeline will automatically:
   - Run tests and linting
   - Build binaries for Windows and Linux
   - Create a GitHub release with the binaries attached
   - Generate release notes

### Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

Apache License 2.0

## Author

WhyWhatHow (https://github.com/WhyWhatHow)