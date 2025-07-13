package codegen

import (
	"fmt"
	"strings"
)

type StructDef struct {
	Name   string
	Fields []FieldDef
}

func (s StructDef) String() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "type %s struct {\n", s.Name)
	for _, field := range s.Fields {
		builder.WriteString("\t")
		fmt.Fprint(&builder, field.String())
		builder.WriteString("\n")
	}
	builder.WriteString("}\n")

	return builder.String()
}

type FieldDef struct {
	Name string
	Type string
	Tag  *FieldTag
}

func (f FieldDef) String() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%s %s", Capitalize(f.Name), f.Type))
	if f.Tag != nil && f.Tag.String() != "" {
		builder.WriteString(fmt.Sprintf(" %s", f.Tag.String()))
	}
	return builder.String()
}

type FieldTag struct {
	Prefix string
	Value  string
	Flags  []string
}

func (f *FieldTag) String() string {
	var sb strings.Builder

	if f.Prefix == "" || f.Value == "" {
		return ""
	}

	fmt.Fprintf(&sb, "`%s:\"%s", f.Prefix, f.Value)
	if len(f.Flags) > 0 {
		for _, f := range f.Flags {
			fmt.Fprintf(&sb, ",%s", f)
		}
	}
	sb.WriteString("\"`")

	return sb.String()
}

func Capitalize(s string) string {
	if s == "" {
		return ""
	}

	// Convert snake_case or kebab-case to PascalCase
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-'
	})

	var result strings.Builder
	for _, part := range parts {
		if len(part) > 0 {
			result.WriteString(strings.ToUpper(string(part[0])))
			if len(part) > 1 {
				result.WriteString(part[1:])
			}
		}
	}

	return result.String()
}
