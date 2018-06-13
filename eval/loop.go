package eval

import (
	"roq/lib/ast"
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
	switch iterable.(type) {
	case *VSEXP:
		return EvalForLoopOverVector(ev, e, identifier, iterable)
	case *TSEXP:
		return EvalForLoopOverStringArray(ev, e, identifier, iterable)
	case *RSEXP:
		return EvalForLoopOverList(ev, e, identifier, iterable)
	default:
		panic("For loop over unknown s-expression\n")
	}
}


func EvalForLoopOverList(ev *Evaluator, e *ast.BlockStmt, identifier string, iterable SEXPItf) SEXPItf {
	if iterable.(*RSEXP).Slice == nil {
		panic("List expected\n")
	}
	defer un(trace(ev, "LoopOverList"))
	var evloop Evaluator
	evloop = *ev
	evloop.state = loopState
	var rstate LoopState
	for _, v := range iterable.(*RSEXP).Slice {
		evloop.state = loopState
		ev.topFrame.Insert(identifier, v)
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

func EvalForLoopOverVector(ev *Evaluator, e *ast.BlockStmt, identifier string, iterable SEXPItf) SEXPItf {
	if iterable.(*VSEXP).Slice == nil {
		panic("Vector expected\n")
	}
	defer un(trace(ev, "LoopOverVector"))
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

func EvalForLoopOverStringArray(ev *Evaluator, e *ast.BlockStmt, identifier string, iterable SEXPItf) SEXPItf {
	if iterable.(*TSEXP).Slice == nil {
		panic("String array expected\n")
	}
	defer un(trace(ev, "LoopOverStringArray"))
	var evloop Evaluator
	evloop = *ev
	evloop.state = loopState
	var rstate LoopState
	for _, v := range iterable.(*TSEXP).Slice {
		evloop.state = loopState
		// TODO: make use of cached position in map
		ev.topFrame.Insert(identifier, &TSEXP{String: v})
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
