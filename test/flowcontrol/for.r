for (e in c(1,2,3)) {print(e)}

cat("\n")

for (e in list(1,2,3)) {print(e)}

cat("\n")

for (n in c(1,2,3,4,5,6)) {
	if(n==3)next
	print(n)
	if(n==5)break
}
