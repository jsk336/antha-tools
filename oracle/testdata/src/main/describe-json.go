// antha-tools/oracle/testdata/src/main/describe-json.go: Part of the Antha language
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

package describe // @describe pkgdecl "describe"

// Tests of 'describe' query, -format=json.
// See go.tools/oracle/oracle_test.go for explanation.
// See describe-json.golden for expected query results.

func main() { //
	var s struct{ x [3]int }
	p := &s.x[0] // @describe desc-val-p "p"
	_ = p

	var i I = C(0)
	if i == nil {
		i = new(D)
	}
	print(i) // @describe desc-val-i "\\bi\\b"

	go main() // @describe desc-stmt "go"
}

type I interface {
	f()
}

type C int // @describe desc-type-C "C"
type D struct{}

func (c C) f()  {}
func (d *D) f() {}