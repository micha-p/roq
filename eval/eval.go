// interfaces:
// ast.Expr
// ast.Stmt

// unevaluated expressions: ast.EXPR
// evaluated expressions and values bound: SEXP

package eval

import (
	"lib/ast"
	"lib/parser"
	"lib/token"
	"math"
	"fmt"
	"reflect"
	"strconv"
)

// Frames are derived from ast.Scopes:
// insert method inserts link struct Object with value set.
// however, data field must be set before insertion

type Frame struct {
	Outer   *Frame
	Objects map[string] SEXPItf
}

// NewFrame creates a new scope nested in the outer scope.
func NewFrame(outer *Frame) *Frame {
	const n = 4 // initial frame capacity
	return &Frame{outer, make(map[string]SEXPItf, n)}
}

// Lookup returns the object with the given name if it is
// found in frame s, otherwise it returns nil. Outer frames
// are ignored. TODO!!!
//
func (f *Frame) Lookup(name string) SEXPItf {
	return f.Objects[name]
}

func (f *Frame) Recursive(name string) (r SEXPItf) {
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
func (s *Frame) Insert(identifier string, obj SEXPItf) (alt SEXPItf) {
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
	Trace  bool
	Debug  bool
	indent int // indentation

	Invisible bool
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
	if e.Trace {
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
	if e.Trace {
		e.indent--
		trace(e, args...)
		e.indent++
	}
	return e
}

// Usage pattern: defer un(trace(p, "..."))
func un(e *Evaluator) {
	if e.Trace {
		e.indent--
	}
}

func EvalInit(fset *token.FileSet, filename string, src interface{}, mode parser.Mode, traceflag bool, debugflag bool) (r *Evaluator, err error) {

	if fset == nil {
		panic("eval.evalInit: no token.FileSet provided (fset == nil)")
	}

	e := Evaluator{Trace: traceflag, Debug: debugflag, indent: 0, topFrame: nil}
	e.topFrame = NewFrame(e.topFrame)
	return &e, err
}

// https://cran.r-project.org/doc/manuals/R-lang.html#if

// If value1 is a logical vector with first element TRUE then statement2 is evaluated.
// If the first element of value1 is FALSE then statement3 is evaluated.
// If value1 is a numeric vector then statement3 is evaluated when the first element of value1 is zero and otherwise statement2 is evaluated.
// Only the first element of value1 is used. All other elements are ignored.
// If value1 has any type other than a logical or a numeric vector an error is signalled.

func isTrue(e SEXPItf) bool {
	if e == nil {
		return false
	}
	if e.(*VSEXP).Kind == token.TRUE {
		return true
	}
	if e.(*VSEXP).Kind == token.FLOAT && e.(*VSEXP).Immediate != 0 {
		return true
	}
	if e.(*VSEXP).Kind == token.FLOAT && e.(*VSEXP).Slice != nil && e.Length()>0 {
		if e.(*VSEXP).Slice[0] == 0 {
			return false
		} else {
			return true
		}
	}

	//  THIS IS A MAIN DIFFERENCE
	//  TODO: better documentation on zero=true/false
	if e.(*VSEXP).Kind == token.FLOAT && e.(*VSEXP).Immediate == 0 {
		return true
	}
	return false
}

// evaluator
// https://go-book.appspot.com/interfaces.html
// an empty interface accepts all pointers

func EvalLoop(ev *Evaluator, e *ast.BlockStmt, cond ast.Expr) SEXPItf {
	defer un(trace(ev, "LoopBody"))
	var evloop Evaluator
	evloop = *ev
	evloop.state = loopState
	var rstate LoopState
	for cond == nil || isTrue(EvalExpr(&evloop, cond)) {
		evloop.state = loopState
		for n := 0; n < len(e.List); n++ {
			EvalStmt(&evloop, e.List[n])
			rstate = evloop.state
			if rstate == nextState {
				break
			}
		}
		if rstate == nextState {
			continue
		}
		if rstate == breakState {
			break
		}
	}
	ev.Invisible = true
	return &NSEXP{}
}
func EvalFor(ev *Evaluator, e *ast.BlockStmt, identifier string, iterable SEXPItf) SEXPItf {
	if iterable.(*VSEXP).Slice == nil {
		panic("Vector expected")
	}
	defer un(trace(ev, "LoopBody"))
	var evloop Evaluator
	evloop = *ev
	evloop.state = loopState
	var rstate LoopState
	for _,v := range iterable.(*VSEXP).Slice {
		evloop.state = loopState
		// TODO: make use of cached position in map
		ev.topFrame.Insert(identifier, &VSEXP{Kind: token.FLOAT, Immediate: v})
		for n := 0; n < len(e.List); n++ {
			EvalStmt(&evloop, e.List[n])
			rstate = evloop.state
			if rstate == nextState {
				break
			}
		}
		if rstate == nextState {
			continue
		}
		if rstate == breakState {
			break
		}
	}
	ev.Invisible = true
	return &NSEXP{}
}

func EvalStmt(ev *Evaluator, s ast.Stmt) (r SEXPItf) {
	TRACE := ev.Trace
	DEBUG := ev.Debug
	switch s.(type) {
	case *ast.AssignStmt:
		return EvalAssignment(ev, s.(*ast.AssignStmt))
	case *ast.ExprStmt:
		e := s.(*ast.ExprStmt)
		return EvalExprOrAssignment(ev, e.X)
	case *ast.EmptyStmt:
		if DEBUG {
			println("emptyStmt")
		}
		return nil
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
		e := s.(*ast.ForStmt)
		return EvalFor(ev, e.Body, e.Parameter.String(), EvalExpr(ev, e.Iterable))
	case *ast.BreakStmt:
		ev.state = breakState
		return &VSEXP{Kind: token.BREAK}
	case *ast.NextStmt:
		ev.state = nextState
		return &VSEXP{Kind: token.NEXT}
	case *ast.BlockStmt:
		if TRACE {
			println("blockStmt")
		}
		e := s.(*ast.BlockStmt)
		for _, stmt := range e.List {
			switch stmt.(type) {
			case *ast.EmptyStmt:
			default:
				r = EvalStmt(ev, stmt)
			}
		}
		if TRACE {
			println("return: "/*, r.Kind.String()*/) // TODO typoeofstring
		}
		return r
	case *ast.VersionStmt:
		ev.state = nextState
		return &VSEXP{Kind: token.VERSION}
	default:
		givenType := reflect.TypeOf(s)
		println("?Stmt:", givenType.String())
	}
	return &VSEXP{Kind: token.ILLEGAL}
}


func doAssignment(ev *Evaluator,lhs ast.Expr, rhs ast.Expr) SEXPItf {
	var value SEXPItf
	switch lhs.(type){
	case *ast.CallExpr:
		doAttributeReplacement(ev,lhs.(*ast.CallExpr),rhs)
	case *ast.BasicLit:
		target := getIdent(ev, lhs)
		defer un(trace(ev, "assignment: "+target+" <- "))
		value = EvalExpr(ev, rhs)
		defer un(trace(ev, value.(*VSEXP).Immediate, " ", value.(*VSEXP).Kind.String()))
		ev.topFrame.Insert(target, value)
	}
	ev.Invisible = true // just for the following print
	return value
}

// Assignments might be Expressions or Stmts, the first return a SEXP during evaluation,
// the latter an invisible object

func EvalAssignment(ev *Evaluator, e *ast.AssignStmt) SEXPItf {

	//	defer un(trace(ev, "assignStmt"))

	var target,value ast.Expr
	if e.Tok == token.RIGHTASSIGNMENT {
		target = e.Rhs
		value = e.Lhs
	} else {
		target = e.Lhs
		value = e.Rhs
	}
	return doAssignment(ev, target, value)
}


func EvalExprOrAssignment(ev *Evaluator, ex ast.Expr) SEXPItf {
	TRACE := ev.Trace
	if TRACE {
		println("Expr or assignment:")
	}
	switch ex.(type) {
	case *ast.BinaryExpr:
		node := ex.(*ast.BinaryExpr)
		switch node.Op {
		case token.SHORTASSIGNMENT:
			return doAssignment(ev, node.X, node.Y)
		case token.LEFTASSIGNMENT:
			return doAssignment(ev, node.X, node.Y)
		case token.RIGHTASSIGNMENT:
			return doAssignment(ev, node.X, node.Y)
		}
	}
	return EvalExpr(ev, ex)
}

func EvalExpr(ev *Evaluator, ex ast.Expr) SEXPItf {
	DEBUG := ev.Debug

	//	defer un(trace(ev, "EvalExpr"))
	switch ex.(type) {
	case *ast.FuncLit:
		node := ex.(*ast.FuncLit)
		defer un(trace(ev, "FuncLit"))
		return &VSEXP{Kind: token.FUNCTION, Fieldlist: node.Type.Params.List, Body: node.Body}
	case *ast.BasicLit:
		ev.Invisible = false
		node := ex.(*ast.BasicLit)
		defer un(trace(ev, "BasicLit ", node.Kind.String()))
		switch node.Kind {
		case token.FLOAT:
			vfloat, err := strconv.ParseFloat(node.Value, 64) // TODO: support for all R formatted values
			if err != nil {
				print("ERROR:")
				println(err)
			}
			return &VSEXP{ValuePos: node.ValuePos,Kind: node.Kind, Immediate: vfloat}
		case token.INT:
			vint, err := strconv.Atoi(node.Value)
			if err != nil {
				print("ERROR:")
				println(err)
			}
			var vfloat float64
			vfloat, err = strconv.ParseFloat(node.Value, 64) // TODO: support for all R formatted values
			if err != nil {
				print("ERROR:")
				println(err)
			}
			return &ISEXP{ValuePos: node.ValuePos,Integer: vint, Immediate: vfloat}
		case token.STRING:
			return &TSEXP{ValuePos: node.ValuePos,String: node.Value}
		case token.NULL, token.NA, token.FALSE:
			// TODO jus return nil?
			return &NSEXP{ValuePos: node.ValuePos}
		case token.INF:
			return &VSEXP{ValuePos: node.ValuePos, Immediate: math.Inf(+1)}
		case token.NAN:
			return &VSEXP{ValuePos: node.ValuePos, Immediate: math.NaN()}
		case token.IDENT:
			sexprec := ev.topFrame.Recursive(node.Value)
			if sexprec == nil {
				print("error: object '", node.Value, "' not found\n")
				return &VSEXP{ValuePos: node.ValuePos,Kind: token.ILLEGAL}
			} else {
				return sexprec
			}
		default:
			println("Unknown node.Kind")
		}
	case *ast.BinaryExpr:
		ev.Invisible = false
		return evalBinary(ev, ex.(*ast.BinaryExpr))
	case *ast.CallExpr:
		ev.Invisible = false
		return EvalCall(ev, ex.(*ast.CallExpr))
	case *ast.IndexExpr:
		ev.Invisible = false
		return EvalIndexExpr(ev, ex.(*ast.IndexExpr))
	case *ast.ParenExpr:
		ev.Invisible = false
		node := ex.(*ast.ParenExpr)
		if DEBUG {
			println("ParenExpr")
		}
		return EvalExpr(ev, node.X)
	default:
		ev.Invisible = false
		givenType := reflect.TypeOf(ex)
		println("?Expr:", givenType.String())
	}
	return &VSEXP{Kind: token.ILLEGAL}
}

func evalBinary(ev *Evaluator, node *ast.BinaryExpr) SEXPItf {
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
	case token.SEQUENCE:
		// TODO: same for INDEXDOMAIN
		low := EvalExpr(ev, node.X)
		high := EvalExpr(ev, node.Y)
		slice := make([]float64,1+high.IntegerGet()-low.IntegerGet())
		start := low.FloatGet()
		for n,_ := range slice {
			slice[n]=start
			start=start+1
		}
		return &VSEXP{Kind: token.FLOAT, Slice: slice}
	case token.LESS, token.LESSEQUAL, token.GREATER, token.GREATEREQUAL, token.EQUAL, token.UNEQUAL:
		return EvalComp(node.Op, x.(*VSEXP), EvalExpr(ev, node.Y).(*VSEXP))
	default:
		return EvalOp(node.Op, x.(*VSEXP), EvalExpr(ev, node.Y).(*VSEXP))
	}
}

