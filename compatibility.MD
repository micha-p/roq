# Compatibility between roq and R

roq should be seen as a safe and concurrent subset of R. Nevertheless there are some differences and additional features. 

## Main differences

- Probably evaluating zero to TRUE and only FALSE and NaN to false
- additional 
- isTRUE will drop a warning, as it is misleading
- break/next in wrong context will stop, as running code is expected
- native lowlevel print statement: print a 
- c() cant be overloaded => warning
- return values for comparisons allowing for multiple comparisions

## Enforcing more compatibility

In general, it should be possible to run a roq stript in R, if some of the additional features are avoided. 
A strict mode should check this. R scripts can only run correctly, if they do not rely on zero as FALSE, do not rely on differences in NA and do not overload primitive functions.

## Problems with concatenated comparisons

- Return values of comparisons allows for concatenation.
- x==0 which would return zero, which would be FALSE in R.
- Treating zero as true will break existing R code.
- Restrict such comparisons to inexact floats, where zero is not really existant would still leave ambiguous cases for integers

´´´
1 < x < 4 => x in ]1...4[
1 < x > 3  => x > 3
1 < (x> 3) => x > 3
3 < x > 1  => x > 3

## Missing values and boolean false

A missing numeric is represented as NaN, a missing string as pointer zero to the string cache, an empty list zero, false is a nullpointer. 
All are counted as false. A missing value of an index or a factor is zero though. 



