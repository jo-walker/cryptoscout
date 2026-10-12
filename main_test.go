package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadTargets(t *testing.T) {
	tests := []struct {
		filename    string
		dirname     string
		fileContent string
		wantErr     bool
	}{
		{"good config", "targets.json", `{"timeoutSeconds": 30, "source": [{"host": "example.com", "port": 443}]}`, false},
		{"no targets", "targets.json", `{"timeoutSeconds": 30, "source": []}`, true},
		{"invalid JSON", "targets.json", `{"timeoutSeconds": 30, "source": [}`, true},
		{"file does not exist", "nonexistent.json", "", true},
	}

	dirname := t.TempDir()

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			if tt.fileContent != "" {
				if err := os.WriteFile(filepath.Join(dirname, tt.filename), []byte(tt.fileContent), 0o644); err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}
			}
			_, err := loadTargets(filepath.Join(dirname, tt.filename))
			if (err != nil) != tt.wantErr {
				t.Errorf("loadTargets(%q) error = %v, wantErr %v", filepath.Join(dirname, tt.filename), err, tt.wantErr)
			}
		})
	}
}
