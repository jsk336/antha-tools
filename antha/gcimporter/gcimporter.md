# gcimporter
--
    import "."

Package gcimporter implements Import for gc-generated object files. Importing
this package installs Import as antha/types.DefaultImport.

## Usage

#### func  FindExportData

```go
func FindExportData(r *bufio.Reader) (err error)
```
FindExportData positions the reader r at the beginning of the export data
section of an underlying GC-created object/archive file by reading from it. The
reader must be positioned at the start of the file before calling this function.

#### func  FindPkg

```go
func FindPkg(path, srcDir string) (filename, id string)
```
FindPkg returns the filename and unique package id for an import path based on
package information provided by build.Import (using the build.Default
build.Context). If no file was found, an empty filename is returned.

#### func  Import

```go
func Import(imports map[string]*types.Package, path string) (pkg *types.Package, err error)
```
Import imports a gc-generated package given its import path, adds the
corresponding package object to the imports map, and returns the object. Local
import paths are interpreted relative to the current working directory. The
imports map must contains all packages already imported.

#### func  ImportData

```go
func ImportData(imports map[string]*types.Package, filename, id string, data io.Reader) (pkg *types.Package, err error)
```
ImportData imports a package by reading the gc-generated export data, adds the
corresponding package object to the imports map indexed by id, and returns the
object.

The imports map must contains all packages already imported. The data reader
position must be the beginning of the export data section. The filename is only
used in error messages.

If imports[id] contains the completely imported package, that package can be
used directly, and there is no need to call this function (but there is also no
harm but for extra time used).
