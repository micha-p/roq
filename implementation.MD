# Internal representation

In R, the types for storing data objects are not orthogonal: There are two different types 
for single strings (CHARSXP) and arrays of pointers to strings (STRSXP). In contrast, even atomic 
floats, integers and logical values are stored in arrays (REALSXP, INTSXP, LGLSXP). 
Tagged pairs (LISTSXP, LANGSXP), generic vectors and expressions (VECSXP, EXPRSXP).

This implementation tries to use different domains with very restricted types (Values, Indices, Strings) and evaluates the usefulness 
of a common interface to them.

## Value and index domain

Strictly, there should be only one type of numbers in each domain. float64 on toplevel and unsigned integers within brackets. 

To support wider compatibility, a more relaxed mode is neccesssary: Value domain defaults to double. Index domain defaults to int. 
This way, integer vectors might be composed at the toplevel. 
Functions on numbers should not rely on the types.

# Types and their representation

- Values are organized in domains: numeric, text, indices, recursive structures
- Domains are used during scanning
- sequences might be used in a numeric or index context
- So far, user defined functions are restricted to numerics and therefore share the same structure
- The Null/false object is only used, where nil is not possible

## Objects supporting arrays and an immediate value

VSEXP	numeric		float64
TSEXP	strings		[]char
ISEXP	indices		uint

## Object supporting arrays and several immediate values

RSEXP	cons		pointer

## Other objects

NSEXP	false		bool
ESEXP	token		token (various errors and break conditions)

# Input

Parser and scanner of golang use unbuffered input until end of file. Interactive mode would therefore require a completely different scanner structure.
Therefore input is restricted to files, unbuffered stdin and single expressions, which might be lists of statements. 


# Output

For printing, output for arrays is withoul linebreaks, so linecount is always [1] 
