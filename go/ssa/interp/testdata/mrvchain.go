// antha-tools/go/ssa/interp/testdata/mrvchain.go: Part of the Antha language
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

// Tests of call chaining f(g()) when g has multiple return values (MRVs).
// See https://code.google.com/p/go/issues/detail?id=4573.

package main

func assert(actual, expected int) {
	if actual != expected {
		panic(actual)
	}
}

func g() (int, int) {
	return 5, 7
}

func g2() (float64, float64) {
	return 5, 7
}

func f1v(x int, v ...int) {
	assert(x, 5)
	assert(v[0], 7)
}

func f2(x, y int) {
	assert(x, 5)
	assert(y, 7)
}

func f2v(x, y int, v ...int) {
	assert(x, 5)
	assert(y, 7)
	assert(len(v), 0)
}

func complexArgs() (float64, float64) {
	return 5, 7
}

func appendArgs() ([]string, string) {
	return []string{"foo"}, "bar"
}

func h() (i interface{}, ok bool) {
	m := map[int]string{1: "hi"}
	i, ok = m[1] // string->interface{} conversion within multi-valued expression
	return
}

func h2() (i interface{}, ok bool) {
	ch := make(chan string, 1)
	ch <- "hi"
	i, ok = <-ch // string->interface{} conversion within multi-valued expression
	return
}

func main() {
	f1v(g())
	f2(g())
	f2v(g())
	if c := complex(complexArgs()); c != 5+7i {
		panic(c)
	}
	if s := append(appendArgs()); len(s) != 2 || s[0] != "foo" || s[1] != "bar" {
		panic(s)
	}
	i, ok := h()
	if !ok || i.(string) != "hi" {
		panic(i)
	}
	i, ok = h2()
	if !ok || i.(string) != "hi" {
		panic(i)
	}
}