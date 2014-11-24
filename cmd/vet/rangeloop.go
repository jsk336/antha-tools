// antha-tools/cmd/vet/rangeloop.go: Part of the Antha language
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


/*
This file contains the code to check range loop variables bound inside function
literals that are deferred or launched in new goroutines. We only check
instances where the defer or antha statement is the last statement in the loop
body, as otherwise we would need whole program analysis.

For example:

	for i, v := range s {
		go func() {
			println(i, v) // not what you might expect
		}()
	}

See: http://golang.org/doc/go_faq.html#closures_and_goroutines
*/

package main

import "github.com/antha-lang/antha/ast"

// checkRangeLoop walks the body of the provided range statement, checking if
// its index or value variables are used unsafely inside goroutines or deferred
// function literals.
func checkRangeLoop(f *File, n *ast.RangeStmt) {
	if !vet("rangeloops") {
		return
	}
	key, _ := n.Key.(*ast.Ident)
	val, _ := n.Value.(*ast.Ident)
	if key == nil && val == nil {
		return
	}
	sl := n.Body.List
	if len(sl) == 0 {
		return
	}
	var last *ast.CallExpr
	switch s := sl[len(sl)-1].(type) {
	case *ast.GoStmt:
		last = s.Call
	case *ast.DeferStmt:
		last = s.Call
	default:
		return
	}
	lit, ok := last.Fun.(*ast.FuncLit)
	if !ok {
		return
	}
	ast.Inspect(lit.Body, func(n ast.Node) bool {
		id, ok := n.(*ast.Ident)
		if !ok || id.Obj == nil {
			return true
		}
		if key != nil && id.Obj == key.Obj || val != nil && id.Obj == val.Obj {
			f.Bad(id.Pos(), "range variable", id.Name, "enclosed by function")
		}
		return true
	})
}