package main

import (
	"math"
	"roq/eval"
	"testing"
)

func testEqSlice(a, b []float64, epsilon float64) bool {

	if a == nil || b == nil {
		panic("Testing empty slice for equality")
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		diff := math.Abs(a[i] - b[i])
		if diff > epsilon {
			return false
		}
	}
	return true
}
func quicktestSlice(t *testing.T, s string, expected []float64, epsilon float64) {
	r := eval.EvalStringForValue(s).(*eval.VSEXP).Slice
	if !testEqSlice(r, expected, epsilon) {
		t.Error("Error in Slice Test \"", s, "\" :", r, "!=", expected)
	}
}

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
