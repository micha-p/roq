f<-function(ma,mb,c,d){1000*ma + 100*mb + 10*c + d}
f(ma=1,mb=2,c=3,d=4)
f(ma=1,m=2,c=3,d=4)		# ma is satisfied, so m matches only mb
f(mb=2,m=1,3,4)			# mb is satisfied, so m matches only ma
f(m=1,2,c=3,d=4)
