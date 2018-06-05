package eval

import (
	"roq/lib/parser"
	"roq/lib/token"
)

func EvalStringForValue(src string) SEXPItf{
	filename:=""
	TRACE := false
	DEBUG := false
	PRINT := false
	return EvalMain(&filename, src, parser.AllErrors, TRACE, DEBUG, PRINT)
}

func EvalStringForTest(src string){
	filename:=""
	TRACE := false
	DEBUG := false
	PRINT := true
	EvalMain(&filename, src, parser.AllErrors, TRACE, DEBUG, PRINT)
}

func EvalFileForTest(filename string){
	TRACE := false
	DEBUG := false
	PRINT := true
	EvalMain(&filename, nil, parser.AllErrors, TRACE, DEBUG, PRINT)
}

// parser might be started with filename or various other sources (string, []byte, *bytes.Buffer, io.Reader)
func EvalMain(filePtr *string, src interface{}, parserOpts parser.Mode, TRACE bool, DEBUG bool, PRINT bool) SEXPItf{
	var returnExpression SEXPItf

	fset := token.NewFileSet() // positions are relative to fset

	p, errp := parser.ParseInit(fset, *filePtr, src, parserOpts)
	if errp != nil {
		panic(errp)
	}
	ev, erre := EvalInit(fset, *filePtr, src, parser.AllErrors, TRACE, DEBUG)
	if erre != nil {
		panic(erre)
	}

	for true {
		stmt, tok := parser.ParseIter(p) 	// main iterator calls parse.stmt
		if tok == token.EOF {
			if DEBUG {
				println("EOF token found")
			}
			break
		}
		sexp := EvalStmt(ev, stmt)
		if sexp != nil {
			if ev.Invisible { 				// invisibility is stored in the evaluator and is set during assignment
				ev.Invisible = false		// unsetting invisiblity again
			} else if PRINT{
				PrintResult(sexp)
			}
			returnExpression = sexp
			if ev.state == eofState {
				if DEBUG {
					println("terminating...")
				}
				break
			}
		}
	}
	return returnExpression
}
