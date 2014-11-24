# types
--
    import "."

Package types declares the data types and implements the algorithms for
type-checking of Go packages. Use Check and Config.Check to invoke the
type-checker.

Type-checking consists of several interdependent phases:

Name resolution maps each identifier (ast.Ident) in the program to the language
object (Object) it denotes. Use Info.{Defs,Uses,Implicits} for the results of
name resolution.

Constant folding computes the exact constant value (exact.Value) for every
expression (ast.Expr) that is a compile-time constant. Use
Info.Types[expr].Value for the results of constant folding.

Type inference computes the type (Type) of every expression (ast.Expr) and
checks for compliance with the language specification. Use Info.Types[expr].Type
for the results of type inference.

## Usage

```go
var (
	Universe *Scope
	Unsafe   *Package

	UniverseByte *Basic // uint8 alias, but has name "byte"
	UniverseRune *Basic // int32 alias, but has name "rune"
)
```

```go
var GcCompatibilityMode bool
```
If GcCompatibilityMode is set, printing of types is modified to match the
representation of some types in the gc compiler:

    - byte and rune lose their alias name and simply stand for
      uint8 and int32 respectively
    - embedded interfaces get flattened (the embedding info is lost,
      and certain recursive interface types cannot be printed anymore)

This makes it easier to compare packages computed with the type- checker vs
packages imported from gc export data.

Caution: This flag affects all uses of WriteType, globally. It is only provided
for testing in conjunction with gc-generated data. It may be removed at any
time.

```go
var Typ = [...]*Basic{
	Invalid: {Invalid, 0, "invalid type"},

	Bool:          {Bool, IsBoolean, "bool"},
	Int:           {Int, IsInteger, "int"},
	Int8:          {Int8, IsInteger, "int8"},
	Int16:         {Int16, IsInteger, "int16"},
	Int32:         {Int32, IsInteger, "int32"},
	Int64:         {Int64, IsInteger, "int64"},
	Uint:          {Uint, IsInteger | IsUnsigned, "uint"},
	Uint8:         {Uint8, IsInteger | IsUnsigned, "uint8"},
	Uint16:        {Uint16, IsInteger | IsUnsigned, "uint16"},
	Uint32:        {Uint32, IsInteger | IsUnsigned, "uint32"},
	Uint64:        {Uint64, IsInteger | IsUnsigned, "uint64"},
	Uintptr:       {Uintptr, IsInteger | IsUnsigned, "uintptr"},
	Float32:       {Float32, IsFloat, "float32"},
	Float64:       {Float64, IsFloat, "float64"},
	Complex64:     {Complex64, IsComplex, "complex64"},
	Complex128:    {Complex128, IsComplex, "complex128"},
	String:        {String, IsString, "string"},
	UnsafePointer: {UnsafePointer, 0, "Pointer"},

	UntypedBool:    {UntypedBool, IsBoolean | IsUntyped, "untyped bool"},
	UntypedInt:     {UntypedInt, IsInteger | IsUntyped, "untyped int"},
	UntypedRune:    {UntypedRune, IsInteger | IsUntyped, "untyped rune"},
	UntypedFloat:   {UntypedFloat, IsFloat | IsUntyped, "untyped float"},
	UntypedComplex: {UntypedComplex, IsComplex | IsUntyped, "untyped complex"},
	UntypedString:  {UntypedString, IsString | IsUntyped, "untyped string"},
	UntypedNil:     {UntypedNil, IsUntyped, "untyped nil"},
}
```

#### func  AssertableTo

```go
func AssertableTo(V *Interface, T Type) bool
```
AssertableTo reports whether a value of type V can be asserted to have type T.

#### func  AssignableTo

```go
func AssignableTo(V, T Type) bool
```
AssignableTo reports whether a value of type V is assignable to a variable of
type T.

#### func  Comparable

```go
func Comparable(T Type) bool
```
Comparable reports whether values of type T are comparable.

#### func  ConvertibleTo

```go
func ConvertibleTo(V, T Type) bool
```
ConvertibleTo reports whether a value of type V is convertible to a value of
type T.

#### func  DefPredeclaredTestFuncs

```go
func DefPredeclaredTestFuncs()
```
DefPredeclaredTestFuncs defines the assert and trace built-ins. These built-ins
are intended for debugging and testing of this package only.

#### func  ExprString

```go
func ExprString(x ast.Expr) string
```
ExprString returns the (possibly simplified) string representation for x.

#### func  Id

```go
func Id(pkg *Package, name string) string
```
Id returns name if it is exported, otherwise it returns the name qualified with
the package path.

#### func  Identical

```go
func Identical(x, y Type) bool
```
Identical reports whether x and y are identical.

#### func  Implements

```go
func Implements(V Type, T *Interface) bool
```
Implements reports whether type V implements interface T.

#### func  NewChecker

```go
func NewChecker(conf *Config, fset *token.FileSet, pkg *Package, info *Info) *checker
```
NewChecker returns a new Checker instance for a given package. Package files may
be incrementally added via checker.Files.

#### func  ObjectString

```go
func ObjectString(this *Package, obj Object) string
```
ObjectString returns the string form of obj. Object and type names are printed
package-qualified only if they do not belong to this package.

#### func  SelectionString

```go
func SelectionString(this *Package, s *Selection) string
```
SelectionString returns the string form of s. Type names are printed
package-qualified only if they do not belong to this package.

Examples:

    "field (T) f int"
    "method (T) f(X) Y"
    "method expr (T) f(X) Y"
    "qualified ident var math.Pi float64"

#### func  TypeString

