# build
--
    import "."


## Usage

```go
const (
	PerfRunLength = 1024
)
```

#### func  AddCommitToPerfTodo

```go
func AddCommitToPerfTodo(c appengine.Context, com *Commit) error
```
AddCommitToPerfTodo adds the commit to all existing PerfTodo entities.

#### func  AuthHandler

```go
func AuthHandler(h dashHandler) http.HandlerFunc
```
AuthHandler wraps a http.HandlerFunc with a handler that validates the supplied
key and builder query parameters.

#### func  GetCommits

```go
func GetCommits(c appengine.Context, startCommitNum, n int) ([]*Commit, error)
```
GetCommits returns [startCommitNum, startCommitNum+n) commits. Commits
information is partial (obtained from CommitRun), do not store them back into
datastore.

#### func  GetPerfMetricsForCommits

```go
func GetPerfMetricsForCommits(c appengine.Context, builder, benchmark, metric string, startCommitNum, n int) ([]uint64, error)
```
GetPerfMetricsForCommits returns perf metrics for builder/benchmark/metric and
commits [startCommitNum, startCommitNum+n).

#### func  Packages

```go
func Packages(c appengine.Context, kind string) ([]*Package, error)
```
Packages returns packages of the specified kind. Kind must be one of "external"
or "subrepo".

#### func  PerfConfigKey

```go
func PerfConfigKey(c appengine.Context) *datastore.Key
```

#### func  PutLog

```go
func PutLog(c appengine.Context, text string) (hash string, err error)
```

#### func  UpdatePerfConfig

```go
func UpdatePerfConfig(c appengine.Context, r *http.Request, req *PerfRequest) (newBenchmark bool, err error)
```
UpdatePerfConfig updates the PerfConfig entity with results of benchmarking.
Returns whether it's a benchmark that we have not yet seem on the builder.

#### type Commit

```go
type Commit struct {
	PackagePath string // (empty for main repo commits)
	Hash        string
	ParentHash  string
	Num         int // Internal monotonic counter unique to this package.

	User              string
	Desc              string `datastore:",noindex"`
	Time              time.Time
	NeedsBenchmarking bool
	TryPatch          bool

	// ResultData is the Data string of each build Result for this Commit.
	// For non-Go commits, only the Results for the current Go tip, weekly,
	// and release Tags are stored here. This is purely de-normalized data.
	// The complete data set is stored in Result entities.
	ResultData []string `datastore:",noindex"`

	// PerfResults holds a set of “builder|benchmark” tuples denoting
	// what benchmarks have been executed on the commit.
	PerfResults []string `datastore:",noindex"`

	FailNotificationSent bool
}
```

A Commit describes an individual commit in a package.

Each Commit entity is a descendant of its associated Package entity. In other
words, all Commits with the same PackagePath belong to the same datastore entity
group.

#### func (*Commit) AddPerfResult

```go
func (com *Commit) AddPerfResult(c appengine.Context, builder, benchmark string) error
```
AddPerfResult remembers that the builder has run the benchmark on the commit. It
must be called from inside a datastore transaction.

#### func (*Commit) AddResult

```go
func (com *Commit) AddResult(c appengine.Context, r *Result) error
```
AddResult adds the denormalized Result data to the Commit's Result field. It
must be called from inside a datastore transaction.

#### func (*Commit) Key

```go
func (com *Commit) Key(c appengine.Context) *datastore.Key
```

#### func (*Commit) Result

```go
func (c *Commit) Result(builder, goHash string) *Result
```
Result returns the build Result for this Commit for the given builder/goHash.

#### func (*Commit) ResultGoHashes

```go
func (c *Commit) ResultGoHashes() []string
```

#### func (*Commit) Results

```go
func (c *Commit) Results() (results []*Result)
```
Results returns the build Results for this Commit.

#### func (*Commit) Valid

```go
func (c *Commit) Valid() error
```

#### type CommitRun

