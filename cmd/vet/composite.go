// antha-tools/cmd/vet/composite.go: Part of the Antha language
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


// This file contains the test for unkeyed struct literals.

package main

import (
	"flag"
	"github.com/antha-lang/antha/ast"
	"strings"

	"github.com/antha-lang/antha-tools/cmd/vet/whitelist"
)

var compositeWhiteList = flag.Bool("compositewhitelist", true, "use composite white list; for testing only")

// checkUnkeyedLiteral checks if a composite literal is a struct literal with
// unkeyed fields.
func (f *File) checkUnkeyedLiteral(c *ast.CompositeLit) {
	if !vet("composites") {
		return
	}

	typ := c.Type
	for {
		if typ1, ok := c.Type.(*ast.ParenExpr); ok {
			typ = typ1
			continue
		}
		break
	}

	switch typ.(type) {
	case *ast.ArrayType:
		return
	case *ast.MapType:
		return
	case *ast.StructType:
		return // a literal struct type does not need to use keys
	case *ast.Ident:
		// A simple type name like t or T does not need keys either,
		// since it is almost certainly declared in the current package.
		// (The exception is names being used via import . "pkg", but
		// those are already breaking the Go 1 compatibility promise,
		// so not reporting potential additional breakage seems okay.)
		return
	}

	// Otherwise the type is a selector like pkg.Name.
	// We only care if pkg.Name is a struct, not if it's a map, array, or slice.
	isStruct, typeString := f.pkg.isStruct(c)
	if !isStruct {
		return
	}

	if typeString == "" { // isStruct doesn't know
		typeString = f.gofmt(typ)
	}

	// It's a struct, or we can't tell it's not a struct because we don't have types.

	// Check if the CompositeLit contains an unkeyed field.
	allKeyValue := true
	for _, e := range c.Elts {
		if _, ok := e.(*ast.KeyValueExpr); !ok {
			allKeyValue = false
			break
		}
	}
	if allKeyValue {
		return
	}

	// Check that the CompositeLit's type has the form pkg.Typ.
	s, ok := c.Type.(*ast.SelectorExpr)
	if !ok {
		return
	}
	pkg, ok := s.X.(*ast.Ident)
	if !ok {
		return
	}

	// Convert the package name to an import path, and compare to a whitelist.
	path := pkgPath(f, pkg.Name)
	if path == "" {
		f.Badf(c.Pos(), "unresolvable package for %s.%s literal", pkg.Name, s.Sel.Name)
		return
	}
	typeName := path + "." + s.Sel.Name
	if *compositeWhiteList && whitelist.UnkeyedLiteral[typeName] {
		return
	}

	f.Bad(c.Pos(), typeString+" composite literal uses unkeyed fields")
}

// pkgPath returns the import path "image/png" for the package name "png".
//
// This is based purely on syntax and convention, and not on the imported
// package's contents. It will be incorrect if a package name differs from the
// leaf element of the import path, or if the package was a dot import.
func pkgPath(f *File, pkgName string) (path string) {
	for _, x := range f.file.Imports {
		s := strings.Trim(x.Path.Value, `"`)
		if x.Name != nil {
			// Catch `import pkgName "foo/bar"`.
			if x.Name.Name == pkgName {
				return s
			}
		} else {
			// Catch `import "pkgName"` or `import "foo/bar/pkgName"`.
			if s == pkgName || strings.HasSuffix(s, "/"+pkgName) {
				return s
			}
		}
	}
	return ""
}