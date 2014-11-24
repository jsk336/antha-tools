// antha-tools/oracle/testdata/src/main/referrers-json.go: Part of the Antha language
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

package referrers

// Tests of 'referrers' query.
// See go.tools/oracle/oracle_test.go for explanation.
// See referrers.golden for expected query results.

import "lib"

type s struct {
	f int
}

func main() {
	var v lib.Type = lib.Const // @referrers ref-package "lib"
	_ = v.Method               // @referrers ref-method "Method"
	_ = v.Method
	v++ //@referrers ref-local "v"
	v++

	_ = s{}.f // @referrers ref-field "f"

	var s2 s
	s2.f = 1
}