// https://cran.r-project.org/doc/manuals/R-lang.html#Argument-matching
// 1.) Exact matching on tags
// 2.) Partial matching on tags
// 3.) Positional matching

package eval

import (
	"lib/ast"
	"lib/token"
	"strings"
	"strconv"
)

// function definition -> formal arguments
// function call -> actual arguments

type formalindex int

func tryPartialMatch(partial string, argNames map[formalindex]string, collectedArgs map[formalindex]ast.Expr) map[int]formalindex {
	//  println("  try to match:",partial)
	matches := make(map[int]formalindex, len(argNames))
	i := 0
	for n, name := range argNames {
		if strings.Contains(name, partial) {
			//          print("    found: ",name)
			if collectedArgs[n] == nil {
				matches[i] = n
				i += 1
			}
		}
	}
	return matches
}

func arityOK(funcname string, arity int, node *ast.CallExpr) bool {
	if len(node.Args) == arity {
		return true
	} else {
		print(len(node.Args), " arguments passed to '", funcname, "' which requires ", arity, "\n")
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
		for _,arg := range(node.Args){
			ev.topFrame.Delete(arg.(*ast.BasicLit).Value,DEBUG)
		}
	case "quit":
		panic("quit")
	default:
		println("Error: could not find function \"" + funcname + "\"")
		return &ESEXP{Kind: token.ILLEGAL}
	}
	return
}

func EvalApply(ev *Evaluator, funcname string, f *VSEXP, argNames map[formalindex]string, collectedArgs map[formalindex]ast.Expr) (r SEXPItf) {
	TRACE := ev.Trace
	DEBUG := ev.Debug

	evaluatedArgs := make(map[formalindex]SEXPItf, len(collectedArgs))

	// eval args
	if TRACE {
		println("Eval args " + funcname)
	}
	for n, v := range collectedArgs { // TODO: strictly left to right
		val := EvalExpr(ev, v)
		evaluatedArgs[formalindex(n)] = val
	}


	ev.openFrame()
	{
		if TRACE {
			println("Apply function " + funcname)
		}

		if DEBUG {println("apply function",funcname, "to call:")}
		for n, fieldname := range argNames {
			value := evaluatedArgs[formalindex(n)]
			if DEBUG {
				print("arg[",n,"]\t",fieldname,"\t")
				PrintResult(ev,value)
			} 
			ev.topFrame.Insert(fieldname, value)
		}
		r = EvalStmt(ev, f.Body)
	}
	ev.closeFrame()
	return
}



func EvalCallEllipsisFunction(ev *Evaluator, node *ast.CallExpr, funcname string, f *VSEXP) (r SEXPItf) {
	TRACE := ev.Trace
	DEBUG := ev.Debug

	argNames := make(map[formalindex]string)

	// collect field names of formals in function definition
	for n, field := range f.Fieldlist {
		i := formalindex(n)
		switch field.Type.(type){
		case *ast.Ident:
			identifier := field.Type.(*ast.Ident)
			argNames[i] = identifier.Name
		case *ast.Ellipsis:
			argNames[i] = "..."
		}
	}
	
	// collect tagged arguments (unevaluated)
	taggedArgs := make(map[string]int, len(node.Args))
	for n, arg := range(node.Args) {
		switch arg.(type) {
		case *ast.TaggedExpr:
			a := arg.(*ast.TaggedExpr)
			taggedArgs[a.Tag] = n + 1  // one above default zero value
			if DEBUG {println("parameter collected:", a.Tag, n)}
		}
	}

	// this map uses the same index as argNames (instead of using a structure)
	collectedArgs := make(map[formalindex]ast.Expr, len(node.Args))

	// this slice covers the actual arguments of the call
	usedArgs := make([]bool, len(node.Args))

	// match tagged arguments
	for fieldindex, fieldname := range argNames { // order of n in map not fix
		if DEBUG {println("searching parameter:", fieldname)}
		callindex := taggedArgs[fieldname]
		if callindex != 0 {  // missing index return default zero value
 			collectedArgs[fieldindex] = node.Args[callindex-1].(*ast.TaggedExpr).Rhs
			usedArgs[callindex-1] = true
			if DEBUG {println("tagged parameter found:", fieldname, callindex-1)}
			delete(taggedArgs, fieldname)
		}
	}

	// find partially matching tags
	for fieldname, callindex := range taggedArgs {
		matchList := tryPartialMatch(fieldname, argNames, collectedArgs)
		if len(matchList) == 1 {
			fieldindex := matchList[0]
			if TRACE {
				println("argument", fieldname, "matches one formal argument:", argNames[fieldindex])
			}
			collectedArgs[fieldindex] = node.Args[callindex-1]
			usedArgs[callindex-1] = true
			delete(taggedArgs, fieldname)
		} else if len(matchList) > 1 {
			println("argument", fieldname, "matches multiple formal arguments")
		}
	}

	// match positional arguments up to ellipsis
	j := 0
	for n,fieldname := range(argNames) {
		if fieldname=="..." {break}
		if DEBUG {println("collecting positional argument:", n, j, fieldname)}
		if collectedArgs[n] == nil {
			for usedArgs[j]==true {j++}
			expr := node.Args[j]
			collectedArgs[n] = expr
			usedArgs[j]=true
		}
	}

	// collect unused arguments
	j = 1
	for n,isUsed := range(usedArgs) {
		if isUsed != true {
			fieldname := ".." + strconv.Itoa(j)
			if DEBUG {println("parameter appended:", fieldname)}
			collectedArgs[formalindex(len(argNames))] = node.Args[n]
			argNames[formalindex(len(argNames))]=fieldname
			j++
		}
	}
	return EvalApply(ev,funcname,f, argNames, collectedArgs)
}


