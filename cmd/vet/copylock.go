// antha-tools/cmd/vet/copylock.go: Part of the Antha language
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


// This file contains the code to check that locks are not passed by value.

package main

import (
	"bytes"
	"fmt"
	"github.com/antha-lang/antha/ast"

	"github.com/antha-lang/antha-tools/antha/types"
)

// checkCopyLocks checks whether a function might
// inadvertently copy a lock, by checking whether
// its receiver, parameters, or return values
// are locks.
func (f *File) checkCopyLocks(d *ast.FuncDecl) {
	if !vet("copylocks") {
		return
	}

	if d.Recv != nil && len(d.Recv.List) > 0 {
		expr := d.Recv.List[0].Type
		if path := lockPath(f.pkg.typesPkg, f.pkg.types[expr].Type); path != nil {
			f.Badf(expr.Pos(), "%s passes Lock by value: %v", d.Name.Name, path)
		}
	}

	if d.Type.Params != nil {
		for _, field := range d.Type.Params.List {
			expr := field.Type
			if path := lockPath(f.pkg.typesPkg, f.pkg.types[expr].Type); path != nil {
				f.Badf(expr.Pos(), "%s passes Lock by value: %v", d.Name.Name, path)
			}
		}
	}

	if d.Type.Results != nil {
		for _, field := range d.Type.Results.List {
			expr := field.Type
			if path := lockPath(f.pkg.typesPkg, f.pkg.types[expr].Type); path != nil {
				f.Badf(expr.Pos(), "%s returns Lock by value: %v", d.Name.Name, path)
			}
		}
	}
}

type typePath []types.Type

// pathString pretty-prints a typePath.
func (path typePath) String() string {
	n := len(path)
	var buf bytes.Buffer
	for i := range path {
		if i > 0 {
			fmt.Fprint(&buf, " contains ")
		}
		// The human-readable path is in reverse order, outermost to innermost.
		fmt.Fprint(&buf, path[n-i-1].String())
	}
	return buf.String()
}

// lockPath returns a typePath describing the location of a lock value
// contained in typ. If there is no contained lock, it returns nil.
func lockPath(tpkg *types.Package, typ types.Type) typePath {
	if typ == nil {
		return nil
	}

	// We're only interested in the case in which the underlying
	// type is a struct. (Interfaces and pointers are safe to copy.)
	styp, ok := typ.Underlying().(*types.Struct)
	if !ok {
		return nil
	}

	// We're looking for cases in which a reference to this type
	// can be locked, but a value cannot. This differentiates
	// embedded interfaces from embedded values.
	if plock := types.NewMethodSet(types.NewPointer(typ)).Lookup(tpkg, "Lock"); plock != nil {
		if lock := types.NewMethodSet(typ).Lookup(tpkg, "Lock"); lock == nil {
			return []types.Type{typ}
		}
	}

	nfields := styp.NumFields()
	for i := 0; i < nfields; i++ {
		ftyp := styp.Field(i).Type()
		subpath := lockPath(tpkg, ftyp)
		if subpath != nil {
			return append(subpath, typ)
		}
	}

	return nil
}