# socket
--
    import "."

Package socket implements an WebSocket-based playground backend. Clients connect
to a websocket handler and send run/kill commands, and the server sends the
output and exit status of the running processes. Multiple clients running
multiple processes may be served concurrently. The wire format is JSON and is
described by the Message type.

This will not run on App Engine as WebSockets are not supported there.

## Usage

```go
var Environ func() []string = os.Environ
```
Environ provides an environment when a binary, such as the go tool, is invoked.

```go
var RunScripts = true
```
RunScripts specifies whether the socket handler should execute shell scripts
(snippets that start with a shebang).

#### func  NewHandler

```go
func NewHandler(origin *url.URL) websocket.Server
```
NewHandler returns a websocket server which checks the origin of requests.

#### type Message

```go
type Message struct {
	Id      string // client-provided unique id for the process
	Kind    string // in: "run", "kill" out: "stdout", "stderr", "end"
	Body    string
	Options *Options `json:",omitempty"`
}
```

Message is the wire format for the websocket connection to the browser. It is
used for both sending output messages and receiving commands, as distinguished
by the Kind field.

#### type Options

```go
type Options struct {
	Race bool // use -race flag when building code (for "run" only)
}
```

Options specify additional message options.
