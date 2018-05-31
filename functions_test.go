package main

import (
	"roq/eval"
//	"testing" // not needed for Examples
)

// TODO length(vector()) 		length(NULL)

func ExampleFunctions() {
	eval.EvalStringForTest(`
		f<-function(x){1}
		f(1)
		g<-function(x){x+1}
		g(1)
		f<-function(x){x+2}
		f(1)
		f()
		f(1,2)
		f(x=1,2)
		unknown(1)
		`)
// Output:
//[1] 1
//[1] 2
//[1] 3
//Error in f() : argument "x" is missing, with no default
//Error in f() : unused argument (pos 2, pos 3)
//Error in f() : unused argument (pos 2)
//Error: could not find function "unknown"
}

func ExampleArguments() {
	eval.EvalStringForTest(`
		b<-function(c,d){c+d}
		b(c=1,d=2)
		b(3,4)
		b(3,4+1)
		a <- 5
		b(c=3,d=a+1)
		`)
// Output:
//[1] 3
//[1] 7
//[1] 8
//[1] 9
}

func ExampleFunctionBody() {
	eval.EvalStringForTest(`
		f<-function(a,b)a+b
		f(1,2)
		f<-function(a,b) a+b
		f(1,2)
		f<-function(a,b) a+b; 10+11 # 10*11 is an extra statement
		f(1,2)
		`)
// Output:
//[1] 3
//[1] 3
//[1] 21
//[1] 3
}

func ExampleMissingReturnValue() {
	eval.EvalStringForTest(`
a<-133
f<-function(a,b){a<-2}
f(1,2) 		# there is no return value
f<-function(a,b) a<-2
f(1,2) 		# there is no return value
a
		`)
// Output:
//[1] 133
}
