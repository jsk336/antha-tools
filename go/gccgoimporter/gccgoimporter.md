# gccgoimporter
--
    import "."

Package gccgoimporter implements Import for gccgo-generated object files.

## Usage

#### func  GetImporter

```go
func GetImporter(searchpaths []string) types.Importer
```

#### type GccgoInstallation

```go
type GccgoInstallation struct {
	// Version of gcc (e.g. 4.8.0).
	GccVersion string

	// Target triple (e.g. x86_64-unknown-linux-gnu).
	TargetTriple string

	// Built-in library paths used by this installation.
	LibPaths []string
}
```

Information about a specific installation of gccgo.

#### func (*GccgoInstallation) GetImporter

```go
func (inst *GccgoInstallation) GetImporter(incpaths []string) types.Importer
```
Return an importer that searches incpaths followed by the gcc installation's
built-in search paths and the current directory.

#### func (*GccgoInstallation) InitFromDriver

```go
func (inst *GccgoInstallation) InitFromDriver(gccgoPath string) (err error)
```
Ask the driver at the given path for information for this GccgoInstallation.

#### func (*GccgoInstallation) SearchPaths

```go
func (inst *GccgoInstallation) SearchPaths() (paths []string)
```
Return the list of export search paths for this GccgoInstallation.
