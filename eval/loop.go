package eval

import (
	"lib/ast"
)

func EvalLoop(ev *Evaluator, e *ast.BlockStmt, cond ast.Expr) SEXPItf {
	defer un(trace(ev, "LoopBody"))
	var evloop Evaluator
	evloop = *ev
	evloop.state = loopState
	var rstate LoopState
	for cond == nil || isTrue(EvalExpr(&evloop, cond)) {
		evloop.state = loopState
		for n := 0; n < len(e.List); n++ {
			EvalStmt(&evloop, e.List[n])
			rstate = evloop.state
			if rstate == nextState {
				break
			}
		}
		if rstate == nextState {
			continue
		}
		if rstate == breakState {
			break
		}
	}
	ev.Invisible = true
	return &NSEXP{}
}
func EvalFor(ev *Evaluator, e *ast.BlockStmt, identifier string, iterable SEXPItf) SEXPItf {
	if iterable.(*VSEXP).Slice == nil {
		panic("Vector expected")
	}
	defer un(trace(ev, "LoopBody"))
	var evloop Evaluator
	evloop = *ev
	evloop.state = loopState
	var rstate LoopState
	for _, v := range iterable.(*VSEXP).Slice {
		evloop.state = loopState
		// TODO: make use of cached position in map
		ev.topFrame.Insert(identifier, &VSEXP{Immediate: v})
		for n := 0; n < len(e.List); n++ {
			EvalStmt(&evloop, e.List[n])
			rstate = evloop.state
			if rstate == nextState {
				break
			}
		}
		if rstate == nextState {
			continue
		}
		if rstate == breakState {
			break
		}
	}
	ev.Invisible = true
	return &NSEXP{}
}

