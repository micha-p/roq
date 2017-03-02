// interfaces:
// ast.Expr
// ast.Stmt

// unevaluated expressions: ast.EXPR
// evaluated expressions and values bound: SEXP

package eval

import (
	"fmt"
	"lib/ast"
	"lib/parser"
	"lib/token"
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
	normalState LoopState = iota
	loopState
	breakState
	nextState
)

// derived of type parser
type Evaluator struct {

	// Tracing/debugging
	trace     bool
	debug     bool
	indent    int // indentation

	invisible bool
	state     LoopState

	// frame
	topFrame *Frame // top-most frame; may be pkgFrame
}

func (e *Evaluator) openFrame() {
	e.topFrame = NewFrame(e.topFrame)
}

func (e *Evaluator) closeFrame() {
	e.topFrame = e.topFrame.Outer
}



func trace(e *Evaluator, args ...interface{}) *Evaluator {
	if e.trace {
		i := 2 * e.indent
		for i > 0 {
			print(" ")
			i--
		}
		fmt.Print(args...)
		print("\n")
		e.indent++
	}
	return e
}

func traceff(e *Evaluator, args ...interface{}) *Evaluator {
	if e.trace {
		e.indent--
		trace(e,args...)
		e.indent++
	}
	return e
}

// Usage pattern: defer un(trace(p, "..."))
func un(e *Evaluator) {
	if e.trace {
		e.indent--
	}
}

func EvalInit(fset *token.FileSet, filename string, src interface{}, mode parser.Mode, traceflag bool, debugflag bool) (r *Evaluator, err error) {

	if fset == nil {
		panic("eval.evalInit: no token.FileSet provided (fset == nil)")
	}

	e := Evaluator{trace: traceflag, debug: debugflag, indent: 0, topFrame: nil}
	e.topFrame = NewFrame(e.topFrame)
	return &e, err
}

// https://cran.r-project.org/doc/manuals/R-lang.html#if

// If value1 is a logical vector with first element TRUE then statement2 is evaluated.
// If the first element of value1 is FALSE then statement3 is evaluated.
// If value1 is a numeric vector then statement3 is evaluated when the first element of value1 is zero and otherwise statement2 is evaluated.
// Only the first element of value1 is used. All other elements are ignored.
// If value1 has any type other than a logical or a numeric vector an error is signalled.

func isTrue(e *SEXP) bool {
	if e == nil {
		return false
	}
	if e.Kind == token.TRUE {
		return true
	}
	if e.Kind == token.FLOAT && e.Value != 0 {
		return true
	}
//  TODO: better documentation on zero=true/false 
	if e.Kind == token.FLOAT && e.Value == 0 {
		return true
	}
	return false
}

// evaluator
// https://go-book.appspot.com/interfaces.html
// an empty interface accepts all pointers

func EvalLoop(ev *Evaluator, e *ast.BlockStmt, cond ast.Expr) *SEXP {
		defer un(trace(ev, "LoopBody"))
		var evloop Evaluator
		evloop = *ev
		evloop.state = loopState
		var rstate LoopState
		for (cond==nil || isTrue(EvalExpr(&evloop, cond))){
			evloop.state=loopState
			for n := 0; n<len(e.List); n++ {
				EvalStmt(&evloop, e.List[n])
				rstate = evloop.state
				if rstate == nextState {break}
			}
			if rstate == nextState {continue}
			if rstate == breakState {break}
		}
		ev.invisible = true
		return &SEXP{Kind: token.NULL}
}

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
		defer un(trace(ev, "ifStmt"))
		e := s.(*ast.IfStmt)
		if isTrue(EvalExpr(ev, e.Cond)) {
			return EvalStmt(ev, e.Body)
		} else if e.Else != nil {
			return EvalStmt(ev, e.Else)
		}
	case *ast.WhileStmt:
		defer un(trace(ev, "whileStmt"))
		e := s.(*ast.WhileStmt)
		return EvalLoop(ev, e.Body, e.Cond)
	case *ast.RepeatStmt:
		defer un(trace(ev, "repeatStmt"))
		e := s.(*ast.RepeatStmt)
		return EvalLoop(ev, e.Body, nil)
	case *ast.ForStmt:
		defer un(trace(ev, "forStmt"))
		//		return EvalLoop(ev, s.(*ast.ForStmt))
	case *ast.BreakStmt:
		ev.state=breakState
		return &SEXP{Kind: token.BREAK}
	case *ast.NextStmt:
		ev.state=nextState
		return &SEXP{Kind: token.NEXT}
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
	defer un(trace(ev, "assignment: " + identifier + " <- "))

	result := EvalExpr(ev, ex)
	defer un(trace(ev, result.Value," ", result.Kind.String()))

	ev.topFrame.Insert(identifier, result)
	ev.invisible = true // just for the following print
	return result
}

