// antha-tools/antha/ssa/interp/testdata/ifaceprom.go: Part of the Antha language
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

// Test of promotion of methods of an interface embedded within a
// struct.  In particular, this test exercises that the correct
// method is called.

type I interface {
	one() int
	two() string
}

type S struct {
	I
}

type impl struct{}

func (impl) one() int {
	return 1
}

func (impl) two() string {
	return "two"
}

func main() {
	var s S
	s.I = impl{}
	if one := s.I.one(); one != 1 {
		panic(one)
	}
	if one := s.one(); one != 1 {
		panic(one)
	}
	closOne := s.I.one
	if one := closOne(); one != 1 {
		panic(one)
	}
	closOne = s.one
	if one := closOne(); one != 1 {
		panic(one)
	}

	if two := s.I.two(); two != "two" {
		panic(two)
	}
	if two := s.two(); two != "two" {
		panic(two)
	}
	closTwo := s.I.two
	if two := closTwo(); two != "two" {
		panic(two)
	}
	closTwo = s.two
	if two := closTwo(); two != "two" {
		panic(two)
	}
}