```go
func TypeString(this *Package, typ Type) string
```
TypeString returns the string representation of typ. Named types are printed
package-qualified if they do not belong to this package.

#### func  WriteExpr

```go
func WriteExpr(buf *bytes.Buffer, x ast.Expr)
```
WriteExpr writes the (possibly simplified) string representation for x to buf.

#### func  WriteSignature

```go
func WriteSignature(buf *bytes.Buffer, this *Package, sig *Signature)
```
WriteSignature writes the representation of the signature sig to buf, without a
leading "func" keyword. Named types are printed package-qualified if they do not
belong to this package.

#### func  WriteType

```go
func WriteType(buf *bytes.Buffer, this *Package, typ Type)
```
WriteType writes the string representation of typ to buf. Named types are
printed package-qualified if they do not belong to this package.

#### type Array

```go
type Array struct {
}
```

An Array represents an array type.

#### func  NewArray

```go
func NewArray(elem Type, len int64) *Array
```
NewArray returns a new array type for the given element type and length.

#### func (*Array) Elem

```go
func (a *Array) Elem() Type
```
Elem returns element type of array a.

#### func (*Array) Len

```go
func (a *Array) Len() int64
```
Len returns the length of array a.

#### func (*Array) String

```go
func (t *Array) String() string
```

#### func (*Array) Underlying

```go
func (t *Array) Underlying() Type
```

#### type Basic

```go
type Basic struct {
}
```

A Basic represents a basic type.

#### func (*Basic) Info

```go
func (b *Basic) Info() BasicInfo
```
Info returns information about properties of basic type b.

#### func (*Basic) Kind

```go
func (b *Basic) Kind() BasicKind
```
Kind returns the kind of basic type b.

#### func (*Basic) Name

```go
func (b *Basic) Name() string
```
Name returns the name of basic type b.

#### func (*Basic) String

```go
func (t *Basic) String() string
```

#### func (*Basic) Underlying

```go
func (t *Basic) Underlying() Type
```

#### type BasicInfo

```go
type BasicInfo int
```

BasicInfo is a set of flags describing properties of a basic type.

```go
const (
	IsBoolean BasicInfo = 1 << iota
	IsInteger
	IsUnsigned
	IsFloat
	IsComplex
	IsString
	IsUntyped

	IsOrdered   = IsInteger | IsFloat | IsString
	IsNumeric   = IsInteger | IsFloat | IsComplex
	IsConstType = IsBoolean | IsNumeric | IsString
)
```
Properties of basic types.

#### type BasicKind

```go
type BasicKind int
```

BasicKind describes the kind of basic type.

```go
const (
	Invalid BasicKind = iota // type is invalid

	// predeclared types
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	String
	UnsafePointer

	// types for untyped values
	UntypedBool
	UntypedInt
	UntypedRune
	UntypedFloat
	UntypedComplex
	UntypedString
	UntypedNil

	// aliases
	Byte = Uint8
	Rune = Int32
)
```

#### type Builtin

```go
type Builtin struct {
}
```

A Builtin represents a built-in function. Builtins don't have a valid type.

#### func (*Builtin) Exported

```go
func (obj *Builtin) Exported() bool
```

#### func (*Builtin) Id

```go
func (obj *Builtin) Id() string
```

#### func (*Builtin) Name

```go
func (obj *Builtin) Name() string
```

#### func (*Builtin) Parent

```go
func (obj *Builtin) Parent() *Scope
```

#### func (*Builtin) Pkg

```go
func (obj *Builtin) Pkg() *Package
```

#### func (*Builtin) Pos

```go
func (obj *Builtin) Pos() token.Pos
```

#### func (*Builtin) String

```go
func (obj *Builtin) String() string
```

#### func (*Builtin) Type

```go
func (obj *Builtin) Type() Type
```

#### type Chan

```go
type Chan struct {
}
```

A Chan represents a channel type.

#### func  NewChan

```go
func NewChan(dir ChanDir, elem Type) *Chan
```
NewChan returns a new channel type for the given direction and element type.

#### func (*Chan) Dir

```go
func (c *Chan) Dir() ChanDir
```
Dir returns the direction of channel c.

#### func (*Chan) Elem

```go
func (c *Chan) Elem() Type
```
Elem returns the element type of channel c.

#### func (*Chan) String

```go
func (t *Chan) String() string
```

#### func (*Chan) Underlying

```go
func (t *Chan) Underlying() Type
```

#### type ChanDir

```go
type ChanDir int
```

A ChanDir value indicates a channel direction.

```go
const (
	SendRecv ChanDir = iota
	SendOnly
	RecvOnly
)
```
The direction of a channel is indicated by one of the following constants.

#### type Config

```go
type Config struct {
	// If IgnoreFuncBodies is set, function bodies are not
	// type-checked.
	IgnoreFuncBodies bool

	// If FakeImportC is set, `import "C"` (for packages requiring Cgo)
	// declares an empty "C" package and errors are omitted for qualified
	// identifiers referring to package C (which won't find an object).
	// This feature is intended for the standard library cmd/api tool.
	//
	// Caution: Effects may be unpredictable due to follow-up errors.
	//          Do not use casually!
	FakeImportC bool

	// Packages is used to look up (and thus canonicalize) packages by
	// package path. If Packages is nil, it is set to a new empty map.
	// During type-checking, imported packages are added to the map.
	Packages map[string]*Package

	// If Error != nil, it is called with each error found
	// during type checking; err has dynamic type Error.
	// Secondary errors (for instance, to enumerate all types
	// involved in an invalid recursive type declaration) have
	// error strings that start with a '\t' character.
	Error func(err error)

	// If Import != nil, it is called for each imported package.
	// Otherwise, DefaultImport is called.
	Import Importer

	// If Sizes != nil, it provides the sizing functions for package unsafe.
	// Otherwise &StdSize{WordSize: 8, MaxAlign: 8} is used instead.
	Sizes Sizes
}
```