func EvalAssignment(ev *Evaluator, e *ast.AssignStmt) *SEXP {

//	defer un(trace(ev, "assignStmt"))

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

	DEBUG := ev.debug

	if DEBUG {
		givenType := reflect.TypeOf(r)
		print("print: ", givenType.String(), ": ", r.Kind.String(), ": ")
	}

	if ev.invisible {
		ev.invisible = false
		return
	} else if r == nil {
		println("FALSE")
	} else {
		switch r.Kind {
		case token.SEMICOLON:
			if DEBUG {
				println("Semicolon")
			}
		case token.ILLEGAL:
			if DEBUG {
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
			if DEBUG {
				println("default print")
			}
			println(r.Kind.String())
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
	DEBUG := ev.debug

//	defer un(trace(ev, "EvalExpr"))
	switch ex.(type) {
	case *ast.FuncLit:
		node := ex.(*ast.FuncLit)
		defer un(trace(ev, "FuncLit"))
		return &SEXP{Kind: token.FUNCTION, Fieldlist: node.Type.Params.List, Body: node.Body}
	case *ast.BasicLit:
		ev.invisible = false
		node := ex.(*ast.BasicLit)
		defer un(trace(ev, "BasicLit ", node.Kind.String()))
		switch node.Kind {
		case token.INT:
			v, err := strconv.ParseFloat(node.Value, 64)
			if err != nil {
				print("ERROR:")
				println(err)
			}
			defer un(traceff(ev, v))
			return &SEXP{ValuePos: node.ValuePos, Kind: token.FLOAT, Value: v}
		case token.FLOAT:
			v, err := strconv.ParseFloat(node.Value, 64) // TODO: support for all R formatted values
			if err != nil {
				print("ERROR:")
				println(err)
			}
			defer un(traceff(ev, v))
			return &SEXP{ValuePos: node.ValuePos, Kind: node.Kind, Value: v}
		case token.STRING:
			return &SEXP{ValuePos: node.ValuePos, Kind: node.Kind, String: node.Value}
		case token.NULL, token.NA, token.INF, token.NAN, token.TRUE, token.FALSE:
			return &SEXP{ValuePos: node.ValuePos, Kind: node.Kind}
		case token.IDENT:
			sexprec := ev.topFrame.Recursive(node.Value)
			if sexprec == nil {
				print("error: object '", node.Value, "' not found\n")
				return &SEXP{ValuePos: node.ValuePos, Kind: token.ILLEGAL, Value: math.NaN()}
			} else {
				defer un(traceff(ev, fmt.Sprintf("%g", sexprec.Value)))
				return sexprec
			}
		default:
			println("Unknown node.kind")
		}
	case *ast.BinaryExpr:
		ev.invisible = false
		return evalBinary(ev,ex.(*ast.BinaryExpr))
	case *ast.CallExpr:
		ev.invisible = false
		return EvalCall(ev, ex.(*ast.CallExpr))
	case *ast.ParenExpr:
		ev.invisible = false
		node := ex.(*ast.ParenExpr)
		if DEBUG {
			println("ParenExpr")
		}
		return EvalExpr(ev, node.X)
	default:
		ev.invisible = false
		givenType := reflect.TypeOf(ex)
		println("?Expr:", givenType.String())
	}
	return &SEXP{Kind: token.ILLEGAL}
}

func evalBinary(ev *Evaluator,node *ast.BinaryExpr) *SEXP {
		defer un(trace(ev, "BinaryExpr"))
		x := EvalExpr(ev, node.X)
		un(traceff(ev, node.Op.String()))
		switch node.Op {
		case token.AND, token.ANDVECTOR:
			if isTrue(x) {
				y := EvalExpr(ev, node.Y)
				if isTrue(y) {
					return y
				} else {
					return nil
				}
			} else {
				return nil
			}
		case token.OR, token.ORVECTOR:
			if isTrue(x) {
				return x
			} else {
				y := EvalExpr(ev, node.Y)
				if isTrue(y) {
					return y
				} else {
					return nil
				}
			}
		case token.LESS, token.LESSEQUAL, token.GREATER, token.GREATEREQUAL, token.EQUAL, token.UNEQUAL:
			return EvalComp(node.Op, x, EvalExpr(ev, node.Y))
		default:
			return EvalOp(node.Op, x, EvalExpr(ev, node.Y))
		}
}

func EvalComp(op token.Token, x *SEXP, y *SEXP) *SEXP {
	if x.Kind == token.ILLEGAL || y.Kind == token.ILLEGAL {
		return &SEXP{Kind: token.ILLEGAL}
	}
	switch op {
	case token.LESS:
		if x.Value < y.Value {
			return x
		} else {
			return nil
		}
	case token.LESSEQUAL:
		if x.Value <= y.Value {
			return x
		} else {
			return nil
		}
	case token.GREATER:
		if x.Value > y.Value {
			return x
		} else {
			return nil
		}
	case token.GREATEREQUAL:
		if x.Value >= y.Value {
			return x
		} else {
			return nil
		}
	case token.EQUAL:
		if x.Value == y.Value {
			return x
		} else {
			return nil
		}
	case token.UNEQUAL:
		if x.Value != y.Value {
			return x
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
		println("?Op: " + op.String())
		return &SEXP{Kind: token.ILLEGAL}
	}
	return &SEXP{Kind: token.FLOAT, Value: val}
}
