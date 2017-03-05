package eval

import (
	"lib/ast"
	"lib/token"
	"math"
)

// strongly stripped down call to c()
// Therefore, all elements are evaluated withon the context of the call
// TODO recursive=TRUE/FALSE
// TODO faster vector literals, composed just of floats

func EvalC(ev *Evaluator, node *ast.CallExpr) (r *SEXP) {
	TRACE := ev.trace
	if TRACE {
		println("VectorExpr")
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

	return &SEXP{ValuePos: node.Fun.Pos(), TypeOf: REALSXP, Kind: token.FLOAT, Array: &c}
}

func intMin(x int, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

func intMax(x int, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}

// TODO work on slices of same length instead of single values
func fPLUS(x float64, y float64) float64 { 
	return x + y
}
func fMINUS(x float64, y float64) float64 { 
	return x - y
}
func fMULTIPLICATION(x float64, y float64) float64 { 
	return x * y
}
func fDIVISION(x float64, y float64) float64 { 
	return x / y
}
func fMODULUS(x float64, y float64) float64 { 
	return math.Mod(x, y)
}
func fEXPONENTIATION(x float64, y float64) float64 { 
	return math.Pow(x, y)
}

func EvalVectorOp(x *SEXP, y *SEXP, FUN func(float64, float64) float64) *SEXP {
	xv := *x.Array
	yv := *y.Array
	lenx := len(xv)
	leny := len(yv)
	sliceLen := intMin(len(xv),len(yv))
	resultLen := intMax(len(xv),len(yv))
	
	r := make([]float64,resultLen)

	for base := 0; base < resultLen; base += sliceLen {
		for i := base ; (i < (base+sliceLen) && i < resultLen); i++ {
			r[i] = FUN(xv[i % lenx], yv[i % leny])
		}
	}
	return &SEXP{Kind: token.FLOAT, Array: &r}
}