A Config specifies the configuration for type checking. The zero value for
Config is a ready-to-use default configuration.

#### func (*Config) Check

```go
func (conf *Config) Check(path string, fset *token.FileSet, files []*ast.File, info *Info) (*Package, error)
```
Check type-checks a package and returns the resulting package object, the first
error if any, and if info != nil, additional type information. The package is
marked as complete if no errors occurred, otherwise it is incomplete.

The package is specified by a list of *ast.Files and corresponding file set, and
the package path the package is identified with. The clean path must not be
empty or dot (".").

#### type Const

```go
type Const struct {
}
```

A Const represents a declared constant.

#### func  NewConst

```go
func NewConst(pos token.Pos, pkg *Package, name string, typ Type, val exact.Value) *Const
```

#### func (*Const) Exported

```go
func (obj *Const) Exported() bool
```

#### func (*Const) Id

```go
func (obj *Const) Id() string
```

#### func (*Const) Name

```go
func (obj *Const) Name() string
```

#### func (*Const) Parent

```go
func (obj *Const) Parent() *Scope
```

#### func (*Const) Pkg

```go
func (obj *Const) Pkg() *Package
```

#### func (*Const) Pos

```go
func (obj *Const) Pos() token.Pos
```

#### func (*Const) String

```go
func (obj *Const) String() string
```

#### func (*Const) Type

```go
func (obj *Const) Type() Type
```

#### func (*Const) Val

```go
func (obj *Const) Val() exact.Value
```

#### type Error

```go
type Error struct {
	Fset *token.FileSet // file set for interpretation of Pos
	Pos  token.Pos      // error position
	Msg  string         // error message
	Soft bool           // if set, error is "soft"
}
```

An Error describes a type-checking error; it implements the error interface. A
"soft" error is an error that still permits a valid interpretation of a package
(such as "unused variable"); "hard" errors may lead to unpredictable behavior if
ignored.

#### func (Error) Error

```go
func (err Error) Error() string
```
Error returns an error string formatted as follows: filename:line:column:
message

#### type Func

```go
type Func struct {
}
```

A Func represents a declared function, concrete method, or abstract (interface)
method. Its Type() is always a *Signature. An abstract method may belong to many
interfaces due to embedding.

#### func  MissingMethod

```go
func MissingMethod(V Type, T *Interface, static bool) (method *Func, wrongType bool)
```
MissingMethod returns (nil, false) if V implements T, otherwise it returns a
missing method required by T and whether it is missing or just has the wrong
type.

For non-interface types V, or if static is set, V implements T if all methods of
T are present in V. Otherwise (V is an interface and static is not set),
MissingMethod only checks that methods of T which are also present in V have
matching types (e.g., for a type assertion x.(T) where x is of interface type
V).

#### func  NewFunc

```go
func NewFunc(pos token.Pos, pkg *Package, name string, sig *Signature) *Func
```

#### func (*Func) Exported

```go
func (obj *Func) Exported() bool
```

#### func (*Func) FullName

```go
func (obj *Func) FullName() string
```
FullName returns the package- or receiver-type-qualified name of function or
method obj.

#### func (*Func) Id

```go
func (obj *Func) Id() string
```

#### func (*Func) Name

```go
func (obj *Func) Name() string
```

#### func (*Func) Parent

```go
func (obj *Func) Parent() *Scope
```

#### func (*Func) Pkg

```go
func (obj *Func) Pkg() *Package
```

#### func (*Func) Pos

```go
func (obj *Func) Pos() token.Pos
```

#### func (*Func) Scope

```go
func (obj *Func) Scope() *Scope
```

#### func (*Func) String

```go
func (obj *Func) String() string
```

#### func (*Func) Type

```go
func (obj *Func) Type() Type
```

#### type Importer

```go
type Importer func(map[string]*Package, string) (*Package, error)
```

An importer resolves import paths to Packages. The imports map records packages
already known, indexed by package path. The type-checker will invoke Import with
Config.Packages. An importer must determine the canonical package path and check
imports to see if it is already present in the map. If so, the Importer can
return the map entry. Otherwise, the importer must load the package data for the
given path into a new *Package, record it in imports map, and return the
package. TODO(gri) Need to be clearer about requirements of completeness.

```go
var DefaultImport Importer
```
DefaultImport is the default importer invoked if Config.Import == nil. The
declaration:

    import _ "antha-tools/antha/gcimporter"

in a client of antha/types will initialize DefaultImport to gcimporter.Import.

#### type Info

