
f<-function(a,b){a + b}
f(1,b=2)

f<-function(a,b,...){a + b +..1 + ..2}
f(1,2,3,4)

f<-function(a,...,c){a + ..1 + c}
f(1,2,3)

f(1,2,c=3)
