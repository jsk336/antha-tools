// antha-tools/go/ssa/interp/testdata/initorder.go: Part of the Antha language
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

// Test of initialization order of package-level vars.

var counter int

func next() int {
	c := counter
	counter++
	return c
}

func next2() (x int, y int) {
	x = next()
	y = next()
	return
}

func makeOrder() int {
	_, _, _, _ = f, b, d, e
	return 0
}

func main() {
	// Initialization constraints:
	// - {f,b,c/d,e} < order  (ref graph traversal)
	// - order < {a}          (lexical order)
	// - b < c/d < e < f      (lexical order)
	// Solution: b c/d e f a
	abcdef := [6]int{a, b, c, d, e, f}
	if abcdef != [6]int{5, 0, 1, 2, 3, 4} {
		panic(abcdef)
	}
}

var order = makeOrder()

var a, b = next(), next()
var c, d = next2()
var e, f = next(), next()