```go
type CommitRun struct {
	PackagePath       string // (empty for main repo commits)
	StartCommitNum    int
	Hash              []string    `datastore:",noindex"`
	User              []string    `datastore:",noindex"`
	Desc              []string    `datastore:",noindex"` // Only first line.
	Time              []time.Time `datastore:",noindex"`
	NeedsBenchmarking []bool      `datastore:",noindex"`
}
```

A CommitRun provides summary information for commits [StartCommitNum,
StartCommitNum + PerfRunLength). Descendant of Package.

#### func  GetCommitRun

```go
func GetCommitRun(c appengine.Context, commitNum int) (*CommitRun, error)
```
GetCommitRun loads and returns CommitRun that contains information for commit
commitNum.

#### func (*CommitRun) AddCommit

```go
func (cr *CommitRun) AddCommit(c appengine.Context, com *Commit) error
```

#### func (*CommitRun) Key

```go
func (cr *CommitRun) Key(c appengine.Context) *datastore.Key
```

#### type Dashboard

```go
type Dashboard struct {
	Name     string     // This dashboard's name and namespace
	RelPath  string     // The relative url path
	Packages []*Package // The project's packages to build
}
```

Dashboard describes a unique build dashboard.

#### func (*Dashboard) Context

```go
func (d *Dashboard) Context(c appengine.Context) appengine.Context
```
Context returns a namespaced context for this dashboard, or panics if it fails
to create a new context.

#### type Log

```go
type Log struct {
	CompressedLog []byte
}
```

A Log is a gzip-compressed log file stored under the SHA1 hash of the
uncompressed log text.

#### func (*Log) Text

```go
func (l *Log) Text() ([]byte, error)
```

#### type Package

```go
type Package struct {
	Kind    string // "subrepo", "external", or empty for the main Go tree
	Name    string
	Path    string // (empty for the main Go tree)
	NextNum int    // Num of the next head Commit
}
```

A Package describes a package that is listed on the dashboard.

#### func  GetPackage

```go
func GetPackage(c appengine.Context, path string) (*Package, error)
```
GetPackage fetches a Package by path from the datastore.

#### func (*Package) Key

```go
func (p *Package) Key(c appengine.Context) *datastore.Key
```

#### func (*Package) LastCommit

```go
func (p *Package) LastCommit(c appengine.Context) (*Commit, error)
```
LastCommit returns the most recent Commit for this Package.

#### func (*Package) String

```go
func (p *Package) String() string
```

#### type PackageState

```go
type PackageState struct {
	Package *Package
	Commit  *Commit
}
```

PackageState represents the state of a Package at a Tag.

#### type Pagination

```go
type Pagination struct {
	Next, Prev int
	HasPrev    bool
}
```


#### type ParsedPerfResult

```go
type ParsedPerfResult struct {
	OK        bool
	Metrics   map[string]uint64
	Artifacts map[string]string
}
```


#### type PerfArtifact

```go
type PerfArtifact struct {
	Type string
	Body string
}
```


#### type PerfChange

```go
type PerfChange struct {
}
```


#### type PerfChangeBenchmark

```go
type PerfChangeBenchmark struct {
	Name    string
	Metrics []*PerfChangeMetric
}
```


#### type PerfChangeBenchmarkSlice

```go
type PerfChangeBenchmarkSlice []*PerfChangeBenchmark
```


#### func (PerfChangeBenchmarkSlice) Len

```go
func (l PerfChangeBenchmarkSlice) Len() int
```

#### func (PerfChangeBenchmarkSlice) Less

```go
func (l PerfChangeBenchmarkSlice) Less(i, j int) bool
```

#### func (PerfChangeBenchmarkSlice) Swap

```go
func (l PerfChangeBenchmarkSlice) Swap(i, j int)
```

#### type PerfChangeMetric

```go
type PerfChangeMetric struct {
	Name  string
	Old   uint64
	New   uint64
	Delta float64
}
```


#### type PerfChangeMetricSlice

