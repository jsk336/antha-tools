// antha-tools/refactor/eg/eg.go: Part of the Antha language
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

// Package eg implements the example-based refactoring tool whose
// command-line is defined in antha-tools/cmd/eg.
package eg

import (
	"bytes"
	"fmt"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/printer"
	"github.com/antha-lang/antha/token"
	"os"

	"github.com/antha-lang/antha-tools/antha/loader"
	"github.com/antha-lang/antha-tools/antha/types"
)

const Help = `
This tool implements example-based refactoring of expressions.

The transformation is specified as a Go file defining two functions,
'before' and 'after', of identical types.  Each function body consists
of a single statement: either a return statement with a single
(possibly multi-valued) expression, or an expression statement.  The
'before' expression specifies a pattern and the 'after' expression its
replacement.

	package P
 	import ( "errors"; "fmt" )
 	func before(s string) error { return fmt.Errorf("%s", s) }
 	func after(s string)  error { return errors.New(s) }

The expression statement form is useful when the expression has no
result, for example:

 	func before(msg string) { log.Fatalf("%s", msg) }
 	func after(msg string)  { log.Fatal(msg) }

The parameters of both functions are wildcards that may match any
expression assignable to that type.  If the pattern contains multiple
occurrences of the same parameter, each must match the same expression
in the input for the pattern to match.  If the replacement contains
multiple occurrences of the same parameter, the expression will be
duplicated, possibly changing the side-effects.

The tool analyses all Go code in the packages specified by the
arguments, replacing all occurrences of the pattern with the
substitution.

So, the transform above would change this input:
	err := fmt.Errorf("%s", "error: " + msg)
to this output:
	err := errors.New("error: " + msg)

Identifiers, including qualified identifiers (p.X) are considered to
match only if they denote the same object.  This allows correct
matching even in the presence of dot imports, named imports and
locally shadowed package names in the input program.

Matching of type syntax is semantic, not syntactic: type syntax in the
pattern matches type syntax in the input if the types are identical.
Thus, func(x int) matches func(y int).

This tool was inspired by other example-based refactoring tools,
'gofmt -r' for Go and Refaster for Java.


LIMITATIONS
===========

EXPRESSIVENESS

Only refactorings that replace one expression with another, regardless
of the expression's context, may be expressed.  Refactoring arbitrary
statements (or sequences of statements) is a less well-defined problem
and is less amenable to this approach.

A pattern that contains a function literal (and hence statements)
never matches.

There is no way to generalize over related types, e.g. to express that
a wildcard may have any integer type, for example.

It is not possible to replace an expression by one of a different
type, even in contexts where this is legal, such as x in fmt.Print(x).


SAFETY

Verifying that a transformation does not introduce type errors is very
complex in the general case.  An innocuous-looking replacement of one
constant by another (e.g. 1 to 2) may cause type errors relating to
array types and indices, for example.  The tool performs only very
superficial checks of type preservation.


IMPORTS

Although the matching algorithm is fully aware of scoping rules, the
replacement algorithm is not, so the replacement code may contain
incorrect identifier syntax for imported objects if there are dot
imports, named imports or locally shadowed package names in the input
program.

Imports are added as needed, but they are not removed as needed.
Run 'goimports' on the modified file for now.

Dot imports are forbidden in the template.
`

// TODO(adonovan): allow the tool to be invoked using relative package
// directory names (./foo).  Requires changes to antha/loader.

// TODO(adonovan): expand upon the above documentation as an HTML page.

// TODO(adonovan): eliminate dependency on loader.PackageInfo.
// Move its ObjectOf/IsType/TypeOf methods into antha/types.

// A Transformer represents a single example-based transformation.
type Transformer struct {
	fset           *token.FileSet
	verbose        bool
	info           loader.PackageInfo // combined type info for template/input/output ASTs
	seenInfos      map[*types.Info]bool
	wildcards      map[*types.Var]bool                // set of parameters in func before()
	env            map[string]ast.Expr                // maps parameter name to wildcard binding
	importedObjs   map[types.Object]*ast.SelectorExpr // objects imported by after().
	before, after  ast.Expr
	allowWildcards bool

	// Working state of Transform():
	nsubsts    int            // number of substitutions made
	currentPkg *types.Package // package of current call
}

