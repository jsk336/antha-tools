// antha-tools/anthadoc/search.go: Part of the Antha language
// Copyright (C) 2014 The Antha authors. All rights reserved.
// 
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
// 
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
// 
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o 
// Synthace Ltd. The London Bioscience Innovation Centre
// 1 Royal College St, London NW1 0NH UK


package godoc

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type SearchResult struct {
	Query string
	Alert string // error or warning message

	// identifier matches
	Pak HitList       // packages matching Query
	Hit *LookupResult // identifier matches of Query
	Alt *AltWords     // alternative identifiers to look for

	// textual matches
	Found    int         // number of textual occurrences found
	Textual  []FileLines // textual matches of Query
	Complete bool        // true if all textual occurrences of Query are reported
	Idents   map[SpotKind][]Ident
}

func (c *Corpus) Lookup(query string) SearchResult {
	result := &SearchResult{Query: query}

	index, timestamp := c.CurrentIndex()
	if index != nil {
		// identifier search
		if r, err := index.Lookup(query); err == nil {
			result = r
		} else if err != nil && !c.IndexFullText {
			// ignore the error if full text search is enabled
			// since the query may be a valid regular expression
			result.Alert = "Error in query string: " + err.Error()
			return *result
		}

		// full text search
		if c.IndexFullText && query != "" {
			rx, err := regexp.Compile(query)
			if err != nil {
				result.Alert = "Error in query regular expression: " + err.Error()
				return *result
			}
			// If we get maxResults+1 results we know that there are more than
			// maxResults results and thus the result may be incomplete (to be
			// precise, we should remove one result from the result set, but
			// nobody is going to count the results on the result page).
			result.Found, result.Textual = index.LookupRegexp(rx, c.MaxResults+1)
			result.Complete = result.Found <= c.MaxResults
			if !result.Complete {
				result.Found-- // since we looked for maxResults+1
			}
		}
	}

	// is the result accurate?
	if c.IndexEnabled {
		if ts := c.FSModifiedTime(); timestamp.Before(ts) {
			// The index is older than the latest file system change under godoc's observation.
			result.Alert = "Indexing in progress: result may be inaccurate"
		}
	} else {
		result.Alert = "Search index disabled: no results available"
	}

	return *result
}

// SearchResultDoc optionally specifies a function returning an HTML body
// displaying search results matching godoc documentation.
func (p *Presentation) SearchResultDoc(result SearchResult) []byte {
	return applyTemplate(p.SearchDocHTML, "searchDocHTML", result)
}

// SearchResultCode optionally specifies a function returning an HTML body
// displaying search results matching source code.
func (p *Presentation) SearchResultCode(result SearchResult) []byte {
	return applyTemplate(p.SearchCodeHTML, "searchCodeHTML", result)
}

// SearchResultTxt optionally specifies a function returning an HTML body
// displaying search results of textual matches.
func (p *Presentation) SearchResultTxt(result SearchResult) []byte {
	return applyTemplate(p.SearchTxtHTML, "searchTxtHTML", result)
}

// HandleSearch obtains results for the requested search and returns a page
// to display them.
func (p *Presentation) HandleSearch(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.FormValue("q"))
	result := p.Corpus.Lookup(query)

	if p.GetPageInfoMode(r)&NoHTML != 0 {
		p.ServeText(w, applyTemplate(p.SearchText, "searchText", result))
		return
	}
	contents := bytes.Buffer{}
	for _, f := range p.SearchResults {
		contents.Write(f(p, result))
	}

	var title string
	if haveResults := contents.Len() > 0; haveResults {
		title = fmt.Sprintf(`Results for query %q`, query)
		if !p.Corpus.IndexEnabled {
			result.Alert = ""
		}
	} else {
		title = fmt.Sprintf(`No results found for query %q`, query)
	}

	body := bytes.NewBuffer(applyTemplate(p.SearchHTML, "searchHTML", result))
	body.Write(contents.Bytes())

	p.ServePage(w, Page{
		Title:    title,
		Tabtitle: query,
		Query:    query,
		Body:     body.Bytes(),
	})
}

func (p *Presentation) serveSearchDesc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/opensearchdescription+xml")
	data := map[string]interface{}{
		"BaseURL": fmt.Sprintf("http://%s", r.Host),
	}
	applyTemplateToResponseWriter(w, p.SearchDescXML, &data)
}