package visitor

import (
	"github.com/goccy/go-yaml/ast"
	"github.com/richerve/yaml2go/pkg/codegen"
	"github.com/richerve/yaml2go/pkg/inference"
)

type ASTVisitor struct {
	structs   map[string]codegen.StructDef
	path      []string
	tagPrefix string
}

func NewASTVisitor(structs map[string]codegen.StructDef, path []string) *ASTVisitor {
	return &ASTVisitor{
		structs: structs,
		path:    path,
	}
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
	var fields []codegen.FieldDef

	for _, mappingValue := range node.Values {
		keyNode := mappingValue.Key
		keyValue := keyNode.String()
		fieldName := keyValue
		fieldType := inference.DetermineType(mappingValue.Value, keyValue, v.structs, v.path)

		flags := []string{}
		// Check if value is empty and add omitempty tag
		if inference.IsEmptyValue(mappingValue.Value) {
			flags = append(flags, "omitempty")
		}

		fd := codegen.FieldDef{
			Name: fieldName,
			Type: fieldType,
			Tag: &codegen.FieldTag{
				Prefix: v.tagPrefix,
				Value:  fieldName,
				Flags:  flags,
			},
		}

		fields = append(fields, fd)
	}

	// Store struct definition
	v.structs[structName] = codegen.StructDef{
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
		structs: v.structs,
		path:    newPath,
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
	return codegen.Capitalize(lastKey)
}
