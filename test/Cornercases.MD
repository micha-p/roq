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



