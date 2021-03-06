// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ast declares the types used to represent syntax trees for Go
// packages.
//
package ast

import (
	"roq/lib/token"
	"unicode"
	"unicode/utf8"
)

// ----------------------------------------------------------------------------
// Interfaces
//
// There are 3 main classes of nodes: Expressions and type nodes,
// statement nodes, and declaration nodes. The node names usually
// match the corresponding Go spec production names to which they
// correspond. The node fields correspond to the individual parts
// of the respective productions.
//
// All nodes contain position information marking the beginning of
// the corresponding source text segment; it is accessible via the
// Pos accessor method. Nodes may contain additional position info
// for language constructs where comments may be found between parts
// of the construct (typically any larger, parenthesized subpart).
// That position information is needed to properly position comments
// when printing the construct.

// All node types implement the Node interface.
type Node interface {
	Pos() token.Pos // position of first character belonging to the node
	End() token.Pos // position of first character immediately after the node
}

// All expression nodes implement the Expr interface.
type Expr interface {
	Node
	exprNode()
}

// All statement nodes implement the Stmt interface.
type Stmt interface {
	Node
	stmtNode()
}


// ----------------------------------------------------------------------------
// Expressions and types

// A Field represents a Field declaration list in a struct type,
// a method list in an interface type, or a parameter/result declaration
// in a signature.
//
type Field struct {
	Names []*Ident  // field/method/parameter names; or nil if anonymous field
	Type  Expr      // field/method/parameter type
	Tag   *BasicLit // field tag; or nil
	Default  Expr
}

func (f *Field) Pos() token.Pos {
	if len(f.Names) > 0 {
		return f.Names[0].Pos()
	}
	return f.Type.Pos()
}

func (f *Field) End() token.Pos {
	if f.Tag != nil {
		return f.Tag.End()
	}
	return f.Type.End()
}

// A FieldList represents a list of Fields, enclosed by parentheses or braces.
type FieldList struct {
	Opening token.Pos // position of opening parenthesis/brace, if any
	List    []*Field  // field list; or nil
	Closing token.Pos // position of closing parenthesis/brace, if any
}

func (f *FieldList) Pos() token.Pos {
	if f.Opening.IsValid() {
		return f.Opening
	}
	// the list should not be empty in this case;
	// be conservative and guard against bad ASTs
	if len(f.List) > 0 {
		return f.List[0].Pos()
	}
	return token.NoPos
}

func (f *FieldList) End() token.Pos {
	if f.Closing.IsValid() {
		return f.Closing + 1
	}
	// the list should not be empty in this case;
	// be conservative and guard against bad ASTs
	if n := len(f.List); n > 0 {
		return f.List[n-1].End()
	}
	return token.NoPos
}

// NumFields returns the number of (named and anonymous fields) in a FieldList.
func (f *FieldList) NumFields() int {
	n := 0
	if f != nil {
		for _, g := range f.List {
			m := len(g.Names)
			if m == 0 {
				m = 1 // anonymous field
			}
			n += m
		}
	}
	return n
}

