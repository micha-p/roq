package eval

import (
	"lib/ast"
	"lib/token"
	"math"
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
	Pos() token.Pos
	Dim() []int
	DimSet([]int)
	Dimnames() *RSEXP
	DimnamesSet(*RSEXP)
	Class() *string
	ClassSet(*string)
	//	Atom()		interface{} // TODO Length=1 => Atom(), is this dispatching really faster?
	IntegerGet() int
	FloatGet() float64
	Length() int
}

// promoted fields can't be used in Promoted fields act like ordinary fields of a struct except that they
// cannot be used as field names in composite literals of the struct.
type SEXP struct {
	names    []string
	dim      []int
	dimnames *RSEXP
	class    *string
}

// value domain
type VSEXP struct {
	ValuePos token.Pos
	TypeOf   SEXPTYPE
	Kind     token.Token     // token.INT, token.FLOAT, token.IMAG, token.CHAR, or token.STRING
	                         // TODO: get rid of it
	SEXP
	Fieldlist []*ast.Field   // only if function
	Body      *ast.BlockStmt // only if function: BlockStmt or single Stmt
	Immediate float64        // single value FLOAT
	Integer   int            // single value INT
	Slice     []float64      // "A slice is a reference to an array"
}

// Index domain
type ISEXP struct {
	ValuePos token.Pos
	TypeOf   SEXPTYPE
	kind     token.Token     // token.INT, token.FLOAT, token.IMAG, token.CHAR, or token.STRING

	SEXP
	Immediate float64        // single value FLOAT
	Integer int              // single value INT
	Slice   []int            // "A slice is a reference to an array"
}

// Recursive domain
type RSEXP struct {
	ValuePos token.Pos
	TypeOf   SEXPTYPE
	kind     token.Token
	SEXP
	CAR   SEXPItf
	CDR   SEXPItf
	TAG   SEXPItf
	Slice []SEXPItf
}

// NULL, NA, FALSE
type NSEXP struct {
	SEXP
	ValuePos token.Pos
	TypeOf   SEXPTYPE
	kind     token.Token
}


// Text domain: pointer to cached strings , factors, symbols
type TSEXP struct {
	SEXP
	ValuePos token.Pos
	TypeOf   SEXPTYPE
	kind     token.Token
	String string
	Slice  []string // "A slice is a reference to an array"
}

// Errors and exceptions
type ESEXP struct {
	SEXP
	ValuePos token.Pos
	TypeOf   SEXPTYPE
}

func (x *SEXP) Dim() []int {
	return x.dim
}
func (x *SEXP) DimSet(v []int) {
	x.dim = v
}
func (x *SEXP) Dimnames() *RSEXP {
	return x.dimnames
}
func (x *SEXP) DimnamesSet(v *RSEXP) {
	x.dimnames = v
}
func (x *SEXP) Class() *string {
	return x.class
}
func (x *SEXP) ClassSet(v *string) {
	x.class = v
}

func (x *VSEXP) Pos() token.Pos {
	return x.ValuePos
}
func (x *VSEXP) Atom() interface{} {
	return x.Immediate
}
func (x *VSEXP) Length() int {
	if x.Slice == nil {return 1} else {return len(x.Slice)}
}
func (x *VSEXP) IntegerGet() int {
	// TODO check conversion to integer
	return int(math.Floor(x.Immediate))
}
func (x *VSEXP) FloatGet() float64 {
	return x.Immediate
}


func (x *ISEXP) Pos() token.Pos {
	return x.ValuePos
}
func (x *ISEXP) Length() int {
	if x.Slice == nil {return 1} else {return len(x.Slice)}
}
func (x *ISEXP) IntegerGet() int {
	return x.Integer
}
func (x *ISEXP) FloatGet() float64 {
	return x.Immediate
}



func (x *RSEXP) Pos() token.Pos {
	return x.ValuePos
}
func (x *RSEXP) Length() int {
	if x.Slice == nil {return 1} else {return len(x.Slice)}
}
func (x *RSEXP) IntegerGet() int {
	panic("Trying to get an integer from a list")
	return 0
}
func (x *RSEXP) FloatGet() float64 {
	panic("Trying to get a float from a list")
	return 0
}

func (x *TSEXP) Pos() token.Pos {
	return x.ValuePos
}
func (x *TSEXP) Length() int {
	if x.Slice == nil {return 1} else {return len(x.Slice)}
}
func (x *TSEXP) IntegerGet() int {
	panic("Trying to get an integer from characters")
	return 0
}
func (x *TSEXP) FloatGet() float64 {
	panic("Trying to get a float from characters")
	return 0
}

func (x *NSEXP) Pos() token.Pos {
	return x.ValuePos
}
func (x *NSEXP) Length() int {
	return 0
}
func (x *NSEXP) IntegerGet() int {
	panic("Trying to get a float from NULL")
	return 0
}
func (x *NSEXP) FloatGet() float64 {
	panic("Trying to get a float from NULL")
	return 0
}

// Errors and exceptions
func (x *ESEXP) Pos() token.Pos {
	return x.ValuePos
}
func (x *ESEXP) Length() int {
	return 0
}
func (x *ESEXP) IntegerGet() int {
	panic("Trying to get a float from an ESEXP")
	return 0
}
func (x *ESEXP) FloatGet() float64 {
	panic("Trying to get a float from an ESEXP")
	return 0
}
