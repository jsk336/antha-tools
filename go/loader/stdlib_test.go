// antha-tools/go/loader/stdlib_test.go: Part of the Antha language
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


package loader_test

// This file enumerates all packages beneath $GOROOT, loads them, plus
// their external tests if any, runs the type checker on them, and
// prints some summary information.
//
// Run test with GOMAXPROCS=8.

import (
	"fmt"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/token"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/antha-lang/antha-tools/antha/loader"
	"github.com/antha-lang/antha-tools/antha/types"
)

func allPackages() []string {
	var pkgs []string
	root := filepath.Join(runtime.GOROOT(), "src/pkg") + string(os.PathSeparator)
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// Prune the search if we encounter any of these names:
		switch filepath.Base(path) {
		case "testdata", ".hg":
			return filepath.SkipDir
		}
		if info.IsDir() {
			pkg := filepath.ToSlash(strings.TrimPrefix(path, root))
			switch pkg {
			case "builtin", "pkg":
				return filepath.SkipDir // skip these subtrees
			case "":
				return nil // ignore root of tree
			}
			pkgs = append(pkgs, pkg)
		}

		return nil
	})
	return pkgs
}

func TestStdlib(t *testing.T) {
	runtime.GC()
	t0 := time.Now()
	var memstats runtime.MemStats
	runtime.ReadMemStats(&memstats)
	alloc := memstats.Alloc

	// Load, parse and type-check the program.
	var conf loader.Config
	for _, path := range allPackages() {
		if err := conf.ImportWithTests(path); err != nil {
			t.Error(err)
		}
	}

	prog, err := conf.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	t1 := time.Now()
	runtime.GC()
	runtime.ReadMemStats(&memstats)

	numPkgs := len(prog.AllPackages)
	if want := 205; numPkgs < want {
		t.Errorf("Loaded only %d packages, want at least %d", numPkgs, want)
	}

	// Dump package members.
	if false {
		for pkg := range prog.AllPackages {
			fmt.Printf("Package %s:\n", pkg.Path())
			scope := pkg.Scope()
			for _, name := range scope.Names() {
				if ast.IsExported(name) {
					fmt.Printf("\t%s\n", types.ObjectString(pkg, scope.Lookup(name)))
				}
			}
			fmt.Println()
		}
	}

	// Check that Test functions for io/ioutil, regexp and
	// compress/bzip2 are all simultaneously present.
	// (The apparent cycle formed when augmenting all three of
	// these packages by their tests was the original motivation
	// for reporting b/7114.)
	//
	// compress/bzip2.TestBitReader in bzip2_test.go    imports io/ioutil
	// io/ioutil.TestTempFile       in tempfile_test.go imports regexp
	// regexp.TestRE2Search         in exec_test.go     imports compress/bzip2
	for _, test := range []struct{ pkg, fn string }{
		{"io/ioutil", "TestTempFile"},
		{"regexp", "TestRE2Search"},
		{"compress/bzip2", "TestBitReader"},
	} {
		info := prog.Imported[test.pkg]
		if info == nil {
			t.Errorf("failed to load package %q", test.pkg)
			continue
		}
		obj, _ := info.Pkg.Scope().Lookup(test.fn).(*types.Func)
		if obj == nil {
			t.Errorf("package %q has no func %q", test.pkg, test.fn)
			continue
		}
	}

	// Dump some statistics.

	// determine line count
	var lineCount int
	prog.Fset.Iterate(func(f *token.File) bool {
		lineCount += f.LineCount()
		return true
	})

	t.Log("GOMAXPROCS:           ", runtime.GOMAXPROCS(0))
	t.Log("#Source lines:        ", lineCount)
	t.Log("Load/parse/typecheck: ", t1.Sub(t0))
	t.Log("#MB:                  ", int64(memstats.Alloc-alloc)/1000000)
}