# blog
--
    import "."

Package blog implements a web server for articles written in present format.

## Usage

#### type Config

```go
type Config struct {
	ContentPath  string // Relative or absolute location of article files and related content.
	TemplatePath string // Relative or absolute location of template files.

	BaseURL  string // Absolute base URL (for permalinks; no trailing slash).
	BasePath string // Base URL path relative to server root (no trailing slash).
	GodocURL string // The base URL of godoc (for menu bar; no trailing slash).
	Hostname string // Server host name, used for rendering ATOM feeds.

	HomeArticles int    // Articles to display on the home page.
	FeedArticles int    // Articles to include in Atom and JSON feeds.
	FeedTitle    string // The title of the Atom XML feed

	PlayEnabled bool
}
```

Config specifies Server configuration values.

#### type Doc

```go
type Doc struct {
	*present.Doc
	Permalink string        // Canonical URL for this document.
	Path      string        // Path relative to server root (including base).
	HTML      template.HTML // rendered article

	Related      []*Doc
	Newer, Older *Doc
}
```

Doc represents an article adorned with presentation data.

#### type Server

```go
type Server struct {
}
```

Server implements an http.Handler that serves blog articles.

#### func  NewServer

```go
func NewServer(cfg Config) (*Server, error)
```
NewServer constructs a new Server using the specified config.

#### func (*Server) ServeHTTP

```go
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request)
```
ServeHTTP serves the front, index, and article pages as well as the ATOM and
JSON feeds.
