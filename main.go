package main

import (
	"fmt"
	//	"io/ioutil"
	//	"log"
	//s	"strings"
//	"os"
	"lib/go/parser"
	"lib/go/token"
	"lib/go/ast"
)

func main() {
	fset := token.NewFileSet() // positions are relative to fset

	// Parse the file containing this very example
	// but stop after processing the imports.
	f, err := parser.ParseFile(fset, "example.src", nil,0)
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
