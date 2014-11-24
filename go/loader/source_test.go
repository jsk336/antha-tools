// antha-tools/go/loader/source_test.go: Part of the Antha language
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


package loader_test

// This file defines tests of source utilities.

import (
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/parser"
	"github.com/antha-lang/antha/token"
	"strings"
	"testing"

	"github.com/antha-lang/antha-tools/astutil"
	"github.com/antha-lang/antha-tools/antha/loader"
	"github.com/antha-lang/antha-tools/antha/ssa"
)

// findInterval parses input and returns the [start, end) positions of
// the first occurrence of substr in input.  f==nil indicates failure;
// an error has already been reported in that case.
//
func findInterval(t *testing.T, fset *token.FileSet, input, substr string) (f *ast.File, start, end token.Pos) {
	f, err := parser.ParseFile(fset, "<input>", input, 0)
	if err != nil {
		t.Errorf("parse error: %s", err)
		return
	}

	i := strings.Index(input, substr)
	if i < 0 {
		t.Errorf("%q is not a substring of input", substr)
		f = nil
		return
	}

	filePos := fset.File(f.Package)
	return f, filePos.Pos(i), filePos.Pos(i + len(substr))
}

func TestEnclosingFunction(t *testing.T) {
	tests := []struct {
		input  string // the input file
		substr string // first occurrence of this string denotes interval
		fn     string // name of expected containing function
	}{
		// We use distinctive numbers as syntactic landmarks.

		// Ordinary function:
		{`package main
		  func f() { println(1003) }`,
			"100", "main.f"},
		// Methods:
		{`package main
                  type T int
		  func (t T) f() { println(200) }`,
			"200", "(main.T).f"},
		// Function literal:
		{`package main
		  func f() { println(func() { print(300) }) }`,
			"300", "f$1"},
		// Doubly nested
		{`package main
		  func f() { println(func() { print(func() { print(350) })})}`,
			"350", "f$1$1"},
		// Implicit init for package-level var initializer.
		{"package main; var a = 400", "400", "main.init"},
		// No code for constants:
		{"package main; const a = 500", "500", "(none)"},
		// Explicit init()
		{"package main; func init() { println(600) }", "600", "main.init$1"},
		// Multiple explicit init functions:
		{`package main
		  func init() { println("foo") }
		  func init() { println(800) }`,
			"800", "main.init$2"},
		// init() containing FuncLit.
		{`package main
		  func init() { println(func(){print(900)}) }`,
			"900", "init$1$1"},
	}
	for _, test := range tests {
		conf := loader.Config{Fset: token.NewFileSet()}
		f, start, end := findInterval(t, conf.Fset, test.input, test.substr)
		if f == nil {
			continue
		}
		path, exact := astutil.PathEnclosingInterval(f, start, end)
		if !exact {
			t.Errorf("EnclosingFunction(%q) not exact", test.substr)
			continue
		}

		conf.CreateFromFiles("main", f)

		iprog, err := conf.Load()
		if err != nil {
			t.Error(err)
			continue
		}
		prog := ssa.Create(iprog, 0)
		pkg := prog.Package(iprog.Created[0].Pkg)
		pkg.Build()

		name := "(none)"
		fn := ssa.EnclosingFunction(pkg, path)
		if fn != nil {
			name = fn.String()
		}

		if name != test.fn {
			t.Errorf("EnclosingFunction(%q in %q) got %s, want %s",
				test.substr, test.input, name, test.fn)
			continue
		}

		// While we're here: test HasEnclosingFunction.
		if has := ssa.HasEnclosingFunction(pkg, path); has != (fn != nil) {
			t.Errorf("HasEnclosingFunction(%q in %q) got %v, want %v",
				test.substr, test.input, has, fn != nil)
			continue
		}
	}
}