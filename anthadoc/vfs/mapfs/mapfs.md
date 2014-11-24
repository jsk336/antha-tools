# mapfs
--
    import "."

Package mapfs file provides an implementation of the FileSystem interface based
on the contents of a map[string]string.

## Usage

#### func  New

```go
func New(m map[string]string) vfs.FileSystem
```
New returns a new FileSystem from the provided map. Map keys should be forward
slash-separated pathnames and not contain a leading slash.
