# redirect
--
    import "."

Package redirect provides hooks to register HTTP handlers that redirect old
godoc paths to their new equivalents and assist in accessing the issue tracker,
wiki, code review system, etc.

## Usage

#### func  Handler

```go
func Handler(target string) http.Handler
```

#### func  PrefixHandler

```go
func PrefixHandler(prefix, baseURL string) http.Handler
```

#### func  Register

```go
func Register(mux *http.ServeMux)
```
Register registers HTTP handlers that redirect old godoc paths to their new
equivalents and assist in accessing the issue tracker, wiki, code review system,
etc. If mux is nil it uses http.DefaultServeMux.