// An expression is represented by a tree consisting of one
// or more of the following concrete expression nodes.
//
type (
	// A BadExpr node is a placeholder for expressions containing
	// syntax errors for which no correct expression nodes can be
	// created.
	//
	BadExpr struct {
		From, To token.Pos // position range of bad expression
	}

	// An Ident node represents an identifier.
	Ident struct {
		NamePos token.Pos // identifier position
		Name    string    // identifier name
	}

	// An Ellipsis node stands for the "..." type in a
	// parameter list or the "..." length in an array type.
	//
	Ellipsis struct {
		ValuePos token.Pos // position of "..."
		Elt      Expr      // ellipsis element type (parameter lists only); or nil
	}

	// A BasicLit node represents a literal of basic type.
	BasicLit struct {
		ValuePos token.Pos   // literal position
		Kind     token.Token // token.INT, token.FLOAT, token.IMAG, token.CHAR, or token.STRING
		Value    string      // literal string; e.g. 42, 0x7f, 3.14, 1e-9, 2.4i, 'a', '\x7f', "foo" or `\m\n\o`
	}

	// A FuncLit node represents a function literal.
	FuncLit struct {
		Type *FuncType  // function type
		Body *BlockStmt // function body
	}

	// A CompositeLit node represents a composite literal.
	CompositeLit struct {
		Type   Expr      // literal type; or nil
		Left   token.Pos // position of "{"
		Elts   []Expr    // list of composite elements; or nil
		Right  token.Pos // position of "}"
	}

	// A ParenExpr node represents a parenthesized expression.
	ParenExpr struct {
		Left   token.Pos // position of "("
		X      Expr      // parenthesized expression
		Right  token.Pos // position of ")"
	}
	EvalExpr struct {
		Left  token.Pos
		X     Expr
		Right token.Pos // position of ")"
	}

	// A SelectorExpr node represents an expression followed by a selector.
	SelectorExpr struct {
		X   Expr   // expression
		Sel *Ident // field selector
	}

	// An IndexExpr node represents an expression followed by an index.
	IndexExpr struct {
		Array  Expr      // expression
		Left   token.Pos // position of "["
		Index  Expr      // index expression
		Right  token.Pos // position of "]"
	}

	// An IndexExpr node represents an expression followed by an index.
	ListIndexExpr struct {
		Array  Expr      // expression
		Left   token.Pos // position of "["
		Index  Expr      // index expression
		Right  token.Pos // position of "]"
	}

	// A CallExpr node represents an expression followed by an argument list.
	CallExpr struct {
		Fun      Expr      // function expression
		Left     token.Pos // position of "("
		Args     []Expr    // function arguments; or nil
		Ellipsis token.Pos // position of "...", if any
		Right    token.Pos // position of ")"
	}

	// Same as above, but expecting string literal as function identifier
	ArbitraryCallExpr struct {
		Fun      Expr      // expression for function name should result in TSEXP
		Left     token.Pos // position of "("
		Args     []Expr    // function arguments; or nil
		Ellipsis token.Pos // position of "...", if any
		Right    token.Pos // position of ")"
	}

	// A UnaryExpr node represents a unary expression.
	//
	UnaryExpr struct {
		OpPos token.Pos   // position of Op
		Op    token.Token // operator
		X     Expr        // operand
	}

	// A BinaryExpr node represents a binary expression.
	BinaryExpr struct {
		X     Expr        // left operand
		OpPos token.Pos   // position of Op
		Op    token.Token // operator
		Y     Expr        // right operand
	}

	QuotedExpr struct {
		Left  token.Pos
		X     Expr
		Right token.Pos // position of ")"
	}

	TaggedExpr struct {
		X     Expr // left operand
		Tag   string
		OpPos token.Pos // position of "="
		Rhs   Expr
	}

	// A KeyValueExpr node represents (key : value) pairs
	// in composite literals.
	//
	KeyValueExpr struct {
		Key   Expr
		OpPos token.Pos // position of ":"
		Value Expr
	}
)

// A type is represented by a tree consisting of one
// or more of the following type-specific expression
// nodes.
//
type (
	// An ArrayType node represents an array or slice type.
	ArrayType struct {
		Lbrack token.Pos // position of "["
		Len    Expr      // Ellipsis node for [...]T array types, nil for slice types
		Elt    Expr      // element type
	}

	// A FuncType node represents a function type.
	FuncType struct {
		Func    token.Pos  // position of "func" keyword (token.NoPos if there is no "func")
		Params  *FieldList // (incoming) parameters; non-nil
		Results *FieldList // (outgoing) results; or nil*/
	}
)

// Pos and End implementations for expression/type nodes.

