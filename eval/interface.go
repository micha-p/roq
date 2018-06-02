package eval

import (
	"fmt"
	"roq/lib/parser"
	"roq/lib/token"
)

// TODO EvalStringForValues => needed for tests with assert

func EvalStringForTest(src string){
	filename:=""
	EvalMain(&filename, src, parser.AllErrors, false, false,"0","0.0")
}

func EvalFileForTest(filename string){
	EvalMain(&filename, nil, parser.AllErrors, false, false,"0","0.0")
}

// parser might be started with filename or various other sources (string, []byte, *bytes.Buffer, io.Reader)
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
		if tok == token.EOF {
			if DEBUG {
				println("EOF token found")
			}
			break
		}
		sexp := EvalStmt(ev, stmt)
		if sexp != nil {
			PrintResult(ev, sexp)
			if ev.state == eofState {
				if DEBUG {
					println("terminating...")
				}
				break
			}
		}
	}
}

