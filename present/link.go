// antha-tools/present/link.go: Part of the Antha language
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


package present

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

func init() {
	Register("link", parseLink)
}

type Link struct {
	URL   *url.URL
	Label string
}

func (l Link) TemplateName() string { return "link" }

func parseLink(ctx *Context, fileName string, lineno int, text string) (Elem, error) {
	args := strings.Fields(text)
	url, err := url.Parse(args[1])
	if err != nil {
		return nil, err
	}
	label := ""
	if len(args) > 2 {
		label = strings.Join(args[2:], " ")
	} else {
		scheme := url.Scheme + "://"
		if url.Scheme == "mailto" {
			scheme = "mailto:"
		}
		label = strings.Replace(url.String(), scheme, "", 1)
	}
	return Link{url, label}, nil
}

func renderLink(href, text string) string {
	text = font(text)
	if text == "" {
		text = href
	}
	// Open links in new window only when their url is absolute.
	target := "_blank"
	if u, err := url.Parse(href); err != nil {
		log.Println("rendernLink parsing url:", err)
	} else if !u.IsAbs() || u.Scheme == "javascript" {
		target = "_self"
	}

	return fmt.Sprintf(`<a href="%s" target="%s">%s</a>`, href, target, text)
}

// parseInlineLink parses an inline link at the start of s, and returns
// a rendered HTML link and the total length of the raw inline link.
// If no inline link is present, it returns all zeroes.
func parseInlineLink(s string) (link string, length int) {
	if !strings.HasPrefix(s, "[[") {
		return
	}
	end := strings.Index(s, "]]")
	if end == -1 {
		return
	}
	urlEnd := strings.Index(s, "]")
	rawURL := s[2:urlEnd]
	const badURLChars = `<>"{}|\^[] ` + "`" // per RFC2396 section 2.4.3
	if strings.ContainsAny(rawURL, badURLChars) {
		return
	}
	if urlEnd == end {
		simpleUrl := ""
		url, err := url.Parse(rawURL)
		if err == nil {
			// If the URL is http://foo.com, drop the http://
			// In other words, render [[http://golang.org]] as:
			//   <a href="http://golang.org">golang.org</a>
			if strings.HasPrefix(rawURL, url.Scheme+"://") {
				simpleUrl = strings.TrimPrefix(rawURL, url.Scheme+"://")
			} else if strings.HasPrefix(rawURL, url.Scheme+":") {
				simpleUrl = strings.TrimPrefix(rawURL, url.Scheme+":")
			}
		}
		return renderLink(rawURL, simpleUrl), end + 2
	}
	if s[urlEnd:urlEnd+2] != "][" {
		return
	}
	text := s[urlEnd+2 : end]
	return renderLink(rawURL, text), end + 2
}