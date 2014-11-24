// antha-tools/cmd/vet/testdata/rangeloop.go: Part of the Antha language
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


// This file contains tests for the rangeloop checker.

package testdata

func RangeLoopTests() {
	var s []int
	for i, v := range s {
		go func() {
			println(i) // ERROR "range variable i enclosed by function"
			println(v) // ERROR "range variable v enclosed by function"
		}()
	}
	for i, v := range s {
		defer func() {
			println(i) // ERROR "range variable i enclosed by function"
			println(v) // ERROR "range variable v enclosed by function"
		}()
	}
	for i := range s {
		go func() {
			println(i) // ERROR "range variable i enclosed by function"
		}()
	}
	for _, v := range s {
		go func() {
			println(v) // ERROR "range variable v enclosed by function"
		}()
	}
	for i, v := range s {
		go func() {
			println(i, v)
		}()
		println("unfortunately, we don't catch the error above because of this statement")
	}
	for i, v := range s {
		go func(i, v int) {
			println(i, v)
		}(i, v)
	}
	for i, v := range s {
		i, v := i, v
		go func() {
			println(i, v)
		}()
	}
	// If the key of the range statement is not an identifier
	// the code should not panic (it used to).
	var x [2]int
	var f int
	for x[0], f = range s {
		go func() {
			_ = f // ERROR "range variable f enclosed by function"
		}()
	}
}