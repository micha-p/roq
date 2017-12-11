package eval

import (
	"math"
	"fmt"
	"roq/lib/ast"
	"roq/lib/token"
)

// length is used instead of linecount

// visibility is stored in the evaluator and unset after every print
// TODO typeswitch should depend on Kind
func PrintResult(ev *Evaluator, r SEXPItf) {
	if ev.Invisible {
		ev.Invisible = false
		return
	} else if r == nil {
		println("FALSE/NULL")
	} else {
		switch r.(type) {
		case *VSEXP:
			PrintResultV(ev, r.(*VSEXP))
		case *ISEXP:
			PrintResultI(ev, r.(*ISEXP))
		case *RSEXP:
			PrintResultR(ev, r.(*RSEXP))
		case *TSEXP:
			PrintResultT(ev, r.(*TSEXP))
		case *ESEXP:
			PrintResultE(ev, r.(*ESEXP))
		case *NSEXP:
			println("NULL")
		default:
			panic("?prnt")
		}
	}
}

func PrintResultR(ev *Evaluator, r *RSEXP) {
	if r == nil {
		println("ERROR: uncatched NULL pointer: ", r)
		return
	}
	if r.Slice == nil {
		println("[[1]]")
		PrintResult(ev, r.CAR)
		println()
		println("[[2]]")
		PrintResult(ev, r.CDR)
		println()
	} else {
		for n, v := range r.Slice {
			print("[[", n+1, "]]\n")
			PrintResult(ev, v)
			println()
		}
	}
}

func PrintResultT(ev *Evaluator, r *TSEXP) {
	if r.Slice == nil {
		println("[1]", "\""+r.String+"\"")
	} else {
		print("[", len(r.Slice), "]")
		for _, v := range r.Slice {
			print(" \"", v, "\"")
		}
		println()
	}
}

func PrintResultI(ev *Evaluator, r *ISEXP) {
	rdim := r.Dim()
	if rdim == nil {
		println("[1]", r.Integer)
	} else {
		print("[", len(rdim), "]")
		for _, v := range rdim {
			fmt.Printf(" %d", v)
		}
		println()
	}
}

func PrintResultE(ev *Evaluator, r *ESEXP) {
	switch r.Kind {
	case token.ILLEGAL:
		DEBUG := ev.Debug
		if DEBUG {
			println("ILLEGAL RESULT")
		}
	case token.VERSION:
		PrintVersion()
	case token.EOF:
		DEBUG := ev.Debug
		if DEBUG {
			println("EOF")
		}
	default:
		println(r.Message)
	}
}

func PrintResultV(ev *Evaluator, r *VSEXP) {

	DEBUG := ev.Debug
	if r== nil {
		println("nil")
	} else if r.Body != nil {
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
	} else {
		if r.Slice == nil {
			if r.Immediate == math.NaN(){
				println("NAN")
			} else {
				fmt.Printf("[1] %g\n", r.Immediate) // R has small e for exponential format
			}
		} else {
			rdim := r.Dim()
			if rdim == nil {
				print("[", r.Length(), "]")
				printArray(r.Slice)
			} else if len(rdim) == 2 && r.Dimnames() != nil {
				printMatrixDimnames(r.Slice,
					rdim[0],
					rdim[1],
					r.Dimnames().Slice[0].(*TSEXP).Slice,
					r.Dimnames().Slice[1].(*TSEXP).Slice)
			} else if len(rdim) == 2 {
				printMatrix(r.Slice, rdim[0], rdim[1])
			} else {
				print("[")
				for n, v := range rdim {
					if n > 0 {
						print(",")
					}
					fmt.Printf("%d", v)
				}
				print("]")
				printArray(r.Slice)
			}
		}
	}
}

func printArray(slice []float64) {
	for _, v := range slice {
		fmt.Printf(" %g", v)
	}
	println()
}

func printMatrixDimnames(slice []float64, rows int, cols int, rownames []string, colnames []string) {
	for col := 0; col < cols; col++ {
		if col < len(colnames) {
			print("\t", colnames[col])
		} else {
			print("\t[,", col+1, "]")
		}
	}
	println()
	for row := 0; row < rows; row++ {
		if row < len(rownames) {
			print(rownames[row])
		} else {
			print("[", row+1, ",]")
		}
		for col := 0; col < cols; col++ {
			fmt.Printf(" %7g", slice[row+rows*col])
		}
		println()
	}
}

func printMatrix(slice []float64, rows int, cols int) {
	for col := 0; col < cols; col++ {
		print("\t[,", col+1, "]")
	}
	println()
	for row := 0; row < rows; row++ {
		print("[", row+1, ",]")
		for col := 0; col < cols; col++ {
			fmt.Printf(" %7g", slice[row+rows*col])
		}
		println()
	}
}
