package codegen

import (
	"testing"
)

func TestFieldTag_String(t *testing.T) {
	tests := []struct {
		name     string
		fieldTag FieldTag
		expected string
	}{
		{
			name: "basic tag",
			fieldTag: FieldTag{
				Prefix: "json",
				Value:  "username",
				Flags:  []string{},
			},
			expected: "`json:\"username\"`",
		},
		{
			name: "tag with single flag",
			fieldTag: FieldTag{
				Prefix: "json",
				Value:  "email",
				Flags:  []string{"omitempty"},
			},
			expected: "`json:\"email,omitempty\"`",
		},
		{
			name: "tag with multiple flags",
			fieldTag: FieldTag{
				Prefix: "json",
				Value:  "password",
				Flags:  []string{"omitempty", "readonly"},
			},
			expected: "`json:\"password,omitempty,readonly\"`",
		},
		{
			name: "empty prefix",
			fieldTag: FieldTag{
				Prefix: "",
				Value:  "username",
				Flags:  []string{},
			},
			expected: "",
		},
		{
			name: "empty value",
			fieldTag: FieldTag{
				Prefix: "json",
				Value:  "",
				Flags:  []string{},
			},
			expected: "",
		},
		{
			name: "empty prefix and value",
			fieldTag: FieldTag{
				Prefix: "",
				Value:  "",
				Flags:  []string{"omitempty"},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fieldTag.String()
			if result != tt.expected {
				t.Errorf("FieldTag.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFieldDef_String(t *testing.T) {
	tests := []struct {
		name     string
		field    FieldDef
		expected string
	}{
		{
			name: "basic field",
			field: FieldDef{
				Name: "username",
				Type: "string",
				Tag: &FieldTag{
					Prefix: "json",
					Value:  "username",
					Flags:  []string{},
				},
			},
			expected: "Username string `json:\"username\"`",
		},
		{
			name: "field with flags",
			field: FieldDef{
				Name: "email",
				Type: "string",
				Tag: &FieldTag{
					Prefix: "json",
					Value:  "email",
					Flags:  []string{"omitempty"},
				},
			},
			expected: "Email string `json:\"email,omitempty\"`",
		},
		{
			name: "field with multiple flags",
			field: FieldDef{
				Name: "password",
				Type: "string",
				Tag: &FieldTag{
					Prefix: "json",
					Value:  "password",
					Flags:  []string{"omitempty", "readonly"},
				},
			},
			expected: "Password string `json:\"password,omitempty,readonly\"`",
		},
		{
			name: "snake_case field name",
			field: FieldDef{
				Name: "user_name",
				Type: "string",
				Tag: &FieldTag{
					Prefix: "yaml",
					Value:  "user_name",
					Flags:  []string{},
				},
			},
			expected: "UserName string `yaml:\"user_name\"`",
		},
		{
			name: "kebab-case field name",
			field: FieldDef{
				Name: "user-name",
				Type: "string",
				Tag: &FieldTag{
					Prefix: "json",
					Value:  "user-name",
					Flags:  []string{},
				},
			},
			expected: "UserName string `json:\"user-name\"`",
		},
		{
			name: "complex type",
			field: FieldDef{
				Name: "items",
				Type: "[]Item",
				Tag: &FieldTag{
					Prefix: "json",
					Value:  "items",
					Flags:  []string{},
				},
			},
			expected: "Items []Item `json:\"items\"`",
		},
		{
			name: "field with nil tag",
			field: FieldDef{
				Name: "username",
				Type: "string",
				Tag:  nil,
			},
			expected: "Username string",
		},
		{
			name: "field with empty tag",
			field: FieldDef{
				Name: "username",
				Type: "string",
				Tag: &FieldTag{
					Prefix: "",
					Value:  "",
					Flags:  []string{},
				},
			},
			expected: "Username string",
		},
		{
			name: "empty field name",
			field: FieldDef{
				Name: "",
				Type: "string",
				Tag: &FieldTag{
					Prefix: "json",
					Value:  "field",
					Flags:  []string{},
				},
			},
			expected: " string `json:\"field\"`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.field.String()
			if result != tt.expected {
				t.Errorf("FieldDef.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestStructDef_String(t *testing.T) {
	tests := []struct {
		name      string
		structDef StructDef
		expected  string
	}{
		{
			name: "basic struct",
			structDef: StructDef{
				Name: "User",
				Fields: []FieldDef{
					{
						Name: "username",
						Type: "string",
						Tag: &FieldTag{
							Prefix: "json",
							Value:  "username",
							Flags:  []string{},
						},
					},
					{
						Name: "email",
						Type: "string",
						Tag: &FieldTag{
							Prefix: "json",
							Value:  "email",
							Flags:  []string{},
						},
					},
				},
			},
			expected: `type User struct {
	Username string ` + "`json:\"username\"`" + `
	Email string ` + "`json:\"email\"`" + `
}
`,
		},
		{
			name: "struct with flags",
			structDef: StructDef{
				Name: "User",
				Fields: []FieldDef{
					{
						Name: "username",
						Type: "string",
						Tag: &FieldTag{
							Prefix: "json",
							Value:  "username",
							Flags:  []string{},
						},
					},
					{
						Name: "email",
						Type: "string",
						Tag: &FieldTag{
							Prefix: "json",
							Value:  "email",
							Flags:  []string{"omitempty"},
						},
					},
				},
			},
			expected: `type User struct {
	Username string ` + "`json:\"username\"`" + `
	Email string ` + "`json:\"email,omitempty\"`" + `
}
`,
		},
		{
			name: "empty struct",
			structDef: StructDef{
				Name:   "Empty",
				Fields: []FieldDef{},
			},
			expected: `type Empty struct {
}
`,
		},
		{
			name: "struct with complex types",
			structDef: StructDef{
				Name: "Document",
				Fields: []FieldDef{
					{
						Name: "title",
						Type: "string",
						Tag: &FieldTag{
							Prefix: "yaml",
							Value:  "title",
							Flags:  []string{},
						},
					},
					{
						Name: "tags",
						Type: "[]string",
						Tag: &FieldTag{
							Prefix: "yaml",
							Value:  "tags",
							Flags:  []string{},
						},
					},
					{
						Name: "metadata",
						Type: "map[string]interface{}",
						Tag: &FieldTag{
							Prefix: "yaml",
							Value:  "metadata",
							Flags:  []string{},
						},
					},
					{
						Name: "nested",
						Type: "NestedStruct",
						Tag: &FieldTag{
							Prefix: "yaml",
							Value:  "nested",
							Flags:  []string{},
						},
					},
				},
			},
			expected: `type Document struct {
	Title string ` + "`yaml:\"title\"`" + `
	Tags []string ` + "`yaml:\"tags\"`" + `
	Metadata map[string]interface{} ` + "`yaml:\"metadata\"`" + `
	Nested NestedStruct ` + "`yaml:\"nested\"`" + `
}
`,
		},
		{
			name: "struct with snake_case fields",
			structDef: StructDef{
				Name: "Config",
				Fields: []FieldDef{
					{
						Name: "api_key",
						Type: "string",
						Tag: &FieldTag{
							Prefix: "json",
							Value:  "api_key",
							Flags:  []string{},
						},
					},
					{
						Name: "base_url",
						Type: "string",
						Tag: &FieldTag{
							Prefix: "json",
							Value:  "base_url",
							Flags:  []string{},
						},
					},
					{
						Name: "timeout_seconds",
						Type: "int",
						Tag: &FieldTag{
							Prefix: "json",
							Value:  "timeout_seconds",
							Flags:  []string{},
						},
					},
				},
			},
			expected: `type Config struct {
	ApiKey string ` + "`json:\"api_key\"`" + `
	BaseUrl string ` + "`json:\"base_url\"`" + `
	TimeoutSeconds int ` + "`json:\"timeout_seconds\"`" + `
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.structDef.String()
			if result != tt.expected {
				t.Errorf("StructDef.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCapitalize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic string",
			input:    "hello",
			expected: "Hello",
		},
		{
			name:     "snake_case",
			input:    "user_name",
			expected: "UserName",
		},
		{
			name:     "kebab-case",
			input:    "user-name",
			expected: "UserName",
		},
		{
			name:     "mixed separators",
			input:    "user_name-id",
			expected: "UserNameId",
		},
		{
			name:     "multiple underscores",
			input:    "user__name",
			expected: "UserName",
		},
		{
			name:     "multiple dashes",
			input:    "user--name",
			expected: "UserName",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single character",
			input:    "a",
			expected: "A",
		},
		{
			name:     "already capitalized",
			input:    "UserName",
			expected: "UserName",
		},
		{
			name:     "with numbers",
			input:    "user_id_123",
			expected: "UserId123",
		},
		{
			name:     "starts with separator",
			input:    "_user_name",
			expected: "UserName",
		},
		{
			name:     "ends with separator",
			input:    "user_name_",
			expected: "UserName",
		},
		{
			name:     "only separators",
			input:    "___",
			expected: "",
		},
		{
			name:     "complex case",
			input:    "api_key_config_v2",
			expected: "ApiKeyConfigV2",
		},
		{
			name:     "camelCase input",
			input:    "userName",
			expected: "UserName",
		},
		{
			name:     "PascalCase input",
			input:    "UserName",
			expected: "UserName",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Capitalize(tt.input)
			if result != tt.expected {
				t.Errorf("Capitalize(%v) = %v, want %v", result, result, tt.expected)
			}
		})
	}
}
