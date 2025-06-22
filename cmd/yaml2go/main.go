package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/goccy/go-yaml/parser"
	"github.com/richerve/yaml2go/generator"
)

func main() {
	var tagPrefix string
	flag.StringVar(&tagPrefix, "tag-prefix", "json", "tag prefix to use, default is json")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s <options> [yaml-file]\n", os.Args[0])
		os.Exit(1)
	}

	filename := flag.Arg(0)
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Parse YAML into AST
	file, err := parser.ParseBytes(data, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing YAML: %v\n", err)
		os.Exit(1)
	}

	gen := generator.New()
	fmt.Print(gen.Generate(file, tagPrefix))
}
