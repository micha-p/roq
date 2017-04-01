// https://cran.r-project.org/doc/manuals/R-lang.html#Argument-matching
// 1.) Exact matching on tags
// 2.) Partial matching on tags
// 3.) Positional matching

package eval

import (
	"lib/ast"
	"lib/token"
	"strings"
)

// argindex is running along the expected arguments taken from function definition
// callindex is running along the actual arguments given by the call
type argindex int
type callindex int

func tryPartialMatch(partial string, argNames map[argindex]string, bound map[argindex]bool) map[int]argindex {
	//	println("  try to match:",partial)
	matches := make(map[int]argindex, 16)
	i := 0
	for n, name := range argNames {
		if strings.Contains(name, partial) {
			//			print("    found: ",name)
			if bound[n] {
				//				println(" (bound)")
			} else {
				//				println("match:",n)
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
func EvalCall(ev *Evaluator, node *ast.CallExpr) (r SEXPItf) {
	TRACE := ev.Trace
	funcobject := node.Fun
	funcname := funcobject.(*ast.BasicLit).Value
	if TRACE {
		println("CallExpr " + funcname)
	}
	f := ev.topFrame.Lookup(funcname)
	if f == nil {
		switch funcname {
		case "c":
			return EvalColumn(ev, node)
		case "list":
			return EvalList(ev, node)
		case "cat":
			return EvalCat(ev, node)
		// TODO eval arg
		case "length":
			if arityOK(funcname, 1, node) {
				return EvalLength(ev, node)
			} else {
				return &VSEXP{Kind: token.ILLEGAL}
			}
		case "dimnames":
			if arityOK(funcname, 1, node) {
				object := EvalExpr(ev, node.Args[0])
				r := object.Dimnames()
				return r
			} else {
				return &VSEXP{Kind: token.ILLEGAL}
			}
		case "dim":
			if arityOK(funcname, 1, node) {
				object := EvalExpr(ev, node.Args[0])
				r := new(ISEXP)
				r.DimSet(object.Dim())
				r.Test=1
				return r
			} else {
				return &VSEXP{Kind: token.ILLEGAL}
			}
		case "class":
			if arityOK(funcname, 1, node) {
				object := EvalExpr(ev, node.Args[0])
				s := object.Class()
				if s==nil{
					var r string
					switch object.(type){
						case *VSEXP:
							r="numeric"
						case *ISEXP:
							r="numeric"
//						case *LSEXP:
//							r="logical"
						case *TSEXP:
							r="character"
						case *RSEXP:  // TODO pairlist
							r="list"
						case *NSEXP:
							r="NULL"
						default:
							panic("unknown type")
					}
					return &TSEXP{String: r}
				}
				return &TSEXP{String: *s}
			} else {
				return &VSEXP{Kind: token.ILLEGAL}
			}
		default:
			println("\nError: could not find function \"" + funcname + "\"")
			return &VSEXP{Kind: token.ILLEGAL}
		}
	} else {
		argNames := make(map[argindex]string)

		// collect field names
		for n, field := range f.(*VSEXP).Fieldlist {
			i := argindex(n)
			identifier := field.Type.(*ast.Ident)
			argNames[i] = identifier.Name
		}

		argnum := len(argNames)
		// these maps use the same index as argNames (instead of using a structure)
		// might be downgarded to arrays
		boundArgs := make(map[argindex]bool, argnum)
		collectedArgs := make(map[argindex]ast.Expr, argnum) // ast.Expr contains pointers to ast.nodes
		evaluatedArgs := make(map[argindex]SEXPItf, argnum)

		// collect tagged and untagged arguments (unevaluated)
		taggedArgs := make(map[string]ast.Expr, argnum)
		untaggedArgs := make(map[int]ast.Expr, argnum)
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
				boundArgs[n] = true
				collectedArgs[n] = expr
				delete(taggedArgs, v)
			}
		}

		// find partially matching tags
		for k, v := range taggedArgs {
			matchList := tryPartialMatch(k, argNames, boundArgs)
			if len(matchList) == 1 {
				argindex := matchList[0]
				if TRACE {
					println("argument", k, "matches one formal argument:", argNames[argindex])
				}
				collectedArgs[argindex] = v
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
			return &VSEXP{Kind: token.ILLEGAL}
		}

		// match positional arguments
		j := 0
		for n := argindex(0); n < argindex(argnum); n++ {
			if collectedArgs[n] == nil {
				expr := untaggedArgs[j]
				collectedArgs[n] = expr // TODO check length
				j = j + 1
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
			return &VSEXP{Kind: token.ILLEGAL}
		}

		// eval args
		if TRACE {
			println("Eval args " + funcname)
		}
		for n, v := range collectedArgs { // TODO: strictly left to right
			val := EvalExpr(ev, v)
			evaluatedArgs[n] = val
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
			r = EvalStmt(ev, f.(*VSEXP).Body)
		}
		ev.closeFrame()
	}
	return
}
