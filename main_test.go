package main // same package as main.go, so can call unexported functions like loadTargets to test.

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadTargets(t *testing.T) { // Go test runs everry func named TestXxx in a _test.g file
	// t is how we report errors and log in tests
	tests := []struct { // table-driven tests
		name        string // name of the test case
		fileContent string // content of the file to create
		wantErr     bool   // whether we expect an error from loadTargets
	}{
		{"good config", `{"timeoutSeconds": 30, "source": [{"host": "example.com", "port": 443}]}`, false},
		{"no targets", `{"timeoutSeconds": 30, "source": []}`, true},
		{"invalid JSON", `{"timeoutSeconds": 30, "source": [}`, true},
		{"file does not exist", "", true},
	}

	name := t.TempDir() // create a temporary directory for the test files

	for _, tt := range tests { // _ throws away the index. tt is a copy of the struct for each test case
		t.Run(tt.name, func(t *testing.T) { // subtest for each test case
			if tt.fileContent != "" { // missing file content means we don't create the file, simulating a non-existent file
				// declare and check in one line
				if err := os.WriteFile(filepath.Join(name, tt.name), []byte(tt.fileContent), 0o644); err != nil {
					t.Fatalf("Failed to write test file: %v", err) // fatal AND stop this case, without the file, there's no point in continuing
				}
			}
			_, err := loadTargets(filepath.Join(name, tt.name)) // _ throws away the Target. Only care about the error for this test
			if (err != nil) != tt.wantErr {                     // "Did it fail?" vs. "Should it have failed?" If they don't match, we have a problem.
				// if they disagree, report the error. t.Errorf reports an error but continues the test, unlike t.Fatalf which stops the test.
				t.Errorf("loadTargets(%q) error = %v, wantErr %v", filepath.Join(name, tt.name), err, tt.wantErr)
			}
		})
	}
}
