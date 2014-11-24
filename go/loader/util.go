// antha-tools/go/loader/util.go: Part of the Antha language
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
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/build"
	"github.com/antha-lang/antha/parser"
	"github.com/antha-lang/antha/token"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// parseFiles parses the Go source files files within directory dir
// and returns their ASTs, or the first parse error if any.
//
// I/O is done via ctxt, which may specify a virtual file system.
// displayPath is used to transform the filenames attached to the ASTs.
//
func parseFiles(fset *token.FileSet, ctxt *build.Context, displayPath func(string) string, dir string, files []string, mode parser.Mode) ([]*ast.File, error) {
	if displayPath == nil {
		displayPath = func(path string) string { return path }
	}
	isAbs := filepath.IsAbs
	if ctxt.IsAbsPath != nil {
		isAbs = ctxt.IsAbsPath
	}
	joinPath := filepath.Join
	if ctxt.JoinPath != nil {
		joinPath = ctxt.JoinPath
	}
	var wg sync.WaitGroup
	n := len(files)
	parsed := make([]*ast.File, n)
	errors := make([]error, n)
	for i, file := range files {
		if !isAbs(file) {
			file = joinPath(dir, file)
		}
		wg.Add(1)
		go func(i int, file string) {
			defer wg.Done()
			var rd io.ReadCloser
			var err error
			if ctxt.OpenFile != nil {
				rd, err = ctxt.OpenFile(file)
			} else {
				rd, err = os.Open(file)
			}
			defer rd.Close()
			if err != nil {
				errors[i] = err
				return
			}
			parsed[i], errors[i] = parser.ParseFile(fset, displayPath(file), rd, mode)
		}(i, file)
	}
	wg.Wait()

	for _, err := range errors {
		if err != nil {
			return nil, err
		}
	}
	return parsed, nil
}

// ---------- Internal helpers ----------

// TODO(adonovan): make this a method: func (*token.File) Contains(token.Pos)
func tokenFileContainsPos(f *token.File, pos token.Pos) bool {
	p := int(pos)
	base := f.Base()
	return base <= p && p < base+f.Size()
}