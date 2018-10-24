package eval

import (
	"roq/lib/ast"
	"roq/lib/token"
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
	return x.Length()
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
			i := int(a[x.Counter]) // TODO ISEXP
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



func EvalRangeExpressionToIterator(ev *Evaluator, a SEXPItf, b SEXPItf) IteratorItf {
	r := new(RangeIterator)
	r.Start = a.IntegerGet() - 2
	r.Counter = r.Start
	r.End = b.IntegerGet() -2
	return r
}

func IndexValueAsInt(node *ast.BasicLit) int{
	switch node.Kind {
	case token.FLOAT:
		vfloat, err := strconv.ParseFloat(node.Value, 64) // TODO: support for all R formatted values
		if err != nil {
			print("ERROR:")
			println(err)
		}
		return int(math.Floor(vfloat))
	case token.INT:
		vint, err := strconv.Atoi(node.Value)
		if err != nil {
			print("ERROR:")
			println(err)
		}
		return vint
	case token.NULL:
		return 0
	case token.IDENT:
		return 0
	default:
		println("Unknown node.Kind for index")
		return 0
	}
}

func EvalSexpressionToIterator(sexp SEXPItf) IteratorItf {
	switch sexp.(type) {
	case *ISEXP:
		r := new(OnceIterator)
		r.Offset=int(math.Floor(sexp.(*ISEXP).Immediate))
		return r
	case *VSEXP:
		if sexp.(*VSEXP).Slice == nil {
			r := new(OnceIterator)
			r.Offset=int(math.Floor(sexp.(*VSEXP).Immediate))
			return r
		} else {
			r := new(ArrayIterator)
			r.Slice=sexp.(*VSEXP).Slice
			return r
		}
	default:
		givenType := reflect.TypeOf(sexp)
		println("?IndexSExpr:", givenType.String())
		r := new(EmptyIterator)
		return r
	}
}

func EvalIndexExpressionToIterator(ev *Evaluator, ex ast.Expr) IteratorItf {
	defer un(trace(ev, "EvalIndexExpressionToIterator"))
	switch ex.(type) {
	case *ast.Ident:
		sexp:=EvalExpr(ev,ex)
		return EvalSexpressionToIterator(sexp)
	case *ast.BasicLit:
		ev.Invisible = false
		node := ex.(*ast.BasicLit)
		defer un(trace(ev, "BasicLit ", node.Kind.String()))
		index := IndexValueAsInt(node)
		if index == 0 {
			obj := ev.topFrame.Recursive(node.Value)
			if obj == nil {
				print("error: object '", node.Value, "' not found\n")
				return new(EmptyIterator)
			} else {
				r := new(ArrayIterator)
				r.Slice = obj.(*VSEXP).Slice // TODO: check this
				r.Len=r.Length()
				return r
			}
		} else {
			r := new(OnceIterator)
			r.Offset=index
			return r
		}
	case *ast.BinaryExpr:
		ev.Invisible = false
		node := ex.(*ast.BinaryExpr)
		if node.Op == token.SEQUENCE {
			return EvalRangeExpressionToIterator(ev, EvalExpr(ev,node.X).(*VSEXP),EvalExpr(ev,node.Y).(*VSEXP))
		} else {
			sexp:=evalBinary(ev,node)
			return EvalSexpressionToIterator(sexp)
		}
	default:
		ev.Invisible = false
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


// TODO consistant naming for index, value and toplevel domain:
// evalExprI -> ISEXPR
func EvalIndexedArray(ev *Evaluator, node *ast.IndexExpr) SEXPItf {
	array := EvalExpr(ev,node.Array)
	if array == nil {
		panic("array not found\n")
	} else {
		iterator := EvalIndexExpressionToIterator(ev, node.Index)
		r := make([]float64,0,array.Length())
		var n int
		for true {
			n = iterator.Next()
			if n >= 0 {
				element := array.(*VSEXP).Slice[n]
				r = append(r,element)
			} else {
				break
			}
		}  
		return &VSEXP{ValuePos: array.Pos(), Slice:r}
	}
}

func EvalIndexedList(ev *Evaluator, node *ast.ListIndexExpr) SEXPItf {
	list := EvalExpr(ev,node.Array)
	if list == nil {
		panic("list not found\n")
	} else {
		return list.(*RSEXP).Slice[IndexValueAsInt(node.Index.(*ast.BasicLit))-1]
	}
}


