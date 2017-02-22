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
