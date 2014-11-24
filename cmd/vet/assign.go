// antha-tools/cmd/vet/assign.go: Part of the Antha language
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
This file contains the code to check for useless assignments.
*/

package main

import (
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/token"
	"reflect"
)

// TODO: should also check for assignments to struct fields inside methods
// that are on T instead of *T.

// checkAssignStmt checks for assignments of the form "<expr> = <expr>".
// These are almost always useless, and even when they aren't they are usually a mistake.
func (f *File) checkAssignStmt(stmt *ast.AssignStmt) {
	if !vet("assign") {
		return
	}
	if stmt.Tok != token.ASSIGN {
		return // ignore :=
	}
	if len(stmt.Lhs) != len(stmt.Rhs) {
		// If LHS and RHS have different cardinality, they can't be the same.
		return
	}
	for i, lhs := range stmt.Lhs {
		rhs := stmt.Rhs[i]
		if reflect.TypeOf(lhs) != reflect.TypeOf(rhs) {
			continue // short-circuit the heavy-weight gofmt check
		}
		le := f.gofmt(lhs)
		re := f.gofmt(rhs)
		if le == re {
			f.Badf(stmt.Pos(), "self-assignment of %s to %s", re, le)
		}
	}
}