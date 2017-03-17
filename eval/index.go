package eval

import (
	"lib/ast"
	"lib/token"
	"math"
	"strconv"
	"reflect" // TODO: make obsolete
)

// nodes inside brackets are evaluated within a special domain, the index domain. 
// this domain is based on integers and should be recognized already during scanning.
//
// several iterators are needed for arrays on LHS and RHS
// full array
// sequence
// subset vector
// boolean vector

// iterators return either a positive number, which should be used as offset or -1, indicating the end.


type IteratorItf interface {
	Next() int
	Length()  int
}

type (
	OnceIterator struct {
		Done	bool
		Offset	int
	}
	EmptyIterator struct {
	}
	FullIterator struct {
		Max		int
		Counter	int
	}
	RangeIterator struct{
		Start 	int
		End 	int
		Counter	int
	}
	ArrayIterator struct{
		Slice 	[]float64
		Len 	int  // cached length
		Counter	int
	}
)

func (x *FullIterator) Length() int { 
	return x.Max
}
func (x *RangeIterator) Length() int {
	return 1+(x.End - x.Start)
}
func (x *ArrayIterator) Length() int {
	return len(x.Slice)
}
func (x *OnceIterator) Length() int { 
	return 1
}
func (x *EmptyIterator) Length() int { 
	return 0
}

func (x *FullIterator) Next() int { 
	x.Counter +=1
	if (x.Counter < x.Max){ 
		return x.Counter
	} else {
		return -1
	}
}
func (x *RangeIterator) Next() int {
	if (x.Counter <= x.End){ 
			x.Counter +=1
			return x.Counter
	} else {
		return -1
	}
}
func (x *ArrayIterator) Next() int {
	a := x.Slice
	if (x.Counter < len(a)){ 
			i := int(a[x.Counter])
			x.Counter +=1
			return i-1
	} else {
		return -1
	}
}
func (x *OnceIterator) Next() int { 
	if x.Done { 
		return -1
	} else {
		x.Done = true
		return x.Offset -1
	}
}
func (x *EmptyIterator) Next() int { 
	return -1
}



func IndexDomainEvalRange(ev *Evaluator, a *SEXP, b *SEXP) IteratorItf {
	r := new(RangeIterator)
	r.Start = a.Offset-1
	r.Counter = a.Offset -1
	r.End = b.Offset -1
	return r
}

func IndexDomainEval(ev *Evaluator, ex ast.Expr) IteratorItf {

	defer un(trace(ev, "IndexDomainEval"))
	switch ex.(type) {
	case *ast.BasicLit:
		ev.invisible = false
		node := ex.(*ast.BasicLit)
		defer un(trace(ev, "BasicLit ", node.Kind.String()))
		switch node.Kind {
		case token.FLOAT:
			vfloat, err := strconv.ParseFloat(node.Value, 64) // TODO: support for all R formatted values
			if err != nil {
				print("ERROR:")
				println(err)
			}
			// TODO check conversion to integer
			r := new(OnceIterator)
			r.Offset=int(math.Floor(vfloat))
			return r
		case token.INT:
			vint, err := strconv.Atoi(node.Value)
			if err != nil {
				print("ERROR:")
				println(err)
			}
			r := new(OnceIterator)
			r.Offset=vint
			return r
		case token.NULL:
			r := new(EmptyIterator)
			return r
		case token.IDENT:
			sexprec := ev.topFrame.Recursive(node.Value)
			if sexprec == nil {
				print("error: object '", node.Value, "' not found\n")
				r := new(EmptyIterator)
				return r
			} else {
				r := new(ArrayIterator)
				r.Slice = sexprec.Slice
				r.Len=len(sexprec.Slice)
				return r
			}
		default:
			println("Unknown node.kind")
		}
    case *ast.BinaryExpr:
		ev.invisible = false
		node := ex.(*ast.BinaryExpr)
		if node.Op == token.SEQUENCE {
			return IndexDomainEvalRange(ev, EvalExpr(ev,node.X),EvalExpr(ev,node.Y))
		} else {
//			return sexp2iterator(evalBinary(ev, ex.(*ast.BinaryExpr)))
			r := new(EmptyIterator)
			return r
		}
	/*case *ast.CallExpr:
		ev.invisible = false
		return EvalCall(ev, ex.(*ast.CallExpr))
	case *ast.IndexExpr:
		ev.invisible = false
		return EvalIndexExpr(ev, ex.(*ast.IndexExpr))
	case *ast.ParenExpr:
		ev.invisible = false
		node := ex.(*ast.ParenExpr)
		if DEBUG {
			println("ParenExpr")
		}
		return EvalExpr(ev, node.X) */
	default:
		ev.invisible = false
		givenType := reflect.TypeOf(ex)
		println("?IndexExpr:", givenType.String())
	}
	r := new(EmptyIterator)
	return r
}


// https://cran.r-project.org/doc/manuals/R-lang.html#Indexing-by-vectors
// A special case is the zero index, which has null effects: 
// x[0] is an empty vector and otherwise including zeros among positive or 
// negative indices has the same effect as if they were omitted.

func EvalIndexExpr(ev *Evaluator, node *ast.IndexExpr) *SEXP {
	arrayPart := node.Array.(*ast.BasicLit)
	array := ev.topFrame.Recursive(arrayPart.Value)
	if array == nil {
		print("error: object '", arrayPart.Value, "' not found\n")
		return &SEXP{ValuePos: arrayPart.ValuePos, Kind: token.ILLEGAL, Immediate: math.NaN()}
	} else {
		iterator := IndexDomainEval(ev, node.Index)
		r := make([]float64,0,len(array.Slice))
		var n int
		for true {
			n = iterator.Next()
			if n >= 0 {
				element := array.Slice[n]
				r = append(r,element)
			} else {
				break
			}
		}  
		return &SEXP{ValuePos: arrayPart.ValuePos, Kind: token.FLOAT, Slice:r}
	}
}
