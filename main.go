package main

// libraries forked from Go 1.7.3

import (
	"flag"
	"fmt"
	"bufio"
	"os"
	"strconv"
	"io/ioutil"
	"roq/eval"
	"roq/lib/ast"
	"roq/lib/parser"
	"roq/lib/scanner"
	"roq/lib/token"
	"roq/version"
)

var TRACE bool
var DEBUG bool
var ECHO  bool

func myerrorhandler(pos token.Position, msg string) {
	println("SCANNER ERROR", pos.Filename, pos.Line, pos.Column, msg)
}

func mainScan(filePtr *string, srcString string, ECHO bool) {
	fset := token.NewFileSet() // positions are relative to fset
	var src []byte
	if srcString != "" {
		src = []byte(srcString)
	} else {
		src, _ = ioutil.ReadFile(*filePtr)
	}
	var s scanner.Scanner
	file := fset.AddFile(*filePtr, fset.Base(), len(src)) // register input "file"
	s.Init(file, src, myerrorhandler, ECHO)

	// Repeated calls to Scan yield the token sequence found in the input
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		if lit == "" {
			fmt.Printf("%s\t%s\t \n", fset.Position(pos), tok)
		} else {
			fmt.Printf("%s\t%s\t%q\n", fset.Position(pos), tok, lit)
		}
	}
}

func mainParse(filePtr *string, src interface{}, parserOpts parser.Mode) {
	fset := token.NewFileSet() // positions are relative to fset
	p, err := parser.ParseInit(fset, *filePtr, src, parserOpts)
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
		}
		if tok == token.EOF {
			return
		}
	}
}

func main() {

	scanPtr := flag.Bool("scan", false, "scan (file only)") // TODO stdin
	parsePtr := flag.Bool("parse", false, "parse (file only)") // TODO stdin
	versionPtr := flag.Bool("version", false, "version")
	traceLongPtr := flag.Bool("trace", false, "trace")
	traceFlagPtr := flag.Bool("T", false, "trace")
	debugLongPtr := flag.Bool("debug", false, "debug")
	debugFlagPtr := flag.Bool("D", false, "debug")
	echoLongPtr := flag.Bool("echo", false, "echo")
	echoFlagPtr := flag.Bool("E", false, "echo")
	filePtr := flag.String("file", "", "filename to process")
	exprPtr := flag.String("expr", "", "expression to process")
	flag.Parse()

	TRACE = *traceFlagPtr || *traceLongPtr
	DEBUG = *debugFlagPtr || *debugLongPtr
	ECHO = *echoFlagPtr || *echoLongPtr
	PRINT := true

	if *versionPtr {
		version.PrintVersion()
	}

	if DEBUG == false {  // FIXME
		defer func() {
			if x := recover(); x != nil && x != "quit" {
				fmt.Printf("run time panic: %v", x)
			}
		}()
	}

	var src interface{}
	if *exprPtr == "" {
		src = nil
	} else {
		filename := "EXPRESSION"
		filePtr = &filename
		src = *exprPtr
	}

	if *filePtr == "" && *exprPtr == "" {  // interactive
		println("roq version", 
			version.MAJOR+"."+version.MINOR, 
			"(" +
			strconv.Itoa(version.YEAR) + "-" +
			strconv.Itoa(version.MONTH) + "-" + 
			strconv.Itoa(version.DAY) + 
			")")
		reader := bufio.NewReader(os.Stdin)
		//for s, _, err:= reader.ReadLine(); err==nil; {
			//eval.EvalStringForTest(string(s))
		//}
		for true {
			fmt.Print("> ")
			s, _, err := reader.ReadLine()
			if err != nil{
				break
			}
			if ECHO {
				println("> "+string(s))
			}
			eval.EvalStringForTest(string(s))
		}
	} else if *scanPtr {
		mainScan(filePtr, *exprPtr, ECHO)
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
		mainParse(filePtr, src, parserOpts)
	} else {
		var parserOpts parser.Mode
		parserOpts = parser.AllErrors

		if TRACE {
			parserOpts = parserOpts | parser.Trace
		}

		if ECHO {
			parserOpts = parserOpts | parser.Echo
		}
		eval.EvalMain(filePtr, src, parserOpts, TRACE, DEBUG, PRINT)
	}
}
