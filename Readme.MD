Roq -- an implementation of R in Go

https://blog.gopheracademy.com/advent-2014/parsers-lexers/

https://cran.r-project.org/doc/manuals/R-lang.html

*EXTREMELY ALPHA*

# Usage

```
cd $HOME
mkdir -p go/src
cd go
export GOPATH=$PWD
cd src
git clone https://github.com/micha-p/roq.git
cd roq

go run main.go -expr "1+2"
go run main.go
go test
go build
./roq -file demo.r
```

# Short summary

## Preliminary differences to R

- No welcome message
- No interactive mode, just file or expression
- numeric values are float, integers restricted to index oprations
- single dot means NA
- NA, NaN, ... should be unified to NULL and NaN (less decicions)
- false values are treated differently
- limited access to workspace and environment
- simplified output


## Additional features (TODO)

- multiline strings
- multiline arrays and matrices
- concurrency on the interpreter-level of course

## Dirty hacks (TO BE REMOVED)

- exit as recovered panic

# Notes on the planned implementation

## Internal representation

In R, the types for storing data objects are not orthogonal: There are two different types 
for single strings (CHARSXP) and arrays of pointers to strings (STRSXP). In contrast, even atomic 
floats, integers and logical values are stored in arrays (REALSXP, INTSXP, LGLSXP). 
Tagged pairs (LISTSXP, LANGSXP), generic vectors and expressions (VECSXP, EXPRSXP).

This implementation tries to use different domains with very restricted types (Values, Indices, Strings) and evaluates the usefulness 
of a common interface to them.

## Value and index domain

Strictly, there should be only one type of numbers in each domain. float64 on toplevel and unsigned integers within brackets. 

To support wider compatibility, a more relaxed mode is neccesssary: Value domain defaults to double. Index domain defaults to int. 
This way, integer vecors might be coposed in toplevel. 
Functions on numbers should not rely on the types.

## Output

For printing, length is shown instead of linecount, but this information might be dropped completely (TO BE EVALUATED). 
