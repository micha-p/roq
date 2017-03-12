// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/* 10.3.1 Constants
 * There are five types of constants: integer, logical, numeric, complex and string.
 * In addition, there are four special constants, NULL, NA, Inf, and NaN.
 * */

package token

import "strconv"

// https://cran.r-project.org/doc/manuals/R-lang.html#Tokens
// Token is the set of lexical tokens

type Token int

// The list of tokens.
const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	COMMENT

	literal_beg
	// Identifiers and basic type literals
	// (these tokens stand for classes of literals)
	IDENT // main
	INT   // 12345
	FLOAT // 123.45
	IMAG  // 123.45i
	//	CHAR   // 'a'
	//	STRING // "abc"

	// R Literals
	IDENTIFIER
	NUMERIC
	INTEGER
	DOUBLE
	LOGICAL
	COMPLEX // wont be implemented in version 1
	STRING

	literal_end

	// https://cran.r-project.org/doc/manuals/R-lang.html#Literal-constants
	// 10.3.1 Literal constants

	constant_beg

	NULL
	NA // TODO: extension documentation: Single dot is treated as missing value
	INF
	NAN
	TRUE
	FALSE

	constant_end

	/* 3.1.4 Operators
	   	   R contains a number of operators. They are listed in the table below.
	   	       -	Minus, can be unary or binary
	   	       +	Plus, can be unary or binary
	   	       !	Unary not
	   	       ~	Tilde, used for model formulae, can be either unary or binary
	   	       ?	Help
	   	       :	Sequence, binary (in model formulae: interaction)
	   	       *	Multiplication, binary
	   	       /	Division, binary
	   	       ^	Exponentiation, binary
	   	       %x%	Special binary operators, x can be replaced by any valid name
	   	       %%	Modulus, binary
	   	       %/%	Integer divide, binary
	   	       %*%	Matrix product, binary
	   	       %o%	Outer product, binary
	   	       %x%	Kronecker product, binary
	   	       %in%	Matching operator, binary (in model formulae: nesting)
	   	       <	Less than, binary
	   	       >	Greater than, binary
	   	       ==	Equal to, binary
	   	       >=	Greater than or equal to, binary
	   	       <=	Less than or equal to, binary
	   	       &	And, binary, vectorized
	   	       &&	And, binary, not vectorized
	   	       |	Or, binary, vectorized
	   	       ||	Or, binary, not vectorized
	   	       <-	Left assignment, binary
	   	       ->	Right assignment, binary
	   	       $	List subset, binary

	   10.3.6 Operator tokens
	   	   R uses the following operator tokens
	   	       + - * / %% ^	arithmetic
	   	       > >= < <= == !=	relational
	   	       ! & |	logical
	   	       ~	model formulae
	   	       -> <-	assignment
	   	       $	list indexing
	   	       :	sequence
	*/

	LPAREN // (
	LBRACK // [
	LBRACE // {
	COMMA  // ,

	RPAREN    // )
	RBRACK    // ]
	RBRACE    // }
	SEMICOLON // ;

	operator_beg
	// Operators and delimiters

	ELLIPSIS // ...

	// R Operators

	MINUS                // -	Minus, can be unary or binary
	PLUS                 // +	Plus, can be unary or binary
	UNARYMINUS           // -	Minus, can be unary or binary
	UNARYPLUS            // +	Plus, can be unary or binary
	NOT                  // !	Unary not
	TILDE                // ~	Tilde, used for model formulae, can be either unary or binary
	HELP                 // ?	Help
	SEQUENCE             // :	Sequence, binary (in model formulae: interaction)
	MULTIPLICATION       // *	Multiplication, binary
	DIVISION             // /	Division, binary
	MODULUS              // %%	Modulus, binary
	EXPONENTIATION       // ^	Exponentiation, binary
	LESS                 // <	Less than, binary
	GREATER              // >	Greater than, binary
	EQUAL                // ==	Equal to, binary
	UNEQUAL              // !=	ADDITIONAL TO DOCUMENTATION
	GREATEREQUAL         // >=	Greater than or equal to, binary
	LESSEQUAL            // <=	Less than or equal to, binary
	ANDVECTOR            // &	And, binary, vectorized
	AND                  // &&	And, binary, not vectorized
	ORVECTOR             // |	Or, binary, vectorized
	OR                   // ||	Or, binary, not vectorized
	SHORTASSIGNMENT      // =	SOMEHOW FORGOTTEN IN DOCUMENTATION
	LEFTASSIGNMENT       // <-	Left assignment, binary
	RIGHTASSIGNMENT      // ->	Right assignment, binary
	SUPERLEFTASSIGNMENT  // <-	Left assignment, binary
	SUPERRIGHTASSIGNMENT // ->	Right assignment, binary
	SUBSET               // $	List subset, binary
	SLOT                 // @	List subset, binary
	DOUBLECOLON          // ::	List subset, binary

	// R SPECIALOPERATORS

	/*
	   %x%	Special binary operators, x can be replaced by any valid name
	   %/%	Integer divide, binary
	   %*%	Matrix product, binary
	   %o%	Outer product, binary
	   %x%	Kronecker product, binary
	   %in%	Matching operator, binary (in model formulae: nesting)
	*/
	operator_end

	keyword_beg
	// Keywords
	//	BREAK
	CASE
	CHAN
	CONST

	IMPORT

	INTERFACE
	MAP
	PACKAGE
	RANGE
	RETURN // not well covered in R language definition

	SELECT
	STRUCT
	SWITCH
	TYPE
	VAR

	// 10.3.3 Reserved words
	IF
	ELSE
	REPEAT
	WHILE
	FUNCTION
	FOR
	IN
	NEXT
	BREAK

	// special command
	VERSION

	keyword_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT: "IDENT",
	INT:   "INT",
	FLOAT: "FLOAT",
	IMAG:  "IMAG",

	// R Literals
	IDENTIFIER: "IDENTIFIER",
	NUMERIC:    "NUMERIC",
	INTEGER:    "INTEGER",
	DOUBLE:     "DOUBLE",
	COMPLEX:    "COMPLEX",
	STRING:     "STRING", // single or double quoted

	NULL:  "NULL",
	NA:    "NA",
	INF:   "Inf",
	NAN:   "NaN",
	TRUE:  "TRUE",
	FALSE: "FALSE",

	// R Operators
	MINUS:                "-",  // -	Minus, can be unary or binary
	PLUS:                 "+",  // +	Plus, can be unary or binary
	NOT:                  "!",  // !	Unary not
	TILDE:                "~",  // ~	Tilde, used for model formulae, can be either unary or binary
	HELP:                 "?",  // ?	Help
	SEQUENCE:             ":",  // :	Sequence, binary (in model formulae: interaction)
	MULTIPLICATION:       "*",  // *	Multiplication, binary
	DIVISION:             "/",  // /	Division, binary
	MODULUS:              "%%", // %%	Modulus, binary
	EXPONENTIATION:       "^",  // ^	Exponentiation, binary
	LESS:                 "<",  // <	Less than, binary
	GREATER:              ">",  // >	Greater than, binary
	EQUAL:                "==", // ==	Equal to, binary
	UNEQUAL:              "!=", // !=	ADDITIONAL TO DOCUMENTATION
	GREATEREQUAL:         ">=", // >=	Greater than or equal to, binary
	LESSEQUAL:            "<=", // <=	Less than or equal to, binary
	ANDVECTOR:            "&",  // &	And, binary, vectorized
	AND:                  "&&", // &&	And, binary, not vectorized
	ORVECTOR:             "|",  // |	Or, binary, vectorized
	OR:                   "||", // ||	Or, binary, not vectorized
	SHORTASSIGNMENT:      "=",  // =	NOT STRITCLY AN OPERATOR, ALSO USED AS ASSIGNMENT
	LEFTASSIGNMENT:       "<-", // <-	Left assignment, binary
	RIGHTASSIGNMENT:      "->", // ->	Right assignment, binary
	SUPERLEFTASSIGNMENT:  "<<-",
	SUPERRIGHTASSIGNMENT: "->>",
	SUBSET:               "$",  // $	List subset, binary
	SLOT:                 "@",  // $	List subset, binary
	DOUBLECOLON:          "::", // Namespace

	ELLIPSIS: "...",

	LPAREN: "(",
	LBRACK: "[",
	LBRACE: "{",
	COMMA:  ",",
	//	PERIOD: ".",

	RPAREN:    ")",
	RBRACK:    "]",
	RBRACE:    "}",
	SEMICOLON: ";",
	//	COLON:     ":",

	CASE:     "case",
	CHAN:     "chan",
	CONST:    "const",

	IMPORT: "import",

	INTERFACE: "interface",
	MAP:       "map",
	PACKAGE:   "package",
	RANGE:     "range",
	RETURN:    "return",

	SELECT: "select",
	STRUCT: "struct",
	SWITCH: "switch",
	TYPE:   "type",
	VAR:    "var",

	//R keywords
	IF:       "if",
	ELSE:     "else",
	REPEAT:   "repeat",
	WHILE:    "while",
	FUNCTION: "function",
	FOR:      "for",
	IN:       "in",
	NEXT:     "next",
	BREAK:    "break",
	VERSION:  "version",
}

