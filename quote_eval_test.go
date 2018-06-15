package main

import (
	"roq/eval"
	"testing"
)

// eval unquotes only directly given identifers but evaluates normal expressions
func TestQuoteEval(t *testing.T) {
	quicktestValue(t, "a<-quote(1+2);eval(a)",3, 0)
	quicktestValue(t, "a<-quote(1+2);eval(10+eval(a))",13, 0)
	quicktestValue(t, "a<-quote(f<-function(a){a*2});eval(a);f(3)",6,0)
}

func ExampleWrongQuoteLevel() {
	eval.EvalStringForTest("a<-quote(1+2);eval(10+a)")
// Output:
// Error: non-numeric argument to binary operator
}
