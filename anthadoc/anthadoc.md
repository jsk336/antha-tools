# godoc
--
    import "."

Copyright 2013 The Go Authors. All rights reserved. Use of this source code is
governed by a BSD-style license that can be found in the LICENSE file.

Package godoc is a work-in-progress (2013-07-17) package to begin splitting up
the godoc binary into multiple pieces.

This package comment will evolve over time as this package splits into smaller
pieces.

## Usage

```go
var ErrFileIndexVersion = errors.New("file index version out of date")
```

#### func  CommandLine

```go
func CommandLine(w io.Writer, fs vfs.NameSpace, pres *Presentation, args []string) error
```
CommandLine returns godoc results to w. Note that it may add a /target path to
fs.

#### func  FormatSelections

```go
func FormatSelections(w io.Writer, text []byte, lw LinkWriter, links Selection, sw SegmentWriter, selections ...Selection)
```
FormatSelections takes a text and writes it to w using link and segment writers
lw and sw as follows: lw is invoked for consecutive segment starts and ends as
specified through the links selection, and sw is invoked for consecutive
segments of text overlapped by the same selections as specified by selections.
The link writer lw may be nil, in which case the links Selection is ignored.

#### func  FormatText

```go
func FormatText(w io.Writer, text []byte, line int, goSource bool, pattern string, selection Selection)
```
FormatText HTML-escapes text and writes it to w. Consecutive text segments are
wrapped in HTML spans (with tags as defined by startTags and endTag) as follows:

    - if line >= 0, line number (ln) spans are inserted before each line,
      starting with the value of line
    - if the text is Go source, comments get the "comment" span class
    - each occurrence of the regular expression pattern gets the "highlight"
      span class
    - text segments covered by selection get the "selection" span class

Comments, highlights, and selections may overlap arbitrarily; the respective
HTML span classes are specified in the startTags variable.

#### func  Linkify

```go
func Linkify(out io.Writer, src []byte)
```

#### func  LinkifyText

```go
func LinkifyText(w io.Writer, text []byte, n ast.Node)
```
LinkifyText HTML-escapes source text and writes it to w. Identifiers that are in
a "use" position (i.e., that are not being declared), are wrapped with HTML
links pointing to the respective declaration, if possible. Comments are
formatted the same way as with FormatText.

#### type AltWords

```go
type AltWords struct {
	Canon string   // canonical word spelling (all lowercase)
	Alts  []string // alternative spelling for the same word
}
```

An AltWords describes a list of alternative spellings for a canonical (all
lowercase) spelling of a word.

#### type Corpus

```go
type Corpus struct {

	// Verbose logging.
	Verbose bool

	// IndexEnabled controls whether indexing is enabled.
	IndexEnabled bool

	// IndexFiles specifies a glob pattern specifying index files.
	// If not empty, the index is read from these files in sorted
	// order.
	IndexFiles string

	// IndexThrottle specifies the indexing throttle value
	// between 0.0 and 1.0. At 0.0, the indexer always sleeps.
	// At 1.0, the indexer never sleeps. Because 0.0 is useless
	// and redundant with setting IndexEnabled to false, the
	// zero value for IndexThrottle means 0.9.
	IndexThrottle float64

	// IndexInterval specifies the time to sleep between reindexing
	// all the sources.
	// If zero, a default is used. If negative, the index is only
	// built once.
	IndexInterval time.Duration

	// IndexDocs enables indexing of Go documentation.
	// This will produce search results for exported types, functions,
	// methods, variables, and constants, and will link to the godoc
	// documentation for those identifiers.
	IndexDocs bool

	// IndexGoCode enables indexing of Go source code.
	// This will produce search results for internal and external identifiers
	// and will link to both declarations and uses of those identifiers in
	// source code.
	IndexGoCode bool

	// IndexFullText enables full-text indexing.
	// This will provide search results for any matching text in any file that
	// is indexed, including non-Go files (see whitelisted in index.go).
	// Regexp searching is supported via full-text indexing.
	IndexFullText bool

	// MaxResults optionally specifies the maximum results for indexing.
	MaxResults int

	// SummarizePackage optionally specifies a function to
	// summarize a package. It exists as an optimization to
	// avoid reading files to parse package comments.
	//
	// If SummarizePackage returns false for ok, the caller
	// ignores all return values and parses the files in the package
	// as if SummarizePackage were nil.
	//
	// If showList is false, the package is hidden from the
	// package listing.
	SummarizePackage func(pkg string) (summary string, showList, ok bool)

	// IndexDirectory optionally specifies a function to determine
	// whether the provided directory should be indexed.  The dir
	// will be of the form "/src/cmd/6a", "/doc/play",
	// "/src/pkg/io", etc.
	// If nil, all directories are indexed if indexing is enabled.
	IndexDirectory func(dir string) bool

	// Analysis is the result of type and pointer analysis.
	Analysis analysis.Result
}
```

