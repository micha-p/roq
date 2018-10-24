package main

import (
	"testing"
)

// indexing an array returns a slice
func TestIndex(t *testing.T) {
	quicktestSlice(t, "a=c(11,22,33,44,55,66); a[1]", []float64{11}, 0)
	quicktestSlice(t, "a=c(11,22,33,44,55,66); a[1.1]", []float64{11}, 0)
	quicktestSlice(t, "a=c(11,22,33,44,55,66); a[1.1+1]", []float64{22}, 0)
	quicktestSlice(t, "a=c(11,22,33);b=1;a[b]", []float64{11}, 0)
	quicktestSlice(t, "a=c(11,22,33,44,55,66); a[6]", []float64{66}, 0)
}

func TestListIndex(t *testing.T) {
	quicktestValue(t, "a=list(11,22,33,44,55,66); a[[1]]",11, 0)
}

func TestIndexSlicing(t *testing.T) {
	quicktestSlice(t, "a=c(11,22,33,44,55,66); a[2:4]", []float64{22,33,44}, 0)
	quicktestSlice(t, "a=c(11,22,33,44,55,66); b=c(1,3,6);a[b]", []float64{11,33,66}, 0)
	quicktestSlice(t, "a=c(11,22,33,44,55,66); b=c(1,3);a[2*b]", []float64{22,66}, 0)
}


func TestChainedIndex(t *testing.T) {
	quicktestSlice(t, "a=c(11,22,33,44,55,66); a[2:4][2]",[]float64{33}, 0)
	quicktestValue(t, "a=list(11,list(22,23,24),33,44,55,66); a[[2]][[3]]",24,0)
	quicktestSlice(t, "a=list(11,list(22,23,24,c(100,101,102)),33,44,55,66);a[[2]][[4]][2]",[]float64{101},0)
}


