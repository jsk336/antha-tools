// antha-tools/go/pointer/testdata/maps.go: Part of the Antha language
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

// Test of maps.

var a, b, c int

func maps1() {
	m1 := map[*int]*int{&a: &b} // @line m1m1
	m2 := make(map[*int]*int)   // @line m1m2
	m2[&b] = &a

	print(m1[nil]) // @pointsto main.b | main.c
	print(m2[nil]) // @pointsto main.a

	print(m1) // @pointsto makemap@m1m1:21
	print(m2) // @pointsto makemap@m1m2:12

	m1[&b] = &c

	for k, v := range m1 {
		print(k) // @pointsto main.a | main.b
		print(v) // @pointsto main.b | main.c
	}

	for k, v := range m2 {
		print(k) // @pointsto main.b
		print(v) // @pointsto main.a
	}

	// Lookup doesn't create any aliases.
	print(m2[&c]) // @pointsto main.a
	if _, ok := m2[&a]; ok {
		print(m2[&c]) // @pointsto main.a
	}
}

func maps2() {
	m1 := map[*int]*int{&a: &b}
	m2 := map[*int]*int{&b: &c}
	_ = []map[*int]*int{m1, m2} // (no spurious merging of m1, m2)

	print(m1[nil]) // @pointsto main.b
	print(m2[nil]) // @pointsto main.c
}

func main() {
	maps1()
	maps2()
}