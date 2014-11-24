// antha-tools/oracle/testdata/src/main/callgraph.go: Part of the Antha language
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

// Tests of call-graph queries.
// See go.tools/oracle/oracle_test.go for explanation.
// See callgraph.golden for expected query results.

import "lib"

func A() {}

func B() {}

// call is not (yet) treated context-sensitively.
func call(f func()) {
	f()
}

// nop *is* treated context-sensitively.
func nop() {}

func call2(f func()) {
	f()
	f()
}

func main() {
	call(A)
	call(B)

	nop()
	nop()

	call2(func() {
		// called twice from main.call2,
		// but call2 is not context sensitive (yet).
	})

	print("builtin")
	_ = string("type conversion")
	call(nil)
	if false {
		main()
	}
	var nilFunc func()
	nilFunc()
	var i interface {
		f()
	}
	i.f()

	lib.Func()
}

func deadcode() {
	main()
}

// @callgraph callgraph-main "^"

// @callgraph callgraph-complete "nopos"