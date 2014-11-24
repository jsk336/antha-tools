// antha-tools/go/pointer/testdata/chanreflect1.go: Part of the Antha language
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

import "reflect"

//
// This test is very sensitive to line-number perturbations!

// Test of channels with reflection.

var a, b int

func chanreflect1() {
	ch := make(chan *int, 0)
	crv := reflect.ValueOf(ch)
	crv.Send(reflect.ValueOf(&a))
	print(crv.Interface())             // @types chan *int
	print(crv.Interface().(chan *int)) // @pointsto makechan@testdata/chanreflect.go:15:12
	print(<-ch)                        // @pointsto main.a
}

func chanreflect2() {
	ch := make(chan *int, 0)
	ch <- &b
	crv := reflect.ValueOf(ch)
	r, _ := crv.Recv()
	print(r.Interface())        // @types *int
	print(r.Interface().(*int)) // @pointsto main.b
}

func main() {
	chanreflect1()
	chanreflect2()
}