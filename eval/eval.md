# Notes on evaluation

This is, where the real work is done.
The Evaluator takes an ast-node and returns a s-expressions depending on type
Besides this, the evaluator holds a state, if it is inside a loop and a flag for invisible output.


## Dispatching
entry is always eval.go
builtin commands are located in primitives.go, 
user defined-functions are processed in call.go


eval.go     call.go         vector.go       float.go
ast->SEXP   SEXP<->SEXP     SEXP<->float    float

            primitive.go    index.go        integer.go
            SEXP<->SEXP     SEXP<->int      int

            


