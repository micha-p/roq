package eval

import (
	"fmt"
	"strings"
	"lib/ast"
	"lib/token"
)
	
func EvalLength(ev *Evaluator, node *ast.CallExpr) (r *SEXP) {
	TRACE := ev.Trace
	if TRACE {
		println("Length")
	}
	l:=len(node.Args)
	if l == 1 {
		ex := node.Args[0]
		switch ex.(type) {
		case *ast.IndexExpr:
			iterator := IndexDomainEval(ev, ex.(*ast.IndexExpr).Index)
			return &SEXP{ValuePos: node.Fun.Pos(), TypeOf: REALSXP, Kind: token.INT, Offset: iterator.Length()}
		default:
			val := EvalExpr(ev,node.Args[0])
			if val.Slice ==nil {
				return &SEXP{ValuePos: node.Fun.Pos(), TypeOf: REALSXP, Kind: token.INT, Offset: 1}
			} else {
				return &SEXP{ValuePos: node.Fun.Pos(), TypeOf: REALSXP, Kind: token.INT, Offset: len(val.Slice)}
			}
		}
	} else {
		println(l,"arguments passed to 'length' which requires 1") 
		return &SEXP{Kind: token.ILLEGAL}
	}
}

func EvalCat(ev *Evaluator, node *ast.CallExpr) (r *SEXP) {
	TRACE := ev.Trace
	if TRACE {
		println("PrintExpr")
	}
	for n := 0; n < len(node.Args); n++ {
		r = EvalExpr(ev, node.Args[n])
		if n > 0 {
			print(" ")
		}
		switch r.Kind {
		case token.STRING:
			print(strings.Replace(r.String, "\\n", "\n", -1)) // needs strings.Map
		case token.INT:
			fmt.Printf("%g", r.Immediate)
		case token.FLOAT:
			if r.Slice==nil {
				print(r.Immediate)
			} else {
				for n, v := range r.Slice {
					if n>0 {
						print(" ")
					}
					fmt.Printf(" %g", v) // R has small e for exponential format
				}
			}
		default:
			println("?CAT", r.Kind.String())
		}
	}
	ev.Invisible = true
	return
}

// strongly stripped down call of c()
// Therefore, all elements are evaluated withon the context of the call
// TODO recursive=TRUE/FALSE
// TODO faster vector literals, composed just of floats

func EvalCombine(ev *Evaluator, node *ast.CallExpr) (r *SEXP) {
	TRACE := ev.Trace
	if TRACE {
		println("Combine")
	}

	evaluatedArgs := make(map[int]float64)
	for n, v := range node.Args { // TODO: strictly left to right
		val := EvalExprOrAssignment(ev, v)
		evaluatedArgs[n] = val.Immediate
	}
	c := make([]float64, len(evaluatedArgs))
	for n,v := range evaluatedArgs {
		c[n] = v
	}

	return &SEXP{ValuePos: node.Fun.Pos(), TypeOf: REALSXP, Kind: token.FLOAT, Slice: c}
}
