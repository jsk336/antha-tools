// antha-tools/antha/types/api.go: Part of the Antha language
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


// Package types declares the data types and implements
// the algorithms for type-checking of Go packages.
// Use Check and Config.Check to invoke the type-checker.
//
// Type-checking consists of several interdependent phases:
//
// Name resolution maps each identifier (ast.Ident) in the program to the
// language object (Object) it denotes.
// Use Info.{Defs,Uses,Implicits} for the results of name resolution.
//
// Constant folding computes the exact constant value (exact.Value) for
// every expression (ast.Expr) that is a compile-time constant.
// Use Info.Types[expr].Value for the results of constant folding.
//
// Type inference computes the type (Type) of every expression (ast.Expr)
// and checks for compliance with the language specification.
// Use Info.Types[expr].Type for the results of type inference.
//
package types

import (
	"bytes"
	"fmt"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/token"

	"github.com/antha-lang/antha-tools/antha/exact"
)

// Check type-checks a package and returns the resulting complete package
// object, or a nil package and the first error. The package is specified
// by a list of *ast.Files and corresponding file set, and the import path
// the package is identified with. The clean path must not be empty or dot (".").
//
// For more control over type-checking and results, use Config.Check.
func Check(path string, fset *token.FileSet, files []*ast.File) (*Package, error) {
	var conf Config
	pkg, err := conf.Check(path, fset, files, nil)
	if err != nil {
		return nil, err
	}
	return pkg, nil
}

// An Error describes a type-checking error; it implements the error interface.
// A "soft" error is an error that still permits a valid interpretation of a
// package (such as "unused variable"); "hard" errors may lead to unpredictable
// behavior if ignored.
type Error struct {
	Fset *token.FileSet // file set for interpretation of Pos
	Pos  token.Pos      // error position
	Msg  string         // error message
	Soft bool           // if set, error is "soft"
}

// Error returns an error string formatted as follows:
// filename:line:column: message
func (err Error) Error() string {
	return fmt.Sprintf("%s: %s", err.Fset.Position(err.Pos), err.Msg)
}

// An importer resolves import paths to Packages.
// The imports map records packages already known,
// indexed by package path. The type-checker
// will invoke Import with Config.Packages.
// An importer must determine the canonical package path and
// check imports to see if it is already present in the map.
// If so, the Importer can return the map entry.  Otherwise,
// the importer must load the package data for the given path
// into a new *Package, record it in imports map, and return
// the package.
// TODO(gri) Need to be clearer about requirements of completeness.
type Importer func(map[string]*Package, string) (*Package, error)

// A Config specifies the configuration for type checking.
// The zero value for Config is a ready-to-use default configuration.
type Config struct {
	// If IgnoreFuncBodies is set, function bodies are not
	// type-checked.
	IgnoreFuncBodies bool

	// If FakeImportC is set, `import "C"` (for packages requiring Cgo)
	// declares an empty "C" package and errors are omitted for qualified
	// identifiers referring to package C (which won't find an object).
	// This feature is intended for the standard library cmd/api tool.
	//
	// Caution: Effects may be unpredictable due to follow-up errors.
	//          Do not use casually!
	FakeImportC bool

	// Packages is used to look up (and thus canonicalize) packages by
	// package path. If Packages is nil, it is set to a new empty map.
	// During type-checking, imported packages are added to the map.
	Packages map[string]*Package

	// If Error != nil, it is called with each error found
	// during type checking; err has dynamic type Error.
	// Secondary errors (for instance, to enumerate all types
	// involved in an invalid recursive type declaration) have
	// error strings that start with a '\t' character.
	Error func(err error)

	// If Import != nil, it is called for each imported package.
	// Otherwise, DefaultImport is called.
	Import Importer

	// If Sizes != nil, it provides the sizing functions for package unsafe.
	// Otherwise &StdSize{WordSize: 8, MaxAlign: 8} is used instead.
	Sizes Sizes
}

// DefaultImport is the default importer invoked if Config.Import == nil.
// The declaration:
//
//	import _ "github.com/antha-lang/antha-tools/antha/gcimporter"
//
// in a client of antha/types will initialize DefaultImport to gcimporter.Import.
var DefaultImport Importer

type TypeAndValue struct {
	Type  Type
	Value exact.Value
}

