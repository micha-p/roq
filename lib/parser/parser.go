// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package parser implements a parser for Go source files. Input may be
// provided in a variety of forms (see the various Parse* functions); the
// output is an abstract syntax tree (AST) representing the Go source. The
// parser is invoked through one of the Parse* functions.
//
// The parser accepts a larger language than is syntactically permitted by
// the Go spec, for simplicity, and for improved robustness in the presence
// of syntax errors. For instance, in method declarations, the receiver is
// treated like an ordinary parameter list and thus may contain multiple
// entries where the spec permits exactly one. Consequently, the corresponding
// field in the AST (ast.FuncDecl.Recv) field is not restricted to one entry.
//
package parser

import (
	"fmt"
	"roq/lib/ast"
	"roq/lib/scanner"
	"roq/lib/token"
	"strconv"
)

// The parser structure holds the parser's internal state.
type Parser struct {
	file    *token.File
	errors  scanner.ErrorList
	scanner scanner.Scanner

	// Tracing/debugging
	trace  bool // == (mode & Trace != 0)
	debug  bool // == (mode & Trace != 0)
	echo   bool // == (mode & Echo  != 0)
	indent int  // indentation used for tracing output

	// Next token
	pos token.Pos   // token position
	tok token.Token // one token look-ahead
	lit string      // token literal

	// Error recovery
	// (used to limit the number of calls to syncXXX functions
	// w/o making scanning progress - avoids potential endless
	// loops across multiple parser functions during error recovery)
	syncPos token.Pos // last synchronization position
	syncCnt int       // number of calls to syncXXX without progress

	// Non-syntactic parser control
	exprLev int  // < 0: in control clause, >= 0: in expression
	inRhs   bool // if set, the parser is parsing a rhs expression

	// Identifier scopes
	pkgScope   *ast.Scope        // pkgScope.Outer == nil
	topScope   *ast.Scope        // top-most scope; may be pkgScope
	unresolved []*ast.Ident      // unresolved identifiers
	imports    []*ast.ImportSpec // list of imports
}

func (p *Parser) init(fset *token.FileSet, filename string, src []byte, mode Mode) {
	p.file = fset.AddFile(filename, -1, len(src))

	p.trace = mode&Trace != 0 // for convenience (p.trace is used frequently)
	p.debug = mode&Debug != 0 // for convenience (p.debug is used frequently)
	p.echo  = mode&Echo  != 0

	eh := func(pos token.Position, msg string) { p.errors.Add(pos, msg) }
	p.scanner.Init(p.file, src, eh, p.echo)


	p.next()
}

// ----------------------------------------------------------------------------
// Scoping support

func (p *Parser) openScope() {
	p.topScope = ast.NewScope(p.topScope)
}

func (p *Parser) closeScope() {
	p.topScope = p.topScope.Outer
}

func (p *Parser) declare(decl, data interface{}, scope *ast.Scope, kind ast.ObjKind, idents ...*ast.Ident) {
	for _, ident := range idents {
		assert(ident.Obj == nil, "identifier already declared or resolved")
		obj := ast.NewObj(kind, ident.Name)
		// remember the corresponding declaration for redeclaration
		// errors and global variable resolution/typechecking phase
		obj.Decl = decl
		obj.Data = data
		ident.Obj = obj
		if ident.Name != "_" {
			if alt := scope.Insert(obj); alt != nil {
				prevDecl := ""
				if pos := alt.Pos(); pos.IsValid() {
					prevDecl = fmt.Sprintf("\n\tprevious declaration at %s", p.file.Position(pos))
				}
				p.error(ident.Pos(), fmt.Sprintf("%s redeclared in this block%s", ident.Name, prevDecl))
			}
		}
	}
}

// The unresolved object is a sentinel to mark identifiers that have been added
// to the list of unresolved identifiers. The sentinel is only used for verifying
// internal consistency.
var unresolved = new(ast.Object)

