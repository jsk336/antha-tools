# importer
--
    import "."

package importer implements an exporter and importer for Go export data.

## Usage

#### func  ExportData

```go
func ExportData(pkg *types.Package) []byte
```
ExportData serializes the interface (exported package objects) of package pkg
and returns the corresponding data. The export format is described elsewhere
(TODO).

#### func  ImportData

```go
func ImportData(imports map[string]*types.Package, data []byte) (*types.Package, error)
```
ImportData imports a package from the serialized package data. If data is
obviously malformed, an error is returned but in general it is not recommended
to call ImportData on untrusted data.
