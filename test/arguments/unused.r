options(error=expression(NULL))

f<-function(ma,mb,c,d){1000*ma + 100*mb + 10*c + d}
f(ma=1,mb=2,c=3,d=4)
f(ma=1,mb=2,c=3,d=4,x=1)
f(ma=1,mb=2,c=3,d=4,x=a)
f(ma=1,mb=2,c=3,d=4,x=h(a))
f(ma=1,mb=2,c=3,d=4,x=10,y=100)
f(ma=1,mb=2,c=3,d=4,55)
f(ma=1,mb=2,c=3,d=4,55,66)
