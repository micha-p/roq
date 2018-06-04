TRUE && FALSE
TRUE && TRUE
TRUE || FALSE
FALSE || FALSE

cat("\n")

1 && 2
TRUE && 1
TRUE && 0
FALSE && 1
FALSE && 0
0 || FALSE
1 || FALSE

cat("\nleft to right\n")


1 || 2 || 0
1 && 2 && 0
0 || 2 && 1
0 || 2 && 0
1 && 2 || 0
1 && 0 || 0
1 && 0 || 1
