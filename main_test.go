package main

import (
	"roq/eval"
//	"testing" // not needed for Examples
)

// Example testing does not work with print, only with fmt

func ExampleArithmeticOperators() {
	eval.EvalStringForTest("1+2\n2*3.1\n5-1\n7/8\n11%%2\n3^2")
	// Output:
	// [1] 3
	// [1] 6.2
	// [1] 4
	// [1] 0.875
	// [1] 1
	// [1] 9
}

func ExamplePrecedenceOperators() {
	eval.EvalStringForTest(`
		1+2*3
		(1+2)*3
		10-3^2
		2*3/4`)
	// Output:
	// [1] 7
	// [1] 9
	// [1] 1
	// [1] 1.5
}

func ExampleMissingValues() {
	eval.EvalStringForTest(`
		r = 1  +2 %% 10 + .
		r
		r = 1  +2 %% 10 + NaN
		r
		r = 1  +2 %% 10 + NA
		r`)
	// Output:
	// [1] NaN
	// [1] NaN
	// [1] NaN
}


func ExampleAssignment() {
	eval.EvalStringForTest(`
		a <- 3
		4 -> d
		b = 2 +a 
		a
		d
		b`)
	// Output:
	// [1] 3
	// [1] 4
	// [1] 5
}
