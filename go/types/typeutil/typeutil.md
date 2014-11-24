# typeutil
--
    import "."

Package typeutil defines various utilities for types, such as Map, a mapping
from types.Type to interface{} values.

## Usage

#### func  IntuitiveMethodSet

```go
func IntuitiveMethodSet(T types.Type, msets *types.MethodSetCache) []*types.Selection
```
IntuitiveMethodSet returns the intuitive method set of a type, T.

The result contains MethodSet(T) and additionally, if T is a concrete type,
methods belonging to *T if there is no identically named method on T itself.
This corresponds to user intuition about method sets; this function is intended
only for user interfaces.

The order of the result is as for types.MethodSet(T).

#### type Hasher

```go
type Hasher struct {
}
```

A Hasher maps each type to its hash value. For efficiency, a hasher uses
memoization; thus its memory footprint grows monotonically over time. Hashers
are not thread-safe. Hashers have reference semantics. Call MakeHasher to create
a Hasher.

#### func  MakeHasher

```go
func MakeHasher() Hasher
```
MakeHasher returns a new Hasher instance.

#### func (Hasher) Hash

```go
func (h Hasher) Hash(t types.Type) uint32
```
Hash computes a hash value for the given type t such that Identical(t, t') =>
Hash(t) == Hash(t').

#### type Map

```go
type Map struct {
}
```

Map is a hash-table-based mapping from types (types.Type) to arbitrary
interface{} values. The concrete types that implement the Type interface are
pointers. Since they are not canonicalized, == cannot be used to check for
equivalence, and thus we cannot simply use a Go map.

Just as with map[K]V, a nil *Map is a valid empty map.

Not thread-safe.

#### func (*Map) At

```go
func (m *Map) At(key types.Type) interface{}
```
At returns the map entry for the given key. The result is nil if the entry is
not present.

#### func (*Map) Delete

```go
func (m *Map) Delete(key types.Type) bool
```
Delete removes the entry with the given key, if any. It returns true if the
entry was found.

#### func (*Map) Iterate

```go
func (m *Map) Iterate(f func(key types.Type, value interface{}))
```
Iterate calls function f on each entry in the map in unspecified order.

If f should mutate the map, Iterate provides the same guarantees as Go maps: if
f deletes a map entry that Iterate has not yet reached, f will not be invoked
for it, but if f inserts a map entry that Iterate has not yet reached, whether
or not f will be invoked for it is unspecified.

#### func (*Map) Keys

```go
func (m *Map) Keys() []types.Type
```
Keys returns a new slice containing the set of map keys. The order is
unspecified.

#### func (*Map) KeysString

```go
func (m *Map) KeysString() string
```
KeysString returns a string representation of the map's key set. Order is
unspecified.

#### func (*Map) Len

```go
func (m *Map) Len() int
```
Len returns the number of map entries.

#### func (*Map) Set

```go
func (m *Map) Set(key types.Type, value interface{}) (prev interface{})
```
Set sets the map entry for key to val, and returns the previous entry, if any.

#### func (*Map) SetHasher

```go
func (m *Map) SetHasher(hasher Hasher)
```
SetHasher sets the hasher used by Map.

All Hashers are functionally equivalent but contain internal state used to cache
the results of hashing previously seen types.

A single Hasher created by MakeHasher() may be shared among many Maps. This is
recommended if the instances have many keys in common, as it will amortize the
cost of hash computation.

A Hasher may grow without bound as new types are seen. Even when a type is
deleted from the map, the Hasher never shrinks, since other types in the map may
reference the deleted type indirectly.

Hashers are not thread-safe, and read-only operations such as Map.Lookup require
updates to the hasher, so a full Mutex lock (not a read-lock) is require around
all Map operations if a shared hasher is accessed from multiple threads.

If SetHasher is not called, the Map will create a private hasher at the first
call to Insert.

#### func (*Map) String

```go
func (m *Map) String() string
```
String returns a string representation of the map's entries. Values are printed
using fmt.Sprintf("%v", v). Order is unspecified.
