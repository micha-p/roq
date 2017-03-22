package print

import (
	"fmt"
	"lib/token"
	"lib/ast"
	"reflect"
	"eval"
)

// visibility is stored in the evaluator and unset after every print
// TODO typeswitch should depend on Kind
func PrintResult(ev *eval.Evaluator, r *eval.SEXP) {

	DEBUG := ev.Debug
	if DEBUG {
		givenType := reflect.TypeOf(r)
		print("print: ", givenType.String(), ": ", r.Kind.String(), ": ")
	}

	if ev.Invisible {
		ev.Invisible = false
		return
	} else if r == nil {
		println("FALSE")
	} else {
		switch r.Kind {
		case token.SEMICOLON:
			if DEBUG {
				println("Semicolon")
			}
		case token.ILLEGAL:
			if DEBUG {
				println("ILLEGAL RESULT")
			}
		case token.FLOAT:
			if r.Slice==nil {
				fmt.Printf("[1] %g\n", r.Immediate) // R has small e for exponential format
			} else {
				printArray(r.Slice)
			}
		case token.INT:
			if r.Dim==nil {
				println("[1]", r.Offset)
			} else {
				print("[", len(r.Dim), "]")
				for _, v := range r.Dim {
					fmt.Printf(" %d", v)
				}
			}
			println()
		case token.FUNCTION:
			if DEBUG {
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
				println(")")
			}
		case token.VERSION:
			PrintVersion()
		default:
			if DEBUG {
				println("default print")
			}
			println(r.Kind.String())
		}
	}
}

func printArray(slice []float64){
	print("[", len(slice), "]")
	for _, v := range r.slice {
		fmt.Printf(" %g", v) // R has small e for exponential format
	}
	println()
}


