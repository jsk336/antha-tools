// antha-tools/go/types/package.go: Part of the Antha language
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


package types

import "fmt"

// A Package describes a Go package.
type Package struct {
	path     string
	name     string
	scope    *Scope
	complete bool
	imports  []*Package
	fake     bool // scope lookup errors are silently dropped if package is fake (internal use only)
}

// NewPackage returns a new Package for the given package path and name;
// the name must not be the blank identifier.
// The package is not complete and contains no explicit imports.
func NewPackage(path, name string) *Package {
	if name == "_" {
		panic("invalid package name _")
	}
	scope := NewScope(Universe, fmt.Sprintf("package %q", path))
	return &Package{path: path, name: name, scope: scope}
}

// Path returns the package path.
func (pkg *Package) Path() string { return pkg.path }

// Name returns the package name.
func (pkg *Package) Name() string { return pkg.name }

// Scope returns the (complete or incomplete) package scope
// holding the objects declared at package level (TypeNames,
// Consts, Vars, and Funcs).
func (pkg *Package) Scope() *Scope { return pkg.scope }

// A package is complete if its scope contains (at least) all
// exported objects; otherwise it is incomplete.
func (pkg *Package) Complete() bool { return pkg.complete }

// MarkComplete marks a package as complete.
func (pkg *Package) MarkComplete() { pkg.complete = true }

// Imports returns the list of packages explicitly imported by
// pkg; the list is in source order. Package unsafe is excluded.
func (pkg *Package) Imports() []*Package { return pkg.imports }

// SetImports sets the list of explicitly imported packages to list.
// It is the caller's responsibility to make sure list elements are unique.
func (pkg *Package) SetImports(list []*Package) { pkg.imports = list }

func (pkg *Package) String() string {
	return fmt.Sprintf("package %s (%q)", pkg.name, pkg.path)
}