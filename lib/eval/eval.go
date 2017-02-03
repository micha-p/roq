// interfaces:
// ast.Expr
// ast.Stmt


// unevaluated expressions: ast.EXPR 
// evaluated expressions and values bound: SEXPREC

package eval

import (
	"fmt"
	"lib/go/ast"
	"lib/go/parser"
	"lib/go/token"
	"math"
	"strconv"
	"reflect"
)

// Frames are derived from ast.Scopes:
// insert method inserts link struct Object with value set.
// however, data field must be set before insertion


type Frame struct {
	Outer   *Frame
	Objects map[string]*SEXPREC
}

// NewFrame creates a new scope nested in the outer scope.
func NewFrame(outer *Frame) *Frame {
	const n = 4 // initial frame capacity
	return &Frame{outer, make(map[string]*SEXPREC, n)}
}

// Lookup returns the object with the given name if it is
// found in frame s, otherwise it returns nil. Outer frames
// are ignored. TODO!!!
//
func (f *Frame) Lookup(name string) *SEXPREC {
	return f.Objects[name]
}

func (f *Frame) Recursive(name string) (r *SEXPREC) {
	r = f.Objects[name]
	if r == nil {
		if f.Outer != nil {
			return f.Outer.Recursive(name)
		} 
	}
	return 
}


func getIdent(ev *Evaluator, ex ast.Expr) string {
	node := ex.(*ast.BasicLit)
	return node.Value
}

// Insert attempts to insert a named object obj into the frame s.
// If the frame already contains an object alt with the same name, this object is overwritten
func (s *Frame) Insert(identifier string, obj *SEXPREC) (alt *SEXPREC) {
	s.Objects[identifier] = obj
	return
}

// derived of type parser
type Evaluator struct {

	// Tracing/debugging
	trace  bool // == (mode & Trace != 0)
	indent int  // indentation used for tracing output

	// frame
	topFrame *Frame // top-most frame; may be pkgFrame
}

func (e *Evaluator) openFrame() {
	e.topFrame = NewFrame(e.topFrame)
}

func (e *Evaluator) closeFrame() {
	e.topFrame = e.topFrame.Outer
}

func EvalInit(fset *token.FileSet, filename string, src interface{}, mode parser.Mode, traceflag bool) (r *Evaluator, err error) {

	if fset == nil {
		panic("eval.evalInit: no token.FileSet provided (fset == nil)")
	}

	e := Evaluator{traceflag, 0, nil}
	e.topFrame = NewFrame(e.topFrame)
	return &e, err
}

// evaluator
// https://go-book.appspot.com/interfaces.html
// an empty interface accepts all pointers

func EvalStmt(ev *Evaluator, s ast.Stmt) SEXPREC {
	TRACE := ev.trace
	DEBUG := false
	switch s.(type) {
	case *ast.AssignStmt:
		return EvalAssignment(ev, s.(*ast.AssignStmt))
	case *ast.ExprStmt:
		e := s.(*ast.ExprStmt)
		return EvalExprOrShortAssign(ev, e.X)
	case *ast.EmptyStmt:
		if DEBUG {
			println("emptyStmt")
		}
		return SEXPREC{Kind:  token.INVISIBLE}
	case *ast.IfStmt:
		if TRACE {
			println("ifStmt")
		}
	case *ast.ForStmt:
		if TRACE {
			println("forStmt")
		}
	case *ast.BlockStmt:
		if TRACE {
			println("blockStmt")
		}
		e := s.(*ast.BlockStmt)
		var r SEXPREC
		for _, stmt := range e.List {
			r = EvalStmt(ev, stmt)
		}
		return r 
	default:
		givenType := reflect.TypeOf(s)
		println("?Stmt:",givenType.String())
	}
	return SEXPREC{Kind:  token.ILLEGAL}
}

func doAssignment(ev *Evaluator,identifier string, ex ast.Expr) SEXPREC {
	TRACE := ev.trace

	if TRACE {
		print("assignment: ",identifier," := ")
	}
	result := EvalExpr(ev, ex)
	if TRACE {
		println(result.Kind.String())
	}
	ev.topFrame.Insert(identifier, &result)
	return result
}

func EvalAssignment(ev *Evaluator,e *ast.AssignStmt) SEXPREC {
	TRACE := ev.trace

		if TRACE {
			println("assignStmt: ")
		}

		var identifier string
		if e.Tok == token.RIGHTASSIGNMENT {
			identifier = getIdent(ev, e.Rhs)
		} else {
			identifier = getIdent(ev, e.Lhs)
		}

		var nodepointer ast.Expr
		if e.Tok == token.RIGHTASSIGNMENT {
			nodepointer = e.Lhs
		} else {
			nodepointer = e.Rhs
		}

		return doAssignment(ev,identifier, nodepointer)
}

