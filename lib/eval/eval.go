package eval

import (
	"fmt"
	"lib/go/ast"
	"lib/go/parser"
	"lib/go/token"
	"math"
	"strconv"
)

// -> scope.go
// ast.Scopes are used as frames:
// insert method inserts link struct Object with value set.
// however, data field must be set before insertion

type Frame struct {
	Outer   *Frame
	Objects map[string]*ast.Object
}

// NewFrame creates a new scope nested in the outer scope.
func NewFrame(outer *Frame) *Frame {
	const n = 4 // initial frame capacity
	return &Frame{outer, make(map[string]*ast.Object, n)}
}

// Lookup returns the object with the given name if it is
// found in frame s, otherwise it returns nil. Outer frames
// are ignored. TODO!!!
//
func (s *Frame) Lookup(name string) *ast.Object {
	return s.Objects[name]
}

// Insert attempts to insert a named object obj into the frame s.
// If the frame already contains an object alt with the same name, this object is overwritten
func (s *Frame) Insert(obj *ast.Object) (alt *ast.Object) {
	s.Objects[obj.Name] = obj
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

func EvalStmt(ev *Evaluator, s interface{}) {
	TRACE := ev.trace
	switch s.(type) {
	case *ast.AssignStmt:
		if TRACE {
			println("assignStmt")
		}
		e := s.(*ast.AssignStmt)
		value := EvalExpr(ev, e.Rhs)
		identifier := EvalIdent(ev, e.Lhs)
		fmt.Printf("%s <- %g", identifier, value)
		obj := ast.Object{Name: identifier, Data: value}
		ev.topFrame.Insert(&obj)
	case *ast.ExprStmt:
		if TRACE {
			println("exprStmt")
		}
		e := s.(*ast.ExprStmt)
		fmt.Printf("%g", EvalExpr(ev, e.X)) // R has small e for exponential format
	case *ast.EmptyStmt:
		println("")
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
	default:
		println("? Stmt")
	}
}

func EvalIdent(ev *Evaluator, ex ast.Expr) string {
	node := ex.(*ast.BasicLit)
	return node.Value
}

func EvalExpr(ev *Evaluator, ex ast.Expr) float64 {
	TRACE := ev.trace
	switch ex.(type) {
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
			return v
		case token.FLOAT:
			v, err := strconv.ParseFloat(node.Value, 64) // TODO: support for all R formatted values
			if err != nil {
				print("ERROR:")
				println(err)
			}
			if TRACE {
				println(v)
			}
			return v
		case token.IDENT:
			obj := ev.topFrame.Lookup(node.Value)
			v := obj.Data.(float64)
			if TRACE {
				println(v)
			}
			return v
		default:
			println("Unknown node.kind")
		}
	case *ast.BinaryExpr:
		node := ex.(*ast.BinaryExpr)
		if TRACE {
			println("BinaryExpr " + " " + node.Op.String())
		}
		return EvalOp(node.Op, EvalExpr(ev, node.X), EvalExpr(ev, node.Y))
	case *ast.ParenExpr:
		node := ex.(*ast.ParenExpr)
		if TRACE {
			println("ParenExpr")
		}
		return EvalExpr(ev, node.X)
	default:
		println("? Expr")
	}
	return math.NaN()
}

func EvalOp(op token.Token, x float64, y float64) float64 {
	switch op {
	case token.PLUS:
		return x + y
	case token.MINUS:
		return x - y
	case token.MULTIPLICATION:
		return x * y
	case token.DIVISION:
		return x / y
	case token.EXPONENTIATION:
		return math.Pow(x, y)
	case token.MODULUS:
		return math.Mod(x, y)
	default:
		println("? Op: " + op.String())
	}
	return math.NaN()
}
