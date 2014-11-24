// antha-tools/go/ssa/interp/testdata/boundmeth.go: Part of the Antha language
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

// Tests of bound method closures.

package main

import "fmt"

func assert(b bool) {
	if !b {
		panic("oops")
	}
}

type I int

func (i I) add(x int) int {
	return int(i) + x
}

func valueReceiver() {
	var three I = 3
	assert(three.add(5) == 8)
	var add3 func(int) int = three.add
	assert(add3(5) == 8)
}

type S struct{ x int }

func (s *S) incr() {
	s.x++
}

func (s *S) get() int {
	return s.x
}

func pointerReceiver() {
	ps := new(S)
	incr := ps.incr
	get := ps.get
	assert(get() == 0)
	incr()
	incr()
	incr()
	assert(get() == 3)
}

func addressibleValuePointerReceiver() {
	var s S
	incr := s.incr
	get := s.get
	assert(get() == 0)
	incr()
	incr()
	incr()
	assert(get() == 3)
}

type S2 struct {
	S
}

func promotedReceiver() {
	var s2 S2
	incr := s2.incr
	get := s2.get
	assert(get() == 0)
	incr()
	incr()
	incr()
	assert(get() == 3)
}

func anonStruct() {
	var s struct{ S }
	incr := s.incr
	get := s.get
	assert(get() == 0)
	incr()
	incr()
	incr()
	assert(get() == 3)
}

func typeCheck() {
	var i interface{}
	i = (*S).incr
	_ = i.(func(*S)) // type assertion: receiver type prepended to params

	var s S
	i = s.incr
	_ = i.(func()) // type assertion: receiver type disappears
}

type errString string

func (err errString) Error() string {
	return string(err)
}

// Regression test for a builder crash.
func regress1(x error) func() string {
	return x.Error
}

// Regression test for b/7269:
// taking the value of an interface method performs a nil check.
func nilInterfaceMethodValue() {
	err := fmt.Errorf("ok")
	f := err.Error
	if got := f(); got != "ok" {
		panic(got)
	}

	err = nil
	if got := f(); got != "ok" {
		panic(got)
	}

	defer func() {
		r := fmt.Sprint(recover())
		// runtime panic string varies across toolchains
		if r != "runtime error: interface conversion: interface is nil, not error" &&
			r != "runtime error: invalid memory address or nil pointer dereference" {
			panic("want runtime panic from nil interface method value, got " + r)
		}
	}()
	f = err.Error // runtime panic: err is nil
	panic("unreachable")
}

func main() {
	valueReceiver()
	pointerReceiver()
	addressibleValuePointerReceiver()
	promotedReceiver()
	anonStruct()
	typeCheck()

	if e := regress1(errString("hi"))(); e != "hi" {
		panic(e)
	}

	nilInterfaceMethodValue()
}