package main

import (
	"roq/eval"
)

func ExampleFunctionCalls() {
	eval.EvalFileForTest("test/functions/call.r")
// Output:
//[1] 1
//[1] 2
//[1] 3
//Error in f() : argument "x" is missing, with no default
//Error in f() : unused argument (pos 2, pos 3)
//Error in f() : unused argument (pos 2)
//Error in unknown(): could not find function "unknown"
}

func ExampleFunctionArguments() {
	eval.EvalFileForTest("test/functions/arguments.r")
// Output:
//[1] 3
//[1] 7
//[1] 8
//[1] 9
}

func ExampleFunctionBody() {
	eval.EvalFileForTest("test/functions/body.r")
// Output:
//[1] 3
//[1] 3
//[1] 21
//[1] 3
}

func ExampleMissingReturnValue() {
	eval.EvalFileForTest("test/functions/missing_return_value.r")
// Output:
//[1] 133
}

func ExampleEllipsis() {
	eval.EvalFileForTest("test/functions/ellipsis.r")
// Output:
//test 1:
//[1] 1
//[1] 2
//[1] 3
//test 2:
//[1] 11
//[1] 22
//[1] 33
//test 3:
//[1] 11
//[1] 22
//test 4:
//[[1]]
//[1] 22

//[[2]]
//[1] 33

//test 5:
//[[1]]
//[1] -2

//[[2]]
//[1] -1

//[[3]]
//[1] 0

//[[4]]
//[1] 22

//[[5]]
//[1] 33

//[[6]]
//[1] 5

//[[7]]
//[1] 6

}
