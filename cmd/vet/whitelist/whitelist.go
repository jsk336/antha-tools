// antha-tools/cmd/vet/whitelist/whitelist.go: Part of the Antha language
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


// Package whitelist defines exceptions for the vet tool.
package whitelist

// UnkeyedLiteral are types that are actually slices, but
// syntactically, we cannot tell whether the Typ in pkg.Typ{1, 2, 3}
// is a slice or a struct, so we whitelist all the standard package
// library's exported slice types.
var UnkeyedLiteral = map[string]bool{
	/*
		find $GOROOT/src/pkg -type f | grep -v _test.go | xargs grep '^type.*\[\]' | \
			grep -v ' map\[' | sed 's,/[^/]*go.type,,' | sed 's,.*src/pkg/,,' | \
			sed 's, ,.,' |  sed 's, .*,,' | grep -v '\.[a-z]' | \
			sort | awk '{ print "\"" $0 "\": true," }'
	*/
	"crypto/x509/pkix.RDNSequence":                  true,
	"crypto/x509/pkix.RelativeDistinguishedNameSET": true,
	"database/sql.RawBytes":                         true,
	"debug/macho.LoadBytes":                         true,
	"encoding/asn1.ObjectIdentifier":                true,
	"encoding/asn1.RawContent":                      true,
	"encoding/json.RawMessage":                      true,
	"encoding/xml.CharData":                         true,
	"encoding/xml.Comment":                          true,
	"encoding/xml.Directive":                        true,
	"github.com/antha-lang/antha/scanner.ErrorList":                          true,
	"image/color.Palette":                           true,
	"net.HardwareAddr":                              true,
	"net.IP":                                        true,
	"net.IPMask":                                    true,
	"sort.Float64Slice":                             true,
	"sort.IntSlice":                                 true,
	"sort.StringSlice":                              true,
	"unicode.SpecialCase":                           true,

	// These image and image/color struct types are frozen. We will never add fields to them.
	"image/color.Alpha16": true,
	"image/color.Alpha":   true,
	"image/color.Gray16":  true,
	"image/color.Gray":    true,
	"image/color.NRGBA64": true,
	"image/color.NRGBA":   true,
	"image/color.RGBA64":  true,
	"image/color.RGBA":    true,
	"image/color.YCbCr":   true,
	"image.Point":         true,
	"image.Rectangle":     true,
	"image.Uniform":       true,
}