A Corpus holds all the state related to serving and indexing a collection of Go
code.

Construct a new Corpus with NewCorpus, then modify options, then call its Init
method.

#### func  NewCorpus

```go
func NewCorpus(fs vfs.FileSystem) *Corpus
```
NewCorpus returns a new Corpus from a filesystem. The returned corpus has all
indexing enabled and MaxResults set to 1000. Change or set any options on Corpus
before calling the Corpus.Init method.

#### func (*Corpus) CurrentIndex

```go
func (c *Corpus) CurrentIndex() (*Index, time.Time)
```

#### func (*Corpus) FSModifiedTime

```go
func (c *Corpus) FSModifiedTime() time.Time
```

#### func (*Corpus) Init

```go
func (c *Corpus) Init() error
```
Init initializes Corpus, once options on Corpus are set. It must be called
before any subsequent method calls.

#### func (*Corpus) Lookup

```go
func (c *Corpus) Lookup(query string) SearchResult
```

#### func (*Corpus) MetadataFor

```go
func (c *Corpus) MetadataFor(relpath string) *Metadata
```
MetadataFor returns the *Metadata for a given relative path or nil if none
exists.

#### func (*Corpus) NewIndex

```go
func (c *Corpus) NewIndex() *Index
```
NewIndex creates a new index for the .go files provided by the corpus.

#### func (*Corpus) RunIndexer

```go
func (c *Corpus) RunIndexer()
```
RunIndexer runs forever, indexing.

#### func (*Corpus) UpdateIndex

```go
func (c *Corpus) UpdateIndex()
```

#### type DirEntry

```go
type DirEntry struct {
	Depth    int    // >= 0
	Height   int    // = DirList.MaxHeight - Depth, > 0
	Path     string // directory path; includes Name, relative to DirList root
	Name     string // directory name
	HasPkg   bool   // true if the directory contains at least one package
	Synopsis string // package documentation, if any
}
```

DirEntry describes a directory entry. The Depth and Height values are useful for
presenting an entry in an indented fashion.

#### type DirList

```go
type DirList struct {
	MaxHeight int // directory tree height, > 0
	List      []DirEntry
}
```


#### type Directory

```go
type Directory struct {
	Depth    int
	Path     string       // directory path; includes Name
	Name     string       // directory name
	HasPkg   bool         // true if the directory contains at least one package
	Synopsis string       // package documentation, if any
	Dirs     []*Directory // subdirectories
}
```


#### type File

```go
type File struct {
	Name string // directory-local file name
	Pak  *Pak   // the package to which the file belongs
}
```

A File describes a Go file.

#### func (*File) Path

```go
func (f *File) Path() string
```
Path returns the file path of f.

#### type FileLines

```go
type FileLines struct {
	Filename string
	Lines    []int
}
```

A FileLines value specifies a file and line numbers within that file.

#### type FileRun

```go
type FileRun struct {
	File   *File
	Groups []KindRun
}
```

A FileRun is a list of KindRuns belonging to the same file.

#### type HitList

```go
type HitList []*PakRun
```

A HitList describes a list of PakRuns.

#### type Ident

```go
type Ident struct {
	Path    string // e.g. "net/http"
	Package string // e.g. "http"
	Name    string // e.g. "NewRequest"
	Doc     string // e.g. "NewRequest returns a new Request..."
}
```

Ident stores information about external identifiers in order to create links to
package documentation.

#### type Index

```go
type Index struct {
}
```


#### func (*Index) CompatibleWith

```go
func (x *Index) CompatibleWith(c *Corpus) bool
```
CompatibleWith reports whether the Index x is compatible with the corpus
indexing options set in c.

