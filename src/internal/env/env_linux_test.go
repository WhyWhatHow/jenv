//go:build linux
// +build linux

package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDoSetEnv(t *testing.T) {
	// Test setting environment variable
	key := "TEST_JENV_VAR"
	value := "test_value"
	
	err := doSetEnv(key, value)
	if err != nil {
		t.Errorf("doSetEnv failed: %v", err)
	}
	
	// Verify the environment variable is set in current process
	if os.Getenv(key) != value {
		t.Errorf("Environment variable not set correctly. Expected: %s, Got: %s", value, os.Getenv(key))
	}
	
	// Clean up
	os.Unsetenv(key)
}

func TestDetectUserShells(t *testing.T) {
	// Create a temporary home directory
	tempDir, err := os.MkdirTemp("", "jenv_test_home")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create shell config files
	bashrc := filepath.Join(tempDir, ".bashrc")
	zshrc := filepath.Join(tempDir, ".zshrc")
	
	// Create .bashrc
	if err := os.WriteFile(bashrc, []byte("# bashrc"), 0644); err != nil {
		t.Fatalf("Failed to create .bashrc: %v", err)
	}
	
	// Create .zshrc
	if err := os.WriteFile(zshrc, []byte("# zshrc"), 0644); err != nil {
		t.Fatalf("Failed to create .zshrc: %v", err)
	}
	
	shells := detectUserShells(tempDir)
	
	// Should detect bash and zsh
	expectedShells := []string{"bash", "zsh"}
	if len(shells) != len(expectedShells) {
		t.Errorf("Expected %d shells, got %d", len(expectedShells), len(shells))
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

func TestUpdateBashEnvironment(t *testing.T) {
	// Create a temporary home directory
	tempDir, err := os.MkdirTemp("", "jenv_test_bash")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	key := "TEST_JAVA_HOME"
	value := "/opt/java/jdk11"
	
	err = updateBashEnvironment(tempDir, key, value)
	if err != nil {
		t.Errorf("updateBashEnvironment failed: %v", err)
	}
	
	// Verify the .bashrc file was created and contains the export
	bashrcPath := filepath.Join(tempDir, ".bashrc")
	content, err := os.ReadFile(bashrcPath)
	if err != nil {
		t.Fatalf("Failed to read .bashrc: %v", err)
	}
	
	expectedLine := "export " + key + "=\"" + value + "\""
	if !strings.Contains(string(content), expectedLine) {
		t.Errorf("Expected line '%s' not found in .bashrc content: %s", expectedLine, string(content))
	}
}

func TestUpdateFishEnvironment(t *testing.T) {
	// Create a temporary home directory
	tempDir, err := os.MkdirTemp("", "jenv_test_fish")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	key := "TEST_JAVA_HOME"
	value := "/opt/java/jdk11"
	
	err = updateFishEnvironment(tempDir, key, value)
	if err != nil {
		t.Errorf("updateFishEnvironment failed: %v", err)
	}
	
	// Verify the fish config file was created and contains the set command
	configPath := filepath.Join(tempDir, ".config", "fish", "config.fish")
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read fish config: %v", err)
	}
	
	expectedLine := "set -gx " + key + " \"" + value + "\""
	if !strings.Contains(string(content), expectedLine) {
		t.Errorf("Expected line '%s' not found in fish config content: %s", expectedLine, string(content))
	}
}

func TestCleanPathLinux(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		binPath  string
		expected string
	}{
		{
			name:     "Remove from middle",
			path:     "/usr/bin:/opt/java/bin:/usr/local/bin",
			binPath:  "/opt/java/bin",
			expected: "/usr/bin:/usr/local/bin",
		},
		{
			name:     "Remove from beginning",
			path:     "/opt/java/bin:/usr/bin:/usr/local/bin",
			binPath:  "/opt/java/bin",
			expected: "/usr/bin:/usr/local/bin",
		},
		{
			name:     "Remove from end",
			path:     "/usr/bin:/usr/local/bin:/opt/java/bin",
			binPath:  "/opt/java/bin",
			expected: "/usr/bin:/usr/local/bin",
		},
		{
			name:     "Path not found",
			path:     "/usr/bin:/usr/local/bin",
			binPath:  "/opt/java/bin",
			expected: "/usr/bin:/usr/local/bin",
		},
		{
			name:     "Empty path",
			path:     "",
			binPath:  "/opt/java/bin",
			expected: "",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := cleanPathLinux(tc.path, tc.binPath)
			if result != tc.expected {
				t.Errorf("Expected: %s, Got: %s", tc.expected, result)
			}
		})
	}
}

func TestQuerySystemEnvironmentVariable(t *testing.T) {
	// Test with a known environment variable
	key := "PATH"
	value, err := QuerySystemEnvironmentVariable(key)
	if err != nil {
		t.Errorf("QuerySystemEnvironmentVariable failed: %v", err)
	}
	
	if value == "" {
		t.Errorf("Expected non-empty PATH, got empty string")
	}
	
	// Test with non-existent variable
	nonExistentKey := "JENV_TEST_NON_EXISTENT_VAR"
	value, err = QuerySystemEnvironmentVariable(nonExistentKey)
	if err != nil {
		// This is expected for non-existent variables
		t.Logf("Expected error for non-existent variable: %v", err)
	}
	
	if value != "" {
		t.Errorf("Expected empty value for non-existent variable, got: %s", value)
	}
}

func TestUpdateShellConfigFile(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "jenv_test_config")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()
	
	key := "TEST_VAR"
	value := "test_value"
	
	err = updateShellConfigFile(tempFile.Name(), key, value, "export")
	if err != nil {
		t.Errorf("updateShellConfigFile failed: %v", err)
	}
	
	// Read the file and verify content
	content, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}
	
	expectedLine := "export " + key + "=\"" + value + "\""
	if !strings.Contains(string(content), expectedLine) {
		t.Errorf("Expected line '%s' not found in file content: %s", expectedLine, string(content))
	}
	
	// Test updating existing variable
	newValue := "new_test_value"
	err = updateShellConfigFile(tempFile.Name(), key, newValue, "export")
	if err != nil {
		t.Errorf("updateShellConfigFile update failed: %v", err)
	}
	
	// Read again and verify the value was updated
	content, err = os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read temp file after update: %v", err)
	}
	
	newExpectedLine := "export " + key + "=\"" + newValue + "\""
	if !strings.Contains(string(content), newExpectedLine) {
		t.Errorf("Expected updated line '%s' not found in file content: %s", newExpectedLine, string(content))
	}
	
	// Ensure old value is not present
	oldExpectedLine := "export " + key + "=\"" + value + "\""
	if strings.Contains(string(content), oldExpectedLine) {
		t.Errorf("Old line '%s' should not be present in file content: %s", oldExpectedLine, string(content))
	}
}
