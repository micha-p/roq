package main

import (
	"roq/eval"
	"testing"
)

// eval unquotes only directly given identifers but evaluates normal expressions
func TestQuoteEval(t *testing.T) {
	quicktestValue(t, "a<-quote(1+2);eval(a)",3, 0)
	quicktestValue(t, "eval(quote(1+2))",3, 0)
	quicktestValue(t, "a<-quote(1+2);eval(10+eval(a))",13, 0)
	quicktestValue(t, "a<-quote(f<-function(a){a*2});eval(a);f(3)",6,0)
}

func ExampleWrongQuoteLevel() {
	eval.EvalStringForTest("a<-quote(1+2);eval(10+a)")
// Output:
// Error: non-numeric argument to binary operator
}

func ExamplePrintQuoted() {
	eval.EvalStringForTest("a<-quote(1+2);a")
// Output:
//1  BinaryExpr{
//2  .  1
//3  .  +
//4  .  2
//5  }
}

func ExampleArbitraryCall() {
	eval.EvalStringForTest("f<-function(a){a*2};call(\"f\",3)")
// Output:
//1  ArbitraryCallExpr{
//2  .  "f"
//3  .  [
//4  .  .  1: 3
//5  .  ]
//6  }
}

// eval unquotes only directly given identifers but evaluates normal expressions
func TestConstructedCall(t *testing.T) {
	quicktestValue(t, "f<-function(a){a*2};eval(call(\"f\",3))",6, 0)
	quicktestValue(t, "f<-function(a){a*2};b=\"f\";eval(call(b,3))",6, 0)
}
