package main

import (
	"fmt"
	//	"io/ioutil"
	//	"log"
	//s	"strings"
//	"os"
	"go/parser"
	"go/token"
	"go/ast"
)

func main() {
	fset := token.NewFileSet() // positions are relative to fset

	// Parse the file containing this very example
	// but stop after processing the imports.
	f, err := parser.ParseFile(fset, "lib/go/parser/example_test.go", nil, parser.ImportsOnly)
	if err != nil {
		fmt.Println(err)
		return
	}

// Print the AST.
	ast.Print(fset, f)

/*
	// Print the imports from the file's AST.
	for _, s := range f.Imports {
		fmt.Println(s.Path.Value)
	}
*/
}

//	p := parser.NewParser(os.Stdin)