```go
type Info struct {
	// Types maps expressions to their types, and for constant
	// expressions, their values.
	// Identifiers are collected in Defs and Uses, not Types.
	//
	// For an expression denoting a predeclared built-in function
	// the recorded signature is call-site specific. If the call
	// result is not a constant, the recorded type is an argument-
	// specific signature. Otherwise, the recorded type is invalid.
	Types map[ast.Expr]TypeAndValue

	// Defs maps identifiers to the objects they define (including
	// package names, dots "." of dot-imports, and blank "_" identifiers).
	// For identifiers that do not denote objects (e.g., the package name
	// in package clauses, or symbolic variables t in t := x.(type) of
	// type switch headers), the corresponding objects are nil.
	//
	// For an anonymous field, Defs returns the field *Var it defines.
	//
	// Invariant: Defs[id] == nil || Defs[id].Pos() == id.Pos()
	Defs map[*ast.Ident]Object

	// Uses maps identifiers to the objects they denote.
	//
	// For an anonymous field, Uses returns the *TypeName it denotes.
	//
	// Invariant: Uses[id].Pos() != id.Pos()
	Uses map[*ast.Ident]Object

	// Implicits maps nodes to their implicitly declared objects, if any.
	// The following node and object types may appear:
	//
	//	node               declared object
	//
	//	*ast.ImportSpec    *PkgName for dot-imports and imports without renames
	//	*ast.CaseClause    type-specific *Var for each type switch case clause (incl. default)
	//      *ast.Field         anonymous struct field or parameter *Var
	//
	Implicits map[ast.Node]Object

	// Selections maps selector expressions to their corresponding selections.
	Selections map[*ast.SelectorExpr]*Selection

	// Scopes maps ast.Nodes to the scopes they define. Package scopes are not
	// associated with a specific node but with all files belonging to a package.
	// Thus, the package scope can be found in the type-checked Package object.
	// Scopes nest, with the Universe scope being the outermost scope, enclosing
	// the package scope, which contains (one or more) files scopes, which enclose
	// function scopes which in turn enclose statement and function literal scopes.
	// Note that even though package-level functions are declared in the package
	// scope, the function scopes are embedded in the file scope of the file
	// containing the function declaration.
	//
	// The following node types may appear in Scopes:
	//
	//	*ast.File
	//	*ast.FuncType
	//	*ast.BlockStmt
	//	*ast.IfStmt
	//	*ast.SwitchStmt
	//	*ast.TypeSwitchStmt
	//	*ast.CaseClause
	//	*ast.CommClause
	//	*ast.ForStmt
	//	*ast.RangeStmt
	//
	Scopes map[ast.Node]*Scope

	// InitOrder is the list of package-level initializers in the order in which
	// they must be executed. Initializers referring to variables related by an
	// initialization dependency appear in topological order, the others appear
	// in source order. Variables without an initialization expression do not
	// appear in this list.
	InitOrder []*Initializer
}
```

Info holds result type information for a type-checked package. Only the
information for which a map is provided is collected. If the package has type
errors, the collected information may be incomplete.

#### type Initializer

```go
type Initializer struct {
	Lhs []*Var // var Lhs = Rhs
	Rhs ast.Expr
}
```

An Initializer describes a package-level variable, or a list of variables in
case of a multi-valued initialization expression, and the corresponding
initialization expression.

#### func (*Initializer) String

```go
func (init *Initializer) String() string
```

#### type Interface

```go
type Interface struct {
}
```

An Interface represents an interface type.

#### func  NewInterface

```go
func NewInterface(methods []*Func, embeddeds []*Named) *Interface
```
NewInterface returns a new interface for the given methods and embedded types.

#### func (*Interface) Embedded

```go
func (t *Interface) Embedded(i int) *Named
```
Embedded returns the i'th embedded type of interface t for 0 <= i <
t.NumEmbeddeds(). The types are ordered by the corresponding TypeName's unique
Id.

#### func (*Interface) Empty

```go
func (t *Interface) Empty() bool
```
Empty returns true if t is the empty interface.

#### func (*Interface) ExplicitMethod

```go
func (t *Interface) ExplicitMethod(i int) *Func
```
ExplicitMethod returns the i'th explicitly declared method of interface t for 0
<= i < t.NumExplicitMethods(). The methods are ordered by their unique Id.

#### func (*Interface) Method

```go
func (t *Interface) Method(i int) *Func
```
Method returns the i'th method of interface t for 0 <= i < t.NumMethods(). The
methods are ordered by their unique Id.

#### func (*Interface) NumEmbeddeds

```go
func (t *Interface) NumEmbeddeds() int
```
NumEmbeddeds returns the number of embedded types in interface t.

#### func (*Interface) NumExplicitMethods

```go
func (t *Interface) NumExplicitMethods() int
```
NumExplicitMethods returns the number of explicitly declared methods of
interface t.

#### func (*Interface) NumMethods

```go
func (t *Interface) NumMethods() int
```
NumMethods returns the total number of methods of interface t.

#### func (*Interface) String

```go
func (t *Interface) String() string
```

#### func (*Interface) Underlying

```go
func (t *Interface) Underlying() Type
```

#### type Label

```go
type Label struct {
}
```

A Label represents a declared label.

#### func  NewLabel

```go
func NewLabel(pos token.Pos, name string) *Label
```

#### func (*Label) Exported

```go
func (obj *Label) Exported() bool
```

#### func (*Label) Id

```go
func (obj *Label) Id() string
```

#### func (*Label) Name

```go
func (obj *Label) Name() string
```

#### func (*Label) Parent

```go
func (obj *Label) Parent() *Scope
```

#### func (*Label) Pkg

```go
func (obj *Label) Pkg() *Package
```

#### func (*Label) Pos

```go
func (obj *Label) Pos() token.Pos
```

#### func (*Label) String

```go
func (obj *Label) String() string
```

#### func (*Label) Type

```go
func (obj *Label) Type() Type
```

#### type Map

```go
type Map struct {
}
```

A Map represents a map type.

#### func  NewMap

```go
func NewMap(key, elem Type) *Map
```
NewMap returns a new map for the given key and element types.

#### func (*Map) Elem

```go
func (m *Map) Elem() Type
```
Elem returns the element type of map m.

#### func (*Map) Key

```go
func (m *Map) Key() Type
```
Key returns the key type of map m.

#### func (*Map) String

```go
func (t *Map) String() string
```

#### func (*Map) Underlying

```go
func (t *Map) Underlying() Type
```

