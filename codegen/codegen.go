package codegen

import (
	"fmt"
	"strings"
)

type StructDef struct {
	Name   string
	Fields []FieldDef
}

type FieldDef struct {
	Name  string
	Type  string
	Flags []string
}

func NewFieldDef(name, ftype string) FieldDef {
	return FieldDef{
		Name: name,
		Type: ftype,
	}
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

func WriteField(builder *strings.Builder, field FieldDef, tagPrefix string) {
	ft := FieldTag{
		Prefix: tagPrefix,
		Value:  field.Name,
		Flags:  field.Flags,
	}
	tag := ft.String()

	fmt.Fprintf(builder, "%s %s %s", Capitalize(field.Name), field.Type, tag)
}

func WriteStruct(builder *strings.Builder, structDef StructDef, tagPrefix string) {
	fmt.Fprintf(builder, "type %s struct {\n", structDef.Name)
	for _, field := range structDef.Fields {
		builder.WriteString("\t")
		WriteField(builder, field, tagPrefix)
		builder.WriteString("\n")
	}
	builder.WriteString("}\n")
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
