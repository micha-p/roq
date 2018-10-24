package main

import (
	"roq/eval"
	"testing"
)


func TestArrayArithmetic(t *testing.T) {
	quicktestSlice(t, "c(1,2,3) + c(4,5,6)", []float64{5, 7, 9}, 0)
}

func ExampleArrayArithmentic() {
	eval.EvalFileForTest("test/operator/array.r")
// Output:
//[1] 5 7 9
//[1] 2 3 4 2 3
//[1] 2 3 4 5 6
//[1] 2 4 6 8 10
//[1] 2 4 6 8 10
//[1] 0.1 0.2 0.3 0.4 0.5
//[1] 0.1 0.2 0.3 0.4 0.5
//[1] 1 2 0 1 2 0 1 2 0
//[1] 1 2 0 1 2 0 1 2 0
//[1] 1 4 9 16
//[1] 1 4 9 16
}

// TODO these comparisons need further considerations
func ExampleArrayComparison() {
	eval.EvalFileForTest("test/operator/array_comparison.r")
// Output:
//[1] NaN NaN NaN
//[1] 1 2 3
//[1] 1 NaN 3
//[1] 1 NaN 3
//[1] NaN 2 NaN
//
//[1] 4 5 6
//[1] 1 5 6
//[1] NaN NaN NaN
//[1] NaN 2 NaN
//
//[1] 4 5 6
//[1] NaN NaN NaN
//[1] NaN NaN NaN
//[1] 1 NaN NaN NaN 1
}
