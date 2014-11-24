// antha-tools/antha/pointer/testdata/flow.go: Part of the Antha language
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

// Demonstration of directionality of flow edges.

func f1() {}
func f2() {}

var somepred bool

// Tracking functions.
func flow1() {
	s := f1
	p := f2
	q := p
	r := q
	if somepred {
		r = s
	}
	print(s) // @pointsto main.f1
	print(p) // @pointsto main.f2
	print(q) // @pointsto main.f2
	print(r) // @pointsto main.f1 | main.f2
}

// Tracking concrete types in interfaces.
func flow2() {
	var s interface{} = 1
	var p interface{} = "foo"
	q := p
	r := q
	if somepred {
		r = s
	}
	print(s) // @types int
	print(p) // @types string
	print(q) // @types string
	print(r) // @types int | string
}

var g1, g2 int

// Tracking addresses of globals.
func flow3() {
	s := &g1
	p := &g2
	q := p
	r := q
	if somepred {
		r = s
	}
	print(s) // @pointsto main.g1
	print(p) // @pointsto main.g2
	print(q) // @pointsto main.g2
	print(r) // @pointsto main.g2 | main.g1
}

func main() {
	flow1()
	flow2()
	flow3()
}