// Info holds result type information for a type-checked package.
// Only the information for which a map is provided is collected.
// If the package has type errors, the collected information may
// be incomplete.
type Info struct {
	// Types maps expressions to their types, and for constant
	// expressions, their values.
	// Identifiers are collected in Defs and Uses, not Types.
	//
	// For an expression denoting a predeclared built-in function
	// the recorded signature is call-site specific. If the call
	// result is not a constant, the recorded type is an argument-
	// specific signature. Otherwise, the recorded type is invalid.
	Types map[ast.Expr]TypeAndValue

	// Defs maps identifiers to the objects they define (including
	// package names, dots "." of dot-imports, and blank "_" identifiers).
	// For identifiers that do not denote objects (e.g., the package name
	// in package clauses, or symbolic variables t in t := x.(type) of
	// type switch headers), the corresponding objects are nil.
	//
	// For an anonymous field, Defs returns the field *Var it defines.
	//
	// Invariant: Defs[id] == nil || Defs[id].Pos() == id.Pos()
	Defs map[*ast.Ident]Object

	// Uses maps identifiers to the objects they denote.
	//
	// For an anonymous field, Uses returns the *TypeName it denotes.
	//
	// Invariant: Uses[id].Pos() != id.Pos()
	Uses map[*ast.Ident]Object

	// Implicits maps nodes to their implicitly declared objects, if any.
	// The following node and object types may appear:
	//
	//	node               declared object
	//
	//	*ast.ImportSpec    *PkgName for dot-imports and imports without renames
	//	*ast.CaseClause    type-specific *Var for each type switch case clause (incl. default)
	//      *ast.Field         anonymous struct field or parameter *Var
	//
	Implicits map[ast.Node]Object

	// Selections maps selector expressions to their corresponding selections.
	Selections map[*ast.SelectorExpr]*Selection

	// Scopes maps ast.Nodes to the scopes they define. Package scopes are not
	// associated with a specific node but with all files belonging to a package.
	// Thus, the package scope can be found in the type-checked Package object.
	// Scopes nest, with the Universe scope being the outermost scope, enclosing
	// the package scope, which contains (one or more) files scopes, which enclose
	// function scopes which in turn enclose statement and function literal scopes.
	// Note that even though package-level functions are declared in the package
	// scope, the function scopes are embedded in the file scope of the file
	// containing the function declaration.
	//
	// The following node types may appear in Scopes:
	//
	//	*ast.File
	//	*ast.FuncType
	//	*ast.BlockStmt
	//	*ast.IfStmt
	//	*ast.SwitchStmt
	//	*ast.TypeSwitchStmt
	//	*ast.CaseClause
	//	*ast.CommClause
	//	*ast.ForStmt
	//	*ast.RangeStmt
	//
	Scopes map[ast.Node]*Scope

	// InitOrder is the list of package-level initializers in the order in which
	// they must be executed. Initializers referring to variables related by an
	// initialization dependency appear in topological order, the others appear
	// in source order. Variables without an initialization expression do not
	// appear in this list.
	InitOrder []*Initializer
}

// An Initializer describes a package-level variable, or a list of variables in case
// of a multi-valued initialization expression, and the corresponding initialization
// expression.
type Initializer struct {
	Lhs []*Var // var Lhs = Rhs
	Rhs ast.Expr
}

func (init *Initializer) String() string {
	var buf bytes.Buffer
	for i, lhs := range init.Lhs {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(lhs.Name())
	}
	buf.WriteString(" = ")
	WriteExpr(&buf, init.Rhs)
	return buf.String()
}

// Check type-checks a package and returns the resulting package object,
// the first error if any, and if info != nil, additional type information.
// The package is marked as complete if no errors occurred, otherwise it is
// incomplete.
//
// The package is specified by a list of *ast.Files and corresponding
// file set, and the package path the package is identified with.
// The clean path must not be empty or dot (".").
func (conf *Config) Check(path string, fset *token.FileSet, files []*ast.File, info *Info) (*Package, error) {
	pkg := NewPackage(path, "")
	return pkg, NewChecker(conf, fset, pkg, info).Files(files)
}

// AssertableTo reports whether a value of type V can be asserted to have type T.
func AssertableTo(V *Interface, T Type) bool {
	f, _ := MissingMethod(T, V, false)
	return f == nil
}

// AssignableTo reports whether a value of type V is assignable to a variable of type T.
func AssignableTo(V, T Type) bool {
	x := operand{mode: value, typ: V}
	return x.assignableTo(nil, T) // config not needed for non-constant x
}

// ConvertibleTo reports whether a value of type V is convertible to a value of type T.
func ConvertibleTo(V, T Type) bool {
	x := operand{mode: value, typ: V}
	return x.convertibleTo(nil, T) // config not needed for non-constant x
}

// Implements reports whether type V implements interface T.
func Implements(V Type, T *Interface) bool {
	f, _ := MissingMethod(V, T, true)
	return f == nil
}