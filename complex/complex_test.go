package complex

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// https://github.com/knsh14/gocc/blob/9078b24a5eb4377455473212ec67b8034de1439f/complexity/complexity_test.go
func GetAST(t *testing.T, code string) ast.Node {
	t.Helper()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", code, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, decl := range file.Decls {
		if fd, ok := decl.(*ast.FuncDecl); ok {
			return fd
		}
	}
	t.Fatal("no function declear found")
	return nil
}

func TestCount(t *testing.T) {
	tests := []struct {
		name string
		code string
		want int
	}{
		{
			name: "if",
			code: `package main
			func Double(n int) int {
				return n * 2
			}`,
			want: 1,
		},
		{
			name: "if statement",
			code: `package main
			func Double(n int) int {
				if n%2 == 0 {
					return 0
				}
				return n
			}`,
			want: 2,
		},
		{
			name: "for statement",
			code: `package main
			func Sum(n int) int {
				c := 0
				for i := 0; i < n; i++ {
					c += i
				}
				return c
			}`,
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast := GetAST(t, tt.code)
			if got := Count(ast); got != tt.want {
				t.Errorf("Count() = %v, want %v", got, tt.want)
			}
		})
	}
}
