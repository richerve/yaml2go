# Yaml 2 Go

The purpose of this program is to generate Go code from an input in yaml format.

These are the rules that the program will follow:

* When the input is a string, the program will generate a string literal.
* When the input is a number, the program will generate a number literal.
* When the input is a list, the program will generate a slice.
* When the input is a boolean, the program will generate a boolean literal.
* When the input is a map, if it is populated and have items under it, the program will generate a struct.
* When the input is an empty map (`{}`), the program will generated a `map[string]any`
* For each yaml document read from the input, a root level struct "Document#" will be created, where # is an int starting from 1.
* If the document has only one map key and all remaining items are under that key. The name of the initial struct will be the name of that key.
* Any empty yaml value (`""`, `[]`, `{}`) will have the `omitempty` json tag flag.

## Example

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
empty: {}
```

### Output

```Go
type Document struct {
	MyKey MyKey `json:"mykey"`
	Foo string `json:"foo"`
	Items []string `json:"items"`
	Empty map[string]any `json:"empty"`
}

type MyKey struct {
	MyValue1 int `json:"myvalue1"`
	MyValue2 string `json:"myvalue2"`
}
```