#### type MethodSet

```go
type MethodSet struct {
}
```

A MethodSet is an ordered set of concrete or abstract (interface) methods; a
method is a MethodVal selection, and they are ordered by ascending m.Obj().Id().
The zero value for a MethodSet is a ready-to-use empty method set.

#### func  NewMethodSet

```go
func NewMethodSet(T Type) *MethodSet
```
NewMethodSet returns the method set for the given type T. It always returns a
non-nil method set, even if it is empty.

A MethodSetCache handles repeat queries more efficiently.

#### func (*MethodSet) At

```go
func (s *MethodSet) At(i int) *Selection
```
At returns the i'th method in s for 0 <= i < s.Len().

#### func (*MethodSet) Len

```go
func (s *MethodSet) Len() int
```
Len returns the number of methods in s.

#### func (*MethodSet) Lookup

```go
func (s *MethodSet) Lookup(pkg *Package, name string) *Selection
```
Lookup returns the method with matching package and name, or nil if not found.

#### func (*MethodSet) String

```go
func (s *MethodSet) String() string
```

#### type MethodSetCache

```go
type MethodSetCache struct {
}
```

A MethodSetCache records the method set of each type T for which MethodSet(T) is
called so that repeat queries are fast. The zero value is a ready-to-use cache
instance.

#### func (*MethodSetCache) MethodSet

```go
func (cache *MethodSetCache) MethodSet(T Type) *MethodSet
```
MethodSet returns the method set of type T. It is thread-safe.

If cache is nil, this function is equivalent to NewMethodSet(T). Utility
functions can thus expose an optional *MethodSetCache parameter to clients that
care about performance.

#### type Named

```go
type Named struct {
}
```

A Named represents a named type.

#### func  NewNamed

```go
func NewNamed(obj *TypeName, underlying Type, methods []*Func) *Named
```
NewNamed returns a new named type for the given type name, underlying type, and
associated methods. The underlying type must not be a *Named.

#### func (*Named) AddMethod

```go
func (t *Named) AddMethod(m *Func)
```
AddMethod adds method m unless it is already in the method list. TODO(gri) find
a better solution instead of providing this function

#### func (*Named) Method

```go
func (t *Named) Method(i int) *Func
```
Method returns the i'th method of named type t for 0 <= i < t.NumMethods().

#### func (*Named) NumMethods

```go
func (t *Named) NumMethods() int
```
NumMethods returns the number of explicit methods whose receiver is named type
t.

#### func (*Named) Obj

```go
func (t *Named) Obj() *TypeName
```
TypeName returns the type name for the named type t.

#### func (*Named) SetUnderlying

```go
func (t *Named) SetUnderlying(underlying Type)
```
SetUnderlying sets the underlying type and marks t as complete. TODO(gri)
determine if there's a better solution rather than providing this function

#### func (*Named) String

```go
func (t *Named) String() string
```

#### func (*Named) Underlying

```go
func (t *Named) Underlying() Type
```

#### type Nil

```go
type Nil struct {
}
```

Nil represents the predeclared value nil.

#### func (*Nil) Exported

```go
func (obj *Nil) Exported() bool
```

#### func (*Nil) Id

```go
func (obj *Nil) Id() string
```

#### func (*Nil) Name

```go
func (obj *Nil) Name() string
```

#### func (*Nil) Parent

```go
func (obj *Nil) Parent() *Scope
```

#### func (*Nil) Pkg

```go
func (obj *Nil) Pkg() *Package
```

#### func (*Nil) Pos

```go
func (obj *Nil) Pos() token.Pos
```

#### func (*Nil) String

```go
func (obj *Nil) String() string
```

#### func (*Nil) Type

```go
func (obj *Nil) Type() Type
```

#### type Object

```go
type Object interface {
	Parent() *Scope // scope in which this object is declared
	Pos() token.Pos // position of object identifier in declaration
	Pkg() *Package  // nil for objects in the Universe scope and labels
	Name() string   // package local object name
	Type() Type     // object type
	Exported() bool // reports whether the name starts with a capital letter
	Id() string     // object id (see Id below)

	// String returns a human-readable string of the object.
	String() string
	// contains filtered or unexported methods
}
```

An Object describes a named language entity such as a package, constant, type,
variable, function (incl. methods), or label. All objects implement the Object
interface.

#### func  LookupFieldOrMethod

```go
func LookupFieldOrMethod(T Type, pkg *Package, name string) (obj Object, index []int, indirect bool)
```
LookupFieldOrMethod looks up a field or method with given package and name in T
and returns the corresponding *Var or *Func, an index sequence, and a bool
indicating if there were any pointer indirections on the path to the field or
method.

The last index entry is the field or method index in the (possibly embedded)
type where the entry was found, either:

    1) the list of declared methods of a named type; or
    2) the list of all methods (method set) of an interface type; or
    3) the list of fields of a struct type.

The earlier index entries are the indices of the embedded fields traversed to
get to the found entry, starting at depth 0.

If no entry is found, a nil object is returned. In this case, the returned index
sequence points to an ambiguous entry if it exists, or it is nil.

#### type Package

```go
type Package struct {
}
```

A Package describes a Go package.

#### func  Check

```go
func Check(path string, fset *token.FileSet, files []*ast.File) (*Package, error)
```
Check type-checks a package and returns the resulting complete package object,
or a nil package and the first error. The package is specified by a list of
*ast.Files and corresponding file set, and the import path the package is
identified with. The clean path must not be empty or dot (".").

For more control over type-checking and results, use Config.Check.

#### func  NewPackage

