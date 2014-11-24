# util
--
    import "."

Package util contains utility types and functions for anthadoc.

## Usage

#### func  IsText

```go
func IsText(s []byte) bool
```
IsText reports whether a significant prefix of s looks like correct UTF-8; that
is, if it is likely that s is human-readable text.

#### func  IsTextFile

```go
func IsTextFile(fs vfs.Opener, filename string) bool
```
IsTextFile reports whether the file has a known extension indicating a text
file, or if a significant chunk of the specified file looks like correct UTF-8;
that is, if it is likely that the file contains human- readable text.

#### type RWValue

```go
type RWValue struct {
}
```

An RWValue wraps a value and permits mutually exclusive access to it and records
the time the value was last set.

#### func (*RWValue) Get

```go
func (v *RWValue) Get() (interface{}, time.Time)
```

#### func (*RWValue) Set

```go
func (v *RWValue) Set(value interface{})
```

#### type Throttle

```go
type Throttle struct {
}
```

A Throttle permits throttling of a goroutine by calling the Throttle method
repeatedly.

#### func  NewThrottle

```go
func NewThrottle(r float64, dt time.Duration) *Throttle
```
NewThrottle creates a new Throttle with a throttle value r and a minimum
allocated run time slice of dt:

    r == 0: "empty" throttle; the goroutine is always sleeping
    r == 1: full throttle; the goroutine is never sleeping

A value of r == 0.6 throttles a goroutine such that it runs approx. 60% of the
time, and sleeps approx. 40% of the time. Values of r < 0 or r > 1 are clamped
down to values between 0 and 1. Values of dt < 0 are set to 0.

#### func (*Throttle) Throttle

```go
func (p *Throttle) Throttle()
```
Throttle calls time.Sleep such that over time the ratio tr/ts between
accumulated run (tr) and sleep times (ts) approximates the value 1/(1-r) where r
is the throttle value. Throttle returns immediately (w/o sleeping) if less than
tm ns have passed since the last call to Throttle.
