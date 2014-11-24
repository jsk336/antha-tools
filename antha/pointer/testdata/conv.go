// antha-tools/antha/pointer/testdata/conv.go: Part of the Antha language
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

import "unsafe"

var a int

func conv1() {
	// Conversions of channel direction.
	ch := make(chan int)    // @line c1make
	print((<-chan int)(ch)) // @pointsto makechan@c1make:12
	print((chan<- int)(ch)) // @pointsto makechan@c1make:12
}

func conv2() {
	// string -> []byte/[]rune conversion
	s := "foo"
	ba := []byte(s) // @line c2ba
	ra := []rune(s) // @line c2ra
	print(ba)       // @pointsto convert@c2ba:14
	print(ra)       // @pointsto convert@c2ra:14
}

func conv3() {
	// Conversion of same underlying types.
	type PI *int
	pi := PI(&a)
	print(pi) // @pointsto main.a

	pint := (*int)(pi)
	print(pint) // @pointsto main.a

	// Conversions between pointers to identical base types.
	var y *PI = &pi
	var x **int = (**int)(y)
	print(*x) // @pointsto main.a
	print(*y) // @pointsto main.a
	y = (*PI)(x)
	print(*y) // @pointsto main.a
}

// @warning "main.conv4 contains an unsafe.Pointer conversion"
func conv4() {
	// Handling of unsafe.Pointer conversion is unsound:
	// we lose the alias to main.a and get something like new(int) instead.
	// We require users to provide aliasing summaries.
	p := (*int)(unsafe.Pointer(&a)) // @line c2p
	print(p)                        // @pointsto convert@c2p:13
}

// Regression test for b/8231.
func conv5() {
	type P unsafe.Pointer
	var i *struct{}
	_ = P(i)
}

func main() {
	conv1()
	conv2()
	conv3()
	conv4()
	conv5()
}