// String returns the string corresponding to the token tok.
// For operators, delimiters, and keywords the string is the actual
// token character sequence (e.g., for the token ADD, the string is
// "+"). For all other tokens the string corresponds to the token
// constant name (e.g. for the token IDENT, the string is "IDENT").
//
func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

// A set of constants for precedence-based expression parsing.
// Non-operators have lowest precedence, followed by operators
// starting with precedence 1 up to unary operators. The highest
// precedence serves as "catch-all" precedence for selector,
// indexing, and other operator and delimiter tokens.
//
const (
	LowestPrec  = 2 // non-operators
	UnaryPrec   = 13
	HighestPrec = 16
)

// 10.4.2 Infix and prefix operators
//
// The order of precedence (highest first) of the operators is
//
// ::
// $ @
// ^
// - +                (unary)
// :                  (precedes binary +/-, but not ^)
// %xyz%
// * /
// + -                (binary)
// > >= < <= == !=
// !
// & &&
// | ||
// ~                  (unary and binary)
// -> ->>
// =                  (as short assignment)
// <- <<-

// Precedence returns the operator precedence of the binary
// operator op. If op is not a binary operator, the result
// is LowestPrecedence.
//
func (op Token) Precedence() int {
	switch op {
	case DOUBLECOLON:
		return 16
	case SUBSET, SLOT:
		return 15
	case EXPONENTIATION:
		return 14
	case UNARYMINUS, UNARYPLUS:
		return 13
	case SEQUENCE:
		return 12
	case MODULUS:
		return 11
	case MULTIPLICATION, DIVISION:
		return 10
	case PLUS, MINUS:
		return 9
	case GREATER, GREATEREQUAL, LESS, LESSEQUAL, EQUAL, UNEQUAL:
		return 8
	case AND, ANDVECTOR:
		return 7
	case NOT:
		return 6
	case OR, ORVECTOR:
		return 5
	case TILDE:
		return 4
	case SHORTASSIGNMENT: // TODO CHECK
		return 3
	case LEFTASSIGNMENT, SUPERLEFTASSIGNMENT:
		return 2
	case RIGHTASSIGNMENT, SUPERRIGHTASSIGNMENT:
		return 1
	}
	return LowestPrec
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
	for i := constant_beg + 1; i < constant_end; i++ {
		keywords[tokens[i]] = i
	}
}

