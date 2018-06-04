package main

import (
	"math"
	"testing"
	"roq/eval"
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

func quicktest(t *testing.T, s string, expected float64, epsilon float64){
  r := eval.EvalStringforValue(s)
  diff := math.Abs(r - expected)
  if diff > epsilon {
    t.Error("Error in Test",s ,":", r,"!=", expected)
  }
}

func TestArithmetic(t *testing.T) {
  quicktest(t, "2.0+3.0*4.0", 14, 0)
  quicktest(t, "2.0+(3.0*4.0)", 14, 0)
  quicktest(t, "10000000*0.0000001", 1, 0)
  quicktest(t, "2^2", 4, 0)
  quicktest(t, "3.4 %% 1.0", 0.4, 0.0000001)
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
