package main

import (
	"roq/eval"
	"testing"
)

func TestIndex(t *testing.T) {
	quicktestSlice(t, "a=c(11,22,33,44,55,66);a[1]", []float64{11}, 0)
	quicktestSlice(t, "a=c(11,22,33,44,55,66);a[6.5]", []float64{66}, 0)
	quicktestSlice(t, "a=c(11,22,33,44,55,66);a[2:4]", []float64{22,33,44}, 0)
}

func TestListIndex(t *testing.T) {
	quicktestValue(t, "a=list(11,22,33,44,55,66);a[[1]]",11, 0)
}


// TODO testing manually is fine for all, there seems to be a problem with EvalStringforTest
func TestChainedIndex(t *testing.T) {
	quicktestValue(t, "a=list(11,22,33,44,55,66);a[[1]]",11, 0)
//	quicktestValue(t, "a=c(11,22,33,44,55,66);a[2:4][2]",33,0)
	quicktestValue(t, "a=list(11,list(22,23,24),33,44,55,66);a[[2]][[3]]",24,0)
//	quicktestValue(t, "a=list(11,list(22,23,24,c(100,101,102)),33,44,55,66);a[[2]][[4]][2]",101,0)
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
