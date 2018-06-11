package eval

import (
	"fmt"
)

func EvalApply(ev *Evaluator, funcname string, f *VSEXP, argNames []string, evaluatedArgs []SEXPItf) (r SEXPItf) {
	TRACE := ev.Trace
	DEBUG := ev.Debug

	ev.openFrame()
	defer ev.closeFrame()

	if (TRACE || DEBUG) {
		println("Insert arguments of call to function \"" + funcname + "\" into new top frame:")
	}
	for n, fieldname := range argNames {
		if (TRACE || DEBUG) {
			print("\targ[",n+1, "]\t", fieldname)
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
				print("\t= ")
				PrintResult(value)
			}
			ev.topFrame.Insert(fieldname, value)
		} else {
			if (TRACE || DEBUG) {
				println()
			}
		}
	}
	if DEBUG {
		DumpFrames(ev)
	}
	if (TRACE || DEBUG) {
		println("Eval body of function \"" + funcname + "\":")
	}
	if f.Body==nil{
		panic("EvalCall: body==nil")
	}
	r=EvalStmt(ev, f.Body)
	if r != nil {
		if (TRACE || DEBUG) {
			println("Return from function \"" + funcname + "\" with result: ")
			PrintResult(r)
			println("End of result")
		}
		return r
	} else {
		return &NSEXP{}
	}
}

func EvalApplyFrameToBody(ev *Evaluator, funcname string, f *VSEXP, frame *Frame) (r SEXPItf) {
	TRACE := ev.Trace
	DEBUG := ev.Debug
	if (TRACE || DEBUG) {
		println("EvalApplyFrameToBody \"" + funcname + "\" ENTERING Frame")
	}

	frame.Outer = ev.topFrame
	ev.topFrame = frame
	defer ev.closeFrame()

	if DEBUG {
		DumpFrames(ev)
	}
	if f.Body==nil{
		panic("EvalCall: function body==nil")
	}
	r=EvalStmt(ev, f.Body)
	if r != nil {
		if (TRACE || DEBUG) {
			println("Return from function \"" + funcname + "\" with result: ")
			PrintResult(r)
			println("End of result")
		}
		return r
	} else {
		return &NSEXP{}
	}
}
