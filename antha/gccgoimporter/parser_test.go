// antha-tools/antha/gccgoimporter/parser_test.go: Part of the Antha language
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


package gccgoimporter

import (
	"bytes"
	"strings"
	"testing"
	"text/scanner"

	"github.com/antha-lang/antha-tools/antha/types"
)

var typeParserTests = []struct {
	id, typ, want, underlying, methods string
}{
	{id: "foo", typ: "<type -1>", want: "int8"},
	{id: "foo", typ: "<type 1 *<type -19>>", want: "*error"},
	{id: "foo", typ: "<type 1 *any>", want: "unsafe.Pointer"},
	{id: "foo", typ: "<type 1 \"Bar\" <type 2 *<type 1>>>", want: "foo.Bar", underlying: "*foo.Bar"},
	{id: "foo", typ: "<type 1 \"bar.Foo\" \"bar\" <type -1> func (? <type 1>) M (); >", want: "bar.Foo", underlying: "int8", methods: "func (bar.Foo).M()"},
	{id: "foo", typ: "<type 1 \".bar.foo\" \"bar\" <type -1>>", want: "bar.foo", underlying: "int8"},
	{id: "foo", typ: "<type 1 []<type -1>>", want: "[]int8"},
	{id: "foo", typ: "<type 1 [42]<type -1>>", want: "[42]int8"},
	{id: "foo", typ: "<type 1 map [<type -1>] <type -2>>", want: "map[int8]int16"},
	{id: "foo", typ: "<type 1 chan <type -1>>", want: "chan int8"},
	{id: "foo", typ: "<type 1 chan <- <type -1>>", want: "<-chan int8"},
	{id: "foo", typ: "<type 1 chan -< <type -1>>", want: "chan<- int8"},
	{id: "foo", typ: "<type 1 struct { I8 <type -1>; I16 <type -2> \"i16\"; }>", want: "struct{I8 int8; I16 int16 \"i16\"}"},
	{id: "foo", typ: "<type 1 interface { Foo (a <type -1>, b <type -2>) <type -1>; Bar (? <type -2>, ? ...<type -1>) (? <type -2>, ? <type -1>); Baz (); }>", want: "interface{Bar(int16, ...int8) (int16, int8); Baz(); Foo(a int8, b int16) int8}"},
	{id: "foo", typ: "<type 1 (? <type -1>) <type -2>>", want: "func(int8) int16"},
}

func TestTypeParser(t *testing.T) {
	for _, test := range typeParserTests {
		var p parser
		p.init("test.gox", strings.NewReader(test.typ), make(map[string]*types.Package))
		p.pkgname = test.id
		p.pkgpath = test.id
		p.maybeCreatePackage()
		typ := p.parseType(p.pkg)

		if p.tok != scanner.EOF {
			t.Errorf("expected full parse, stopped at %q", p.lit)
		}

		got := typ.String()
		if got != test.want {
			t.Errorf("got type %q, expected %q", got, test.want)
		}

		if test.underlying != "" {
			underlying := typ.Underlying().String()
			if underlying != test.underlying {
				t.Errorf("got underlying type %q, expected %q", underlying, test.underlying)
			}
		}

		if test.methods != "" {
			nt := typ.(*types.Named)
			var buf bytes.Buffer
			for i := 0; i != nt.NumMethods(); i++ {
				buf.WriteString(nt.Method(i).String())
			}
			methods := buf.String()
			if methods != test.methods {
				t.Errorf("got methods %q, expected %q", methods, test.methods)
			}
		}
	}
}