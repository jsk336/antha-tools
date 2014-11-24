# callgraph
--
    import "."

Package callgraph defines the call graph and various algorithms and utilities to
operate on it.

A call graph is a labelled directed graph whose nodes represent functions and
whose edge labels represent syntactic function call sites. The presence of a
labelled edge (caller, site, callee) indicates that caller may call callee at
the specified call site.

A call graph is a multigraph: it may contain multiple edges (caller, *, callee)
connecting the same pair of nodes, so long as the edges differ by label; this
occurs when one function calls another function from multiple call sites. Also,
it may contain multiple edges (caller, site, *) that differ only by callee; this
indicates a polymorphic call.

A SOUND call graph is one that overapproximates the dynamic calling behaviors of
the program in all possible executions. One call graph is more PRECISE than
another if it is a smaller overapproximation of the dynamic behavior.

All call graphs have a synthetic root node which is responsible for calling
main() and init().

Calls to built-in functions (e.g. panic, println) are not represented in the
call graph; they are treated like built-in operators of the language.

## Usage

#### func  AddEdge

```go
func AddEdge(caller *Node, site ssa.CallInstruction, callee *Node)
```
AddEdge adds the edge (caller, site, callee) to the call graph. Elimination of
duplicate edges is the caller's responsibility.

#### func  CalleesOf

```go
func CalleesOf(caller *Node) map[*Node]bool
```
CalleesOf returns a new set containing all direct callees of the caller node.

#### func  GraphVisitEdges

```go
func GraphVisitEdges(g *Graph, edge func(*Edge) error) error
```
GraphVisitEdges visits all the edges in graph g in depth-first order. The edge
function is called for each edge in postorder. If it returns non-nil, visitation
stops and GraphVisitEdges returns that value.

#### func  PathSearch

```go
func PathSearch(start *Node, isEnd func(*Node) bool) []*Edge
```
PathSearch finds an arbitrary path starting at node start and ending at some
node for which isEnd() returns true. On success, PathSearch returns the path as
an ordered list of edges; on failure, it returns nil.

#### type Edge

```go
type Edge struct {
	Caller *Node
	Site   ssa.CallInstruction
	Callee *Node
}
```

A Edge represents an edge in the call graph.

Site is nil for edges originating in synthetic or intrinsic functions, e.g.
reflect.Call or the root of the call graph.

#### func (Edge) String

```go
func (e Edge) String() string
```

#### type Graph

```go
type Graph struct {
	Root  *Node                   // the distinguished root node
	Nodes map[*ssa.Function]*Node // all nodes by function
}
```

A Graph represents a call graph.

A graph may contain nodes that are not reachable from the root. If the call
graph is sound, such nodes indicate unreachable functions.

#### func  New

```go
func New(root *ssa.Function) *Graph
```
New returns a new Graph with the specified root node.

#### func (*Graph) CreateNode

```go
func (g *Graph) CreateNode(fn *ssa.Function) *Node
```
CreateNode returns the Node for fn, creating it if not present.

#### func (*Graph) DeleteNode

```go
func (g *Graph) DeleteNode(n *Node)
```
DeleteNode removes node n and its edges from the graph g. (NB: not efficient for
batch deletion.)

#### func (*Graph) DeleteSyntheticNodes

```go
func (g *Graph) DeleteSyntheticNodes()
```
DeleteSyntheticNodes removes from call graph g all nodes for synthetic functions
(except g.Root and package initializers), preserving the topology. In effect,
calls to synthetic wrappers are "inlined".

#### type Node

```go
type Node struct {
	Func *ssa.Function // the function this node represents
	ID   int           // 0-based sequence number
	In   []*Edge       // unordered set of incoming call edges (n.In[*].Callee == n)
	Out  []*Edge       // unordered set of outgoing call edges (n.Out[*].Caller == n)
}
```

A Node represents a node in a call graph.

#### func (*Node) String

```go
func (n *Node) String() string
```
