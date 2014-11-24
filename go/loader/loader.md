# loader
--
    import "."

Package loader loads, parses and type-checks packages of Go code plus their
transitive closure, and retains both the ASTs and the derived facts.

THIS INTERFACE IS EXPERIMENTAL AND IS LIKELY TO CHANGE.

The package defines two primary types: Config, which specifies a set of initial
packages to load and various other options; and Program, which is the result of
successfully loading the packages specified by a configuration.

The configuration can be set directly, but *Config provides various convenience
methods to simplify the common cases, each of which can be called any number of
times. Finally, these are followed by a call to Load() to actually load and
type-check the program.

    var conf loader.Config

    // Use the command-line arguments to specify
    // a set of initial packages to load from source.
    // See FromArgsUsage for help.
    rest, err := conf.FromArgs(os.Args[1:], wantTests)

    // Parse the specified files and create an ad-hoc package with path "foo".
    // All files must have the same 'package' declaration.
    err := conf.CreateFromFilenames("foo", "foo.go", "bar.go")

    // Create an ad-hoc package with path "foo" from
    // the specified already-parsed files.
    // All ASTs must have the same 'package' declaration.
    err := conf.CreateFromFiles("foo", parsedFiles)

    // Add "runtime" to the set of packages to be loaded.
    conf.Import("runtime")

    // Adds "fmt" and "fmt_test" to the set of packages
    // to be loaded.  "fmt" will include *_test.go files.
    err := conf.ImportWithTests("fmt")

    // Finally, load all the packages specified by the configuration.
    prog, err := conf.Load()


### CONCEPTS AND TERMINOLOGY

An AD-HOC package is one specified as a set of source files on the command line.
In the simplest case, it may consist of a single file such as
src/pkg/net/http/triv.go.

EXTERNAL TEST packages are those comprised of a set of *_test.go files all with
the same 'package foo_test' declaration, all in the same directory.
(antha/build.Package calls these files XTestFiles.)

An IMPORTABLE package is one that can be referred to by some import spec. The
Path() of each importable package is unique within a Program.

Ad-hoc packages and external test packages are NON-IMPORTABLE. The Path() of an
ad-hoc package is inferred from the package declarations of its files and is
therefore not a unique package key. For example, Config.CreatePkgs may specify
two initial ad-hoc packages both called "main".

An AUGMENTED package is an importable package P plus all the *_test.go files
with same 'package foo' declaration as P. (antha/build.Package calls these files
TestFiles.)

The INITIAL packages are those specified in the configuration. A DEPENDENCY is a
package loaded to satisfy an import in an initial package or another dependency.

## Usage

```go
const FromArgsUsage = `
<args> is a list of arguments denoting a set of initial packages.
It may take one of two forms:

1. A list of *.go source files.

   All of the specified files are loaded, parsed and type-checked
   as a single package.  All the files must belong to the same directory.

2. A list of import paths, each denoting a package.

   The package's directory is found relative to the $GOROOT and
   $GOPATH using similar logic to 'go build', and the *.go files in
   that directory are loaded, parsed and type-checked as a single
   package.

   In addition, all *_test.go files in the directory are then loaded
   and parsed.  Those files whose package declaration equals that of
   the non-*_test.go files are included in the primary package.  Test
   files whose package declaration ends with "_test" are type-checked
   as another package, the 'external' test package, so that a single
   import path may denote two packages.  (Whether this behaviour is
   enabled is tool-specific, and may depend on additional flags.)

   Due to current limitations in the type-checker, only the first
   import path of the command line will contribute any tests.

A '--' argument terminates the list of packages.
`
```
FromArgsUsage is a partial usage message that applications calling FromArgs may
wish to include in their -help output.

#### type Config

