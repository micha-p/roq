package main

import (
	"testing"
	"roq/eval"
)


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
    t.Error("Error in overwrituing version")
  }
}
