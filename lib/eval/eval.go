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

func EvalStmt(ev *Evaluator, s interface{}) SEXPREC {
	TRACE := ev.trace
	switch s.(type) {
	case *ast.AssignStmt:
		if TRACE {
			print("assignStmt: ")
		}
		e := s.(*ast.AssignStmt)
		var identifier string
		var result SEXPREC
		if e.Tok == token.RIGHTASSIGNMENT {
			identifier = GetIdent(ev, e.Rhs)
			if TRACE {
				println(identifier + " " + e.Tok.String() + " ")
			}
			result = EvalExpr(ev, e.Lhs)
		} else {
			identifier = GetIdent(ev, e.Lhs)
			if TRACE {
				println(identifier + " " + e.Tok.String() + " ")
			}
			result = EvalExpr(ev, e.Rhs)
		}
		ev.topFrame.Insert(identifier, &result)
		return SEXPREC{Kind:  token.ILLEGAL}
	case *ast.ExprStmt:
		if TRACE {
			println("exprStmt")
		}
		e := s.(*ast.ExprStmt)
		return EvalExpr(ev, e.X)
	case *ast.EmptyStmt:
		if TRACE {
			println("emptyStmt")
		}
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
		println("? Stmt")
	}
	return SEXPREC{Kind:  token.ILLEGAL}
}

func PrintResult(ev *Evaluator,r *SEXPREC) {
	TRACE := ev.trace
	switch r.Kind {
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
		print(")")
	default:
		println("SEXPREC with unknown TOKEN")
	}
}

func GetIdent(ev *Evaluator, ex ast.Expr) string {
	node := ex.(*ast.BasicLit)
	return node.Value
}

func EvalExpr(ev *Evaluator, ex ast.Expr) SEXPREC {
	TRACE := ev.trace
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
		println("? Expr")
	}
	node := ex.(*ast.BadExpr)
	return SEXPREC{ValuePos: node.From, Kind: token.FLOAT, Value: math.NaN()}
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

// https://cran.r-project.org/doc/manuals/R-lang.html#Argument-matching
// 1.) Exact matching on tags
// 2.) Partial matching on tags
// 3.) Positional matching

func EvalCall(ev *Evaluator, node *ast.CallExpr) (r SEXPREC) {
	TRACE := ev.trace
	funcobject := node.Fun
	funcname := funcobject.(*ast.BasicLit).Value
	if TRACE {
		println("CallExpr " + funcname)
	}
	f := ev.topFrame.Lookup(funcname)
	if f == nil {
		println("\nError: could not find function \"" + funcname + "\"")
		return SEXPREC{Kind:  token.ILLEGAL}
	} else {
		argNames := make(map[int]string, 3)

		// collect field names
		for n, field := range f.Fieldlist {
			identifier := field.Type.(*ast.Ident)
			argNames[n] = identifier.Name
		}

		argnum := len(argNames)
		taggedArgs := make(map[string]ast.Expr, argnum)
		untaggedArgs := make(map[int]ast.Expr, argnum)
		collectedArgs := make(map[int]*ast.Expr, argnum)
		evaluatedArgs := make(map[int]*SEXPREC, argnum)

		// collect tagged and untagged arguments (unevaluated)
		i := 0
		for n := 0; n < len(node.Args); n++ {
			arg := node.Args[n]
			switch arg.(type) {
			case *ast.TaggedExpr:
				a := arg.(*ast.TaggedExpr)
				taggedArgs[a.Tag] = a.Rhs
			default:
				untaggedArgs[i] = arg
				i = i + 1
			}
		}

		// match tagged arguments
		for n, v := range argNames { // order of n not fix
			expr := taggedArgs[v]
			if expr != nil {
				collectedArgs[n] = &expr
				delete(taggedArgs, v)
			}
		}
		
		// check unused tagged arguments
		if len(taggedArgs) > 0 {
			print("unused argument")
			if len(taggedArgs) > 1 {
				print("s")
			}
			print(" (")
			start:=true
			for k,_ := range taggedArgs{
				if !start {
					print(", ")
				}
				print(k)
				start=false
			}
			print(")\n")
			return SEXPREC{Kind:  token.ILLEGAL}
		}

		// match positional arguments
		j := 0
		for n := 0; n < argnum; n++ {
			if collectedArgs[n] == nil {
				expr := untaggedArgs[j]
				collectedArgs[n] = &expr // TODO check length
				j = j + 1
			}
		}
		
		// check unused positional arguments
		if len(untaggedArgs) > j { // CONT
			
			print("unused argument")
			if (len(untaggedArgs) - j > 1) {
				print("s")
			}
			print(" (")
			start:=true
			// TODO: some caching
			for n := len(argNames)+1 ; n < len(argNames) + len(untaggedArgs) +1 ; n++ {
				if !start {
					print(", ")
				}
				print(n)
				start=false
			}
			print(")\n")
			return SEXPREC{Kind:  token.ILLEGAL}
		}
		
		
		// eval args
		if TRACE {
			println("Eval args " + funcname)
		}
		for n, v := range collectedArgs {
			val := EvalExpr(ev, *v)
			evaluatedArgs[n] = &val
		}

		ev.openFrame()
		{
			if TRACE {
				println("Apply function " + funcname)
			}

			for n, v := range argNames {
				value := evaluatedArgs[n]
				ev.topFrame.Insert(v, value)
			}
			r = EvalStmt(ev, f.Body)
		}
		ev.closeFrame()
	}
	return
}
