// antha-tools/antha/pointer/testdata/arrays.go: Part of the Antha language
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

var unknown bool // defeat dead-code elimination

var a, b int

func array1() {
	sliceA := make([]*int, 10) // @line a1make
	sliceA[0] = &a

	var sliceB []*int
	sliceB = append(sliceB, &b) // @line a1append

	print(sliceA)    // @pointsto makeslice@a1make:16
	print(sliceA[0]) // @pointsto main.a

	print(sliceB)      // @pointsto append@a1append:17
	print(sliceB[100]) // @pointsto main.b
}

func array2() {
	sliceA := make([]*int, 10) // @line a2make
	sliceA[0] = &a

	sliceB := sliceA[:]

	print(sliceA)    // @pointsto makeslice@a2make:16
	print(sliceA[0]) // @pointsto main.a

	print(sliceB)    // @pointsto makeslice@a2make:16
	print(sliceB[0]) // @pointsto main.a
}

func array3() {
	a := []interface{}{"", 1}
	b := []interface{}{true, func() {}}
	print(a[0]) // @types string | int
	print(b[0]) // @types bool | func()
}

// Test of append, copy, slice.
func array4() {
	var s2 struct { // @line a4L0
		a [3]int
		b struct{ c, d int }
	}
	var sl1 = make([]*int, 10) // @line a4make
	var someint int            // @line a4L1
	sl1[1] = &someint
	sl2 := append(sl1, &s2.a[1]) // @line a4append1
	print(sl1)                   // @pointsto makeslice@a4make:16
	print(sl2)                   // @pointsto append@a4append1:15 | makeslice@a4make:16
	print(sl1[0])                // @pointsto someint@a4L1:6 | s2.a[*]@a4L0:6
	print(sl2[0])                // @pointsto someint@a4L1:6 | s2.a[*]@a4L0:6

	// In z=append(x,y) we should observe flow from y[*] to x[*].
	var sl3 = make([]*int, 10) // @line a4L2
	_ = append(sl3, &s2.a[1])
	print(sl3)    // @pointsto makeslice@a4L2:16
	print(sl3[0]) // @pointsto s2.a[*]@a4L0:6

	var sl4 = []*int{&a} // @line a4L3
	sl4a := append(sl4)  // @line a4L4
	print(sl4a)          // @pointsto slicelit@a4L3:18 | append@a4L4:16
	print(&sl4a[0])      // @pointsto slicelit[*]@a4L3:18 | append[*]@a4L4:16
	print(sl4a[0])       // @pointsto main.a

	var sl5 = []*int{&b} // @line a4L5
	copy(sl5, sl4)
	print(sl5)     // @pointsto slicelit@a4L5:18
	print(&sl5[0]) // @pointsto slicelit[*]@a4L5:18
	print(sl5[0])  // @pointsto main.b | main.a

	var sl6 = sl5[:0]
	print(sl6)     // @pointsto slicelit@a4L5:18
	print(&sl6[0]) // @pointsto slicelit[*]@a4L5:18
	print(sl6[0])  // @pointsto main.b | main.a
}

func array5() {
	var arr [2]*int
	arr[0] = &a
	arr[1] = &b

	var n int
	print(arr[n]) // @pointsto main.a | main.b
}

func main() {
	array1()
	array2()
	array3()
	array4()
	array5()
}