// If x is an identifier, tryResolve attempts to resolve x by looking up
// the object it denotes. If no object is found and collectUnresolved is
// set, x is marked as unresolved and collected in the list of unresolved
// identifiers.
//
func (p *Parser) tryResolve(x ast.Expr, collectUnresolved bool) {
	// nothing to do if x is not an identifier or the blank identifier
	ident, _ := x.(*ast.Ident)
	if ident == nil {
		return
	}
	assert(ident.Obj == nil, "identifier already declared or resolved")
	if ident.Name == "_" {
		return
	}
	// try to resolve the identifier
	for s := p.topScope; s != nil; s = s.Outer {
		if obj := s.Lookup(ident.Name); obj != nil {
			ident.Obj = obj
			return
		}
	}
	// all local scopes are known, so any unresolved identifier
	// must be found either in the file scope, package scope
	// (perhaps in another file), or universe scope --- collect
	// them so that they can be resolved later
	if collectUnresolved {
		ident.Obj = unresolved
		p.unresolved = append(p.unresolved, ident)
	}
}

func (p *Parser) resolve(x ast.Expr) {
	p.tryResolve(x, true)
}

// ----------------------------------------------------------------------------
// Parsing support

func (p *Parser) printTrace(a ...interface{}) {
	const dots = ". . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . "
	const n = len(dots)
	pos := p.file.Position(p.pos)
	fmt.Printf("%5d:%3d: ", pos.Line, pos.Column)
	i := 2 * p.indent
	for i > n {
		fmt.Print(dots)
		i -= n
	}
	// i <= n
	fmt.Print(dots[0:i])
	fmt.Println(a...)
}

func trace0(p *Parser, msg string) *Parser {
	p.printTrace(msg)
	return p
}

func trace(p *Parser, msg string) *Parser {
	p.printTrace(msg, "(")
	p.indent++
	return p
}

// Usage pattern: defer un(trace(p, "..."))
func un(p *Parser) {
	p.indent--
	p.printTrace(")")
}

// Advance to the next token.
func (p *Parser) next() {
	// Because of one-token look-ahead, print the previous token
	// when tracing as it provides a more readable output. The
	// very first token (!p.pos.IsValid()) is not initialized
	// (it is token.ILLEGAL), so don't print it .
	if p.trace && p.pos.IsValid() {
		s := p.tok.String()
		switch {
		case p.tok.IsLiteral():
			p.printTrace(s, p.lit)
		case p.tok.IsOperator(), p.tok.IsKeyword():
			p.printTrace("\"" + s + "\"")
		default:
			p.printTrace(s)
		}
	}

	p.pos, p.tok, p.lit = p.scanner.Scan()
}

// A bailout panic is raised to indicate early termination.
type bailout struct{}

func (p *Parser) error(pos token.Pos, msg string) {
	epos := p.file.Position(pos)
	p.errors.Add(epos, msg)
}

func (p *Parser) errorExpected(pos token.Pos, msg string) {
	msg = "expected " + msg
	if pos == p.pos {
		// the error happened at the current position;
		// make the error message more specific
		if p.tok == token.SEMICOLON && p.lit == "\n" {
			msg += ", found newline"
		} else {
			msg += ", found '" + p.tok.String() + "'"
			if p.tok.IsLiteral() {
				msg += " " + p.lit
			}
		}
	}
	p.error(pos, msg)
}

func (p *Parser) expect(tok token.Token) token.Pos {
	if p.debug {
		defer trace0(p, "EXPECT: "+tok.String()+"  GIVEN: "+p.tok.String())
	}

	pos := p.pos
	if p.tok != tok {
		p.errorExpected(pos, "'"+tok.String()+"'")
	}
	p.next() // make progress
	return pos
}

// expectClosing is like expect but provides a better error message
// for the common case of a missing comma before a newline.
//
func (p *Parser) expectClosing(tok token.Token, context string) token.Pos {
	if p.tok != tok && p.tok == token.SEMICOLON && p.lit == "\n" {
		p.error(p.pos, "missing ',' before newline in "+context)
		p.next()
	}
	return p.expect(tok)
}

