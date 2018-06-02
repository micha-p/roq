// https://cran.r-project.org/doc/manuals/R-lang.html#Argument-matching
// 1.) Exact matching on tags
// 2.) Partial matching on tags
// 3.) Positional matching

package eval

import (
	"roq/lib/ast"
	"roq/lib/token"
	"fmt"
	"strconv"
	"strings"
)

// function definition -> formal arguments
// function call -> actual arguments

func tryPartialMatch(partial string, argNames []string, collectedArgs []ast.Expr, DEBUG bool) (int,int) {
	if DEBUG {
		println("Search for partial match: " + partial)
	}
	i := 0
	fieldindex := 0
	for n, name := range argNames {
		if strings.Contains("^"+name, "^"+partial) {
			if DEBUG {
				println("    found: ",name)
			}
			if collectedArgs[n] == nil { 
				fieldindex = n
				i += 1
			}
		}
	}
	return i,fieldindex 
}

func arityOK(funcname string, arity int, node *ast.CallExpr) bool {
	if len(node.Args) == arity {
		return true
	} else {
		fmt.Printf("%d arguments passed to '%s' which requires %d\n", len(node.Args), funcname, arity)
		return false
	}
}

// TODO use results field of funcType
func EvalCallBuiltin(ev *Evaluator, node *ast.CallExpr, funcname string) (r SEXPItf) {
	TRACE := ev.Trace
	DEBUG := ev.Debug
	if TRACE {
		println("Search for builtin: " + funcname)
	}
	switch funcname {
	case "print":
		if arityOK(funcname, 1, node) {
			value := EvalExpr(ev, node.Args[0])
			PrintResult(ev, value)
			return nil
		} else {
			return &ESEXP{Kind: token.ILLEGAL}
		}
	case "list":
		return EvalList(ev, node)
	case "pairlist":
		return EvalPairlist(ev, node)
	case "cat":
		return EvalCat(ev, node)
	// TODO eval arg
	case "length":
		if arityOK(funcname, 1, node) {
			return EvalLength(ev, node)
		} else {
			return &ESEXP{Kind: token.ILLEGAL}
		}
	case "dimnames":
		if arityOK(funcname, 1, node) {
			object := EvalExpr(ev, node.Args[0])
			r := object.Dimnames()
			return r
		} else {
			return &ESEXP{Kind: token.ILLEGAL}
		}
	case "dim":
		if arityOK(funcname, 1, node) {
			object := EvalExpr(ev, node.Args[0])
			r := new(ISEXP)
			r.DimSet(object.Dim())
			r.Test = 1
			return r
		} else {
			return &ESEXP{Kind: token.ILLEGAL}
		}
	case "typeof":
		return EvalTypeof(ev, node)
	case "class":
		return EvalClass(ev, node)
	case "remove":
		for _, arg := range node.Args {
			ev.topFrame.Delete(arg.(*ast.BasicLit).Value, DEBUG)
		}
	case "quit":
		ev.state = eofState
		return &ESEXP{Kind: token.EOF}
	default:
		fmt.Printf("Error in %s(): could not find function \"%s\"\n",funcname,funcname)
		return &ESEXP{Kind: token.ILLEGAL}
	}
	return
}

