# cache
--
    import "."


## Usage

```go
var TimeKey = "cachetime"
```
TimeKey specifies the memcache entity that keeps the logical datastore time.

#### func  Get

```go
func Get(r *http.Request, now uint64, name string, value interface{}) bool
```
Get fetches data for name at time now from memcache and unmarshals it into
value. It reports whether it found the cache record and logs any errors to the
admin console.

#### func  Now

```go
func Now(c appengine.Context) uint64
```
Now returns the current logical datastore time to use for cache lookups.

#### func  Set

```go
func Set(r *http.Request, now uint64, name string, value interface{})
```
Set puts value into memcache under name at time now. It logs any errors to the
admin console.

#### func  Tick

```go
func Tick(c appengine.Context) uint64
```
Tick sets the current logical datastore time to a never-before-used time and
returns that time. It should be called to invalidate the cache.
