package generator

import (
	"strings"
	"testing"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
)

// New is a simple constructor - no test needed

func TestGenerator_Generate(t *testing.T) {
	tests := []struct {
		name      string
		yamlInput string
		tagPrefix string
		expected  string
	}{
		{
			name: "simple single document",
			yamlInput: `
name: "john"
age: 30
`,
			tagPrefix: "json",
			expected: `type Document struct {
	Name string ` + "`json:\"name\"`" + `
	Age int ` + "`json:\"age\"`" + `
}
`,
		},
		{
			name: "single document with single key",
			yamlInput: `
user:
  name: "john"
  age: 30
`,
			tagPrefix: "json",
			expected: `type User struct {
	Name string ` + "`json:\"name\"`" + `
	Age int ` + "`json:\"age\"`" + `
}
`,
		},
		{
			name: "multiple documents",
			yamlInput: `
name: "doc1"
---
title: "doc2"
`,
			tagPrefix: "json",
			expected: `type Name struct {
	Name string ` + "`json:\"name\"`" + `
}

type Title struct {
	Title string ` + "`json:\"title\"`" + `
}
`,
		},
		{
			name: "nested structures",
			yamlInput: `
user:
  name: "john"
  profile:
    age: 30
    active: true
`,
			tagPrefix: "yaml",
			expected: `type User struct {
	Name string ` + "`yaml:\"name\"`" + `
	Profile Profile ` + "`yaml:\"profile\"`" + `
}

type Profile struct {
	Age int ` + "`yaml:\"age\"`" + `
	Active bool ` + "`yaml:\"active\"`" + `
}
`,
		},
		{
			name: "arrays and complex types",
			yamlInput: `
items: ["item1", "item2"]
count: 42
active: true
metadata: {}
`,
			tagPrefix: "json",
			expected: `type Document struct {
	Items []string ` + "`json:\"items\"`" + `
	Count int ` + "`json:\"count\"`" + `
	Active bool ` + "`json:\"active\"`" + `
	Metadata map[string]any ` + "`json:\"metadata,omitempty\"`" + `
}
`,
		},
		{
			name: "empty values with omitempty",
			yamlInput: `
name: "test"
empty_string: ""
empty_array: []
normal_field: "value"
`,
			tagPrefix: "json",
			expected: `type Document struct {
	Name string ` + "`json:\"name\"`" + `
	EmptyString string ` + "`json:\"empty_string,omitempty\"`" + `
	EmptyArray []interface{} ` + "`json:\"empty_array,omitempty\"`" + `
	NormalField string ` + "`json:\"normal_field\"`" + `
}
`,
		},
		{
			name: "snake_case to PascalCase conversion",
			yamlInput: `
user_name: "john"
email_address: "john@example.com"
is_active: true
user_id: 123
`,
			tagPrefix: "json",
			expected: `type Document struct {
	UserName string ` + "`json:\"user_name\"`" + `
	EmailAddress string ` + "`json:\"email_address\"`" + `
	IsActive bool ` + "`json:\"is_active\"`" + `
	UserId int ` + "`json:\"user_id\"`" + `
}
`,
		},
		{
			name: "kebab-case to PascalCase conversion",
			yamlInput: `
user-name: "john"
email-address: "john@example.com"
is-active: true
`,
			tagPrefix: "yaml",
			expected: `type Document struct {
	UserName string ` + "`yaml:\"user-name\"`" + `
	EmailAddress string ` + "`yaml:\"email-address\"`" + `
	IsActive bool ` + "`yaml:\"is-active\"`" + `
}
`,
		},
		{
			name: "multiple documents with single keys",
			yamlInput: `
config:
  database: "localhost"
  port: 5432
---
user:
  name: "admin"
  role: "administrator"
`,
			tagPrefix: "json",
			expected: `type Config struct {
	Database string ` + "`json:\"database\"`" + `
	Port int ` + "`json:\"port\"`" + `
}

type User struct {
	Name string ` + "`json:\"name\"`" + `
	Role string ` + "`json:\"role\"`" + `
}
`,
		},
		{
			name: "complex nested with arrays",
			yamlInput: `
users:
  - name: "john"
    age: 30
  - name: "jane"
    age: 25
settings:
  theme: "dark"
  features: ["feature1", "feature2"]
`,
			tagPrefix: "json",
			expected: `type Document struct {
	Users []NestedStruct ` + "`json:\"users\"`" + `
	Settings Settings ` + "`json:\"settings\"`" + `
}

type Settings struct {
	Theme string ` + "`json:\"theme\"`" + `
	Features []string ` + "`json:\"features\"`" + `
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := parser.ParseBytes([]byte(tt.yamlInput), 0)
			if err != nil {
				t.Fatalf("Failed to parse YAML: %v", err)
			}

			gen := New()
			result := gen.Generate(file, tt.tagPrefix)

			// Normalize whitespace for comparison
			normalizeWhitespace := func(s string) string {
				return strings.TrimSpace(strings.ReplaceAll(s, "\t", "    "))
			}

			expectedNorm := normalizeWhitespace(tt.expected)
			resultNorm := normalizeWhitespace(result)

			if resultNorm != expectedNorm {
				t.Errorf("Generate() result mismatch:\nExpected:\n%s\n\nGot:\n%s", tt.expected, result)
			}
		})
	}
}

func TestGenerator_determineDocumentName(t *testing.T) {
	tests := []struct {
		name      string
		yamlInput string
		index     int
		totalDocs int
		expected  string
	}{
		{
			name: "single document with multiple keys",
			yamlInput: `
name: "test"
age: 30
`,
			index:     0,
			totalDocs: 1,
			expected:  "Document",
		},
		{
			name: "single document with single key",
			yamlInput: `
user:
  name: "john"
  age: 30
`,
			index:     0,
			totalDocs: 1,
			expected:  "User",
		},
		{
			name: "multiple documents - first doc with single key",
			yamlInput: `
config:
  database: "localhost"
`,
			index:     0,
			totalDocs: 2,
			expected:  "Config",
		},
		{
			name: "multiple documents - second doc with multiple keys",
			yamlInput: `
name: "test"
value: 42
`,
			index:     1,
			totalDocs: 2,
			expected:  "Document2",
		},
		{
			name: "snake_case key name",
			yamlInput: `
user_config:
  setting: "value"
`,
			index:     0,
			totalDocs: 1,
			expected:  "UserConfig",
		},
		{
			name: "kebab-case key name",
			yamlInput: `
user-settings:
  theme: "dark"
`,
			index:     0,
			totalDocs: 1,
			expected:  "UserSettings",
		},
		{
			name:      "empty document",
			yamlInput: ``,
			index:     0,
			totalDocs: 1,
			expected:  "Document",
		},
		{
			name:      "document with empty mapping",
			yamlInput: `{}`,
			index:     0,
			totalDocs: 2,
			expected:  "Document1",
		},
		{
			name: "third document in sequence",
			yamlInput: `
name: "doc3"
`,
			index:     2,
			totalDocs: 3,
			expected:  "Name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := parser.ParseBytes([]byte(tt.yamlInput), 0)
			if err != nil {
				t.Fatalf("Failed to parse YAML: %v", err)
			}

			gen := New()

			var doc *ast.DocumentNode
			if len(file.Docs) > 0 {
				doc = file.Docs[0]
			} else {
				// Create empty document for empty input case
				doc = &ast.DocumentNode{}
			}

			result := gen.determineDocumentName(doc, tt.index, tt.totalDocs)

			if result != tt.expected {
				t.Errorf("determineDocumentName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGenerator_Generate_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		yamlInput string
		tagPrefix string
		expected  string
	}{
		{
			name:      "completely empty YAML",
			yamlInput: ``,
			tagPrefix: "json",
			expected:  "",
		},
		{
			name: "only comments",
			yamlInput: `
# This is a comment
# Another comment
`,
			tagPrefix: "json",
			expected:  "",
		},
		{
			name:      "null document",
			yamlInput: `null`,
			tagPrefix: "json",
			expected:  "",
		},
		{
			name: "array at root level",
			yamlInput: `
- item1
- item2
`,
			tagPrefix: "json",
			expected:  "",
		},
		{
			name:      "string at root level",
			yamlInput: `"just a string"`,
			tagPrefix: "json",
			expected:  "",
		},
		{
			name:      "number at root level",
			yamlInput: `42`,
			tagPrefix: "json",
			expected:  "",
		},
		{
			name:      "boolean at root level",
			yamlInput: `true`,
			tagPrefix: "json",
			expected:  "",
		},
		{
			name: "mixed document types",
			yamlInput: `
user:
  name: "john"
---
"just a string"
---
count: 42
active: true
`,
			tagPrefix: "json",
			expected: `type User struct {
	Name string ` + "`json:\"name\"`" + `
}

type Document3 struct {
	Count int ` + "`json:\"count\"`" + `
	Active bool ` + "`json:\"active\"`" + `
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := parser.ParseBytes([]byte(tt.yamlInput), 0)
			if err != nil {
				t.Fatalf("Failed to parse YAML: %v", err)
			}

			gen := New()
			result := gen.Generate(file, tt.tagPrefix)

			// Normalize whitespace for comparison
			normalizeWhitespace := func(s string) string {
				return strings.TrimSpace(strings.ReplaceAll(s, "\t", "    "))
			}

			expectedNorm := normalizeWhitespace(tt.expected)
			resultNorm := normalizeWhitespace(result)

			if resultNorm != expectedNorm {
				t.Errorf("Generate() result mismatch:\nExpected:\n%s\n\nGot:\n%s", tt.expected, result)
			}
		})
	}
}

func TestGenerator_Generate_StructOrdering(t *testing.T) {
	tests := []struct {
		name      string
		yamlInput string
		tagPrefix string
		checkFunc func(t *testing.T, result string)
	}{
		{
			name: "root struct appears first",
			yamlInput: `
user:
  profile:
    settings:
      theme: "dark"
`,
			tagPrefix: "json",
			checkFunc: func(t *testing.T, result string) {
				lines := strings.Split(result, "\n")
				if len(lines) < 1 {
					t.Error("Expected at least one line of output")
					return
				}
				if !strings.Contains(lines[0], "type User struct") {
					t.Errorf("Expected User struct to be first, got: %s", lines[0])
				}
			},
		},
		{
			name: "nested structs appear after root",
			yamlInput: `
app:
  database:
    host: "localhost"
    port: 5432
  cache:
    redis:
      url: "redis://localhost"
`,
			tagPrefix: "json",
			checkFunc: func(t *testing.T, result string) {
				appIndex := strings.Index(result, "type App struct")
				databaseIndex := strings.Index(result, "type Database struct")
				cacheIndex := strings.Index(result, "type Cache struct")

				if appIndex == -1 {
					t.Error("App struct not found")
					return
				}
				if databaseIndex != -1 && databaseIndex < appIndex {
					t.Error("Database struct should appear after App struct")
				}
				if cacheIndex != -1 && cacheIndex < appIndex {
					t.Error("Cache struct should appear after App struct")
				}
			},
		},
		{
			name: "multiple documents preserve order",
			yamlInput: `
config:
  setting: "value"
---
user:
  name: "john"
---
data:
  items: ["a", "b"]
`,
			tagPrefix: "json",
			checkFunc: func(t *testing.T, result string) {
				configIndex := strings.Index(result, "type Config struct")
				userIndex := strings.Index(result, "type User struct")
				dataIndex := strings.Index(result, "type Data struct")

				if configIndex == -1 || userIndex == -1 || dataIndex == -1 {
					t.Error("Not all expected structs found")
					return
				}

				if configIndex > userIndex || userIndex > dataIndex {
					t.Error("Document structs should appear in order")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := parser.ParseBytes([]byte(tt.yamlInput), 0)
			if err != nil {
				t.Fatalf("Failed to parse YAML: %v", err)
			}

			gen := New()
			result := gen.Generate(file, tt.tagPrefix)

			tt.checkFunc(t, result)
		})
	}
}

func TestGenerator_Generate_DifferentTagPrefixes(t *testing.T) {
	yamlInput := `
name: "test"
age: 30
`

	tests := []struct {
		name      string
		tagPrefix string
		expected  string
	}{
		{
			name:      "json tag prefix",
			tagPrefix: "json",
			expected:  "`json:\"name\"`",
		},
		{
			name:      "yaml tag prefix",
			tagPrefix: "yaml",
			expected:  "`yaml:\"name\"`",
		},
		{
			name:      "xml tag prefix",
			tagPrefix: "xml",
			expected:  "`xml:\"name\"`",
		},
		{
			name:      "custom tag prefix",
			tagPrefix: "custom",
			expected:  "`custom:\"name\"`",
		},
		{
			name:      "empty tag prefix",
			tagPrefix: "",
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := parser.ParseBytes([]byte(yamlInput), 0)
			if err != nil {
				t.Fatalf("Failed to parse YAML: %v", err)
			}

			gen := New()
			result := gen.Generate(file, tt.tagPrefix)

			if tt.expected == "" {
				// For empty tag prefix, tags should be empty
				if strings.Contains(result, "`") {
					t.Error("Expected no tags for empty tag prefix")
				}
			} else {
				if !strings.Contains(result, tt.expected) {
					t.Errorf("Expected result to contain %s, got: %s", tt.expected, result)
				}
			}
		})
	}
}
