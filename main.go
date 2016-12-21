package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	//	"log"
	//	"strings"
	//	"os"
	"lib/go/ast"
	"lib/go/parser"
	"lib/go/scanner"
	"lib/go/token"
)

func myerrorhandler(pos token.Position, msg string) {
	println("SCANNER ERROR", pos.Filename, pos.Line, pos.Column, msg)
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

		p, err := parser.ParseInit(fset, *filePtr, nil, parser.AllErrors|parser.Trace)

		if err != nil {
			fmt.Println(err)
			return
		}

		
		for true {
			stmt, tok := parser.ParseIter(p) // main iterator calls parse.stmt
			ast.Print(fset, stmt)
			if tok == token.EOF {
				return
			}
		}
	}
}
