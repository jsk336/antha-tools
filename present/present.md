# present
--
    import "."

The present file format

Present files have the following format. The first non-blank non-comment line is
the title, so the header looks like

    Title of document
    Subtitle of document
    15:04 2 Jan 2006
    Tags: foo, bar, baz
    <blank line>
    Author Name
    Job title, Company
    joe@example.com
    http://url/
    @twitter_name

The subtitle, date, and tags lines are optional.

The date line may be written without a time:

    2 Jan 2006

In this case, the time will be interpreted as 10am UTC on that date.

The tags line is a comma-separated list of tags that may be used to categorize
the document.

The author section may contain a mixture of text, twitter names, and links. For
slide presentations, only the plain text lines will be displayed on the first
slide.

Multiple presenters may be specified, separated by a blank line.

After that come slides/sections, each after a blank line:

    * Title of slide or section (must have asterisk)

    Some Text

    ** Subsection

    - bullets
    - more bullets
    - a bullet with

    *** Sub-subsection

    Some More text

      Preformatted text
      is indented (however you like)

    Further Text, including invocations like:

    .code x.go /^func main/,/^}/
    .play y.go
    .image image.jpg
    .iframe http://foo
    .link http://foo label
    .html file.html

    Again, more text

Blank lines are OK (not mandatory) after the title and after the text. Text,
bullets, and .code etc. are all optional; title is not.

Lines starting with # in column 1 are commentary.

Fonts:

