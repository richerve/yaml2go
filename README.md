# Yaml 2 Go

The purpose of this program is to generate Go code from an input in yaml format.

The project is also used as a semi-testbed for agentic coding trying different models and approaches.

## Specification

- When the input is a string, the program will generate a string literal.
- When the input is a number, the program will generate a number literal.
- When the input is a list, the program will generate a slice.
- When the input is a boolean, the program will generate a boolean literal.
- When the input is a map, if it is populated and have items under it, the program will generate a struct.
- For each yaml document read from the input, a root level struct "Document#" will be created, where # is an int starting from 1.
- If the document has only one map key and all remaining items are under that key. The name of the initial struct will be the name of that key.
- The yaml types `string`, `number` or `boolean` are represented as a pointer to the corresponding Go type.
- Empty yaml values: `""`, `[]`, `{}`, `0`, must have an `omitempty` json tag flag. When passing the `-use-omitzero` cli flag, the `omitzero` json tag flag is used instead.
  - if the yaml value is `[]` is represented as a `[]any` in Go.
  - if the yaml value `{}` is represented as a `map[string]any`

## Examples

### Input

```yaml
mykey:
  myvalue1: 1
  myvalue2: two
foo: bar
items:
  - a
  - b
  - c
emptymap: {}
emptyslice: []
emptystring: ""
emptyint: 0
```

### Output

```Go
type Document struct {
	MyKey MyKey `json:"mykey"`
	Foo *string `json:"foo"`
	Items []string `json:"items"`
	EmptyMap map[string]any `json:"emptymap,omitempty"`
	EmptySlice []string `json:"emptyslice,omitempty"`
	EmptyString *string `json:"emptystring,omitempty"`
	EmptyInt *int `json:"emptyint,omitempty"`
}

type MyKey struct {
	MyValue1 int `json:"myvalue1"`
	MyValue2 string `json:"myvalue2"`
}
```