#### func (*Index) Exports

```go
func (x *Index) Exports() map[string]map[string]SpotKind
```
Exports returns a map from full package path to exported symbol name to its
type.

#### func (*Index) Idents

```go
func (x *Index) Idents() map[SpotKind]map[string][]Ident
```
Idents returns a map from identifier type to exported symbol name to the list of
identifiers matching that name.

#### func (*Index) ImportCount

```go
func (x *Index) ImportCount() map[string]int
```
ImportCount returns a map from import paths to how many times they were seen.

#### func (*Index) Lookup

```go
func (x *Index) Lookup(query string) (*SearchResult, error)
```
For a given query, which is either a single identifier or a qualified
identifier, Lookup returns a SearchResult containing packages, a LookupResult, a
list of alternative spellings, and identifiers, if any. Any and all results may
be nil. If the query syntax is wrong, an error is reported.

#### func (*Index) LookupRegexp

```go
func (x *Index) LookupRegexp(r *regexp.Regexp, n int) (found int, result []FileLines)
```
LookupRegexp returns the number of matches and the matches where a regular
expression r is found in the full text index. At most n matches are returned
(thus found <= n).

#### func (*Index) PackagePath

```go
func (x *Index) PackagePath() map[string]map[string]bool
```
PackagePath returns a map from short package name to a set of full package path
names that use that short package name.

#### func (*Index) ReadFrom

```go
func (x *Index) ReadFrom(r io.Reader) (n int64, err error)
```
ReadFrom reads the index from r into x; x must not be nil. If r does not also
implement io.ByteReader, it will be wrapped in a bufio.Reader. If the index is
from an old version, the error is ErrFileIndexVersion.

#### func (*Index) Snippet

```go
func (x *Index) Snippet(i int) *Snippet
```

#### func (*Index) Stats

```go
func (x *Index) Stats() Statistics
```
Stats returns index statistics.

#### func (*Index) WriteTo

```go
func (x *Index) WriteTo(w io.Writer) (n int64, err error)
```
WriteTo writes the index x to w.

#### type IndexResult

```go
type IndexResult struct {
	Decls  RunList // package-level declarations (with snippets)
	Others RunList // all other occurrences
}
```


#### type Indexer

```go
type Indexer struct {
}
```

An Indexer maintains the data structures and provides the machinery for indexing
.go files under a file tree. It implements the path.Visitor interface for
walking file trees, and the ast.Visitor interface for walking Go ASTs.

#### func (*Indexer) Visit

```go
func (x *Indexer) Visit(node ast.Node) ast.Visitor
```

#### type KindRun

```go
type KindRun []SpotInfo
```

A KindRun is a run of SpotInfos of the same kind in a given file. The kind (3
bits) is stored in each SpotInfo element; to find the kind of a KindRun, look at
any of its elements.

#### func (KindRun) Len

```go
func (k KindRun) Len() int
```
KindRuns are sorted by line number or index. Since the isIndex bit is always the
same for all infos in one list we can compare lori's.

#### func (KindRun) Less

```go
func (k KindRun) Less(i, j int) bool
```

#### func (KindRun) Swap

```go
func (k KindRun) Swap(i, j int)
```

#### type LinkWriter

```go
type LinkWriter func(w io.Writer, offs int, start bool)
```

A LinkWriter writes some start or end "tag" to w for the text offset offs. It is
called by FormatSelections at the start or end of each link segment.

#### type LookupResult

```go
type LookupResult struct {
	Decls  HitList // package-level declarations (with snippets)
	Others HitList // all other occurrences
}
```


#### type Metadata

```go
type Metadata struct {
	Title    string
	Subtitle string
	Template bool   // execute as template
	Path     string // canonical path for this page
}
```

TODO(adg): why are some exported and some aren't? -brad

#### func (*Metadata) FilePath

```go
func (m *Metadata) FilePath() string
```

#### type Page

```go
type Page struct {
	Title    string
	Tabtitle string
	Subtitle string
	Query    string
	Body     []byte

	// filled in by servePage
	SearchBox  bool
	Playground bool
	Version    string
}
```

Page describes the contents of the top-level godoc webpage.

#### type PageInfo