// NewTransformer returns a transformer based on the specified template,
// a package containing "before" and "after" functions as described
// in the package documentation.
//
func NewTransformer(fset *token.FileSet, template *loader.PackageInfo, verbose bool) (*Transformer, error) {
	// Check the template.
	beforeSig := funcSig(template.Pkg, "before")
	if beforeSig == nil {
		return nil, fmt.Errorf("no 'before' func found in template")
	}
	afterSig := funcSig(template.Pkg, "after")
	if afterSig == nil {
		return nil, fmt.Errorf("no 'after' func found in template")
	}

	// TODO(adonovan): should we also check the names of the params match?
	if !types.Identical(afterSig, beforeSig) {
		return nil, fmt.Errorf("before %s and after %s functions have different signatures",
			beforeSig, afterSig)
	}

	templateFile := template.Files[0]
	for _, imp := range templateFile.Imports {
		if imp.Name != nil && imp.Name.Name == "." {
			// Dot imports are currently forbidden.  We
			// make the simplifying assumption that all
			// imports are regular, without local renames.
			//TODO document
			return nil, fmt.Errorf("dot-import (of %s) in template", imp.Path.Value)
		}
	}
	var beforeDecl, afterDecl *ast.FuncDecl
	for _, decl := range templateFile.Decls {
		if decl, ok := decl.(*ast.FuncDecl); ok {
			switch decl.Name.Name {
			case "before":
				beforeDecl = decl
			case "after":
				afterDecl = decl
			}
		}
	}

	before, err := soleExpr(beforeDecl)
	if err != nil {
		return nil, fmt.Errorf("before: %s", err)
	}
	after, err := soleExpr(afterDecl)
	if err != nil {
		return nil, fmt.Errorf("after: %s", err)
	}

	wildcards := make(map[*types.Var]bool)
	for i := 0; i < beforeSig.Params().Len(); i++ {
		wildcards[beforeSig.Params().At(i)] = true
	}

	// checkExprTypes returns an error if Tb (type of before()) is not
	// safe to replace with Ta (type of after()).
	//
	// Only superficial checks are performed, and they may result in both
	// false positives and negatives.
	//
	// Ideally, we would only require that the replacement be assignable
	// to the context of a specific pattern occurrence, but the type
	// checker doesn't record that information and it's complex to deduce.
	// A Go type cannot capture all the constraints of a given expression
	// context, which may include the size, constness, signedness,
	// namedness or constructor of its type, and even the specific value
	// of the replacement.  (Consider the rule that array literal keys
	// must be unique.)  So we cannot hope to prove the safety of a
	// transformation in general.
	Tb := template.TypeOf(before)
	Ta := template.TypeOf(after)
	if types.AssignableTo(Tb, Ta) {
		// safe: replacement is assignable to pattern.
	} else if tuple, ok := Tb.(*types.Tuple); ok && tuple.Len() == 0 {
		// safe: pattern has void type (must appear in an ExprStmt).
	} else {
		return nil, fmt.Errorf("%s is not a safe replacement for %s", Ta, Tb)
	}

	tr := &Transformer{
		fset:           fset,
		verbose:        verbose,
		wildcards:      wildcards,
		allowWildcards: true,
		seenInfos:      make(map[*types.Info]bool),
		importedObjs:   make(map[types.Object]*ast.SelectorExpr),
		before:         before,
		after:          after,
	}

	// Combine type info from the template and input packages, and
	// type info for the synthesized ASTs too.  This saves us
	// having to book-keep where each ast.Node originated as we
	// construct the resulting hybrid AST.
	//
	// TODO(adonovan): move type utility methods of PackageInfo to
	// types.Info, or at least into antha/types.typeutil.
	tr.info.Info = types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
	}
	mergeTypeInfo(&tr.info.Info, &template.Info)

	// Compute set of imported objects required by after().
	// TODO reject dot-imports in pattern
	ast.Inspect(after, func(n ast.Node) bool {
		if n, ok := n.(*ast.SelectorExpr); ok {
			sel := tr.info.Selections[n]
			if sel.Kind() == types.PackageObj {
				tr.importedObjs[sel.Obj()] = n
				return false // prune
			}
		}
		return true // recur
	})

	return tr, nil
}

// WriteAST is a convenience function that writes AST f to the specified file.
func WriteAST(fset *token.FileSet, filename string, f *ast.File) (err error) {
	fh, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if err2 := fh.Close(); err != nil {
			err = err2 // prefer earlier error
		}
	}()
	return printer.Fprint(fh, fset, f)
}

// -- utilities --------------------------------------------------------

// funcSig returns the signature of the specified package-level function.
func funcSig(pkg *types.Package, name string) *types.Signature {
	if f, ok := pkg.Scope().Lookup(name).(*types.Func); ok {
		return f.Type().(*types.Signature)
	}
	return nil
}

// soleExpr returns the sole expression in the before/after template function.
func soleExpr(fn *ast.FuncDecl) (ast.Expr, error) {
	if fn.Body == nil {
		return nil, fmt.Errorf("no body")
	}
	if len(fn.Body.List) != 1 {
		return nil, fmt.Errorf("must contain a single statement")
	}
	switch stmt := fn.Body.List[0].(type) {
	case *ast.ReturnStmt:
		if len(stmt.Results) != 1 {
			return nil, fmt.Errorf("return statement must have a single operand")
		}
		return stmt.Results[0], nil

	case *ast.ExprStmt:
		return stmt.X, nil
	}

	return nil, fmt.Errorf("must contain a single return or expression statement")
}

// mergeTypeInfo adds type info from src to dst.
func mergeTypeInfo(dst, src *types.Info) {
	for k, v := range src.Types {
		dst.Types[k] = v
	}
	for k, v := range src.Defs {
		dst.Defs[k] = v
	}
	for k, v := range src.Uses {
		dst.Uses[k] = v
	}
	for k, v := range src.Selections {
		dst.Selections[k] = v
	}
}

// (debugging only)
func astString(fset *token.FileSet, n ast.Node) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, fset, n)
	return buf.String()
}