func EvalCallEllipsisFunction(ev *Evaluator, node *ast.CallExpr, funcname string, f *VSEXP) (r SEXPItf) {
	TRACE := ev.Trace
	DEBUG := ev.Debug

	// collect field names of formals in function definition
	argNames := make([]string, 0, len(f.Fieldlist)+len(node.Args))
	for _, field := range f.Fieldlist {
		switch field.Type.(type) {
		case *ast.Ident:
			identifier := field.Type.(*ast.Ident)
			argNames = append(argNames, identifier.Name)
		case *ast.Ellipsis:
			argNames = append(argNames, "...")
		}
	}

	// collect tagged arguments (unevaluated)
	taggedArgs := make(map[string]int, len(node.Args))
	for n, arg := range node.Args {
		switch arg.(type) {
		case *ast.TaggedExpr:
			a := arg.(*ast.TaggedExpr)
			taggedArgs[a.Tag] = n + 1 // one above default zero value
			if DEBUG {
				println("parameter collected:", a.Tag, n)
			}
		}
	}

	// this map uses the same index as argNames (instead of using a structure)
	// and is filled with correctly identified args during this procedure
	collectedArgs := make([]ast.Expr, len(argNames))

	// this slice covers the actual arguments of the call
	usedArgs := make([]bool, len(node.Args))

	// match tagged arguments
	for fieldindex, fieldname := range argNames {
		if DEBUG {
			println("searching parameter:", fieldname)
		}
		callindex := taggedArgs[fieldname]
		if callindex != 0 { // missing index return default zero value
			collectedArgs[fieldindex] = node.Args[callindex-1].(*ast.TaggedExpr).Rhs
			usedArgs[callindex-1] = true
			if DEBUG {
				println("tagged parameter found:", fieldname, fieldindex, callindex-1)
			}
			delete(taggedArgs, fieldname)
		}
	}

	// find partially matching tags
	for fieldname, callindex := range taggedArgs {
		matches, fieldindex := tryPartialMatch(fieldname, argNames, collectedArgs, DEBUG)
		if matches > 1 {
			fmt.Printf("argument %s matches multiple formal arguments\n", fieldname)
			return &ESEXP{Kind: token.ILLEGAL}
		} else if matches == 1 {
			if TRACE {
				println("argument", fieldname, "matches one formal argument:", argNames[fieldindex])
			}
			collectedArgs[fieldindex] = node.Args[callindex-1]
			usedArgs[callindex-1] = true
			delete(taggedArgs, fieldname)
		}
	}

	// match positional arguments up to ellipsis
	j := 0
	for n, fieldname := range argNames {
		if fieldname == "..." {
			break
		}
		if collectedArgs[n] == nil {
			for usedArgs[j] == true {
				j++
			}
			expr := node.Args[j]
			if DEBUG {
				println("collecting positional argument:", n, j, fieldname)
			}
			collectedArgs[n] = expr
			usedArgs[j] = true
		}
	}

	// collect unused arguments
	j = 1
	for n, isUsed := range usedArgs {
		if isUsed != true {
			fieldname := ".." + strconv.Itoa(j)
			if DEBUG {
				println("appending parameter:", fieldname)
			}
			arg := node.Args[n]
			collectedArgs = append(collectedArgs, arg)
			argNames = append(argNames, fieldname)
			j++
		}
	}
	return EvalApply(ev, funcname, f, argNames, collectedArgs)
}

func EvalCallFunction(ev *Evaluator, node *ast.CallExpr, funcname string, f *VSEXP) (r SEXPItf) {
	DEBUG := ev.Debug
	TRACE := ev.Trace

	// collect field names
	argNames := make([]string, len(f.Fieldlist), len(f.Fieldlist)+len(node.Args))
	for i, field := range f.Fieldlist {
		identifier := field.Type.(*ast.Ident)
		argNames[i] = identifier.Name
	}

	// this map uses the same index as argNames (instead of using a structure)
	// and is filled with correctly identified args during this procedure
	collectedArgs := make([]ast.Expr, len(argNames))

	// collect tagged and untagged arguments (unevaluated)
	taggedArgs := make(map[string]ast.Expr, len(argNames))
	untaggedArgs := make(map[int]ast.Expr, len(argNames))
	i := 0
	for _, arg := range node.Args {
		switch arg.(type) {
		case *ast.TaggedExpr:
			a := arg.(*ast.TaggedExpr)
			taggedArgs[a.Tag] = a.Rhs
		default:
			untaggedArgs[i] = arg
			i++
		}
	}

	// match tagged arguments
	for n, v := range argNames { // order of n not fix
		expr := taggedArgs[v]
		if expr != nil {
			collectedArgs[n] = expr
			delete(taggedArgs, v)
		}
	}

	// find partially matching tags
	for k, v := range taggedArgs {
		matches, fieldindex := tryPartialMatch(k, argNames, collectedArgs, DEBUG)
		if matches > 1 {
			fmt.Printf("Error in %s() : ",funcname)
			fmt.Printf("argument %s matches multiple formal arguments\n", k)
			return &ESEXP{Kind: token.ILLEGAL}
		} else if matches == 1 {
			if TRACE {
				println("argument", k, "matches one formal argument:", argNames[fieldindex])
			}
			collectedArgs[fieldindex] = v
			delete(taggedArgs, k)
		}
	}

	// check unused named arguments // TODO double check
	if len(taggedArgs) > 0 {
		fmt.Printf("Error in %s() : ",funcname)
		fmt.Printf("unused argument")
		if len(taggedArgs) > 1 {
			fmt.Printf("s")
		}
		fmt.Printf(" (")
		start := true
		for k, _ := range taggedArgs {
			if !start {
				fmt.Printf(", ")
			}
			fmt.Printf("%s =", k) // TODO: should ast.expressions carry their input string?
			start = false
		}
		fmt.Printf(")\n")
		return &ESEXP{Kind: token.ILLEGAL}
	}

	// match positional arguments
	j := 0
	for n, _ := range argNames {
		if collectedArgs[n] == nil {
			expr := untaggedArgs[j]
			collectedArgs[n] = expr
			j++
		}
	}

	// check unused positional arguments
	if len(untaggedArgs) > j { // CONT
		fmt.Printf("Error in %s() : ",funcname)
		fmt.Printf("unused argument")
		if len(untaggedArgs)-j > 1 {
			fmt.Printf("s")
		}
		fmt.Printf(" (")
		start := true
		// TODO: some caching
		for n := len(argNames) + 1; n < len(argNames)+len(untaggedArgs)+1; n++ {
			if !start {
				fmt.Printf(", ")
			}
			fmt.Printf("pos %d",n)
			start = false
		}
		fmt.Printf(")\n")
		return &ESEXP{Kind: token.ILLEGAL}
	}
	return EvalApply(ev, funcname, f, argNames, collectedArgs)
}