```go
type PerfChangeMetricSlice []*PerfChangeMetric
```


#### func (PerfChangeMetricSlice) Len

```go
func (l PerfChangeMetricSlice) Len() int
```

#### func (PerfChangeMetricSlice) Less

```go
func (l PerfChangeMetricSlice) Less(i, j int) bool
```

#### func (PerfChangeMetricSlice) Swap

```go
func (l PerfChangeMetricSlice) Swap(i, j int)
```

#### type PerfConfig

```go
type PerfConfig struct {
	BuilderBench []string `datastore:",noindex"` // "builder|benchmark" pairs
	BuilderProcs []string `datastore:",noindex"` // "builder|proc" pairs
	BenchMetric  []string `datastore:",noindex"` // "benchmark|metric" pairs
	NoiseLevels  []string `datastore:",noindex"` // "builder|benchmark|metric1=noise1|metric2=noise2"
}
```

PerfConfig holds read-mostly configuration related to benchmarking. There is
only one PerfConfig entity.

#### func  GetPerfConfig

```go
func GetPerfConfig(c appengine.Context, r *http.Request) (*PerfConfig, error)
```

#### func (*PerfConfig) BenchmarkProcList

```go
func (pc *PerfConfig) BenchmarkProcList() (res []string)
```

#### func (*PerfConfig) BenchmarksForBuilder

```go
func (pc *PerfConfig) BenchmarksForBuilder(builder string) []string
```

#### func (*PerfConfig) BuildersForBenchmark

```go
func (pc *PerfConfig) BuildersForBenchmark(bench string) []string
```

#### func (*PerfConfig) MetricsForBenchmark

```go
func (pc *PerfConfig) MetricsForBenchmark(bench string) []string
```

#### func (*PerfConfig) NoiseLevel

```go
func (pc *PerfConfig) NoiseLevel(builder, benchmark, metric string) float64
```

#### func (*PerfConfig) ProcList

```go
func (pc *PerfConfig) ProcList(builder string) []int
```

#### type PerfMetric

```go
type PerfMetric struct {
	Type string
	Val  uint64
}
```


#### type PerfMetricRun

```go
type PerfMetricRun struct {
	PackagePath    string
	Builder        string
	Benchmark      string
	Metric         string // e.g. realtime, cputime, gc-pause
	StartCommitNum int
	Vals           []int64 `datastore:",noindex"`
}
```

A PerfMetricRun entity holds a set of metric values for builder/benchmark/metric
for commits [StartCommitNum, StartCommitNum + PerfRunLength). Descendant of
Package.

#### func  GetPerfMetricRun

```go
func GetPerfMetricRun(c appengine.Context, builder, benchmark, metric string, commitNum int) (*PerfMetricRun, error)
```
GetPerfMetricRun loads and returns PerfMetricRun that contains information for
commit commitNum.

#### func (*PerfMetricRun) AddMetric

```go
func (m *PerfMetricRun) AddMetric(c appengine.Context, commitNum int, v uint64) error
```

#### func (*PerfMetricRun) Key

```go
func (m *PerfMetricRun) Key(c appengine.Context) *datastore.Key
```

#### type PerfRequest

```go
type PerfRequest struct {
	Builder   string
	Benchmark string
	Hash      string
	OK        bool
	Metrics   []PerfMetric
	Artifacts []PerfArtifact
}
```

perf-result request payload

#### type PerfResult

```go
type PerfResult struct {
	PackagePath string
	CommitHash  string
	CommitNum   int
	Data        []string `datastore:",noindex"` // "builder|benchmark|ok|metric1=val1|metric2=val2|file:log=hash|file:cpuprof=hash"
}
```

A PerfResult describes all benchmarking result for a Commit. Descendant of
Package.

#### func (*PerfResult) AddResult

```go
func (r *PerfResult) AddResult(req *PerfRequest) bool
```
AddResult add the benchmarking result to r. Existing result for the same
builder/benchmark is replaced if already exists. Returns whether the result was
already present.

