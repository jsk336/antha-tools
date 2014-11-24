# gatefs
--
    import "."

Package gatefs provides an implementation of the FileSystem interface that wraps
another FileSystem and limits its concurrency.

## Usage

#### func  New

```go
func New(fs vfs.FileSystem, gateCh chan bool) vfs.FileSystem
```
New returns a new FileSystem that delegates to fs. If gateCh is non-nil and
buffered, it's used as a gate to limit concurrency on calls to fs.
