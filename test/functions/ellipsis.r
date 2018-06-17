cat("test 1:\n")
f<-function(a,b,c){for (e in list(a,b,c)) print(e)}
f(11,22,33)

cat("\ntest 2:\n")
f<-function(...){print(..1);print(..2)}
f(11,22,33,b=44)

cat("\ntest 3:\n")
f<-function(a,...,b){c(...)}
f(11,11+11,33,b=44)

cat("\ntest 4:\n")
f<-function(...)c(1,2,3,...,4,5)
f(11,22,33)
