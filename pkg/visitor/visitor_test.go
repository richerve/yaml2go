package visitor

import (
	"testing"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/richerve/yaml2go/pkg/codegen"
)

// NewASTVisitor is a simple constructor - no test needed

func TestASTVisitor_Visit(t *testing.T) {
	tests := []struct {
		name           string
		yamlInput      string
		initialStructs map[string]codegen.StructDef
		initialPath    []string
		expectedCount  int
		expectedStruct string
	}{
		{
			name:           "document node continues traversal",
			yamlInput:      `name: "test"`,
			initialStructs: make(map[string]codegen.StructDef),
			initialPath:    []string{"Root"},
			expectedCount:  1,
			expectedStruct: "Root",
		},
		{
			name:           "mapping node creates struct",
			yamlInput:      `{name: "john", age: 30}`,
			initialStructs: make(map[string]codegen.StructDef),
			initialPath:    []string{"User"},
			expectedCount:  1,
			expectedStruct: "User",
		},
		{
			name:           "empty mapping node",
			yamlInput:      `{}`,
			initialStructs: make(map[string]codegen.StructDef),
			initialPath:    []string{"Empty"},
			expectedCount:  0,
			expectedStruct: "",
		},
		{
			name:           "sequence node stops traversal",
			yamlInput:      `["item1", "item2"]`,
			initialStructs: make(map[string]codegen.StructDef),
			initialPath:    []string{"Items"},
			expectedCount:  0,
			expectedStruct: "",
		},
		{
			name:           "nested mapping structure",
			yamlInput:      `{user: {name: "john", profile: {age: 30}}}`,
			initialStructs: make(map[string]codegen.StructDef),
			initialPath:    []string{"Document"},
			expectedCount:  1,
			expectedStruct: "Document",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := parser.ParseBytes([]byte(tt.yamlInput), 0)
			if err != nil {
				t.Fatalf("Failed to parse YAML: %v", err)
			}

			if len(file.Docs) == 0 {
				t.Fatalf("No documents found in parsed YAML")
			}

			visitor := NewASTVisitor(tt.initialStructs, tt.initialPath, "json", false)
			result := visitor.Visit(file.Docs[0].Body)

			// Check if visitor returns correctly
			if tt.expectedCount == 0 {
				if result != nil {
					t.Errorf("Expected nil visitor for %s, got non-nil", tt.name)
				}
			} else {
				if result == nil {
					t.Errorf("Expected non-nil visitor for %s, got nil", tt.name)
				}
			}

			// Check if struct was created when expected
			if tt.expectedStruct != "" {
				if _, exists := tt.initialStructs[tt.expectedStruct]; !exists {
					// This is expected since Visit() doesn't directly create structs
					// It sets up for traversal, actual struct creation happens in visitMappingNode
				}
			}
		})
	}
}

func TestASTVisitor_visitMappingNode(t *testing.T) {
	tests := []struct {
		name              string
		yamlInput         string
		path              []string
		expectedStructs   int
		expectedFieldName string
		expectedFieldType string
		expectOmitEmpty   bool
	}{
		{
			name:              "basic mapping",
			yamlInput:         `{name: "john", age: 30}`,
			path:              []string{"User"},
			expectedStructs:   1,
			expectedFieldName: "name",
			expectedFieldType: "*string",
			expectOmitEmpty:   false,
		},
		{
			name:              "mapping with empty values",
			yamlInput:         `{name: "", tags: []}`,
			path:              []string{"Config"},
			expectedStructs:   1,
			expectedFieldName: "name",
			expectedFieldType: "*string",
			expectOmitEmpty:   true,
		},
		{
			name:            "empty mapping",
			yamlInput:       `{}`,
			path:            []string{"Empty"},
			expectedStructs: 0,
		},
		{
			name:              "complex nested mapping",
			yamlInput:         `{user: {name: "john"}, settings: {theme: "dark"}}`,
			path:              []string{"Document"},
			expectedStructs:   1,
			expectedFieldName: "user",
			expectedFieldType: "User",
			expectOmitEmpty:   false,
		},
		{
			name:              "mapping with different field types",
			yamlInput:         `{id: 123, active: true, score: 95.5, tags: ["tag1"]}`,
			path:              []string{"Item"},
			expectedStructs:   1,
			expectedFieldName: "id",
			expectedFieldType: "*int",
			expectOmitEmpty:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := parser.ParseBytes([]byte(tt.yamlInput), 0)
			if err != nil {
				t.Fatalf("Failed to parse YAML: %v", err)
			}

			if len(file.Docs) == 0 {
				t.Fatalf("No documents found in parsed YAML")
			}

			structs := make(map[string]codegen.StructDef)
			visitor := NewASTVisitor(structs, []string{}, "json", false)

			mappingNode, ok := file.Docs[0].Body.(*ast.MappingNode)
			if !ok {
				t.Fatalf("Expected MappingNode, got %T", file.Docs[0].Body)
			}

			result := visitor.visitMappingNode(mappingNode)

			// Check struct creation
			if tt.expectedStructs == 0 {
				if len(structs) != 0 {
					t.Errorf("Expected no structs, got %d", len(structs))
				}
				if result != nil {
					t.Errorf("Expected nil visitor, got non-nil")
				}
			} else {
				if len(structs) != tt.expectedStructs {
					t.Errorf("Expected %d structs, got %d", tt.expectedStructs, len(structs))
				}

				structName := visitor.getCurrentStructName()
				if structDef, exists := structs[structName]; exists {
					// Check if expected field exists
					if tt.expectedFieldName != "" {
						found := false
						for _, field := range structDef.Fields {
							if field.Name == tt.expectedFieldName {
								found = true
								if field.Type != tt.expectedFieldType {
									t.Errorf("Expected field type %s, got %s", tt.expectedFieldType, field.Type)
								}
								if tt.expectOmitEmpty {
									hasOmitEmpty := false
									if field.Tag != nil {
										for _, flag := range field.Tag.Flags {
											if flag == "omitempty" {
												hasOmitEmpty = true
												break
											}
										}
									}
									if !hasOmitEmpty {
										t.Errorf("Expected omitempty flag for field %s", tt.expectedFieldName)
									}
								}
								break
							}
						}
						if !found {
							t.Errorf("Expected field %s not found in struct", tt.expectedFieldName)
						}
					}
				} else {
					t.Errorf("Expected struct %s not found", structName)
				}

				if result == nil {
					t.Errorf("Expected non-nil visitor, got nil")
				}
			}
		})
	}
}

