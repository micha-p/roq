package main

import (
	"testing"
	"roq/eval"
)

func ExampleIdentifiers() {
	eval.EvalFileForTest("test/parser/identifiers.r")
// Output:
//Error: object 'a_2' not found
//Error: object '..b' not found
//Error: object '._c' not found
//Error: object 'a...b' not found
}


func ExampleNan() {
	eval.EvalFileForTest("test/parser/nan.r")
// Output:
//[1] NaN
//[1] NaN
//[1] NaN
//[1] NaN
}

func ExampleReturnFunction() {
	eval.EvalFileForTest("test/parser/functions.r")
// Output:
// function(x)
}

func ExampleVersionOverwrite() {
	eval.EvalFileForTest("test/parser/version.r")
// Output:
// Error in version(): could not find function "version"
// [1] 1
}

func ExampleVersionCall() {
	eval.EvalStringForTest("version()")
// Output:
// Error in version(): could not find function "version"
}
func TestVersionOverwrite(t *testing.T) {
  r := eval.EvalStringForValue("version<-1\nversion")
  if r.(*eval.VSEXP).Immediate != 1 {
    t.Error("Error in overwriting version")
  }
}
