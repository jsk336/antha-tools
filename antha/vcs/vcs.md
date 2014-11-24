# vcs
--
    import "."


## Usage

```go
var ShowCmd bool
```
ShowCmd controls whether VCS commands are printed.

```go
var Verbose bool
```
Verbose enables verbose operation logging.

#### type Cmd

```go
type Cmd struct {
	Name string
	Cmd  string // name of binary to invoke command

	CreateCmd   string // command to download a fresh copy of a repository
	DownloadCmd string // command to download updates into an existing repository

	TagCmd         []TagCmd // commands to list tags
	TagLookupCmd   []TagCmd // commands to lookup tags before running tagSyncCmd
	TagSyncCmd     string   // command to sync to specific tag
	TagSyncDefault string   // command to sync to default tag

	LogCmd string // command to list repository changelogs in an XML format

	Scheme  []string
	PingCmd string
}
```

A Cmd describes how to use a version control system like Mercurial, Git, or
Subversion.

#### func  ByCmd

```go
func ByCmd(cmd string) *Cmd
```
ByCmd returns the version control system for the given command name (hg, git,
svn, bzr).

#### func  FromDir

```go
func FromDir(dir, srcRoot string) (vcs *Cmd, root string, err error)
```
FromDir inspects dir and its parents to determine the version control system and
code repository to use. On return, root is the import path corresponding to the
root of the repository (thus root is a prefix of importPath).

#### func (*Cmd) Create

```go
func (v *Cmd) Create(dir, repo string) error
```
Create creates a new copy of repo in dir. The parent of dir must exist; dir must
not.

#### func (*Cmd) CreateAtRev

```go
func (v *Cmd) CreateAtRev(dir, repo, rev string) error
```
CreateAtRev creates a new copy of repo in dir at revision rev. The parent of dir
must exist; dir must not. rev must be a valid revision in repo.

#### func (*Cmd) Download

```go
func (v *Cmd) Download(dir string) error
```
Download downloads any new changes for the repo in dir. dir must be a valid VCS
repo compatible with v.

#### func (*Cmd) Log

```go
func (v *Cmd) Log(dir, logTemplate string) ([]byte, error)
```
Log logs the changes for the repo in dir. dir must be a valid VCS repo
compatible with v.

#### func (*Cmd) LogAtRev

```go
func (v *Cmd) LogAtRev(dir, rev, logTemplate string) ([]byte, error)
```
LogAtRev logs the change for repo in dir at the rev revision. dir must be a
valid VCS repo compatible with v. rev must be a valid revision for the repo in
dir.

#### func (*Cmd) Ping

```go
func (v *Cmd) Ping(scheme, repo string) error
```
Ping pings the repo to determine if scheme used is valid. This repo must be
pingable with this scheme and VCS.

#### func (*Cmd) String

```go
func (v *Cmd) String() string
```

#### func (*Cmd) TagSync

```go
func (v *Cmd) TagSync(dir, tag string) error
```
TagSync syncs the repo in dir to the named tag, which either is a tag returned
by tags or is v.TagDefault. dir must be a valid VCS repo compatible with v and
the tag must exist.

#### func (*Cmd) Tags

```go
func (v *Cmd) Tags(dir string) ([]string, error)
```
Tags returns the list of available tags for the repo in dir. dir must be a valid
VCS repo compatible with v.

#### type RepoRoot

```go
type RepoRoot struct {
	VCS *Cmd

	// repo is the repository URL, including scheme
	Repo string

	// root is the import path corresponding to the root of the
	// repository
	Root string
}
```

RepoRoot represents a version control system, a repo, and a root of where to put
it on disk.

#### func  RepoRootForImportDynamic

```go
func RepoRootForImportDynamic(importPath string, verbose bool) (*RepoRoot, error)
```
RepoRootForImportDynamic finds a *RepoRoot for a custom domain that's not
statically known by RepoRootForImportPathStatic.

This handles "vanity import paths" like "name.tld/pkg/foo".

#### func  RepoRootForImportPath

```go
func RepoRootForImportPath(importPath string, verbose bool) (*RepoRoot, error)
```
RepoRootForImportPath analyzes importPath to determine the version control
system, and code repository to use.

#### func  RepoRootForImportPathStatic

```go
func RepoRootForImportPathStatic(importPath, scheme string) (*RepoRoot, error)
```
RepoRootForImportPathStatic attempts to map importPath to a RepoRoot using the
commonly-used VCS hosting sites in vcsPaths (github.com/user/dir), or from a
fully-qualified importPath already containing its VCS type
(foo.com/repo.git/dir)

If scheme is non-empty, that scheme is forced.

#### type TagCmd

```go
type TagCmd struct {
	Cmd     string // command to list tags
	Pattern string // regexp to extract tags from list
}
```

A TagCmd describes a command to list available tags that can be passed to
Cmd.TagSyncCmd.