```go
func NewPackage(path, name string) *Package
```
NewPackage returns a new Package for the given package path and name; the name
must not be the blank identifier. The package is not complete and contains no
explicit imports.

#### func (*Package) Complete

```go
func (pkg *Package) Complete() bool
```
A package is complete if its scope contains (at least) all exported objects;
otherwise it is incomplete.

#### func (*Package) Imports

```go
func (pkg *Package) Imports() []*Package
```
Imports returns the list of packages explicitly imported by pkg; the list is in
source order. Package unsafe is excluded.

#### func (*Package) MarkComplete

```go
func (pkg *Package) MarkComplete()
```
MarkComplete marks a package as complete.

#### func (*Package) Name

```go
func (pkg *Package) Name() string
```
Name returns the package name.

#### func (*Package) Path

```go
func (pkg *Package) Path() string
```
Path returns the package path.

#### func (*Package) Scope

```go
func (pkg *Package) Scope() *Scope
```
Scope returns the (complete or incomplete) package scope holding the objects
declared at package level (TypeNames, Consts, Vars, and Funcs).

#### func (*Package) SetImports

```go
func (pkg *Package) SetImports(list []*Package)
```
SetImports sets the list of explicitly imported packages to list. It is the
caller's responsibility to make sure list elements are unique.

#### func (*Package) String

```go
func (pkg *Package) String() string
```

#### type PkgName

```go
type PkgName struct {
}
```

A PkgName represents an imported Go package.

#### func  NewPkgName

```go
func NewPkgName(pos token.Pos, pkg *Package, name string) *PkgName
```

#### func (*PkgName) Exported

```go
func (obj *PkgName) Exported() bool
```

#### func (*PkgName) Id

```go
func (obj *PkgName) Id() string
```

#### func (*PkgName) Name

```go
func (obj *PkgName) Name() string
```

#### func (*PkgName) Parent

```go
func (obj *PkgName) Parent() *Scope
```

#### func (*PkgName) Pkg

```go
func (obj *PkgName) Pkg() *Package
```

#### func (*PkgName) Pos

```go
func (obj *PkgName) Pos() token.Pos
```

#### func (*PkgName) String

```go
func (obj *PkgName) String() string
```

#### func (*PkgName) Type

```go
func (obj *PkgName) Type() Type
```

#### type Pointer

```go
type Pointer struct {
}
```

A Pointer represents a pointer type.

#### func  NewPointer

```go
func NewPointer(elem Type) *Pointer
```
NewPointer returns a new pointer type for the given element (base) type.

#### func (*Pointer) Elem

```go
func (p *Pointer) Elem() Type
```
Elem returns the element type for the given pointer p.

#### func (*Pointer) String

```go
func (t *Pointer) String() string
```

#### func (*Pointer) Underlying

```go
func (t *Pointer) Underlying() Type
```

#### type Scope

```go
type Scope struct {
}
```

A Scope maintains a set of objects and links to its containing (parent) and
contained (children) scopes. Objects may be inserted and looked up by name. The
zero value for Scope is a ready-to-use empty scope.

#### func  NewScope

```go
func NewScope(parent *Scope, comment string) *Scope
```
NewScope returns a new, empty scope contained in the given parent scope, if any.
The comment is for debugging only.

#### func (*Scope) Child

```go
func (s *Scope) Child(i int) *Scope
```
Child returns the i'th child scope for 0 <= i < NumChildren().

#### func (*Scope) Insert

```go
func (s *Scope) Insert(obj Object) Object
```
Insert attempts to insert an object obj into scope s. If s already contains an
alternative object alt with the same name, Insert leaves s unchanged and returns
alt. Otherwise it inserts obj, sets the object's parent scope if not already
set, and returns nil.

#### func (*Scope) Len

```go
func (s *Scope) Len() int
```
Len() returns the number of scope elements.

#### func (*Scope) Lookup

```go
func (s *Scope) Lookup(name string) Object
```
Lookup returns the object in scope s with the given name if such an object
exists; otherwise the result is nil.

#### func (*Scope) LookupParent

```go
func (s *Scope) LookupParent(name string) Object
```
LookupParent follows the parent chain of scopes starting with s until it finds a
scope where Lookup(name) returns a non-nil object, and then returns that object.
If no such scope exists, the result is nil.

#### func (*Scope) Names

```go
func (s *Scope) Names() []string
```
Names returns the scope's element names in sorted order.

#### func (*Scope) NumChildren

```go
func (s *Scope) NumChildren() int
```
NumChildren() returns the number of scopes nested in s.

#### func (*Scope) Parent

```go
func (s *Scope) Parent() *Scope
```
Parent returns the scope's containing (parent) scope.

#### func (*Scope) String

```go
func (s *Scope) String() string
```
String returns a string representation of the scope, for debugging.

#### func (*Scope) WriteTo

```go
func (s *Scope) WriteTo(w io.Writer, n int, recurse bool)
```
WriteTo writes a string representation of the scope to w, with the scope
elements sorted by name. The level of indentation is controlled by n >= 0, with
n == 0 for no indentation. If recurse is set, it also writes nested (children)
scopes.

#### type Selection

```go
type Selection struct {
}
```

A Selection describes a selector expression x.f. For the declarations:

    type T struct{ x int; E }
    type E struct{}
    func (e E) m() {}
    var p *T

the following relations exist:

    Selector    Kind          Recv    Obj    Type               Index     Indirect

    p.x         FieldVal      T       x      int                {0}       true
    p.m         MethodVal     *T      m      func (e *T) m()    {1, 0}    true
    T.m         MethodExpr    T       m      func m(_ T)        {1, 0}    false
    math.Pi     PackageObj    nil     Pi     untyped numeric    nil       false

