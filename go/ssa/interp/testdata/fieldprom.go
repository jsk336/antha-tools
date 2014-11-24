// antha-tools/go/ssa/interp/testdata/fieldprom.go: Part of the Antha language
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

// Tests of field promotion logic.

type A struct {
	x int
	y *int
}

type B struct {
	p int
	q *int
}

type C struct {
	A
	*B
}

type D struct {
	a int
	C
}

func assert(cond bool) {
	if !cond {
		panic("failed")
	}
}

func f1(c C) {
	assert(c.x == c.A.x)
	assert(c.y == c.A.y)
	assert(&c.x == &c.A.x)
	assert(&c.y == &c.A.y)

	assert(c.p == c.B.p)
	assert(c.q == c.B.q)
	assert(&c.p == &c.B.p)
	assert(&c.q == &c.B.q)

	c.x = 1
	*c.y = 1
	c.p = 1
	*c.q = 1
}

func f2(c *C) {
	assert(c.x == c.A.x)
	assert(c.y == c.A.y)
	assert(&c.x == &c.A.x)
	assert(&c.y == &c.A.y)

	assert(c.p == c.B.p)
	assert(c.q == c.B.q)
	assert(&c.p == &c.B.p)
	assert(&c.q == &c.B.q)

	c.x = 1
	*c.y = 1
	c.p = 1
	*c.q = 1
}

func f3(d D) {
	assert(d.x == d.C.A.x)
	assert(d.y == d.C.A.y)
	assert(&d.x == &d.C.A.x)
	assert(&d.y == &d.C.A.y)

	assert(d.p == d.C.B.p)
	assert(d.q == d.C.B.q)
	assert(&d.p == &d.C.B.p)
	assert(&d.q == &d.C.B.q)

	d.x = 1
	*d.y = 1
	d.p = 1
	*d.q = 1
}

func f4(d *D) {
	assert(d.x == d.C.A.x)
	assert(d.y == d.C.A.y)
	assert(&d.x == &d.C.A.x)
	assert(&d.y == &d.C.A.y)

	assert(d.p == d.C.B.p)
	assert(d.q == d.C.B.q)
	assert(&d.p == &d.C.B.p)
	assert(&d.q == &d.C.B.q)

	d.x = 1
	*d.y = 1
	d.p = 1
	*d.q = 1
}

func main() {
	y := 123
	c := C{
		A{x: 42, y: &y},
		&B{p: 42, q: &y},
	}

	assert(&c.x == &c.A.x)

	f1(c)
	f2(&c)

	d := D{C: c}
	f3(d)
	f4(&d)
}