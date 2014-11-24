// antha-tools/cmd/vet/testdata/unsafeptr.go: Part of the Antha language
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


package testdata

import (
	"reflect"
	"unsafe"
)

func f() {
	var x unsafe.Pointer
	var y uintptr
	x = unsafe.Pointer(y) // ERROR "possible misuse of unsafe.Pointer"
	y = uintptr(x)

	// only allowed pointer arithmetic is ptr +/- num.
	// num+ptr is technically okay but still flagged: write ptr+num instead.
	x = unsafe.Pointer(uintptr(x) + 1)
	x = unsafe.Pointer(1 + uintptr(x))          // ERROR "possible misuse of unsafe.Pointer"
	x = unsafe.Pointer(uintptr(x) + uintptr(x)) // ERROR "possible misuse of unsafe.Pointer"
	x = unsafe.Pointer(uintptr(x) - 1)
	x = unsafe.Pointer(1 - uintptr(x)) // ERROR "possible misuse of unsafe.Pointer"

	// certain uses of reflect are okay
	var v reflect.Value
	x = unsafe.Pointer(v.Pointer())
	x = unsafe.Pointer(v.UnsafeAddr())
	var s1 *reflect.StringHeader
	x = unsafe.Pointer(s1.Data)
	var s2 *reflect.SliceHeader
	x = unsafe.Pointer(s2.Data)
	var s3 reflect.StringHeader
	x = unsafe.Pointer(s3.Data) // ERROR "possible misuse of unsafe.Pointer"
	var s4 reflect.SliceHeader
	x = unsafe.Pointer(s4.Data) // ERROR "possible misuse of unsafe.Pointer"

	// but only in reflect
	var vv V
	x = unsafe.Pointer(vv.Pointer())    // ERROR "possible misuse of unsafe.Pointer"
	x = unsafe.Pointer(vv.UnsafeAddr()) // ERROR "possible misuse of unsafe.Pointer"
	var ss1 *StringHeader
	x = unsafe.Pointer(ss1.Data) // ERROR "possible misuse of unsafe.Pointer"
	var ss2 *SliceHeader
	x = unsafe.Pointer(ss2.Data) // ERROR "possible misuse of unsafe.Pointer"

}

type V interface {
	Pointer() uintptr
	UnsafeAddr() uintptr
}

type StringHeader struct {
	Data uintptr
}

type SliceHeader struct {
	Data uintptr
}