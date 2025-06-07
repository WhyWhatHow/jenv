package java

import (
	// "github.com/whywhathow/jenv/internal/config" // No longer directly used by tests in this file
	"os"
	"path/filepath"
	"strings"
	"testing"
	// "time" // No longer directly used by tests in this file
)

// TestScanJDK (the exported function version) was moved to sdk_windows_test.go as it was Windows-centric.
// TestInit was removed due to issues with global state management and Windows-centric assertions.
// The tests below focus on the unexported scanForJDKs logic for Linux and cross-platform behavior.

// Helper function to create a dummy java executable for testing
func createDummyJavaExecutable(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("Failed to create parent dirs for dummy java at %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte("dummy java"), 0755); err != nil {
		t.Fatalf("Failed to create dummy java at %s: %v", path, err)
	}
	// On non-Windows, ensure it's executable - WriteFile with 0755 should handle this.
}

func TestScanForJDKs_LinuxSkipsAndDetection(t *testing.T) {
	baseDir, err := os.MkdirTemp("", "test_scan_")
	if err != nil {
		t.Fatalf("Failed to create temp base dir: %v", err)
	}
	defer os.RemoveAll(baseDir)

	// Structure definition
	jdkDirs := map[string]bool{
		"real_jdk":       true,
		"another_jdk":    true,
		"deep_jdk/level1/level2/level3/level4/level5/level6/actual_jdk": true,
	}
	nonJdkDirs := []string{
		".hidden_dir",        // Should be skipped if top-level hidden skip is effective (current logic only skips baseDir if hidden)
		"proc_sim",           // Will NOT be skipped by "/proc" rule, as rule is absolute. Will be scanned.
		"my_app/node_modules/somepackage", // node_modules itself should be skipped
		"no_bin_jdk",
		"too_deep_jdk/level1/level2/level3/level4/level5/level6/level7/actual_jdk", // should be skipped by depth
	}

	// Create valid JDKs
	for relPath := range jdkDirs {
		fullPath := filepath.Join(baseDir, relPath, "bin", "java")
		createDummyJavaExecutable(t, fullPath)
	}

	// Create non-JDK and special directories
	for _, relPath := range nonJdkDirs {
		// For node_modules, create the structure up to node_modules, then a dummy java inside.
		// The 'node_modules' dir itself should cause the skip.
		if strings.Contains(relPath, "node_modules") {
			nmPath := filepath.Join(baseDir, relPath) // e.g. baseDir/my_app/node_modules/somepackage
			javaInNmPath := filepath.Join(nmPath, "bin", "java")
			createDummyJavaExecutable(t, javaInNmPath)
		} else if strings.HasPrefix(relPath, ".") { // hidden dir
			hiddenDirPath := filepath.Join(baseDir, relPath)
			if err := os.MkdirAll(hiddenDirPath, 0755); err != nil {
				t.Fatalf("Failed to create dir %s: %v", hiddenDirPath, err)
			}
			// also place a java in it to see if it's found
			createDummyJavaExecutable(t, filepath.Join(hiddenDirPath, "bin", "java"))
		} else {
			fullPath := filepath.Join(baseDir, relPath)
			if err := os.MkdirAll(fullPath, 0755); err != nil {
				t.Fatalf("Failed to create dir %s: %v", fullPath, err)
			}
		}
	}
	// Specifically for proc_sim, create a file in it, not a JDK structure
	if err := os.WriteFile(filepath.Join(baseDir, "proc_sim", "somefile"), []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create file in proc_sim: %v", err)
	}


	// Execution
	foundJDKs := scanForJDKs(baseDir)

	// Assertions
	// Expected: real_jdk, another_jdk, .hidden_dir/bin/java (named .hidden_dir)
	// Not Expected: deep_jdk/.../actual_jdk (depth 8), too_deep_jdk/.../actual_jdk (depth 9)
	expectedFoundCount := 3
	if len(foundJDKs) != expectedFoundCount {
		t.Errorf("Expected to find %d JDKs, but found %d.", expectedFoundCount, len(foundJDKs))
		for _, jdk := range foundJDKs {
			t.Logf("Found JDK: Name='%s', Path='%s'", jdk.Name, jdk.Path)
		}
	}

	foundNames := make(map[string]string) // Map name to path for easier checking
	for _, jdk := range foundJDKs {
		foundNames[jdk.Name] = jdk.Path
	}

	// Check for expected JDKs
	expectedJdks := []string{"real_jdk", "another_jdk", ".hidden_dir"}
	for _, name := range expectedJdks {
		if _, ok := foundNames[name]; !ok {
			t.Errorf("Expected JDK '%s' to be found, but it wasn't.", name)
		}
	}

	// Check for JDKs that should have been skipped
	skippedJdkNames := []string{"actual_jdk", "proc_sim", "no_bin_jdk", "somepackage"}
	// "actual_jdk" from deep_jdk and too_deep_jdk
	// "somepackage" is the name of the dir inside node_modules containing a dummy java

	for _, name := range skippedJdkNames {
		if _, ok := foundNames[name]; ok {
			t.Errorf("JDK '%s' (path: %s) should have been skipped, but was found.", name, foundNames[name])
		}
	}
	// Specifically check that the 'actual_jdk' from the 'deep_jdk' path was not found.
    // The name "actual_jdk" could be found if it was from a shallower path.
    // The count check is the main guard for depth issues.
    // If we want to be more precise, we'd check the full path of found JDKs.
    // For now, the count + specific name checks cover the main scenarios.
}

func TestScanForJDKs_EmptyDir(t *testing.T) {
	baseDir, err := os.MkdirTemp("", "test_scan_empty_")
	if err != nil {
		t.Fatalf("Failed to create temp base dir: %v", err)
	}
	defer os.RemoveAll(baseDir)

	foundJDKs := scanForJDKs(baseDir)
	if len(foundJDKs) != 0 {
		t.Errorf("Expected to find 0 JDKs in an empty directory, but found %d. Found: %v", len(foundJDKs), foundJDKs)
	}
}

func TestScanForJDKs_DirNotFound(t *testing.T) {
	// Using a path that is highly unlikely to exist
	nonExistentPath := filepath.Join(os.TempDir(), "jenv_test_non_existent_dir_abc123xyz")

	// Ensure it really doesn't exist (in case of very rare collision)
	os.RemoveAll(nonExistentPath) // Not strictly necessary if truly random, but good practice

	foundJDKs := scanForJDKs(nonExistentPath) // scanForJDKs prints an error but should return empty slice
	if len(foundJDKs) != 0 {
		t.Errorf("Expected to find 0 JDKs for a non-existent directory, but found %d. Found: %v", len(foundJDKs), foundJDKs)
	}
}
