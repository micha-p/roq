package main

import (
	"strconv"
	"fmt"
	"math"
	"lib/go/token"
	"lib/go/ast"
)


// ----------------------------------------------------------------------------
// Tracing support

// https://go-book.appspot.com/interfaces.html
// an empty interface accepts all pointers

func evalStmt(s interface{}){
	switch s.(type) {
	case *ast.ExprStmt:
	  if TRACE {println("exprStmt")}
	  e := s.(*ast.ExprStmt)
	  fmt.Printf("%g",evalExpr(e.X))   // R has small e
	case *ast.EmptyStmt:
	  println("")
	case *ast.IfStmt:
	  if TRACE {println("ifStmt")}
	case *ast.ForStmt:
	  if TRACE {println("forStmt")}
	case *ast.BlockStmt:
	  if TRACE {println("blockStmt")}
	default:
	  println("? Stmt")
	}
}

func evalExpr(e ast.Expr) float64 {
	switch e.(type) {
	case *ast.BasicLit:
	  node := e.(*ast.BasicLit)
	  if TRACE {println("BasicLit " + " " + node.Value +" ("+ node.Kind.String() + ")")}
	  v, err := strconv.ParseFloat(node.Value, 64) // TODO: support for all R formatted values
	  if err != nil {println (err)}
	  return v
	case *ast.BinaryExpr:
	  node := e.(*ast.BinaryExpr)
	  if TRACE {println("BinaryExpr " + " " + node.Op.String())}
	  return evalOp(node.Op, evalExpr(node.X), evalExpr(node.Y))
	case *ast.ParenExpr:
	  node := e.(*ast.ParenExpr)
	  if TRACE {println("ParenExpr")}
	  return evalExpr(node.X )
	default:
	  println("? Expr")
	  return math.NaN()
	}
}

func evalOp(op token.Token, x float64, y float64) float64 {
	switch op {
	case token.PLUS:
	  return x + y
	case token.MINUS:
	  return x - y
	case token.MULTIPLICATION:
	  return x * y
	case token.DIVISION:
	  return x / y
	case token.EXPONENTIATION:
	  return math.Pow(x,y)
	case token.MODULUS:
	  return math.Mod(x,y)
	default:
	  println("? Op: " + op.String())
	}
	return math.NaN()
}
