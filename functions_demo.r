		f<-function(x){1}
		f(1)
		g<-function(x){x+1}
		g(1)
		f<-function(x){x+2}
		f(1)

cat("\n")

		b<-function(c,d){c+d}
		b(c=1,d=2)
		b(3,4)
		b(3,4+1)
		a <- 5
		b(c=3,d=a+1)

cat("\n")




f<-function(a,b)a+b
f(1,2)
f<-function(a,b) a+b
f(1,2)
f<-function(a,b) a+b; 10+11
f(1,2)

cat("\n")

a<-133
f<-function(a,b){a<-2}
f(1,2) # there is no return value
f<-function(a,b)a<-2
f(1,2)
a