#### func (*PerfResult) Key

```go
func (r *PerfResult) Key(c appengine.Context) *datastore.Key
```

#### func (*PerfResult) ParseData

```go
func (r *PerfResult) ParseData() map[string]map[string]*ParsedPerfResult
```

#### type PerfResultCache

```go
type PerfResultCache struct {
}
```

PerfResultCache caches a set of PerfResults so that it's easy to access them
without lots of duplicate accesses to datastore. It allows to iterate over newer
or older results for some base commit.

#### func  MakePerfResultCache

```go
func MakePerfResultCache(c appengine.Context, com *Commit, newer bool) *PerfResultCache
```

#### func (*PerfResultCache) Get

```go
func (rc *PerfResultCache) Get(commitNum int) *PerfResult
```

#### func (*PerfResultCache) Next

```go
func (rc *PerfResultCache) Next(commitNum int) (*PerfResult, error)
```
Next returns the next PerfResult for the commit commitNum. It does not care
whether the result has any data, failed or whatever.

#### func (*PerfResultCache) NextForComparison

```go
func (rc *PerfResultCache) NextForComparison(commitNum int, builder string) (*PerfResult, error)
```
NextForComparison returns PerfResult which we need to use for performance
comprison. It skips failed results, but does not skip results with no data.

#### type PerfTodo

```go
type PerfTodo struct {
	PackagePath string // (empty for main repo commits)
	Builder     string
	CommitNums  []int `datastore:",noindex"` // LIFO queue of commits to benchmark.
}
```

A PerfTodo contains outstanding commits for benchmarking for a builder.
Descendant of Package.

#### func (*PerfTodo) Key

```go
func (todo *PerfTodo) Key(c appengine.Context) *datastore.Key
```

#### type Result

```go
type Result struct {
	PackagePath string // (empty for Go commits)
	Builder     string // "os-arch[-note]"
	Hash        string

	// The Go Commit this was built against (empty for Go commits).
	GoHash string

	OK      bool
	Log     string `datastore:"-"`        // for JSON unmarshaling only
	LogHash string `datastore:",noindex"` // Key to the Log record.

	RunTime int64 // time to build+test in nanoseconds
}
```

A Result describes a build result for a Commit on an OS/architecture.

Each Result entity is a descendant of its associated Package entity.

#### func (*Result) Data

```go
func (r *Result) Data() string
```
Data returns the Result in string format to be stored in Commit's ResultData
field.

#### func (*Result) Key

```go
func (r *Result) Key(c appengine.Context) *datastore.Key
```

#### func (*Result) Valid

```go
func (r *Result) Valid() error
```

#### type Tag

```go
type Tag struct {
	Kind string // "weekly", "release", or "tip"
	Name string // the tag itself (for example: "release.r60")
	Hash string
}
```

A Tag is used to keep track of the most recent Go weekly and release tags.
Typically there will be one Tag entity for each kind of hg tag.

#### func  GetTag

```go
func GetTag(c appengine.Context, tag string) (*Tag, error)
```
GetTag fetches a Tag by name from the datastore.

#### func (*Tag) Commit

```go
func (t *Tag) Commit(c appengine.Context) (*Commit, error)
```
Commit returns the Commit that corresponds with this Tag.

#### func (*Tag) Key

```go
func (t *Tag) Key(c appengine.Context) *datastore.Key
```

#### func (*Tag) Valid

```go
func (t *Tag) Valid() error
```

#### type TagState

```go
type TagState struct {
	Tag      *Commit
	Packages []*PackageState
}
```

TagState represents the state of all Packages at a Tag.

#### func  TagStateByName

```go
func TagStateByName(c appengine.Context, name string) (*TagState, error)
```
TagStateByName fetches the results for all Go subrepos at the specified Tag.

#### type Todo

```go
type Todo struct {
	Kind string // "build-go-commit" or "build-package"
	Data interface{}
}
```

Todo is a todoHandler response.
