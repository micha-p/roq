a=c(11,22,33)
length(a)
length(c(1,2))
length(a[1])
length(a[2.2])
length(a[1:2])
length(a[1:3])

cat("\n")
x <- c(1,2,3,4,5,6)
dim(x) <- c(2,3)
dim(x)
x

cat("\n")
x <- 1:6
x
x[2]
dim(x) <- c(2,3)
cat("\n")
dim(x)
x

cat("\n")
y <- 1:24
dim(y) <- c(2,3,2,2)
y

cat("\n")
list("a1","a2")
list(1,2.0)

cat("\n")
x <- 1:6
dim(x)
dim(x) <- c(2,3)
dimnames(x) <- list(c("a1","a2"),c("b1","b2","b3"))
x
