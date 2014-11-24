# interp
--
    import "."

Package ssa/interp defines an interpreter for the SSA representation of Go
programs.

This interpreter is provided as an adjunct for testing the SSA construction
algorithm. Its purpose is to provide a minimal metacircular implementation of
the dynamic semantics of each SSA instruction. It is not, and will never be, a
production-quality Go interpreter.

The following is a partial list of Go features that are currently unsupported or
incomplete in the interpreter.

* Unsafe operations, including all uses of unsafe.Pointer, are impossible to
support given the "boxed" value representation we have chosen.

* The reflect package is only partially implemented.

* "sync/atomic" operations are not currently atomic due to the "boxed" value
representation: it is not possible to read, modify and write an interface value
atomically. As a consequence, Mutexes are currently broken. TODO(adonovan):
provide a metacircular implementation of Mutex avoiding the broken atomic
primitives.

* recover is only partially implemented. Also, the interpreter makes no attempt
to distinguish target panics from interpreter crashes.

* map iteration is asymptotically inefficient.

* the sizes of the int, uint and uintptr types in the target program are assumed
to be the same as those of the interpreter itself.

* all values occupy space, even those of types defined by the spec to have zero
size, e.g. struct{}. This can cause asymptotic performance degradation.

* os.Exit is implemented using panic, causing deferred functions to run.

## Usage

```go
var CapturedOutput *bytes.Buffer
```
If CapturedOutput is non-nil, all writes by the interpreted program to file
descriptors 1 and 2 will also be written to CapturedOutput.

(The $GOROOT/test system requires that the test be considered a failure if "BUG"
appears in the combined stdout/stderr output, even if it exits zero. This is a
global variable shared by all interpreters in the same process.)

#### func  Interpret

```go
func Interpret(mainpkg *ssa.Package, mode Mode, sizes types.Sizes, filename string, args []string) (exitCode int)
```
Interpret interprets the Go program whose main package is mainpkg. mode
specifies various interpreter options. filename and args are the initial values
of os.Args for the target program. sizes is the effective type-sizing function
for this program.

Interpret returns the exit code of the program: 2 for panic (like gc does), or
the argument to os.Exit for normal termination.

The SSA program must include the "runtime" package.

#### type Mode

```go
type Mode uint
```

Mode is a bitmask of options affecting the interpreter.

```go
const (
	DisableRecover Mode = 1 << iota // Disable recover() in target programs; show interpreter crash instead.
	EnableTracing                   // Print a trace of all instructions as they are interpreted.
)
```
