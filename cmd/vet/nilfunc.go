// antha-tools/cmd/vet/nilfunc.go: Part of the Antha language
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
This file contains the code to check for useless function comparisons.
A useless comparison is one like f == nil as opposed to f() == nil.
*/

package main

import (
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/token"

	"github.com/antha-lang/antha-tools/antha/types"
)

func (f *File) checkNilFuncComparison(e *ast.BinaryExpr) {
	if !vet("nilfunc") {
		return
	}

	// Only want == or != comparisons.
	if e.Op != token.EQL && e.Op != token.NEQ {
		return
	}

	// Only want comparisons with a nil identifier on one side.
	var e2 ast.Expr
	switch {
	case f.isNil(e.X):
		e2 = e.Y
	case f.isNil(e.Y):
		e2 = e.X
	default:
		return
	}

	// Only want identifiers or selector expressions.
	var obj types.Object
	switch v := e2.(type) {
	case *ast.Ident:
		obj = f.pkg.uses[v]
	case *ast.SelectorExpr:
		obj = f.pkg.uses[v.Sel]
	default:
		return
	}

	// Only want functions.
	if _, ok := obj.(*types.Func); !ok {
		return
	}

	f.Badf(e.Pos(), "comparison of function %v %v nil is always %v", obj.Name(), e.Op, e.Op == token.NEQ)
}

// isNil reports whether the provided expression is the built-in nil
// identifier.
func (f *File) isNil(e ast.Expr) bool {
	return f.pkg.types[e].Type == types.Typ[types.UntypedNil]
}