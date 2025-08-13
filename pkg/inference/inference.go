package inference

import (
	"github.com/goccy/go-yaml/ast"
	"github.com/richerve/yaml2go/pkg/codegen"
)

func DetermineType(node ast.Node, fieldName string, structs map[string]codegen.StructDef, path []string) string {
	switch n := node.(type) {
	case *ast.StringNode:
		return "*string"

	case *ast.IntegerNode:
		return "*int"

	case *ast.FloatNode:
		return "*float64"

	case *ast.BoolNode:
		return "*bool"

	case *ast.NullNode:
		return "interface{}"

	case *ast.SequenceNode:
		if len(n.Values) == 0 {
			return "[]any"
		}

		// Determine element type from first element
		elementType := DetermineType(n.Values[0], "", structs, path)
		// For arrays, use non-pointer versions of basic types
		switch elementType {
		case "*string":
			elementType = "string"
		case "*int":
			elementType = "int"
		case "*float64":
			elementType = "float64"
		case "*bool":
			elementType = "bool"
		}
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
	case *ast.IntegerNode:
		// Check if the string representation is "0"
		return n.String() == "0"
	case *ast.FloatNode:
		return n.String() == "0" || n.String() == "0.0"
	case *ast.SequenceNode:
		return len(n.Values) == 0
	case *ast.MappingNode:
		return len(n.Values) == 0
	default:
		return false
	}
}
