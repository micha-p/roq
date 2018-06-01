package main

import (
	"roq/eval"
)

// TODO length(vector()) 		length(NULL)

func ExampleLength() {
	eval.EvalFileForTest("test/dimensions/length.r")
// Output:
// [1] 3
// [1] 2
// [1] 1
// [1] 1
// [1] 2
// [1] 3
}

func ExampleDim() {
	eval.EvalFileForTest("test/dimensions/dim.r")
// Output:
//[2] 2 3
//	[,1]	[,2]	[,3]
//[1]	1	3	5      
//[2]	2	4	6
}		

func ExampleDimNames() {
	eval.EvalFileForTest("test/dimensions/dimnames.r")
// Output:
//[1] 0
//ERROR: 'dimnames' applied to non-array
//	b1	b2	b3
//a1	1	3	5
//a2	2	4	6
//ERROR: length of 'dimnames' [3] must match that of 'dims' [2]
//ERROR: length of 'dimnames' [2] not equal to array extent
//ERROR: length of 'dimnames' [2] not equal to array extent
//ERROR: length of 'dimnames' [1] not equal to array extent
}
