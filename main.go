package main

import (
	"fmt"
	"flag"
	"io/ioutil"
//	"log"
//	"strings"
//	"os"
	"lib/go/scanner"
	"lib/go/parser"
	"lib/go/token"
	"lib/go/ast"
)

func myerrorhandler (pos token.Position, msg string){
	println("SCANNER ERROR",pos.Filename,pos.Line,pos.Column,msg)
}


func main() {
	fset := token.NewFileSet() // positions are relative to fset

	scanPtr := flag.Bool("scan", false, "scan instead of parse")
	filePtr := flag.String("file", "example.src", "filename to process")
	flag.Parse()

	if *scanPtr {

		src, _ := ioutil.ReadFile(*filePtr)

		var s scanner.Scanner
		file := fset.AddFile(*filePtr, fset.Base(), len(src)) // register input "file"
		s.Init(file, src, myerrorhandler)


		// Repeated calls to Scan yield the token sequence found in the input
		for {
			pos, tok, lit := s.Scan()
			if tok == token.EOF {
				break
			}
			fmt.Printf("%s\t%s\t%q\n", fset.Position(pos), tok, lit)
		}
	} else {
		// Parse the file containing this very example
		// but stop after processing the imports.
		f, err := parser.ParseFile(fset, *filePtr, nil,0)
		if err != nil {
			fmt.Println(err)
			return
		}
	
		// Print the AST.
		ast.Print(fset, f)
	}
}

//	p := parser.NewParser(os.Stdin)
