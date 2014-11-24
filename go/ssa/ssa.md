# ssa
--
    import "."

Package ssa defines a representation of the elements of Go programs (packages,
types, functions, variables and constants) using a static single-assignment
(SSA) form intermediate representation (IR) for the bodies of functions.

THIS INTERFACE IS EXPERIMENTAL AND IS LIKELY TO CHANGE.

For an introduction to SSA form, see
http://en.wikipedia.org/wiki/Static_single_assignment_form. This page provides a
broader reading list: http://www.dcs.gla.ac.uk/~jsinger/ssa.html.

The level of abstraction of the SSA form is intentionally close to the source
language to facilitate construction of source analysis tools. It is not intended
for machine code generation.

All looping, branching and switching constructs are replaced with unstructured
control flow. Higher-level control flow constructs such as multi-way branch can
be reconstructed as needed; see ssautil.Switches() for an example.

To construct an SSA-form program, call ssa.Create on a loader.Program, a set of
type-checked packages created from parsed Go source files. The resulting
ssa.Program contains all the packages and their members, but SSA code is not
created for function bodies until a subsequent call to (*Package).Build.

The builder initially builds a naive SSA form in which all local variables are
addresses of stack locations with explicit loads and stores. Registerisation of
eligible locals and φ-node insertion using dominance and dataflow are then
performed as a second pass called "lifting" to improve the accuracy and
performance of subsequent analyses; this pass can be skipped by setting the
NaiveForm builder flag.

The primary interfaces of this package are:

    - Member: a named member of a Go package.
    - Value: an expression that yields a value.
    - Instruction: a statement that consumes values and performs computation.

A computation that yields a result implements both the Value and Instruction
interfaces. The following table shows for each concrete type which of these
interfaces it implements.

                       Value?          Instruction?    Member?
    *Alloc             ✔               ✔
    *BinOp             ✔               ✔
    *Builtin           ✔
    *Call              ✔               ✔
    *Capture           ✔
    *ChangeInterface   ✔               ✔
    *ChangeType        ✔               ✔
    *Const             ✔
    *Convert           ✔               ✔
    *DebugRef                          ✔
    *Defer                             ✔
    *Extract           ✔               ✔
    *Field             ✔               ✔
    *FieldAddr         ✔               ✔
    *Function          ✔                               ✔ (func)
    *Global            ✔                               ✔ (var)
    *Go                                ✔
    *If                                ✔
    *Index             ✔               ✔
    *IndexAddr         ✔               ✔
    *Jump                              ✔
    *Lookup            ✔               ✔
    *MakeChan          ✔               ✔
    *MakeClosure       ✔               ✔
    *MakeInterface     ✔               ✔
    *MakeMap           ✔               ✔
    *MakeSlice         ✔               ✔
    *MapUpdate                         ✔
    *NamedConst                                        ✔ (const)
    *Next              ✔               ✔
    *Panic                             ✔
    *Parameter         ✔
    *Phi               ✔               ✔
    *Range             ✔               ✔
    *Return                            ✔
    *RunDefers                         ✔
    *Select            ✔               ✔
    *Send                              ✔
    *Slice             ✔               ✔
    *Store                             ✔
    *Type                                              ✔ (type)
    *TypeAssert        ✔               ✔
    *UnOp              ✔               ✔

Other key types in this package include: Program, Package, Function and
BasicBlock.

The program representation constructed by this package is fully resolved
internally, i.e. it does not rely on the names of Values, Packages, Functions,
Types or BasicBlocks for the correct interpretation of the program. Only the
identities of objects and the topology of the SSA and type graphs are
semantically significant. (There is one exception: Ids, used to identify field
and method names, contain strings.) Avoidance of name-based operations
simplifies the implementation of subsequent passes and can make them very
efficient. Many objects are nonetheless named to aid in debugging, but it is not
essential that the names be either accurate or unambiguous. The public API
exposes a number of name-based maps for client convenience.

The ssa/ssautil package provides various utilities that depend only on the
public API of this package.

TODO(adonovan): Consider the exceptional control-flow implications of defer and
recover().

TODO(adonovan): write a how-to document for all the various cases of trying to
determine corresponding elements across the four domains of source locations,
ast.Nodes, types.Objects, ssa.Values/Instructions.

## Usage

#### func  DefaultType

```go
func DefaultType(typ types.Type) types.Type
```
DefaultType returns the default "typed" type for an "untyped" type; it returns
the incoming type for all other types. The default type for untyped nil is
untyped nil.

Exported to ssa/interp.

TODO(gri): this is a copy of antha/types.defaultType; export that function.

#### func  HasEnclosingFunction

```go
func HasEnclosingFunction(pkg *Package, path []ast.Node) bool
```
HasEnclosingFunction returns true if the AST node denoted by path is contained
within the declaration of some function or package-level variable.

Unlike EnclosingFunction, the behaviour of this function does not depend on
whether SSA code for pkg has been built, so it can be used to quickly reject
check inputs that will cause EnclosingFunction to fail, prior to SSA building.

#### func  WriteFunction

```go
func WriteFunction(buf *bytes.Buffer, f *Function)
```
WriteFunction writes to buf a human-readable "disassembly" of f.

#### func  WritePackage

```go
func WritePackage(buf *bytes.Buffer, p *Package)
```
WritePackage writes to buf a human-readable summary of p.

#### type Alloc

```go
type Alloc struct {
	Comment string
	Heap    bool
}
```

The Alloc instruction reserves space for a value of the given type,
zero-initializes it, and yields its address.

Alloc values are always addresses, and have pointer types, so the type of the
allocated space is actually indirect(Type()).

If Heap is false, Alloc allocates space in the function's activation record
(frame); we refer to an Alloc(Heap=false) as a "local" alloc. Each local Alloc
returns the same address each time it is executed within the same activation;
the space is re-initialized to zero.

If Heap is true, Alloc allocates space in the heap, and returns; we refer to an
Alloc(Heap=true) as a "new" alloc. Each new Alloc returns a different address
each time it is executed.

When Alloc is applied to a channel, map or slice type, it returns the address of
an uninitialized (nil) reference of that kind; store the result of MakeSlice,
MakeMap or MakeChan in that location to instantiate these types.

Pos() returns the ast.CompositeLit.Lbrace for a composite literal, or the
ast.CallExpr.Lparen for a call to new() or for a call that allocates a varargs
slice.

Example printed form:

    t0 = local int
    t1 = new int

#### func (*Alloc) Name

```go
func (v *Alloc) Name() string
```

#### func (*Alloc) Operands

```go
func (v *Alloc) Operands(rands []*Value) []*Value
```

#### func (*Alloc) Pos

```go
func (v *Alloc) Pos() token.Pos
```

#### func (*Alloc) Referrers

```go
func (v *Alloc) Referrers() *[]Instruction
```

#### func (*Alloc) String

```go
func (v *Alloc) String() string
```

#### func (*Alloc) Type

```go
func (v *Alloc) Type() types.Type
```

#### type BasicBlock

```go
type BasicBlock struct {
	Index   int    // index of this block within Parent().Blocks
	Comment string // optional label; no semantic significance

	Instrs       []Instruction // instructions in order
	Preds, Succs []*BasicBlock // predecessors and successors
}
```

An SSA basic block.

The final element of Instrs is always an explicit transfer of control (If, Jump,
Return or Panic).

A block may contain no Instructions only if it is unreachable, i.e. Preds is
nil. Empty blocks are typically pruned.

BasicBlocks and their Preds/Succs relation form a (possibly cyclic) graph
independent of the SSA Value graph: the control-flow graph or CFG. It is illegal
for multiple edges to exist between the same pair of blocks.

Each BasicBlock is also a node in the dominator tree of the CFG. The tree may be
navigated using Idom()/Dominees() and queried using Dominates().

The order of Preds and Succs is significant (to Phi and If instructions,
respectively).

#### func (*BasicBlock) Dominates

```go
func (b *BasicBlock) Dominates(c *BasicBlock) bool
```
Dominates reports whether b dominates c.

#### func (*BasicBlock) Dominees

```go
func (b *BasicBlock) Dominees() []*BasicBlock
```
Dominees returns the list of blocks that b immediately dominates: its children
in the dominator tree.

#### func (*BasicBlock) Idom

```go
func (b *BasicBlock) Idom() *BasicBlock
```
Idom returns the block that immediately dominates b: its parent in the dominator
tree, if any. Neither the entry node (b.Index==0) nor recover node
(b==b.Parent().Recover()) have a parent.

#### func (*BasicBlock) Parent

```go
func (b *BasicBlock) Parent() *Function
```
Parent returns the function that contains block b.

#### func (*BasicBlock) String

```go
func (b *BasicBlock) String() string
```
String returns a human-readable label of this block. It is not guaranteed unique
within the function.

#### type BinOp

```go
type BinOp struct {

	// One of:
	// ADD SUB MUL QUO REM          + - * / %
	// AND OR XOR SHL SHR AND_NOT   & | ^ << >> &~
	// EQL LSS GTR NEQ LEQ GEQ      == != < <= < >=
	Op   token.Token
	X, Y Value
}
```

The BinOp instruction yields the result of binary operation X Op Y.

Pos() returns the ast.BinaryExpr.OpPos, if explicit in the source.

Example printed form:

    t1 = t0 + 1:int

#### func (*BinOp) Name

```go
func (v *BinOp) Name() string
```

#### func (*BinOp) Operands

```go
func (v *BinOp) Operands(rands []*Value) []*Value
```

