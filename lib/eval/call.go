// https://cran.r-project.org/doc/manuals/R-lang.html#Argument-matching
// 1.) Exact matching on tags
// 2.) Partial matching on tags
// 3.) Positional matching

package eval

import (
	"strings"
	"lib/go/ast"
	"lib/go/token"
)

func tryPartialMatch(partial string, argNames map[int]string, bound map[int]bool) map[int]int {
//	println("  try to match:",partial)
	matches := make(map[int]int, 16)
	i := 0
	for n,name := range argNames {
		if strings.Contains(name,partial) {
//			print("    found: ",name)
			if bound[n] {
//				println(" (bound)")
			} else {
//				println("match:",n)
				matches[i]=n
				i += 1
			}
		}
	}
	return matches
}
	
	
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
		argNames := make(map[int]string, 0)

		// collect field names
		for n, field := range f.Fieldlist {
			identifier := field.Type.(*ast.Ident)
			argNames[n] = identifier.Name
		}

		argnum := len(argNames)
		// these maps use the same index as argNames (instead of using a structure)
		// might be downgarded to arrays
		boundArgs     := make(map[int]bool, argnum)
		collectedArgs := make(map[int]*ast.Expr, argnum)
		evaluatedArgs := make(map[int]*SEXPREC, argnum)

		// collect tagged and untagged arguments (unevaluated)
		taggedArgs    := make(map[string]ast.Expr, argnum)
		untaggedArgs  := make(map[int]ast.Expr, argnum)
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
		for argindex, v := range argNames { // order of n not fix
			expr := taggedArgs[v]
			if expr != nil {
				boundArgs[argindex] = true
				collectedArgs[argindex] = &expr
				delete(taggedArgs, v)
			}
		}
		
		// find partially matching tags
		for k,v := range taggedArgs {
			matchList := tryPartialMatch(k,argNames, boundArgs)
			if len(matchList)==1 {
				argindex := matchList[0]
				if TRACE {
					println("argument",k,"matches one formal argument:",argNames[argindex])
				}
				collectedArgs[argindex] = &v
				delete(taggedArgs, k)
			} else if len(matchList)>1 {
				println("argument",k,"matches multiple formal arguments")
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
