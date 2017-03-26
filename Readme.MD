Roq -- an implementation of R in Go

https://blog.gopheracademy.com/advent-2014/parsers-lexers/

https://cran.r-project.org/doc/manuals/R-lang.html




Internal representation

In R, the types for storing data objects are not orthogonal: There are two different types 
for single strings (CHARSXP) and arrays of pointers to strings (STRSXP). In contrast, even atomic 
floats, integers and logical values are stored in arrays (REALSXP, INTSXP, LGLSXP). 
Tagged pairs (LISTSXP, LANGSXP), generic vectors and expressions (VECSXP, EXPRSXP).

This implementation tries to use different domains with very restricted types (Values, Indices, Strings) and evaluates the usefulness 
of a common interface to them.
