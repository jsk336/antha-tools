// antha-tools/oracle/testdata/src/main/pointsto.go: Part of the Antha language
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

package pointsto

// Tests of 'pointsto' query.
// See go.tools/oracle/oracle_test.go for explanation.
// See pointsto.golden for expected query results.

const pi = 3.141 // @pointsto const "pi"

var global = new(string) // NB: ssa.Global is indirect, i.e. **string

func main() {
	livecode()

	// func objects
	_ = main   // @pointsto func-ref-main "main"
	_ = (*C).f // @pointsto func-ref-*C.f "..C..f"
	_ = D.f    // @pointsto func-ref-D.f "D.f"
	_ = I.f    // @pointsto func-ref-I.f "I.f"
	var d D
	var i I
	_ = d.f // @pointsto func-ref-d.f "d.f"
	_ = i.f // @pointsto func-ref-i.f "i.f"

	// var objects
	anon := func() {
		_ = d.f // @pointsto ref-lexical-d.f "d.f"
	}
	_ = anon   // @pointsto ref-anon "anon"
	_ = global // @pointsto ref-global "global"

	// SSA affords some local flow sensitivity.
	var a, b int
	var x = &a // @pointsto var-def-x-1 "x"
	_ = x      // @pointsto var-ref-x-1 "x"
	x = &b     // @pointsto var-def-x-2 "x"
	_ = x      // @pointsto var-ref-x-2 "x"

	i = new(C) // @pointsto var-ref-i-C "i"
	if i != nil {
		i = D{} // @pointsto var-ref-i-D "i"
	}
	print(i) // @pointsto var-ref-i "\\bi\\b"

	m := map[string]*int{"a": &a}
	mapval, _ := m["a"] // @pointsto map-lookup,ok "m..a.."
	_ = mapval          // @pointsto mapval "mapval"
	_ = m               // @pointsto m "m"

	panic(3) // @pointsto builtin-panic "panic"
}

func livecode() {} // @pointsto func-live "livecode"

func deadcode() { // @pointsto func-dead "deadcode"
	// Pointer analysis can't run on dead code.
	var b = new(int) // @pointsto b "b"
	_ = b
}

type I interface {
	f()
}

type C int
type D struct{}

func (c *C) f() {}
func (d D) f()  {}