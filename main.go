package main

// libraries forked from Go 1.7.3

import (
	"flag"
	"fmt"
	"io/ioutil"
	//	"log"
	//	"strings"
	//	"os"
	"eval"
	"print"
	"lib/ast"
	"lib/parser"
	"lib/scanner"
	"lib/token"
)

var TRACE bool
var DEBUG bool
var ECHO bool

func myerrorhandler(pos token.Position, msg string) {
	println("SCANNER ERROR", pos.Filename, pos.Line, pos.Column, msg)
}

func main() {

	defer func() {
		if x := recover(); x !=nil && x != "quit" {
			fmt.Printf("run time panic: %v", x)
			println()
		}
	}()


	fset := token.NewFileSet() // positions are relative to fset

	scanPtr := flag.Bool("scan", false, "scan")
	parsePtr := flag.Bool("parse", false, "parse")
	traceLongPtr := flag.Bool("trace", false, "trace")
	traceFlagPtr := flag.Bool("T", false, "trace")
	debugLongPtr := flag.Bool("debug", false, "debug")
	debugFlagPtr := flag.Bool("D", false, "debug")
	echoLongPtr := flag.Bool("echo", false, "echo")
	echoFlagPtr := flag.Bool("E", false, "echo")
	filePtr := flag.String("file", "example.src", "filename to process")
	flag.Parse()

	TRACE = *traceFlagPtr || *traceLongPtr
	DEBUG = *debugFlagPtr || *debugLongPtr
	ECHO  = *echoFlagPtr || *echoLongPtr

	if *scanPtr {
		src, _ := ioutil.ReadFile(*filePtr)

		var s scanner.Scanner
		file := fset.AddFile(*filePtr, fset.Base(), len(src)) // register input "file"
		s.Init(file, src, myerrorhandler, ECHO)

		// Repeated calls to Scan yield the token sequence found in the input
		for {
			pos, tok, lit := s.Scan()
			if tok == token.EOF {
				break
			}
			fmt.Printf("%s\t%s\t%q\n", fset.Position(pos), tok, lit)
		}
	} else if *parsePtr {

		var parserOpts parser.Mode
		parserOpts = parser.AllErrors
		
		if TRACE {
			parserOpts = parserOpts | parser.Trace
		}
		if DEBUG {
			parserOpts = parserOpts | parser.Debug
		}
		if ECHO {
			parserOpts = parserOpts | parser.Echo
		}

		p, err := parser.ParseInit(fset, *filePtr, nil, parserOpts)
		if err != nil {
			fmt.Println(err)
			return
		}

		for true {
			stmt, tok := parser.ParseIter(p) // main iterator calls parse.stmt
			switch stmt.(type) {
			case *ast.EmptyStmt:
			default:
				ast.Print(fset, stmt)
				parser.StartLine(p)
			}
			if tok == token.EOF {
				return
			}
		}
	} else { // eval
		var parserOpts parser.Mode
		parserOpts = parser.AllErrors
		
		if TRACE {
			parserOpts = parserOpts | parser.Trace
		}
		if DEBUG {
			parserOpts = parserOpts | parser.Debug
		}
		if ECHO {
			parserOpts = parserOpts | parser.Echo
		}

		p, errp := parser.ParseInit(fset, *filePtr, nil, parserOpts)
		if errp != nil {
			fmt.Println(errp)
			return
		}
		ev, erre := eval.EvalInit(fset, *filePtr, nil, parser.AllErrors, TRACE, DEBUG)
		if erre != nil {
			fmt.Println(erre)
			return
		}

		for true {
			stmt, tok := parser.ParseIter(p) // main iterator calls parse.stmt
			sexp := eval.EvalStmt(ev, stmt)
			if ! (sexp==nil) {
				print.PrintResult(ev, sexp)
				parser.StartLine(p)
			}
			if tok == token.EOF {
				return
			}
		}
	}
}
