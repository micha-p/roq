package eval

import (
	"roq/version"
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
		fmt.Printf("FALSE/NULL")
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
		println("ERROR: uncatched NULL pointer: ", r) // TODO fatalState
		return
	}
	if r.Slice == nil {
		fmt.Printf("[[1]]\n")
		PrintResult(ev, r.CAR)
		fmt.Printf("\n")
		fmt.Printf("[[2]]\n")
		PrintResult(ev, r.CDR)
		fmt.Printf("\n")
	} else {
		for n, v := range r.Slice {
			fmt.Printf("[[%d]]\n", n+1)
			PrintResult(ev, v)
			fmt.Printf("\n")
		}
	}
}

func PrintResultT(ev *Evaluator, r *TSEXP) {
	if r.Slice == nil {
		fmt.Printf("[1] \"%s\"",r.String)
	} else {
		fmt.Printf("[%d]", len(r.Slice))
		for _, v := range r.Slice {
			fmt.Printf(" \"%s\"", v)
		}
		fmt.Printf("\n")
	}
}

func PrintResultI(ev *Evaluator, r *ISEXP) {
	rdim := r.Dim()
	if rdim == nil {
		fmt.Printf("[1] %d\n", r.Integer)
	} else {
		fmt.Printf("[%d]", len(rdim))
		for _, v := range rdim {
			fmt.Printf(" %d", v)
		}
		fmt.Printf("\n")
	}
}

func PrintResultE(ev *Evaluator, r *ESEXP) {
	DEBUG := ev.Debug
	switch r.Kind {
	case token.ILLEGAL:
		if DEBUG {
			println("ILLEGAL RESULT")
		}
	case token.VERSION:
		version.PrintVersion(ev.Major,ev.Minor)
	case token.EOF:
		if DEBUG {
			println("EOF AS RESULT")
		}
	default:
		fmt.Printf("%s",r.Message)
	}
}

func PrintResultV(ev *Evaluator, r *VSEXP) {
	DEBUG := ev.Debug
	if r== nil {
		fmt.Printf("nil")
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
				fmt.Printf("NAN\n")
			} else {
				fmt.Printf("[1] %g\n", r.Immediate) // R has small e for exponential format
			}
		} else {
			rdim := r.Dim()
			if rdim == nil {
				fmt.Printf("[%d]", r.Length())
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
				fmt.Printf("[")
				for n, v := range rdim {
					if n > 0 {
						fmt.Printf(",")
					}
					fmt.Printf("%d", v)
				}
				fmt.Printf("]")
				printArray(r.Slice)
			}
		}
	}
}

func printArray(slice []float64) {
	for _, v := range slice {
		fmt.Printf(" %g", v)
	}
	fmt.Printf("\n")
}

func printMatrixDimnames(slice []float64, rows int, cols int, rownames []string, colnames []string) {
	for col := 0; col < cols; col++ {
		if col < len(colnames) {
			fmt.Printf("\t%s", colnames[col])
		} else {
			fmt.Printf("\t[,%d]", col+1)
		}
	}
	fmt.Printf("\n")
	for row := 0; row < rows; row++ {
		if row < len(rownames) {
			fmt.Printf("%s",rownames[row])
		} else {
			fmt.Printf("[%d]", row+1)
		}
		for col := 0; col < cols; col++ {
			fmt.Printf("\t%s", fmt.Sprintf("%g",slice[row+rows*col]))
		}
		fmt.Printf("\n")
	}
}

func printMatrix(slice []float64, rows int, cols int) {
	for col := 0; col < cols; col++ {
		fmt.Printf("\t[,%d]", col+1)
	}
	fmt.Printf("\n")
	for row := 0; row < rows; row++ {
		fmt.Printf("[%d]", row+1)
		for col := 0; col < cols; col++ {
			fmt.Printf("\t%s", fmt.Sprintf("%g",slice[row+rows*col]))
		}
		fmt.Printf("\n")
	}
}
