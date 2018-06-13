
package eval

import (
	"roq/lib/ast"
	"fmt"
)


func doAttributeReplacement(ev *Evaluator,lhs *ast.CallExpr, rhs ast.Expr) SEXPItf {
	TRACE := ev.Trace
	if TRACE {
		println("attribute replacement:")
	}
	funcobject := lhs.Fun
	attribute := funcobject.(*ast.Ident).Name
	defer un(trace(ev, attribute + "<-"))				// TODO len(lhs.Args) != 1
	identifier := getIdent(ev, lhs.Args[0])
	object := ev.topFrame.Lookup(identifier)
	value := EvalExpr(ev, rhs)
	switch attribute{
	case "dim":
		// TODO instead of converting to float and back, parsing should support ints
		dim := make([]int,value.Length())
		switch value.(type){
			case *VSEXP:
				for n,v := range value.(*VSEXP).Slice {
					dim[n]=int(v)
				}
			case *ISEXP:
				dim=value.(*ISEXP).Slice
			default:
				panic("error in dim<-")
		}
		object.DimSet(dim)
	case "dimnames":
		vlen := value.Length()
		if object.Dim()==nil {
			fmt.Printf("ERROR: 'dimnames' applied to non-array\n")
			return nil
		} else if vlen != len(object.Dim()) {
			fmt.Printf("ERROR: length of 'dimnames' [%d] must match that of 'dims' [%d]\n",vlen,len(object.Dim()))
			return nil
		} else {
			slice := value.(*RSEXP).Slice
			for n,v := range object.Dim() {
				if slice[n].Length() != v {
					fmt.Printf("ERROR: length of 'dimnames' [%d] not equal to array extent\n",n+1)
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


