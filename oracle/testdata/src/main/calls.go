// antha-tools/oracle/testdata/src/main/calls.go: Part of the Antha language
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
// See calls.golden for expected query results.

func A(x *int) { // @pointsto pointsto-A-x "x"
	// @callers callers-A "^"
	// @callstack callstack-A "^"
}

func B(x *int) { // @pointsto pointsto-B-x "x"
	// @callers callers-B "^"
}

// apply is not (yet) treated context-sensitively.
func apply(f func(x *int), x *int) {
	f(x) // @callees callees-apply "f"
	// @callers callers-apply "^"
}

// store *is* treated context-sensitively,
// so the points-to sets for pc, pd are precise.
func store(ptr **int, value *int) {
	*ptr = value
	// @callers callers-store "^"
}

func call(f func() *int) {
	// Result points to anon function.
	f() // @pointsto pointsto-result-f "f"

	// Target of call is anon function.
	f() // @callees callees-main.call-f "f"

	// @callers callers-main.call "^"
}

func main() {
	var a, b int
	apply(A, &a) // @callees callees-main-apply1 "app"
	apply(B, &b)

	var c, d int
	var pc, pd *int // @pointsto pointsto-pc "pc"
	store(&pc, &c)
	store(&pd, &d)
	_ = pd // @pointsto pointsto-pd "pd"

	call(func() *int {
		// We are called twice from main.call
		// @callers callers-main.anon "^"
		return &a
	})

	// Errors
	_ = "no function call here"   // @callees callees-err-no-call "no"
	print("builtin")              // @callees callees-err-builtin "builtin"
	_ = string("type conversion") // @callees callees-err-conversion "str"
	call(nil)                     // @callees callees-err-bad-selection "call\\(nil"
	if false {
		main() // @callees callees-err-deadcode1 "main"
	}
	var nilFunc func()
	nilFunc() // @callees callees-err-nil-func "nilFunc"
	var i interface {
		f()
	}
	i.f() // @callees callees-err-nil-interface "i.f"

	i = new(myint)
	i.f() // @callees callees-not-a-wrapper "f"
}

type myint int

func (myint) f() {
	// @callers callers-not-a-wrapper "^"
}

var dynamic = func() {}

func deadcode() {
	main() // @callees callees-err-deadcode2 "main"
	// @callers callers-err-deadcode "^"
	// @callstack callstack-err-deadcode "^"

	// Within dead code, dynamic calls have no callees.
	dynamic() // @callees callees-err-deadcode3 "dynamic"
}

// This code belongs to init.
var global = 123 // @callers callers-global "global"

// The package initializer may be called by other packages' inits, or
// in this case, the root of the callgraph.  The source-level init functions
// are in turn called by it.
func init() {
	// @callstack callstack-init "^"
}