#### func (*BinOp) Pos

```go
func (v *BinOp) Pos() token.Pos
```

#### func (*BinOp) Referrers

```go
func (v *BinOp) Referrers() *[]Instruction
```

#### func (*BinOp) String

```go
func (v *BinOp) String() string
```

#### func (*BinOp) Type

```go
func (v *BinOp) Type() types.Type
```

#### type BuilderMode

```go
type BuilderMode uint
```

BuilderMode is a bitmask of options for diagnostics and checking.

```go
const (
	LogPackages          BuilderMode = 1 << iota // Dump package inventory to stderr
	LogFunctions                                 // Dump function SSA code to stderr
	LogSource                                    // Show source locations as SSA builder progresses
	SanityCheckFunctions                         // Perform sanity checking of function bodies
	NaiveForm                                    // Build naïve SSA form: don't replace local loads/stores with registers
	BuildSerially                                // Build packages serially, not in parallel.
	GlobalDebug                                  // Enable debug info for all packages
)
```

#### type Builtin

```go
type Builtin struct {
}
```

A Builtin represents a specific use of a built-in function, e.g. len.

Builtins are immutable values. Builtins do not have addresses. Builtins can only
appear in CallCommon.Func.

Object() returns a *types.Builtin.

Type() returns a *types.Signature representing the effective signature of the
built-in for this call.

#### func (*Builtin) Name

```go
func (v *Builtin) Name() string
```

#### func (*Builtin) Object

```go
func (v *Builtin) Object() types.Object
```

#### func (*Builtin) Pos

```go
func (v *Builtin) Pos() token.Pos
```

#### func (*Builtin) Referrers

```go
func (*Builtin) Referrers() *[]Instruction
```

#### func (*Builtin) String

```go
func (v *Builtin) String() string
```

#### func (*Builtin) Type

```go
func (v *Builtin) Type() types.Type
```

#### type Call

```go
type Call struct {
	Call CallCommon
}
```

The Call instruction represents a function or method call.

The Call instruction yields the function result, if there is exactly one, or a
tuple (empty or len>1) whose components are accessed via Extract.

See CallCommon for generic function call documentation.

Pos() returns the ast.CallExpr.Lparen, if explicit in the source.

Example printed form:

    t2 = println(t0, t1)
    t4 = t3()
    t7 = invoke t5.Println(...t6)

#### func (*Call) Common

```go
func (s *Call) Common() *CallCommon
```

#### func (*Call) Name

```go
func (v *Call) Name() string
```

#### func (*Call) Operands

```go
func (s *Call) Operands(rands []*Value) []*Value
```

#### func (*Call) Pos

```go
func (v *Call) Pos() token.Pos
```

#### func (*Call) Referrers

```go
func (v *Call) Referrers() *[]Instruction
```

#### func (*Call) String

```go
func (v *Call) String() string
```

#### func (*Call) Type

```go
func (v *Call) Type() types.Type
```

#### func (*Call) Value

```go
func (s *Call) Value() *Call
```

#### type CallCommon

```go
type CallCommon struct {
	Value  Value       // receiver (invoke mode) or func value (call mode)
	Method *types.Func // abstract method (invoke mode)
	Args   []Value     // actual parameters (in static method call, includes receiver)
}
```

CallCommon is contained by Go, Defer and Call to hold the common parts of a
function or method call.

Each CallCommon exists in one of two modes, function call and interface method
invocation, or "call" and "invoke" for short.

1. "call" mode: when Method is nil (!IsInvoke), a CallCommon represents an
ordinary function call of the value in Value, which may be a *Builtin, a
*Function or any other value of kind 'func'.

Value may be one of:

    (a) a *Function, indicating a statically dispatched call
        to a package-level function, an anonymous function, or
        a method of a named type.
    (b) a *MakeClosure, indicating an immediately applied
        function literal with free variables.
    (c) a *Builtin, indicating a statically dispatched call
        to a built-in function.
    (d) any other value, indicating a dynamically dispatched
        function call.

StaticCallee returns the identity of the callee in cases (a) and (b), nil
otherwise.

Args contains the arguments to the call. If Value is a method, Args[0] contains
the receiver parameter.

Example printed form:

    t2 = println(t0, t1)
    go t3()
    defer t5(...t6)

2. "invoke" mode: when Method is non-nil (IsInvoke), a CallCommon represents a
dynamically dispatched call to an interface method. In this mode, Value is the
interface value and Method is the interface's abstract method. Note: an abstract
method may be shared by multiple interfaces due to embedding; Value.Type()
provides the specific interface used for this call.

Value is implicitly supplied to the concrete method implementation as the
receiver parameter; in other words, Args[0] holds not the receiver but the first
true argument.

Example printed form:

    t1 = invoke t0.String()
    go invoke t3.Run(t2)
    defer invoke t4.Handle(...t5)

For all calls to variadic functions (Signature().Variadic()), the last element
of Args is a slice.

#### func (*CallCommon) Description

```go
func (c *CallCommon) Description() string
```
Description returns a description of the mode of this call suitable for a user
interface, e.g. "static method call".

#### func (*CallCommon) IsInvoke

```go
func (c *CallCommon) IsInvoke() bool
```
IsInvoke returns true if this call has "invoke" (not "call") mode.

#### func (*CallCommon) Operands

```go
func (c *CallCommon) Operands(rands []*Value) []*Value
```

#### func (*CallCommon) Pos

```go
func (c *CallCommon) Pos() token.Pos
```

#### func (*CallCommon) Signature

```go
func (c *CallCommon) Signature() *types.Signature
```
Signature returns the signature of the called function.

For an "invoke"-mode call, the signature of the interface method is returned.

In either "call" or "invoke" mode, if the callee is a method, its receiver is
represented by sig.Recv, not sig.Params().At(0).

#### func (*CallCommon) StaticCallee

```go
func (c *CallCommon) StaticCallee() *Function
```
StaticCallee returns the callee if this is a trivially static "call"-mode call
to a function.

#### func (*CallCommon) String

```go
func (c *CallCommon) String() string
```

#### type CallInstruction

```go
type CallInstruction interface {
	Instruction
	Common() *CallCommon // returns the common parts of the call
	Value() *Call        // returns the result value of the call (*Call) or nil (*Go, *Defer)
}
```

The CallInstruction interface, implemented by *Go, *Defer and *Call, exposes the
common parts of function calling instructions, yet provides a way back to the
Value defined by *Call alone.

#### type Capture

```go
type Capture struct {
}
```

A Capture represents a free variable of the function to which it belongs.

Captures are used to implement anonymous functions, whose free variables are
lexically captured in a closure formed by MakeClosure. The referent of such a
capture is an Alloc or another Capture and is considered a potentially escaping
heap address, with pointer type.

Captures are also used to implement bound method closures. Such a capture
represents the receiver value and may be of any type that has concrete methods.

Pos() returns the position of the value that was captured, which belongs to an
enclosing function.

#### func (*Capture) Name

```go
func (v *Capture) Name() string
```

#### func (*Capture) Parent

```go
func (v *Capture) Parent() *Function
```

#### func (*Capture) Pos

```go
func (v *Capture) Pos() token.Pos
```

#### func (*Capture) Referrers

```go
func (v *Capture) Referrers() *[]Instruction
```

#### func (*Capture) String

```go
func (v *Capture) String() string
```

#### func (*Capture) Type

```go
func (v *Capture) Type() types.Type
```

#### type ChangeInterface

```go
type ChangeInterface struct {
	X Value
}
```

ChangeInterface constructs a value of one interface type from a value of another
interface type known to be assignable to it. This operation cannot fail.

Pos() returns the ast.CallExpr.Lparen if the instruction arose from an explicit
T(e) conversion; the ast.TypeAssertExpr.Lparen if the instruction arose from an
explicit e.(T) operation; or token.NoPos otherwise.

Example printed form:

    t1 = change interface interface{} <- I (t0)

#### func (*ChangeInterface) Name

```go
func (v *ChangeInterface) Name() string
```

#### func (*ChangeInterface) Operands

```go
func (v *ChangeInterface) Operands(rands []*Value) []*Value
```

#### func (*ChangeInterface) Pos

```go
func (v *ChangeInterface) Pos() token.Pos
```

#### func (*ChangeInterface) Referrers

```go
func (v *ChangeInterface) Referrers() *[]Instruction
```

#### func (*ChangeInterface) String

```go
func (v *ChangeInterface) String() string
```

#### func (*ChangeInterface) Type

```go
func (v *ChangeInterface) Type() types.Type
```

#### type ChangeType

```go
type ChangeType struct {
	X Value
}
```

The ChangeType instruction applies to X a value-preserving type change to
Type().

Type changes are permitted:

    - between a named type and its underlying type.
    - between two named types of the same underlying type.
    - between (possibly named) pointers to identical base types.
    - between f(T) functions and (T) func f() methods.
    - from a bidirectional channel to a read- or write-channel,
      optionally adding/removing a name.

This operation cannot fail dynamically.

Pos() returns the ast.CallExpr.Lparen, if the instruction arose from an explicit
conversion in the source.

Example printed form:

    t1 = changetype *int <- IntPtr (t0)

#### func (*ChangeType) Name

```go
func (v *ChangeType) Name() string
```

#### func (*ChangeType) Operands

```go
func (v *ChangeType) Operands(rands []*Value) []*Value
```

#### func (*ChangeType) Pos

```go
func (v *ChangeType) Pos() token.Pos
```

#### func (*ChangeType) Referrers

