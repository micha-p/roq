// interfaces:
// ast.Expr
// ast.Stmt

// unevaluated expressions: ast.EXPR
// evaluated expressions and values bound: SEXP

package eval

import (
	"fmt"
	"roq/lib/ast"
	"roq/lib/parser"
	"roq/lib/token"
	"math"
	"reflect"
	"strconv"
)


type LoopState int

const (
	normalState LoopState = iota
	loopState
	breakState
	nextState
	eofState
)

// derived of type parser
type Evaluator struct {

	// Tracing/debugging
	Trace  bool
	Debug  bool
	Major  string
	Minor  string
	indent int // indentation

	Invisible bool
	state     LoopState

	// frame
	topFrame *Frame // top-most frame; may be pkgFrame
	globalFrame *Frame // top-most frame; may be pkgFrame
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

func EvalInit(fset *token.FileSet, filename string, src interface{}, mode parser.Mode, traceflag bool, debugflag bool, MAJOR string, MINOR string) (r *Evaluator, err error) {

	if fset == nil {
		panic("roq/eval.evalInit: no token.FileSet provided (fset == nil)")
	}

	e := Evaluator{Trace: traceflag, Debug: debugflag, indent: 0, topFrame: nil, Major: MAJOR, Minor: MINOR}
	e.topFrame = NewFrame(e.topFrame)
	e.globalFrame = e.topFrame
	return &e, err
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
			if s.(*ast.EmptyStmt).Implicit{
				// println("emptyStmt (implicit)") // too many messages
			} else {
				println("emptyStmt")
			}
		}
		return nil
	case *ast.IfStmt:
		defer un(trace(ev, "ifStmt"))
		e := s.(*ast.IfStmt)
		testresult := EvalExpr(ev, e.Cond)
		if testresult != nil && isTrue(testresult) {
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
		ex := EvalExpr(ev, e.Iterable)
		return EvalFor(ev, e.Body, e.Parameter.String(), ex)
	case *ast.BreakStmt:
		ev.state = breakState
		return &ESEXP{Kind: token.BREAK}
	case *ast.NextStmt:
		ev.state = nextState
		return &ESEXP{Kind: token.NEXT}
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
			print("blockStmt return: ")
			PrintResult(ev,r)
		}
		return r
	case *ast.VersionStmt:
		ev.state = nextState
		return &ESEXP{Kind: token.VERSION}
	case *ast.EOFStmt:
		ev.state = eofState
		return &ESEXP{Kind: token.EOF}
	default:
		givenType := reflect.TypeOf(s)
		println("?Stmt:", givenType.String())
	}
	return &ESEXP{Kind: token.ILLEGAL}
}

func doAssignment(ev *Evaluator, lhs ast.Expr, rhs ast.Expr) SEXPItf {
	var value SEXPItf
	switch lhs.(type) {
	case *ast.CallExpr:
		doAttributeReplacement(ev, lhs.(*ast.CallExpr), rhs)
	case *ast.BasicLit:
		target := getIdent(ev, lhs)
		defer un(trace(ev, "assignment: "+target+" <- "))
		value = EvalExpr(ev, rhs)
		ev.topFrame.Insert(target, value)
	}
	ev.Invisible = true // just for the following print
	return value
}

func doSuperAssignment(ev *Evaluator, lhs ast.Expr, rhs ast.Expr) SEXPItf {
	var value SEXPItf
	switch lhs.(type) {
	case *ast.CallExpr:
		doAttributeReplacement(ev, lhs.(*ast.CallExpr), rhs)
	case *ast.BasicLit:
		target := getIdent(ev, lhs)
		defer un(trace(ev, "superassignment: "+target+" <<- "))
		value = EvalExpr(ev, rhs)
		ev.globalFrame.Insert(target, value)
	}
	ev.Invisible = true // just for the following print
	return value
}

// Assignments might be Expressions or Stmts, the first return a SEXP during evaluation,
// the latter an invisible object

