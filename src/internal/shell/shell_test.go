//go:build !windows

package shell

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDetectUserShells(t *testing.T) {
	// Create a temporary home directory
	tempDir, err := os.MkdirTemp("", "shell_test_home")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Save original HOME
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	
	// Set temporary HOME
	os.Setenv("HOME", tempDir)
	
	// Create shell config files
	bashrc := filepath.Join(tempDir, ".bashrc")
	zshrc := filepath.Join(tempDir, ".zshrc")
	fishConfigDir := filepath.Join(tempDir, ".config", "fish")
	fishConfig := filepath.Join(fishConfigDir, "config.fish")
	
	// Create .bashrc
	if err := os.WriteFile(bashrc, []byte("# bashrc"), 0644); err != nil {
		t.Fatalf("Failed to create .bashrc: %v", err)
	}
	
	// Create .zshrc
	if err := os.WriteFile(zshrc, []byte("# zshrc"), 0644); err != nil {
		t.Fatalf("Failed to create .zshrc: %v", err)
	}
	
	// Create fish config
	if err := os.MkdirAll(fishConfigDir, 0755); err != nil {
		t.Fatalf("Failed to create fish config dir: %v", err)
	}
	if err := os.WriteFile(fishConfig, []byte("# fish config"), 0644); err != nil {
		t.Fatalf("Failed to create fish config: %v", err)
	}
	
	shells, err := DetectUserShells()
	if err != nil {
		t.Fatalf("DetectUserShells failed: %v", err)
	}
	
	// Should detect bash, zsh, fish, and profile
	expectedShells := []ShellType{Bash, Zsh, Fish, Profile}
	if len(shells) != len(expectedShells) {
		t.Errorf("Expected %d shells, got %d: %v", len(expectedShells), len(shells), shells)
	}
	
	for _, expected := range expectedShells {
		found := false
		for _, shell := range shells {
			if shell == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected shell %s not found in detected shells: %v", expected, shells)
		}
	}
}

func TestSetEnvironmentVariableForShell(t *testing.T) {
	// Create a temporary home directory
	tempDir, err := os.MkdirTemp("", "shell_test_setenv")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Save original HOME
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	
	// Set temporary HOME
	os.Setenv("HOME", tempDir)
	
	key := "TEST_JAVA_HOME"
	value := "/opt/java/jdk11"
	
	testCases := []struct {
		shell        ShellType
		configFile   string
		expectedLine string
	}{
		{
			shell:        Bash,
			configFile:   ".bashrc",
			expectedLine: "export " + key + "=\"" + value + "\"",
		},
		{
			shell:        Zsh,
			configFile:   ".zshrc",
			expectedLine: "export " + key + "=\"" + value + "\"",
		},
		{
			shell:        Fish,
			configFile:   ".config/fish/config.fish",
			expectedLine: "set -gx " + key + " \"" + value + "\"",
		},
		{
			shell:        Profile,
			configFile:   ".profile",
			expectedLine: "export " + key + "=\"" + value + "\"",
		},
	}
	
	for _, tc := range testCases {
		t.Run(string(tc.shell), func(t *testing.T) {
			err := SetEnvironmentVariableForShell(tc.shell, key, value)
			if err != nil {
				t.Errorf("SetEnvironmentVariableForShell failed for %s: %v", tc.shell, err)
				return
			}
			
			// Verify the config file was created and contains the expected line
			configPath := filepath.Join(tempDir, tc.configFile)
			content, err := os.ReadFile(configPath)
			if err != nil {
				t.Fatalf("Failed to read config file %s: %v", configPath, err)
			}
			
			if !strings.Contains(string(content), tc.expectedLine) {
				t.Errorf("Expected line '%s' not found in %s content: %s", tc.expectedLine, tc.configFile, string(content))
			}
		})
	}
}

