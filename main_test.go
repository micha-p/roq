package main_test 

import (
	"roq/eval"
	"testing"
)


func TestOK(t *testing.T) {
	x := 1
	if x == 2 {
		t.Error("1 equal 2")
	}
}

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