func (x *BadExpr) Pos() token.Pos    { return x.From }
func (x *Ident) Pos() token.Pos      { return x.NamePos }
func (x *Ellipsis) Pos() token.Pos   { return x.ValuePos }
func (x *BasicLit) Pos() token.Pos   { return x.ValuePos }
func (x *FuncLit) Pos() token.Pos    { return x.Type.Pos() }
func (x *CompositeLit) Pos() token.Pos {
	if x.Type != nil {
		return x.Type.Pos()
	}
	return x.Left
}
func (x *ParenExpr) Pos() token.Pos      { return x.Left }
func (x *QuotedExpr) Pos() token.Pos     { return x.Left }
func (x *EvalExpr) Pos() token.Pos       { return x.Left }
func (x *SelectorExpr) Pos() token.Pos   { return x.X.Pos() }
func (x *IndexExpr) Pos() token.Pos      { return x.Array.Pos() }
func (x *ListIndexExpr) Pos() token.Pos  { return x.Array.Pos() }
func (x *CallExpr) Pos() token.Pos       { return x.Fun.Pos() }
func (x *ArbitraryCallExpr) Pos() token.Pos { return x.Fun.Pos() }
func (x *UnaryExpr) Pos() token.Pos      { return x.OpPos }
func (x *BinaryExpr) Pos() token.Pos     { return x.X.Pos() }
func (x *TaggedExpr) Pos() token.Pos     { return x.X.Pos() }
func (x *KeyValueExpr) Pos() token.Pos   { return x.Key.Pos() }
func (x *ArrayType) Pos() token.Pos      { return x.Lbrack }
func (x *FuncType) Pos() token.Pos {
	if x.Func.IsValid() || x.Params == nil { // see issue 3870
		return x.Func
	}
	return x.Params.Pos() // interface method declarations have no "func" keyword
}
func (x *BadExpr) End() token.Pos        { return x.To }
func (x *Ident) End() token.Pos          { return token.Pos(int(x.NamePos) + len(x.Name)) }
func (x *Ellipsis) End() token.Pos       { return x.ValuePos + 2 }
func (x *BasicLit) End() token.Pos       { return token.Pos(int(x.ValuePos) + len(x.Value)) }
func (x *FuncLit) End() token.Pos        { return x.Body.End() }
func (x *CompositeLit) End() token.Pos   { return x.Right + 1 }
func (x *ParenExpr) End() token.Pos      { return x.Right + 1 }
func (x *QuotedExpr) End() token.Pos     { return x.Right + 1 }
func (x *EvalExpr) End() token.Pos       { return x.Right + 1 }
func (x *SelectorExpr) End() token.Pos   { return x.Sel.End() }
func (x *IndexExpr) End() token.Pos      { return x.Right + 1 }
func (x *ListIndexExpr) End() token.Pos  { return x.Right + 1 }
func (x *CallExpr) End() token.Pos       { return x.Right + 1 }
func (x *ArbitraryCallExpr) End() token.Pos { return x.Right + 1 }
func (x *UnaryExpr) End() token.Pos      { return x.X.End() }
func (x *BinaryExpr) End() token.Pos     { return x.Y.End() }
func (x *TaggedExpr) End() token.Pos     { return x.Rhs.End() }
func (x *KeyValueExpr) End() token.Pos   { return x.Value.End() }
func (x *ArrayType) End() token.Pos      { return x.Elt.End() }
func (x *FuncType) End() token.Pos {
	/*	if x.Results != nil {
		return x.Results.End()
	}*/
	return x.Params.End()
}

// exprNode() ensures that only expression/type nodes can be
// assigned to an Expr.
//
func (*BadExpr) exprNode()        {}
func (*Ident) exprNode()          {}
func (*Ellipsis) exprNode()       {}
func (*BasicLit) exprNode()       {}
func (*FuncLit) exprNode()        {}
func (*CompositeLit) exprNode()   {}
func (*ParenExpr) exprNode()      {}
func (*QuotedExpr) exprNode()     {}
func (*EvalExpr) exprNode()       {}
func (*SelectorExpr) exprNode()   {}
func (*IndexExpr) exprNode()      {}
func (*ListIndexExpr) exprNode()  {}
func (*CallExpr) exprNode()       {}
func (*ArbitraryCallExpr) exprNode() {}
func (*UnaryExpr) exprNode()      {}
func (*BinaryExpr) exprNode()     {}
func (*TaggedExpr) exprNode()     {}
func (*KeyValueExpr) exprNode()   {}

