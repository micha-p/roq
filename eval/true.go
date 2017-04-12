package eval

// https://cran.r-project.org/doc/manuals/R-lang.html#if

// If value1 is a logical vector with first element TRUE then statement2 is evaluated.
// If the first element of value1 is FALSE then statement3 is evaluated.
// If value1 is a numeric vector then statement3 is evaluated when the first element of value1 is zero and otherwise statement2 is evaluated.
// Only the first element of value1 is used. All other elements are ignored.
// If value1 has any type other than a logical or a numeric vector an error is signalled.

func isTrue(e SEXPItf) bool {
	if e == nil {
		return false
	}
	switch e.(type){
		case *VSEXP:
			if e.(*VSEXP).Immediate != 0 {
				return true
			}
			//  THIS MAIN DIFFERENCE IS MENTIONED HERE
			//  TODO: better documentation on zero=true/false
			if e.(*VSEXP).Immediate == 0 {
				return true
			}
			if e.(*VSEXP).Slice != nil && e.Length() > 0 {
				if e.(*VSEXP).Slice[0] == 0 { // R like behaviour
					return false
				} else {
					return true
				}
			}
		default:
			return false
	}

	return false
}