func (p *Parser) expectSemi() {
	// semicolon is optional before a closing ')' or '}'
	if p.tok != token.RPAREN && p.tok != token.RBRACE {
		switch p.tok {
		case token.COMMA:
			// permit a ',' instead of a ';' but complain
			p.errorExpected(p.pos, "';'")
			fallthrough
		case token.SEMICOLON:
			p.next()
		default:
			p.errorExpected(p.pos, "';'")
			syncStmt(p)
		}
	}
}

func (p *Parser) atComma(context string, follow token.Token) bool {
	if p.tok == token.COMMA {
		return true
	}
	if p.tok != follow {
		msg := "missing ','"
		if p.tok == token.SEMICOLON && p.lit == "\n" {
			msg += " before newline"
		}
		p.error(p.pos, msg+" in "+context)
		return true // "insert" comma and continue
	}
	return false
}

func assert(cond bool, msg string) {
	if !cond {
		panic("roq/lib/parser internal error: " + msg)
	}
}

// syncStmt advances to the next statement.
// Used for synchronization after an error.
//
func syncStmt(p *Parser) {
	for {
		switch p.tok {
		case token.BREAK, token.CONST, token.NEXT, token.FOR,
			token.IF, token.RETURN, token.SELECT, token.SWITCH,
			token.TYPE, token.VAR:
			// Return only if parser made some progress since last
			// sync or if it has not reached 10 sync calls without
			// progress. Otherwise consume at least one token to
			// avoid an endless parser loop (it is possible that
			// both parseOperand and parseStmt call syncStmt and
			// correctly do not advance, thus the need for the
			// invocation limit p.syncCnt).
			if p.pos == p.syncPos && p.syncCnt < 10 {
				p.syncCnt++
				return
			}
			if p.pos > p.syncPos {
				p.syncPos = p.pos
				p.syncCnt = 0
				return
			}
			// Reaching here indicates a parser bug, likely an
			// incorrect token list in this function, but it only
			// leads to skipping of possibly correct code if a
			// previous error is present, and thus is preferred
			// over a non-terminating parse.
		case token.EOF:
			return
		}
		p.next()
	}
}

// syncDecl advances to the next declaration.
// Used for synchronization after an error.
//
func syncDecl(p *Parser) {
	for {
		switch p.tok {
		case token.CONST, token.TYPE, token.VAR:
			// see comments in syncStmt
			if p.pos == p.syncPos && p.syncCnt < 10 {
				p.syncCnt++
				return
			}
			if p.pos > p.syncPos {
				p.syncPos = p.pos
				p.syncCnt = 0
				return
			}
		case token.EOF:
			return
		}
		p.next()
	}
}

// safePos returns a valid file position for a given position: If pos
// is valid to begin with, safePos returns pos. If pos is out-of-range,
// safePos returns the EOF position.
//
// This is hack to work around "artificial" end positions in the AST which
// are computed by adding 1 to (presumably valid) token positions. If the
// token positions are invalid due to parse errors, the resulting end position
// may be past the file's EOF position, which would lead to panics if used
// later on.
//
func (p *Parser) safePos(pos token.Pos) (res token.Pos) {
	defer func() {
		if recover() != nil {
			res = token.Pos(p.file.Base() + p.file.Size()) // EOF position
		}
	}()
	_ = p.file.Offset(pos) // trigger a panic if position is out-of-range
	return pos
}

// ----------------------------------------------------------------------------
// Identifiers

func (p *Parser) parseIdent() *ast.Ident {
	pos := p.pos
	name := "_"
	if p.tok == token.IDENT {
		name = p.lit
		p.next()
	} else {
		p.expect(token.IDENT) // use expect() error handling
	}
	return &ast.Ident{NamePos: pos, Name: name}
}

func (p *Parser) parseIdentList() (list []*ast.Ident) {
	if p.trace {
		defer un(trace(p, "IdentList"))
	}

	list = append(list, p.parseIdent())
	for p.tok == token.COMMA {
		p.next()
		list = append(list, p.parseIdent())
	}

	return
}

