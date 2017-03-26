package eval

import (
	"lib/ast"
	"lib/token"
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
	INTSXP     = 13   //	13	integer vectors
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
type SEXP struct {
	header sxpinfo_struct
	attributes interface{}
	previous *SEXP
	next *SEXP
	data interface{}
}
*/

// generic -> should be changed to SEXP
type SEXPItf interface {
	Pos()		token.Pos
	Kind()		token.Token
	Dim()		[]int
	Atom()		interface{}
	Length()	int
}

// TODO split into several types 
type SEXP struct {
	ValuePos  token.Pos
	TypeOf    SEXPTYPE
	kind      token.Token    // token.INT, token.FLOAT, token.IMAG, token.CHAR, or token.STRING

	Names     []string
	dim       []int
	Dimnames  [][]string 

	Fieldlist []*ast.Field   // only if function
	Body      *ast.BlockStmt // only if function: BlockStmt or single Stmt
	String    string
	Immediate float64        // single value FLOAT
	Integer   int            // single value INT
	Offset    int            // single value INT (zerobased); TODO change to uint in indexdomain?
	Slice     []float64      // "A slice is a reference to an array"
}


// value domain
type VSEXP struct {
	TypeOf    SEXPTYPE
	kind      token.Token    // token.INT, token.FLOAT, token.IMAG, token.CHAR, or token.STRING

	Names     []string
	dim       []int
	Dimnames  [][]string 

	Fieldlist []*ast.Field   // only if function
	Body      *ast.BlockStmt // only if function: BlockStmt or single Stmt
	Immediate float64        // single value FLOAT
	slice     []float64      // "A slice is a reference to an array"
}

// index domain
type ISEXP struct {
	TypeOf    SEXPTYPE
	kind      token.Token    // token.INT, token.FLOAT, token.IMAG, token.CHAR, or token.STRING

	Names     []string
	dim       []int
	Dimnames  [][]string 

	Integer   int            // single value INT
	Offset    int            // single value INT (zerobased); TODO change to uint in indexdomain?
	slice     []float64      // "A slice is a reference to an array"
}

// recursive domain
type RSEXP struct {
	TypeOf		SEXPTYPE
	kind		token.Token    // token.INT, token.FLOAT, token.IMAG, token.CHAR, or token.STRING

	Names		[]string
	dim			[]int
	Dimnames	[][]string 

	CAR			*SEXP
	CDR			*SEXP
	TAG			*SEXP
	slice		[]*SEXP
}

func (x *SEXP) Pos() token.Pos {
	return x.Pos()
}

func (x *SEXP) Atom() interface{} {
	return x.Immediate
}
func (x *SEXP) Dim() []int {
	return x.dim
}
/*
func (x *SEXP) Slice() []interface{} {
	return x.slice.([]interface{})
}
*/
func (x *SEXP) Length() int {
	return len(x.Slice)
}
func (x *SEXP) Kind() token.Token {
	return x.kind
}
