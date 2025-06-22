package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMain_Integration(t *testing.T) {
	tests := []struct {
		name        string
		yamlContent string
		args        []string
		expectError bool
	}{
		{
			name: "basic YAML with default json tags",
			yamlContent: `
name: "john"
age: 30
active: true
`,
			args:        []string{"test.yaml"},
			expectError: false,
		},
		{
			name: "YAML with custom tag prefix",
			yamlContent: `
user:
  name: "jane"
  email: "jane@example.com"
`,
			args:        []string{"-tag-prefix", "yaml", "test.yaml"},
			expectError: false,
		},
		{
			name: "single key YAML",
			yamlContent: `
key: "value"
`,
			args:        []string{"test.yaml"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpFile, err := os.CreateTemp("", "test*.yaml")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			// Write YAML content
			if _, err := tmpFile.WriteString(tt.yamlContent); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			tmpFile.Close()

			// Replace "test.yaml" in args with actual temp file name
			args := make([]string, len(tt.args))
			copy(args, tt.args)
			for i, arg := range args {
				if arg == "test.yaml" {
					args[i] = tmpFile.Name()
				}
			}

			// Build the command to run the main program
			cmd := exec.Command("go", append([]string{"run", "main.go"}, args...)...)
			cmd.Dir = "."

			// Run the command and capture output
			output, err := cmd.CombinedOutput()

			if tt.expectError {
				if err == nil {
					t.Error("Expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v\nOutput: %s", err, string(output))
				}

				// Just verify that we got some output
				if len(output) == 0 {
					t.Error("Expected some output but got none")
				}

				// Verify output contains basic Go struct syntax
				outputStr := string(output)
				if !strings.Contains(outputStr, "type") || !strings.Contains(outputStr, "struct") {
					t.Errorf("Output doesn't appear to be valid Go struct code: %s", outputStr)
				}
			}
		})
	}
}

func TestMain_ErrorCases(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		yamlContent string
		expectError bool
	}{
		{
			name:        "no arguments",
			args:        []string{},
			expectError: true,
		},
		{
			name:        "non-existent file",
			args:        []string{"nonexistent.yaml"},
			expectError: true,
		},
		{
			name:        "invalid YAML syntax",
			args:        []string{"invalid.yaml"},
			yamlContent: "invalid: yaml: content: [unclosed",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tmpFileName string

			// Create temp file if yamlContent is provided
			if tt.yamlContent != "" {
				tmpFile, err := os.CreateTemp("", "test*.yaml")
				if err != nil {
					t.Fatalf("Failed to create temp file: %v", err)
				}
				defer os.Remove(tmpFile.Name())

				if _, err := tmpFile.WriteString(tt.yamlContent); err != nil {
					t.Fatalf("Failed to write to temp file: %v", err)
				}
				tmpFile.Close()
				tmpFileName = tmpFile.Name()

				// Update args to use actual temp file name
				for i, arg := range tt.args {
					if strings.HasSuffix(arg, ".yaml") && arg != "nonexistent.yaml" {
						tt.args[i] = tmpFileName
					}
				}
			}

			// Build the command to run the main program
			cmd := exec.Command("go", append([]string{"run", "main.go"}, tt.args...)...)
			cmd.Dir = "."

			// Run the command and capture output
			output, err := cmd.CombinedOutput()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error but got none. Output: %s", string(output))
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v\nOutput: %s", err, string(output))
				}
			}
		})
	}
}
