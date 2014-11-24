// antha-tools/present/link_test.go: Part of the Antha language
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

import "testing"

func TestInlineParsing(t *testing.T) {
	var tests = []struct {
		in     string
		link   string
		text   string
		length int
	}{
		{"[[http://golang.org]]", "http://golang.org", "golang.org", 21},
		{"[[http://golang.org][]]", "http://golang.org", "http://golang.org", 23},
		{"[[http://golang.org]] this is ignored", "http://golang.org", "golang.org", 21},
		{"[[http://golang.org][link]]", "http://golang.org", "link", 27},
		{"[[http://golang.org][two words]]", "http://golang.org", "two words", 32},
		{"[[http://golang.org][*link*]]", "http://golang.org", "<b>link</b>", 29},
		{"[[http://bad[url]]", "", "", 0},
		{"[[http://golang.org][a [[link]] ]]", "http://golang.org", "a [[link", 31},
		{"[[http:// *spaces* .com]]", "", "", 0},
		{"[[http://bad`char.com]]", "", "", 0},
		{" [[http://google.com]]", "", "", 0},
		{"[[mailto:gopher@golang.org][Gopher]]", "mailto:gopher@golang.org", "Gopher", 36},
		{"[[mailto:gopher@golang.org]]", "mailto:gopher@golang.org", "gopher@golang.org", 28},
	}

	for i, test := range tests {
		link, length := parseInlineLink(test.in)
		if length == 0 && test.length == 0 {
			continue
		}
		if a := renderLink(test.link, test.text); length != test.length || link != a {
			t.Errorf("#%d: parseInlineLink(%q):\ngot\t%q, %d\nwant\t%q, %d", i, test.in, link, length, a, test.length)
		}
	}
}