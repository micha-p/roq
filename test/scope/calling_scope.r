# an extremely  difficult case, is R doing it correctly?

a <- function(b){c <<- c+1; b+c}
c=0
a(0) #[1] 1
a(0) #[1] 2
a(0) #[1] 3

c = 0
a(c <- 100)     #[1] 200; c is updated in the scope of the caller
c               #[1] 100, not 101, as function body
