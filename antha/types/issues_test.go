// antha-tools/antha/types/issues_test.go: Part of the Antha language
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


// This file implements tests for various issues.

package types_test

import (
	"fmt"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/parser"
	"sort"
	"strings"
	"testing"

	_ "github.com/antha-lang/antha-tools/antha/gcimporter"
	. "github.com/antha-lang/antha-tools/antha/types"
)

func TestIssue5770(t *testing.T) {
	src := `package p; type S struct{T}`
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		t.Fatal(err)
	}

	_, err = Check(f.Name.Name, fset, []*ast.File{f}) // do not crash
	want := "undeclared name: T"
	if err == nil || !strings.Contains(err.Error(), want) {
		t.Errorf("got: %v; want: %s", err, want)
	}
}

func TestIssue5849(t *testing.T) {
	src := `
package p
var (
	s uint
	_ = uint8(8)
	_ = uint16(16) << s
	_ = uint32(32 << s)
	_ = uint64(64 << s + s)
	_ = (interface{})("foo")
	_ = (interface{})(nil)
)`
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		t.Fatal(err)
	}

	var conf Config
	types := make(map[ast.Expr]TypeAndValue)
	_, err = conf.Check(f.Name.Name, fset, []*ast.File{f}, &Info{Types: types})
	if err != nil {
		t.Fatal(err)
	}

	for x, tv := range types {
		var want Type
		switch x := x.(type) {
		case *ast.BasicLit:
			switch x.Value {
			case `8`:
				want = Typ[Uint8]
			case `16`:
				want = Typ[Uint16]
			case `32`:
				want = Typ[Uint32]
			case `64`:
				want = Typ[Uint] // because of "+ s", s is of type uint
			case `"foo"`:
				want = Typ[String]
			}
		case *ast.Ident:
			if x.Name == "nil" {
				want = Typ[UntypedNil]
			}
		}
		if want != nil && !Identical(tv.Type, want) {
			t.Errorf("got %s; want %s", tv.Type, want)
		}
	}
}

func TestIssue6413(t *testing.T) {
	src := `
package p
func f() int {
	defer f()
	go f()
	return 0
}
`
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		t.Fatal(err)
	}

	var conf Config
	types := make(map[ast.Expr]TypeAndValue)
	_, err = conf.Check(f.Name.Name, fset, []*ast.File{f}, &Info{Types: types})
	if err != nil {
		t.Fatal(err)
	}

	want := Typ[Int]
	n := 0
	for x, tv := range types {
		if _, ok := x.(*ast.CallExpr); ok {
			if tv.Type != want {
				t.Errorf("%s: got %s; want %s", fset.Position(x.Pos()), tv.Type, want)
			}
			n++
		}
	}

	if n != 2 {
		t.Errorf("got %d CallExprs; want 2", n)
	}
}

func TestIssue7245(t *testing.T) {
	src := `
package p
func (T) m() (res bool) { return }
type T struct{} // receiver type after method declaration
`
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		t.Fatal(err)
	}

	var conf Config
	defs := make(map[*ast.Ident]Object)
	_, err = conf.Check(f.Name.Name, fset, []*ast.File{f}, &Info{Defs: defs})
	if err != nil {
		t.Fatal(err)
	}

	m := f.Decls[0].(*ast.FuncDecl)
	res1 := defs[m.Name].(*Func).Type().(*Signature).Results().At(0)
	res2 := defs[m.Type.Results.List[0].Names[0]].(*Var)

	if res1 != res2 {
		t.Errorf("got %s (%p) != %s (%p)", res1, res2, res1, res2)
	}
}

// This tests that uses of existing vars on the LHS of an assignment
// are Uses, not Defs; and also that the (illegal) use of a non-var on
// the LHS of an assignment is a Use nonetheless.
func TestIssue7827(t *testing.T) {
	const src = `
package p
func _() {
	const w = 1        // defs w
        x, y := 2, 3       // defs x, y
        w, x, z := 4, 5, 6 // uses w, x, defs z; error: cannot assign to w
        _, _, _ = x, y, z  // uses x, y, z
}
`
	const want = `L3 defs func p._()
L4 defs const w untyped int
L5 defs var x int
L5 defs var y int
L6 defs var z int
L6 uses const w untyped int
L6 uses var x int
L7 uses var x int
L7 uses var y int
L7 uses var z int`

	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		t.Fatal(err)
	}

	// don't abort at the first error
	conf := Config{Error: func(err error) { t.Log(err) }}
	defs := make(map[*ast.Ident]Object)
	uses := make(map[*ast.Ident]Object)
	_, err = conf.Check(f.Name.Name, fset, []*ast.File{f}, &Info{Defs: defs, Uses: uses})
	if s := fmt.Sprint(err); !strings.HasSuffix(s, "cannot assign to w") {
		t.Errorf("Check: unexpected error: %s", s)
	}

	var facts []string
	for id, obj := range defs {
		if obj != nil {
			fact := fmt.Sprintf("L%d defs %s", fset.Position(id.Pos()).Line, obj)
			facts = append(facts, fact)
		}
	}
	for id, obj := range uses {
		fact := fmt.Sprintf("L%d uses %s", fset.Position(id.Pos()).Line, obj)
		facts = append(facts, fact)
	}
	sort.Strings(facts)

	got := strings.Join(facts, "\n")
	if got != want {
		t.Errorf("Unexpected defs/uses\ngot:\n%s\nwant:\n%s", got, want)
	}
}