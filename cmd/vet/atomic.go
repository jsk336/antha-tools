// antha-tools/cmd/vet/atomic.go: Part of the Antha language
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

import (
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/token"
)

// checkAtomicAssignment walks the assignment statement checking for common
// mistaken usage of atomic package, such as: x = atomic.AddUint64(&x, 1)
func (f *File) checkAtomicAssignment(n *ast.AssignStmt) {
	if !vet("atomic") {
		return
	}

	if len(n.Lhs) != len(n.Rhs) {
		return
	}

	for i, right := range n.Rhs {
		call, ok := right.(*ast.CallExpr)
		if !ok {
			continue
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		pkg, ok := sel.X.(*ast.Ident)
		if !ok || pkg.Name != "atomic" {
			continue
		}

		switch sel.Sel.Name {
		case "AddInt32", "AddInt64", "AddUint32", "AddUint64", "AddUintptr":
			f.checkAtomicAddAssignment(n.Lhs[i], call)
		}
	}
}

// checkAtomicAddAssignment walks the atomic.Add* method calls checking for assigning the return value
// to the same variable being used in the operation
func (f *File) checkAtomicAddAssignment(left ast.Expr, call *ast.CallExpr) {
	if len(call.Args) != 2 {
		return
	}
	arg := call.Args[0]
	broken := false

	if uarg, ok := arg.(*ast.UnaryExpr); ok && uarg.Op == token.AND {
		broken = f.gofmt(left) == f.gofmt(uarg.X)
	} else if star, ok := left.(*ast.StarExpr); ok {
		broken = f.gofmt(star.X) == f.gofmt(arg)
	}

	if broken {
		f.Bad(left.Pos(), "direct assignment to atomic value")
	}
}