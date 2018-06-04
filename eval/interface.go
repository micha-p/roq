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
	MAJOR := "0"
	MINOR := "0.0"
	EvalMain(&filename, src, parser.AllErrors, TRACE, DEBUG, MAJOR, MINOR)
}

func EvalFileForTest(filename string){
	TRACE := false
	DEBUG := false
	MAJOR := "0"
	MINOR := "0.0"
	EvalMain(&filename, nil, parser.AllErrors, TRACE, DEBUG, MAJOR, MINOR)
}


// parser might be started with filename or various other sources (string, []byte, *bytes.Buffer, io.Reader)
func EvalPreInit(filePtr *string, src interface{}, parserOpts parser.Mode, TRACE bool, DEBUG bool, MAJOR string, MINOR string)(*parser.Parser, *Evaluator){

	fset := token.NewFileSet() // positions are relative to fset

	p, errp := parser.ParseInit(fset, *filePtr, src, parserOpts)
	if errp != nil {
		panic(errp)
	}
	ev, erre := EvalInit(fset, *filePtr, src, parser.AllErrors, TRACE, DEBUG, MAJOR, MINOR)
	if erre != nil {
		panic(erre)
	}
	return p, ev
}

// parser might be started with filename or various other sources (string, []byte, *bytes.Buffer, io.Reader)
func EvalMain(filePtr *string, src interface{}, parserOpts parser.Mode, TRACE bool, DEBUG bool, MAJOR string, MINOR string) {
	
	var p *parser.Parser
	var ev *Evaluator
	p, ev = EvalPreInit(filePtr, src, parserOpts, TRACE, DEBUG, MAJOR, MINOR)
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

func EvalStringforValue(src string) float64{
	filename := ""
	TRACE := false
	DEBUG := false
	MAJOR := "0"
	MINOR := "0.0"
	
	var p *parser.Parser
	var ev *Evaluator
	p, ev = EvalPreInit(&filename, src, parser.AllErrors, TRACE, DEBUG, MAJOR, MINOR)

	stmt, _ := parser.ParseIter(p)
	sexp := EvalStmt(ev, stmt)
	r := sexp.(*VSEXP)
	return r.Immediate
}
