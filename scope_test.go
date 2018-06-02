package main

import (
	"roq/eval"
)

func ExampleScope() {
	eval.EvalFileForTest("test/scope/closure.r")
// Output:
//[1] 3
//[1] 11
//[1] 24
//[1] 22
//[1] 1
//[1] 103
//[1] 1
}

func ExampleBlocks() {
	eval.EvalFileForTest("test/scope/blocks.r")
// Output:
//[1] 11
//[1] 12
//[1] 12

//[1] 3

//[1] 12
//[1] 22
//[1] 22
}

func ExampleGlobals() {
	eval.EvalFileForTest("test/scope/globals.r")
// Output:
//[1] 2
//[1] 4

//[1] 10
//[1] 11
}


// TODO 
// an extremely difficult problem, which needs clarification
// R seems to behave wrong here
func ExampleCallingScope() {
	eval.EvalFileForTest("test/scope/calling_scope.r")
// Output:
//[1] 1
//[1] 2
//[1] 3
//[1] 201
//[1] 101
}
