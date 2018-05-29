package main

import (
	"roq/eval"
//	"testing" // not needed for Examples
)

// TODO length(vector()) 		length(NULL)

func ExampleLength() {
	eval.EvalStringForTest(`a=c(11,22,33)
		length(a)
		length(c(1,2))
		length(c(1,2),c(1,3))
		length()
		length(a[1])
		length(a[2.2])
		length(a[1:2])
		length(a[1:3])
		`)
// Output:
// [1] 3
// [1] 2
// 2 arguments passed to 'length' which requires 1
// 0 arguments passed to 'length' which requires 1
// [1] 1
// [1] 1
// [1] 2
// [1] 3
}

func ExampleDim() {
	eval.EvalStringForTest(`
		x <- c(1,2,3,4,5,6)
		dim(x) <- c(2,3)
		dim(x)
		x
		`)
// Output:
//[2] 2 3
//	[,1]	[,2]	[,3]
//[1]	1	3	5      
//[2]	2	4	6
}		

func ExampleDimNames() {
	eval.EvalStringForTest(`
		x <- c(1,2,3,4,5,6)
		dim(x)
		dimnames(x) <- list(c("a1","a2"),c("b1","b2","b3"))
		dim(x) <- c(2,3)
		dimnames(x) <- list(c("a1","a2"),c("b1","b2","b3"))
		x
		dimnames(x) <- list(c("a1","a2"),c("b1","b2","b3"),c("d"))
		dimnames(x) <- list(c("a1","a2"),c("b1","b2"))
		dimnames(x) <- list(c("a1","a2"),c("b1","b2","b3","b4"))
		dimnames(x) <- list(c("a1"),c("b1","b2","b3"))
		`)
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