func EvalCallFunction(ev *Evaluator, node *ast.CallExpr, funcname string, f *VSEXP) (r SEXPItf) {
	TRACE := ev.Trace

	argNames := make(map[formalindex]string)

	// collect field names
	for n, field := range f.Fieldlist {
		i := formalindex(n)
		identifier := field.Type.(*ast.Ident)
		argNames[i] = identifier.Name
	}
	
	// this map uses the same index as argNames (instead of using a structure)
	collectedArgs := make(map[formalindex]ast.Expr, len(argNames))

	// collect tagged and untagged arguments (unevaluated)
	taggedArgs := make(map[string]ast.Expr, len(argNames))
	untaggedArgs := make(map[int]ast.Expr, len(argNames))
	i := 0

	for _, arg := range(node.Args) {
		switch arg.(type) {
		case *ast.TaggedExpr:
			a := arg.(*ast.TaggedExpr)
			taggedArgs[a.Tag] = a.Rhs
		default:
			untaggedArgs[i] = arg
			i ++
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
		matchList := tryPartialMatch(k, argNames, collectedArgs)
		if len(matchList) == 1 {
			fieldindex := matchList[0]
			if TRACE {
				println("argument", k, "matches one formal argument:", argNames[fieldindex])
			}
			collectedArgs[fieldindex] = v
			delete(taggedArgs, k)
		} else if len(matchList) > 1 {
			println("argument", k, "matches multiple formal arguments")
		}
	}

	// check unused tagged arguments
	if len(taggedArgs) > 0 {
		print("unused argument")
		if len(taggedArgs) > 1 {
			print("s")
		}
		print(" (")
		start := true
		for k, _ := range taggedArgs {
			if !start {
				print(", ")
			}
			print(k)
			start = false
		}
		print(")\n")
		return &ESEXP{Kind: token.ILLEGAL}
	}

	// match positional arguments
	j := 0
	for n,_ := range(argNames) {
		if collectedArgs[n] == nil {
			expr := untaggedArgs[j]
			collectedArgs[n] = expr
			j++
		}
	}

	// check unused positional arguments
	if len(untaggedArgs) > j { // CONT

		print("unused argument")
		if len(untaggedArgs)-j > 1 {
			print("s")
		}
		print(" (")
		start := true
		// TODO: some caching
		for n := len(argNames) + 1; n < len(argNames)+len(untaggedArgs)+1; n++ {
			if !start {
				print(", ")
			}
			print(n)
			start = false
		}
		print(")\n")
		return &ESEXP{Kind: token.ILLEGAL}
	}
	return EvalApply(ev,funcname,f, argNames, collectedArgs)
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
		if f.(*VSEXP).ellipsis{
			return EvalCallEllipsisFunction(ev, node, funcname, f.(*VSEXP))
		} else {
			return EvalCallFunction(ev, node, funcname, f.(*VSEXP))
		}
	}
}
