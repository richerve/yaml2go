package yaml2go

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/goccy/go-yaml/ast"
)

type Generator struct {
	structs map[string]StructDef
}

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

type ASTVisitor struct {
	generator *Generator
	path      []string
}

func NewGenerator() *Generator {
	return &Generator{
		structs: make(map[string]StructDef),
	}
}

func (g *Generator) Generate(file *ast.File, tagPrefix string) string {
	// Process each document in the file using Walk
	for i, doc := range file.Docs {
		rootName := g.determineDocumentName(doc, i, len(file.Docs))

		visitor := &ASTVisitor{
			generator: g,
			path:      []string{rootName},
		}
		ast.Walk(visitor, doc)
	}

	var result strings.Builder

	// Generate root structs first in order
	var rootNames []string
	for i, doc := range file.Docs {
		rootName := g.determineDocumentName(doc, i, len(file.Docs))
		if _, exists := g.structs[rootName]; exists {
			rootNames = append(rootNames, rootName)
		}
	}

	// Write root structs
	for i, rootName := range rootNames {
		if i > 0 {
			result.WriteString("\n")
		}
		WriteStruct(&result, g.structs[rootName], tagPrefix)
	}

	// Generate other structs in sorted order
	var otherNames []string
	for name := range g.structs {
		isRoot := slices.Contains(rootNames, name)
		if !isRoot {
			otherNames = append(otherNames, name)
		}
	}
	sort.Strings(otherNames)

	for _, name := range otherNames {
		result.WriteString("\n")
		WriteStruct(&result, g.structs[name], tagPrefix)
	}

	return result.String()
}

func (g *Generator) determineDocumentName(doc *ast.DocumentNode, index int, totalDocs int) string {
	// Check if document has only one top-level key
	if doc.Body != nil {
		if mappingNode, ok := doc.Body.(*ast.MappingNode); ok && len(mappingNode.Values) == 1 {
			// Single key document - use the key name as struct name
			firstMapping := mappingNode.Values[0]
			keyNode := firstMapping.Key
			keyValue := keyNode.String()
			return capitalize(keyValue)
		}
	}

	// Multiple keys, not a mapping, or no body - use default naming
	if totalDocs == 1 {
		return "Document"
	} else {
		return fmt.Sprintf("Document%d", index+1)
	}
}

func WriteField(builder *strings.Builder, field FieldDef, tagPrefix string) {
	ft := FieldTag{
		Prefix: tagPrefix,
		Value:  field.Name,
		Flags:  field.Flags,
	}
	tag := ft.String()

	fmt.Fprintf(builder, "%s %s %s", capitalize(field.Name), field.Type, tag)
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

func (v *ASTVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.DocumentNode:
		// Continue visiting the document body
		return v

	case *ast.MappingNode:
		return v.visitMappingNode(n)

	case *ast.MappingValueNode:
		return v.visitMappingValueNode(n)

	case *ast.SequenceNode:
		// Sequences are handled in determineType, don't traverse further
		return nil

	default:
		// For other node types, continue traversal
		return v
	}
}

func (v *ASTVisitor) visitMappingNode(node *ast.MappingNode) ast.Visitor {
	if len(node.Values) == 0 {
		// Empty map - don't process further
		return nil
	}

	// Non-empty mapping - create struct
	structName := v.getCurrentStructName()
	var fields []FieldDef

	for _, mappingValue := range node.Values {
		keyNode := mappingValue.Key
		keyValue := keyNode.String()
		fieldName := keyValue
		fieldType := v.determineType(mappingValue.Value, keyValue)

		flags := []string{}
		// Check if value is empty and add omitempty tag
		if v.isEmptyValue(mappingValue.Value) {
			flags = append(flags, "omitempty")
		}

		fd := FieldDef{
			Name:  fieldName,
			Type:  fieldType,
			Flags: flags,
		}

		fields = append(fields, fd)
	}

	// Store struct definition
	v.generator.structs[structName] = StructDef{
		Name:   structName,
		Fields: fields,
	}

	// Continue traversal to handle nested structures
	return v
}

func (v *ASTVisitor) visitMappingValueNode(node *ast.MappingValueNode) ast.Visitor {
	keyNode := node.Key
	keyValue := keyNode.String()
	// Create new visitor with updated path for nested structures
	newPath := make([]string, len(v.path))
	copy(newPath, v.path)
	newPath = append(newPath, keyValue)

	newVisitor := &ASTVisitor{
		generator: v.generator,
		path:      newPath,
	}

	// Walk the value with the updated path context
	ast.Walk(newVisitor, node.Value)

	// Don't continue walking from this visitor to avoid double processing
	return nil
}

func (v *ASTVisitor) getCurrentStructName() string {
	if len(v.path) == 0 {
		return "Document"
	}

	// Use the last element in path as struct name
	lastKey := v.path[len(v.path)-1]
	return capitalize(lastKey)
}

func (v *ASTVisitor) isEmptyValue(node ast.Node) bool {
	switch n := node.(type) {
	case *ast.StringNode:
		return n.Value == ""
	case *ast.SequenceNode:
		return len(n.Values) == 0
	case *ast.MappingNode:
		return len(n.Values) == 0
	default:
		return false
	}
}

func (v *ASTVisitor) determineType(node ast.Node, fieldName string) string {
	switch n := node.(type) {
	case *ast.StringNode:
		return "string"

	case *ast.IntegerNode:
		return "int"

	case *ast.FloatNode:
		return "float64"

	case *ast.BoolNode:
		return "bool"

	case *ast.NullNode:
		return "interface{}"

	case *ast.SequenceNode:
		if len(n.Values) == 0 {
			return "[]interface{}"
		}

		// Determine element type from first element
		elementType := v.determineType(n.Values[0], "")
		return "[]" + elementType

	case *ast.MappingNode:
		if len(n.Values) == 0 {
			return "map[string]any"
		}

		// Non-empty mapping becomes a struct
		structName := capitalize(fieldName)
		if structName == "" {
			structName = "NestedStruct"
		}

		// Create visitor for this nested structure with updated path
		newPath := make([]string, len(v.path))
		copy(newPath, v.path)
		newPath = append(newPath, fieldName)

		nestedVisitor := &ASTVisitor{
			generator: v.generator,
			path:      newPath,
		}

		// Walk the nested mapping to process its structure
		ast.Walk(nestedVisitor, n)

		return structName

	default:
		return "interface{}"
	}
}

func capitalize(s string) string {
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
