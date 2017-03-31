package print

import (
	"fmt"
	"lib/token"
	"lib/ast"
	"reflect"
	"eval"
)

// length is used instead of linecount

// visibility is stored in the evaluator and unset after every print
// TODO typeswitch should depend on Kind
func PrintResult(ev *eval.Evaluator, r eval.SEXPItf) {
	DEBUG := ev.Debug
	if DEBUG {
		print("print: ")
	}

	if ev.Invisible {
		ev.Invisible = false
		return
	} else if r == nil {
		println("FALSE")
	} else {
		if DEBUG {
			givenType := reflect.TypeOf(r)
			println(givenType.String(),r)
		}
		switch r.(type) {
		case *eval.VSEXP:
			PrintResultV(ev, r.(*eval.VSEXP))
		case *eval.RSEXP:
			PrintResultR(ev, r.(*eval.RSEXP))
		case *eval.TSEXP:
			PrintResultT(ev, r.(*eval.TSEXP))
		case *eval.NSEXP:
			println("NULL")
		default:
			println("?prnt")
		}
	}
}

func PrintResultR(ev *eval.Evaluator, r *eval.RSEXP) {
	if r == nil {
		println("ERROR: uncatched NULL pointer: ",r)
		return
	}
	if r.Slice==nil {
		print("[1] ")
		PrintResult(ev, r.CAR)
		print("[2] ")
		PrintResult(ev, r.CDR)
	} else {
		for n,v := range r.Slice {
			print("[[",n+1,"]]\n")
			PrintResult(ev,v)
			println()
		}
	}
}

func PrintResultT(ev *eval.Evaluator, r *eval.TSEXP) {
			if r.Slice==nil {
				println("[1]", "\""+r.String+"\"")
			} else {
				print("[", len(r.Slice), "]")
				for _, v := range r.Slice {
					print(" \"",v,"\"")
				}
			}
			println()
}

func PrintResultV(ev *eval.Evaluator, r *eval.VSEXP) {

	DEBUG := ev.Debug
		switch r.Kind() {
		case token.ILLEGAL:
			if DEBUG {
				println("ILLEGAL RESULT")
			}
		case token.FLOAT:
			if r.Slice==nil {
				fmt.Printf("[1] %g\n", r.Immediate) // R has small e for exponential format
			} else {
				rdim := r.Dim()
				if rdim==nil {
					print("[", r.Length(), "]")
					printArray(r.Slice)
				} else if (len(rdim)==2 && r.Dimnames() != nil) {
					printMatrixDimnames(r.Slice, 
										rdim[0],
										rdim[1],
										r.Dimnames().Slice[0].(*eval.TSEXP).Slice,
										r.Dimnames().Slice[1].(*eval.TSEXP).Slice)
				} else if (len(rdim)==2) {
					printMatrix(r.Slice, rdim[0],rdim[1])
				} else {
					print("[")
					for n, v := range rdim {
						if n>0 {print(",")}
						fmt.Printf("%d", v)
					}
					print("]")
					printArray(r.Slice)
				}
			}
		case token.INT:
			rdim := r.Dim()
			if rdim==nil {
				println("[1]", r.Integer)
			} else {
				print("[", len(rdim), "]")
				for _, v := range rdim {
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
			println(r.Kind().String())
		}
}

func printArray(slice []float64){
	for _, v := range slice {
		fmt.Printf(" %g", v)
	}
	println()
}

func printMatrixDimnames(slice []float64, rows int, cols int, rownames []string, colnames []string){
	for col:=0;col<cols;col++ {
		if col<len(colnames) {
			print("\t",colnames[col])
		} else {
			print("\t[,",col+1, "]")
		}
	}
	println()
	for row:=0; row< rows; row++ {
		if row<len(rownames) {
			print(rownames[row])
		} else {
			print("[",row+1, ",]")
		}
		for col:=0;col<cols;col++ {
			fmt.Printf(" %7g", slice[row+rows*col])
		}
		println()
	}
}

func printMatrix(slice []float64, rows int, cols int){
	for col:=0;col<cols;col++ {
		print("\t[,",col+1, "]")
	}
	println()
	for row:=0; row< rows; row++ {
		print("[",row+1, ",]")
		for col:=0;col<cols;col++ {
			fmt.Printf(" %7g", slice[row+rows*col])
		}
		println()
	}
}
