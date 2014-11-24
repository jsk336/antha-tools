// antha-tools/imports/mkindex.go: Part of the Antha language
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

// +build ignore


// Command mkindex creates the file "pkgindex.go" containing an index of the Go
// standard library. The file is intended to be built as part of the imports
// package, so that the package may be used in environments where a GOROOT is
// not available (such as App Engine).
package main

import (
	"bytes"
	"fmt"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/build"
	"github.com/antha-lang/antha/format"
	"github.com/antha-lang/antha/parser"
	"github.com/antha-lang/antha/token"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	pkgIndex = make(map[string][]pkg)
	exports  = make(map[string]map[string]bool)
)

func main() {
	// Don't use GOPATH.
	ctx := build.Default
	ctx.GOPATH = ""

	// Populate pkgIndex global from GOROOT.
	for _, path := range ctx.SrcDirs() {
		f, err := os.Open(path)
		if err != nil {
			log.Print(err)
			continue
		}
		children, err := f.Readdir(-1)
		f.Close()
		if err != nil {
			log.Print(err)
			continue
		}
		for _, child := range children {
			if child.IsDir() {
				loadPkg(path, child.Name())
			}
		}
	}
	// Populate exports global.
	for _, ps := range pkgIndex {
		for _, p := range ps {
			e := loadExports(p.dir)
			if e != nil {
				exports[p.dir] = e
			}
		}
	}

	// Construct source file.
	var buf bytes.Buffer
	fmt.Fprint(&buf, pkgIndexHead)
	fmt.Fprintf(&buf, "var pkgIndexMaster = %#v\n", pkgIndex)
	fmt.Fprintf(&buf, "var exportsMaster = %#v\n", exports)
	src := buf.Bytes()

	// Replace main.pkg type name with pkg.
	src = bytes.Replace(src, []byte("main.pkg"), []byte("pkg"), -1)
	// Replace actual GOROOT with "/go".
	src = bytes.Replace(src, []byte(ctx.GOROOT), []byte("/go"), -1)
	// Add some line wrapping.
	src = bytes.Replace(src, []byte("}, "), []byte("},\n"), -1)
	src = bytes.Replace(src, []byte("true, "), []byte("true,\n"), -1)

	var err error
	src, err = format.Source(src)
	if err != nil {
		log.Fatal(err)
	}

	// Write out source file.
	err = ioutil.WriteFile("pkgindex.go", src, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

const pkgIndexHead = `package imports

func init() {
	pkgIndexOnce.Do(func() {
		pkgIndex.m = pkgIndexMaster
	})
	loadExports = func(dir string) map[string]bool {
		return exportsMaster[dir]
	}
}
`

type pkg struct {
	importpath string // full pkg import path, e.g. "net/http"
	dir        string // absolute file path to pkg directory e.g. "/usr/lib/antha/src/fmt"
}

var fset = token.NewFileSet()

func loadPkg(root, importpath string) {
	shortName := path.Base(importpath)
	if shortName == "testdata" {
		return
	}

	dir := filepath.Join(root, importpath)
	pkgIndex[shortName] = append(pkgIndex[shortName], pkg{
		importpath: importpath,
		dir:        dir,
	})

	pkgDir, err := os.Open(dir)
	if err != nil {
		return
	}
	children, err := pkgDir.Readdir(-1)
	pkgDir.Close()
	if err != nil {
		return
	}
	for _, child := range children {
		name := child.Name()
		if name == "" {
			continue
		}
		if c := name[0]; c == '.' || ('0' <= c && c <= '9') {
			continue
		}
		if child.IsDir() {
			loadPkg(root, filepath.Join(importpath, name))
		}
	}
}

func loadExports(dir string) map[string]bool {
	exports := make(map[string]bool)
	buildPkg, err := build.ImportDir(dir, 0)
	if err != nil {
		if strings.Contains(err.Error(), "no buildable Go source files in") {
			return nil
		}
		log.Printf("could not import %q: %v", dir, err)
		return nil
	}
	for _, file := range buildPkg.GoFiles {
		f, err := parser.ParseFile(fset, filepath.Join(dir, file), nil, 0)
		if err != nil {
			log.Printf("could not parse %q: %v", file, err)
			continue
		}
		for name := range f.Scope.Objects {
			if ast.IsExported(name) {
				exports[name] = true
			}
		}
	}
	return exports
}