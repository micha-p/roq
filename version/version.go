package version

import "runtime"

const MAJOR = "0"
const MINOR = "1.6"
const YEAR = 2018
const MONTH = 06
const DAY = 11

func PrintVersion() {
	print("platform       "); println(runtime.GOARCH + "-pc-" + runtime.GOOS)
	print("arch           "); println(runtime.GOARCH)
	print("os             "); println(runtime.GOOS)
	print("status         "); println("proof of concept")
	print("major          "); println(MAJOR)
	print("minor          "); println(MINOR)
	print("year           "); println(YEAR)
	print("month          "); println(MONTH)
	print("day            "); println(DAY)
	print("language       "); println("roq")
	print("nickname       "); println("R core in go")
	return
}

//platform       x86_64-pc-linux-gnu         
//arch           x86_64                      
//os             linux-gnu                   
//system         x86_64, linux-gnu           
//status                                     
//major          3                           
//minor          4.4                         
//year           2018                        
//month          03                          
//day            15                          
//svn rev        74408                       
//language       R                           
//version.string R version 3.4.4 (2018-03-15)
//nickname       Someone to Lean On      
