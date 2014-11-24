// antha-tools/go/loader/pkginfo.go: Part of the Antha language
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


package loader

import (
	"fmt"
	"github.com/antha-lang/antha/ast"

	"github.com/antha-lang/antha-tools/antha/exact"
	"github.com/antha-lang/antha-tools/antha/types"
)

// PackageInfo holds the ASTs and facts derived by the type-checker
// for a single package.
//
// Not mutated once exposed via the API.
//
type PackageInfo struct {
	Pkg                   *types.Package
	Importable            bool        // true if 'import "Pkg.Path()"' would resolve to this
	TransitivelyErrorFree bool        // true if Pkg and all its dependencies are free of errors
	Files                 []*ast.File // abstract syntax for the package's files
	TypeError             error       // non-nil if the package had type errors
	types.Info                        // type-checker deductions.

	checker interface {
		Files(files []*ast.File) error
	} // transient type-checker state
}

func (info *PackageInfo) String() string {
	return fmt.Sprintf("PackageInfo(%s)", info.Pkg.Path())
}

// TypeOf returns the type of expression e.
// Precondition: e belongs to the package's ASTs.
//
func (info *PackageInfo) TypeOf(e ast.Expr) types.Type {
	if t, ok := info.Types[e]; ok {
		return t.Type
	}
	// Defining ast.Idents (id := expr) get only Ident callbacks
	// but not Expr callbacks.
	if id, ok := e.(*ast.Ident); ok {
		return info.ObjectOf(id).Type()
	}
	panic("no type for expression")
}

// ValueOf returns the value of expression e if it is a constant, nil
// otherwise.
// Precondition: e belongs to the package's ASTs.
//
func (info *PackageInfo) ValueOf(e ast.Expr) exact.Value {
	return info.Types[e].Value
}

// ObjectOf returns the typechecker object denoted by the specified id.
//
// If id is an anonymous struct field, the field (*types.Var) is
// returned, not the type (*types.TypeName).
//
// Precondition: id belongs to the package's ASTs.
//
func (info *PackageInfo) ObjectOf(id *ast.Ident) types.Object {
	obj, ok := info.Defs[id]
	if ok {
		return obj
	}
	return info.Uses[id]
}

// IsType returns true iff expression e denotes a type.
// Precondition: e belongs to the package's ASTs.
//
// TODO(gri): move this into antha/types.
//
func (info *PackageInfo) IsType(e ast.Expr) bool {
	switch e := e.(type) {
	case *ast.SelectorExpr: // pkg.Type
		if sel := info.Selections[e]; sel.Kind() == types.PackageObj {
			_, isType := sel.Obj().(*types.TypeName)
			return isType
		}
	case *ast.StarExpr: // *T
		return info.IsType(e.X)
	case *ast.Ident:
		_, isType := info.ObjectOf(e).(*types.TypeName)
		return isType
	case *ast.ArrayType, *ast.StructType, *ast.FuncType, *ast.InterfaceType, *ast.MapType, *ast.ChanType:
		return true
	case *ast.ParenExpr:
		return info.IsType(e.X)
	}
	return false
}

// TypeCaseVar returns the implicit variable created by a single-type
// case clause in a type switch, or nil if not found.
//
func (info *PackageInfo) TypeCaseVar(cc *ast.CaseClause) *types.Var {
	if v := info.Implicits[cc]; v != nil {
		return v.(*types.Var)
	}
	return nil
}

// ImportSpecPkg returns the PkgName for a given ImportSpec, possibly
// an implicit one for a dot-import or an import-without-rename.
// It returns nil if not found.
//
func (info *PackageInfo) ImportSpecPkg(spec *ast.ImportSpec) *types.PkgName {
	if spec.Name != nil {
		return info.ObjectOf(spec.Name).(*types.PkgName)
	}
	if p := info.Implicits[spec]; p != nil {
		return p.(*types.PkgName)
	}
	return nil
}