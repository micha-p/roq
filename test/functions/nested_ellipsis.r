

cat("test 1:\n")
f<-function(...){c(...,5,6)}
g<-function(...){f(...,55)}
g(11,11+11)

cat("\ntest 2:\n")
f<-function(...){c(1,2,...,5,6)}
g<-function(...){f(-55,...,55)}
g(11,11+11)

cat("\ntest 3:\n")
f<-function(...){c(1,2,...,5,6)}
g<-function(...){f(-55,...,55)}
g(11,11+11,33,b=44)
