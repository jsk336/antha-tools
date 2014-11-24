// antha-tools/antha/ssa/stdlib_test.go: Part of the Antha language
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


package ssa_test

// This file runs the SSA builder in sanity-checking mode on all
// packages beneath $GOROOT and prints some summary information.
//
// Run test with GOMAXPROCS=8.

import (
	"github.com/antha-lang/antha/token"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/antha-lang/antha-tools/antha/loader"
	"github.com/antha-lang/antha-tools/antha/ssa"
	"github.com/antha-lang/antha-tools/antha/ssa/ssautil"
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
	// Load, parse and type-check the program.
	t0 := time.Now()

	var conf loader.Config
	conf.SourceImports = true
	if _, err := conf.FromArgs(allPackages(), true); err != nil {
		t.Errorf("FromArgs failed: %v", err)
		return
	}

	iprog, err := conf.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	t1 := time.Now()

	runtime.GC()
	var memstats runtime.MemStats
	runtime.ReadMemStats(&memstats)
	alloc := memstats.Alloc

	// Create SSA packages.
	var mode ssa.BuilderMode
	// Comment out these lines during benchmarking.  Approx SSA build costs are noted.
	mode |= ssa.SanityCheckFunctions // + 2% space, + 4% time
	mode |= ssa.GlobalDebug          // +30% space, +18% time
	prog := ssa.Create(iprog, mode)

	t2 := time.Now()

	// Build SSA.
	prog.BuildAll()

	t3 := time.Now()

	runtime.GC()
	runtime.ReadMemStats(&memstats)

	numPkgs := len(prog.AllPackages())
	if want := 140; numPkgs < want {
		t.Errorf("Loaded only %d packages, want at least %d", numPkgs, want)
	}

	// Dump some statistics.
	allFuncs := ssautil.AllFunctions(prog)
	var numInstrs int
	for fn := range allFuncs {
		for _, b := range fn.Blocks {
			numInstrs += len(b.Instrs)
		}
	}

	// determine line count
	var lineCount int
	prog.Fset.Iterate(func(f *token.File) bool {
		lineCount += f.LineCount()
		return true
	})

	// NB: when benchmarking, don't forget to clear the debug +
	// sanity builder flags for better performance.

	t.Log("GOMAXPROCS:           ", runtime.GOMAXPROCS(0))
	t.Log("#Source lines:        ", lineCount)
	t.Log("Load/parse/typecheck: ", t1.Sub(t0))
	t.Log("SSA create:           ", t2.Sub(t1))
	t.Log("SSA build:            ", t3.Sub(t2))

	// SSA stats:
	t.Log("#Packages:            ", numPkgs)
	t.Log("#Functions:           ", len(allFuncs))
	t.Log("#Instructions:        ", numInstrs)
	t.Log("#MB:                  ", int64(memstats.Alloc-alloc)/1000000)
}