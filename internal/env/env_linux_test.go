//go:build linux
// +build linux

package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUpdateEnvironmentVariable_UserBashrc(t *testing.T) {
	// Setup: Create a temporary home directory
	testHome, err := os.MkdirTemp("", "test_home_bashrc_")
	if err != nil {
		t.Fatalf("Failed to create temp home dir: %v", err)
	}
	defer os.RemoveAll(testHome)

	// Save original HOME and SHELL, and restore them after the test
	originalHome := os.Getenv("HOME")
	originalShell := os.Getenv("SHELL")
	defer os.Setenv("HOME", originalHome)
	defer os.Setenv("SHELL", originalShell)

	// Set HOME and SHELL for the test's context
	if err := os.Setenv("HOME", testHome); err != nil {
		t.Fatalf("Failed to set HOME env var for test: %v", err)
	}
	if err := os.Setenv("SHELL", "/bin/bash"); err != nil {
		t.Fatalf("Failed to set SHELL env var for test: %v", err)
	}

	// Create a dummy .bashrc file
	bashrcPath := filepath.Join(testHome, ".bashrc")
	initialBashrcContent := "# Original bashrc content\nexport OLD_VAR=\"old_value\"\n# Another line"
	err = os.WriteFile(bashrcPath, []byte(initialBashrcContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write initial .bashrc: %v", err)
	}

	// Execution: Call UpdateEnvironmentVariable
	// Assuming UpdateEnvironmentVariable is in the same 'env' package
	err = UpdateEnvironmentVariable("NEW_VAR", "new_value")
	if err != nil {
		// The function prints warnings to stderr if /etc/profile.d is not writable,
		// but it should not return an error if user-level update succeeds.
		// Check if the error is unexpected. For now, we log it.
		// If user-level update failed, that's a test failure.
		t.Logf("UpdateEnvironmentVariable for NEW_VAR returned error (expected nil if user update is successful): %v", err)
	}

	err = UpdateEnvironmentVariable("OLD_VAR", "updated_value")
	if err != nil {
		t.Logf("UpdateEnvironmentVariable for OLD_VAR returned error (expected nil if user update is successful): %v", err)
	}

	// Assertions: Read the content of .bashrc
	updatedBashrcContentBytes, err := os.ReadFile(bashrcPath)
	if err != nil {
		t.Fatalf("Failed to read updated .bashrc: %v", err)
	}
	updatedBashrcContent := string(updatedBashrcContentBytes)
	t.Logf("Updated .bashrc content:\n%s", updatedBashrcContent)

	// Verify NEW_VAR
	expectedNewVarLine := "export NEW_VAR=\"new_value\""
	if !strings.Contains(updatedBashrcContent, expectedNewVarLine) {
		t.Errorf(".bashrc missing expected line: %s", expectedNewVarLine)
	}

	// Verify OLD_VAR updated
	expectedOldVarUpdatedLine := "export OLD_VAR=\"updated_value\""
	if !strings.Contains(updatedBashrcContent, expectedOldVarUpdatedLine) {
		t.Errorf(".bashrc missing expected updated line: %s", expectedOldVarUpdatedLine)
	}
	// Verify OLD_VAR old value is not there (or commented out, depending on updateScriptFile logic)
	// updateScriptFile replaces the line, so the old value shouldn't be there as an active export.
	oldVarOldValueLine := "export OLD_VAR=\"old_value\""
	if strings.Contains(updatedBashrcContent, oldVarOldValueLine) {
		t.Errorf(".bashrc still contains old value line for OLD_VAR: %s", oldVarOldValueLine)
	}


	// Verify original content (other lines)
	if !strings.Contains(updatedBashrcContent, "# Original bashrc content") {
		t.Errorf(".bashrc missing original content: '# Original bashrc content'")
	}
	if !strings.Contains(updatedBashrcContent, "# Another line") {
		t.Errorf(".bashrc missing original content: '# Another line'")
	}

	// Verify current process environment (os.Setenv should have been called)
	if val := os.Getenv("NEW_VAR"); val != "new_value" {
		t.Errorf("Getenv(NEW_VAR) got %s, want new_value", val)
	}
	if val := os.Getenv("OLD_VAR"); val != "updated_value" {
		t.Errorf("Getenv(OLD_VAR) got %s, want updated_value", val)
	}
}

func TestUpdateEnvironmentVariable_UserProfileFallback(t *testing.T) {
	// Setup: Create a temporary home directory
	testHome, err := os.MkdirTemp("", "test_home_profile_")
	if err != nil {
		t.Fatalf("Failed to create temp home dir: %v", err)
	}
	defer os.RemoveAll(testHome)

	originalHome := os.Getenv("HOME")
	originalShell := os.Getenv("SHELL")
	defer os.Setenv("HOME", originalHome)
	defer os.Setenv("SHELL", originalShell)

	if err := os.Setenv("HOME", testHome); err != nil {
		t.Fatalf("Failed to set HOME env var for test: %v", err)
	}
	// Set SHELL to something that should cause a fallback to .profile
	if err := os.Setenv("SHELL", "/bin/sh"); err != nil {
		t.Fatalf("Failed to set SHELL env var for test: %v", err)
	}

	// Create a dummy .profile file
	profilePath := filepath.Join(testHome, ".profile")
	initialProfileContent := "# Original .profile content\nexport EXISTING_PROFILE_VAR=\"initial_value\""
	err = os.WriteFile(profilePath, []byte(initialProfileContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write initial .profile: %v", err)
	}

	// Execution
	err = UpdateEnvironmentVariable("PROFILE_VAR", "profile_value")
	if err != nil {
		t.Logf("UpdateEnvironmentVariable for PROFILE_VAR returned error: %v", err)
	}
	err = UpdateEnvironmentVariable("EXISTING_PROFILE_VAR", "updated_profile_value")
	if err != nil {
		t.Logf("UpdateEnvironmentVariable for EXISTING_PROFILE_VAR returned error: %v", err)
	}

	// Assertions
	updatedProfileContentBytes, err := os.ReadFile(profilePath)
	if err != nil {
		t.Fatalf("Failed to read updated .profile: %v", err)
	}
	updatedProfileContent := string(updatedProfileContentBytes)
	t.Logf("Updated .profile content:\n%s", updatedProfileContent)

	expectedProfileVarLine := "export PROFILE_VAR=\"profile_value\""
	if !strings.Contains(updatedProfileContent, expectedProfileVarLine) {
		t.Errorf(".profile missing expected line: %s", expectedProfileVarLine)
	}

	expectedExistingProfileVarUpdatedLine := "export EXISTING_PROFILE_VAR=\"updated_profile_value\""
	if !strings.Contains(updatedProfileContent, expectedExistingProfileVarUpdatedLine) {
		t.Errorf(".profile missing expected updated line for EXISTING_PROFILE_VAR: %s", expectedExistingProfileVarUpdatedLine)
	}

	oldExistingProfileVarLine := "export EXISTING_PROFILE_VAR=\"initial_value\""
	if strings.Contains(updatedProfileContent, oldExistingProfileVarLine) {
		t.Errorf(".profile still contains old value line for EXISTING_PROFILE_VAR: %s", oldExistingProfileVarLine)
	}

	if !strings.Contains(updatedProfileContent, "# Original .profile content") {
		t.Errorf(".profile missing original content: '# Original .profile content'")
	}

	// Verify current process environment
	if val := os.Getenv("PROFILE_VAR"); val != "profile_value" {
		t.Errorf("Getenv(PROFILE_VAR) got %s, want profile_value", val)
	}
	if val := os.Getenv("EXISTING_PROFILE_VAR"); val != "updated_profile_value" {
		t.Errorf("Getenv(EXISTING_PROFILE_VAR) got %s, want updated_profile_value", val)
	}
}

func TestUpdateEnvironmentVariable_SystemProfileD_NoAccessFallback(t *testing.T) {
	// This test verifies the fallback behavior when /etc/profile.d is not writable.
	// It's a partial test for the system-level scenario. Testing the actual
	// write to /etc/profile.d/jenv.sh is complex due to hardcoded paths and
	// permissions, and is noted as a limitation for this test suite.

	testHome, err := os.MkdirTemp("", "test_home_system_fallback_")
	if err != nil {
		t.Fatalf("Failed to create temp home dir: %v", err)
	}
	defer os.RemoveAll(testHome)

	originalHome := os.Getenv("HOME")
	originalShell := os.Getenv("SHELL")
	defer os.Setenv("HOME", originalHome)
	defer os.Setenv("SHELL", originalShell)

	if err := os.Setenv("HOME", testHome); err != nil {
		t.Fatalf("Failed to set HOME env var for test: %v", err)
	}
	// Use .bashrc for fallback path in this specific test
	if err := os.Setenv("SHELL", "/bin/bash"); err != nil {
		t.Fatalf("Failed to set SHELL env var for test: %v", err)
	}

	bashrcPath := filepath.Join(testHome, ".bashrc")
	// No need to create .bashrc beforehand, updateScriptFile should create it.

	// Execution: Call UpdateEnvironmentVariable
	// canWriteToDir("/etc/profile.d") should be false in test env.
	err = UpdateEnvironmentVariable("SYSTEM_TEST_VAR", "system_value")
	if err != nil {
		t.Logf("UpdateEnvironmentVariable for SYSTEM_TEST_VAR returned error: %v", err)
	}

	// Assertions: Verify it fell back to user's .bashrc
	updatedBashrcContentBytes, err := os.ReadFile(bashrcPath)
	if err != nil {
		t.Fatalf("Failed to read .bashrc (expected fallback): %v", err)
	}
	updatedBashrcContent := string(updatedBashrcContentBytes)
	t.Logf("Fallback .bashrc content for system test:\n%s", updatedBashrcContent)

	expectedLine := "export SYSTEM_TEST_VAR=\"system_value\""
	if !strings.Contains(updatedBashrcContent, expectedLine) {
		t.Errorf(".bashrc (fallback) missing expected line: %s", expectedLine)
	}

	// Verify current process environment
	if val := os.Getenv("SYSTEM_TEST_VAR"); val != "system_value" {
		t.Errorf("Getenv(SYSTEM_TEST_VAR) got %s, want system_value", val)
	}
}
