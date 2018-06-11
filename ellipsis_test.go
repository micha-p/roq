package main

import (
	"roq/eval"
)

// TODO enable example test
// for unknown reasons these examples are not recognized and need to run manually

func ExampleFunctionEllipsis() {
	eval.EvalFileForTest("test/functions/ellipsis.r")
// Output:
//test 1:
//[1] 1
//[1] 2
//[1] 3
//
//test 2:
//[1] 11
//[1] 22
//[1] 33
//
//test 3:
//[1] 11
//[1] 22
//
//test 4:
//[[1]]
//[1] 22

//[[2]]
//[1] 33

//
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
//
//test 6:
//[1] 1 2 3 11 22 33 4 5
}

func ExampleFunctionNestedEllipsis() {
	eval.EvalFileForTest("test/functions/nested_ellipsis.r")
// Output:
//test 1:
//[1] 11 22 55 5 6

//test 2:
//[1] 1 2 -55 11 22 55 5 6

//test 3:
//[1] 1 2 -55 11 22 33 44 55 5 6

//test 4:
//[1] -1 0 0 0 11 22 33 0 0 0 99
}
