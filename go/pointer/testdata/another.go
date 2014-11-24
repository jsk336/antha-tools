// antha-tools/go/pointer/testdata/another.go: Part of the Antha language
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

var unknown bool

type S string

func incr(x int) int { return x + 1 }

func main() {
	var i interface{}
	i = 1
	if unknown {
		i = S("foo")
	}
	if unknown {
		i = (func(int, int))(nil) // NB type compares equal to that below.
	}
	// Look, the test harness can handle equal-but-not-String-equal
	// types because we parse types and using a typemap.
	if unknown {
		i = (func(x int, y int))(nil)
	}
	if unknown {
		i = incr
	}
	print(i) // @types int | S | func(int, int) | func(int) int

	// NB, an interface may never directly alias any global
	// labels, even though it may contain pointers that do.
	print(i)                 // @pointsto makeinterface:func(x int) int | makeinterface:func(x int, y int) | makeinterface:func(int, int) | makeinterface:int | makeinterface:main.S
	print(i.(func(int) int)) // @pointsto main.incr
}