func EvalCall(ev *Evaluator, node *ast.CallExpr) (r SEXPItf) {
	TRACE := ev.Trace
	funcobject := node.Fun
	funcname := funcobject.(*ast.BasicLit).Value
	if TRACE {
		println("CallExpr: " + funcname)
	}
	if funcname == "c" {
		if TRACE {
			println("Call to special: " + funcname)
		}
		return EvalColumn(ev, node)
	}
	f := ev.topFrame.Lookup(funcname)
	if f == nil {
		return EvalCallBuiltin(ev, node, funcname)
	} else {
		if f.(*VSEXP).ellipsis {
			return EvalCallEllipsisFunction(ev, node, funcname, f.(*VSEXP))
		} else {
			return EvalCallFunction(ev, node, funcname, f.(*VSEXP))
		}
	}
}

func EvalApply(ev *Evaluator, funcname string, f *VSEXP, argNames []string, collectedArgs []ast.Expr) (r SEXPItf) {
	TRACE := ev.Trace
	DEBUG := ev.Debug

	evaluatedArgs := make([]SEXPItf, len(collectedArgs))

	// eval args
	if TRACE {
		println("Eval args for function \"" + funcname + "\"")
	}
	for n, v := range collectedArgs {
		if v != nil {
			val := EvalExprOrAssignment(ev, v)
			evaluatedArgs[int(n)] = val
		}
	}

	ev.openFrame()
	defer ev.closeFrame()

	if (TRACE || DEBUG) {
		println("Apply function \"" + funcname + "\" to call:")
	}
	for n, fieldname := range argNames {
		if (TRACE || DEBUG) {
			print("\targ[",n, "]\t", fieldname)
		}
		if fieldname != "..." {
			value := evaluatedArgs[n]
			if value == nil {
				defaultExpr := f.Fieldlist[n].Default
				if defaultExpr == nil { 
					fmt.Printf("Error in %s() : ", funcname)
					fmt.Printf("argument \"%s\" is missing, with no default\n", fieldname)
					return nil
				} else {
					if DEBUG {
						print("\tDEFAULT")
					}
					value = EvalExpr(ev, defaultExpr)
				}
			} 
			if (TRACE || DEBUG) {
				print("\t")
				PrintResult(ev, value)
			}
			ev.topFrame.Insert(fieldname, value)
		}
	}
	if (TRACE || DEBUG) {
		println("Eval body of function \"" + funcname + "\":")
	}
	r=EvalStmt(ev, f.Body)
	if (TRACE || DEBUG) {
		print("Return from function \"" + funcname + "\" with: ")
		PrintResult(ev,r)
	}
	return r
}
