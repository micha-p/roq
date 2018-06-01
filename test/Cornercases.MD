1. Return values of assignments

a <- function(b)b*2

b=2			# this is a shorthand assignment
a(b=4)  	# this is a tagged argument and not an assignemnt
#[1] 8

b
#[1] 2
a(b<-4) value of an assignment (in the calling frame ) is rhs
# [1] 8
b
# [1] 4


3. Calling and evaluated frames  

c <- 1
a <- function(b)b+c

a(1)
#[1] 2

c <- 2
a(1)		# c is lexically scoped (looked up dynamically)
#[1] 3

a(c<-3)		# c is modified in the calling frame
# [1] 6

4. Assignemnts to outer scopes

a <- function(b){c <<- c+1; b+c}
c = 0

a(0)
#[1] 1

a(0)
#[1] 2

a(0)
#[1] 3

c = 0
a(c <- 100)
# [1] 200

c
#[1] 100 	# and not 101

5. Assignment precedence

> a <- b <- 2
> a <- (b <- 2)
> (a <- b) <- 2
Error in (a <- b) <- 2 : could not find function "(<-"

> a = b = 2
> a = b <- 2
> a <- b = 2
Error in (a <- b) = 2 : could not find function "<-<-"

> a = (b = 2)
> (a = b) = 2
Error in (a = b) = 2 : could not find function "(<-"

> 2 -> b -> a
> (2 -> b) -> a
> 2 -> (b -> a)
Error in (a <- b) <- 2 : could not find function "(<-"