Within the input for plain text or lists, text bracketed by font markers will be
presented in italic, bold, or program font. Marker characters are _ (italic), *
(bold) and ` (program font). Unmatched markers appear as plain text. Within
marked text, a single marker character becomes a space and a doubled single
marker quotes the marker character.

    _italic_
    *bold*
    `program`
    _this_is_all_italic_
    _Why_use_scoped__ptr_? Use plain ***ptr* instead.

Inline links:

Links can be included in any text with the form [[url][label]], or [[url]] to
use the URL itself as the label.

Functions:

A number of template functions are available through invocations in the input
text. Each such invocation contains a period as the first character on the line,
followed immediately by the name of the function, followed by any arguments. A
typical invocation might be

    .play demo.go /^func show/,/^}/

(except that the ".play" must be at the beginning of the line and not be
indented like this.)

Here follows a description of the functions:

code:

Injects program source into the output by extracting code from files and
injecting them as HTML-escaped <pre> blocks. The argument is a file name
followed by an optional address that specifies what section of the file to
display. The address syntax is similar in its simplest form to that of ed, but
comes from sam and is more general. See

    http://plan9.bell-labs.com/sys/doc/sam/sam.html Table II

for full details. The displayed block is always rounded out to a full line at
both ends.

If no pattern is present, the entire file is displayed.

Any line in the program that ends with the four characters

    OMIT

is deleted from the source before inclusion, making it easy to write things like

    .code test.go /START OMIT/,/END OMIT/

to find snippets like this

    tedious_code = boring_function()
    // START OMIT
    interesting_code = fascinating_function()
    // END OMIT

and see only this:

    interesting_code = fascinating_function()

Also, inside the displayed text a line that ends

    // HL

will be highlighted in the display; the 'h' key in the browser will toggle extra
emphasis of any highlighted lines. A highlighting mark may have a suffix word,
such as

    // HLxxx

Such highlights are enabled only if the code invocation ends with "HL" followed
by the word:

    .code test.go /^type Foo/,/^}/ HLxxx

The .code function may take one or more flags immediately preceding the
filename. This command shows test.go in an editable text area:

    .code -edit test.go

This command shows test.go with line numbers:

    .code -numbers test.go

play:

The function "play" is the same as "code" but puts a button on the displayed
source so the program can be run from the browser. Although only the selected
text is shown, all the source is included in the HTML output so it can be
presented to the compiler.

link:

Create a hyperlink. The syntax is 1 or 2 space-separated arguments. The first
argument is always the HTTP URL. If there is a second argument, it is the text
label to display for this link.

    .link http://golang.org golang.org

image:

The template uses the function "image" to inject picture files.

The syntax is simple: 1 or 3 space-separated arguments. The first argument is
always the file name. If there are more arguments, they are the height and
width; both must be present, or substituted with an underscore. Replacing a
dimension argument with the underscore parameter preserves the aspect ratio of
the image when scaling.

    .image images/betsy.jpg 100 200

    .image images/janet.jpg _ 300

iframe:

The function "iframe" injects iframes (pages inside pages). Its syntax is the
same as that of image.

html:

The function html includes the contents of the specified file as unescaped HTML.
This is useful for including custom HTML elements that cannot be created using
only the slide format. It is your responsibilty to make sure the included HTML
is valid and safe.

    .html file.html

## Usage

```go
var PlayEnabled = false
```
Is the playground available?

#### func  Register

```go
func Register(name string, parser ParseFunc)
```
Register binds the named action, which does not begin with a period, to the
specified parser to be invoked when the name, with a period, appears in the
present input text.

#### func  Style

```go
func Style(s string) template.HTML
```
Style returns s with HTML entities escaped and font indicators turned into HTML
font tags.

#### func  Template

```go
func Template() *template.Template
```
Template returns an empty template with the action functions in its FuncMap.

#### type Author

```go
type Author struct {
	Elem []Elem
}
```

Author represents the person who wrote and/or is presenting the document.

#### func (*Author) TextElem

```go
func (p *Author) TextElem() (elems []Elem)
```
TextElem returns the first text elements of the author details. This is used to
display the author' name, job title, and company without the contact details.

#### type Code

```go
type Code struct {
	Text     template.HTML
	Play     bool   // runnable code
	FileName string // file name
	Ext      string // file extension
	Raw      []byte // content of the file
}
```


#### func (Code) TemplateName

```go
func (c Code) TemplateName() string
```

#### type Context

```go
type Context struct {
	// ReadFile reads the file named by filename and returns the contents.
	ReadFile func(filename string) ([]byte, error)
}
```

A Context specifies the supporting context for parsing a presentation.

#### func (*Context) Parse

```go
func (ctx *Context) Parse(r io.Reader, name string, mode ParseMode) (*Doc, error)
```
Parse parses a document from r.

#### type Doc

```go
type Doc struct {
	Title    string
	Subtitle string
	Time     time.Time
	Authors  []Author
	Sections []Section
	Tags     []string
}
```

Doc represents an entire document.

#### func  Parse

```go
func Parse(r io.Reader, name string, mode ParseMode) (*Doc, error)
```
Parse parses a document from r. Parse reads assets used by the presentation from
the file system using ioutil.ReadFile.

#### func (*Doc) Render

```go
func (d *Doc) Render(w io.Writer, t *template.Template) error
```
Render renders the doc to the given writer using the provided template.

#### type Elem

```go
type Elem interface {
	TemplateName() string
}
```

Elem defines the interface for a present element. That is, something that can
provide the name of the template used to render the element.

#### type HTML

```go
type HTML struct {
	template.HTML
}
```


#### func (HTML) TemplateName

```go
func (s HTML) TemplateName() string
```

#### type Iframe

```go
type Iframe struct {
	URL    string
	Width  int
	Height int
}
```


#### func (Iframe) TemplateName

```go
func (i Iframe) TemplateName() string
```

#### type Image

```go
type Image struct {
	URL    string
	Width  int
	Height int
}
```


#### func (Image) TemplateName

```go
func (i Image) TemplateName() string
```

#### type Lines

```go
type Lines struct {
}
```

Lines is a helper for parsing line-based input.

#### type Link

```go
type Link struct {
	URL   *url.URL
	Label string
}
```


#### func (Link) TemplateName

```go
func (l Link) TemplateName() string
```

#### type List

```go
type List struct {
	Bullet []string
}
```

List represents a bulleted list.

#### func (List) TemplateName

```go
func (l List) TemplateName() string
```

#### type ParseFunc

```go
type ParseFunc func(ctx *Context, fileName string, lineNumber int, inputLine string) (Elem, error)
```


#### type ParseMode

```go
type ParseMode int
```

ParseMode represents flags for the Parse function.

```go
const (
	// If set, parse only the title and subtitle.
	TitlesOnly ParseMode = 1
)
```

#### type Section

```go
type Section struct {
	Number []int
	Title  string
	Elem   []Elem
}
```

Section represents a section of a document (such as a presentation slide)
comprising a title and a list of elements.

#### func (Section) FormattedNumber

```go
func (s Section) FormattedNumber() string
```
FormattedNumber returns a string containing the concatenation of the numbers
identifying a Section.

#### func (Section) Level

```go
func (s Section) Level() int
```
Level returns the level of the given section. The document title is level 1,
main section 2, etc.

#### func (*Section) Render

```go
func (s *Section) Render(w io.Writer, t *template.Template) error
```
Render renders the section to the given writer using the provided template.

#### func (Section) Sections

```go
func (s Section) Sections() (sections []Section)
```

#### func (Section) TemplateName

```go
func (s Section) TemplateName() string
```

#### type Text

```go
type Text struct {
	Lines []string
	Pre   bool
}
```

Text represents an optionally preformatted paragraph.

#### func (Text) TemplateName

```go
func (t Text) TemplateName() string
```
