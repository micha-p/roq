package main_test 

import (
	"fmt"
	"testing"
)

func TestFail(t *testing.T) {
	x := 1
	t.Error("Fail ", x)
}

func TestOK(t *testing.T) {
	x := 1
	if x == 2 {
		t.Error("1 equal 2")
	}
}

func ExampleOK() {
        fmt.Println("Hello")
        // Output:
        // Hello
}

func ExampleFail() {
        fmt.Println("Bye")
        // Output:
        // Hello
}
