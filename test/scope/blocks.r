a<-11
a
{
	a<-12
	a
}
a

cat("\n")

{
	1 	# invisible
	2	# invisible
	3
}


cat("\n")

f<-function(a,b)a+b
{
	a<-12
	print(a)
	f(a<-22,2)
	a
}
a