func EvalAssignment(ev *Evaluator, e *ast.AssignStmt) SEXPItf {

	//	defer un(trace(ev, "assignStmt"))

	switch e.Tok {
	case token.LEFTASSIGNMENT:
		return doAssignment(ev, e.Lhs, e.Rhs)
	case token.RIGHTASSIGNMENT:
		return doAssignment(ev, e.Rhs, e.Lhs)
	case token.SUPERLEFTASSIGNMENT:
		return doSuperAssignment(ev, e.Lhs, e.Rhs)
	case token.SUPERRIGHTASSIGNMENT:
		return doSuperAssignment(ev, e.Rhs, e.Lhs)
	default:
		panic("panic during assignment")
	}
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
		case token.SUPERLEFTASSIGNMENT:
			return doSuperAssignment(ev, node.X, node.Y)
		case token.SUPERRIGHTASSIGNMENT:
			return doSuperAssignment(ev, node.Y, node.X)
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
		
		withEllipsis := false
		for _, field := range node.Type.Params.List {
			switch field.Type.(type){
			case *ast.Ellipsis:
				withEllipsis=true
				break
			}
		}
		return &VSEXP{Fieldlist: node.Type.Params.List, Body: node.Body, ellipsis: withEllipsis}
	case *ast.BasicLit:
		ev.Invisible = false
		node := ex.(*ast.BasicLit)
		defer un(trace(ev, "BasicLit ", node.Kind.String()))
		switch node.Kind {
		case token.FLOAT:
			vfloat, err := strconv.ParseFloat(node.Value, 64) 		// TODO: support for all R formatted values
			if err != nil {
				print("ERROR:")
				println(err)
			}
			return &VSEXP{ValuePos: node.ValuePos, Immediate: vfloat}
		case token.INT: 											// in value domain, all numbers should be double float
			vfloat, err := strconv.ParseFloat(node.Value, 64) 		// TODO: support for all R formatted values
			if err != nil {
				print("ERROR:")
				println(err)
			}
			vint, err := strconv.Atoi(node.Value) 					// TODO: support for all R formatted values
			if err != nil {
				print("ERROR:")
				println(err)
			}
			return &ISEXP{ValuePos: node.ValuePos, Immediate: vfloat, Integer: vint}
		case token.STRING:
			return &TSEXP{ValuePos: node.ValuePos, String: node.Value}
		case token.TRUE:
			return &VSEXP{ValuePos: node.ValuePos, Immediate: 1}   	// in R: TRUE+1 = 2
		case token.NULL, token.FALSE:								// TODO just return nil?
			return &NSEXP{ValuePos: node.ValuePos}
		case token.INF:
			return &VSEXP{ValuePos: node.ValuePos, Immediate: math.Inf(+1)}
		case token.NAN, token.NA:
			return &VSEXP{ValuePos: node.ValuePos, Immediate: math.NaN()}
		case token.IDENT:
			if DEBUG {
				print("Retrieving ident: ", node.Value," = ")
			}
			return ev.topFrame.Recursive(node.Value)
		default:
			panic("Unknown basic literal:"+node.Kind.String())
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
	return &ESEXP{Kind: token.ILLEGAL}
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
		slice := make([]float64, 1+high.IntegerGet()-low.IntegerGet())
		start := low.FloatGet()
		for n, _ := range slice {
			slice[n] = start
			start = start + 1
		}
		return &VSEXP{Slice: slice}
	case token.LESS, token.LESSEQUAL, token.GREATER, token.GREATEREQUAL, token.EQUAL, token.UNEQUAL:
		y := EvalExpr(ev, node.Y)
		if x == nil || y == nil {
			return nil
		} else {
			return EvalComp(node.Op, x.(*VSEXP), y.(*VSEXP))
		}
	default:
		y := EvalExpr(ev, node.Y)
		if x == nil || y == nil {
			return nil
		} else {
			return EvalOp(node.Op, x.(*VSEXP), y.(*VSEXP))
		}
	}
}
