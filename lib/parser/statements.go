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
	"roq/lib/token"
)


// ----------------------------------------------------------------------------
// Blocks

func (p *Parser) parseStmtList() (list []ast.Stmt) {
	if p.debug {
		defer un(trace(p, "StatementList"))
	}

	for p.tok != token.RBRACE && p.tok != token.EOF {
		list = append(list, p.parseStmt())
	}

	return
}

func (p *Parser) parseBody() *ast.BlockStmt {
	if p.trace {
		defer un(trace(p, "Body"))
	}

	var r *ast.BlockStmt
	if p.tok == token.LBRACE {
		r = p.parseBlockStmt()
	} else {
		r = p.parseBlockStmt1()
	}
	return r
}

func (p *Parser) parseBlockStmt() *ast.BlockStmt {
	if p.trace {
		defer un(trace(p, "BlockStmt"))
	}

	lbrace := p.expect(token.LBRACE)
	list := p.parseStmtList()
	rbrace := p.expect(token.RBRACE)

	return &ast.BlockStmt{Lbrace: lbrace, List: list, Rbrace: rbrace}
}

// return one single statement as list of length 1
func (p *Parser) parseBlockStmt1() *ast.BlockStmt {
	if p.trace {
		defer un(trace(p, "BlockStmt1"))
	}

	var list []ast.Stmt
	stmt := p.parseStmt()
	list = append(list, stmt)

	return &ast.BlockStmt{List: list}
}


// ----------------------------------------------------------------------------
// Statements

// Parsing modes for parseAssignment.
const (
	basic = iota
	labelOk
	rangeOk
)

// parseAssignment returns true as 2nd result if it parsed the assignment
// of a range clause (with mode == rangeOk). The returned statement is an
// assignment with a right-hand side that is a single unary expression of
// the form "range x". No guarantees are given for the left-hand side.
func (p *Parser) parseAssignment() ast.Stmt {
	if p.trace {
		defer un(trace(p, "Assignment (or expr)"))
	}

	x := p.parseExpr(true)
	var y ast.Expr
	pos, tok := p.pos, p.tok
	var s *ast.AssignStmt

	switch p.tok {
	case token.SHORTASSIGNMENT, token.LEFTASSIGNMENT, token.RIGHTASSIGNMENT, 
		 token.SUPERLEFTASSIGNMENT, token.SUPERRIGHTASSIGNMENT:
		p.next()
		y = p.parseRhs()
		s = &ast.AssignStmt{Lhs: x, TokPos: pos, Tok: tok, Rhs: y}
		return s
	default:
		e := &ast.ExprStmt{X: x}
		return e
	}
}

func (p *Parser) parseReturnStmt() *ast.ReturnStmt {
	if p.trace {
		defer un(trace(p, "ReturnStmt"))
	}

	pos := p.pos
	p.expect(token.RETURN)
	p.expect(token.LPAREN)
	var x ast.Expr
	if p.tok != token.RPAREN {
		x = p.parseRhs()
	} else {
		x = nil
	}
	p.expect(token.RPAREN)
	p.expectSemi()

	return &ast.ReturnStmt{Return: pos, Result: x}
}


func (p *Parser) makeExpr(s ast.Stmt, kind string) ast.Expr {
	if s == nil {
		return nil
	}
	if es, isExpr := s.(*ast.ExprStmt); isExpr {
		return es.X
	}
	p.error(s.Pos(), fmt.Sprintf("expected %s, found simple statement (missing parentheses around composite literal?)", kind))
	return &ast.BadExpr{From: s.Pos(), To: p.safePos(s.End())}
}

func (p *Parser) parseIfStmt() *ast.IfStmt {
	if p.trace {
		defer un(trace(p, "IfStmt"))
	}

	pos := p.expect(token.IF)

	var x ast.Expr // TODO strict flag to insist on parentheses
	x = p.parseRhs()

	var body *ast.BlockStmt
	if p.trace {
		defer un(trace(p, "bodyStmt"))
	}
	if p.tok == token.LBRACE {
		body = p.parseBlockStmt()
	} else {
		body = p.parseBlockStmt1()
	}

	var else_ ast.Stmt
	if p.tok == token.ELSE {
		if p.trace {
			defer un(trace(p, "elseStmt"))
		}
		p.next()
		switch p.tok {
		case token.IF:
			else_ = p.parseIfStmt()
		case token.LBRACE:
			else_ = p.parseBlockStmt()
			p.expectSemi()
		default:
			else_ = p.parseBlockStmt1()
			p.expectSemi()
		}
	} else {
		p.expectSemi()
	}

	return &ast.IfStmt{Keyword: pos, Cond: x, Body: body, Else: else_}
}