func (*ArrayType) exprNode()  {}
func (*FuncType) exprNode()   {}

// ----------------------------------------------------------------------------
// Convenience functions for Idents

// NewIdent creates a new Ident without position.
// Useful for ASTs generated by code other than the Go parser.
//
func NewIdent(name string) *Ident { return &Ident{token.NoPos, name} }

// IsExported reports whether name is an exported Go symbol
// (that is, whether it begins with an upper-case letter).
//
func IsExported(name string) bool {
	ch, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(ch)
}

// IsExported reports whether id is an exported Go symbol
// (that is, whether it begins with an uppercase letter).
//
func (id *Ident) IsExported() bool { return IsExported(id.Name) }

func (id *Ident) String() string {
	if id != nil {
		return id.Name
	}
	return "<nil>"
}

// ----------------------------------------------------------------------------
// Statements

// A statement is represented by a tree consisting of one
// or more of the following concrete statement nodes.
//
type (
	// A BadStmt node is a placeholder for statements containing
	// syntax errors for which no correct statement nodes can be
	// created.
	//
	BadStmt struct {
		From, To token.Pos // position range of bad statement
	}

	EOFStmt struct {
		EOF token.Pos     // position of EOF
	}

	// An EmptyStmt node represents an empty statement.
	// The "position" of the empty statement is the position
	// of the immediately following (explicit or implicit) semicolon.
	//
	EmptyStmt struct {
		Semicolon token.Pos // position of following ";"
		Implicit  bool      // if set, ";" was omitted in the source
	}

	// An ExprStmt node represents a (stand-alone) expression
	// in a statement list.
	//
	ExprStmt struct {
		X Expr // expression
	}

	QuotedStmt struct {
		Left  token.Pos
		X     Stmt
		Right token.Pos // position of ")"
	}

	// A ReturnStmt node represents a return statement.
	ReturnStmt struct {
		Return token.Pos // position of "return" keyword
		Result Expr      // result expression; or nil
	}

	// A BlockStmt node represents a braced statement list.
	BlockStmt struct {
		Lbrace token.Pos // position of "{"
		List   []Stmt
		Rbrace token.Pos // position of "}"
	}

	IfStmt struct {
		Keyword token.Pos // position of keyword
		Cond    Expr      // condition
		Body    *BlockStmt
		Else    Stmt // else branch; or nil
	}
	WhileStmt struct {
		Keyword token.Pos // position of keyword
		Cond    Expr      // condition
		Body    *BlockStmt
	}
	RepeatStmt struct {
		Keyword token.Pos // position of keyword
		Body    *BlockStmt
	}
	BreakStmt struct {
		Keyword token.Pos // position of keyword
	}
	NextStmt struct {
		Keyword token.Pos // position of keyword
	}
	ForStmt struct {
		Keyword   token.Pos // position of keyword
		Parameter *Ident
		Iterable  Expr
		Body      *BlockStmt
	}
	VersionStmt struct {
		Keyword token.Pos // position of keyword
	}
)

// Pos and End implementations for statement nodes.

func (s *BadStmt) Pos() token.Pos     { return s.From }
func (s *EOFStmt) Pos() token.Pos     { return s.EOF }
func (s *EmptyStmt) Pos() token.Pos   { return s.Semicolon }
func (s *ExprStmt) Pos() token.Pos    { return s.X.Pos() }
func (s *QuotedStmt) Pos() token.Pos  { return s.Left }
func (s *ReturnStmt) Pos() token.Pos  { return s.Return }
func (s *BlockStmt) Pos() token.Pos   { return s.Lbrace }
func (s *IfStmt) Pos() token.Pos      { return s.Keyword }
func (s *WhileStmt) Pos() token.Pos   { return s.Keyword }
func (s *RepeatStmt) Pos() token.Pos  { return s.Keyword }
func (s *BreakStmt) Pos() token.Pos   { return s.Keyword }
func (s *NextStmt) Pos() token.Pos    { return s.Keyword }
func (s *ForStmt) Pos() token.Pos     { return s.Keyword }
func (s *VersionStmt) Pos() token.Pos { return s.Keyword }

