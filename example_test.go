package underlying_test

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"

	"github.com/tenntenn/underlying"
)

func ExampleConvert() {
	const src = `package main
	type MyInt int
	type Example struct {
		N MyInt
		S *string
	}
	func main() {}`

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "main.go", src, 0)
	if err != nil {
		panic(err)
	}

	info := &types.Info{
		Defs: map[*ast.Ident]types.Object{},
	}

	config := &types.Config{
		Importer: importer.Default(),
	}

	pkg, err := config.Check("main", fset, []*ast.File{f}, info)
	if err != nil {
		panic(err)
	}

	typ := pkg.Scope().Lookup("Example").Type()
	fmt.Println(underlying.Convert(typ, true))

	// Output: struct{N int; S string}
}
