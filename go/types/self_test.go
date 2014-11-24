// antha-tools/go/types/self_test.go: Part of the Antha language
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


package types_test

import (
	"flag"
	"fmt"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/parser"
	"github.com/antha-lang/antha/token"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/antha-lang/antha-tools/antha/gcimporter"
	. "github.com/antha-lang/antha-tools/antha/types"
)

var benchmark = flag.Bool("b", false, "run benchmarks")

func TestSelf(t *testing.T) {
	fset := token.NewFileSet()
	files, err := pkgFiles(fset, ".")
	if err != nil {
		t.Fatal(err)
	}

	_, err = Check("github.com/antha-lang/antha/types", fset, files)
	if err != nil {
		// Importing go.tools/antha/exact doensn't work in the
		// build dashboard environment. Don't report an error
		// for now so that the build remains green.
		// TODO(gri) fix this
		t.Log(err) // replace w/ t.Fatal eventually
		return
	}
}

func TestBenchmark(t *testing.T) {
	if !*benchmark {
		return
	}

	// We're not using testing's benchmarking mechanism directly
	// because we want custom output.

	for _, p := range []string{"types", "exact", "gcimporter"} {
		path := filepath.Join("..", p)
		runbench(t, path, false)
		runbench(t, path, true)
		fmt.Println()
	}
}

func runbench(t *testing.T, path string, ignoreFuncBodies bool) {
	fset := token.NewFileSet()
	files, err := pkgFiles(fset, path)
	if err != nil {
		t.Fatal(err)
	}

	b := testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			conf := Config{IgnoreFuncBodies: ignoreFuncBodies}
			conf.Check(path, fset, files, nil)
		}
	})

	// determine line count
	lines := 0
	fset.Iterate(func(f *token.File) bool {
		lines += f.LineCount()
		return true
	})

	d := time.Duration(b.NsPerOp())
	fmt.Printf(
		"%s: %s for %d lines (%d lines/s), ignoreFuncBodies = %v\n",
		filepath.Base(path), d, lines, int64(float64(lines)/d.Seconds()), ignoreFuncBodies,
	)
}

func pkgFiles(fset *token.FileSet, path string) ([]*ast.File, error) {
	filenames, err := pkgFilenames(path) // from stdlib_test.go
	if err != nil {
		return nil, err
	}

	var files []*ast.File
	for _, filename := range filenames {
		file, err := parser.ParseFile(fset, filename, nil, 0)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}