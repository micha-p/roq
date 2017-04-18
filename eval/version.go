package eval

import "runtime"

func PrintVersion() {
	print("arch           "); println(runtime.GOARCH)
	print("os             "); println(runtime.GOOS)
	print("status         "); println("proof of concept")
	print("major          "); println("0")
	print("minor          "); println("2.0")
	print("language       "); println("R")
}

//platform       x86_64-pc-linux-gnu         
//arch           x86_64                      
//os             linux-gnu                   
//system         x86_64, linux-gnu           
//status                                     
//major          3                           
//minor          2.3                         
//year           2015                        
//month          12                          
//day            10                          
//svn rev        69752                       
//language       R                           
//version.string R version 3.2.3 (2015-12-10)
//nickname       Wooden Christmas-Tree       
