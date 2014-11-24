// antha-tools/oracle/testdata/src/main/describe.go: Part of the Antha language
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

package describe // @describe pkgdecl "describe"

// Tests of 'describe' query.
// See go.tools/oracle/oracle_test.go for explanation.
// See describe.golden for expected query results.

// TODO(adonovan): more coverage of the (extensive) logic.

type cake float64 // @describe type-ref-builtin "float64"

const c = iota // @describe const-ref-iota "iota"

const pi = 3.141     // @describe const-def-pi "pi"
const pie = cake(pi) // @describe const-def-pie "pie"
const _ = pi         // @describe const-ref-pi "pi"

var global = new(string) // NB: ssa.Global is indirect, i.e. **string

func main() { // @describe func-def-main "main"
	// func objects
	_ = main   // @describe func-ref-main "main"
	_ = (*C).f // @describe func-ref-*C.f "..C..f"
	_ = D.f    // @describe func-ref-D.f "D.f"
	_ = I.f    // @describe func-ref-I.f "I.f"
	var d D    // @describe type-D "D"
	var i I    // @describe type-I "I"
	_ = d.f    // @describe func-ref-d.f "d.f"
	_ = i.f    // @describe func-ref-i.f "i.f"

	// var objects
	anon := func() {
		_ = d // @describe ref-lexical-d "d"
	}
	_ = anon   // @describe ref-anon "anon"
	_ = global // @describe ref-global "global"

	// SSA affords some local flow sensitivity.
	var a, b int
	var x = &a // @describe var-def-x-1 "x"
	_ = x      // @describe var-ref-x-1 "x"
	x = &b     // @describe var-def-x-2 "x"
	_ = x      // @describe var-ref-x-2 "x"

	i = new(C) // @describe var-ref-i-C "i"
	if i != nil {
		i = D{} // @describe var-ref-i-D "i"
	}
	print(i) // @describe var-ref-i "\\bi\\b"

	// const objects
	const localpi = 3.141     // @describe const-local-pi "localpi"
	const localpie = cake(pi) // @describe const-local-pie "localpie"
	const _ = localpi         // @describe const-ref-localpi "localpi"

	// type objects
	type T int      // @describe type-def-T "T"
	var three T = 3 // @describe type-ref-T "T"
	_ = three

	print(1 + 2*3)        // @describe const-expr " 2.3"
	print(real(1+2i) - 3) // @describe const-expr2 "real.*3"

	m := map[string]*int{"a": &a}
	mapval, _ := m["a"] // @describe map-lookup,ok "m..a.."
	_ = mapval          // @describe mapval "mapval"
	_ = m               // @describe m "m"

	defer main() // @describe defer-stmt "defer"
	go main()    // @describe go-stmt "go"

	panic(3) // @describe builtin-ref-panic "panic"

	var a2 int // @describe var-decl-stmt "var a2 int"
	_ = a2
	var _ int // @describe var-decl-stmt2 "var _ int"
	var _ int // @describe var-def-blank "_"
}

type I interface { // @describe def-iface-I "I"
	f() // @describe def-imethod-I.f "f"
}

type C int
type D struct{}

func (c *C) f() {}
func (d D) f()  {}