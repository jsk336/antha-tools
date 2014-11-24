// antha-tools/astutil/enclosing_test.go: Part of the Antha language
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


package astutil_test

// This file defines tests of PathEnclosingInterval.

// TODO(adonovan): exhaustive tests that run over the whole input
// tree, not just handcrafted examples.

import (
	"bytes"
	"fmt"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/parser"
	"github.com/antha-lang/antha/token"
	"strings"
	"testing"

	"code.google.com/p/go.tools/astutil"
)

// pathToString returns a string containing the concrete types of the
// nodes in path.
func pathToString(path []ast.Node) string {
	var buf bytes.Buffer
	fmt.Fprint(&buf, "[")
	for i, n := range path {
		if i > 0 {
			fmt.Fprint(&buf, " ")
		}
		fmt.Fprint(&buf, strings.TrimPrefix(fmt.Sprintf("%T", n), "*ast."))
	}
	fmt.Fprint(&buf, "]")
	return buf.String()
}

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

// Common input for following tests.
const input = `
// Hello.
package main
import "fmt"
func f() {}
func main() {
	z := (x + y) // add them
        f() // NB: ExprStmt and its CallExpr have same Pos/End
}
`

func TestPathEnclosingInterval_Exact(t *testing.T) {
	// For the exact tests, we check that a substring is mapped to
	// the canonical string for the node it denotes.
	tests := []struct {
		substr string // first occurrence of this string indicates interval
		node   string // complete text of expected containing node
	}{
		{"package",
			input[11 : len(input)-1]},
		{"\npack",
			input[11 : len(input)-1]},
		{"main",
			"main"},
		{"import",
			"import \"fmt\""},
		{"\"fmt\"",
			"\"fmt\""},
		{"\nfunc f() {}\n",
			"func f() {}"},
		{"x ",
			"x"},
		{" y",
			"y"},
		{"z",
			"z"},
		{" + ",
			"x + y"},
		{" :=",
			"z := (x + y)"},
		{"x + y",
			"x + y"},
		{"(x + y)",
			"(x + y)"},
		{" (x + y) ",
			"(x + y)"},
		{" (x + y) // add",
			"(x + y)"},
		{"func",
			"func f() {}"},
		{"func f() {}",
			"func f() {}"},
		{"\nfun",
			"func f() {}"},
		{" f",
			"f"},
	}
	for _, test := range tests {
		f, start, end := findInterval(t, new(token.FileSet), input, test.substr)
		if f == nil {
			continue
		}

		path, exact := astutil.PathEnclosingInterval(f, start, end)
		if !exact {
			t.Errorf("PathEnclosingInterval(%q) not exact", test.substr)
			continue
		}

		if len(path) == 0 {
			if test.node != "" {
				t.Errorf("PathEnclosingInterval(%q).path: got [], want %q",
					test.substr, test.node)
			}
			continue
		}

		if got := input[path[0].Pos():path[0].End()]; got != test.node {
			t.Errorf("PathEnclosingInterval(%q): got %q, want %q (path was %s)",
				test.substr, got, test.node, pathToString(path))
			continue
		}
	}
}

func TestPathEnclosingInterval_Paths(t *testing.T) {
	// For these tests, we check only the path of the enclosing
	// node, but not its complete text because it's often quite
	// large when !exact.
	tests := []struct {
		substr string // first occurrence of this string indicates interval
		path   string // the pathToString(),exact of the expected path
	}{
		{"// add",
			"[BlockStmt FuncDecl File],false"},
		{"(x + y",
			"[ParenExpr AssignStmt BlockStmt FuncDecl File],false"},
		{"x +",
			"[BinaryExpr ParenExpr AssignStmt BlockStmt FuncDecl File],false"},
		{"z := (x",
			"[AssignStmt BlockStmt FuncDecl File],false"},
		{"func f",
			"[FuncDecl File],false"},
		{"func f()",
			"[FuncDecl File],false"},
		{" f()",
			"[FuncDecl File],false"},
		{"() {}",
			"[FuncDecl File],false"},
		{"// Hello",
			"[File],false"},
		{" f",
			"[Ident FuncDecl File],true"},
		{"func ",
			"[FuncDecl File],true"},
		{"mai",
			"[Ident File],true"},
		{"f() // NB",
			"[CallExpr ExprStmt BlockStmt FuncDecl File],true"},
	}
	for _, test := range tests {
		f, start, end := findInterval(t, new(token.FileSet), input, test.substr)
		if f == nil {
			continue
		}

		path, exact := astutil.PathEnclosingInterval(f, start, end)
		if got := fmt.Sprintf("%s,%v", pathToString(path), exact); got != test.path {
			t.Errorf("PathEnclosingInterval(%q): got %q, want %q",
				test.substr, got, test.path)
			continue
		}
	}
}