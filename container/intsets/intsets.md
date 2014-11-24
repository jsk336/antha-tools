# intsets
--
    import "."

Package intsets provides Sparse, a compact and fast representation for sparse
sets of int values.

The time complexity of the operations Len, Insert, Remove and Has is in O(n) but
in practice those methods are faster and more space-efficient than equivalent
operations on sets based on the Go map type. The IsEmpty, Min, Max, Clear and
TakeMin operations require constant time.

## Usage

```go
const (
	MaxInt = int(^uint(0) >> 1)
	MinInt = -MaxInt - 1
)
```
Limit values of implementation-specific int type.

#### type Sparse

```go
type Sparse struct {
}
```

A Sparse is a set of int values. Sparse operations (even queries) are not
concurrency-safe.

The zero value for Sparse is a valid empty set.

Sparse sets must be copied using the Copy method, not by assigning a Sparse
value.

#### func (*Sparse) AppendTo

```go
func (s *Sparse) AppendTo(slice []int) []int
```
AppendTo returns the result of appending the elements of s to slice in order.

#### func (*Sparse) BitString

```go
func (s *Sparse) BitString() string
```
BitString returns the set as a string of 1s and 0s denoting the sum of the i'th
powers of 2, for each i in s. A radix point, always preceded by a digit, appears
if the sum is non-integral.

Examples:

                 {}.BitString() =      "0"
    	     {4,5}.BitString() = "110000"
               {-3}.BitString() =      "0.001"
         {-3,0,4,5}.BitString() = "110001.001"

#### func (*Sparse) Clear

```go
func (s *Sparse) Clear()
```
Clear removes all elements from the set s.

#### func (*Sparse) Copy

```go
func (s *Sparse) Copy(x *Sparse)
```
Copy sets s to the value of x.

#### func (*Sparse) Difference

```go
func (s *Sparse) Difference(x, y *Sparse)
```
Difference sets s to the difference x ∖ y.

#### func (*Sparse) DifferenceWith

```go
func (s *Sparse) DifferenceWith(x *Sparse)
```
DifferenceWith sets s to the difference s ∖ x.

#### func (*Sparse) Equals

```go
func (s *Sparse) Equals(t *Sparse) bool
```
Equals reports whether the sets s and t have the same elements.

#### func (*Sparse) GoString

```go
func (s *Sparse) GoString() string
```
GoString returns a string showing the internal representation of the set s.

#### func (*Sparse) Has

```go
func (s *Sparse) Has(x int) bool
```
Has reports whether x is an element of the set s.

#### func (*Sparse) Insert

```go
func (s *Sparse) Insert(x int) bool
```
Insert adds x to the set s, and reports whether the set grew.

#### func (*Sparse) Intersection

```go
func (s *Sparse) Intersection(x, y *Sparse)
```
Intersection sets s to the intersection x ∩ y.

#### func (*Sparse) IntersectionWith

```go
func (s *Sparse) IntersectionWith(x *Sparse)
```
IntersectionWith sets s to the intersection s ∩ x.

#### func (*Sparse) IsEmpty

```go
func (s *Sparse) IsEmpty() bool
```
IsEmpty reports whether the set s is empty.

#### func (*Sparse) Len

```go
func (s *Sparse) Len() int
```
Len returns the number of elements in the set s.

#### func (*Sparse) Max

```go
func (s *Sparse) Max() int
```
Max returns the maximum element of the set s, or MinInt if s is empty.

#### func (*Sparse) Min

```go
func (s *Sparse) Min() int
```
Min returns the minimum element of the set s, or MaxInt if s is empty.

#### func (*Sparse) Remove

```go
func (s *Sparse) Remove(x int) bool
```
Remove removes x from the set s, and reports whether the set shrank.

#### func (*Sparse) String

```go
func (s *Sparse) String() string
```
String returns a human-readable description of the set s.

#### func (*Sparse) TakeMin

```go
func (s *Sparse) TakeMin(p *int) bool
```
If set s is non-empty, TakeMin sets *p to the minimum element of the set s,
removes that element from the set and returns true. Otherwise, it returns false
and *p is undefined.

This method may be used for iteration over a worklist like so:

    var x int
    for worklist.TakeMin(&x) { use(x) }

#### func (*Sparse) Union

```go
func (s *Sparse) Union(x, y *Sparse)
```
Union sets s to the union x ∪ y.

#### func (*Sparse) UnionWith

```go
func (s *Sparse) UnionWith(x *Sparse) bool
```
UnionWith sets s to the union s ∪ x, and reports whether s grew.
