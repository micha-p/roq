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
	DimSet([]int)
	Dimnames()	*RSEXP
	DimnamesSet(*RSEXP)
//	Atom()		interface{} // TODO Length=1 => Atom(), is this dispatching really faster?
	Length()	int
}

// TODO split into several types 

// value domain
type VSEXP struct {
	ValuePos  token.Pos
	TypeOf    SEXPTYPE
	kind      token.Token    // token.INT, token.FLOAT, token.IMAG, token.CHAR, or token.STRING

	Names     []string
	dim       []int
	dimnames  *RSEXP

	Fieldlist []*ast.Field   // only if function
	Body      *ast.BlockStmt // only if function: BlockStmt or single Stmt
	String    string
	Immediate float64        // single value FLOAT
	Integer   int            // single value INT // TODO move into indexdomain?
	Offset    int            // single value INT (zerobased); TODO move into indexdomain?
	Slice     []float64      // "A slice is a reference to an array"
}


// index domain
type ISEXP struct {
	ValuePos  token.Pos
	TypeOf    SEXPTYPE
	kind      token.Token    // token.INT, token.FLOAT, token.IMAG, token.CHAR, or token.STRING

	Names     []string
	dim       []int
	dimnames  *RSEXP

	Integer   int            // single value INT
	Offset    int            // single value INT (zerobased); TODO change to uint in indexdomain?
	Slice     []int          // "A slice is a reference to an array"
}

// recursive domain
type RSEXP struct {
	ValuePos  token.Pos
	TypeOf		SEXPTYPE
	kind		token.Token    // token.INT, token.FLOAT, token.IMAG, token.CHAR, or token.STRING

	Names		[]string
	dim			[]int
	dimnames	*RSEXP

	CAR			SEXPItf
	CDR			SEXPItf
	TAG			SEXPItf
	Slice		[]SEXPItf
}


// text domain (strings, factors and symbols)
type TSEXP struct {
	ValuePos  token.Pos
	TypeOf    SEXPTYPE
	kind      token.Token    // token.INT, token.FLOAT, token.IMAG, token.CHAR, or token.STRING

	Names     []string
	dim       []int
	dimnames  *RSEXP

	String    string
	Slice     []string      // "A slice is a reference to an array"
}


func (x *VSEXP) Pos() token.Pos {
	return x.ValuePos
}
func (x *VSEXP) Atom() interface{} {
	return x.Immediate
}
func (x *VSEXP) Dim() []int {
	return x.dim
}
func (x *VSEXP) DimSet(v []int) {
	x.dim=v
}
func (x *VSEXP) Dimnames() *RSEXP {
	return x.dimnames
}
func (x *VSEXP) DimnamesSet(v *RSEXP) {
	x.dimnames=v
}
func (x *VSEXP) Length() int {
	return len(x.Slice)
}
func (x *VSEXP) Kind() token.Token {
	return x.kind
}

func (x *RSEXP) Pos() token.Pos {
	return x.ValuePos
}
func (x *RSEXP) Dim() []int {
	return x.dim
}
func (x *RSEXP) DimSet(v []int) {
	x.dim=v
}
func (x *RSEXP) Dimnames() *RSEXP {
	return x.dimnames
}
func (x *RSEXP) DimnamesSet(v *RSEXP) {
	x.dimnames=v
}
func (x *RSEXP) Length() int {
	return len(x.Slice)
}
func (x *RSEXP) Kind() token.Token {
	return x.kind
}

func (x *TSEXP) Pos() token.Pos {
	return x.ValuePos
}
func (x *TSEXP) Dim() []int {
	return x.dim
}
func (x *TSEXP) DimSet(v []int) {
	x.dim=v
}
func (x *TSEXP) Dimnames() *RSEXP {
	return x.dimnames
}
func (x *TSEXP) DimnamesSet(v *RSEXP) {
	x.dimnames=v
}
func (x *TSEXP) Length() int {
	return len(x.Slice)
}
func (x *TSEXP) Kind() token.Token {
	return x.kind
}

func (x *ISEXP) Pos() token.Pos {
	return x.ValuePos
}
func (x *ISEXP) Dim() []int {
	return x.dim
}
func (x *ISEXP) DimSet(v []int) {
	x.dim=v
}
func (x *ISEXP) Dimnames() *RSEXP {
	return x.dimnames
}
func (x *ISEXP) DimnamesSet(v *RSEXP) {
	x.dimnames=v
}
func (x *ISEXP) Length() int {
	return len(x.Slice)
}
func (x *ISEXP) Kind() token.Token {
	return x.kind
}
