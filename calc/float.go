package calc

import (
	"math"
)


// TODO map on slices of same length instead of single values
func FEQUAL(x float64, y float64) float64 {
	if x == y {
		return x
	} else {
		return math.NaN()
	}
}
func FUNEQUAL(x float64, y float64) float64 {
	if x != y {
		return x
	} else {
		return math.NaN()
	}
}
func FLESS(x float64, y float64) float64 { 
	if x < y {
		return y
	} else {
		return math.NaN()
	}
}
func FLESSEQUAL(x float64, y float64) float64 { 
	if x <= y {
		return y
	} else {
		return math.NaN()
	}
}
func FPLUS(x float64, y float64) float64 { 
	return x + y
}
func FMINUS(x float64, y float64) float64 { 
	return x - y
}
func FMULTIPLICATION(x float64, y float64) float64 { 
	return x * y
}
func FDIVISION(x float64, y float64) float64 { 
	return x / y
}
func FMODULUS(x float64, y float64) float64 { 
	return math.Mod(x, y)
}
func FEXPONENTIATION(x float64, y float64) float64 { 
	return math.Pow(x, y)
}

func MapIA(FUN func(float64, float64) float64, x float64, y []float64) []float64 {
	resultLen := len(y)
	r := make([]float64,resultLen)
	for n,value := range y {
		r[n]=FUN(x,value)
	}
	return r
}

func MapAI(FUN func(float64, float64) float64, x []float64, y float64) []float64 {
	resultLen := len(x)
	r := make([]float64,resultLen)
	for n,value := range x {
		r[n]=FUN(value,y)
	}
	return r
}

func MapAA(FUN func(float64, float64) float64, x []float64, y []float64) []float64 {
	lenx := len(x)
	leny := len(y)
	sliceLen := IntMin(lenx,leny)
	resultLen := IntMax(lenx,leny)
	
	r := make([]float64,resultLen)

	for base := 0; base < resultLen; base += sliceLen {
		for i := base ; (i < (base+sliceLen) && i < resultLen); i++ {
			r[i] = FUN(x[i % lenx], y[i % leny])
		}
	}
	return r
}
