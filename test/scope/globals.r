
a<<-2
a
4->>b
b

cat("\n")

sum <- 10
add <- function(x){sum <<- sum+x}
sum
add(1) 		# just for the sideeffect, superassignment has no return value
sum