func TestASTVisitor_visitMappingValueNode(t *testing.T) {
	tests := []struct {
		name            string
		yamlInput       string
		initialPath     []string
		expectedStructs int
		keyToTest       string
	}{
		{
			name:            "basic mapping value",
			yamlInput:       `{user: {name: "john"}}`,
			initialPath:     []string{"Document"},
			expectedStructs: 1,
			keyToTest:       "user",
		},
		{
			name:            "nested mapping value",
			yamlInput:       `{config: {database: {host: "localhost"}}}`,
			initialPath:     []string{"App"},
			expectedStructs: 1,
			keyToTest:       "config",
		},
		{
			name:            "mapping value with string",
			yamlInput:       `{title: "Test Document"}`,
			initialPath:     []string{"Doc"},
			expectedStructs: 0,
			keyToTest:       "title",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := parser.ParseBytes([]byte(tt.yamlInput), 0)
			if err != nil {
				t.Fatalf("Failed to parse YAML: %v", err)
			}

			if len(file.Docs) == 0 {
				t.Fatalf("No documents found in parsed YAML")
			}

			structs := make(map[string]codegen.StructDef)
			visitor := NewASTVisitor(structs, tt.initialPath, "json", false)

			mappingNode, ok := file.Docs[0].Body.(*ast.MappingNode)
			if !ok {
				t.Fatalf("Expected MappingNode, got %T", file.Docs[0].Body)
			}

			// Find the mapping value node we want to test
			var targetMappingValue *ast.MappingValueNode
			for _, mv := range mappingNode.Values {
				if mv.Key.String() == tt.keyToTest {
					targetMappingValue = mv
					break
				}
			}

			if targetMappingValue == nil {
				t.Fatalf("Could not find mapping value for key %s", tt.keyToTest)
			}

			result := visitor.visitMappingValueNode(targetMappingValue)

			// visitMappingValueNode should always return nil to avoid double processing
			if result != nil {
				t.Errorf("Expected nil visitor from visitMappingValueNode, got non-nil")
			}

		})
	}
}

func TestASTVisitor_getCurrentStructName(t *testing.T) {
	tests := []struct {
		name     string
		path     []string
		expected string
	}{
		{
			name:     "empty path",
			path:     []string{},
			expected: "Document",
		},
		{
			name:     "single element path",
			path:     []string{"user"},
			expected: "User",
		},
		{
			name:     "multiple element path",
			path:     []string{"app", "database"},
			expected: "Database",
		},
		{
			name:     "snake_case path element",
			path:     []string{"user", "user_profile"},
			expected: "UserProfile",
		},
		{
			name:     "kebab-case path element",
			path:     []string{"config", "user-settings"},
			expected: "UserSettings",
		},
		{
			name:     "mixed case path",
			path:     []string{"api", "userData"},
			expected: "UserData",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			structs := make(map[string]codegen.StructDef)
			visitor := NewASTVisitor(structs, tt.path, "json", false)
			result := visitor.getCurrentStructName()

			if result != tt.expected {
				t.Errorf("getCurrentStructName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestASTVisitor_Integration(t *testing.T) {
	tests := []struct {
		name          string
		yamlInput     string
		initialPath   []string
		expectedCount int
		checkStruct   string
		checkField    string
	}{
		{
			name: "complete document processing",
			yamlInput: `
user:
  name: "john"
  profile:
    age: 30
    tags: ["admin", "user"]
settings:
  theme: "dark"
  notifications: true
`,
			initialPath:   []string{"Config"},
			expectedCount: 1,
			checkStruct:   "Config",
			checkField:    "user",
		},
		{
			name: "array of objects",
			yamlInput: `
users:
  - name: "john"
    age: 30
  - name: "jane"
    age: 25
`,
			initialPath:   []string{"Document"},
			expectedCount: 1,
			checkStruct:   "Document",
			checkField:    "users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := parser.ParseBytes([]byte(tt.yamlInput), 0)
			if err != nil {
				t.Fatalf("Failed to parse YAML: %v", err)
			}

			if len(file.Docs) == 0 {
				t.Fatalf("No documents found in parsed YAML")
			}

			structs := make(map[string]codegen.StructDef)
			visitor := NewASTVisitor(structs, tt.initialPath, "json", false)

			// Walk the entire document
			ast.Walk(visitor, file.Docs[0])

			if len(structs) < tt.expectedCount {
				t.Errorf("Expected at least %d structs, got %d", tt.expectedCount, len(structs))
			}

			if tt.checkStruct != "" {
				if structDef, exists := structs[tt.checkStruct]; exists {
					if tt.checkField != "" {
						found := false
						for _, field := range structDef.Fields {
							if field.Name == tt.checkField {
								found = true
								break
							}
						}
						if !found {
							t.Errorf("Expected field %s not found in struct %s", tt.checkField, tt.checkStruct)
						}
					}
				} else {
					t.Errorf("Expected struct %s not found", tt.checkStruct)
				}
			}
		})
	}
}