```go
func (v *ChangeType) Referrers() *[]Instruction
```

#### func (*ChangeType) String

```go
func (v *ChangeType) String() string
```

#### func (*ChangeType) Type

```go
func (v *ChangeType) Type() types.Type
```

#### type Const

```go
type Const struct {
	Value exact.Value
}
```

A Const represents the value of a constant expression.

The underlying type of a constant may be any boolean, numeric, or string type.
In addition, a Const may represent the nil value of any reference type:
interface, map, channel, pointer, slice, or function---but not "untyped nil".

All source-level constant expressions are represented by a Const of equal type
and value.

Value holds the exact value of the constant, independent of its Type(), using
the same representation as package antha/exact uses for constants, or nil for a
typed nil value.

Pos() returns token.NoPos.

Example printed form:

    42:int
    "hello":untyped string
    3+4i:MyComplex

#### func  NewConst

```go
func NewConst(val exact.Value, typ types.Type) *Const
```
NewConst returns a new constant of the specified value and type. val must be
valid according to the specification of Const.Value.

#### func (*Const) Complex128

```go
func (c *Const) Complex128() complex128
```
Complex128 returns the complex value of this constant truncated to fit a
complex128.

#### func (*Const) Float64

```go
func (c *Const) Float64() float64
```
Float64 returns the numeric value of this constant truncated to fit a float64.

#### func (*Const) Int64

```go
func (c *Const) Int64() int64
```
Int64 returns the numeric value of this constant truncated to fit a signed
64-bit integer.

#### func (*Const) IsNil

```go
func (c *Const) IsNil() bool
```
IsNil returns true if this constant represents a typed or untyped nil value.

#### func (*Const) Name

```go
func (c *Const) Name() string
```

#### func (*Const) Pos

```go
func (c *Const) Pos() token.Pos
```

#### func (*Const) Referrers

```go
func (c *Const) Referrers() *[]Instruction
```

#### func (*Const) RelString

```go
func (c *Const) RelString(from *types.Package) string
```

#### func (*Const) String

```go
func (c *Const) String() string
```

#### func (*Const) Type

```go
func (c *Const) Type() types.Type
```

#### func (*Const) Uint64

```go
func (c *Const) Uint64() uint64
```
Uint64 returns the numeric value of this constant truncated to fit an unsigned
64-bit integer.

#### type Convert

```go
type Convert struct {
	X Value
}
```

The Convert instruction yields the conversion of value X to type Type(). One or
both of those types is basic (but possibly named).

A conversion may change the value and representation of its operand. Conversions
are permitted:

    - between real numeric types.
    - between complex numeric types.
    - between string and []byte or []rune.
    - between pointers and unsafe.Pointer.
    - between unsafe.Pointer and uintptr.
    - from (Unicode) integer to (UTF-8) string.

A conversion may imply a type name change also.

This operation cannot fail dynamically.

Conversions of untyped string/number/bool constants to a specific representation
are eliminated during SSA construction.

Pos() returns the ast.CallExpr.Lparen, if the instruction arose from an explicit
conversion in the source.

Example printed form:

    t1 = convert []byte <- string (t0)

#### func (*Convert) Name

```go
func (v *Convert) Name() string
```

#### func (*Convert) Operands

```go
func (v *Convert) Operands(rands []*Value) []*Value
```

#### func (*Convert) Pos

```go
func (v *Convert) Pos() token.Pos
```

#### func (*Convert) Referrers

```go
func (v *Convert) Referrers() *[]Instruction
```

#### func (*Convert) String

```go
func (v *Convert) String() string
```

#### func (*Convert) Type

```go
func (v *Convert) Type() types.Type
```

#### type DebugRef

```go
type DebugRef struct {
	Expr ast.Expr // the referring expression (never *ast.ParenExpr)

	IsAddr bool  // Expr is addressable and X is the address it denotes
	X      Value // the value or address of Expr
}
```

A DebugRef instruction maps a source-level expression Expr to the SSA value X
that represents the value (!IsAddr) or address (IsAddr) of that expression.

DebugRef is a pseudo-instruction: it has no dynamic effect.

Pos() returns Expr.Pos(), the start position of the source-level expression.
This is not the same as the "designated" token as documented at Value.Pos().
e.g. CallExpr.Pos() does not return the position of the ("designated") Lparen
token.

If Expr is an *ast.Ident denoting a var or func, Object() returns the object;
though this information can be obtained from the type checker, including it here
greatly facilitates debugging. For non-Ident expressions, Object() returns nil.

DebugRefs are generated only for functions built with debugging enabled; see
Package.SetDebugMode() and the GlobalDebug builder mode flag.

DebugRefs are not emitted for ast.Idents referring to constants or predeclared
identifiers, since they are trivial and numerous. Nor are they emitted for
ast.ParenExprs.

(By representing these as instructions, rather than out-of-band, consistency is
maintained during transformation passes by the ordinary SSA renaming machinery.)

Example printed form:

    ; *ast.CallExpr @ 102:9 is t5
    ; var x float64 @ 109:72 is x
    ; address of *ast.CompositeLit @ 216:10 is t0

#### func (*DebugRef) Block

```go
func (v *DebugRef) Block() *BasicBlock
```

#### func (*DebugRef) Operands

```go
func (s *DebugRef) Operands(rands []*Value) []*Value
```

#### func (*DebugRef) Parent

```go
func (v *DebugRef) Parent() *Function
```

#### func (*DebugRef) Pos

```go
func (s *DebugRef) Pos() token.Pos
```

#### func (*DebugRef) String

```go
func (s *DebugRef) String() string
```

#### type Defer

```go
type Defer struct {
	Call CallCommon
}
```

The Defer instruction pushes the specified call onto a stack of functions to be
called by a RunDefers instruction or by a panic.

See CallCommon for generic function call documentation.

Pos() returns the ast.DeferStmt.Defer.

Example printed form:

    defer println(t0, t1)
    defer t3()
    defer invoke t5.Println(...t6)

#### func (*Defer) Block

```go
func (v *Defer) Block() *BasicBlock
```

#### func (*Defer) Common

```go
func (s *Defer) Common() *CallCommon
```

#### func (*Defer) Operands

```go
func (s *Defer) Operands(rands []*Value) []*Value
```

#### func (*Defer) Parent

```go
func (v *Defer) Parent() *Function
```

#### func (*Defer) Pos

```go
func (s *Defer) Pos() token.Pos
```

#### func (*Defer) String

```go
func (s *Defer) String() string
```

#### func (*Defer) Value

```go
func (s *Defer) Value() *Call
```

#### type Extract

```go
type Extract struct {
	Tuple Value
	Index int
}
```

The Extract instruction yields component Index of Tuple.

This is used to access the results of instructions with multiple return values,
such as Call, TypeAssert, Next, UnOp(ARROW) and IndexExpr(Map).

Example printed form:

    t1 = extract t0 #1

#### func (*Extract) Name

```go
func (v *Extract) Name() string
```

#### func (*Extract) Operands

```go
func (v *Extract) Operands(rands []*Value) []*Value
```

#### func (*Extract) Pos

```go
func (v *Extract) Pos() token.Pos
```

#### func (*Extract) Referrers

```go
func (v *Extract) Referrers() *[]Instruction
```

#### func (*Extract) String

```go
func (v *Extract) String() string
```

#### func (*Extract) Type

```go
func (v *Extract) Type() types.Type
```

#### type Field

```go
type Field struct {
	X     Value // struct
	Field int   // index into X.Type().(*types.Struct).Fields
}
```

The Field instruction yields the Field of struct X.

The field is identified by its index within the field list of the struct type of
X; by using numeric indices we avoid ambiguity of package-local identifiers and
permit compact representations.

Pos() returns the position of the ast.SelectorExpr.Sel for the field, if
explicit in the source.

