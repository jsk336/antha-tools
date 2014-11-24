// antha-tools/go/ssa/interp/testdata/methprom.go: Part of the Antha language
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

// Tests of method promotion logic.

type A struct{ magic int }

func (a A) x() {
	if a.magic != 1 {
		panic(a.magic)
	}
}
func (a *A) y() *A {
	return a
}

type B struct{ magic int }

func (b B) p() {
	if b.magic != 2 {
		panic(b.magic)
	}
}
func (b *B) q() {
	if b != theC.B {
		panic("oops")
	}
}

type I interface {
	f()
}

type impl struct{ magic int }

func (i impl) f() {
	if i.magic != 3 {
		panic("oops")
	}
}

type C struct {
	A
	*B
	I
}

func assert(cond bool) {
	if !cond {
		panic("failed")
	}
}

var theC = C{
	A: A{1},
	B: &B{2},
	I: impl{3},
}

func addr() *C {
	return &theC
}

func value() C {
	return theC
}

func main() {
	// address
	addr().x()
	if addr().y() != &theC.A {
		panic("oops")
	}
	addr().p()
	addr().q()
	addr().f()

	// addressable value
	var c C = value()
	c.x()
	if c.y() != &c.A {
		panic("oops")
	}
	c.p()
	c.q()
	c.f()

	// non-addressable value
	value().x()
	// value().y() // not in method set
	value().p()
	value().q()
	value().f()
}