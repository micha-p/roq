package eval

import (
	"fmt"
	"strings"
	"lib/ast"
	"lib/token"
)

// https://cran.r-project.org/doc/manuals/R-ints.html#g_t_002eInternal-vs-_002ePrimitive

func EvalLength(ev *Evaluator, node *ast.CallExpr) (r *ISEXP) {
	TRACE := ev.Trace
	if TRACE {
		println("Length")
	}
	ex := node.Args[0]
	switch ex.(type) {
	case *ast.IndexExpr:
		iterator := IndexDomainEval(ev, ex.(*ast.IndexExpr).Index)
		return &ISEXP{ValuePos: node.Fun.Pos(), TypeOf: REALSXP, Integer: iterator.Length()}
	default:
		val := EvalExpr(ev,node.Args[0])
		if val.(*VSEXP).Slice ==nil {
			return &ISEXP{ValuePos: node.Fun.Pos(), TypeOf: REALSXP, Integer: 0}
		} else {
			return &ISEXP{ValuePos: node.Fun.Pos(), TypeOf: REALSXP, Integer: val.Length()}
		}
	}
}

func EvalCat(ev *Evaluator, node *ast.CallExpr) (r SEXPItf) {
	TRACE := ev.Trace
	if TRACE {
		println("PrintExpr")
	}
	for n := 0; n < len(node.Args); n++ {
		r = EvalExpr(ev, node.Args[n])
		if n > 0 {
			print(" ")
		}
		switch r.(type) {
		case *TSEXP:
			print(strings.Replace(r.(*TSEXP).String, "\\n", "\n", -1)) // needs strings.Map
		case *ISEXP:
			fmt.Printf("%g", r.(*ISEXP).Immediate)  // TODO
		case *VSEXP:
			if r.(*VSEXP).Slice==nil {
				print(r.(*VSEXP).Immediate)
			} else {
				for n, v := range r.(*VSEXP).Slice {
					if n>0 {
						print(" ")
					}
					fmt.Printf(" %g", v) // R has small e for exponential format
				}
			}
		default:
			println("?CAT")
		}
	}
	ev.Invisible = true
	return
}

// strongly stripped down call of c()
// Therefore, all elements are evaluated withon the context of the call
// TODO recursive=TRUE/FALSE
// TODO faster vector literals, composed just of floats

func EvalColumn(ev *Evaluator, node *ast.CallExpr) (r SEXPItf) {
	TRACE := ev.Trace
	if TRACE {
		println("Column")
	}

	if len(node.Args)>0 {
		evaluatedArgs := make(map[int]SEXPItf)
		for n, v := range node.Args { // TODO: strictly left to right
			val := EvalExprOrAssignment(ev, v)
			evaluatedArgs[n] = val
		}
		switch evaluatedArgs[0].(type){
		case *VSEXP:
			c := make([]float64, len(evaluatedArgs))
			for n,v := range evaluatedArgs {
				c[n] = v.(*VSEXP).Immediate
			}
			return &VSEXP{ValuePos: node.Fun.Pos(), TypeOf: REALSXP, Kind: token.FLOAT, Slice: c}
		case *TSEXP:
			c := make([]string, len(evaluatedArgs))
			for n,v := range evaluatedArgs {
				c[n] = v.(*TSEXP).String
			}
			return &TSEXP{ValuePos: node.Fun.Pos(), TypeOf: STRSXP, Slice: c}
		default:
			println("Error in c") // TODO
			return nil
		}
	} else {
		return nil
	}
}

func EvalList(ev *Evaluator, node *ast.CallExpr) (r *RSEXP) {
	TRACE := ev.Trace
	if TRACE {
		println("List")
	}

	evaluatedArgs := make([]SEXPItf,len(node.Args))
	for n, v := range node.Args { // TODO: strictly left to right
		val := EvalExprOrAssignment(ev, v)
		evaluatedArgs[n] = val
	}

	return &RSEXP{ValuePos: node.Fun.Pos(), TypeOf: VECSXP, Slice: evaluatedArgs}
}
