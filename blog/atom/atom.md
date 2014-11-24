# atom
--
    import "."

Package atom defines XML data structures for an Atom feed.

## Usage

#### type Entry

```go
type Entry struct {
	Title     string  `xml:"title"`
	ID        string  `xml:"id"`
	Link      []Link  `xml:"link"`
	Published TimeStr `xml:"published"`
	Updated   TimeStr `xml:"updated"`
	Author    *Person `xml:"author"`
	Summary   *Text   `xml:"summary"`
	Content   *Text   `xml:"content"`
}
```


#### type Feed

```go
type Feed struct {
	XMLName xml.Name `xml:"http://www.w3.org/2005/Atom feed"`
	Title   string   `xml:"title"`
	ID      string   `xml:"id"`
	Link    []Link   `xml:"link"`
	Updated TimeStr  `xml:"updated"`
	Author  *Person  `xml:"author"`
	Entry   []*Entry `xml:"entry"`
}
```


#### type Link

```go
type Link struct {
	Rel  string `xml:"rel,attr"`
	Href string `xml:"href,attr"`
}
```


#### type Person

```go
type Person struct {
	Name     string `xml:"name"`
	URI      string `xml:"uri,omitempty"`
	Email    string `xml:"email,omitempty"`
	InnerXML string `xml:",innerxml"`
}
```


#### type Text

```go
type Text struct {
	Type string `xml:"type,attr"`
	Body string `xml:",chardata"`
}
```


#### type TimeStr

```go
type TimeStr string
```


#### func  Time

```go
func Time(t time.Time) TimeStr
```
