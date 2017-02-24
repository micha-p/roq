# Performance

## Avoiding construction of s-expressions

- replacing scalar FALSE with nil
- logical operators return appropriate argument instead of new boolean s-expression
- immediate floats as scalars

## literal vectors

- avoiding function parsing and call to c(,,)
