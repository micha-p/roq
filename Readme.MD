# Roq -- an implementation of R in Go

*EXTREMELY ALPHA*

This is an implementation of the core language of R, the popular implementation of the S scripting language. It is not a complete computing environment with base, stats and graphics libraries, but covers all features of a functional programming language with C-like syntax. So far, parser and scanner are taken from Go 1.7, but this might change.

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
go run main.go -file demo.r
go run main.go -E -file demo.r

go test -v
go build
./roq -file example.r
./roq                  # interactive (one expression per line)
rlwrap ./roq           # with readline support
```

# Short summary

## Preliminary differences to R

- numeric values are float64
- integers restricted to index oprations
- single dot means NA
- limited support for attribute propagation
- simplified trueness for less decisions: NULL and NaN are false, all values including 0 are true (TODO: evaluation of performance)
- very limited access to filesystem (like chroot) and environment
- simplified output with arrays in one line


## Additional features (TODO)

- restrictions on types should improve performance
- concurrency on the interpreter-level of course

## links

https://blog.gopheracademy.com/advent-2014/parsers-lexers/

https://cran.r-project.org/doc/manuals/R-lang.html

