package eval

import (
	"roq/lib/ast"
)

// Frames are derived from ast.Scopes:
// insert method inserts link struct Object with value set.
// however, data field must be set before insertion

type Frame struct {
	Outer   *Frame
	Objects map[string]SEXPItf
}


func DumpFrames(ev *Evaluator) {
	top := ev.topFrame
	if top != nil {
		top.Dump(ev, 0)
	}
}

func (f *Frame) Dump(ev *Evaluator, level int) {
	n := 1
	for key,value := range f.Objects {
		print("DUMP\t",level,"\t",n,":\t",key,"\t")
		PrintResult(value)
		n++
	}
	if f.Outer != nil {
		f.Outer.Dump(ev, level -1)
	}
}


// NewFrame creates a new scope nested in the outer scope.
func NewFrame(outer *Frame) *Frame {
	const n = 4 // initial frame capacity
	return &Frame{outer, make(map[string]SEXPItf, n)}
}

// Lookup returns the object with the given name if it is
// found in frame s, otherwise it returns nil. Outer frames
// are ignored. TODO!!!
//
func (f *Frame) Lookup(name string) SEXPItf {
	return f.Objects[name]
}

func (f *Frame) Recursive(name string) (r SEXPItf) {
	r = f.Objects[name]
	if r == nil {
		if f.Outer != nil {
			return f.Outer.Recursive(name)
		} else {
			return nil
		}
	}
	return
}

func (f *Frame) Delete(name string, DEBUG bool) () {
	r := f.Objects[name]
	if r == nil {
		if f.Outer != nil {
			f.Outer.Delete(name, DEBUG)
		} else {
			print("In remove(",name,") : object '",name,"' not found\n")
		}
	} else {
		delete(f.Objects,name)
		if DEBUG {
			println("Removed object: ",name)
		} 
	}
}

func getIdent(ev *Evaluator, ex ast.Expr) string {
	node := ex.(*ast.Ident)
	return node.Name
}

// Insert attempts to insert a named object obj into the frame s.
// If the frame already contains an object with the same name, this object is overwritten
func (s *Frame) Insert(identifier string, obj SEXPItf) (alt SEXPItf) {
	s.Objects[identifier] = obj
	return
}
