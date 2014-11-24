// antha-tools/antha/ssa/testdata/valueforexpr.go: Part of the Antha language
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

//+build ignore

package main

// This file is the input to TestValueForExpr in source_test.go, which
// ensures that each expression e immediately following a /*@kind*/(x)
// annotation, when passed to Function.ValueForExpr(e), returns a
// non-nil Value of the same type as e and of kind 'kind'.

func f(spilled, unspilled int) {
	_ = /*@UnOp*/ (spilled)
	_ = /*@Parameter*/ (unspilled)
	_ = /*@<nil>*/ (1 + 2) // (constant)
	i := 0
	/*@Call*/ (print( /*@BinOp*/ (i + 1)))
	ch := /*@MakeChan*/ (make(chan int))
	/*@UnOp*/ (<-ch)
	x := /*@UnOp*/ (<-ch)
	_ = x
	select {
	case /*@Extract*/ (<-ch):
	case x := /*@Extract*/ (<-ch):
		_ = x
	}
	defer /*@Function*/ (func() {
	})()
	go /*@Function*/ (func() {
	})()
	y := 0
	if true && /*@BinOp*/ (bool(y > 0)) {
		y = 1
	}
	_ = /*@Phi*/ (y)
	map1 := /*@MakeMap*/ (make(map[string]string))
	_ = map1
	_ = /*@MakeSlice*/ (make([]int, 0))
	_ = /*@MakeClosure*/ (func() { print(spilled) })

	sl := []int{}
	_ = /*@Slice*/ (sl[:0])

	_ = /*@<nil>*/ (new(int)) // optimized away
	tmp := /*@Alloc*/ (new(int))
	_ = tmp
	var iface interface{}
	_ = /*@TypeAssert*/ (iface.(int))
	_ = /*@UnOp*/ (sl[0])
	_ = /*@IndexAddr*/ (&sl[0])
	_ = /*@Index*/ ([2]int{}[0])
	var p *int
	_ = /*@UnOp*/ (*p)

	_ = /*@UnOp*/ (global)
	/*@UnOp*/ (global)[""] = ""
	/*@Global*/ (global) = map[string]string{}

	var local t
	/*UnOp*/ (local.x) = 1

	// Exercise corner-cases of lvalues vs rvalues.
	type N *N
	var n N
	/*@UnOp*/ (n) = /*@UnOp*/ (n)
	/*@ChangeType*/ (n) = /*@Alloc*/ (&n)
	/*@UnOp*/ (n) = /*@UnOp*/ (*n)
	/*@UnOp*/ (n) = /*@UnOp*/ (**n)
}

func complit() {
	// Composite literals.
	// We get different results for
	// - composite literal as value (e.g. operand to print)
	// - composite literal initializer for addressable value
	// - composite literal value assigned to blank var

	// 1. Slices
	print( /*@Slice*/ ([]int{}))
	print( /*@Alloc*/ (&[]int{}))
	print(& /*@Alloc*/ ([]int{}))

	sl1 := /*@Slice*/ ([]int{})
	sl2 := /*@Alloc*/ (&[]int{})
	sl3 := & /*@Alloc*/ ([]int{})
	_, _, _ = sl1, sl2, sl3

	_ = /*@Slice*/ ([]int{})
	_ = /*@<nil>*/ (& /*@Slice*/ ([]int{})) // & optimized away
	_ = & /*@Slice*/ ([]int{})

	// 2. Arrays
	print( /*@UnOp*/ ([1]int{}))
	print( /*@Alloc*/ (&[1]int{}))
	print(& /*@Alloc*/ ([1]int{}))

	arr1 := /*@Alloc*/ ([1]int{})
	arr2 := /*@Alloc*/ (&[1]int{})
	arr3 := & /*@Alloc*/ ([1]int{})
	_, _, _ = arr1, arr2, arr3

	_ = /*@UnOp*/ ([1]int{})
	_ = /*@Alloc*/ (& /*@Alloc*/ ([1]int{})) // & optimized away
	_ = & /*@Alloc*/ ([1]int{})

	// 3. Maps
	type M map[int]int
	print( /*@MakeMap*/ (M{}))
	print( /*@Alloc*/ (&M{}))
	print(& /*@Alloc*/ (M{}))

	m1 := /*@MakeMap*/ (M{})
	m2 := /*@Alloc*/ (&M{})
	m3 := & /*@Alloc*/ (M{})
	_, _, _ = m1, m2, m3

	_ = /*@MakeMap*/ (M{})
	_ = /*@<nil>*/ (& /*@MakeMap*/ (M{})) // & optimized away
	_ = & /*@MakeMap*/ (M{})

	// 4. Structs
	print( /*@UnOp*/ (struct{}{}))
	print( /*@Alloc*/ (&struct{}{}))
	print(& /*@Alloc*/ (struct{}{}))

	s1 := /*@Alloc*/ (struct{}{})
	s2 := /*@Alloc*/ (&struct{}{})
	s3 := & /*@Alloc*/ (struct{}{})
	_, _, _ = s1, s2, s3

	_ = /*@UnOp*/ (struct{}{})
	_ = /*@Alloc*/ (& /*@Alloc*/ (struct{}{}))
	_ = & /*@Alloc*/ (struct{}{})
}

type t struct{ x int }

// Ensure we can locate methods of named types.
func (t) f(param int) {
	_ = /*@Parameter*/ (param)
}

// Ensure we can locate init functions.
func init() {
	m := /*@MakeMap*/ (make(map[string]string))
	_ = m
}

// Ensure we can locate variables in initializer expressions.
var global = /*@MakeMap*/ (make(map[string]string))