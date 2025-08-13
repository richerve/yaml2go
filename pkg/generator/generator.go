package generator

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/goccy/go-yaml/ast"
	"github.com/richerve/yaml2go/pkg/codegen"
	"github.com/richerve/yaml2go/pkg/visitor"
)

type Generator struct {
	structs map[string]codegen.StructDef
}

func New() *Generator {
	return &Generator{
		structs: make(map[string]codegen.StructDef),
	}
}

func (g *Generator) Generate(file *ast.File, tagPrefix string) string {
	// Process each document in the file using Walk
	for i, doc := range file.Docs {
		rootName := g.determineDocumentName(doc, i, len(file.Docs))

		v := visitor.NewASTVisitor(g.structs, []string{rootName}, tagPrefix)
		ast.Walk(v, doc)
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
		_, err := result.WriteString(g.structs[rootName].String())
		if err != nil {
			return ""
		}
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
		_, err := result.WriteString(g.structs[name].String())
		if err != nil {
			return ""
		}
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
			return codegen.Capitalize(keyValue)
		}
	}

	// Multiple keys, not a mapping, or no body - use default naming
	if totalDocs == 1 {
		return "Document"
	} else {
		return fmt.Sprintf("Document%d", index+1)
	}
}