// ----------------------------------------------------------------------------
// Common productions

// If lhs is set, result list elements which are identifiers are not resolved.
func (p *Parser) parseExprList(lhs bool) (list []ast.Expr) {
	if p.trace {
		defer un(trace(p, "ExpressionList"))
	}

	list = append(list, p.parseExpr(lhs))
	for p.tok == token.COMMA {
		p.next()
		list = append(list, p.parseExpr(lhs))
	}

	return
}

// ----------------------------------------------------------------------------
// Types

func (p *Parser) parseType() ast.Expr {
	if p.trace {
		defer un(trace(p, "Type"))
	}

	typ := p.tryType()

	if typ == nil {
		pos := p.pos
		p.errorExpected(pos, "type")
		p.next() // make progress
		return &ast.BadExpr{From: pos, To: p.pos}
	}

	return typ
}

// If the result is an identifier, it is not resolved.
func (p *Parser) parseTypeName() ast.Expr {
	if p.trace {
		defer un(trace(p, "TypeName"))
	}

	ident := p.parseIdent()
	// don't resolve ident yet - it may be a parameter or field name

	if p.tok == token.NA /*TODO: CHECK*/ {
		// ident is a package name
		p.next()
		p.resolve(ident)
		sel := p.parseIdent()
		return &ast.SelectorExpr{X: ident, Sel: sel}
	}

	return ident
}


func (p *Parser) parseFuncParameterList(scope *ast.Scope, ellipsisOk bool) (list []*ast.Field) {
	if p.trace {
		defer un(trace(p, "ParameterList: "+p.lit))
	}

	for {
		if p.tok == token.RPAREN {
			break
		}
		//		list = append(list, p.parseVarType(ellipsisOk))
		if p.tok == token.ELLIPSIS {
			list = append(list, &ast.Field{Type: &ast.Ellipsis{}})
			p.next()
		} else {
			identifier := p.parseIdent()
			if p.tok == token.SHORTASSIGNMENT {
				p.next()
				list = append(list, &ast.Field{Type: identifier, Default: p.parseRhs()})
			}else{
				list = append(list, &ast.Field{Type: identifier})
			}
		}
		if p.tok != token.COMMA {
			break
		}
		p.next()
		if p.tok == token.RPAREN {
			break
		}
	}
	return
}

func (p *Parser) parseFuncParameters(scope *ast.Scope, ellipsisOk bool) *ast.FieldList {

	if p.trace {
		defer un(trace(p, "Parameters: "+p.lit))
	}

	var params []*ast.Field
	lparen := p.expect(token.LPAREN)

	if p.tok != token.RPAREN {
		params = p.parseFuncParameterList(scope, ellipsisOk)
	}
	rparen := p.expect(token.RPAREN)

	return &ast.FieldList{Opening: lparen, List: params, Closing: rparen}
}

func (p *Parser) parseFuncType() (*ast.FuncType, *ast.Scope) {
	if p.trace {
		defer un(trace(p, "FuncType"))
	}

	pos := p.expect(token.FUNCTION)
	scope := ast.NewScope(p.topScope) // function scope
	params := p.parseFuncParameters(scope, true)

	return &ast.FuncType{Func: pos, Params: params}, scope
}

// If the result is an identifier, it is not resolved.
func (p *Parser) tryIdentOrType() ast.Expr {
	switch p.tok {
	case token.IDENT:
		return p.parseTypeName()
	case token.FUNCTION:
		typ, _ := p.parseFuncType()
		return typ
		/*	case token.INTERFACE:
			return p.parseInterfaceType()*/
	case token.LPAREN:
		lparen := p.pos
		p.next()
		typ := p.parseType()
		rparen := p.expect(token.RPAREN)
		return &ast.ParenExpr{Lparen: lparen, X: typ, Rparen: rparen}
	}

	// no type found
	return nil
}

