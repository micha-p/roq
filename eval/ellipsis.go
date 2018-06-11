package eval

import (
	"roq/lib/ast"
	"strconv"
)


func EvalArgswithDotDotArguments(ev *Evaluator, funcname string, arglist []ast.Expr)[]SEXPItf{
	DEBUG := ev.Debug
	evaluatedArgs := make([]SEXPItf, 0, len(arglist))
	if DEBUG {
		println("EvalArgswithDotDotArguments")
	}
	for n, arg := range arglist { // TODO: strictly left to right
		if arg != nil {
			var val SEXPItf
			if DEBUG {
				print("\tProcessing: ", n, ":\t")
			}
			switch arg.(type) {
			case *ast.BasicLit:
				val=EvalExprOrAssignment(ev, arg)
				if DEBUG {
					print("appending evaluated argument:\t")
					PrintResult(val)
				}
				evaluatedArgs=append(evaluatedArgs,val)
			case *ast.Ellipsis:
				if DEBUG {
					println(" ELLIPSIS")
					DumpFrames(ev)
				}
				for k:=1; k<=len(ev.topFrame.Objects); k++ {
					key := ".." + strconv.Itoa(k)
					obj := ev.topFrame.Objects[key] 
					if obj != nil{
						if DEBUG {
							print("\t\tappending dotdotvalue (evaluated):\t", key, "=")
							PrintResult(obj)
						}
						evaluatedArgs=append(evaluatedArgs,obj)
					}
				}
				//for key,obj := range ev.topFrame.Objects {
					//if key == ".." + strconv.Itoa(n){
						//if DEBUG {
							//print("\t\t\tappending dotdotvalues to arguments:\t", key, "=")
							//PrintResult(obj)
						//}
						//n++
						//evaluatedArgs=append(evaluatedArgs,obj)
					//}
				//}
				if DEBUG {
					DumpFrames(ev)
				}
			default:
				val = EvalExprOrAssignment(ev, arg)
				if DEBUG {
					print("appending evaluated argument (non-literal):\t")
					PrintResult(val)
				}
				evaluatedArgs=append(evaluatedArgs,val)
			}
		}
	}
	return evaluatedArgs
}


func EvalCallwithEllipsis(ev *Evaluator, node *ast.CallExpr, thefunction SEXPItf) (r SEXPItf) {
	TRACE := ev.Trace
	DEBUG := ev.Debug
	funcobject := node.Fun
	funcname := funcobject.(*ast.BasicLit).Value
	if TRACE || DEBUG {
		println("EvalCallwithEllipsis: " + funcname)
	}
	argNames := getArgNames(thefunction.(*VSEXP))
	if DEBUG {
		println("\tList of arg names of function: " + funcname)
		PrintArgNames(argNames)
		println("\tList of supplied args to call to function: " + funcname)
		PrintListofAstExpressions(ev,node.Args)
	}
	frame := CollectArgsIntoFrameWithVariableArity(ev, node, argNames)
	return EvalApplyFrameToBody(ev, funcname, thefunction.(*VSEXP), frame)
}

func CollectArgsIntoFrameWithVariableArity(ev *Evaluator, node *ast.CallExpr, argNames []string) *Frame {
	DEBUG := ev.Debug
	funcobject := node.Fun
	funcname := funcobject.(*ast.BasicLit).Value
	frame := NewFrame(nil)

	if DEBUG {
		println("\tCollectArgsIntoFrameWithVariableArity:", funcname)
	}
	
	// collect tagged arguments (unevaluated) in an array of call position numbers
	taggedArgs := make(map[string]int, len(node.Args))
	for n, arg := range node.Args {
		switch arg.(type) {
		case *ast.TaggedExpr:
			a := arg.(*ast.TaggedExpr)
			if isValidArgName(argNames,a.Tag) {
				taggedArgs[a.Tag] = n + 1 // one above default zero value
				if DEBUG {
					println("\t\ttagged argument collected:", a.Tag, "  pos: ",n+1)
				}
			} else {
				if DEBUG {
					println("\t\ttagged argument found:", a.Tag, "  pos: ",n+1, "BUT NOT IN LIST OF ARGUMENT NAMES")
				}
			}
		}
	}

	// this parallel slice covers the actual arguments of the call and is set to false by default
	// for every used argument it will be set to true
	usedArgs := make([]bool, len(node.Args))

	// match defined argument names against tagged arguments
	// fieldindex: position in argument list of function definition
	// callindex:  position in parameter list of function call
	for fieldindex, fieldname := range argNames {
		if DEBUG {
			print("\t\tsearching parameter: '", fieldname,"' ")
		}
		callindex := taggedArgs[fieldname] 
		if callindex != 0 { // missing index return default zero value
			frame.Insert(fieldname, EvalExpr(ev, node.Args[callindex-1]))
			usedArgs[callindex-1] = true
			if DEBUG {
				println("=> found at position:", callindex-1, "argument number:", fieldindex)
			}
			delete(taggedArgs, fieldname)
		} else {
			if DEBUG {
				println("=> not found")
			}
		}
	}

	// TODO variadic functions do not support partial matching
	//// find partially matching tags in the remaining tagged args
	//for fieldname, callindex := range taggedArgs {
		//matches, fieldindex := tryPartialMatch(fieldname, argNames, collectedArgs, DEBUG)
		//if matches > 1 {
			//return nil, nil, errors.New("Error: argument matches multiple formal arguments" )
		//} else if matches == 1 {
			//if DEBUG {
				//println("argument", fieldname, "matches one formal argument:", argNames[fieldindex])
			//}
			//frame.Insert(fieldname, EvalExpr(ev, node.Args[callindex-1].(*ast.TaggedExpr).Rhs))
			//usedArgs[callindex-1] = true
			//delete(taggedArgs, fieldname)
		//}
	//}

	// match positional arguments up to ellipsis
	j := 0
	for n, fieldname := range argNames {
		if fieldname == "..." {
			break
		}
		for usedArgs[j] == true {
			j++
		}
		if DEBUG {
			println("\t\tcollecting positional argument:   pos:", n+1, j, fieldname)
		}
		frame.Insert(fieldname, EvalExpr(ev, node.Args[j]))
		usedArgs[j] = true
	}
	
	// collect ellipsis and remaining arguments
	n := 1
	for callindex, isUsed := range usedArgs {
		if isUsed != true {
			switch node.Args[callindex].(type) {
			case *ast.Ellipsis:
				for k:=1;k<=len(ev.topFrame.Objects);k++ {
					key := ".." + strconv.Itoa(k)
					new := ".." + strconv.Itoa(n)
					obj := ev.topFrame.Objects[key] 
					if obj != nil{
						if DEBUG {
							print("\t\tappending dotdotvalue (evaluated):", new, "= ")
							PrintResult(obj)
						}
						frame.Insert(new, obj)
						n++
					}
				}
			default:
				fieldname := ".." + strconv.Itoa(n)
				if DEBUG {
					print("\t\tappending unused argument from call: ", fieldname, "= ")
					PrintResult(EvalExpr(ev,node.Args[callindex]))
				}
				frame.Insert(fieldname, EvalExpr(ev, node.Args[callindex]))
				n++
			}
		}
	}
	if DEBUG {
		frame.Dump(ev,1)
	}
	return frame
}