func (s *BadStmt) End() token.Pos { return s.To }
func (s *EOFStmt) End() token.Pos { return s.EOF }
func (s *EmptyStmt) End() token.Pos {
	if s.Implicit {
		return s.Semicolon
	}
	return s.Semicolon + 1 /* len(";") */
}
func (s *ExprStmt) End() token.Pos   { return s.X.End() }
func (s *QuotedStmt) End() token.Pos { return s.Right }
func (s *ReturnStmt) End() token.Pos {
	if s.Result != nil {
		return s.Result.End()
	}
	return s.Return + 6 // len("return")
}

func (s *BlockStmt) End() token.Pos { return s.Rbrace + 1 }
func (s *IfStmt) End() token.Pos {
	if s.Else != nil {
		return s.Else.End()
	}
	return s.Body.End()
}
func (s *WhileStmt) End() token.Pos   { return s.Body.End() }
func (s *RepeatStmt) End() token.Pos  { return s.Body.End() }
func (s *NextStmt) End() token.Pos    { return s.Keyword + 4 }
func (s *BreakStmt) End() token.Pos   { return s.Keyword + 5 }
func (s *VersionStmt) End() token.Pos { return s.Keyword + 7 }
func (s *ForStmt) End() token.Pos     { return s.Body.End() }

// stmtNode() ensures that only statement nodes can be
// assigned to a Stmt.
//
func (*BadStmt) stmtNode()     {}
func (*EOFStmt) stmtNode()     {}
func (*EmptyStmt) stmtNode()   {}
func (*ExprStmt) stmtNode()    {}
func (*QuotedStmt) stmtNode()  {}
func (*ReturnStmt) stmtNode()  {}
func (*BlockStmt) stmtNode()   {}
func (*IfStmt) stmtNode()      {}
func (*WhileStmt) stmtNode()   {}
func (*RepeatStmt) stmtNode()  {}
func (*BreakStmt) stmtNode()   {}
func (*NextStmt) stmtNode()    {}
func (*ForStmt) stmtNode()     {}
func (*VersionStmt) stmtNode() {}

// ----------------------------------------------------------------------------
// Declarations

// A Spec node represents a single (non-parenthesized) import,
// constant, type, or variable declaration.
//
type (
	// The Spec type stands for any of *ImportSpec, *ValueSpec, and *TypeSpec.
	Spec interface {
		Node
		specNode()
	}

	// An ImportSpec node represents a single package import.
	ImportSpec struct {
		Name   *Ident    // local package name (including "."); or nil
		Path   *BasicLit // import path
		EndPos token.Pos // end of spec (overrides Path.Pos if nonzero)
	}
)

// Pos and End implementations for spec nodes.

func (s *ImportSpec) Pos() token.Pos {
	if s.Name != nil {
		return s.Name.Pos()
	}
	return s.Path.Pos()
}
func (s *ImportSpec) End() token.Pos {
	if s.EndPos != 0 {
		return s.EndPos
	}
	return s.Path.End()
}

// specNode() ensures that only spec nodes can be
// assigned to a Spec.
//
func (*ImportSpec) specNode() {}

// ----------------------------------------------------------------------------
// Files and packages

// A File node represents a Go source file.
//
// The Comments list contains all comments in the source file in order of
// appearance, including the comments that are pointed to from other nodes
// via Doc and Comment fields.
//
type File struct {
	Package    token.Pos     // position of "package" keyword
	Name       *Ident        // package name
	Imports    []*ImportSpec // imports in this file
	Unresolved []*Ident      // unresolved identifiers in this file
}

func (f *File) Pos() token.Pos { return f.Package }
func (f *File) End() token.Pos { return f.Name.End() }