func (p *Parser) tryType() ast.Expr {
	typ := p.tryIdentOrType()
	if typ != nil {
		p.resolve(typ)
	}
	return typ
}


// ----------------------------------------------------------------------------
// Expressions

func (p *Parser) parseFuncLit() ast.Expr {
	if p.trace {
		defer un(trace(p, "FuncLit"))
	}

	typ, scope := p.parseFuncType()
	//if p.tok != token.LBRACE {
	//// function type only
	//return typ
	//}

	p.exprLev++
	body := p.parseBody(scope)
	p.exprLev--

	return &ast.FuncLit{Type: typ, Body: body}
}

// parseOperand may return an expression or a raw type (incl. array
// types of the form [...]T. Callers must verify the result.
// If lhs is set and the result is an identifier, it is not resolved.
//
func (p *Parser) parseOperand(lhs bool) ast.Expr {
	if p.trace {
		defer un(trace(p, "Operand"))
	}

	if p.tok.IsLiteral() || p.tok.IsConstant() {
		x := &ast.BasicLit{ValuePos: p.pos, Kind: p.tok, Value: p.lit}
		p.next()
		return x
	} else {
		switch p.tok {
		case token.IDENT:
			x := p.parseIdent()
			if !lhs {
				p.resolve(x)
			}
			return x
		case token.LPAREN:
			lparen := p.pos
			p.next()
			p.exprLev++
			x := p.parseRhs()
			p.exprLev--
			rparen := p.expect(token.RPAREN)
			return &ast.ParenExpr{Lparen: lparen, X: x, Rparen: rparen}
		case token.FUNCTION:
			return p.parseFuncLit()
		}
	}
	if typ := p.tryIdentOrType(); typ != nil {
		// could be type for composite literal or conversion
		_, isIdent := typ.(*ast.Ident)
		assert(!isIdent, "type cannot be identifier")
		return typ
	}

	// we have an error
	pos := p.pos
	p.errorExpected(pos, "operand")
	syncStmt(p)
	return &ast.BadExpr{From: pos, To: p.pos}
}

func (p *Parser) parseSelector(x ast.Expr) ast.Expr {
	if p.trace {
		defer un(trace(p, "Selector"))
	}

	sel := p.parseIdent()

	return &ast.SelectorExpr{X: x, Sel: sel}
}

func (p *Parser) parseIndex(x ast.Expr) ast.Expr {
	if p.trace {
		defer un(trace(p, "IndexOrSlice"))
	}

	lbrack := p.expect(token.LBRACK)
	var index ast.Expr
	if p.tok != token.SEQUENCE {  // TODO comma for dims
		index = p.parseRhs()
	}
	p.exprLev--
	rbrack := p.expect(token.RBRACK)

	return &ast.IndexExpr{Array: x, Lbrack: lbrack, Index: index, Rbrack: rbrack}
}

func (p *Parser) parseCall(fun ast.Expr) *ast.CallExpr {
	if p.trace {
		defer un(trace(p, "Call"))
	}

	lparen := p.expect(token.LPAREN)
	p.exprLev++
	var list []ast.Expr
	var ellipsis token.Pos
	for p.tok != token.RPAREN && p.tok != token.EOF && !ellipsis.IsValid() {
		list = append(list, p.parseParameter())
		if p.tok == token.ELLIPSIS {
			ellipsis = p.pos
			p.next()
		}
		if !p.atComma("argument list", token.RPAREN) {
			break
		}
		p.next()
	}
	p.exprLev--
	rparen := p.expectClosing(token.RPAREN, "argument list")

	return &ast.CallExpr{Fun: fun, Lparen: lparen, Args: list, Ellipsis: ellipsis, Rparen: rparen}
}

// TODO vector literals consisting just of floats

