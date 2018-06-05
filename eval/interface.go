package eval

import (
	"roq/lib/parser"
	"roq/lib/token"
)

// TODO EvalStringForValues => needed for tests with assert


func EvalStringForTest(src string){
	filename:=""
	TRACE := false
	DEBUG := false
	EvalMain(&filename, src, parser.AllErrors, TRACE, DEBUG)
}

func EvalFileForTest(filename string){
	TRACE := false
	DEBUG := false
	EvalMain(&filename, nil, parser.AllErrors, TRACE, DEBUG)
}


// parser might be started with filename or various other sources (string, []byte, *bytes.Buffer, io.Reader)
func EvalPreInit(filePtr *string, src interface{}, parserOpts parser.Mode, TRACE bool, DEBUG bool)(*parser.Parser, *Evaluator){

	fset := token.NewFileSet() // positions are relative to fset

	p, errp := parser.ParseInit(fset, *filePtr, src, parserOpts)
	if errp != nil {
		panic(errp)
	}
	ev, erre := EvalInit(fset, *filePtr, src, parser.AllErrors, TRACE, DEBUG)
	if erre != nil {
		panic(erre)
	}
	return p, ev
}

// parser might be started with filename or various other sources (string, []byte, *bytes.Buffer, io.Reader)
func EvalMain(filePtr *string, src interface{}, parserOpts parser.Mode, TRACE bool, DEBUG bool){
	var p *parser.Parser
	var ev *Evaluator
	p, ev = EvalPreInit(filePtr, src, parserOpts, TRACE, DEBUG)
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
			if ev.Invisible { 				// invisibility is stored in the evaluator and is set during assignment
				ev.Invisible = false		// unsetting invisiblity again
			} else {
				PrintResult(sexp)
			}
			if ev.state == eofState {
				if DEBUG {
					println("terminating...")
				}
				break
			}
		}
	}
}

func EvalStringForValue(src string) SEXPItf{
	filename := ""
	TRACE := false
	DEBUG := false
	
	var p *parser.Parser
	var ev *Evaluator
	p, ev = EvalPreInit(&filename, src, parser.AllErrors, TRACE, DEBUG)

	stmt, _ := parser.ParseIter(p)
	sexp := EvalStmt(ev, stmt)
	return sexp
}
