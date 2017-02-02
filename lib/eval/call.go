// https://cran.r-project.org/doc/manuals/R-lang.html#Argument-matching
// 1.) Exact matching on tags
// 2.) Partial matching on tags
// 3.) Positional matching

package eval

import (
	"fmt"
	"lib/go/ast"
	"lib/go/parser"
	"lib/go/token"
	"math"
	"strconv"
)

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
