package eval

import (
	"roq/lib/token"
	"math"
)


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

func EvalVectorOp(x *VSEXP, y *VSEXP, FUN func(float64, float64) float64) *VSEXP {
	if x.Slice==nil && y.Slice==nil {
		return &VSEXP{Immediate: FUN(x.Immediate,y.Immediate)}
	} else if x.Slice==nil {
		return &VSEXP{Slice: mapIA(FUN,x.Immediate,y.Slice)}
	} else if y.Slice==nil {
		return &VSEXP{Slice: mapAI(FUN,x.Slice,y.Immediate)}
	} else {
		return &VSEXP{Slice: mapAA(FUN,x.Slice,y.Slice)}
	}
}

// FALSE is counted as zero, 
// TRUE as 1 in comparisons (this will cause different behaviour; TODO: Warnings
//
// Concatenation of comparisons:
// As evaluation is from left to right, y value has to be returned

func EvalComp(op token.Token, x *VSEXP, y *VSEXP) *VSEXP {
	// false and true are not really the same. false is rather the base level.
	if x == nil || y == nil {
		return nil
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
		panic("?Comp: " + op.String())
	}
}

func EvalOp(op token.Token, x *VSEXP, y *VSEXP) *VSEXP {
	if x == nil || y == nil {
		return nil
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
		panic("?Op: " + op.String())
	}
}
