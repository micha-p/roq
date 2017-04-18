package eval

import (
	"lib/ast"
)

// Frames are derived from ast.Scopes:
// insert method inserts link struct Object with value set.
// however, data field must be set before insertion

type Frame struct {
	Outer   *Frame
	Objects map[string]SEXPItf
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
			print("Error: object '", name, "' not found\n")
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
	node := ex.(*ast.BasicLit)
	return node.Value
}

// Insert attempts to insert a named object obj into the frame s.
// If the frame already contains an object alt with the same name, this object is overwritten
func (s *Frame) Insert(identifier string, obj SEXPItf) (alt SEXPItf) {
	s.Objects[identifier] = obj
	return
}
