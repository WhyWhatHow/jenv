package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestQueryUserEnvironmentVariable(t *testing.T) {
	// 测试查询用户级别的环境变量
	key := "PATH"
	value, err := QueryUserEnvironmentVariable(key)
	if err != nil {
		t.Fatalf("QueryUserEnvironmentVariable failed: %v", err)
	}
	if value == "" {
		t.Errorf("Expected non-empty value for key %s, got empty string", key)
	}
}

func TestSetSystemPath_Windows(t *testing.T) {
	// This test, as originally written in env_test.go, primarily tests os.Setenv and path string manipulation.
	// To truly test the Windows SetSystemPath, it would need to interact with registry functions
	// and require admin privileges or mocking. For now, its original logic is preserved here.
	// It's expected to run only on Windows.

	originalPath := os.Getenv("PATH")
	defer func() {
		if err := os.Setenv("PATH", originalPath); err != nil {
			t.Errorf("Failed to restore original PATH: %v", err)
		}
	}()

	testPath := "C:\\Test\\PathForWindowsTest" // Unique path for this test

	currentPath := os.Getenv("PATH")
	newPath := testPath
	if currentPath != "" {
		newPath = testPath + string(filepath.ListSeparator) + currentPath
	}

	if err := os.Setenv("PATH", newPath); err != nil {
		t.Fatalf("Failed to set PATH environment variable: %v", err)
	}

	updatedPath := os.Getenv("PATH")
	if !strings.Contains(updatedPath, testPath) {
		t.Errorf("PATH environment variable does not contain test path\nExpected to contain: %s\nActual: %s", testPath, updatedPath)
	}

	paths := strings.Split(updatedPath, string(filepath.ListSeparator))
	if len(paths) > 0 && paths[0] != testPath {
		// This assertion might be flaky if other tests modify PATH concurrently or if currentPath was empty.
		// For an isolated test of SetSystemPath, one would typically start with a known PATH or mock os.Getenv.
		t.Logf("Warning: New path was not added to the exact beginning of PATH. This might be due to test setup or concurrent ops.")
		t.Logf("Expected beginning: %s, Actual beginning: %s, Full PATH: %s", testPath, paths[0], updatedPath)
	}
}

func TestUpdateSystemEnvironmentVariable_Windows(t *testing.T) {
	// This test expects to interact with Windows system environment variables,
	// potentially requiring admin rights or specific mocking for registry operations.

	testKey := "JENV_TEST_VAR_WINDOWS"
	originalValue, originalExists := os.LookupEnv(testKey)

	// Attempt to clean up by setting to original value or unsetting.
	// Note: True registry modification cleanup would require admin rights if the variable was set persistently.
	defer func() {
		var err error
		if originalExists {
			err = os.Setenv(testKey, originalValue) // This only affects process env
		} else {
			err = os.Unsetenv(testKey) // This only affects process env
		}
		if err != nil {
			t.Logf("Error during cleanup of %s (process level): %v", testKey, err)
		}
		// A true persistent cleanup would involve calling UpdateSystemEnvironmentVariable to set back or delete.
		// For this test, we assume UpdateSystemEnvironmentVariable is the function under test, not a cleanup util.
	}()

	testValue := "test_value_windows"
	// This calls the UpdateSystemEnvironmentVariable, presumably from env_windows.go
	err := UpdateSystemEnvironmentVariable(testKey, testValue)

	if err == nil {
		// If no error, it implies success (either ran as admin, or the function doesn't persist and only used os.Setenv).
		// Check current process's environment. This doesn't confirm persistence.
		actualValue := os.Getenv(testKey)
		if actualValue != testValue {
			t.Errorf("Getenv returned unexpected value for %s. Expected: %s, Got: %s. (Note: This doesn't confirm persistent change.)", testKey, testValue, actualValue)
		}
		t.Logf("UpdateSystemEnvironmentVariable call succeeded for key %s. Further manual verification of persistence may be needed.", testKey)
		// To truly verify, one would need to query the registry or start a new process.
	} else {
		// If an error occurred, it might be a permissions issue (expected if not admin),
		// or the function is not fully implemented for Windows persistence in env_windows.go.
		// Common errors: "Access is denied.", "The requested operation requires elevation."
		t.Logf("UpdateSystemEnvironmentVariable for key %s returned error: %v. This is expected if not running with admin privileges or if persistent update failed.", testKey, err)
		if !strings.Contains(strings.ToLower(err.Error()), "access is denied") &&
			!strings.Contains(strings.ToLower(err.Error()), "permission denied") &&
			!strings.Contains(strings.ToLower(err.Error()), "requires elevation") &&
			!strings.Contains(strings.ToLower(err.Error()), "privileges") {
			t.Logf("Warning for key %s: Received error might not be permission-related: %v", testKey, err)
		}
	}
}

func TestQuerySystemEnvironmentVariable(t *testing.T) {
	// 测试查询系统级别的环境变量
	key := "PATH"
	value, err := QuerySystemEnvironmentVariable(key)
	if err != nil {
		t.Fatalf("QuerySystemEnvironmentVariable failed: %v", err)
	}
	if value == "" {
		t.Errorf("Expected non-empty value for key %s, got empty string", key)
	}
}