```go
type Config struct {
	// Fset is the file set for the parser to use when loading the
	// program.  If nil, it will be lazily initialized by any
	// method of Config.
	Fset *token.FileSet

	// ParserMode specifies the mode to be used by the parser when
	// loading source packages.
	ParserMode parser.Mode

	// TypeChecker contains options relating to the type checker.
	//
	// The supplied IgnoreFuncBodies is not used; the effective
	// value comes from the TypeCheckFuncBodies func below.
	//
	// TypeChecker.Packages is lazily initialized during Load.
	TypeChecker types.Config

	// TypeCheckFuncBodies is a predicate over package import
	// paths.  A package for which the predicate is false will
	// have its package-level declarations type checked, but not
	// its function bodies; this can be used to quickly load
	// dependencies from source.  If nil, all func bodies are type
	// checked.
	TypeCheckFuncBodies func(string) bool

	// SourceImports determines whether to satisfy dependencies by
	// loading Go source code.
	//
	// If true, the entire program---the initial packages and
	// their transitive closure of dependencies---will be loaded,
	// parsed and type-checked.  This is required for
	// whole-program analyses such as pointer analysis.
	//
	// If false, the TypeChecker.Import mechanism will be used
	// instead.  Since that typically supplies only the types of
	// package-level declarations and values of constants, but no
	// code, it will not yield a whole program.  It is intended
	// for analyses that perform intraprocedural analysis of a
	// single package, e.g. traditional compilation.
	//
	// The initial packages (CreatePkgs and ImportPkgs) are always
	// loaded from Go source, regardless of this flag's setting.
	SourceImports bool

	// If Build is non-nil, it is used to locate source packages.
	// Otherwise &build.Default is used.
	Build *build.Context

	// If DisplayPath is non-nil, it is used to transform each
	// file name obtained from Build.Import().  This can be used
	// to prevent a virtualized build.Config's file names from
	// leaking into the user interface.
	DisplayPath func(path string) string

	// If AllowTypeErrors is true, Load will return a Program even
	// if some of the its packages contained type errors; such
	// errors are accessible via PackageInfo.TypeError.
	// If false, Load will fail if any package had a type error.
	AllowTypeErrors bool

	// CreatePkgs specifies a list of non-importable initial
	// packages to create.  Each element specifies a list of
	// parsed files to be type-checked into a new package, and a
	// path for that package.  If the path is "", the package's
	// name will be used instead.  The path needn't be globally
	// unique.
	//
	// The resulting packages will appear in the corresponding
	// elements of the Program.Created slice.
	CreatePkgs []CreatePkg

	// ImportPkgs specifies a set of initial packages to load from
	// source.  The map keys are package import paths, used to
	// locate the package relative to $GOROOT.  The corresponding
	// values indicate whether to augment the package by *_test.go
	// files in a second pass.
	ImportPkgs map[string]bool
}
```

Config specifies the configuration for a program to load. The zero value for
Config is a ready-to-use default configuration.

#### func (*Config) CreateFromFilenames

```go
func (conf *Config) CreateFromFilenames(path string, filenames ...string) error
```
CreateFromFilenames is a convenience function that parses the specified *.go
files and adds a package entry for them to conf.CreatePkgs.

#### func (*Config) CreateFromFiles

```go
func (conf *Config) CreateFromFiles(path string, files ...*ast.File)
```
CreateFromFiles is a convenience function that adds a CreatePkgs entry to create
package of the specified path and parsed files.

Precondition: conf.Fset is non-nil and was the fileset used to parse the files.
(e.g. the files came from conf.ParseFile().)

#### func (*Config) FromArgs

```go
func (conf *Config) FromArgs(args []string, xtest bool) (rest []string, err error)
```
FromArgs interprets args as a set of initial packages to load from source and
updates the configuration. It returns the list of unconsumed arguments.

It is intended for use in command-line interfaces that require a set of initial
packages to be specified; see FromArgsUsage message for details.

#### func (*Config) Import

```go
func (conf *Config) Import(path string)
```
Import is a convenience function that adds path to ImportPkgs, the set of
initial packages that will be imported from source.

#### func (*Config) ImportWithTests

```go
func (conf *Config) ImportWithTests(path string) error
```
ImportWithTests is a convenience function that adds path to ImportPkgs, the set
of initial source packages located relative to $GOPATH. The package will be
augmented by any *_test.go files in its directory that contain a "package x"
(not "package x_test") declaration.

In addition, if any *_test.go files contain a "package x_test" declaration, an
additional package comprising just those files will be added to CreatePkgs.

#### func (*Config) Load

```go
func (conf *Config) Load() (*Program, error)
```
Load creates the initial packages specified by conf.{Create,Import}Pkgs, loading
their dependencies packages as needed.

On success, Load returns a Program containing a PackageInfo for each package. On
failure, it returns an error.

