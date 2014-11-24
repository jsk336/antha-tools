// antha-tools/antha/vcs/discovery.go: Part of the Antha language
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


package vcs

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// charsetReader returns a reader for the given charset. Currently
// it only supports UTF-8 and ASCII. Otherwise, it returns a meaningful
// error which is printed by antha get, so the user can find why the package
// wasn't downloaded if the encoding is not supported. Note that, in
// order to reduce potential errors, ASCII is treated as UTF-8 (i.e. characters
// greater than 0x7f are not rejected).
func charsetReader(charset string, input io.Reader) (io.Reader, error) {
	switch strings.ToLower(charset) {
	case "ascii":
		return input, nil
	default:
		return nil, fmt.Errorf("can't decode XML document using charset %q", charset)
	}
}

// parseMetaGoImports returns meta imports from the HTML in r.
// Parsing ends at the end of the <head> section or the beginning of the <body>.
func parseMetaGoImports(r io.Reader) (imports []metaImport, err error) {
	d := xml.NewDecoder(r)
	d.CharsetReader = charsetReader
	d.Strict = false
	var t xml.Token
	for {
		t, err = d.Token()
		if err != nil {
			return
		}
		if e, ok := t.(xml.StartElement); ok && strings.EqualFold(e.Name.Local, "body") {
			return
		}
		if e, ok := t.(xml.EndElement); ok && strings.EqualFold(e.Name.Local, "head") {
			return
		}
		e, ok := t.(xml.StartElement)
		if !ok || !strings.EqualFold(e.Name.Local, "meta") {
			continue
		}
		if attrValue(e.Attr, "name") != "go-import" {
			continue
		}
		if f := strings.Fields(attrValue(e.Attr, "content")); len(f) == 3 {
			imports = append(imports, metaImport{
				Prefix:   f[0],
				VCS:      f[1],
				RepoRoot: f[2],
			})
		}
	}
}

// attrValue returns the attribute value for the case-insensitive key
// `name', or the empty string if nothing is found.
func attrValue(attrs []xml.Attr, name string) string {
	for _, a := range attrs {
		if strings.EqualFold(a.Name.Local, name) {
			return a.Value
		}
	}
	return ""
}