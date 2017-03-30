
package eval

import (
	"lib/ast"
)


func doAttributeReplacement(ev *Evaluator,lhs *ast.CallExpr, rhs ast.Expr) SEXPItf {
	TRACE := ev.Trace
	if TRACE {
		println("attribute replacement:")
	}
	funcobject := lhs.Fun
	attribute := funcobject.(*ast.BasicLit).Value
	defer un(trace(ev, attribute + "<-"))
	// TODO len(lhs.Args) != 1
	identifier := getIdent(ev, lhs.Args[0])
	object := ev.topFrame.Lookup(identifier)
	value := EvalExpr(ev, rhs)
	switch attribute{
	case "dim":
		// TODO instead of converting to float and back, parsing should support ints
		dim := make([]int,value.Length())
		for n,v := range value.(*VSEXP).Slice {
			dim[n]=int(v)
		}
		object.DimSet(dim)
	case "dimnames":
		if object.Dim()==nil {
			println("ERROR: 'dimnames' applied to non-array")
			return nil
		} else if value.Length() != len(object.Dim()) {
			print("ERROR: ","length of 'dimnames' [",value.Length(),"] must match that of 'dims' [",len(object.Dim()),"]\n")
			return nil
		} else {
			slice := value.(*RSEXP).Slice
			for n,v := range object.Dim() {
				if slice[n].Length() != v {
					print("ERROR: ","length of 'dimnames' [",n+1,"] not equal to array extent\n")
					return nil
				}
			}
			object.DimnamesSet(value.(*RSEXP))
		}
	case "class":
		switch value.(type){
			case *TSEXP:
				s:= value.(*TSEXP).String
				object.ClassSet(&s)
			case *NSEXP:
				object.ClassSet(nil)
			default:
				panic("attribute replacement") // TODO
		}
	}
	ev.Invisible = true // just for the following print
	return nil
}


