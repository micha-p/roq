
strict mode enforcing more compatibility

Differences:
- only standard  attributes
- probably evaluating zero to TRUE (in progress)
- isTRUE will drop a warning, as it is misleading
- break/next in wrong context will stop, as running code is expected
- native lowlevel print statement: print a 


	// returning values allows for concatenated comparisons
	// No problem UNLESS X == 0!!!
	// TODO documentation of extension
	// Either treat zero as true, which will break existing R code
	// or restrict such comparisons to inexact floats, where zero is not really existant
	// and take zero in integer context as false:
	// - This would enable zero one encoded vectors 
	// - even then, checking integer ranges like 0<i<10 would not be possible
	// TODO: 
	// Final solution: detect concatenated comparisons during parsing and evaluate such multiple comparisons differently
