// antha-tools/cmd/vet/structtag.go: Part of the Antha language
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


// This file contains the test for canonical struct tags.

package main

import (
	"github.com/antha-lang/antha/ast"
	"reflect"
	"strconv"
)

// checkField checks a struct field tag.
func (f *File) checkCanonicalFieldTag(field *ast.Field) {
	if !vet("structtags") {
		return
	}
	if field.Tag == nil {
		return
	}

	tag, err := strconv.Unquote(field.Tag.Value)
	if err != nil {
		f.Badf(field.Pos(), "unable to read struct tag %s", field.Tag.Value)
		return
	}

	// Check tag for validity by appending
	// new key:value to end and checking that
	// the tag parsing code can find it.
	if reflect.StructTag(tag+` _gofix:"_magic"`).Get("_gofix") != "_magic" {
		f.Badf(field.Pos(), "struct field tag %s not compatible with reflect.StructTag.Get", field.Tag.Value)
		return
	}
}