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
	"roq/lib/ast"
	"roq/lib/token"
	"strconv"
)

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
	r := ast.Ident{NamePos: pos, Name: name}
	return &r
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


func (p *Parser) parseFuncParameterList(ellipsisOk bool) (list []*ast.Field) {
	if p.trace {
		defer un(trace(p, "FuncParameterList: "+p.lit))
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

func (p *Parser) parseFuncParameters(ellipsisOk bool) *ast.FieldList {

	if p.trace {
		defer un(trace(p, "FuncParameters: "+p.lit))
	}

	var params []*ast.Field
	lparen := p.expect(token.LPAREN)

	if p.tok != token.RPAREN {
		params = p.parseFuncParameterList(ellipsisOk)
	}
	rparen := p.expect(token.RPAREN)

	return &ast.FieldList{Opening: lparen, List: params, Closing: rparen}
}

func (p *Parser) parseFuncType() (*ast.FuncType) {
	if p.trace {
		defer un(trace(p, "FuncType"))
	}

	pos := p.expect(token.FUNCTION)
	params := p.parseFuncParameters(true)

	return &ast.FuncType{Func: pos, Params: params}
}


// ----------------------------------------------------------------------------
// Expressions

func (p *Parser) parseFuncLit() ast.Expr {
	if p.trace {
		defer un(trace(p, "FuncLit"))
	}

	typ := p.parseFuncType()
	//if p.tok != token.LBRACE {
	//// function type only
	//return typ
	//}

	p.exprLev++
	body := p.parseBody()
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
			return x
		case token.ELLIPSIS:
			panic("ELLIPSIS parsed")
		case token.LPAREN:
			lparen := p.pos
			p.next()
			p.exprLev++
			x := p.parseRhs()
			p.exprLev--
			rparen := p.expect(token.RPAREN)
			return &ast.ParenExpr{Left: lparen, X: x, Right: rparen}
		case token.FUNCTION:
			return p.parseFuncLit()
		}
	}

	// we have an error
	pos := p.pos
	p.errorExpected(pos, "operand")
	return &ast.BadExpr{From: pos, To: p.pos}
}

//   Quote eval
func (p *Parser) parseQuoteExpr() *ast.QuoteExpr {
	if p.trace {
		defer un(trace(p, "QuoteExpr"))
	}
	pos := p.pos
	p.expect(token.QUOTE)
	p.expect(token.LPAREN)

	var x ast.Stmt
	if p.tok != token.RPAREN {
		x = p.parseAssignment()
	} else {
		x = nil
	}
	rparen := p.expect(token.RPAREN)
	return &ast.QuoteExpr{Left: pos, X: x, Right: rparen}
}

func (p *Parser) parseEvalExpr() *ast.EvalExpr {
	if p.trace {
		defer un(trace(p, "EvalExpr"))
	}
	pos := p.pos
	p.expect(token.EVAL)
	p.expect(token.LPAREN)

	var x ast.Expr
	if p.tok != token.RPAREN {
		x = p.parseRhs()
	} else {
		x = nil
	}
	rparen := p.expect(token.RPAREN)
	return &ast.EvalExpr{Left: pos, X: x, Right: rparen}
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

	return &ast.IndexExpr{Array: x, Left: lbrack, Index: index, Right: rbrack}
}

func (p *Parser) parseListIndex(x ast.Expr) ast.Expr {
	if p.trace {
		defer un(trace(p, "ListIndex"))
	}

	dlbrack := p.expect(token.DOUBLELBRACK)
	var index ast.Expr
	if p.tok != token.SEQUENCE {  // TODO comma for dims
		index = p.parseRhs()
	}
	p.exprLev--
	drbrack := p.expect(token.DOUBLERBRACK)

	return &ast.ListIndexExpr{Array: x, Left: dlbrack, Index: index, Right: drbrack}
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

	return &ast.CallExpr{Fun: fun, Left: lparen, Args: list, Ellipsis: ellipsis, Right: rparen}
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
		defer un(trace(p, "Value"))
	}

	if p.tok == token.LBRACE {
		return p.parseLiteralValue(nil)
	}
	x := p.parseExpr(keyOk)
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
	return &ast.CompositeLit{Type: typ, Left: lbrace, Elts: elts, Right: rbrace}
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
//
// A primary expression is
// - literal
// - variable
// - a variable with selector or selectors
// - function call with parentheses
// - variable with an index
// - variable with a listindex
// - FIXME c(1,2,3)[1]

func (p *Parser) parsePrimaryExpr(lhs bool) ast.Expr {
	if p.trace && p.debug {
		if lhs {
			defer un(trace(p, "PrimaryExpr (LHS)"))
		} else {
			defer un(trace(p, "PrimaryExpr (RHS)"))
		}
	}
	
	x := p.parseOperand(lhs)
L:
	for {
		switch p.tok {
		case token.NA: //TODO test for method selectors
			p.next()
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
			x = p.parseIndex(x)
		case token.DOUBLELBRACK:
			x = p.parseListIndex(x)
		case token.LPAREN:
			x = p.parseCall(x)
		case token.LBRACE:
			if isLiteralType(x) && (p.exprLev >= 0) {
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
		if lhs {
			defer un(trace(p, "UnaryExpr (LHS)"))
		} else {
			defer un(trace(p, "UnaryExpr (RHS)"))
		}
	}

	switch p.tok {
	case token.PLUS, token.MINUS, token.NOT:
		pos, op := p.pos, p.tok
		p.next()
		x := p.parseUnaryExpr(lhs)
		return &ast.UnaryExpr{OpPos: pos, Op: op, X: x}
	case token.ELLIPSIS:
		pos := p.pos
		p.next()
		return &ast.Ellipsis{ValuePos: pos}
	case token.QUOTE:
		return p.parseQuoteExpr()
	case token.EVAL:
		return p.parseEvalExpr()
	default:
		return p.parsePrimaryExpr(lhs)
	}
}




func (p *Parser) tokPrec() (token.Token, int) {
	tok := p.tok
	return tok, tok.Precedence()
}

// If lhs is set and the result is an identifier, it is not resolved.
func (p *Parser) parseBinaryExpr(lhs bool, prec1 int) ast.Expr {
	if p.trace {
		if lhs {
			defer un(trace(p, "Binary or UnaryExpr (LHS)  precedence:"+strconv.Itoa(prec1)))
		} else {
			defer un(trace(p, "Binary or UnaryExpr (RHS) precedence:"+strconv.Itoa(prec1)))
		}
	}
	r := p.parseUnaryExpr(lhs)
	//	for p.tok != token.RPAREN && p.tok != token.COMMA && p.tok != token.SEMICOLON && p.tok != token.ELSE{
	for p.tok.IsOperator() {
		operator, oprec := p.tokPrec()
		if oprec < prec1 {
			return r
		}
		pos := p.expect(operator)
		if lhs {
		}
		y := p.parseBinaryExpr(false, oprec+1)
		r = &ast.BinaryExpr{X: r, OpPos: pos, Op: operator, Y: y}
	}
	return r
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
			lhs := e.X.(*ast.Ident) // TODO check for ident
			return &ast.TaggedExpr{X: lhs, Tag: lhs.Name, OpPos: e.OpPos, Rhs: e.Y}
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
		if lhs {
			defer un(trace(p, "Expr (LHS)"))
		} else {
			defer un(trace(p, "Expr (RHS)"))
		}
	}
	return p.parseBinaryExpr(lhs, 3)
}


func (p *Parser) parseRhs() ast.Expr {
	old := p.inRhs
	p.inRhs = true
	x := p.parseBinaryExpr(false, 3)
	p.inRhs = old
	return x
}


