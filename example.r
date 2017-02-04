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
