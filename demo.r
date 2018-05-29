# go run main.go -E -file all.r
# echo ". <- NA" | cat - all.r | R --no-save --interactive --quiet


### Very Basic


1+2

a <- 3
3 -> d


af <- function(a) {
	a+1
}
af(1)

void <- function(a) {
	return()
}
void(2)

zero <- function() {
	return(0)
}
zero()


### Parsing

r = 1  +2 %% 10 + . ;
a.a=23. + .1 * a_2 + ..b * ._c / a...b
C= 1 +2 * 3 ^ 2 / 1 * .



5.5

2.0+(3.0*4.0)


2^2

3.4 %% 1.0


b<-function(c,d){1+2}
b
b(1,2)
b(c=1,d=2)
b(3,4)
b(3,4+1)
unkn(3)

f<-function(a,b,c,d){1000*a + 100*b + 10*c + d}
f(a=1,b=2,c=3,d=4)
f(b=1,a=2,c=3,d=4)
f(1,2,3,4)
f(1,2,3,c=4)
f(c=3,d=4,1,1+1)


f<-function(ma,mb,c,d){1000*a + 100*b + 10*c + d}
f(ma=1,mb=2,c=3,d=4)
a<-1
b<-1
f(ma=1,mb=2,c=3,d=4)
f(ma=1,mb=2,c=3,d=4,x=1)
f(ma=1,mb=2,c=3,d=4,x=1,y=1)

f<-function(ma,mb,c,d){1000*a + 100*b + 10*c + d}
f(ma=1,mb=2,c=3,d=4)
a<-1
b<-1
f(ma=1,mb=2,c=3,d=4)
f(ma=1,mb=2,c=3,d=4,x=1)
f(ma=1,mb=2,c=3,d=4,x=1,y=1)
f(ma=1,mb=2,c=3,d=4,5)
f(ma=1,mb=2,c=3,d=4,5,6)


f<-function(ma,mb,c,d){1000*ma + 100*mb + 10*c + d}
f(ma=1,mb=2,c=3,d=4)
f(ma=1,m=2,c=3,d=4)
f(ma=1,m=2,3,4)


### Scope

c=1
f<-function(a,b){c=100;a+b+c}
c
f(1,2)
c

{
# [1] 3
a<-11
#f(a<-22,2)
# [1] 24
a
# [1] 22
}
{
a<-11
a
}
f<-function(a,b)a+b
f(1,2)
#[1] 3

f<-function(a,b) a+b
f(1,2)
#[1] 3

f<-function(a,b) a+b; 10+11
#[1] 21

f(1,2)
#[1] 3

a<-1
f<-function(a,b) a<-2
f(1,2)
a
#[1] 1



#f<-function(a,b) if(a)a else 33
#f(1,2)
##[1] 1
#f(0,2)
##[1] 33

if(1){1}else{33}
if(1)1 else 33
if(0){1}else{33}
if(0)1 else 33

if (TRUE) {1} else {2}
if (TRUE) {7}

if (TRUE) 3 else {4};
if (TRUE) 3 else {4}
if (TRUE) 6

if (1>0) {3} else 4;
if (1>0) {3} else 4


f<-function(a,b)a+b
f(1,2)
#[1] 3

f<-function(a,b) a+b
f(1,2)
#[1] 3

f<-function(a,b) a+b; 10+11
#[1] 21

f(1,2)
#[1] 3

a<-133
f<-function(a,b){a<-2}
f(1,2)
a
#[1] 133


f<-function(a,b){if(a)a else 33}
f(1,2)
#[1] 1
f(0,2)
#[1] 33

f<-function(a,b)a<-2  # TODO assignment is not recognized in parseBlockStmt1
f(1,2)
f<-function(a,b)if(a)a else 33
f(1,2)
# [1] 1
f<-function(a,b)a<-2
f(1,2)
f<-function(a,b)if(a)a else 33
f(1,2)
# [1] 1
f<-function(a,b)a+b
f(1,2)
#[1] 3

f<-function(a,b) a+b
f(1,2)
#[1] 3

f<-function(a,b) a+b; 10+11
#[1] 21

f(1,2)
#[1] 3

a<-133
f<-function(a,b){a<-2}
f(1,2)
a
#[1] 133


f<-function(a,b){if(a)a else 33}
f(1,2)
#[1] 1
f(0,2)
#[1] 33

