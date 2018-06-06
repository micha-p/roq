package main

import (
	"roq/eval"
)


func ExampleFlowcontrolIf() {
	eval.EvalFileForTest("test/flowcontrol/if.r")
// Output:
//[1] "T"
//[1] "F"
}

func ExampleFlowcontrolFor() {
	eval.EvalFileForTest("test/flowcontrol/for.r")
// Output:
//[1] 1
//[1] 2
//[1] 3

//[1] 1
//[1] 2
//[1] 3

//[1] 1
//[1] 2
//[1] 4
//[1] 5
}
