package inference

import (
	"github.com/richerve/yaml2go/codegen"

	"github.com/goccy/go-yaml/ast"
)

func DetermineType(node ast.Node, fieldName string, structs map[string]codegen.StructDef, path []string) string {
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
		elementType := DetermineType(n.Values[0], "", structs, path)
		return "[]" + elementType

	case *ast.MappingNode:
		if len(n.Values) == 0 {
			return "map[string]any"
		}

		// Non-empty mapping becomes a struct
		structName := codegen.Capitalize(fieldName)
		if structName == "" {
			structName = "NestedStruct"
		}

		return structName

	default:
		return "interface{}"
	}
}

func IsEmptyValue(node ast.Node) bool {
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