func (p *Parser) parseWhileStmt() *ast.WhileStmt {
	if p.trace {
		defer un(trace(p, "WhileStmt"))
	}

	pos := p.expect(token.WHILE)

	var x ast.Expr // TODO strict flag to insist on parentheses
	x = p.parseRhs()

	var body *ast.BlockStmt
	if p.trace {
		defer un(trace(p, "bodyStmt"))
	}
	if p.tok == token.LBRACE {
		body = p.parseBlockStmt()
	} else {
		body = p.parseBlockStmt1()
	}
	p.expectSemi()
	return &ast.WhileStmt{Keyword: pos, Cond: x, Body: body}
}

func (p *Parser) parseRepeatStmt() *ast.RepeatStmt {
	if p.trace {
		defer un(trace(p, "RepeatStmt"))
	}

	pos := p.expect(token.REPEAT)

	var body *ast.BlockStmt
	if p.trace {
		defer un(trace(p, "bodyStmt"))
	}
	if p.tok == token.LBRACE {
		body = p.parseBlockStmt()
	} else {
		body = p.parseBlockStmt1()
	}
	p.expectSemi()
	return &ast.RepeatStmt{Keyword: pos, Body: body}
}

func (p *Parser) parseForStmt() ast.Stmt {
	if p.trace {
		defer un(trace(p, "ForStmt"))
	}

	pos := p.expect(token.FOR)

	// TODO this is quick and dirty
	p.expect(token.LPAREN)
	id := p.parseIdent()
	p.expect(token.IN)
	vec := p.parseRhs()
	p.expect(token.RPAREN)

	var body *ast.BlockStmt
	if p.tok == token.LBRACE {
		body = p.parseBlockStmt()
	} else {
		body = p.parseBlockStmt1()
	}
	p.expectSemi()

	return &ast.ForStmt{
		Keyword:   pos,
		Parameter: id,
		Iterable:  vec,
		Body:      body,
	}
}

func (p *Parser) parseStmt() (s ast.Stmt) {
	if p.trace {
		defer un(trace(p, "Statement"))
	}

	switch p.tok {
	case
		// tokens that may start an expression
		token.IDENT, token.INT, token.FLOAT, token.IMAG, token.STRING, token.FUNCTION, token.LPAREN, // operands
		token.NULL, token.NA, token.INF, token.NAN, token.TRUE, token.FALSE, // constants
		token.PLUS, token.MINUS, token.NOT, // unary operators
		token.LBRACK,
		token.QUOTE, token.EVAL, token.CALL:
		s = p.parseAssignment() // this parses an assignment or an expression!

	case token.IF:
		s = p.parseIfStmt()
	case token.FOR:
		s = p.parseForStmt()
	case token.WHILE:
		s = p.parseWhileStmt()
	case token.REPEAT:
		s = p.parseRepeatStmt()
	case token.RETURN:
		s = p.parseReturnStmt()
	case token.BREAK:
		s = &ast.BreakStmt{Keyword: p.pos}
		p.next()
	case token.NEXT:
		s = &ast.NextStmt{Keyword: p.pos}
		p.next()
	case token.VERSION:
		s = &ast.VersionStmt{Keyword: p.pos}
		p.next()
	case token.LBRACE:
		s = p.parseBlockStmt()
		p.expectSemi()
	case token.SEMICOLON:
		// Is it ever possible to have an implicit semicolon
		// producing an empty statement in a valid program?
		// (handle correctly anyway)
		s = &ast.EmptyStmt{Semicolon: p.pos, Implicit: p.lit == "\n"}
		p.next()
	case token.RBRACE:
		// a semicolon may be omitted before a closing "}"
		s = &ast.EmptyStmt{Semicolon: p.pos, Implicit: true}
	case token.EOF:
		println("EOF encountered during parseStmt")
		s = &ast.EOFStmt{EOF: p.pos}
	default:
		// no statement found
		pos := p.pos
		p.errorExpected(pos, "statement")
		s = &ast.BadStmt{From: pos, To: p.pos}
	}
	return
}