#### func (*Selection) Index

```go
func (s *Selection) Index() []int
```
Index describes the path from x to f in x.f. The result is nil if x.f is a
qualified identifier (PackageObj).

The last index entry is the field or method index of the type declaring f;
either:

    1) the list of declared methods of a named type; or
    2) the list of methods of an interface type; or
    3) the list of fields of a struct type.

The earlier index entries are the indices of the embedded fields implicitly
traversed to get from (the type of) x to f, starting at embedding depth 0.

#### func (*Selection) Indirect

```go
func (s *Selection) Indirect() bool
```
Indirect reports whether any pointer indirection was required to get from x to f
in x.f. The result is false if x.f is a qualified identifier (PackageObj).

#### func (*Selection) Kind

```go
func (s *Selection) Kind() SelectionKind
```
Kind returns the selection kind.

#### func (*Selection) Obj

```go
func (s *Selection) Obj() Object
```
Obj returns the object denoted by x.f. The following object types may appear:

    Kind          Object

    FieldVal      *Var                          field
    MethodVal     *Func                         method
    MethodExpr    *Func                         method
    PackageObj    *Const, *Type, *Var, *Func    imported const, type, var, or func

#### func (*Selection) Recv

```go
func (s *Selection) Recv() Type
```
Recv returns the type of x in x.f. The result is nil if x.f is a qualified
identifier (PackageObj).

#### func (*Selection) String

```go
func (s *Selection) String() string
```

#### func (*Selection) Type

```go
func (s *Selection) Type() Type
```
Type returns the type of x.f, which may be different from the type of f. See
Selection for more information.

#### type SelectionKind

```go
type SelectionKind int
```

SelectionKind describes the kind of a selector expression x.f.

```go
const (
	FieldVal   SelectionKind = iota // x.f is a struct field selector
	MethodVal                       // x.f is a method selector
	MethodExpr                      // x.f is a method expression
	PackageObj                      // x.f is a qualified identifier
)
```

#### type Signature

```go
type Signature struct {
}
```

A Signature represents a (non-builtin) function or method type.

#### func  NewSignature

```go
func NewSignature(scope *Scope, recv *Var, params, results *Tuple, variadic bool) *Signature
```
NewSignature returns a new function type for the given receiver, parameters, and
results, either of which may be nil. If variadic is set, the function is
variadic, it must have at least one parameter, and the last parameter must be of
unnamed slice type.

#### func (*Signature) Params

```go
func (s *Signature) Params() *Tuple
```
Params returns the parameters of signature s, or nil.

#### func (*Signature) Recv

```go
func (s *Signature) Recv() *Var
```
Recv returns the receiver of signature s (if a method), or nil if a function.

For an abstract method, Recv returns the enclosing interface either as a *Named
or an *Interface. Due to embedding, an interface may contain methods whose
receiver type is a different interface.

#### func (*Signature) Results

```go
func (s *Signature) Results() *Tuple
```
Results returns the results of signature s, or nil.

#### func (*Signature) String

```go
func (t *Signature) String() string
```

#### func (*Signature) Underlying

```go
func (t *Signature) Underlying() Type
```

#### func (*Signature) Variadic

```go
func (s *Signature) Variadic() bool
```
Variadic reports whether the signature s is variadic.

#### type Sizes

```go
type Sizes interface {
	// Alignof returns the alignment of a variable of type T.
	// Alignof must implement the alignment guarantees required by the spec.
	Alignof(T Type) int64

	// Offsetsof returns the offsets of the given struct fields, in bytes.
	// Offsetsof must implement the offset guarantees required by the spec.
	Offsetsof(fields []*Var) []int64

	// Sizeof returns the size of a variable of type T.
	// Sizeof must implement the size guarantees required by the spec.
	Sizeof(T Type) int64
}
```

Sizes defines the sizing functions for package unsafe.

#### type Slice

```go
type Slice struct {
}
```

A Slice represents a slice type.

#### func  NewSlice

```go
func NewSlice(elem Type) *Slice
```
NewSlice returns a new slice type for the given element type.

#### func (*Slice) Elem

```go
func (s *Slice) Elem() Type
```
Elem returns the element type of slice s.

#### func (*Slice) String

```go
func (t *Slice) String() string
```

#### func (*Slice) Underlying

```go
func (t *Slice) Underlying() Type
```

#### type StdSizes

```go
type StdSizes struct {
	WordSize int64 // word size in bytes - must be >= 4 (32bits)
	MaxAlign int64 // maximum alignment in bytes - must be >= 1
}
```

StdSizes is a convenience type for creating commonly used Sizes. It makes the
following simplifying assumptions:

    - The size of explicitly sized basic types (int16, etc.) is the
      specified size.
    - The size of strings and interfaces is 2*WordSize.
    - The size of slices is 3*WordSize.
    - All other types have size WordSize.
    - Arrays and structs are aligned per spec definition; all other
      types are naturally aligned with a maximum alignment MaxAlign.

*StdSizes implements Sizes.

#### func (*StdSizes) Alignof

```go
func (s *StdSizes) Alignof(T Type) int64
```

#### func (*StdSizes) Offsetsof

```go
func (s *StdSizes) Offsetsof(fields []*Var) []int64
```

#### func (*StdSizes) Sizeof

```go
func (s *StdSizes) Sizeof(T Type) int64
```

#### type Struct

```go
type Struct struct {
}
```

A Struct represents a struct type.

#### func  NewStruct

