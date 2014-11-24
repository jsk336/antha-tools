// antha-tools/antha/pointer/testdata/panic.go: Part of the Antha language
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

// Test of value flow from panic() to recover().
// We model them as stores/loads of a global location.
// We ignore concrete panic types originating from the runtime.

var someval int

type myPanic struct{}

func f(int) {}

func g() string { return "" }

func deadcode() {
	panic(123) // not reached
}

func main() {
	switch someval {
	case 0:
		panic("oops")
	case 1:
		panic(myPanic{})
	case 2:
		panic(f)
	case 3:
		panic(g)
	}
	ex := recover()
	print(ex)                 // @types myPanic | string | func(int) | func() string
	print(ex.(func(int)))     // @pointsto main.f
	print(ex.(func() string)) // @pointsto main.g
}