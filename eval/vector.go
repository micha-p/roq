package eval

import (
	"roq/calc"
	"roq/lib/token"
)

func EvalVectorOp(x *VSEXP, y *VSEXP, FUN func(float64, float64) float64) *VSEXP {
	if x.Slice==nil && y.Slice==nil {
		return &VSEXP{Immediate: FUN(x.Immediate,y.Immediate)}
	} else if x.Slice==nil {
		return &VSEXP{Slice: calc.MapIA(FUN,x.Immediate,y.Slice)}
	} else if y.Slice==nil {
		return &VSEXP{Slice: calc.MapAI(FUN,x.Slice,y.Immediate)}
	} else {
		return &VSEXP{Slice: calc.MapAA(FUN,x.Slice,y.Slice)}
	}
}

// FALSE is counted as zero, 
// TRUE as 1 in comparisons (this will cause different behaviour; TODO: Warnings
//
// Concatenation of comparisons:
// As evaluation is from left to right, y value has to be returned

func EvalComp(op token.Token, x *VSEXP, y *VSEXP) *VSEXP {
	if x == nil || y == nil {
		return nil
	}
	switch op {
	case token.EQUAL:
		return EvalVectorOp(x,y,calc.FEQUAL)
	case token.UNEQUAL:
		return EvalVectorOp(x,y,calc.FUNEQUAL)
	case token.LESS:
		return EvalVectorOp(x,y,calc.FLESS)
	case token.LESSEQUAL:
		return EvalVectorOp(x,y,calc.FLESSEQUAL)
	case token.GREATER:
		return EvalVectorOp(y,x,calc.FLESS)
	case token.GREATEREQUAL:
		return EvalVectorOp(y,x,calc.FLESSEQUAL)
	default:
		panic("?Vcomp: " + op.String())
	}
}

func EvalOp(op token.Token, x *VSEXP, y *VSEXP) *VSEXP {
	if x == nil || y == nil {
		return nil
	}
	switch op {
	case token.PLUS:
		return EvalVectorOp(x,y,calc.FPLUS)
	case token.MINUS:
		return EvalVectorOp(x,y,calc.FMINUS)
	case token.MULTIPLICATION:
		return EvalVectorOp(x,y,calc.FMULTIPLICATION)
	case token.DIVISION:
		return EvalVectorOp(x,y,calc.FDIVISION)
	case token.EXPONENTIATION:
		return EvalVectorOp(x,y,calc.FEXPONENTIATION)
	case token.MODULUS:
		return EvalVectorOp(x,y,calc.FMODULUS)
	default:
		panic("?Op: " + op.String())
	}
}
