# analysis
--
    import "."

Package analysis performs type and pointer analysis and generates mark-up for
the Go source view.

The Run method populates a Result object by running type and (optionally)
pointer analysis. The Result object is thread-safe and at all times may be
accessed by a serving thread, even as it is progressively populated as analysis
facts are derived.

The Result is a mapping from each godoc file URL (e.g. /src/pkg/fmt/print.go) to
information about that file. The information is a list of HTML markup links and
a JSON array of structured data values. Some of the links call client-side
JavaScript functions that index this array.

The analysis computes mark-up for the following relations:

IMPORTS: for each ast.ImportSpec, the package that it denotes.

RESOLUTION: for each ast.Ident, its kind and type, and the location of its
definition.

METHOD SETS, IMPLEMENTS: for each ast.Ident defining a named type, its
method-set, the set of interfaces it implements or is implemented by, and its
size/align values.

CALLERS, CALLEES: for each function declaration ('func' token), its callers, and
for each call-site ('(' token), its callees.

CALLGRAPH: the package docs include an interactive viewer for the intra-package
call graph of "fmt".

CHANNEL PEERS: for each channel operation make/<-/close, the set of other
channel ops that alias the same channel(s).

ERRORS: for each locus of a static (antha/types) error, the location is
highlighted in red and hover text provides the compiler error message.

## Usage

#### func  Run

```go
func Run(pta bool, result *Result)
```
Run runs program analysis and computes the resulting markup, populating *result
in a thread-safe manner, first with type information then later with pointer
analysis information if enabled by the pta flag.

#### type Link

```go
type Link interface {
	Start() int
	End() int
	Write(w io.Writer, _ int, start bool) // the godoc.LinkWriter signature
}
```

A Link is an HTML decoration of the bytes [Start, End) of a file. Write is
called before/after those bytes to emit the mark-up.

#### type PCGNodeJSON

```go
type PCGNodeJSON struct {
	Func    anchorJSON
	Callees []int // indices within CALLGRAPH of nodes called by this one
}
```

JavaScript's cgAddChild requires a global array of PCGNodeJSON called CALLGRAPH,
representing the intra-package call graph. The first element is special and
represents "all external callers".

#### type Result

```go
type Result struct {
}
```

Result contains the results of analysis. The result contains a mapping from
filenames to a set of HTML links and JavaScript data referenced by the links.

#### func (*Result) FileInfo

```go
func (res *Result) FileInfo(url string) ([]interface{}, []Link)
```
FileInfo returns new slices containing opaque JSON values and the HTML link
markup for the specified godoc file URL. Thread-safe. Callers must not mutate
the elements. It returns "zero" if no data is available.

#### func (*Result) PackageInfo

```go
func (res *Result) PackageInfo(importPath string) ([]*PCGNodeJSON, map[string]int, []*TypeInfoJSON)
```
PackageInfo returns new slices of JSON values for the callgraph and type info
for the specified package. Thread-safe. Callers must not mutate the elements.
PackageInfo returns "zero" if no data is available.

#### type TypeInfoJSON

```go
type TypeInfoJSON struct {
	Name        string // type name
	Size, Align int64
	Methods     []anchorJSON
	ImplGroups  []implGroupJSON
}
```

JavaScript's onClickIdent() expects a TypeInfoJSON.
