# Types and their representation

- Values are organized in domains: numeric, text, indices, recursive structures
- Domains are used during scanning
- sequences might be used in a numeric or index context
- So far, user defined functions are restricted to numerics and therefore share the same structure
- The Null/false object is only used, where nil is not possible

## Objects supporting arrays and an immediate value

VSEXP	numeric		float64
TSEXP	strings		[]char
ISEXP	indices		uint

## Object supporting arrays and several immediate values

RSEXP	cons		pointer

## Other objects

NSEXP	false		bool