//func (p *Parser) parseVector(start token.Pos) *ast.VectorLit {
	//if p.trace {
		//defer un(trace(p, "Vector"))
	//}
	//lparen := p.expect(token.LPAREN)
	//p.exprLev++
	//var list []ast.Expr
	//for p.tok != token.RPAREN && p.tok != token.EOF {
		//list = append(list, p.parseParameter())
		//if !p.atComma("vector element list", token.RPAREN) {
			//break
		//}
		//p.next()
	//}
	//p.exprLev--
	//rparen := p.expectClosing(token.RPAREN, "vector element list")
	//return &ast.VectorLit{Lparen: lparen, Args: list, Rparen: rparen}
//}

func (p *Parser) parseValue(keyOk bool) ast.Expr {
	if p.trace {
		defer un(trace(p, "Element"))
	}

	if p.tok == token.LBRACE {
		return p.parseLiteralValue(nil)
	}

	// Because the parser doesn't know the composite literal type, it cannot
	// know if a key that's an identifier is a struct field name or a name
	// denoting a value. The former is not resolved by the parser or the
	// resolver.
	//
	// Instead, _try_ to resolve such a key if possible. If it resolves,
	// it a) has correctly resolved, or b) incorrectly resolved because
	// the key is a struct field with a name matching another identifier.
	// In the former case we are done, and in the latter case we don't
	// care because the type checker will do a separate field lookup.
	//
	// If the key does not resolve, it a) must be defined at the top
	// level in another file of the same package, the universe scope, or be
	// undeclared; or b) it is a struct field. In the former case, the type
	// checker can do a top-level lookup, and in the latter case it will do
	// a separate field lookup.
	x := p.parseExpr(keyOk)
	if keyOk {
		/*if p.tok == token.COLON {
			// Try to resolve the key but don't collect it
			// as unresolved identifier if it fails so that
			// we don't get (possibly false) errors about
			// undeclared names.
			p.tryResolve(x, false)
		} else {
			// not a key */
		p.resolve(x)
		/*}*/
	}

	return x
}

func (p *Parser) parseElement() ast.Expr {
	if p.trace {
		defer un(trace(p, "Element"))
	}

	x := p.parseValue(true)
	/*if p.tok == token.COLON {
		colon := p.pos
		p.next()
		x = &ast.KeyValueExpr{Key: x, Colon: colon, Value: p.parseValue(false)}
	}*/

	return x
}

func (p *Parser) parseElementList() (list []ast.Expr) {
	if p.trace {
		defer un(trace(p, "ElementList"))
	}

	for p.tok != token.RBRACE && p.tok != token.EOF {
		list = append(list, p.parseElement())
		if !p.atComma("composite literal", token.RBRACE) {
			break
		}
		p.next()
	}

	return
}

func (p *Parser) parseLiteralValue(typ ast.Expr) ast.Expr {
	if p.trace {
		defer un(trace(p, "LiteralValue"))
	}

	lbrace := p.expect(token.LBRACE)
	var elts []ast.Expr
	p.exprLev++
	if p.tok != token.RBRACE {
		elts = p.parseElementList()
	}
	p.exprLev--
	rbrace := p.expectClosing(token.RBRACE, "composite literal")
	return &ast.CompositeLit{Type: typ, Lbrace: lbrace, Elts: elts, Rbrace: rbrace}
}

// isLiteralType reports whether x is a legal composite literal type.
func isLiteralType(x ast.Expr) bool {
	switch t := x.(type) {
	case *ast.BadExpr:
	case *ast.Ident:
	case *ast.SelectorExpr:
		_, isIdent := t.X.(*ast.Ident)
		return isIdent
	case *ast.ArrayType:
	default:
		return false // all other nodes are not legal composite literal types
	}
	return true
}

// If x is of the form (T), unparen returns unparen(T), otherwise it returns x.
func unparen(x ast.Expr) ast.Expr {
	if p, isParen := x.(*ast.ParenExpr); isParen {
		x = unparen(p.X)
	}
	return x
}

