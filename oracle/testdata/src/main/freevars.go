// antha-tools/oracle/testdata/src/main/freevars.go: Part of the Antha language
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

package main

// Tests of 'freevars' query.
// See go.tools/oracle/oracle_test.go for explanation.
// See freevars.golden for expected query results.

// TODO(adonovan): it's hard to test this query in a single line of gofmt'd code.

type T struct {
	a, b int
}

type S struct {
	x int
	t T
}

func f(int) {}

func main() {
	type C int
	x := 1
	const exp = 6
	if y := 2; x+y+int(C(3)) != exp { // @freevars fv1 "if.*{"
		panic("expected 6")
	}

	var s S

	for x, y := range "foo" {
		println(s.x + s.t.a + s.t.b + x + int(y)) // @freevars fv2 "print.*y."
	}

	f(x) // @freevars fv3 "f.x."

	// TODO(adonovan): enable when antha/types supports labels.
loop: // #@freevars fv-def-label "loop:"
	for {
		break loop // #@freevars fv-ref-label "break loop"
	}
}