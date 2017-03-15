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

func EvalCombine(ev *Evaluator, node *ast.CallExpr) (r *SEXP) {
	TRACE := ev.trace
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

func EvalLength(ev *Evaluator, node *ast.CallExpr) (r *SEXP) {
	TRACE := ev.trace
	if TRACE {
		println("Length")
	}

	l:=len(node.Args)
	if l == 1 {
		val := EvalExpr(ev,node.Args[0])
		if val.Slice ==nil {
			return &SEXP{ValuePos: node.Fun.Pos(), TypeOf: REALSXP, Kind: token.FLOAT, Immediate: 1}
		} else {
			return &SEXP{ValuePos: node.Fun.Pos(), TypeOf: REALSXP, Kind: token.FLOAT, Immediate: float64(len(val.Slice))}
		}
	} else {
		println(l,"arguments passed to 'length' which requires 1") 
		return &SEXP{Kind: token.ILLEGAL}
	}
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

func mapIA(FUN func(float64, float64) float64, x float64, y []float64) []float64 {
	resultLen := len(y)
	r := make([]float64,resultLen)
	for n,value := range y {
		r[n]=FUN(x,value)
	}
	return r
}

func mapAI(FUN func(float64, float64) float64, x []float64, y float64) []float64 {
	resultLen := len(x)
	r := make([]float64,resultLen)
	for n,value := range x {
		r[n]=FUN(value,y)
	}
	return r
}

func mapAA(FUN func(float64, float64) float64, x []float64, y []float64) []float64 {
	lenx := len(x)
	leny := len(y)
	sliceLen := intMin(lenx,leny)
	resultLen := intMax(lenx,leny)
	
	r := make([]float64,resultLen)

	for base := 0; base < resultLen; base += sliceLen {
		for i := base ; (i < (base+sliceLen) && i < resultLen); i++ {
			r[i] = FUN(x[i % lenx], y[i % leny])
		}
	}
	return r
}

func EvalVectorOp(x *SEXP, y *SEXP, FUN func(float64, float64) float64) *SEXP {
	if x.Slice==nil && y.Slice==nil {
		return &SEXP{Kind: token.FLOAT, Immediate: FUN(x.Immediate,y.Immediate)}
	} else if x.Slice==nil {
		return &SEXP{Kind: token.FLOAT, Slice: mapIA(FUN,x.Immediate,y.Slice)}
	} else if y.Slice==nil {
		return &SEXP{Kind: token.FLOAT, Slice: mapAI(FUN,x.Slice,y.Immediate)}
	} else {
		return &SEXP{Kind: token.FLOAT, Slice: mapAA(FUN,x.Slice,y.Slice)}
	}
}

// FALSE is counted as zero, 
// TRUE as 1 in comparisons (this will cause different behaviour; TODO: Warnings
//
// Concatenation of comparisons:
// As evaluation is from left to right, y value has to be returned

func EvalComp(op token.Token, x *SEXP, y *SEXP) *SEXP {
	// false and true are not really the same. false is rather the base level.
	if x == nil || y == nil {
		return nil
	}
	if x.Kind == token.ILLEGAL || y.Kind == token.ILLEGAL {
		return &SEXP{Kind: token.ILLEGAL}
	}
	var o1,o2 float64
	if x.Slice==nil {
		o1 = x.Immediate
	} else {
		o1 = x.Slice[0]
	}
	if y.Slice==nil {
		o2 = y.Immediate
	} else {
		o2 = y.Slice[0]
	}
	// println("?",o1,op.String(),o2)
	switch op {
	case token.LESS:
		if o1 < o2 {
			return y
		} else {
			return nil
		}
	case token.LESSEQUAL:
		if o1 <= o2 {
			return y
		} else {
			return nil
		}
	case token.GREATER:
		if o1 > o2 {
			return y
		} else {
			return nil
		}
	case token.GREATEREQUAL:
		if o1 >= o2 {
			return y
		} else {
			return nil
		}
	case token.EQUAL:
		if o1 == o2 {
			return y
		} else {
			return nil
		}
	case token.UNEQUAL:
		if o1 != o2 {
			return y
		} else {
			return nil
		}
	default:
		println("?Comp: " + op.String())
		return &SEXP{Kind: token.ILLEGAL}
	}
}

func EvalOp(op token.Token, x *SEXP, y *SEXP) *SEXP {
	if x.Kind == token.ILLEGAL || y.Kind == token.ILLEGAL {
		return &SEXP{Kind: token.ILLEGAL}
	}
	switch op {
	case token.PLUS:
		return EvalVectorOp(x,y,fPLUS)
	case token.MINUS:
		return EvalVectorOp(x,y,fMINUS)
	case token.MULTIPLICATION:
		return EvalVectorOp(x,y,fMULTIPLICATION)
	case token.DIVISION:
		return EvalVectorOp(x,y,fDIVISION)
	case token.EXPONENTIATION:
		return EvalVectorOp(x,y,fEXPONENTIATION)
	case token.MODULUS:
		return EvalVectorOp(x,y,fMODULUS)
	default:
		println("?Op: " + op.String())
		return &SEXP{Kind: token.ILLEGAL}
	}
}