func PrintResult(ev *Evaluator,r *SEXPREC) {
	TRACE := ev.trace
	switch r.Kind {
	case token.INVISIBLE:
	case token.ILLEGAL:
		if TRACE {
			println("ILLEGAL RESULT")
		}
	case token.FLOAT:
		fmt.Printf("%g\n", r.Value) // R has small e for exponential format
	case token.FUNCTION:
		print("function(")
		for n, field := range r.Fieldlist {
			//for _,ident := range field.Names {
			//	print(ident)
			//}
			identifier := field.Type.(*ast.Ident)
			if n > 0 {
				print(",")
			}
			print(identifier.Name)
		}
		println(")")
	default:
		println("SEXPREC with unknown TOKEN")
	}
}


func EvalExprOrShortAssign(ev *Evaluator, ex ast.Expr) SEXPREC {
	TRACE := ev.trace
	if TRACE {
		println("Expr or short assignment:")
	}
	switch ex.(type) {
	case *ast.BinaryExpr:
		node := ex.(*ast.BinaryExpr)
		if node.Op==token.SHORTASSIGNMENT {
			return doAssignment(ev,getIdent(ev,node.X),node.Y)
		}
	}
	return EvalExpr(ev,ex)
}

func EvalExpr(ev *Evaluator, ex ast.Expr) SEXPREC {
	TRACE := ev.trace
	if TRACE {
		print("EvalExpr ")
	}
	switch ex.(type) {
	case *ast.FuncLit:
		node := ex.(*ast.FuncLit)
		if TRACE {
			print("FuncLit")
		}
		return SEXPREC{Kind: token.FUNCTION, Fieldlist: node.Type.Params.List, Body: node.Body}
	case *ast.BasicLit:
		node := ex.(*ast.BasicLit)
		if TRACE {
			print("BasicLit " + " " + node.Value + " (" + node.Kind.String() + "): ")
		}
		switch node.Kind {
		case token.INT:
			v, err := strconv.ParseFloat(node.Value, 64)
			if err != nil {
				print("ERROR:")
				println(err)
			}
			if TRACE {
				println(v)
			}
			return SEXPREC{ValuePos: node.ValuePos, Kind: token.FLOAT, Value: v}
		case token.FLOAT:
			v, err := strconv.ParseFloat(node.Value, 64) // TODO: support for all R formatted values
			if err != nil {
				print("ERROR:")
				println(err)
			}
			if TRACE {
				println(v)
			}
			return SEXPREC{ValuePos: node.ValuePos, Kind: token.FLOAT, Value: v}
		case token.IDENT:
			sexprec := ev.topFrame.Recursive(node.Value)
			if sexprec == nil {
				print("error: object '",node.Value,"' not found\n")
				return SEXPREC{ValuePos: node.ValuePos, Kind: token.ILLEGAL, Value: math.NaN()}
			} else {
				if TRACE {
					fmt.Printf("%g\n", sexprec.Value)
				}
				return *sexprec
			}
		default:
			println("Unknown node.kind")
		}
	case *ast.BinaryExpr:
		node := ex.(*ast.BinaryExpr)
		if TRACE {
			println("BinaryExpr " + " " + node.Op.String())
		}
		return EvalOp(node.Op, EvalExpr(ev, node.X), EvalExpr(ev, node.Y))
	case *ast.CallExpr:
		return EvalCall(ev, ex.(*ast.CallExpr))
	case *ast.ParenExpr:
		node := ex.(*ast.ParenExpr)
		if TRACE {
			println("ParenExpr")
		}
		return EvalExpr(ev, node.X)
	default:
		givenType := reflect.TypeOf(ex)
		println("?Expr:",givenType.String())
	}
	return SEXPREC{Kind: token.ILLEGAL}
}


func EvalOp(op token.Token, x SEXPREC, y SEXPREC) SEXPREC {
	if (x.Kind == token.ILLEGAL || y.Kind == token.ILLEGAL) {
		return SEXPREC{Kind:  token.ILLEGAL}
	}
	var val float64
	switch op {
	case token.PLUS:
		val = x.Value + y.Value
	case token.MINUS:
		val = x.Value - y.Value
	case token.MULTIPLICATION:
		val = x.Value * y.Value
	case token.DIVISION:
		val = x.Value / y.Value
	case token.EXPONENTIATION:
		val = math.Pow(x.Value, y.Value)
	case token.MODULUS:
		val = math.Mod(x.Value, y.Value)
	default:
		println("? Op: " + op.String())
		return SEXPREC{Kind:  token.ILLEGAL}
	}
    return SEXPREC{Kind:  token.FLOAT, Value: val}
}


