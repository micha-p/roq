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

func ExampleAdd() {
        eval.EvalTest("1+2")
        // Output:
        // [1] 3
}