f<-function(a,b)a<-2
f(1,2)
f<-function(a,b)if(a)a else 33
f(1,2)
# [1] 1


# error in R
# if 3 4 else 5
# if 3+4 4 else 5
# if 3-3 4 else 5

names <- c("alp", "bet", "gamm")
#for(n in names) print(n)
vect <-c(1,2,3,4)
vect

# to be fixed
#for(a in vect) {cat(a+1)}

cat(5)
cat(5,6,"a",7,"\n")


TRUE
FALSE
NULL
Inf
NA
NaN
#1<2
#1>2
#1==2
#1==1
TRUE && FALSE
TRUE && TRUE
TRUE || FALSE
FALSE || FALSE

cat("\n")

1 && 2
TRUE && 1
TRUE && 0

cat("\n")

1 || 2 || 0
1 && 2 && 0
0 || 2 && 1
0 || 2 && 0

cat("\n")

1 && 2 || 0
1 && 0 || 0
1 && 0 || 1

#1<2
#1>2
#1==2
#1==1
1<2
1>2
1==2
1==1
1!=2
1!=1

cat("\n")
3 >= 4
3 >= 3
3 <= 4
3 <= 2

cat("\n")
1 < 2 < 3 < 4
3 < 5 > 1
1 < 2 < 3
1 < 3 < 2
1 < 3 > 2
1 < 2 <= 3
3 > 2 > 1

if(1<2) 3 else 4
if(1<2<5) 3 else 4
if(0<2) 3 else 4

c(1,2,3)
cat("\n")

c(1,1+1,3)
cat("\n")


f<-function(a,b,c){
cat(a)
cat(b)
cat(c)
cat("\n")
}

f(1,2,3)


#b=1
#b
#c(1,(b<-2),3)
#b
#cat("\n")

#b
#c(1,b=5,3)
#b
#cat("\n")
b=1.0
b
c(1,b<-2.0,3)
b

d=33
#c(1,b<-2.2,d=3)
b
d
c(1,2,3) + c(4,5,6)

c(1,2,3) + c(1,1,1,1,1)

c(1,1) + c(1,2,3,4,5)

c(2) * c(1,2,3,4,5)

c(1,2,3,4,5) / c(10)

c(1,2,3,4,5,6,7,8,9,10) %% c(3)

c(1,2,3,4) ^ c(2)
c(1,2,3) + c(4,5,6)

c(1,2,3) + c(1,1,1,1,1)

c(1,1) + c(1,2,3,4,5)

c(2) * c(1,2,3,4,5)
2 * c(1,2,3,4,5)

c(1,2,3,4,5) / c(10)
c(1,2,3,4,5) / 10

c(1,2,3,4,5,6,7,8,9,10) %% c(3)
c(1,2,3,4,5,6,7,8,9,10) %% 3

c(1,2,3,4) ^ c(2)
c(1,2,3,4) ^ 2

1+1
1*1
1 < 2
1 > 2
c(1,2,3) < c(4,5,6)
c(1,2,3) > c(4,5,6)
1 < c(4,5,6)
1 > c(4,5,6)

cat("\n")
1 < 2 < 3
1 > 2 < 3 # 1>2 is FALSE; 0 would be < 3
3 > c(2,2) == 2
# n<-0
# repeat{
# 	cat(n)
# 	cat("\n")
# 	if(n==3)break
# 	n=n+1
# }

cat("\n")

# for (n in c(1,2,3,4,5,6)) {
# 	if(n==3)next
# 	cat(n)
# 	cat("\n")
# 	if(n==5)break
# }

### Indexing

a=c(1,2,3)
a
a[1]
a[1.1]
a[0]
a=c(11,22,33)
b=c(1,2)
d=c(1,3)
a[1]
a[2.2]
a[1:3]
#a[b]
#a[d]


x <- 1.0
class(x)
class(x)<-"myclass"
class(x)
class(x)<-NULL
class(x)
class(1:6)
class("a")
class(NULL)
class(list(1,2,3))
#class(pairlist(1,2))
x <- 1.0
class(x)
class(x)<-"myclass"
class(x)
class(x)<-NULL
class(x)
class(1:6)
class("a")
class(NULL)
#class(list(1,2,3))
#class(pairlist(1,2))

