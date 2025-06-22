package inference

import (
	"testing"

	"github.com/goccy/go-yaml/parser"
	"github.com/richerve/yaml2go/codegen"
)

func TestDetermineType(t *testing.T) {
	tests := []struct {
		name      string
		yamlInput string
		fieldName string
		structs   map[string]codegen.StructDef
		path      []string
		expected  string
	}{
		{
			name:      "string node",
			yamlInput: `"hello world"`,
			fieldName: "message",
			structs:   make(map[string]codegen.StructDef),
			path:      []string{},
			expected:  "string",
		},
		{
			name:      "integer node",
			yamlInput: `42`,
			fieldName: "count",
			structs:   make(map[string]codegen.StructDef),
			path:      []string{},
			expected:  "int",
		},
		{
			name:      "float node",
			yamlInput: `3.14`,
			fieldName: "pi",
			structs:   make(map[string]codegen.StructDef),
			path:      []string{},
			expected:  "float64",
		},
		{
			name:      "boolean node true",
			yamlInput: `true`,
			fieldName: "enabled",
			structs:   make(map[string]codegen.StructDef),
			path:      []string{},
			expected:  "bool",
		},
		{
			name:      "boolean node false",
			yamlInput: `false`,
			fieldName: "disabled",
			structs:   make(map[string]codegen.StructDef),
			path:      []string{},
			expected:  "bool",
		},
		{
			name:      "null node",
			yamlInput: `null`,
			fieldName: "nullable",
			structs:   make(map[string]codegen.StructDef),
			path:      []string{},
			expected:  "interface{}",
		},
		{
			name:      "empty sequence",
			yamlInput: `[]`,
			fieldName: "items",
			structs:   make(map[string]codegen.StructDef),
			path:      []string{},
			expected:  "[]interface{}",
		},
		{
			name:      "string sequence",
			yamlInput: `["item1", "item2"]`,
			fieldName: "items",
			structs:   make(map[string]codegen.StructDef),
			path:      []string{},
			expected:  "[]string",
		},
		{
			name:      "integer sequence",
			yamlInput: `[1, 2, 3]`,
			fieldName: "numbers",
			structs:   make(map[string]codegen.StructDef),
			path:      []string{},
			expected:  "[]int",
		},
		{
			name:      "float sequence",
			yamlInput: `[1.1, 2.2, 3.3]`,
			fieldName: "floats",
			structs:   make(map[string]codegen.StructDef),
			path:      []string{},
			expected:  "[]float64",
		},
		{
			name:      "boolean sequence",
			yamlInput: `[true, false]`,
			fieldName: "flags",
			structs:   make(map[string]codegen.StructDef),
			path:      []string{},
			expected:  "[]bool",
		},
		{
			name:      "empty mapping",
			yamlInput: `{}`,
			fieldName: "config",
			structs:   make(map[string]codegen.StructDef),
			path:      []string{},
			expected:  "map[string]any",
		},
		{
			name:      "non-empty mapping",
			yamlInput: `{name: "test"}`,
			fieldName: "user",
			structs:   make(map[string]codegen.StructDef),
			path:      []string{},
			expected:  "User",
		},
		{
			name:      "mapping with empty field name",
			yamlInput: `{key: "value"}`,
			fieldName: "",
			structs:   make(map[string]codegen.StructDef),
			path:      []string{},
			expected:  "NestedStruct",
		},
		{
			name:      "nested sequence of mappings",
			yamlInput: `[{name: "user1"}, {name: "user2"}]`,
			fieldName: "users",
			structs:   make(map[string]codegen.StructDef),
			path:      []string{},
			expected:  "[]NestedStruct",
		},
		{
			name:      "sequence with mixed types (first is string)",
			yamlInput: `["string", 123]`,
			fieldName: "mixed",
			structs:   make(map[string]codegen.StructDef),
			path:      []string{},
			expected:  "[]string",
		},
		{
			name:      "complex nested structure",
			yamlInput: `{user: {name: "john", age: 30}}`,
			fieldName: "data",
			structs:   make(map[string]codegen.StructDef),
			path:      []string{},
			expected:  "Data",
		},
		{
			name:      "snake_case field name",
			yamlInput: `{key: "value"}`,
			fieldName: "user_data",
			structs:   make(map[string]codegen.StructDef),
			path:      []string{},
			expected:  "UserData",
		},
		{
			name:      "kebab-case field name",
			yamlInput: `{key: "value"}`,
			fieldName: "user-data",
			structs:   make(map[string]codegen.StructDef),
			path:      []string{},
			expected:  "UserData",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse YAML input to get AST node
			file, err := parser.ParseBytes([]byte(tt.yamlInput), 0)
			if err != nil {
				t.Fatalf("Failed to parse YAML: %v", err)
			}

			if len(file.Docs) == 0 {
				t.Fatalf("No documents found in parsed YAML")
			}

			node := file.Docs[0].Body
			result := DetermineType(node, tt.fieldName, tt.structs, tt.path)

			if result != tt.expected {
				t.Errorf("DetermineType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDetermineType_UnknownNode(t *testing.T) {
	// Test with a mock node type that's not handled
	structs := make(map[string]codegen.StructDef)
	path := []string{}

	// For this test, we'll just verify that unknown nodes return interface{}
	// In practice, most nodes will be one of the handled types
	result := DetermineType(nil, "unknown", structs, path)
	expected := "interface{}"

	if result != expected {
		t.Errorf("DetermineType() with unknown node = %v, want %v", result, expected)
	}
}

func TestIsEmptyValue(t *testing.T) {
	tests := []struct {
		name      string
		yamlInput string
		expected  bool
	}{
		{
			name:      "empty string",
			yamlInput: `""`,
			expected:  true,
		},
		{
			name:      "non-empty string",
			yamlInput: `"hello"`,
			expected:  false,
		},
		{
			name:      "empty sequence",
			yamlInput: `[]`,
			expected:  true,
		},
		{
			name:      "non-empty sequence",
			yamlInput: `["item"]`,
			expected:  false,
		},
		{
			name:      "empty mapping",
			yamlInput: `{}`,
			expected:  true,
		},
		{
			name:      "non-empty mapping",
			yamlInput: `{key: "value"}`,
			expected:  false,
		},
		{
			name:      "integer (not empty)",
			yamlInput: `42`,
			expected:  false,
		},
		{
			name:      "zero integer (not empty)",
			yamlInput: `0`,
			expected:  false,
		},
		{
			name:      "float (not empty)",
			yamlInput: `3.14`,
			expected:  false,
		},
		{
			name:      "zero float (not empty)",
			yamlInput: `0.0`,
			expected:  false,
		},
		{
			name:      "boolean true (not empty)",
			yamlInput: `true`,
			expected:  false,
		},
		{
			name:      "boolean false (not empty)",
			yamlInput: `false`,
			expected:  false,
		},
		{
			name:      "null (not empty)",
			yamlInput: `null`,
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse YAML input to get AST node
			file, err := parser.ParseBytes([]byte(tt.yamlInput), 0)
			if err != nil {
				t.Fatalf("Failed to parse YAML: %v", err)
			}

			if len(file.Docs) == 0 {
				t.Fatalf("No documents found in parsed YAML")
			}

			node := file.Docs[0].Body
			result := IsEmptyValue(node)

			if result != tt.expected {
				t.Errorf("IsEmptyValue() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsEmptyValue_UnknownNode(t *testing.T) {
	// Test with nil node
	result := IsEmptyValue(nil)
	expected := false

	if result != expected {
		t.Errorf("IsEmptyValue() with nil node = %v, want %v", result, expected)
	}
}

func TestDetermineType_WithExistingStructs(t *testing.T) {
	tests := []struct {
		name      string
		yamlInput string
		fieldName string
		structs   map[string]codegen.StructDef
		path      []string
		expected  string
	}{
		{
			name:      "mapping with existing struct definitions",
			yamlInput: `{name: "john", age: 30}`,
			fieldName: "user",
			structs: map[string]codegen.StructDef{
				"User": {
					Name: "User",
					Fields: []codegen.FieldDef{
						{Name: "name", Type: "string"},
						{Name: "age", Type: "int"},
					},
				},
			},
			path:     []string{},
			expected: "User",
		},
		{
			name:      "sequence of mappings with complex path",
			yamlInput: `[{id: 1}, {id: 2}]`,
			fieldName: "items",
			structs:   make(map[string]codegen.StructDef),
			path:      []string{"root", "nested"},
			expected:  "[]NestedStruct",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse YAML input to get AST node
			file, err := parser.ParseBytes([]byte(tt.yamlInput), 0)
			if err != nil {
				t.Fatalf("Failed to parse YAML: %v", err)
			}

			if len(file.Docs) == 0 {
				t.Fatalf("No documents found in parsed YAML")
			}

			node := file.Docs[0].Body
			result := DetermineType(node, tt.fieldName, tt.structs, tt.path)

			if result != tt.expected {
				t.Errorf("DetermineType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDetermineType_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		yamlInput string
		fieldName string
		structs   map[string]codegen.StructDef
		path      []string
		expected  string
	}{
		{
			name:      "very long path",
			yamlInput: `{key: "value"}`,
			fieldName: "deep",
			structs:   make(map[string]codegen.StructDef),
			path:      []string{"level1", "level2", "level3", "level4", "level5"},
			expected:  "Deep",
		},
		{
			name:      "nil structs map",
			yamlInput: `"test"`,
			fieldName: "field",
			structs:   nil,
			path:      []string{},
			expected:  "string",
		},
		{
			name:      "sequence with null element",
			yamlInput: `[null]`,
			fieldName: "nulls",
			structs:   make(map[string]codegen.StructDef),
			path:      []string{},
			expected:  "[]interface{}",
		},
		{
			name:      "deeply nested sequence",
			yamlInput: `[[["nested"]]]`,
			fieldName: "nested",
			structs:   make(map[string]codegen.StructDef),
			path:      []string{},
			expected:  "[][][]string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse YAML input to get AST node
			file, err := parser.ParseBytes([]byte(tt.yamlInput), 0)
			if err != nil {
				t.Fatalf("Failed to parse YAML: %v", err)
			}

			if len(file.Docs) == 0 {
				t.Fatalf("No documents found in parsed YAML")
			}

			node := file.Docs[0].Body
			result := DetermineType(node, tt.fieldName, tt.structs, tt.path)

			if result != tt.expected {
				t.Errorf("DetermineType() = %v, want %v", result, tt.expected)
			}
		})
	}
}