```go
type PageInfo struct {
	Dirname string // directory containing the package
	Err     error  // error or nil

	// package info
	FSet       *token.FileSet         // nil if no package documentation
	PDoc       *doc.Package           // nil if no package documentation
	Examples   []*doc.Example         // nil if no example code
	Notes      map[string][]*doc.Note // nil if no package Notes
	PAst       map[string]*ast.File   // nil if no AST with package exports
	IsMain     bool                   // true for package main
	IsFiltered bool                   // true if results were filtered

	// analysis info
	TypeInfoIndex  map[string]int  // index of JSON datum for type T (if -analysis=type)
	AnalysisData   htmltemplate.JS // array of TypeInfoJSON values
	CallGraph      htmltemplate.JS // array of PCGNodeJSON values    (if -analysis=pointer)
	CallGraphIndex map[string]int  // maps func name to index in CallGraph

	// directory info
	Dirs    *DirList  // nil if no directory information
	DirTime time.Time // directory time stamp
	DirFlat bool      // if set, show directory in a flat (non-indented) manner
}
```


#### func (*PageInfo) IsEmpty

```go
func (info *PageInfo) IsEmpty() bool
```

#### type PageInfoMode

```go
type PageInfoMode uint
```


```go
const (
	NoFiltering PageInfoMode = 1 << iota // do not filter exports
	AllMethods                           // show all embedded methods
	ShowSource                           // show source code, do not extract documentation
	NoHTML                               // show result in textual form, do not generate HTML
	FlatDir                              // show directory in a flat (non-indented) manner
	NoTypeAssoc                          // don't associate consts, vars, and factory functions with types
)
```

#### type Pak

```go
type Pak struct {
	Path string // path of directory containing the package
	Name string // package name as declared by package clause
}
```

A Pak describes a Go package.

#### type PakRun

```go
type PakRun struct {
	Pak   *Pak
	Files []*FileRun
}
```

A PakRun describes a run of *FileRuns of a package.

#### func (*PakRun) Len

```go
func (p *PakRun) Len() int
```
Sorting support for files within a PakRun.

#### func (*PakRun) Less

```go
func (p *PakRun) Less(i, j int) bool
```

#### func (*PakRun) Swap

```go
func (p *PakRun) Swap(i, j int)
```

#### type Presentation

```go
type Presentation struct {
	Corpus *Corpus

	CallGraphHTML,
	DirlistHTML,
	ErrorHTML,
	ExampleHTML,
	GodocHTML,
	ImplementsHTML,
	MethodSetHTML,
	PackageHTML,
	PackageText,
	SearchHTML,
	SearchDocHTML,
	SearchCodeHTML,
	SearchTxtHTML,
	SearchText,
	SearchDescXML *template.Template

	// TabWidth optionally specifies the tab width.
	TabWidth int

	ShowTimestamps bool
	ShowPlayground bool
	ShowExamples   bool
	DeclLinks      bool

	// SrcMode outputs source code instead of documentation in command-line mode.
	SrcMode bool
	// HTMLMode outputs HTML instead of plain text in command-line mode.
	HTMLMode bool

	// NotesRx optionally specifies a regexp to match
	// notes to render in the output.
	NotesRx *regexp.Regexp

	// AdjustPageInfoMode optionally specifies a function to
	// modify the PageInfoMode of a request. The default chosen
	// value is provided.
	AdjustPageInfoMode func(req *http.Request, mode PageInfoMode) PageInfoMode

	// URLForSrc optionally specifies a function that takes a source file and
	// returns a URL for it.
	// The source file argument has the form /src/pkg/<path>/<filename>.
	URLForSrc func(src string) string

	// URLForSrcPos optionally specifies a function to create a URL given a
	// source file, a line from the source file (1-based), and low & high offset
	// positions (0-based, bytes from beginning of file). Ideally, the returned
	// URL will be for the specified line of the file, while the high & low
	// positions will be used to highlight a section of the file.
	// The source file argument has the form /src/pkg/<path>/<filename>.
	URLForSrcPos func(src string, line, low, high int) string

	// URLForSrcQuery optionally specifies a function to create a URL given a
	// source file, a query string, and a line from the source file (1-based).
	// The source file argument has the form /src/pkg/<path>/<filename>.
	// The query argument will be escaped for the purposes of embedding in a URL
	// query parameter.
	// Ideally, the returned URL will be for the specified line of the file with
	// the query string highlighted.
	URLForSrcQuery func(src, query string, line int) string

	// SearchResults optionally specifies a list of functions returning an HTML
	// body for displaying search results.
	SearchResults []SearchResultFunc
}
```

