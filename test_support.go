package main

import (
	"math"
	"roq/eval"
	"testing"
)

// For testing we compare the absolute difference between calculated and expected value.
// Such an absolute measure is closer to hardware restrictions of floats than a relative.

func quicktestValue(t *testing.T, s string, expected float64, epsilon float64) {
	TRACE := false
	DEBUG := false
	PRINT := false
	r := eval.EvalStringForValue(s, TRACE, DEBUG, PRINT).(*eval.VSEXP).Immediate
	diff := math.Abs(r - expected)
	if diff > epsilon {
		t.Error("Error in Test \"", s, "\" :", r, "!=", expected)
	}
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

func quicktestSlice(t *testing.T, s string, expected []float64, epsilon float64) {
	TRACE := false
	DEBUG := false
	PRINT := false
	r := eval.EvalStringForValue(s, TRACE, DEBUG, PRINT).(*eval.VSEXP).Slice
	if !testEqSlice(r, expected, epsilon) {
		t.Error("Error in Slice Test \"", s, "\" :", r, "!=", expected)
	}
}
