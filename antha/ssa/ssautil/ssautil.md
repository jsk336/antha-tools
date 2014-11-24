# ssautil
--
    import "."


## Usage

#### func  AllFunctions

```go
func AllFunctions(prog *ssa.Program) map[*ssa.Function]bool
```
AllFunctions finds and returns the set of functions potentially needed by
program prog, as determined by a simple linker-style reachability algorithm
starting from the members and method-sets of each package. The result may
include anonymous functions and synthetic wrappers.

Precondition: all packages are built.

#### func  Switches

```go
func Switches(fn *ssa.Function) []Switch
```
Switches examines the control-flow graph of fn and returns the set of inferred
value and type switches. A value switch tests an ssa.Value for equality against
two or more compile-time constant values. Switches involving link-time constants
(addresses) are ignored. A type switch type-asserts an ssa.Value against two or
more types.

The switches are returned in dominance order.

The resulting switches do not necessarily correspond to uses of the 'switch'
keyword in the source: for example, a single source-level switch statement with
non-constant cases may result in zero, one or many Switches, one per plural
sequence of constant cases. Switches may even be inferred from if/else- or
goto-based control flow. (In general, the control flow constructs of the source
program cannot be faithfully reproduced from the SSA representation.)

#### type ConstCase

```go
type ConstCase struct {
	Block *ssa.BasicBlock // block performing the comparison
	Body  *ssa.BasicBlock // body of the case
	Value *ssa.Const      // case comparand
}
```

A ConstCase represents a single constant comparison. It is part of a Switch.

#### type Switch

```go
type Switch struct {
	Start      *ssa.BasicBlock // block containing start of if/else chain
	X          ssa.Value       // the switch operand
	ConstCases []ConstCase     // ordered list of constant comparisons
	TypeCases  []TypeCase      // ordered list of type assertions
	Default    *ssa.BasicBlock // successor if all comparisons fail
}
```

A Switch is a logical high-level control flow operation (a multiway branch)
discovered by analysis of a CFG containing only if/else chains. It is not part
of the ssa.Instruction set.

One of ConstCases and TypeCases has length >= 2; the other is nil.

In a value switch, the list of cases may contain duplicate constants. A type
switch may contain duplicate types, or types assignable to an interface type
also in the list. TODO(adonovan): eliminate such duplicates.

#### func (*Switch) String

```go
func (sw *Switch) String() string
```

#### type TypeCase

```go
type TypeCase struct {
	Block   *ssa.BasicBlock // block performing the type assert
	Body    *ssa.BasicBlock // body of the case
	Type    types.Type      // case type
	Binding ssa.Value       // value bound by this case
}
```

A TypeCase represents a single type assertion. It is part of a Switch.
