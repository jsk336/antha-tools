# astutil
--
    import "."

Package astutil contains common utilities for working with the Go AST.

## Usage

#### func  AddImport

```go
func AddImport(fset *token.FileSet, f *ast.File, ipath string) (added bool)
```
AddImport adds the import path to the file f, if absent.

#### func  AddNamedImport

```go
func AddNamedImport(fset *token.FileSet, f *ast.File, name, ipath string) (added bool)
```
AddNamedImport adds the import path to the file f, if absent. If name is not
empty, it is used to rename the import.

For example, calling

    AddNamedImport(fset, f, "pathpkg", "path")

adds

    import pathpkg "path"

#### func  DeleteImport

```go
func DeleteImport(fset *token.FileSet, f *ast.File, path string) (deleted bool)
```
DeleteImport deletes the import path from the file f, if present.

#### func  Imports

```go
func Imports(fset *token.FileSet, f *ast.File) [][]*ast.ImportSpec
```
Imports returns the file imports grouped by paragraph.

#### func  NodeDescription

```go
func NodeDescription(n ast.Node) string
```
NodeDescription returns a description of the concrete type of n suitable for a
user interface.

TODO(adonovan): in some cases (e.g. Field, FieldList, Ident, StarExpr) we could
be much more specific given the path to the AST root. Perhaps we should do that.

#### func  PathEnclosingInterval

```go
func PathEnclosingInterval(root *ast.File, start, end token.Pos) (path []ast.Node, exact bool)
```
PathEnclosingInterval returns the node that encloses the source interval [start,
end), and all its ancestors up to the AST root.

The definition of "enclosing" used by this function considers additional
whitespace abutting a node to be enclosed by it. In this example:

    z := x + y // add them
         <-A->
        <----B----->

the ast.BinaryExpr(+) node is considered to enclose interval B even though its
[Pos()..End()) is actually only interval A. This behaviour makes user interfaces
more tolerant of imperfect input.

This function treats tokens as nodes, though they are not included in the
result. e.g. PathEnclosingInterval("+") returns the enclosing ast.BinaryExpr("x
+ y").

If start==end, the 1-char interval following start is used instead.

The 'exact' result is true if the interval contains only path[0] and perhaps
some adjacent whitespace. It is false if the interval overlaps multiple children
of path[0], or if it contains only interior whitespace of path[0]. In this
example:

    z := x + y // add them
      <--C-->     <---E-->
        ^
        D

intervals C, D and E are inexact. C is contained by the z-assignment statement,
because it spans three of its children (:=, x, +). So too is the 1-char interval
D, because it contains only interior whitespace of the assignment. E is
considered interior whitespace of the BlockStmt containing the assignment.

Precondition: [start, end) both lie within the same file as root.
TODO(adonovan): return (nil, false) in this case and remove precond. Requires
FileSet; see loader.tokenFileContainsPos.

Postcondition: path is never nil; it always contains at least 'root'.

#### func  RenameTop

```go
func RenameTop(f *ast.File, old, new string) bool
```
RenameTop renames all references to the top-level name old. It returns true if
it makes any changes.

#### func  RewriteImport

```go
func RewriteImport(fset *token.FileSet, f *ast.File, oldPath, newPath string) (rewrote bool)
```
RewriteImport rewrites any import of path oldPath to path newPath.

#### func  UsesImport

```go
func UsesImport(f *ast.File, path string) (used bool)
```
UsesImport reports whether a given import is used.
