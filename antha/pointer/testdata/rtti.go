// antha-tools/antha/pointer/testdata/rtti.go: Part of the Antha language
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

// Regression test for oracle crash
// https://code.google.com/p/antha/issues/detail?id=6605
//
// Using reflection, methods may be called on types that are not the
// operand of any ssa.MakeInterface instruction.  In this example,
// (Y).F is called by deriving the type Y from *Y.  Prior to the fix,
// no RTTI (or method set) for type Y was included in the program, so
// the F() call would crash.

import "reflect"

var a int

type X struct{}

func (X) F() *int {
	return &a
}

type I interface {
	F() *int
}

func main() {
	type Y struct{ X }
	print(reflect.Indirect(reflect.ValueOf(new(Y))).Interface().(I).F()) // @pointsto main.a
}