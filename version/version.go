package version

import "runtime"

func PrintVersion(MAJOR string, MINOR string) {
	print("arch           "); println(runtime.GOARCH)
	print("os             "); println(runtime.GOOS)
	print("status         "); println("proof of concept")
	print("major          "); println(MAJOR)
	print("minor          "); println(MINOR)
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
//minor          2.3                         
//year           2015                        
//month          12                          
//day            10                          
//svn rev        69752                       
//language       R                           
//version.string R version 3.2.3 (2015-12-10)
//nickname       Wooden Christmas-Tree       


//R --version
//R version 3.4.3 (2017-11-30) -- "Kite-Eating Tree"
//Copyright (C) 2017 The R Foundation for Statistical Computing
//Platform: x86_64-pc-linux-gnu (64-bit)
