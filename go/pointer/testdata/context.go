// antha-tools/go/pointer/testdata/context.go: Part of the Antha language
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

// +build ignore

package main

// Test of context-sensitive treatment of certain function calls,
// e.g. static calls to simple accessor methods.

var a, b int

type T struct{ x *int }

func (t *T) SetX(x *int) { t.x = x }
func (t *T) GetX() *int  { return t.x }

func context1() {
	var t1, t2 T
	t1.SetX(&a)
	t2.SetX(&b)
	print(t1.GetX()) // @pointsto main.a
	print(t2.GetX()) // @pointsto main.b
}

func context2() {
	id := func(x *int) *int {
		print(x) // @pointsto main.a | main.b
		return x
	}
	print(id(&a)) // @pointsto main.a
	print(id(&b)) // @pointsto main.b

	// Same again, but anon func has free vars.
	var c int // @line context2c
	id2 := func(x *int) (*int, *int) {
		print(x) // @pointsto main.a | main.b
		return x, &c
	}
	p, q := id2(&a)
	print(p) // @pointsto main.a
	print(q) // @pointsto c@context2c:6
	r, s := id2(&b)
	print(r) // @pointsto main.b
	print(s) // @pointsto c@context2c:6
}

func main() {
	context1()
	context2()
}