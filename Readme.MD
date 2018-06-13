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
go run main.go -file demo.r
go run main.go -E -file demo.r

go test -v
go build
./roq -file example.r
rlwrap ./roq # for readline support
```

# Short summary

## Preliminary differences to R

- numeric values are float64
- integers restricted to index oprations
- single dot means NA
- NA, NaN, ... should be unified to NULL and NaN (less decicions)
- false values are treated differently
- very limited access to filesystem and environment
- simplified output with arrays in one line


## Additional features (TODO)

- restrictions on type occurance should improve performance
- concurrency on the interpreter-level of course