Presentation generates output from a corpus.

#### func  NewPresentation

```go
func NewPresentation(c *Corpus) *Presentation
```
NewPresentation returns a new Presentation from a corpus. It sets SearchResults
to: [SearchResultDoc SearchResultCode SearchResultTxt].

#### func (*Presentation) CmdFSRoot

```go
func (p *Presentation) CmdFSRoot() string
```

#### func (*Presentation) FileServer

```go
func (p *Presentation) FileServer() http.Handler
```

#### func (*Presentation) FuncMap

```go
func (p *Presentation) FuncMap() template.FuncMap
```
FuncMap defines template functions used in godoc templates.

Convention: template function names ending in "_html" or "_url" produce

    HTML- or URL-escaped strings; all other function results may
    require explicit escaping in the template.

#### func (*Presentation) GetCmdPageInfo

```go
func (p *Presentation) GetCmdPageInfo(abspath, relpath string, mode PageInfoMode) *PageInfo
```
TODO(bradfitz): move this to be a method on Corpus. Just moving code around for
now, but this doesn't feel right.

#### func (*Presentation) GetPageInfoMode

```go
func (p *Presentation) GetPageInfoMode(r *http.Request) PageInfoMode
```
GetPageInfoMode computes the PageInfoMode flags by analyzing the request URL
form value "m". It is value is a comma-separated list of mode names as defined
by modeNames (e.g.: m=src,text).

#### func (*Presentation) GetPkgPageInfo

```go
func (p *Presentation) GetPkgPageInfo(abspath, relpath string, mode PageInfoMode) *PageInfo
```
TODO(bradfitz): move this to be a method on Corpus. Just moving code around for
now, but this doesn't feel right.

#### func (*Presentation) HandleSearch

```go
func (p *Presentation) HandleSearch(w http.ResponseWriter, r *http.Request)
```
HandleSearch obtains results for the requested search and returns a page to
display them.

#### func (*Presentation) NewSnippet

```go
func (p *Presentation) NewSnippet(fset *token.FileSet, decl ast.Decl, id *ast.Ident) *Snippet
```
NewSnippet creates a text snippet from a declaration decl containing an
identifier id. Parts of the declaration not containing the identifier may be
removed for a more compact snippet.

#### func (*Presentation) PkgFSRoot

```go
func (p *Presentation) PkgFSRoot() string
```

#### func (*Presentation) SearchResultCode

```go
func (p *Presentation) SearchResultCode(result SearchResult) []byte
```
SearchResultCode optionally specifies a function returning an HTML body
displaying search results matching source code.

#### func (*Presentation) SearchResultDoc

```go
func (p *Presentation) SearchResultDoc(result SearchResult) []byte
```
SearchResultDoc optionally specifies a function returning an HTML body
displaying search results matching godoc documentation.

#### func (*Presentation) SearchResultTxt

```go
func (p *Presentation) SearchResultTxt(result SearchResult) []byte
```
SearchResultTxt optionally specifies a function returning an HTML body
displaying search results of textual matches.

#### func (*Presentation) ServeError

```go
func (p *Presentation) ServeError(w http.ResponseWriter, r *http.Request, relpath string, err error)
```

#### func (*Presentation) ServeFile

```go
func (p *Presentation) ServeFile(w http.ResponseWriter, r *http.Request)
```

#### func (*Presentation) ServeHTMLDoc

```go
func (p *Presentation) ServeHTMLDoc(w http.ResponseWriter, r *http.Request, abspath, relpath string)
```

#### func (*Presentation) ServeHTTP

```go
func (p *Presentation) ServeHTTP(w http.ResponseWriter, r *http.Request)
```

#### func (*Presentation) ServePage

```go
func (p *Presentation) ServePage(w http.ResponseWriter, page Page)
```

#### func (*Presentation) ServeText

```go
func (p *Presentation) ServeText(w http.ResponseWriter, text []byte)
```

