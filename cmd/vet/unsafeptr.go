// antha-tools/cmd/vet/unsafeptr.go: Part of the Antha language
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


// Check for invalid uintptr -> unsafe.Pointer conversions.

package main

import (
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/token"

	"github.com/antha-lang/antha-tools/antha/types"
)

func (f *File) checkUnsafePointer(x *ast.CallExpr) {
	if !vet("unsafeptr") {
		return
	}
	if len(x.Args) != 1 {
		return
	}
	if f.hasBasicType(x.Fun, types.UnsafePointer) && f.hasBasicType(x.Args[0], types.Uintptr) && !f.isSafeUintptr(x.Args[0]) {
		f.Badf(x.Pos(), "possible misuse of unsafe.Pointer")
	}
}

// isSafeUintptr reports whether x - already known to be a uintptr -
// is safe to convert to unsafe.Pointer. It is safe if x is itself derived
// directly from an unsafe.Pointer via conversion and pointer arithmetic
// or if x is the result of reflect.Value.Pointer or reflect.Value.UnsafeAddr
// or obtained from the Data field of a *reflect.SliceHeader or *reflect.StringHeader.
func (f *File) isSafeUintptr(x ast.Expr) bool {
	switch x := x.(type) {
	case *ast.ParenExpr:
		return f.isSafeUintptr(x.X)

	case *ast.SelectorExpr:
		switch x.Sel.Name {
		case "Data":
			// reflect.SliceHeader and reflect.StringHeader are okay,
			// but only if they are pointing at a real slice or string.
			// It's not okay to do:
			//	var x SliceHeader
			//	x.Data = uintptr(unsafe.Pointer(...))
			//	... use x ...
			//	p := unsafe.Pointer(x.Data)
			// because in the middle the garbage collector doesn't
			// see x.Data as a pointer and so x.Data may be dangling
			// by the time we get to the conversion at the end.
			// For now approximate by saying that *Header is okay
			// but Header is not.
			pt, ok := f.pkg.types[x.X].Type.(*types.Pointer)
			if ok {
				t, ok := pt.Elem().(*types.Named)
				if ok && t.Obj().Pkg().Path() == "reflect" {
					switch t.Obj().Name() {
					case "StringHeader", "SliceHeader":
						return true
					}
				}
			}
		}

	case *ast.CallExpr:
		switch len(x.Args) {
		case 0:
			// maybe call to reflect.Value.Pointer or reflect.Value.UnsafeAddr.
			sel, ok := x.Fun.(*ast.SelectorExpr)
			if !ok {
				break
			}
			switch sel.Sel.Name {
			case "Pointer", "UnsafeAddr":
				t, ok := f.pkg.types[sel.X].Type.(*types.Named)
				if ok && t.Obj().Pkg().Path() == "reflect" && t.Obj().Name() == "Value" {
					return true
				}
			}

		case 1:
			// maybe conversion of uintptr to unsafe.Pointer
			return f.hasBasicType(x.Fun, types.Uintptr) && f.hasBasicType(x.Args[0], types.UnsafePointer)
		}

	case *ast.BinaryExpr:
		switch x.Op {
		case token.ADD, token.SUB:
			return f.isSafeUintptr(x.X) && !f.isSafeUintptr(x.Y)
		}
	}
	return false
}