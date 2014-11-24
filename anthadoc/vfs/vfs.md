# vfs
--
    import "."

Package vfs defines types for abstract file system access and provides an
implementation accessing the file system of the underlying OS.

## Usage

#### func  ReadFile

```go
func ReadFile(fs Opener, path string) ([]byte, error)
```
ReadFile reads the file named by path from fs and returns the contents.

#### type BindMode

```go
type BindMode int
```


```go
const (
	BindReplace BindMode = iota
	BindBefore
	BindAfter
)
```

#### type FileSystem

```go
type FileSystem interface {
	Opener
	Lstat(path string) (os.FileInfo, error)
	Stat(path string) (os.FileInfo, error)
	ReadDir(path string) ([]os.FileInfo, error)
	String() string
}
```

The FileSystem interface specifies the methods godoc is using to access the file
system for which it serves documentation.

#### func  OS

```go
func OS(root string) FileSystem
```
OS returns an implementation of FileSystem reading from the tree rooted at root.
Recording a root is convenient everywhere but necessary on Windows, because the
slash-separated path passed to Open has no way to specify a drive letter. Using
a root lets code refer to OS(`c:\`), OS(`d:\`) and so on.

#### type NameSpace

```go
type NameSpace map[string][]mountedFS
```

A NameSpace is a file system made up of other file systems mounted at specific
locations in the name space.

The representation is a map from mount point locations to the list of file
systems mounted at that location. A traditional Unix mount table would use a
single file system per mount point, but we want to be able to mount multiple
file systems on a single mount point and have the system behave as if the union
of those file systems were present at the mount point. For example, if the OS
file system has a Go installation in c:\Go and additional Go path trees in
d:\Work1 and d:\Work2, then this name space creates the view we want for the
godoc server:

    NameSpace{
    	"/": {
    		{old: "/", fs: OS(`c:\Go`), new: "/"},
    	},
    	"/src/pkg": {
    		{old: "/src/pkg", fs: OS(`c:\Go`), new: "/src/pkg"},
    		{old: "/src/pkg", fs: OS(`d:\Work1`), new: "/src"},
    		{old: "/src/pkg", fs: OS(`d:\Work2`), new: "/src"},
    	},
    }

This is created by executing:

    ns := NameSpace{}
    ns.Bind("/", OS(`c:\Go`), "/", BindReplace)
    ns.Bind("/src/pkg", OS(`d:\Work1`), "/src", BindAfter)
    ns.Bind("/src/pkg", OS(`d:\Work2`), "/src", BindAfter)

A particular mount point entry is a triple (old, fs, new), meaning that to
operate on a path beginning with old, replace that prefix (old) with new and
then pass that path to the FileSystem implementation fs.

Given this name space, a ReadDir of /src/pkg/code will check each prefix of the
path for a mount point (first /src/pkg/code, then /src/pkg, then /src, then /),
stopping when it finds one. For the above example, /src/pkg/code will find the
mount point at /src/pkg:

    {old: "/src/pkg", fs: OS(`c:\Go`), new: "/src/pkg"},
    {old: "/src/pkg", fs: OS(`d:\Work1`), new: "/src"},
    {old: "/src/pkg", fs: OS(`d:\Work2`), new: "/src"},

ReadDir will when execute these three calls and merge the results:

    OS(`c:\Go`).ReadDir("/src/pkg/code")
    OS(`d:\Work1').ReadDir("/src/code")
    OS(`d:\Work2').ReadDir("/src/code")

Note that the "/src/pkg" in "/src/pkg/code" has been replaced by just "/src" in
the final two calls.

OS is itself an implementation of a file system: it implements
OS(`c:\Go`).ReadDir("/src/pkg/code") as ioutil.ReadDir(`c:\Go\src\pkg\code`).

Because the new path is evaluated by fs (here OS(root)), another way to read the
mount table is to mentally combine fs+new, so that this table:

    {old: "/src/pkg", fs: OS(`c:\Go`), new: "/src/pkg"},
    {old: "/src/pkg", fs: OS(`d:\Work1`), new: "/src"},
    {old: "/src/pkg", fs: OS(`d:\Work2`), new: "/src"},

reads as:

    "/src/pkg" -> c:\Go\src\pkg
    "/src/pkg" -> d:\Work1\src
    "/src/pkg" -> d:\Work2\src

An invariant (a redundancy) of the name space representation is that
ns[mtpt][i].old is always equal to mtpt (in the example, ns["/src/pkg"]'s mount
table entries always have old == "/src/pkg"). The 'old' field is useful to
callers, because they receive just a []mountedFS and not any other indication of
which mount point was found.

#### func (NameSpace) Bind

```go
func (ns NameSpace) Bind(old string, newfs FileSystem, new string, mode BindMode)
```
Bind causes references to old to redirect to the path new in newfs. If mode is
BindReplace, old redirections are discarded. If mode is BindBefore, this
redirection takes priority over existing ones, but earlier ones are still
consulted for paths that do not exist in newfs. If mode is BindAfter, this
redirection happens only after existing ones have been tried and failed.

#### func (NameSpace) Fprint

```go
func (ns NameSpace) Fprint(w io.Writer)
```
Fprint writes a text representation of the name space to w.

#### func (NameSpace) Lstat

```go
func (ns NameSpace) Lstat(path string) (os.FileInfo, error)
```

#### func (NameSpace) Open

```go
func (ns NameSpace) Open(path string) (ReadSeekCloser, error)
```
Open implements the FileSystem Open method.

#### func (NameSpace) ReadDir

```go
func (ns NameSpace) ReadDir(path string) ([]os.FileInfo, error)
```
ReadDir implements the FileSystem ReadDir method. It's where most of the magic
is. (The rest is in resolve.)

Logically, ReadDir must return the union of all the directories that are named
by path. In order to avoid misinterpreting Go packages, of all the directories
that contain Go source code, we only include the files from the first, but we
include subdirectories from all.

ReadDir must also return directory entries needed to reach mount points. If the
name space looks like the example in the type NameSpace comment, but c:\Go does
not have a src/pkg subdirectory, we still want to be able to find that
subdirectory, because we've mounted d:\Work1 and d:\Work2 there. So if we don't
see "src" in the directory listing for c:\Go, we add an entry for it before
returning.

#### func (NameSpace) Stat

```go
func (ns NameSpace) Stat(path string) (os.FileInfo, error)
```

#### func (NameSpace) String

```go
func (NameSpace) String() string
```

#### type Opener

```go
type Opener interface {
	Open(name string) (ReadSeekCloser, error)
}
```

Opener is a minimal virtual filesystem that can only open regular files.

#### type ReadSeekCloser

```go
type ReadSeekCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}
```

A ReadSeekCloser can Read, Seek, and Close.
