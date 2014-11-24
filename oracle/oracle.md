# oracle
--
    import "."

Package oracle contains the implementation of the oracle tool whose command-line
is provided by antha-tools/cmd/oracle.

http://golang.org/s/oracle-design http://golang.org/s/oracle-user-manual

## Usage

#### type Oracle

```go
type Oracle struct {
}
```

An Oracle holds the program state required for one or more queries.

#### func  New

```go
func New(iprog *loader.Program, ptalog io.Writer, reflection bool) (*Oracle, error)
```
New constructs a new Oracle that can be used for a sequence of queries.

iprog specifies the program to analyze. ptalog is the (optional)
pointer-analysis log file. reflection determines whether to model reflection
soundly (currently slow).

#### func (*Oracle) Query

```go
func (o *Oracle) Query(mode string, qpos *QueryPos) (*Result, error)
```
Query runs the query of the specified mode and selection.

TODO(adonovan): fix: this function does not currently support the "what" query,
which needs to access the antha/build.Context.

#### type QueryPos

```go
type QueryPos struct {
}
```

A QueryPos represents the position provided as input to a query: a textual
extent in the program's source code, the AST node it corresponds to, and the
package to which it belongs. Instances are created by ParseQueryPos.

#### func  ParseQueryPos

```go
func ParseQueryPos(iprog *loader.Program, posFlag string, needExact bool) (*QueryPos, error)
```
ParseQueryPos parses the source query position pos. If needExact, it must
identify a single AST subtree; this is appropriate for queries that allow fairly
arbitrary syntax, e.g. "describe".

#### func (*QueryPos) ObjectString

```go
func (qpos *QueryPos) ObjectString(obj types.Object) string
```
ObjectString prints object obj relative to the query position.

#### func (*QueryPos) SelectionString

```go
func (qpos *QueryPos) SelectionString(sel *types.Selection) string
```
SelectionString prints selection sel relative to the query position.

#### func (*QueryPos) TypeString

```go
func (qpos *QueryPos) TypeString(T types.Type) string
```
TypeString prints type T relative to the query position.

#### type Result

```go
type Result struct {
}
```

A Result encapsulates the result of an oracle.Query.

#### func  Query

```go
func Query(args []string, mode, pos string, ptalog io.Writer, buildContext *build.Context, reflection bool) (*Result, error)
```
Query runs a single oracle query.

args specify the main package in (*loader.Config).FromArgs syntax. mode is the
query mode ("callers", etc). ptalog is the (optional) pointer-analysis log file.
buildContext is the antha/build configuration for locating packages. reflection
determines whether to model reflection soundly (currently slow).

Clients that intend to perform multiple queries against the same analysis scope
should use this pattern instead:

    conf := loader.Config{Build: buildContext, SourceImports: true}
    ... populate config, e.g. conf.FromArgs(args) ...
    iprog, err := conf.Load()
    if err != nil { ... }
    o, err := oracle.New(iprog, nil, false)
    if err != nil { ... }
    for ... {
    	qpos, err := oracle.ParseQueryPos(imp, pos, needExact)
    	if err != nil { ... }

    	res, err := o.Query(mode, qpos)
    	if err != nil { ... }

    	// use res
    }

TODO(adonovan): the ideal 'needsExact' parameter for ParseQueryPos depends on
the query mode; how should we expose this?

#### func (*Result) Serial

```go
func (res *Result) Serial() *serial.Result
```
Serial returns an instance of serial.Result, which implements the
{xml,json}.Marshaler interfaces so that query results can be serialized as JSON
or XML.

#### func (*Result) WriteTo

```go
func (res *Result) WriteTo(out io.Writer)
```
WriteTo writes the oracle query result res to out in a compiler diagnostic
format.
