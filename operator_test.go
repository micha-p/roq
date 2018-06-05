package main

import (
	"math"
	"roq/eval"
	"testing"
)

func ExampleLogical() {
	eval.EvalFileForTest("test/operator/logical.r")
	// Output:
	//[1] 1
	//[1] 1

	//[1] 2
	//[1] 1
	//[1] 0
	//[1] 0
	//[1] 1

	//left to right
	//[1] 1
	//[1] 0
	//[1] 0
	//[1] 0
	//[1] 2
	//[1] 0
	//[1] 0
}

func ExampleArithmetic() {
	eval.EvalFileForTest("test/operator/arithmetic.r")
	// Output:
	//[1] 14
	//[1] 14
	//[1] 1
	//[1] 4
}

// For testing we compare the absolute difference between calculated and expected value.
// Such an absolute measure is closer to hardware restrictions of floats than a relative.

func quicktestValue(t *testing.T, s string, expected float64, epsilon float64) {
	r := eval.EvalStringForValue(s).(*eval.VSEXP).Immediate
	diff := math.Abs(r - expected)
	if diff > epsilon {
		t.Error("Error in Test \"", s, "\" :", r, "!=", expected)
	}
}

func quicktestSlice(t *testing.T, s string, expected []float64, epsilon float64) {
	r := eval.EvalStringForValue(s).(*eval.VSEXP).Slice
	if !testEqSlice(r, expected, epsilon) {
		t.Error("Error in Slice Test \"", s, "\" :", r, "!=", expected)
	}
}

func TestArithmetic(t *testing.T) {
	quicktestValue(t, "2.0+3.0*4.0", 14, 0)
	quicktestValue(t, "2.0+(3.0*4.0)", 14, 0)
	quicktestValue(t, "10000000*0.0000001", 1, 0)
	quicktestValue(t, "2^2", 4, 0)
	quicktestValue(t, "3.4 %% 1.0", 0.4, 0.0000001)
}

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

func TestArrayArithmetic(t *testing.T) {
	quicktestSlice(t, "c(1,2,3) + c(4,5,6)", []float64{5, 7, 9}, 0)
}

func ExampleComparison() {
	eval.EvalFileForTest("test/operator/comparison.r")
	// Output:
	//[1] 2
	//nil
	//nil
	//[1] 1
	//[1] 2
	//nil

	//nil
	//[1] 3
	//[1] 4
	//nil

	//[1] 4
	//nil
	//[1] 1
	//[1] 4
	//nil
	//[1] 4
}