Example printed form:

    t1 = t0.name [#1]

#### func (*Field) Name

```go
func (v *Field) Name() string
```

#### func (*Field) Operands

```go
func (v *Field) Operands(rands []*Value) []*Value
```

#### func (*Field) Pos

```go
func (v *Field) Pos() token.Pos
```

#### func (*Field) Referrers

```go
func (v *Field) Referrers() *[]Instruction
```

#### func (*Field) String

```go
func (v *Field) String() string
```

#### func (*Field) Type

```go
func (v *Field) Type() types.Type
```

#### type FieldAddr

```go
type FieldAddr struct {
	X     Value // *struct
	Field int   // index into X.Type().Deref().(*types.Struct).Fields
}
```

The FieldAddr instruction yields the address of Field of *struct X.

The field is identified by its index within the field list of the struct type of
X.

Dynamically, this instruction panics if X evaluates to a nil pointer.

Type() returns a (possibly named) *types.Pointer.

Pos() returns the position of the ast.SelectorExpr.Sel for the field, if
explicit in the source.

Example printed form:

    t1 = &t0.name [#1]

#### func (*FieldAddr) Name

```go
func (v *FieldAddr) Name() string
```

#### func (*FieldAddr) Operands

```go
func (v *FieldAddr) Operands(rands []*Value) []*Value
```

#### func (*FieldAddr) Pos

```go
func (v *FieldAddr) Pos() token.Pos
```

#### func (*FieldAddr) Referrers

```go
func (v *FieldAddr) Referrers() *[]Instruction
```

#### func (*FieldAddr) String

```go
func (v *FieldAddr) String() string
```

#### func (*FieldAddr) Type

```go
func (v *FieldAddr) Type() types.Type
```

#### type Function

```go
type Function struct {
	Signature *types.Signature

	Synthetic string // provenance of synthetic function; "" for true source functions

	Enclosing *Function    // enclosing function if anon; nil if global
	Pkg       *Package     // enclosing package; nil for shared funcs (wrappers and error.Error)
	Prog      *Program     // enclosing program
	Params    []*Parameter // function parameters; for methods, includes receiver
	FreeVars  []*Capture   // free variables whose values must be supplied by closure
	Locals    []*Alloc
	Blocks    []*BasicBlock // basic blocks of the function; nil => external
	Recover   *BasicBlock   // optional; control transfers here after recovered panic
	AnonFuncs []*Function   // anonymous functions directly beneath this one
}
```

Function represents the parameters, results and code of a function or method.

If Blocks is nil, this indicates an external function for which no Go source
code is available. In this case, FreeVars and Locals will be nil too. Clients
performing whole-program analysis must handle external functions specially.

Functions are immutable values; they do not have addresses.

Blocks contains the function's control-flow graph (CFG). Blocks[0] is the
function entry point; block order is not otherwise semantically significant,
though it may affect the readability of the disassembly. To iterate over the
blocks in dominance order, use DomPreorder().

Recover is an optional second entry point to which control resumes after a
recovered panic. The Recover block may contain only a load of the function's
named return parameters followed by a return of the loaded values.

A nested function that refers to one or more lexically enclosing local variables
("free variables") has Capture parameters. Such functions cannot be called
directly but require a value created by MakeClosure which, via its Bindings,
supplies values for these parameters.

If the function is a method (Signature.Recv() != nil) then the first element of
Params is the receiver parameter.

Pos() returns the declaring ast.FuncLit.Type.Func or the position of the
ast.FuncDecl.Name, if the function was explicit in the source. Synthetic
wrappers, for which Synthetic != "", may share the same position as the function
they wrap.

Type() returns the function's Signature.

#### func  EnclosingFunction

```go
func EnclosingFunction(pkg *Package, path []ast.Node) *Function
```
EnclosingFunction returns the function that contains the syntax node denoted by
path.

Syntax associated with package-level variable specifications is enclosed by the
package's init() function.

Returns nil if not found; reasons might include:

    - the node is not enclosed by any function.
    - the node is within an anonymous function (FuncLit) and
      its SSA function has not been created yet
      (pkg.Build() has not yet been called).

#### func  NewFunction

```go
func NewFunction(name string, sig *types.Signature, provenance string) *Function
```
NewFunction returns a new synthetic Function instance with its name and
signature fields set as specified.

The caller is responsible for initializing the remaining fields of the function
object, e.g. Pkg, Prog, Params, Blocks.

It is practically impossible for clients to construct well-formed SSA
functions/packages/programs directly, so we assume this is the job of the
Builder alone. NewFunction exists to provide clients a little flexibility. For
example, analysis tools may wish to construct fake Functions for the root of the
callgraph, a fake "reflect" package, etc.

TODO(adonovan): think harder about the API here.

#### func (*Function) DomPreorder

```go
func (f *Function) DomPreorder() []*BasicBlock
```
DomPreorder returns a new slice containing the blocks of f in dominator tree
preorder.

#### func (*Function) Name

```go
func (v *Function) Name() string
```

#### func (*Function) Object

```go
func (v *Function) Object() types.Object
```

#### func (*Function) Package

```go
func (v *Function) Package() *Package
```

#### func (*Function) Pos

```go
func (v *Function) Pos() token.Pos
```

#### func (*Function) Referrers

```go
func (v *Function) Referrers() *[]Instruction
```

#### func (*Function) RelString

```go
func (f *Function) RelString(from *types.Package) string
```
RelString returns the full name of this function, qualified by package name,
receiver type, etc.

The specific formatting rules are not guaranteed and may change.

Examples:

    "math.IsNaN"                // a package-level function
    "IsNaN"                     // intra-package reference to same
    "(*sync.WaitGroup).Add"     // a declared method
    "(*Return).Block"           // a promotion wrapper method (intra-package ref)
    "(Instruction).Block"       // an interface method wrapper (intra-package ref)
    "main$1"                    // an anonymous function
    "init$1"                    // a declared init function
    "init"                      // the synthesized package initializer
    "bound$(*T).f"              // a bound method wrapper

If from==f.Pkg, suppress package qualification.

#### func (*Function) String

```go
func (v *Function) String() string
```

#### func (*Function) Syntax

```go
func (f *Function) Syntax() ast.Node
```
Syntax returns an ast.Node whose Pos/End methods provide the lexical extent of
the function if it was defined by Go source code (f.Synthetic==""), or nil
otherwise.

If f was built with debug information (see Package.SetDebugRef), the result is
the *ast.FuncDecl or *ast.FuncLit that declared the function. Otherwise, it is
an opaque Node providing only position information; this avoids pinning the AST
in memory.

#### func (*Function) Token

```go
func (v *Function) Token() token.Token
```

#### func (*Function) Type

```go
func (v *Function) Type() types.Type
```

#### func (*Function) ValueForExpr

```go
func (f *Function) ValueForExpr(e ast.Expr) (value Value, isAddr bool)
```
ValueForExpr returns the SSA Value that corresponds to non-constant expression
e.

It returns nil if no value was found, e.g.

    - the expression is not lexically contained within f;
    - f was not built with debug information; or
    - e is a constant expression.  (For efficiency, no debug
      information is stored for constants. Use
      loader.PackageInfo.ValueOf(e) instead.)
    - e is a reference to nil or a built-in function.
    - the value was optimised away.

If e is an addressable expression used an an lvalue context, value is the
address denoted by e, and isAddr is true.

The types of e (or &e, if isAddr) and the result are equal (modulo "untyped"
bools resulting from comparisons).

(Tip: to find the ssa.Value given a source position, use
importer.PathEnclosingInterval to locate the ast.Node, then EnclosingFunction to
locate the Function, then ValueForExpr to find the ssa.Value.)

#### func (*Function) WriteTo

```go
func (f *Function) WriteTo(w io.Writer) (int64, error)
```

#### type Global

```go
type Global struct {
	Pkg *Package
}
```

A Global is a named Value holding the address of a package-level variable.

Pos() returns the position of the ast.ValueSpec.Names[*] identifier.

#### func (*Global) Name

```go
func (v *Global) Name() string
```

#### func (*Global) Object

```go
func (v *Global) Object() types.Object
```

#### func (*Global) Package

```go
func (v *Global) Package() *Package
```

#### func (*Global) Pos

```go
func (v *Global) Pos() token.Pos
```

#### func (*Global) Referrers

```go
func (v *Global) Referrers() *[]Instruction
```

#### func (*Global) RelString

```go
func (v *Global) RelString(from *types.Package) string
```

#### func (*Global) String

```go
func (v *Global) String() string
```

#### func (*Global) Token

```go
func (v *Global) Token() token.Token
```

#### func (*Global) Type

```go
func (v *Global) Type() types.Type
```

#### type Go

```go
type Go struct {
	Call CallCommon
}
```

The Go instruction creates a new goroutine and calls the specified function
within it.

See CallCommon for generic function call documentation.

Pos() returns the ast.GoStmt.Go.

Example printed form:

    go println(t0, t1)
    go t3()
    go invoke t5.Println(...t6)

#### func (*Go) Block

```go
func (v *Go) Block() *BasicBlock
```

#### func (*Go) Common

```go
func (s *Go) Common() *CallCommon
```

#### func (*Go) Operands

```go
func (s *Go) Operands(rands []*Value) []*Value
```

#### func (*Go) Parent

```go
func (v *Go) Parent() *Function
```

#### func (*Go) Pos

```go
func (s *Go) Pos() token.Pos
```

#### func (*Go) String

```go
func (s *Go) String() string
```

#### func (*Go) Value

```go
func (s *Go) Value() *Call
```

#### type If

```go
type If struct {
	Cond Value
}
```

The If instruction transfers control to one of the two successors of its owning
block, depending on the boolean Cond: the first if true, the second if false.

An If instruction must be the last instruction of its containing BasicBlock.

Pos() returns NoPos.

Example printed form:

    if t0 goto done else body

#### func (*If) Block

```go
func (v *If) Block() *BasicBlock
```

#### func (*If) Operands

```go
func (s *If) Operands(rands []*Value) []*Value
```

#### func (*If) Parent

```go
func (v *If) Parent() *Function
```

#### func (*If) Pos

```go
func (s *If) Pos() token.Pos
```

#### func (*If) String

```go
func (s *If) String() string
```

#### type Index

```go
type Index struct {
	X     Value // array
	Index Value // integer index
}
```

The Index instruction yields element Index of array X.

Pos() returns the ast.IndexExpr.Lbrack for the index operation, if explicit in
the source.

Example printed form:

    t2 = t0[t1]

#### func (*Index) Name

```go
func (v *Index) Name() string
```

#### func (*Index) Operands

```go
func (v *Index) Operands(rands []*Value) []*Value
```

#### func (*Index) Pos

```go
func (v *Index) Pos() token.Pos
```

#### func (*Index) Referrers

```go
func (v *Index) Referrers() *[]Instruction
```

#### func (*Index) String

```go
func (v *Index) String() string
```

#### func (*Index) Type

```go
func (v *Index) Type() types.Type
```

#### type IndexAddr

```go
type IndexAddr struct {
	X     Value // slice or *array,
	Index Value // numeric index
}
```

The IndexAddr instruction yields the address of the element at index Index of
collection X. Index is an integer expression.

The elements of maps and strings are not addressable; use Lookup or MapUpdate
instead.

Dynamically, this instruction panics if X evaluates to a nil *array pointer.

Type() returns a (possibly named) *types.Pointer.

Pos() returns the ast.IndexExpr.Lbrack for the index operation, if explicit in
the source.

Example printed form:

    t2 = &t0[t1]

#### func (*IndexAddr) Name

```go
func (v *IndexAddr) Name() string
```

#### func (*IndexAddr) Operands

```go
func (v *IndexAddr) Operands(rands []*Value) []*Value
```

#### func (*IndexAddr) Pos

```go
func (v *IndexAddr) Pos() token.Pos
```

#### func (*IndexAddr) Referrers

```go
func (v *IndexAddr) Referrers() *[]Instruction
```

#### func (*IndexAddr) String

```go
func (v *IndexAddr) String() string
```

#### func (*IndexAddr) Type

```go
func (v *IndexAddr) Type() types.Type
```

#### type Instruction

```go
type Instruction interface {
	// String returns the disassembled form of this value.  e.g.
	//
	// Examples of Instructions that define a Value:
	// e.g.  "x + y"     (BinOp)
	//       "len([])"   (Call)
	// Note that the name of the Value is not printed.
	//
	// Examples of Instructions that do define (are) Values:
	// e.g.  "return x"  (Return)
	//       "*y = x"    (Store)
	//
	// (This separation is useful for some analyses which
	// distinguish the operation from the value it
	// defines. e.g. 'y = local int' is both an allocation of
	// memory 'local int' and a definition of a pointer y.)
	String() string

	// Parent returns the function to which this instruction
	// belongs.
	Parent() *Function

	// Block returns the basic block to which this instruction
	// belongs.
	Block() *BasicBlock

	// Operands returns the operands of this instruction: the
	// set of Values it references.
	//
	// Specifically, it appends their addresses to rands, a
	// user-provided slice, and returns the resulting slice,
	// permitting avoidance of memory allocation.
	//
	// The operands are appended in undefined order; the addresses
	// are always non-nil but may point to a nil Value.  Clients
	// may store through the pointers, e.g. to effect a value
	// renaming.
	//
	// Value.Referrers is a subset of the inverse of this
	// relation.  (Referrers are not tracked for all types of
	// Values.)
	Operands(rands []*Value) []*Value

	// Pos returns the location of the AST token most closely
	// associated with the operation that gave rise to this
	// instruction, or token.NoPos if it was not explicit in the
	// source.
	//
	// For each ast.Node type, a particular token is designated as
	// the closest location for the expression, e.g. the Go token
	// for an *ast.GoStmt.  This permits a compact but approximate
	// mapping from Instructions to source positions for use in
	// diagnostic messages, for example.
	//
	// (Do not use this position to determine which Instruction
	// corresponds to an ast.Expr; see the notes for Value.Pos.
	// This position may be used to determine which non-Value
	// Instruction corresponds to some ast.Stmts, but not all: If
	// and Jump instructions have no Pos(), for example.)
	//
	Pos() token.Pos
	// contains filtered or unexported methods
}
```

An Instruction is an SSA instruction that computes a new Value or has some
effect.

An Instruction that defines a value (e.g. BinOp) also implements the Value
interface; an Instruction that only has an effect (e.g. Store) does not.

#### type Jump

```go
type Jump struct {
}
```

The Jump instruction transfers control to the sole successor of its owning
block.

A Jump must be the last instruction of its containing BasicBlock.

Pos() returns NoPos.

Example printed form:

    jump done

#### func (*Jump) Block

```go
func (v *Jump) Block() *BasicBlock
```

#### func (*Jump) Operands

```go
func (*Jump) Operands(rands []*Value) []*Value
```

#### func (*Jump) Parent

```go
func (v *Jump) Parent() *Function
```

#### func (*Jump) Pos

```go
func (s *Jump) Pos() token.Pos
```

#### func (*Jump) String

```go
func (s *Jump) String() string
```

#### type Lookup

```go
type Lookup struct {
	X       Value // string or map
	Index   Value // numeric or key-typed index
	CommaOk bool  // return a value,ok pair
}
```

The Lookup instruction yields element Index of collection X, a map or string.
Index is an integer expression if X is a string or the appropriate key type if X
is a map.

If CommaOk, the result is a 2-tuple of the value above and a boolean indicating
the result of a map membership test for the key. The components of the tuple are
accessed using Extract.

Pos() returns the ast.IndexExpr.Lbrack, if explicit in the source.

Example printed form:

    t2 = t0[t1]
    t5 = t3[t4],ok

#### func (*Lookup) Name

```go
func (v *Lookup) Name() string
```

#### func (*Lookup) Operands

```go
func (v *Lookup) Operands(rands []*Value) []*Value
```

#### func (*Lookup) Pos

```go
func (v *Lookup) Pos() token.Pos
```

#### func (*Lookup) Referrers

```go
func (v *Lookup) Referrers() *[]Instruction
```

#### func (*Lookup) String

```go
func (v *Lookup) String() string
```

#### func (*Lookup) Type

```go
func (v *Lookup) Type() types.Type
```

#### type MakeChan

```go
type MakeChan struct {
	Size Value // int; size of buffer; zero => synchronous.
}
```

The MakeChan instruction creates a new channel object and yields a value of kind
chan.

Type() returns a (possibly named) *types.Chan.

Pos() returns the ast.CallExpr.Lparen for the make(chan) that created it.

Example printed form:

    t0 = make chan int 0
    t0 = make IntChan 0

#### func (*MakeChan) Name

```go
func (v *MakeChan) Name() string
```

#### func (*MakeChan) Operands

```go
func (v *MakeChan) Operands(rands []*Value) []*Value
```

#### func (*MakeChan) Pos

```go
func (v *MakeChan) Pos() token.Pos
```

#### func (*MakeChan) Referrers

```go
func (v *MakeChan) Referrers() *[]Instruction
```

#### func (*MakeChan) String

```go
func (v *MakeChan) String() string
```

#### func (*MakeChan) Type

```go
func (v *MakeChan) Type() types.Type
```

#### type MakeClosure

```go
type MakeClosure struct {
	Fn       Value   // always a *Function
	Bindings []Value // values for each free variable in Fn.FreeVars
}
```

The MakeClosure instruction yields a closure value whose code is Fn and whose
free variables' values are supplied by Bindings.

Type() returns a (possibly named) *types.Signature.

Pos() returns the ast.FuncLit.Type.Func for a function literal closure or the
ast.SelectorExpr.Sel for a bound method closure.

Example printed form:

    t0 = make closure anon@1.2 [x y z]
    t1 = make closure bound$(main.I).add [i]

#### func (*MakeClosure) Name

```go
func (v *MakeClosure) Name() string
```

#### func (*MakeClosure) Operands

```go
func (v *MakeClosure) Operands(rands []*Value) []*Value
```

#### func (*MakeClosure) Pos

```go
func (v *MakeClosure) Pos() token.Pos
```

#### func (*MakeClosure) Referrers

```go
func (v *MakeClosure) Referrers() *[]Instruction
```

#### func (*MakeClosure) String

```go
func (v *MakeClosure) String() string
```

#### func (*MakeClosure) Type

```go
func (v *MakeClosure) Type() types.Type
```

#### type MakeInterface

```go
type MakeInterface struct {
	X Value
}
```

MakeInterface constructs an instance of an interface type from a value of a
concrete type.

Use Program.MethodSets.MethodSet(X.Type()) to find the method-set of X, and
Program.Method(m) to find the implementation of a method.

To construct the zero value of an interface type T, use:

    NewConst(exact.MakeNil(), T, pos)

Pos() returns the ast.CallExpr.Lparen, if the instruction arose from an explicit
conversion in the source.

Example printed form:

    t1 = make interface{} <- int (42:int)
    t2 = make Stringer <- t0

#### func (*MakeInterface) Name

```go
func (v *MakeInterface) Name() string
```

#### func (*MakeInterface) Operands

```go
func (v *MakeInterface) Operands(rands []*Value) []*Value
```

#### func (*MakeInterface) Pos

```go
func (v *MakeInterface) Pos() token.Pos
```

#### func (*MakeInterface) Referrers

```go
func (v *MakeInterface) Referrers() *[]Instruction
```

#### func (*MakeInterface) String

```go
func (v *MakeInterface) String() string
```

#### func (*MakeInterface) Type

```go
func (v *MakeInterface) Type() types.Type
```

#### type MakeMap

```go
type MakeMap struct {
	Reserve Value // initial space reservation; nil => default
}
```

The MakeMap instruction creates a new hash-table-based map object and yields a
value of kind map.

Type() returns a (possibly named) *types.Map.

Pos() returns the ast.CallExpr.Lparen, if created by make(map), or the
ast.CompositeLit.Lbrack if created by a literal.

Example printed form:

    t1 = make map[string]int t0
    t1 = make StringIntMap t0

#### func (*MakeMap) Name

```go
func (v *MakeMap) Name() string
```

#### func (*MakeMap) Operands

```go
func (v *MakeMap) Operands(rands []*Value) []*Value
```

#### func (*MakeMap) Pos

```go
func (v *MakeMap) Pos() token.Pos
```

#### func (*MakeMap) Referrers

```go
func (v *MakeMap) Referrers() *[]Instruction
```

#### func (*MakeMap) String

```go
func (v *MakeMap) String() string
```

#### func (*MakeMap) Type

```go
func (v *MakeMap) Type() types.Type
```

#### type MakeSlice

```go
type MakeSlice struct {
	Len Value
	Cap Value
}
```

The MakeSlice instruction yields a slice of length Len backed by a newly
allocated array of length Cap.

Both Len and Cap must be non-nil Values of integer type.

(Alloc(types.Array) followed by Slice will not suffice because Alloc can only
create arrays of constant length.)

Type() returns a (possibly named) *types.Slice.

Pos() returns the ast.CallExpr.Lparen for the make([]T) that created it.

Example printed form:

    t1 = make []string 1:int t0
    t1 = make StringSlice 1:int t0

#### func (*MakeSlice) Name

```go
func (v *MakeSlice) Name() string
```

#### func (*MakeSlice) Operands

```go
func (v *MakeSlice) Operands(rands []*Value) []*Value
```

#### func (*MakeSlice) Pos

```go
func (v *MakeSlice) Pos() token.Pos
```

#### func (*MakeSlice) Referrers

```go
func (v *MakeSlice) Referrers() *[]Instruction
```

#### func (*MakeSlice) String

```go
func (v *MakeSlice) String() string
```

#### func (*MakeSlice) Type

```go
func (v *MakeSlice) Type() types.Type
```

#### type MapUpdate

```go
type MapUpdate struct {
	Map   Value
	Key   Value
	Value Value
}
```

The MapUpdate instruction updates the association of Map[Key] to Value.

Pos() returns the ast.KeyValueExpr.Colon or ast.IndexExpr.Lbrack, if explicit in
the source.

Example printed form:

    t0[t1] = t2

#### func (*MapUpdate) Block

```go
func (v *MapUpdate) Block() *BasicBlock
```

#### func (*MapUpdate) Operands

```go
func (v *MapUpdate) Operands(rands []*Value) []*Value
```

#### func (*MapUpdate) Parent

```go
func (v *MapUpdate) Parent() *Function
```

#### func (*MapUpdate) Pos

```go
func (s *MapUpdate) Pos() token.Pos
```

#### func (*MapUpdate) String

```go
func (s *MapUpdate) String() string
```

#### type Member

```go
type Member interface {
	Name() string                    // declared name of the package member
	String() string                  // package-qualified name of the package member
	RelString(*types.Package) string // like String, but relative refs are unqualified
	Object() types.Object            // typechecker's object for this member, if any
	Pos() token.Pos                  // position of member's declaration, if known
	Type() types.Type                // type of the package member
	Token() token.Token              // token.{VAR,FUNC,CONST,TYPE}
	Package() *Package               // returns the containing package. (TODO: rename Pkg)
}
```

A Member is a member of a Go package, implemented by *NamedConst, *Global,
*Function, or *Type; they are created by package-level const, var, func and type
declarations respectively.

#### type NamedConst

```go
type NamedConst struct {
	Value *Const
}
```

A NamedConst is a Member of Package representing a package-level named constant
value.

Pos() returns the position of the declaring ast.ValueSpec.Names[*] identifier.

NB: a NamedConst is not a Value; it contains a constant Value, which it augments
with the name and position of its 'const' declaration.

#### func (*NamedConst) Name

```go
func (c *NamedConst) Name() string
```

#### func (*NamedConst) Object

```go
func (c *NamedConst) Object() types.Object
```

#### func (*NamedConst) Package

```go
func (c *NamedConst) Package() *Package
```

#### func (*NamedConst) Pos

```go
func (c *NamedConst) Pos() token.Pos
```

#### func (*NamedConst) RelString

```go
func (c *NamedConst) RelString(from *types.Package) string
```

#### func (*NamedConst) String

```go
func (c *NamedConst) String() string
```

#### func (*NamedConst) Token

```go
func (c *NamedConst) Token() token.Token
```

#### func (*NamedConst) Type

```go
func (c *NamedConst) Type() types.Type
```

#### type Next

```go
type Next struct {
	Iter     Value
	IsString bool // true => string iterator; false => map iterator.
}
```

The Next instruction reads and advances the (map or string) iterator Iter and
returns a 3-tuple value (ok, k, v). If the iterator is not exhausted, ok is true
and k and v are the next elements of the domain and range, respectively.
Otherwise ok is false and k and v are undefined.

Components of the tuple are accessed using Extract.

The IsString field distinguishes iterators over strings from those over maps, as
the Type() alone is insufficient: consider map[int]rune.

Type() returns a *types.Tuple for the triple (ok, k, v). The types of k and/or v
may be types.Invalid.

Example printed form:

    t1 = next t0

#### func (*Next) Name

```go
func (v *Next) Name() string
```

#### func (*Next) Operands

```go
func (v *Next) Operands(rands []*Value) []*Value
```

#### func (*Next) Pos

```go
func (v *Next) Pos() token.Pos
```

#### func (*Next) Referrers

```go
func (v *Next) Referrers() *[]Instruction
```

#### func (*Next) String

```go
func (v *Next) String() string
```

#### func (*Next) Type

```go
func (v *Next) Type() types.Type
```

#### type Package

```go
type Package struct {
	Prog    *Program          // the owning program
	Object  *types.Package    // the type checker's package object for this package
	Members map[string]Member // all package members keyed by name
}
```

A Package is a single analyzed Go package containing Members for all
package-level functions, variables, constants and types it declares. These may
be accessed directly via Members, or via the type-specific accessor methods
Func, Type, Var and Const.

#### func (*Package) Build

```go
func (p *Package) Build()
```
Build builds SSA code for all functions and vars in package p.

Precondition: CreatePackage must have been called for all of p's direct imports
(and hence its direct imports must have been error-free).

Build is idempotent and thread-safe.

#### func (*Package) Const

```go
func (p *Package) Const(name string) (c *NamedConst)
```
Const returns the package-level constant of the specified name, or nil if not
found.

#### func (*Package) Func

```go
func (p *Package) Func(name string) (f *Function)
```
Func returns the package-level function of the specified name, or nil if not
found.

#### func (*Package) SetDebugMode

```go
func (pkg *Package) SetDebugMode(debug bool)
```
SetDebugMode sets the debug mode for package pkg. If true, all its functions
will include full debug info. This greatly increases the size of the instruction
stream, and causes Functions to depend upon the ASTs, potentially keeping them
live in memory for longer.

#### func (*Package) String

```go
func (p *Package) String() string
```

#### func (*Package) Type

```go
func (p *Package) Type(name string) (t *Type)
```
Type returns the package-level type of the specified name, or nil if not found.

#### func (*Package) TypesWithMethodSets

```go
func (pkg *Package) TypesWithMethodSets() []types.Type
```
TypesWithMethodSets returns an unordered slice containing the set of all types
referenced within package pkg and not belonging to some other package, for which
a complete (non-empty) method set is required at run-time.

A type belongs to a package if it is a named type or a pointer to a named type,
and the name was defined in that package. All other types belong to no package.

A type may appear in the TypesWithMethodSets() set of multiple distinct packages
if that type belongs to no package. Typical compilers emit method sets for such
types multiple times (using weak symbols) into each package that references
them, with the linker performing duplicate elimination.

This set includes the types of all operands of some MakeInterface instruction,
the types of all exported members of some package, and all types that are
subcomponents, since even types that aren't used directly may be derived via
reflection.

Callers must not mutate the result.

#### func (*Package) Var

```go
func (p *Package) Var(name string) (g *Global)
```
Var returns the package-level variable of the specified name, or nil if not
found.

#### func (*Package) WriteTo

```go
func (p *Package) WriteTo(w io.Writer) (int64, error)
```

#### type Panic

```go
type Panic struct {
	X Value // an interface{}
}
```

The Panic instruction initiates a panic with value X.

A Panic instruction must be the last instruction of its containing BasicBlock,
which must have no successors.

NB: 'go panic(x)' and 'defer panic(x)' do not use this instruction; they are
treated as calls to a built-in function.

Pos() returns the ast.CallExpr.Lparen if this panic was explicit in the source.

Example printed form:

    panic t0

#### func (*Panic) Block

```go
func (v *Panic) Block() *BasicBlock
```

#### func (*Panic) Operands

```go
func (s *Panic) Operands(rands []*Value) []*Value
```

#### func (*Panic) Parent

```go
func (v *Panic) Parent() *Function
```

#### func (*Panic) Pos

```go
func (s *Panic) Pos() token.Pos
```

#### func (*Panic) String

```go
func (s *Panic) String() string
```

#### type Parameter

```go
type Parameter struct {
}
```

A Parameter represents an input parameter of a function.

#### func (*Parameter) Name

```go
func (v *Parameter) Name() string
```

#### func (*Parameter) Object

```go
func (v *Parameter) Object() types.Object
```

#### func (*Parameter) Parent

```go
func (v *Parameter) Parent() *Function
```

#### func (*Parameter) Pos

```go
func (v *Parameter) Pos() token.Pos
```

#### func (*Parameter) Referrers

```go
func (v *Parameter) Referrers() *[]Instruction
```

#### func (*Parameter) String

```go
func (v *Parameter) String() string
```

#### func (*Parameter) Type

```go
func (v *Parameter) Type() types.Type
```

#### type Phi

```go
type Phi struct {
	Comment string  // a hint as to its purpose
	Edges   []Value // Edges[i] is value for Block().Preds[i]
}
```

The Phi instruction represents an SSA φ-node, which combines values that differ
across incoming control-flow edges and yields a new value. Within a block, all
φ-nodes must appear before all non-φ nodes.

Pos() returns the position of the && or || for short-circuit control-flow joins,
or that of the *Alloc for φ-nodes inserted during SSA renaming.

Example printed form:

    t2 = phi [0.start: t0, 1.if.then: t1, ...]

#### func (*Phi) Name

```go
func (v *Phi) Name() string
```

#### func (*Phi) Operands

```go
func (v *Phi) Operands(rands []*Value) []*Value
```

#### func (*Phi) Pos

```go
func (v *Phi) Pos() token.Pos
```

#### func (*Phi) Referrers

```go
func (v *Phi) Referrers() *[]Instruction
```

#### func (*Phi) String

```go
func (v *Phi) String() string
```

#### func (*Phi) Type

```go
func (v *Phi) Type() types.Type
```

#### type Program

```go
type Program struct {
	Fset *token.FileSet // position information for the files of this Program

	MethodSets types.MethodSetCache // cache of type-checker's method-sets
}
```

A Program is a partial or complete Go program converted to SSA form.

#### func  Create

```go
func Create(iprog *loader.Program, mode BuilderMode) *Program
```
Create returns a new SSA Program. An SSA Package is created for each
transitively error-free package of iprog.

Code for bodies of functions is not built until Build() is called on the result.

mode controls diagnostics and checking during SSA construction.

#### func (*Program) AllPackages

```go
func (prog *Program) AllPackages() []*Package
```
AllPackages returns a new slice containing all packages in the program prog in
unspecified order.

#### func (*Program) BuildAll

```go
func (prog *Program) BuildAll()
```
BuildAll calls Package.Build() for each package in prog. Building occurs in
parallel unless the BuildSerially mode flag was set.

BuildAll is idempotent and thread-safe.

#### func (*Program) ConstValue

```go
func (prog *Program) ConstValue(obj *types.Const) *Const
```
ConstValue returns the SSA Value denoted by the source-level named constant obj.
The result may be a *Const, or nil if not found.

#### func (*Program) CreatePackage

```go
func (prog *Program) CreatePackage(info *loader.PackageInfo) *Package
```
CreatePackage constructs and returns an SSA Package from an error-free package
described by info, and populates its Members mapping.

Repeated calls with the same info return the same Package.

The real work of building SSA form for each function is not done until a
subsequent call to Package.Build().

#### func (*Program) CreateTestMainPackage

```go
func (prog *Program) CreateTestMainPackage(pkgs ...*Package) *Package
```
CreateTestMainPackage creates and returns a synthetic "main" package that runs
all the tests of the supplied packages, similar to the one that would be created
by the 'go test' tool.

It returns nil if the program contains no tests.

#### func (*Program) FuncValue

```go
func (prog *Program) FuncValue(obj *types.Func) *Function
```
FuncValue returns the Function denoted by the source-level named function obj.

#### func (*Program) ImportedPackage

```go
func (prog *Program) ImportedPackage(path string) *Package
```
ImportedPackage returns the importable SSA Package whose import path is path, or
nil if no such SSA package has been created.

Not all packages are importable. For example, no import declaration can resolve
to the x_test package created by 'go test' or the ad-hoc main package created
'go build foo.go'.

#### func (*Program) LookupMethod

```go
func (prog *Program) LookupMethod(T types.Type, pkg *types.Package, name string) *Function
```
LookupMethod returns the implementation of the method of type T identified by
(pkg, name). It panics if there is no such method.

#### func (*Program) Method

```go
func (prog *Program) Method(meth *types.Selection) *Function
```
Method returns the Function implementing method meth, building wrapper methods
on demand.

Thread-safe.

EXCLUSIVE_LOCKS_ACQUIRED(prog.methodsMu)

#### func (*Program) Package

```go
func (prog *Program) Package(obj *types.Package) *Package
```
Package returns the SSA Package corresponding to the specified type-checker
package object. It returns nil if no such SSA package has been created.

#### func (*Program) TypesWithMethodSets

```go
func (prog *Program) TypesWithMethodSets() []types.Type
```
TypesWithMethodSets returns a new unordered slice containing all types in the
program for which a complete (non-empty) method set is required at run-time.

It is the union of pkg.TypesWithMethodSets() for all pkg in prog.AllPackages().

Thread-safe.

EXCLUSIVE_LOCKS_ACQUIRED(prog.methodsMu)

#### func (*Program) VarValue

```go
func (prog *Program) VarValue(obj *types.Var, pkg *Package, ref []ast.Node) (value Value, isAddr bool)
```
VarValue returns the SSA Value that corresponds to a specific identifier
denoting the source-level named variable obj.

VarValue returns nil if a local variable was not found, perhaps because its
package was not built, the debug information was not requested during SSA
construction, or the value was optimized away.

ref is the path to an ast.Ident (e.g. from PathEnclosingInterval), and that
ident must resolve to obj.

pkg is the package enclosing the reference. (A reference to a var always occurs
within a function, so we need to know where to find it.)

The Value of a defining (as opposed to referring) identifier is the value
assigned to it in its definition. Similarly, the Value of an identifier that is
the LHS of an assignment is the value assigned to it in that statement. In all
these examples, VarValue(x) returns the value of x and isAddr==false.

    var x X
    var x = X{}
    x := X{}
    x = X{}

When an identifier appears in an lvalue context other than as the LHS of an
assignment, the resulting Value is the var's address, not its value. This
situation is reported by isAddr, the second component of the result. In these
examples, VarValue(x) returns the address of x and isAddr==true.

    x.y = 0
    x[0] = 0
    _ = x[:]      (where x is an array)
    _ = &x
    x.method()    (iff method is on &x)

#### type Range

```go
type Range struct {
	X Value // string or map
}
```

The Range instruction yields an iterator over the domain and range of X, which
must be a string or map.

Elements are accessed via Next.

Type() returns an opaque and degenerate "rangeIter" type.

Pos() returns the ast.RangeStmt.For.

Example printed form:

    t0 = range "hello":string

#### func (*Range) Name

```go
func (v *Range) Name() string
```

#### func (*Range) Operands

```go
func (v *Range) Operands(rands []*Value) []*Value
```

#### func (*Range) Pos

```go
func (v *Range) Pos() token.Pos
```

#### func (*Range) Referrers

```go
func (v *Range) Referrers() *[]Instruction
```

#### func (*Range) String

```go
func (v *Range) String() string
```

#### func (*Range) Type

```go
func (v *Range) Type() types.Type
```

#### type Return

```go
type Return struct {
	Results []Value
}
```

The Return instruction returns values and control back to the calling function.

len(Results) is always equal to the number of results in the function's
signature.

If len(Results) > 1, Return returns a tuple value with the specified components
which the caller must access using Extract instructions.

There is no instruction to return a ready-made tuple like those returned by a
"value,ok"-mode TypeAssert, Lookup or UnOp(ARROW) or a tail-call to a function
with multiple result parameters.

Return must be the last instruction of its containing BasicBlock. Such a block
has no successors.

Pos() returns the ast.ReturnStmt.Return, if explicit in the source.

Example printed form:

    return
    return nil:I, 2:int

#### func (*Return) Block

```go
func (v *Return) Block() *BasicBlock
```

#### func (*Return) Operands

```go
func (s *Return) Operands(rands []*Value) []*Value
```

#### func (*Return) Parent

```go
func (v *Return) Parent() *Function
```

#### func (*Return) Pos

```go
func (s *Return) Pos() token.Pos
```

#### func (*Return) String

```go
func (s *Return) String() string
```

#### type RunDefers

```go
type RunDefers struct {
}
```

The RunDefers instruction pops and invokes the entire stack of procedure calls
pushed by Defer instructions in this function.

It is legal to encounter multiple 'rundefers' instructions in a single
control-flow path through a function; this is useful in the combined init()
function, for example.

Pos() returns NoPos.

Example printed form:

    rundefers

#### func (*RunDefers) Block

```go
func (v *RunDefers) Block() *BasicBlock
```

#### func (*RunDefers) Operands

```go
func (*RunDefers) Operands(rands []*Value) []*Value
```

#### func (*RunDefers) Parent

```go
func (v *RunDefers) Parent() *Function
```

#### func (*RunDefers) Pos

```go
func (s *RunDefers) Pos() token.Pos
```

#### func (*RunDefers) String

```go
func (*RunDefers) String() string
```

#### type Select

```go
type Select struct {
	States   []*SelectState
	Blocking bool
}
```

The Select instruction tests whether (or blocks until) one of the specified sent
or received states is entered.

Let n be the number of States for which Dir==RECV and T_i (0<=i<n) be the
element type of each such state's Chan. Select returns an n+2-tuple

    (index int, recvOk bool, r_0 T_0, ... r_n-1 T_n-1)

The tuple's components, described below, must be accessed via the Extract
instruction.

If Blocking, select waits until exactly one state holds, i.e. a channel becomes
ready for the designated operation of sending or receiving; select chooses one
among the ready states pseudorandomly, performs the send or receive operation,
and sets 'index' to the index of the chosen channel.

If !Blocking, select doesn't block if no states hold; instead it returns
immediately with index equal to -1.

If the chosen channel was used for a receive, the r_i component is set to the
received value, where i is the index of that state among all n receive states;
otherwise r_i has the zero value of type T_i. Note that the receive index i is
not the same as the state index index.

The second component of the triple, recvOk, is a boolean whose value is true iff
the selected operation was a receive and the receive successfully yielded a
value.

Pos() returns the ast.SelectStmt.Select.

Example printed form:

    t3 = select nonblocking [<-t0, t1<-t2]
    t4 = select blocking []

#### func (*Select) Name

```go
func (v *Select) Name() string
```

#### func (*Select) Operands

```go
func (v *Select) Operands(rands []*Value) []*Value
```

#### func (*Select) Pos

```go
func (v *Select) Pos() token.Pos
```

#### func (*Select) Referrers

```go
func (v *Select) Referrers() *[]Instruction
```

#### func (*Select) String

```go
func (s *Select) String() string
```

#### func (*Select) Type

```go
func (v *Select) Type() types.Type
```

#### type SelectState

```go
type SelectState struct {
	Dir       types.ChanDir // direction of case (SendOnly or RecvOnly)
	Chan      Value         // channel to use (for send or receive)
	Send      Value         // value to send (for send)
	Pos       token.Pos     // position of token.ARROW
	DebugNode ast.Node      // ast.SendStmt or ast.UnaryExpr(<-) [debug mode]
}
```

SelectState is a helper for Select. It represents one goal state and its
corresponding communication.

#### type Send

```go
type Send struct {
	Chan, X Value
}
```

The Send instruction sends X on channel Chan.

Pos() returns the ast.SendStmt.Arrow, if explicit in the source.

Example printed form:

    send t0 <- t1

#### func (*Send) Block

```go
func (v *Send) Block() *BasicBlock
```

#### func (*Send) Operands

```go
func (s *Send) Operands(rands []*Value) []*Value
```

#### func (*Send) Parent

```go
func (v *Send) Parent() *Function
```

#### func (*Send) Pos

```go
func (s *Send) Pos() token.Pos
```

#### func (*Send) String

```go
func (s *Send) String() string
```

#### type Slice

```go
type Slice struct {
	X              Value // slice, string, or *array
	Low, High, Max Value // each may be nil
}
```

The Slice instruction yields a slice of an existing string, slice or *array X
between optional integer bounds Low and High.

Dynamically, this instruction panics if X evaluates to a nil *array pointer.

Type() returns string if the type of X was string, otherwise a *types.Slice with
the same element type as X.

Pos() returns the ast.SliceExpr.Lbrack if created by a x[:] slice operation, the
ast.CompositeLit.Lbrace if created by a literal, or NoPos if not explicit in the
source (e.g. a variadic argument slice).

Example printed form:

    t1 = slice t0[1:]

#### func (*Slice) Name

```go
func (v *Slice) Name() string
```

#### func (*Slice) Operands

```go
func (v *Slice) Operands(rands []*Value) []*Value
```

#### func (*Slice) Pos

```go
func (v *Slice) Pos() token.Pos
```

#### func (*Slice) Referrers

```go
func (v *Slice) Referrers() *[]Instruction
```

#### func (*Slice) String

```go
func (v *Slice) String() string
```

#### func (*Slice) Type

```go
func (v *Slice) Type() types.Type
```

#### type Store

```go
type Store struct {
	Addr Value
	Val  Value
}
```

The Store instruction stores Val at address Addr. Stores can be of arbitrary
types.

Pos() returns the ast.StarExpr.Star, if explicit in the source.

Example printed form:

    *x = y

#### func (*Store) Block

```go
func (v *Store) Block() *BasicBlock
```

#### func (*Store) Operands

```go
func (s *Store) Operands(rands []*Value) []*Value
```

#### func (*Store) Parent

```go
func (v *Store) Parent() *Function
```

#### func (*Store) Pos

```go
func (s *Store) Pos() token.Pos
```

#### func (*Store) String

```go
func (s *Store) String() string
```

#### type Type

```go
type Type struct {
}
```

A Type is a Member of a Package representing a package-level named type.

Type() returns a *types.Named.

#### func (*Type) Name

```go
func (t *Type) Name() string
```

#### func (*Type) Object

```go
func (t *Type) Object() types.Object
```

#### func (*Type) Package

```go
func (t *Type) Package() *Package
```

#### func (*Type) Pos

```go
func (t *Type) Pos() token.Pos
```

#### func (*Type) RelString

```go
func (t *Type) RelString(from *types.Package) string
```

#### func (*Type) String

```go
func (t *Type) String() string
```

#### func (*Type) Token

```go
func (t *Type) Token() token.Token
```

#### func (*Type) Type

```go
func (t *Type) Type() types.Type
```

#### type TypeAssert

```go
type TypeAssert struct {
	X            Value
	AssertedType types.Type
	CommaOk      bool
}
```

The TypeAssert instruction tests whether interface value X has type
AssertedType.

If !CommaOk, on success it returns v, the result of the conversion (defined
below); on failure it panics.

If CommaOk: on success it returns a pair (v, true) where v is the result of the
conversion; on failure it returns (z, false) where z is AssertedType's zero
value. The components of the pair must be accessed using the Extract
instruction.

If AssertedType is a concrete type, TypeAssert checks whether the dynamic type
in interface X is equal to it, and if so, the result of the conversion is a copy
of the value in the interface.

If AssertedType is an interface, TypeAssert checks whether the dynamic type of
the interface is assignable to it, and if so, the result of the conversion is a
copy of the interface value X. If AssertedType is a superinterface of X.Type(),
the operation will fail iff the operand is nil. (Contrast with ChangeInterface,
which performs no nil-check.)

Type() reflects the actual type of the result, possibly a 2-types.Tuple;
AssertedType is the asserted type.

Pos() returns the ast.CallExpr.Lparen if the instruction arose from an explicit
T(e) conversion; the ast.TypeAssertExpr.Lparen if the instruction arose from an
explicit e.(T) operation; or the ast.CaseClause.Case if the instruction arose
from a case of a type-switch statement.

Example printed form:

    t1 = typeassert t0.(int)
    t3 = typeassert,ok t2.(T)

#### func (*TypeAssert) Name

```go
func (v *TypeAssert) Name() string
```

#### func (*TypeAssert) Operands

```go
func (v *TypeAssert) Operands(rands []*Value) []*Value
```

#### func (*TypeAssert) Pos

```go
func (v *TypeAssert) Pos() token.Pos
```

#### func (*TypeAssert) Referrers

```go
func (v *TypeAssert) Referrers() *[]Instruction
```

#### func (*TypeAssert) String

```go
func (v *TypeAssert) String() string
```

#### func (*TypeAssert) Type

```go
func (v *TypeAssert) Type() types.Type
```

#### type UnOp

```go
type UnOp struct {
	Op      token.Token // One of: NOT SUB ARROW MUL XOR ! - <- * ^
	X       Value
	CommaOk bool
}
```

The UnOp instruction yields the result of Op X. ARROW is channel receive. MUL is
pointer indirection (load). XOR is bitwise complement. SUB is negation. NOT is
logical negation.

If CommaOk and Op=ARROW, the result is a 2-tuple of the value above and a
boolean indicating the success of the receive. The components of the tuple are
accessed using Extract.

Pos() returns the ast.UnaryExpr.OpPos or ast.RangeStmt.TokPos (for ranging over
a channel), if explicit in the source.

Example printed form:

    t0 = *x
    t2 = <-t1,ok

#### func (*UnOp) Name

```go
func (v *UnOp) Name() string
```

#### func (*UnOp) Operands

```go
func (v *UnOp) Operands(rands []*Value) []*Value
```

#### func (*UnOp) Pos

```go
func (v *UnOp) Pos() token.Pos
```

#### func (*UnOp) Referrers

```go
func (v *UnOp) Referrers() *[]Instruction
```

#### func (*UnOp) String

```go
func (v *UnOp) String() string
```

#### func (*UnOp) Type

```go
func (v *UnOp) Type() types.Type
```

#### type Value

```go
type Value interface {
	// Name returns the name of this value, and determines how
	// this Value appears when used as an operand of an
	// Instruction.
	//
	// This is the same as the source name for Parameters,
	// Builtins, Functions, Captures, Globals.
	// For constants, it is a representation of the constant's value
	// and type.  For all other Values this is the name of the
	// virtual register defined by the instruction.
	//
	// The name of an SSA Value is not semantically significant,
	// and may not even be unique within a function.
	Name() string

	// If this value is an Instruction, String returns its
	// disassembled form; otherwise it returns unspecified
	// human-readable information about the Value, such as its
	// kind, name and type.
	String() string

	// Type returns the type of this value.  Many instructions
	// (e.g. IndexAddr) change their behaviour depending on the
	// types of their operands.
	Type() types.Type

	// Referrers returns the list of instructions that have this
	// value as one of their operands; it may contain duplicates
	// if an instruction has a repeated operand.
	//
	// Referrers actually returns a pointer through which the
	// caller may perform mutations to the object's state.
	//
	// Referrers is currently only defined for the function-local
	// values Capture, Parameter, Functions (iff anonymous) and
	// all value-defining instructions.
	// It returns nil for named Functions, Builtin, Const and Global.
	//
	// Instruction.Operands contains the inverse of this relation.
	Referrers() *[]Instruction

	// Pos returns the location of the AST token most closely
	// associated with the operation that gave rise to this value,
	// or token.NoPos if it was not explicit in the source.
	//
	// For each ast.Node type, a particular token is designated as
	// the closest location for the expression, e.g. the Lparen
	// for an *ast.CallExpr.  This permits a compact but
	// approximate mapping from Values to source positions for use
	// in diagnostic messages, for example.
	//
	// (Do not use this position to determine which Value
	// corresponds to an ast.Expr; use Function.ValueForExpr
	// instead.  NB: it requires that the function was built with
	// debug information.)
	//
	Pos() token.Pos
}
```

A Value is an SSA value that can be referenced by an instruction.