If conf.AllowTypeErrors is set, a type error does not cause Load to fail, but is
recorded in the PackageInfo.TypeError field.

It is an error if no packages were loaded.

#### func (*Config) ParseFile

```go
func (conf *Config) ParseFile(filename string, src interface{}) (*ast.File, error)
```
ParseFile is a convenience function that invokes the parser using the Config's
FileSet, which is initialized if nil.

#### type CreatePkg

```go
type CreatePkg struct {
	Path  string
	Files []*ast.File
}
```


#### type PackageInfo

```go
type PackageInfo struct {
	Pkg                   *types.Package
	Importable            bool        // true if 'import "Pkg.Path()"' would resolve to this
	TransitivelyErrorFree bool        // true if Pkg and all its dependencies are free of errors
	Files                 []*ast.File // abstract syntax for the package's files
	TypeError             error       // non-nil if the package had type errors
	types.Info                        // type-checker deductions.
}
```

PackageInfo holds the ASTs and facts derived by the type-checker for a single
package.

Not mutated once exposed via the API.

#### func (*PackageInfo) ImportSpecPkg

```go
func (info *PackageInfo) ImportSpecPkg(spec *ast.ImportSpec) *types.PkgName
```
ImportSpecPkg returns the PkgName for a given ImportSpec, possibly an implicit
one for a dot-import or an import-without-rename. It returns nil if not found.

#### func (*PackageInfo) IsType

```go
func (info *PackageInfo) IsType(e ast.Expr) bool
```
IsType returns true iff expression e denotes a type. Precondition: e belongs to
the package's ASTs.

TODO(gri): move this into antha/types.

#### func (*PackageInfo) ObjectOf

```go
func (info *PackageInfo) ObjectOf(id *ast.Ident) types.Object
```
ObjectOf returns the typechecker object denoted by the specified id.

If id is an anonymous struct field, the field (*types.Var) is returned, not the
type (*types.TypeName).

Precondition: id belongs to the package's ASTs.

#### func (*PackageInfo) String

```go
func (info *PackageInfo) String() string
```

#### func (*PackageInfo) TypeCaseVar

```go
func (info *PackageInfo) TypeCaseVar(cc *ast.CaseClause) *types.Var
```
TypeCaseVar returns the implicit variable created by a single-type case clause
in a type switch, or nil if not found.

#### func (*PackageInfo) TypeOf

```go
func (info *PackageInfo) TypeOf(e ast.Expr) types.Type
```
TypeOf returns the type of expression e. Precondition: e belongs to the
package's ASTs.

#### func (*PackageInfo) ValueOf

```go
func (info *PackageInfo) ValueOf(e ast.Expr) exact.Value
```
ValueOf returns the value of expression e if it is a constant, nil otherwise.
Precondition: e belongs to the package's ASTs.

#### type Program

```go
type Program struct {
	Fset *token.FileSet // the file set for this program

	// Created[i] contains the initial package whose ASTs were
	// supplied by Config.CreatePkgs[i].
	Created []*PackageInfo

	// Imported contains the initially imported packages,
	// as specified by Config.ImportPkgs.
	Imported map[string]*PackageInfo

	// ImportMap is the canonical mapping of import paths to
	// packages used by the type-checker (Config.TypeChecker.Packages).
	// It contains all Imported initial packages, but not Created
	// ones, and all imported dependencies.
	ImportMap map[string]*types.Package

	// AllPackages contains the PackageInfo of every package
	// encountered by Load: all initial packages and all
	// dependencies, including incomplete ones.
	AllPackages map[*types.Package]*PackageInfo
}
```

A Program is a Go program loaded from source or binary as specified by a Config.

#### func (*Program) InitialPackages

```go
func (prog *Program) InitialPackages() []*PackageInfo
```
InitialPackages returns a new slice containing the set of initial packages
(Created + Imported) in unspecified order.

#### func (*Program) PathEnclosingInterval

```go
func (prog *Program) PathEnclosingInterval(start, end token.Pos) (pkg *PackageInfo, path []ast.Node, exact bool)
```
PathEnclosingInterval returns the PackageInfo and ast.Node that contain source
interval [start, end), and all the node's ancestors up to the AST root. It
searches all ast.Files of all packages in prog. exact is defined as for
astutil.PathEnclosingInterval.

The zero value is returned if not found.
