package main

import (
	"roq/eval"
)


func ExamplePositionalParameters() {
	eval.EvalFileForTest("test/arguments/positional.r")
// Output:
//[1] 1234
//[1] 1243
//[1] 1234
//[1] 2134
//[1] 1234
}

func ExampleUnusedArguments() {
	eval.EvalFileForTest("test/arguments/unused.r")
// Output:
//Error in options(): could not find function "options"
//[1] 1234
//Error in f() : unused argument (x =)
//Error in f() : unused argument (x =)
//Error in f() : unused argument (x =)
//Error in f() : unused arguments (x =, y =)
//Error in f() : unused argument (pos 5)
//Error in f() : unused arguments (pos 5, pos 6)
}

func ExamplePartialMatching() {
	eval.EvalFileForTest("test/arguments/partial_matching.r")
// Output:
//[1] 1234
//[1] 1234
//[1] 1234
//Error in f() : argument m matches multiple formal arguments
}

// TODO functions with ellipsis

func ExampleDefaultValues() {
	eval.EvalFileForTest("test/arguments/default_values.r")
// Output:
//[1] 12
//Error in g() : argument "b" is missing, with no default
//Error in g() : argument "a" is missing, with no default
//[1] 12
//[1] 16
//[1] 52
//[1] 56
}
