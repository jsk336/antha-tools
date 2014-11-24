// antha-tools/blog/atom/atom.go: Part of the Antha language
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


// Adapted from encoding/xml/read_test.go.

// Package atom defines XML data structures for an Atom feed.
package atom

import (
	"encoding/xml"
	"time"
)

type Feed struct {
	XMLName xml.Name `xml:"http://www.w3.org/2005/Atom feed"`
	Title   string   `xml:"title"`
	ID      string   `xml:"id"`
	Link    []Link   `xml:"link"`
	Updated TimeStr  `xml:"updated"`
	Author  *Person  `xml:"author"`
	Entry   []*Entry `xml:"entry"`
}

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

type Link struct {
	Rel  string `xml:"rel,attr"`
	Href string `xml:"href,attr"`
}

type Person struct {
	Name     string `xml:"name"`
	URI      string `xml:"uri,omitempty"`
	Email    string `xml:"email,omitempty"`
	InnerXML string `xml:",innerxml"`
}

type Text struct {
	Type string `xml:"type,attr"`
	Body string `xml:",chardata"`
}

type TimeStr string

func Time(t time.Time) TimeStr {
	return TimeStr(t.Format("2006-01-02T15:04:05-07:00"))
}