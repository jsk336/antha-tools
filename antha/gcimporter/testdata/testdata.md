# exports
--
    import "."


## Usage

```go
const (
	C0 int = 0
	C1     = 3.14159265
	C2     = 2.718281828i
	C3     = -123.456e-789
	C4     = +123.456E+789
	C5     = 1234i
	C6     = "foo\n"
	C7     = `bar\n`
)
```

```go
var (
	V0 int
	V1 = -991.0
)
```

#### func  F1

```go
func F1()
```

#### func  F2

```go
func F2(x int)
```

#### func  F3

```go
func F3() int
```

#### func  F4

```go
func F4() float32
```

#### func  F5

```go
func F5(a, b, c int, u, v, w struct{ x, y T1 }, more ...interface{}) (p, q, r chan<- T10)
```

#### type T1

```go
type T1 int
```


#### func (*T1) M1

```go
func (p *T1) M1()
```

#### type T10

```go
type T10 struct {
	T8
	T9
}
```


#### type T11

```go
type T11 map[int]string
```


#### type T12

```go
type T12 interface{}
```


#### type T13

```go
type T13 interface {
	// contains filtered or unexported methods
}
```


#### type T14

```go
type T14 interface {
	T12
	T13
	// contains filtered or unexported methods
}
```


#### type T15

```go
type T15 func()
```


#### type T16

```go
type T16 func(int)
```


#### type T17

```go
type T17 func(x int)
```


#### type T18

```go
type T18 func() float32
```


#### type T19

```go
type T19 func() (x float32)
```


#### type T2

```go
type T2 [10]int
```


#### type T20

```go
type T20 func(...interface{})
```


#### type T21

```go
type T21 struct {
}
```


#### type T22

```go
type T22 struct {
}
```


#### type T23

```go
type T23 struct {
}
```


#### type T24

```go
type T24 *T24
```


#### type T25

```go
type T25 *T26
```


#### type T26

```go
type T26 *T27
```


#### type T27

```go
type T27 *T25
```


#### type T28

```go
type T28 func(T28) T28
```


#### type T3

```go
type T3 []int
```


#### type T4

```go
type T4 *int
```


#### type T5

```go
type T5 chan int
```


#### type T6a

```go
type T6a chan<- int
```


#### type T6b

```go
type T6b chan (<-chan int)
```


#### type T6c

```go
type T6c chan<- (chan int)
```


#### type T7

```go
type T7 <-chan *ast.File
```


#### type T8

```go
type T8 struct{}
```


#### type T9

```go
type T9 struct {
}
```