// Lookup maps an identifier to its keyword token or IDENT (if not a keyword).
//
func Lookup(ident string) Token {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	return IDENT
}

// Predicates

// IsLiteral returns true for tokens corresponding to identifiers
// and basic type literals; it returns false otherwise.
//
func (tok Token) IsLiteral() bool { return literal_beg < tok && tok < literal_end }

// IsOperator returns true for tokens corresponding to operators and
// delimiters; it returns false otherwise.
//
func (tok Token) IsOperator() bool { return operator_beg < tok && tok < operator_end }

// IsKeyword returns true for tokens corresponding to keywords;
// it returns false otherwise.
//
func (tok Token) IsKeyword() bool { return keyword_beg < tok && tok < keyword_end }

func (tok Token) IsConstant() bool { return constant_beg < tok && tok < constant_end }

// R predicates
func IsAssignment(t Token) bool {
	return t == SHORTASSIGNMENT ||
		t == LEFTASSIGNMENT || t == RIGHTASSIGNMENT || t == SUPERLEFTASSIGNMENT || t == SUPERRIGHTASSIGNMENT
}

func isCONSTANT(t Token) bool {
	return t == INTEGER || t == LOGICAL || t == NUMERIC || t == COMPLEX || t == STRING
}

func isLOGICAL(t Token) bool {
	return t == TRUE || t == FALSE
}

func isNUMERIC(t Token) bool {
	return t == INTEGER || t == DOUBLE || t == COMPLEX || t == NAN || t == INF
}
