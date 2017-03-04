package eval

import (
	"lib/ast"
	"lib/token"
)

// strongly stripped down call to c()
// Therefore, all elements are evaluated withon the context of the call
// TODO recursive=TRUE/FALSE
// TODO faster vector literals, composed just of floats

func EvalVector(ev *Evaluator, node *ast.VectorExpr) (r *SEXP) {
	TRACE := ev.trace
	if TRACE {
		println("VectorExpr")
	}

	evaluatedArgs := make(map[int]float64)
	for n, v := range node.Args { // TODO: strictly left to right
		val := EvalExpr(ev, v)
		evaluatedArgs[n] = val.Value
	}
	c := make([]float64, len(evaluatedArgs))
	for n,v := range evaluatedArgs {
		c[n] = v
	}

	return &SEXP{ValuePos: node.Start, Kind: token.COLUMN, Array: c}
}