#### func (*Presentation) TemplateFuncs

```go
func (p *Presentation) TemplateFuncs() template.FuncMap
```

#### func (*Presentation) WriteNode

```go
func (p *Presentation) WriteNode(w io.Writer, fset *token.FileSet, x interface{})
```
WriteNode writes x to w. TODO(bgarcia) Is this method needed? It's just a
wrapper for p.writeNode.

#### type RunList

```go
type RunList []interface{}
```

A RunList is a list of entries that can be sorted according to some criteria. A
RunList may be compressed by grouping "runs" of entries which are equal
(according to the sort critera) into a new RunList of runs. For instance, a
RunList containing pairs (x, y) may be compressed into a RunList containing pair
runs (x, {y}) where each run consists of a list of y's with the same x.

#### type SearchResult

```go
type SearchResult struct {
	Query string
	Alert string // error or warning message

	// identifier matches
	Pak HitList       // packages matching Query
	Hit *LookupResult // identifier matches of Query
	Alt *AltWords     // alternative identifiers to look for

	// textual matches
	Found    int         // number of textual occurrences found
	Textual  []FileLines // textual matches of Query
	Complete bool        // true if all textual occurrences of Query are reported
	Idents   map[SpotKind][]Ident
}
```


#### type SearchResultFunc

```go
type SearchResultFunc func(p *Presentation, result SearchResult) []byte
```

SearchResultFunc functions return an HTML body for displaying search results.

#### type Segment

```go
type Segment struct {
}
```

A Segment describes a text segment [start, end). The zero value of a Segment is
a ready-to-use empty segment.

#### type SegmentWriter

```go
type SegmentWriter func(w io.Writer, text []byte, selections int)
```

A SegmentWriter formats a text according to selections and writes it to w. The
selections parameter is a bit set indicating which selections provided to
FormatSelections overlap with the text segment: If the n'th bit is set in
selections, the n'th selection provided to FormatSelections is overlapping with
the text.

#### type Selection

```go
type Selection func() Segment
```

A Selection is an "iterator" function returning a text segment. Repeated calls
to a selection return consecutive, non-overlapping, non-empty segments, followed
by an infinite sequence of empty segments. The first empty segment marks the end
of the selection.

#### func  RangeSelection

```go
func RangeSelection(str string) Selection
```
RangeSelection computes the Selection for a text range described by the argument
str; the range description must match the selRx regular expression.

#### type Snippet

```go
type Snippet struct {
	Line int
	Text string // HTML-escaped
}
```


#### func  NewSnippet

```go
func NewSnippet(fset *token.FileSet, decl ast.Decl, id *ast.Ident) *Snippet
```
NewSnippet creates a text snippet from a declaration decl containing an
identifier id. Parts of the declaration not containing the identifier may be
removed for a more compact snippet.

#### type Spot

```go
type Spot struct {
	File *File
	Info SpotInfo
}
```

A Spot describes a single occurrence of a word.

#### type SpotInfo

```go
type SpotInfo uint32
```

A SpotInfo value describes a particular identifier spot in a given file; It
encodes three values: the SpotKind (declaration or use), a line or snippet index
"lori", and whether it's a line or index.

The following encoding is used:

    bits    32   4    1       0
    value    [lori|kind|isIndex]

#### func (SpotInfo) IsIndex

```go
func (x SpotInfo) IsIndex() bool
```

#### func (SpotInfo) Kind

```go
func (x SpotInfo) Kind() SpotKind
```

#### func (SpotInfo) Lori

```go
func (x SpotInfo) Lori() int
```

#### type SpotKind

```go
type SpotKind uint32
```

SpotKind describes whether an identifier is declared (and what kind of
declaration) or used.

```go
const (
	PackageClause SpotKind = iota
	ImportDecl
	ConstDecl
	TypeDecl
	VarDecl
	FuncDecl
	MethodDecl
	Use
)
```

#### func (SpotKind) Name

```go
func (x SpotKind) Name() string
```

#### type Statistics

```go
type Statistics struct {
	Bytes int // total size of indexed source files
	Files int // number of indexed source files
	Lines int // number of lines (all files)
	Words int // number of different identifiers
	Spots int // number of identifier occurrences
}
```

Statistics provides statistics information for an index.
