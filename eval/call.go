// call: collect args
// args: evaluate args
// apply: jump into closure, insert into topFrame, jump into body

// https://cran.r-project.org/doc/manuals/R-lang.html#Argument-matching
// 1.) Exact matching on tags
// 2.) Partial matching on tags
// 3.) Positional matching


// function definition -> formal arguments
// function call -> actual arguments

package eval

import (
	"roq/lib/ast"
	"roq/lib/token"
	"fmt"
	"strings"
	"errors"
)

func tryPartialMatch(partial string, argNames []string, collectedArgs []ast.Expr) (int,int) {
	i := 0
	fieldindex := 0
	for n, name := range argNames {
		if strings.Contains("^"+name, "^"+partial) {
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
	DEBUG := ev.Debug
	switch funcname {
	case "print": // TODO arity
		if arityOK(funcname, 1, node) {
			return EvalPrint(ev, node)
		} else {
			return &ESEXP{Kind: token.ILLEGAL}
		}
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
			ev.topFrame.Delete(arg.(*ast.Ident).Name, DEBUG)
		}
	case "options":
		return nil
	case "quit":
		ev.state = eofState
		return &ESEXP{Kind: token.EOF}
	default:
		fmt.Printf("Error in %s(): could not find function \"%s\"\n",funcname,funcname)
		return &ESEXP{Kind: token.ILLEGAL}
	}
	return
}

// collect field names of formals in function definition
func getArgNames(thefunction *VSEXP) (argnames []string) {
	argNames := make([]string, 0, len(thefunction.Fieldlist))
	for _, field := range thefunction.Fieldlist {
		switch field.Type.(type) {
		case *ast.Ident:
			identifier := field.Type.(*ast.Ident)
			argNames = append(argNames, identifier.Name)
		case *ast.Ellipsis:
			argNames = append(argNames, "...")
		}
	}
	return argNames
}

func CollectArgs(ev *Evaluator, node *ast.CallExpr, funcname string, argNames []string) ([]ast.Expr, error) {
	TRACE := ev.Trace

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
	for n, v := range argNames { // TODO order of n not guaranteed
		expr := taggedArgs[v]
		if expr != nil {
			collectedArgs[n] = expr
			delete(taggedArgs, v)
		}
	}

	// find partially matching tags
	for k, v := range taggedArgs {
		matches, fieldindex := tryPartialMatch(k, argNames, collectedArgs)
		if matches > 1 {
			fmt.Printf("Error in %s() : ",funcname)
			fmt.Printf("argument %s matches multiple formal arguments\n", k)
			return nil, errors.New("Golang error: argument matches multiple formal arguments" )
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
		return nil, errors.New("Golang error: Unused named arg" )
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
		return nil, errors.New("Golang error: Unused positional arg" )
	}
	return collectedArgs, nil
}

func PrintAstExpression(ev *Evaluator, n int, arg ast.Expr){
	print("\t",n,"=")
	switch arg.(type) {
	case *ast.BasicLit:
		if arg.(*ast.BasicLit).Kind==token.ELLIPSIS {
			println("\t\t","...")
		} else {
			println(" (literal)\t",arg.(*ast.BasicLit).Value)
		}
	default:
		if arg != nil{
			print(" (evaluated)\t")
			PrintResult(EvalExprMute(ev,arg))
		} else{
			println("\tnil")
		}
	}
}

func PrintListofAstExpressions(ev *Evaluator, arglist []ast.Expr){
	for n,arg := range arglist {
		PrintAstExpression(ev,n+1,arg)
	}
}

func PrintListofSExpressions(valuelist []SEXPItf){
	for n,v := range valuelist {
		if v==nil{
			println("\tnil")
		} else {
			print("\t",n+1,"= ")
			PrintResult(v)
		}
	}
}

func PrintArgNames(namelist []string){
	for n,arg := range namelist {
		print("\t",n+1,"=\t")
		println(arg)
	}
}




func EvalCall(ev *Evaluator, funcname string, node *ast.CallExpr) (r SEXPItf) {
	TRACE := ev.Trace
	DEBUG := ev.Debug
	if funcname == "c" {
		if TRACE {
			println("Call to protected function: " + funcname)
		}
		return EvalColumn(ev, node)
	} else if funcname == "list" {
		if TRACE {
			println("Call to protected function: " + funcname)
		}
		return EvalList(ev, node)
	}
	
	thefunction := ev.topFrame.Recursive(funcname)
	if thefunction == nil {
		if TRACE || DEBUG{
			println("Call to builtin: " + funcname)
		}
		return EvalCallBuiltin(ev, node, funcname)
	} else {
		if thefunction.(*VSEXP).ellipsis {
			return EvalCallwithEllipsis(ev, node, thefunction)
		} else {
			if TRACE {
				println("Call to function: " + funcname)
			}
			argNames := getArgNames(thefunction.(*VSEXP))
			collectedArgs, err := CollectArgs(ev, node, funcname, argNames)
			if err != nil {
				return &ESEXP{Kind: token.ILLEGAL}
			} else {
				evaluatedArgs := EvalArgs(ev, funcname, collectedArgs)
				return EvalApply(ev, funcname, thefunction.(*VSEXP), argNames, evaluatedArgs)
			}
			
		}
	}
}


func EvalArgs(ev *Evaluator, funcname string, collectedArgs []ast.Expr) ([]SEXPItf) {
	DEBUG := ev.Debug
	evaluatedArgs := make([]SEXPItf, len(collectedArgs))

	if DEBUG {
		println("Eval args for function \"" + funcname + "\":")
	}
	for n, v := range collectedArgs {
		if v != nil {
			val := EvalExprOrAssignment(ev, v)
			if DEBUG {
				print("\targ[",n,"] = ")
				PrintResult(val)
			}
			evaluatedArgs[int(n)] = val
		}
	}
	return evaluatedArgs
}
