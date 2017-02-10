// interfaces:
// ast.Expr
// ast.Stmt

// unevaluated expressions: ast.EXPR
// evaluated expressions and values bound: SEXP

package eval

import (
	"fmt"
	"lib/go/ast"
	"lib/go/parser"
	"lib/go/token"
	"math"
	"reflect"
	"strconv"
)

// Frames are derived from ast.Scopes:
// insert method inserts link struct Object with value set.
// however, data field must be set before insertion

type Frame struct {
	Outer   *Frame
	Objects map[string]*SEXP
}

// NewFrame creates a new scope nested in the outer scope.
func NewFrame(outer *Frame) *Frame {
	const n = 4 // initial frame capacity
	return &Frame{outer, make(map[string]*SEXP, n)}
}

// Lookup returns the object with the given name if it is
// found in frame s, otherwise it returns nil. Outer frames
// are ignored. TODO!!!
//
func (f *Frame) Lookup(name string) *SEXP {
	return f.Objects[name]
}

func (f *Frame) Recursive(name string) (r *SEXP) {
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
func (s *Frame) Insert(identifier string, obj *SEXP) (alt *SEXP) {
	s.Objects[identifier] = obj
	return
}

type LoopState int

const (
	normalState LoopState=iota
	breakState
	nextState
)
	

// derived of type parser
type Evaluator struct {

	// Tracing/debugging
	trace  bool
	debug  bool
	invisible bool
	state LoopState
	indent int // indentation used for tracing output

	// frame
	topFrame *Frame // top-most frame; may be pkgFrame
}

func (e *Evaluator) openFrame() {
	e.topFrame = NewFrame(e.topFrame)
}

func (e *Evaluator) closeFrame() {
	e.topFrame = e.topFrame.Outer
}

func EvalInit(fset *token.FileSet, filename string, src interface{}, mode parser.Mode, traceflag bool, debugflag bool) (r *Evaluator, err error) {

	if fset == nil {
		panic("eval.evalInit: no token.FileSet provided (fset == nil)")
	}

	e := Evaluator{trace: traceflag, debug: debugflag, indent: 0, topFrame: nil}
	e.topFrame = NewFrame(e.topFrame)
	return &e, err
}

func isTrue(e *SEXP) bool {
	if e.Kind == token.NULL {
		return false
	}
	if e.Kind == token.FLOAT && e.Value != 0 {
		return true
	}

	return false
}

// evaluator
// https://go-book.appspot.com/interfaces.html
// an empty interface accepts all pointers

func EvalStmt(ev *Evaluator, s ast.Stmt) *SEXP {
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
		return &SEXP{Kind: token.SEMICOLON}
	case *ast.IfStmt:
		if TRACE {
			println("ifStmt")
		}
		e := s.(*ast.IfStmt)
		if isTrue(EvalExpr(ev, e.Cond)) {
			return EvalStmt(ev, e.Body)
		} else if e.Else != nil {
			return EvalStmt(ev, e.Else)
		}
	case *ast.WhileStmt:
		if TRACE {
			println("whileStmt")
		}
//		return EvalLoop(ev, s.(*ast.whileStmt))
	case *ast.RepeatStmt:
		if TRACE {
			println("repeatStmt")
		}
//		return EvalLoop(ev, s.(*ast.whileStmt),TRUE)
	case *ast.ForStmt:
		if TRACE {
			println("forStmt")
		}
//		return EvalLoop(ev, s.(*ast.ForStmt))
	case *ast.BlockStmt:
		if TRACE {
			println("blockStmt")
		}
		e := s.(*ast.BlockStmt)
		var r *SEXP
		for _, stmt := range e.List {
			switch stmt.(type) {
			case *ast.EmptyStmt:
			default:
				r = EvalStmt(ev, stmt)
			}
		}
		if TRACE {
			println("return: ", r.Kind.String())
		}
		return r
	default:
		givenType := reflect.TypeOf(s)
		println("?Stmt:", givenType.String())
	}
	return &SEXP{Kind: token.ILLEGAL}
}