func TestRemoveEnvironmentVariableFromShell(t *testing.T) {
	// Create a temporary home directory
	tempDir, err := os.MkdirTemp("", "shell_test_remove")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Save original HOME
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	
	// Set temporary HOME
	os.Setenv("HOME", tempDir)
	
	key := "TEST_JAVA_HOME"
	value := "/opt/java/jdk11"
	
	// First, set the environment variable for bash
	err = SetEnvironmentVariableForShell(Bash, key, value)
	if err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	
	// Verify it was set
	bashrcPath := filepath.Join(tempDir, ".bashrc")
	content, err := os.ReadFile(bashrcPath)
	if err != nil {
		t.Fatalf("Failed to read .bashrc: %v", err)
	}
	
	expectedLine := "export " + key + "=\"" + value + "\""
	if !strings.Contains(string(content), expectedLine) {
		t.Fatalf("Environment variable was not set correctly")
	}
	
	// Now remove it
	err = RemoveEnvironmentVariableFromShell(Bash, key)
	if err != nil {
		t.Errorf("RemoveEnvironmentVariableFromShell failed: %v", err)
		return
	}
	
	// Verify it was removed
	content, err = os.ReadFile(bashrcPath)
	if err != nil {
		t.Fatalf("Failed to read .bashrc after removal: %v", err)
	}
	
	if strings.Contains(string(content), expectedLine) {
		t.Errorf("Environment variable line should have been removed from .bashrc: %s", string(content))
	}
}

func TestGetCurrentShell(t *testing.T) {
	// Save original SHELL
	originalShell := os.Getenv("SHELL")
	defer os.Setenv("SHELL", originalShell)
	
	testCases := []struct {
		shellPath    string
		expectedType ShellType
	}{
		{"/bin/bash", Bash},
		{"/usr/bin/bash", Bash},
		{"/bin/zsh", Zsh},
		{"/usr/bin/zsh", Zsh},
		{"/usr/bin/fish", Fish},
		{"/bin/fish", Fish},
		{"/bin/sh", Profile},
		{"", Profile},
	}
	
	for _, tc := range testCases {
		t.Run(tc.shellPath, func(t *testing.T) {
			os.Setenv("SHELL", tc.shellPath)
			result := GetCurrentShell()
			if result != tc.expectedType {
				t.Errorf("Expected %s, got %s for shell path %s", tc.expectedType, result, tc.shellPath)
			}
		})
	}
}

func TestUpdateShellConfigFile(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "shell_test_config")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()
	
	configs := GetShellConfigs()
	bashConfig := configs[Bash]
	fishConfig := configs[Fish]
	
	key := "TEST_VAR"
	value := "test_value"
	
	// Test bash-style config
	err = updateShellConfigFile(tempFile.Name(), bashConfig, key, value)
	if err != nil {
		t.Errorf("updateShellConfigFile failed for bash: %v", err)
	}
	
	// Read and verify content
	content, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}
	
	expectedLine := "export " + key + "=\"" + value + "\""
	if !strings.Contains(string(content), expectedLine) {
		t.Errorf("Expected line '%s' not found in bash config: %s", expectedLine, string(content))
	}
	
	// Test fish-style config
	tempFile2, err := os.CreateTemp("", "shell_test_fish_config")
	if err != nil {
		t.Fatalf("Failed to create temp file for fish: %v", err)
	}
	defer os.Remove(tempFile2.Name())
	tempFile2.Close()
	
	err = updateShellConfigFile(tempFile2.Name(), fishConfig, key, value)
	if err != nil {
		t.Errorf("updateShellConfigFile failed for fish: %v", err)
	}
	
	// Read and verify fish content
	content, err = os.ReadFile(tempFile2.Name())
	if err != nil {
		t.Fatalf("Failed to read temp fish file: %v", err)
	}
	
	expectedFishLine := "set -gx " + key + " \"" + value + "\""
	if !strings.Contains(string(content), expectedFishLine) {
		t.Errorf("Expected line '%s' not found in fish config: %s", expectedFishLine, string(content))
	}
}
