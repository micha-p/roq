cat("test 1:\n")
f<-function(){for (e in list(1,2,3)) print(e)}
f()
cat("test 2:\n")
f<-function(a,...,b){print(..1);print(..2)}
f(1,1+1,3,b=4)
cat("test 3:\n")
f<-function(a,...,b){list(...)}
f(1,1+1,3,b=4)

cat("test 3:\n")
f<-function(a,...,b){list(...,5,6)}
f(1,1+1,3,b=4)

cat("test 4:\n")
f<-function(a,...,b){list(-2,-1,0,...)}
f(1,1+1,3,b=4)

cat("test 5:\n")
f<-function(a,...,b){list(-2,-1,0,...,5,6)}
f(1,1+1,3,b=4)
