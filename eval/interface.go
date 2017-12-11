package eval

import (
	"fmt"
	"roq/lib/parser"
	"roq/lib/token"
)

func EvalStringForTest(src interface{}){
        filename:=""
	EvalMain(&filename, src, parser.AllErrors, false, false,"0","0.0")
}


func EvalMain(filePtr *string, src interface{}, parserOpts parser.Mode, TRACE bool, DEBUG bool, MAJOR string, MINOR string) {
	fset := token.NewFileSet() // positions are relative to fset
	p, errp := parser.ParseInit(fset, *filePtr, src, parserOpts)
	if errp != nil {
		fmt.Println(errp)
		return
	}
	ev, erre := EvalInit(fset, *filePtr, src, parser.AllErrors, TRACE, DEBUG, MAJOR, MINOR)
	if erre != nil {
		fmt.Println(erre)
		return
	}

	for true {
		stmt, tok := parser.ParseIter(p) // main iterator calls parse.stmt
		sexp := EvalStmt(ev, stmt)
		if sexp != nil {
			PrintResult(ev, sexp)
		}
		parser.StartLine(p)
		if tok == token.EOF {
			return
		}
	}
}

