// antha-tools/godoc/parser.go: Part of the Antha language
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


// This file contains support functions for parsing .go files
// accessed via godoc's file system fs.

package godoc

import (
	"bytes"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/parser"
	"github.com/antha-lang/antha/token"
	pathpkg "path"

	"github.com/antha-lang/antha-tools/anthadoc/vfs"
)

var linePrefix = []byte("//line ")

// This function replaces source lines starting with "//line " with a blank line.
// It does this irrespective of whether the line is truly a line comment or not;
// e.g., the line may be inside a string, or a /*-style comment; however that is
// rather unlikely (proper testing would require a full Go scan which we want to
// avoid for performance).
func replaceLinePrefixCommentsWithBlankLine(src []byte) {
	for {
		i := bytes.Index(src, linePrefix)
		if i < 0 {
			break // we're done
		}
		// 0 <= i && i+len(linePrefix) <= len(src)
		if i == 0 || src[i-1] == '\n' {
			// at beginning of line: blank out line
			for i < len(src) && src[i] != '\n' {
				src[i] = ' '
				i++
			}
		} else {
			// not at beginning of line: skip over prefix
			i += len(linePrefix)
		}
		// i <= len(src)
		src = src[i:]
	}
}

func (c *Corpus) parseFile(fset *token.FileSet, filename string, mode parser.Mode) (*ast.File, error) {
	src, err := vfs.ReadFile(c.fs, filename)
	if err != nil {
		return nil, err
	}

	// Temporary ad-hoc fix for issue 5247.
	// TODO(gri) Remove this in favor of a better fix, eventually (see issue 7702).
	replaceLinePrefixCommentsWithBlankLine(src)

	return parser.ParseFile(fset, filename, src, mode)
}

func (c *Corpus) parseFiles(fset *token.FileSet, relpath string, abspath string, localnames []string) (map[string]*ast.File, error) {
	files := make(map[string]*ast.File)
	for _, f := range localnames {
		absname := pathpkg.Join(abspath, f)
		file, err := c.parseFile(fset, absname, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		files[pathpkg.Join(relpath, f)] = file
	}

	return files, nil
}