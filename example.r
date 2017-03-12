n<-0
repeat{
	cat(n)
	cat("\n")
	if(n==3)break
	n=n+1
}
	cat("\n")
for (n in c(1,2,3,4,5,6)) {
	if(n==3)next
	cat(n)
	cat("\n")
	if(n==5)break
}
