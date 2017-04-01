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

c(2,3)
c(2.0,3.0)
c(2.0,3)
# c(2,3.0) will fail
