// antha-tools/go/ssa/testmain.go: Part of the Antha language
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


package ssa

// CreateTestMainPackage synthesizes a main package that runs all the
// tests of the supplied packages.
// It is closely coupled to src/cmd/antha/test.go and src/pkg/testing.

import (
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/token"
	"os"
	"strings"

	"github.com/antha-lang/antha-tools/antha/exact"
	"github.com/antha-lang/antha-tools/antha/types"
)

// CreateTestMainPackage creates and returns a synthetic "main"
// package that runs all the tests of the supplied packages, similar
// to the one that would be created by the 'go test' tool.
//
// It returns nil if the program contains no tests.
//
func (prog *Program) CreateTestMainPackage(pkgs ...*Package) *Package {
	if len(pkgs) == 0 {
		return nil
	}
	testmain := &Package{
		Prog:    prog,
		Members: make(map[string]Member),
		values:  make(map[types.Object]Value),
		Object:  types.NewPackage("testmain", "testmain"),
	}

	// Build package's init function.
	init := &Function{
		name:      "init",
		Signature: new(types.Signature),
		Synthetic: "package initializer",
		Pkg:       testmain,
		Prog:      prog,
	}
	init.startBody()
	// TODO(adonovan): use lexical order.
	var expfuncs []*Function // all exported functions of *_test.go in pkgs, unordered
	for _, pkg := range pkgs {
		if pkg.Prog != prog {
			panic("wrong Program")
		}
		// Initialize package to test.
		var v Call
		v.Call.Value = pkg.init
		v.setType(types.NewTuple())
		init.emit(&v)

		// Enumerate its possible tests/benchmarks.
		for _, mem := range pkg.Members {
			if f, ok := mem.(*Function); ok &&
				ast.IsExported(f.Name()) &&
				strings.HasSuffix(prog.Fset.Position(f.Pos()).Filename, "_test.go") {
				expfuncs = append(expfuncs, f)
			}
		}
	}
	init.emit(new(Return))
	init.finishBody()
	testmain.init = init
	testmain.Object.MarkComplete()
	testmain.Members[init.name] = init

	testingPkg := prog.ImportedPackage("testing")
	if testingPkg == nil {
		// If the program doesn't import "testing", it can't
		// contain any tests.
		// TODO(adonovan): but it might contain Examples.
		// Support them (by just calling them directly).
		return nil
	}
	testingMain := testingPkg.Func("Main")
	testingMainParams := testingMain.Signature.Params()

	// The generated code is as if compiled from this:
	//
	// func main() {
	//      match      := func(_, _ string) (bool, error) { return true, nil }
	//      tests      := []testing.InternalTest{{"TestFoo", TestFoo}, ...}
	//      benchmarks := []testing.InternalBenchmark{...}
	//      examples   := []testing.InternalExample{...}
	// 	testing.Main(match, tests, benchmarks, examples)
	// }

	main := &Function{
		name:      "main",
		Signature: new(types.Signature),
		Synthetic: "test main function",
		Prog:      prog,
		Pkg:       testmain,
	}

	matcher := &Function{
		name:      "matcher",
		Signature: testingMainParams.At(0).Type().(*types.Signature),
		Synthetic: "test matcher predicate",
		Enclosing: main,
		Pkg:       testmain,
		Prog:      prog,
	}
	main.AnonFuncs = append(main.AnonFuncs, matcher)
	matcher.startBody()
	matcher.emit(&Return{Results: []Value{vTrue, nilConst(types.Universe.Lookup("error").Type())}})
	matcher.finishBody()

	main.startBody()
	var c Call
	c.Call.Value = testingMain

	tests := testMainSlice(main, expfuncs, "Test", testingMainParams.At(1).Type())
	benchmarks := testMainSlice(main, expfuncs, "Benchmark", testingMainParams.At(2).Type())
	examples := testMainSlice(main, expfuncs, "Example", testingMainParams.At(3).Type())
	_, noTests := tests.(*Const) // i.e. nil slice
	_, noBenchmarks := benchmarks.(*Const)
	_, noExamples := examples.(*Const)
	if noTests && noBenchmarks && noExamples {
		return nil
	}

	c.Call.Args = []Value{matcher, tests, benchmarks, examples}
	// Emit: testing.Main(nil, tests, benchmarks, examples)
	emitTailCall(main, &c)
	main.finishBody()

	testmain.Members["main"] = main

	if prog.mode&LogPackages != 0 {
		testmain.WriteTo(os.Stderr)
	}

	if prog.mode&SanityCheckFunctions != 0 {
		sanityCheckPackage(testmain)
	}

	prog.packages[testmain.Object] = testmain

	return testmain
}

// testMainSlice emits to fn code to construct a slice of type slice
// (one of []testing.Internal{Test,Benchmark,Example}) for all
// functions in expfuncs whose name starts with prefix (one of
// "Test", "Benchmark" or "Example") and whose type is appropriate.
// It returns the slice value.
//
func testMainSlice(fn *Function, expfuncs []*Function, prefix string, slice types.Type) Value {
	tElem := slice.(*types.Slice).Elem()
	tFunc := tElem.Underlying().(*types.Struct).Field(1).Type()

	var testfuncs []*Function
	for _, f := range expfuncs {
		if isTest(f.Name(), prefix) && types.Identical(f.Signature, tFunc) {
			testfuncs = append(testfuncs, f)
		}
	}
	if testfuncs == nil {
		return nilConst(slice)
	}

	tString := types.Typ[types.String]
	tPtrString := types.NewPointer(tString)
	tPtrElem := types.NewPointer(tElem)
	tPtrFunc := types.NewPointer(tFunc)

	// Emit: array = new [n]testing.InternalTest
	tArray := types.NewArray(tElem, int64(len(testfuncs)))
	array := emitNew(fn, tArray, token.NoPos)
	array.Comment = "test main"
	for i, testfunc := range testfuncs {
		// Emit: pitem = &array[i]
		ia := &IndexAddr{X: array, Index: intConst(int64(i))}
		ia.setType(tPtrElem)
		pitem := fn.emit(ia)

		// Emit: pname = &pitem.Name
		fa := &FieldAddr{X: pitem, Field: 0} // .Name
		fa.setType(tPtrString)
		pname := fn.emit(fa)

		// Emit: *pname = "testfunc"
		emitStore(fn, pname, NewConst(exact.MakeString(testfunc.Name()), tString))

		// Emit: pfunc = &pitem.F
		fa = &FieldAddr{X: pitem, Field: 1} // .F
		fa.setType(tPtrFunc)
		pfunc := fn.emit(fa)

		// Emit: *pfunc = testfunc
		emitStore(fn, pfunc, testfunc)
	}

	// Emit: slice array[:]
	sl := &Slice{X: array}
	sl.setType(slice)
	return fn.emit(sl)
}

// Plundered from $GOROOT/src/cmd/antha/test.go

// isTest tells whether name looks like a test (or benchmark, according to prefix).
// It is a Test (say) if there is a character after Test that is not a lower-case letter.
// We don't want TesticularCancer.
func isTest(name, prefix string) bool {
	if !strings.HasPrefix(name, prefix) {
		return false
	}
	if len(name) == len(prefix) { // "Test" is ok
		return true
	}
	return ast.IsExported(name[len(prefix):])
}