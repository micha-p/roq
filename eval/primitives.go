package eval

import (
	"fmt"
	"roq/lib/ast"
	"strings"
)

// https://cran.r-project.org/doc/manuals/R-ints.html#g_t_002eInternal-vs-_002ePrimitive
// these functions are responsible for evaluation

func EvalLength(ev *Evaluator, node *ast.CallExpr) (r *ISEXP) {
	TRACE := ev.Trace
	if TRACE {
		println("Length")
	}
	ex := node.Args[0]
	switch ex.(type) {
	case *ast.IndexExpr:
		iterator := IndexDomainEval(ev, ex.(*ast.IndexExpr).Index)
		return &ISEXP{ValuePos: node.Fun.Pos(), Integer: iterator.Length()}
	default:
		val := EvalExpr(ev, node.Args[0])
		return &ISEXP{ValuePos: node.Fun.Pos(), Integer: val.Length()}
	}
}


func EvalPrint(ev *Evaluator, node *ast.CallExpr) (r SEXPItf) {
	TRACE := ev.Trace
	if TRACE {
		println("PrintExpr")
	}
	value := EvalExpr(ev, node.Args[0])
	PrintResult(value)
	ev.Invisible = true
	return nil
}


func EvalCat(ev *Evaluator, node *ast.CallExpr) (r SEXPItf) {
	TRACE := ev.Trace
	if TRACE {
		println("CatExpr")
	}
	for n := 0; n < len(node.Args); n++ {
		r = EvalExpr(ev, node.Args[n])
		if n > 0 {
			fmt.Printf(" ")
		}
		switch r.(type) {
		case *TSEXP:
			fmt.Printf(strings.Replace(r.(*TSEXP).String, "\\n", "\n", -1)) // needs strings.Map
		case *ISEXP:
			fmt.Printf("%g", r.(*ISEXP).Immediate) // TODO
		case *VSEXP:
			if r.(*VSEXP).Slice == nil {
				fmt.Printf("%g",r.(*VSEXP).Immediate)
			} else {
				for n, v := range r.(*VSEXP).Slice {
					if n > 0 {
						fmt.Printf(" ")
					}
					fmt.Printf(" %g", v) // R has small e for exponential format
				}
			}
		default:
			println("?CAT")
		}
	}
	ev.Invisible = true
	return nil
}

// strongly stripped down call of c()
// Therefore, all elements are evaluated within the context of the call
// TODO recursive=TRUE/FALSE
// TODO faster vector literals, composed just of floats

// TODO document difference!
// - R returns double, if there is at least one double, here we decide on first type
// - inside string vectors, there is no conversion of numbers but an error thrown!

/* The output type is determined from the highest type of the
   components in the hierarchy NULL < raw < logical < integer <
   double < complex < character < list < expression.
*/

func EvalColumn(ev *Evaluator, node *ast.CallExpr) (r SEXPItf) {
	TRACE := ev.Trace
	if TRACE {
		println("Column")
	}

	if len(node.Args) > 0 {
		evaluatedArgs := make(map[int]SEXPItf)
		for n, v := range node.Args { // TODO: strictly left to right
			val := EvalExprOrAssignment(ev, v)
			evaluatedArgs[n] = val
		}
		switch evaluatedArgs[0].(type) {
		case *ISEXP:
			c := make([]int, len(evaluatedArgs))
			for n, v := range evaluatedArgs {
				c[n] = v.(*ISEXP).Integer
			}
			return &ISEXP{ValuePos: node.Fun.Pos(), Slice: c}
		case *VSEXP:
			c := make([]float64, len(evaluatedArgs))
			for n, v := range evaluatedArgs {
				switch v.(type) {
				case *VSEXP:
					c[n] = v.(*VSEXP).Immediate
				case *ISEXP:
					c[n] = v.(*ISEXP).Immediate
				default:
					panic("Error in c")
				}
			}
			return &VSEXP{ValuePos: node.Fun.Pos(), Slice: c}
		case *TSEXP:
			c := make([]string, len(evaluatedArgs))
			for n, v := range evaluatedArgs {
				c[n] = v.(*TSEXP).String
			}
			return &TSEXP{ValuePos: node.Fun.Pos(), Slice: c}
		default:
			panic("Error in function c") // TODO
		}
	} else {
		return nil
	}
}

func EvalList(ev *Evaluator, node *ast.CallExpr) (r *RSEXP) {
	TRACE := ev.Trace
	DEBUG := ev.Debug
	if TRACE {
		println("list")
	}
	if DEBUG {
		println("process given arguments for list function")
	}
	evaluatedArgs := EvalArgswithDotDotArguments(ev, "list", node.Args)
	if DEBUG {
		println("List of evaluated args for function: list")
		PrintListofSExpressions(evaluatedArgs)
	}
	return &RSEXP{ValuePos: node.Fun.Pos(), Slice: evaluatedArgs}
}

// TODO documentation and comparison
// pairlist might be called with more than 2 arguments
func EvalPairlist(ev *Evaluator, node *ast.CallExpr) (r *RSEXP) {
	TRACE := ev.Trace
	if TRACE {
		println("Pairlist")
	}

	return &RSEXP{ValuePos: node.Fun.Pos(),
		CAR: EvalExprOrAssignment(ev, node.Args[0]),
		CDR: EvalExprOrAssignment(ev, node.Args[1])}
}

func EvalTypeof(ev *Evaluator, node *ast.CallExpr) (r *TSEXP) {
	if arityOK("typeof", 1, node) {
		object := EvalExpr(ev, node.Args[0])
		var r string
		if object == nil {
			r = "NULL"
		} else {
			switch object.(type) {
			case *VSEXP:
				if object.(*VSEXP).Body == nil {
					r = "double"
				} else {
					r = "closure"
				}
			case *ISEXP:
				r = "integer"
//				case *LSEXP:
//					r="logical"
			case *TSEXP:
				r = "character"
			case *RSEXP:
				if object.(*RSEXP).Slice == nil {
					r = "pairlist"
				} else {
					r = "list"
				}
			case *NSEXP:
				r = "NULL"
			default:
				panic("unknown type")
			}
			return &TSEXP{String: r}
		}
	}
	return
}

func EvalClass(ev *Evaluator, node *ast.CallExpr) (r *TSEXP) {
	if arityOK("class", 1, node) {
		object := EvalExpr(ev, node.Args[0])
		s := object.Class()
		if s == nil {
			var r string
			switch object.(type) {
			case *VSEXP:
				if object.(*VSEXP).Body == nil {
					r = "numeric"
				} else {
					r = "function"
				}
			case *ISEXP:
				r = "numeric"
//				case *LSEXP:
//					r="logical"
			case *TSEXP:
				r = "character"
			case *RSEXP:
				if object.(*RSEXP).Slice == nil {
					r = "pairlist"
				} else {
					r = "list"
				}
			case *NSEXP:
				r = "NULL"
			default:
				panic("unknown type")
			}
			return &TSEXP{String: r}
		}
		return &TSEXP{String: *s}
	}
	return
}
