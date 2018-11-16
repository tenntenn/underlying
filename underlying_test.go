package underlying_test

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"github.com/tenntenn/underlying"
)

func TestConvert(t *testing.T) {
	cases := []struct {
		input  string
		expect string
	}{
		{"string", "string"},
		{"*string", "string"},
		{"MyInt", "int"},
		{"*MyInt", "int"},
		{"[]MyInt", "[]int"},
		{"[2]MyInt", "[2]int"},
		{"map[*string]MyInt", "map[string]int"},
		{"func(MyInt)", "func(int)"},
		{"<-chan MyInt", "<-chan int"},
	}

	for _, tc := range cases {
		tc := tc
		n := fmt.Sprintf("%s -> %s", tc.input, tc.expect)
		t.Run(n, func(t *testing.T) {

			src := fmt.Sprintf(`package main
				type MyInt int
				var input %s
				var expect %s
				func main() {}`, tc.input, tc.expect)

			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "main.go", src, 0)
			if err != nil {
				t.Fatalf("unexpected error:%v", err)
			}

			info := &types.Info{
				Defs: map[*ast.Ident]types.Object{},
			}

			config := &types.Config{
				Importer: importer.Default(),
			}

			pkg, err := config.Check("main", fset, []*ast.File{f}, info)
			if err != nil {
				t.Fatalf("unexpected error:%v", err)
			}

			input := pkg.Scope().Lookup("input").Type()
			actual := underlying.Convert(input, true)
			expect := pkg.Scope().Lookup("expect").Type()

			if !types.Identical(expect, actual) {
				t.Errorf("expect %v but got %v", expect, actual)
			}
		})
	}
}
