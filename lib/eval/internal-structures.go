package eval

import (
	"lib/go/ast"
	"lib/go/token"
)

type SEXPTYPE int

// The list of tokens.
const (
	NILSXP     = iota //	0	NULL
	SYMSXP            //	1	symbols
	LISTSXP           //	2	pairlists
	CLOSXP            //	3	closures
	ENVSXP            //	4	environments
	PROMSXP           //	5	promises
	LANGSXP           //	6	language objects
	SPECIALSXP        //	7	special functions
	BUILTINSXP        //	8	builtin functions
	CHARSXP           //	9	internal character strings
	LGLSXP            //	10	logical vectors
	INTSXP            //	13	integer vectors
	REALSXP           //	14	numeric vectors
	CPLXSXP           //	15	complex vectors
	STRSXP            //	16	character vectors
	DOTSXP            //	17	dot-dot-dot object
	ANYSXP            //	18	make “any” args work
	VECSXP            //	19	list (generic vector)
	EXPRSXP           //	20	expression vector
	BCODESXP          //	21	byte code
	EXTPTRSXP         //	22	external pointer
	WEAKREFSXP        //	23	weak reference
	RAWSXP            //	24	raw vector
	S4SXP             //	25	S4 classes not of simple type
)

type sxpinfo_struct struct {
	sexptype SEXPTYPE //  5;  /* discussed above */
	obj      bool     //  1;  /* is this an object with a class attribute? */
	named    byte     //  2;  /* used to control copying */
	gp       uint16   // 16;  /* general purpose, see below */
	mark     bool     //  1;  /* mark object as ‘in use’ in GC */
	debug    bool     //  1;
	trace    bool     //  1;
	spare    bool     //  1;  /* debug once */
	gcgen    bool     //  1;  /* generation for GC */
	gccls    byte     //  3;  /* class of node for GC */
} /*              Tot: 32 */

/*
type SEXPREC struct {
	header sxpinfo_struct
	attributes interface{}
	previous interface{}
	next interface{}
	data interface{}
}
*/

type SEXPREC struct {
	ValuePos  token.Pos      // literal position
	Kind      token.Token    // token.INT, token.FLOAT, token.IMAG, token.CHAR, or token.STRING
	Fieldlist []*ast.Field   // only if function
	Body      *ast.BlockStmt // only if function
	Value     float64
}