```go
func NewStruct(fields []*Var, tags []string) *Struct
```
NewStruct returns a new struct with the given fields and corresponding field
tags. If a field with index i has a tag, tags[i] must be that tag, but len(tags)
may be only as long as required to hold the tag with the largest index i.
Consequently, if no field has a tag, tags may be nil.

#### func (*Struct) Field

```go
func (s *Struct) Field(i int) *Var
```
Field returns the i'th field for 0 <= i < NumFields().

#### func (*Struct) NumFields

```go
func (s *Struct) NumFields() int
```
NumFields returns the number of fields in the struct (including blank and
anonymous fields).

#### func (*Struct) String

```go
func (t *Struct) String() string
```

#### func (*Struct) Tag

```go
func (s *Struct) Tag(i int) string
```
Tag returns the i'th field tag for 0 <= i < NumFields().

#### func (*Struct) Underlying

```go
func (t *Struct) Underlying() Type
```

#### type Tuple

```go
type Tuple struct {
}
```

A Tuple represents an ordered list of variables; a nil *Tuple is a valid (empty)
tuple. Tuples are used as components of signatures and to represent the type of
multiple assignments; they are not first class types of Go.

#### func  NewTuple

```go
func NewTuple(x ...*Var) *Tuple
```
NewTuple returns a new tuple for the given variables.

#### func (*Tuple) At

```go
func (t *Tuple) At(i int) *Var
```
At returns the i'th variable of tuple t.

#### func (*Tuple) Len

```go
func (t *Tuple) Len() int
```
Len returns the number variables of tuple t.

#### func (*Tuple) String

```go
func (t *Tuple) String() string
```

#### func (*Tuple) Underlying

```go
func (t *Tuple) Underlying() Type
```

#### type Type

```go
type Type interface {
	// Underlying returns the underlying type of a type.
	Underlying() Type

	// String returns a string representation of a type.
	String() string
}
```

A Type represents a type of Go. All types implement the Type interface.

#### func  Eval

```go
func Eval(str string, pkg *Package, scope *Scope) (typ Type, val exact.Value, err error)
```
Eval returns the type and, if constant, the value for the expression or type
literal string str evaluated in scope. If the expression contains function
literals, the function bodies are ignored (though they must be syntactically
correct).

If pkg == nil, the Universe scope is used and the provided scope is ignored.
Otherwise, the scope must belong to the package (either the package scope, or
nested within the package scope).

An error is returned if the scope is incorrect, the string has syntax errors, or
if it cannot be evaluated in the scope. Position info for objects in the result
type is undefined.

Note: Eval should not be used instead of running Check to compute types and
values, but in addition to Check. Eval will re-evaluate its argument each time,
and it also does not know about the context in which an expression is used
(e.g., an assignment). Thus, top- level untyped constants will return an untyped
type rather then the respective context-specific type.

#### func  EvalNode

```go
func EvalNode(fset *token.FileSet, node ast.Expr, pkg *Package, scope *Scope) (typ Type, val exact.Value, err error)
```
EvalNode is like Eval but instead of string it accepts an expression node and
respective file set.

An error is returned if the scope is incorrect if the node cannot be evaluated
in the scope.

#### func  New

```go
func New(str string) Type
```
New is a convenience function to create a new type from a given expression or
type literal string evaluated in Universe scope. New(str) is shorthand for
Eval(str, nil, nil), but only returns the type result, and panics in case of an
error. Position info for objects in the result type is undefined.

#### type TypeAndValue

```go
type TypeAndValue struct {
	Type  Type
	Value exact.Value
}
```


#### type TypeName

```go
type TypeName struct {
}
```

A TypeName represents a declared type.

#### func  NewTypeName

```go
func NewTypeName(pos token.Pos, pkg *Package, name string, typ Type) *TypeName
```

#### func (*TypeName) Exported

```go
func (obj *TypeName) Exported() bool
```

#### func (*TypeName) Id

```go
func (obj *TypeName) Id() string
```

#### func (*TypeName) Name

```go
func (obj *TypeName) Name() string
```

#### func (*TypeName) Parent

```go
func (obj *TypeName) Parent() *Scope
```

#### func (*TypeName) Pkg

```go
func (obj *TypeName) Pkg() *Package
```

#### func (*TypeName) Pos

```go
func (obj *TypeName) Pos() token.Pos
```

#### func (*TypeName) String

```go
func (obj *TypeName) String() string
```

#### func (*TypeName) Type

```go
func (obj *TypeName) Type() Type
```

#### type Var

```go
type Var struct {
}
```

A Variable represents a declared variable (including function parameters and
results, and struct fields).

#### func  NewField

```go
func NewField(pos token.Pos, pkg *Package, name string, typ Type, anonymous bool) *Var
```

#### func  NewParam

```go
func NewParam(pos token.Pos, pkg *Package, name string, typ Type) *Var
```

#### func  NewVar

```go
func NewVar(pos token.Pos, pkg *Package, name string, typ Type) *Var
```

#### func (*Var) Anonymous

```go
func (obj *Var) Anonymous() bool
```

#### func (*Var) Exported

```go
func (obj *Var) Exported() bool
```

#### func (*Var) Id

```go
func (obj *Var) Id() string
```

#### func (*Var) IsField

```go
func (obj *Var) IsField() bool
```

#### func (*Var) Name

```go
func (obj *Var) Name() string
```

#### func (*Var) Parent

```go
func (obj *Var) Parent() *Scope
```

#### func (*Var) Pkg

```go
func (obj *Var) Pkg() *Package
```

#### func (*Var) Pos

```go
func (obj *Var) Pos() token.Pos
```

#### func (*Var) String

```go
func (obj *Var) String() string
```

#### func (*Var) Type

```go
func (obj *Var) Type() Type
```
