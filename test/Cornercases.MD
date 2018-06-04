

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



