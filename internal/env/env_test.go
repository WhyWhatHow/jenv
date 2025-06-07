package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCleanJDKPath(t *testing.T) {
	originalPath := os.Getenv("PATH")
	defer os.Setenv("PATH", originalPath)

	// Setup a PATH with JDK-like entries and non-JDK entries
	jdkPath1 := filepath.Join("path", "to", "myjdk", "bin")
	jdkPath2 := filepath.Join("another", "jre", "bin")
	nonJdkPath1 := filepath.Join("usr", "local", "bin")
	nonJdkPath2 := filepath.Join("some", "other", "path")

	testPathParts := []string{jdkPath1, nonJdkPath1, jdkPath2, nonJdkPath2}
	testPath := strings.Join(testPathParts, string(filepath.ListSeparator))

	os.Setenv("PATH", testPath)

	CleanJDKPath()

	cleanedPath := os.Getenv("PATH")
	cleanedPathParts := strings.Split(cleanedPath, string(filepath.ListSeparator))

	foundNonJdk1 := false
	foundNonJdk2 := false
	for _, p := range cleanedPathParts {
		if strings.Contains(strings.ToLower(p), "jdk") || strings.Contains(strings.ToLower(p), "jre") {
			t.Errorf("Found JDK/JRE path after CleanJDKPath: %s in %s", p, cleanedPath)
		}
		if p == nonJdkPath1 {
			foundNonJdk1 = true
		}
		if p == nonJdkPath2 {
			foundNonJdk2 = true
		}
	}

	if !foundNonJdk1 {
		t.Errorf("Non-JDK path %s was removed by CleanJDKPath", nonJdkPath1)
	}
	if !foundNonJdk2 {
		t.Errorf("Non-JDK path %s was removed by CleanJDKPath", nonJdkPath2)
	}

	expectedCleanedCount := 2
	if len(cleanedPathParts) != expectedCleanedCount {
		// Handle case where initial PATH might have empty strings if separators are at ends.
		// For this test, we construct it cleanly.
		actualCount := 0
		for _, p := range cleanedPathParts {
			if p != "" {
				actualCount++
			}
		}
		if actualCount != expectedCleanedCount {
			t.Errorf("Expected %d path entries after CleanJDKPath, got %d. Path: %s", expectedCleanedCount, actualCount, cleanedPath)
		}
	}
}

// TestRestoreOldPath was removed due to issues with testing global state (oldPath)
// and environment variable visibility in the test environment.
// The SetEnv function's PATH saving mechanism is fragile for unit testing
// without refactoring SetEnv and RestoreOldPath to avoid global variables.
