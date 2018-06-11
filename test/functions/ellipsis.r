cat("test 1:\n")
f<-function(){for (e in list(1,2,3)) print(e)}
f()

cat("test 2:\n")
f<-function(a,b,c){for (e in list(a,b,c)) print(e)}
f(11,22,33)

cat("test 3:\n")
f<-function(...){print(..1);print(..2)}
f(11,22,33,b=44)

cat("test 4:\n")
f<-function(a,...,b){list(...)}
f(11,11+11,33,b=44)


cat("test 5:\n")
f<-function(a,...,b){list(-2,-1,0,...,5,6)}
f(11,11+11,33,b=44)