func doAssignment(ev *Evaluator, identifier string, ex ast.Expr) *SEXP {
	TRACE := ev.trace

	if TRACE {
		print("assignment: ", identifier, " <- ")
	}
	result := EvalExpr(ev, ex)
	if TRACE {
		println(result.Kind.String())
	}
	ev.topFrame.Insert(identifier, result)
	ev.invisible=true  // just for the following print
	return result
}

func EvalAssignment(ev *Evaluator, e *ast.AssignStmt) *SEXP {
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
	return doAssignment(ev, identifier, nodepointer)
}


// visibility is stored in the evaluator and unset after every print
func PrintResult(ev *Evaluator, r *SEXP) {
	TRACE := ev.trace
	DEBUG := ev.debug

	if DEBUG {
		givenType := reflect.TypeOf(r)
		print("print: ", givenType.String(), ": ", r.Kind.String(), ": ")
	}

	if ev.invisible {
		ev.invisible=false
		return
	} else {
		switch r.Kind {
		case token.SEMICOLON:
			if DEBUG {
				println("Semicolon")
			}
		case token.ILLEGAL:
			if TRACE {
				println("ILLEGAL RESULT")
			}
		case token.FLOAT:
			fmt.Printf("[1] %g\n", r.Value) // R has small e for exponential format
		case token.FUNCTION:
			if DEBUG {
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
			}
		default:
			println("?SEXP:", r.Kind.String())
		}
	}
}

func EvalExprOrShortAssign(ev *Evaluator, ex ast.Expr) *SEXP {
	TRACE := ev.trace
	if TRACE {
		println("Expr or short assignment:")
	}
	switch ex.(type) {
	case *ast.BinaryExpr:
		node := ex.(*ast.BinaryExpr)
		if node.Op == token.SHORTASSIGNMENT {
			return doAssignment(ev, getIdent(ev, node.X), node.Y)
		}
	}
	return EvalExpr(ev, ex)
}

func EvalExpr(ev *Evaluator, ex ast.Expr) *SEXP {
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
		return &SEXP{Kind: token.FUNCTION, Fieldlist: node.Type.Params.List, Body: node.Body}
	case *ast.BasicLit:
		ev.invisible=false
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
			return &SEXP{ValuePos: node.ValuePos, Kind: token.FLOAT, Value: v}
		case token.FLOAT:
			v, err := strconv.ParseFloat(node.Value, 64) // TODO: support for all R formatted values
			if err != nil {
				print("ERROR:")
				println(err)
			}
			if TRACE {
				println(v)
			}
			return &SEXP{ValuePos: node.ValuePos, Kind: token.FLOAT, Value: v}
		case token.IDENT:
			sexprec := ev.topFrame.Recursive(node.Value)
			if sexprec == nil {
				print("error: object '", node.Value, "' not found\n")
				return &SEXP{ValuePos: node.ValuePos, Kind: token.ILLEGAL, Value: math.NaN()}
			} else {
				if TRACE {
					fmt.Printf("%g\n", sexprec.Value)
				}
				return sexprec
			}
		default:
			println("Unknown node.kind")
		}
	case *ast.BinaryExpr:
		ev.invisible=false
		node := ex.(*ast.BinaryExpr)
		if TRACE {
			println("BinaryExpr " + " " + node.Op.String())
		}
		return EvalOp(node.Op, EvalExpr(ev, node.X), EvalExpr(ev, node.Y))
	case *ast.CallExpr:
		ev.invisible=false
		return EvalCall(ev, ex.(*ast.CallExpr))
	case *ast.ParenExpr:
		ev.invisible=false
		node := ex.(*ast.ParenExpr)
		if TRACE {
			println("ParenExpr")
		}
		return EvalExpr(ev, node.X)
	default:
		ev.invisible=false
		givenType := reflect.TypeOf(ex)
		println("?Expr:", givenType.String())
	}
	return &SEXP{Kind: token.ILLEGAL}
}

func EvalOp(op token.Token, x *SEXP, y *SEXP) *SEXP {
	if x.Kind == token.ILLEGAL || y.Kind == token.ILLEGAL {
		return &SEXP{Kind: token.ILLEGAL}
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
		return &SEXP{Kind: token.ILLEGAL}
	}
	return &SEXP{Kind: token.FLOAT, Value: val}
}