// If lhs is set and the result is an identifier, it is not resolved.
func (p *Parser) parsePrimaryExpr(lhs bool) ast.Expr {
	if p.trace && p.debug {
		defer un(trace(p, "PrimaryExpr"))
	}
	
	x := p.parseOperand(lhs)
L:
	for {
		switch p.tok {
		case token.NA: //TODO
			p.next()
			if lhs {
				p.resolve(x)
			}
			switch p.tok {
			case token.IDENT:
				x = p.parseSelector(x)
			default:
				pos := p.pos
				p.errorExpected(pos, "selector or type assertion")
				p.next() // make progress
				sel := &ast.Ident{NamePos: pos, Name: "_"}
				x = &ast.SelectorExpr{X: x, Sel: sel}
			}
		case token.LBRACK:
			if lhs {
				p.resolve(x)
			}
			x = p.parseIndex(x)
		case token.LPAREN:
			if lhs {
				p.resolve(x)
			}
			x = p.parseCall(x)
		case token.LBRACE:
			if isLiteralType(x) && (p.exprLev >= 0) {
				if lhs {
					p.resolve(x)
				}
				x = p.parseLiteralValue(x)
			} else {
				break L
			}
		default:
			break L
		}
		lhs = false // no need to try to resolve again
	}

	return x
}

// If lhs is set and the result is an identifier, it is not resolved.
func (p *Parser) parseUnaryExpr(lhs bool) ast.Expr {
	if p.trace && p.debug {
		defer un(trace(p, "UnaryExpr"))
	}

	switch p.tok {
	case token.PLUS, token.MINUS, token.NOT, token.AND:
		pos, op := p.pos, p.tok
		p.next()
		x := p.parseUnaryExpr(false)
		return &ast.UnaryExpr{OpPos: pos, Op: op, X: x}
	}

	return p.parsePrimaryExpr(lhs)
}

func (p *Parser) tokPrec() (token.Token, int) {
	tok := p.tok
	return tok, tok.Precedence()
}

// If lhs is set and the result is an identifier, it is not resolved.
func (p *Parser) parseBinaryExpr(lhs bool, prec1 int) ast.Expr {
	if p.trace {
		defer un(trace(p, "BinaryExpr "+strconv.Itoa(prec1)))

	}
	x := p.parseUnaryExpr(lhs)
	//	for p.tok != token.RPAREN && p.tok != token.COMMA && p.tok != token.SEMICOLON && p.tok != token.ELSE{
	for p.tok.IsOperator() {
		operator, oprec := p.tokPrec()
		if oprec < prec1 {
			return x
		}
		pos := p.expect(operator)
		if lhs {
			p.resolve(x)
			lhs = false
		}
		y := p.parseBinaryExpr(false, oprec+1)
		x = &ast.BinaryExpr{X: x, OpPos: pos, Op: operator, Y: y}
	}
	return x
}

func (p *Parser) parseParameter() ast.Expr {
	if p.trace {
		defer un(trace(p, "Parameter"))
	}
	x := p.parseBinaryExpr(false, 1)
	switch x.(type) {
	case *ast.BinaryExpr:
		e := x.(*ast.BinaryExpr)
		if e.Op == token.SHORTASSIGNMENT {
			lhs := e.X.(*ast.BasicLit) // TODO check for ident
			return &ast.TaggedExpr{X: lhs, Tag: lhs.Value, EqPos: e.OpPos, Rhs: e.Y}
		} else {
			return x
		}
	}
	return x
}

// If lhs is set and the result is an identifier, it is not resolved.
// The result may be a type or even a raw type ([...]int).
func (p *Parser) parseExpr(lhs bool) ast.Expr {
	if p.trace {
		defer un(trace(p, "Parsing expression or shortassignment"))
	}
	return p.parseBinaryExpr(lhs, 3)
}

func (p *Parser) parseExprOrAssignment(lhs bool) ast.Expr {
	if p.trace {
		defer un(trace(p, "Parsing expression or assignment"))
	}
	return p.parseBinaryExpr(lhs, 1)
}

func (p *Parser) parseRhs() ast.Expr {
	old := p.inRhs
	p.inRhs = true
	x := p.parseExprOrAssignment(false)
	p.inRhs = old
	return x
}

