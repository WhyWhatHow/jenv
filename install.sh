#!/bin/bash

set -e

# Define installation directory
JENV_DIR="$HOME/.jenv"
REPO_URL="https://github.com/WhyWhatHow/jenv.git"

# Check for git
if ! command -v git &> /dev/null; then
    echo "Error: git is not installed. Please install git and try again." >&2
    exit 1
fi

# Clone or update the repository
if [ -d "$JENV_DIR" ]; then
    echo "Updating jenv in $JENV_DIR..."
    cd "$JENV_DIR"
    git pull
    cd - > /dev/null
else
    echo "Cloning jenv into $JENV_DIR..."
    git clone "$REPO_URL" "$JENV_DIR"
fi

# Detect shell and update profile
SHELL_NAME=$(basename "$SHELL")
PROFILE_FILE=""

if [ "$SHELL_NAME" = "bash" ]; then
    PROFILE_FILE="$HOME/.bashrc"
elif [ "$SHELL_NAME" = "zsh" ]; then
    PROFILE_FILE="$HOME/.zshrc"
else
    echo "Warning: Unsupported shell '$SHELL_NAME'. Please add the following line to your shell's configuration file:" >&2
    echo "export PATH=\"\$HOME/.jenv/bin:\$PATH\"" >&2
    exit 0
fi

# Add to PATH if not already there
if ! grep -q 'export PATH="$HOME/.jenv/bin:$PATH"' "$PROFILE_FILE"; then
    echo "Adding jenv to PATH in $PROFILE_FILE..."
    echo '' >> "$PROFILE_FILE"
    echo '# jenv: Java Environment Manager' >> "$PROFILE_FILE"
    echo 'export PATH="$HOME/.jenv/bin:$PATH"' >> "$PROFILE_FILE"
    echo "Successfully added jenv to PATH."
else
    echo "jenv is already in your PATH."
fi

echo ""
echo "Installation complete!"
echo "Please restart your shell or run the following command to use jenv:"
echo "source $PROFILE_FILE"
