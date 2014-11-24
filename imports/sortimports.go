// antha-tools/imports/sortimports.go: Part of the Antha language
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

// +build go1.2


// Hacked up copy of antha/ast/import.go

package imports

import (
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/token"
	"sort"
	"strconv"
)

// sortImports sorts runs of consecutive import lines in import blocks in f.
// It also removes duplicate imports when it is possible to do so without data loss.
func sortImports(fset *token.FileSet, f *ast.File) {
	for i, d := range f.Decls {
		d, ok := d.(*ast.GenDecl)
		if !ok || d.Tok != token.IMPORT {
			// Not an import declaration, so we're done.
			// Imports are always first.
			break
		}

		if len(d.Specs) == 0 {
			// Empty import block, remove it.
			f.Decls = append(f.Decls[:i], f.Decls[i+1:]...)
		}

		if !d.Lparen.IsValid() {
			// Not a block: sorted by default.
			continue
		}

		// Identify and sort runs of specs on successive lines.
		i := 0
		specs := d.Specs[:0]
		for j, s := range d.Specs {
			if j > i && fset.Position(s.Pos()).Line > 1+fset.Position(d.Specs[j-1].End()).Line {
				// j begins a new run.  End this one.
				specs = append(specs, sortSpecs(fset, f, d.Specs[i:j])...)
				i = j
			}
		}
		specs = append(specs, sortSpecs(fset, f, d.Specs[i:])...)
		d.Specs = specs

		// Deduping can leave a blank line before the rparen; clean that up.
		if len(d.Specs) > 0 {
			lastSpec := d.Specs[len(d.Specs)-1]
			lastLine := fset.Position(lastSpec.Pos()).Line
			if rParenLine := fset.Position(d.Rparen).Line; rParenLine > lastLine+1 {
				fset.File(d.Rparen).MergeLine(rParenLine - 1)
			}
		}
	}
}

func importPath(s ast.Spec) string {
	t, err := strconv.Unquote(s.(*ast.ImportSpec).Path.Value)
	if err == nil {
		return t
	}
	return ""
}

func importName(s ast.Spec) string {
	n := s.(*ast.ImportSpec).Name
	if n == nil {
		return ""
	}
	return n.Name
}

func importComment(s ast.Spec) string {
	c := s.(*ast.ImportSpec).Comment
	if c == nil {
		return ""
	}
	return c.Text()
}

// collapse indicates whether prev may be removed, leaving only next.
func collapse(prev, next ast.Spec) bool {
	if importPath(next) != importPath(prev) || importName(next) != importName(prev) {
		return false
	}
	return prev.(*ast.ImportSpec).Comment == nil
}

type posSpan struct {
	Start token.Pos
	End   token.Pos
}

func sortSpecs(fset *token.FileSet, f *ast.File, specs []ast.Spec) []ast.Spec {
	// Can't short-circuit here even if specs are already sorted,
	// since they might yet need deduplication.
	// A lone import, however, may be safely ignored.
	if len(specs) <= 1 {
		return specs
	}

	// Record positions for specs.
	pos := make([]posSpan, len(specs))
	for i, s := range specs {
		pos[i] = posSpan{s.Pos(), s.End()}
	}

	// Identify comments in this range.
	// Any comment from pos[0].Start to the final line counts.
	lastLine := fset.Position(pos[len(pos)-1].End).Line
	cstart := len(f.Comments)
	cend := len(f.Comments)
	for i, g := range f.Comments {
		if g.Pos() < pos[0].Start {
			continue
		}
		if i < cstart {
			cstart = i
		}
		if fset.Position(g.End()).Line > lastLine {
			cend = i
			break
		}
	}
	comments := f.Comments[cstart:cend]

	// Assign each comment to the import spec preceding it.
	importComment := map[*ast.ImportSpec][]*ast.CommentGroup{}
	specIndex := 0
	for _, g := range comments {
		for specIndex+1 < len(specs) && pos[specIndex+1].Start <= g.Pos() {
			specIndex++
		}
		s := specs[specIndex].(*ast.ImportSpec)
		importComment[s] = append(importComment[s], g)
	}

	// Sort the import specs by import path.
	// Remove duplicates, when possible without data loss.
	// Reassign the import paths to have the same position sequence.
	// Reassign each comment to abut the end of its spec.
	// Sort the comments by new position.
	sort.Sort(byImportSpec(specs))

	// Dedup. Thanks to our sorting, we can just consider
	// adjacent pairs of imports.
	deduped := specs[:0]
	for i, s := range specs {
		if i == len(specs)-1 || !collapse(s, specs[i+1]) {
			deduped = append(deduped, s)
		} else {
			p := s.Pos()
			fset.File(p).MergeLine(fset.Position(p).Line)
		}
	}
	specs = deduped

	// Fix up comment positions
	for i, s := range specs {
		s := s.(*ast.ImportSpec)
		if s.Name != nil {
			s.Name.NamePos = pos[i].Start
		}
		s.Path.ValuePos = pos[i].Start
		s.EndPos = pos[i].End
		for _, g := range importComment[s] {
			for _, c := range g.List {
				c.Slash = pos[i].End
			}
		}
	}

	sort.Sort(byCommentPos(comments))

	return specs
}

type byImportSpec []ast.Spec // slice of *ast.ImportSpec

func (x byImportSpec) Len() int      { return len(x) }
func (x byImportSpec) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x byImportSpec) Less(i, j int) bool {
	ipath := importPath(x[i])
	jpath := importPath(x[j])

	igroup := importGroup(ipath)
	jgroup := importGroup(jpath)
	if igroup != jgroup {
		return igroup < jgroup
	}

	if ipath != jpath {
		return ipath < jpath
	}
	iname := importName(x[i])
	jname := importName(x[j])

	if iname != jname {
		return iname < jname
	}
	return importComment(x[i]) < importComment(x[j])
}

type byCommentPos []*ast.CommentGroup

func (x byCommentPos) Len() int           { return len(x) }
func (x byCommentPos) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x byCommentPos) Less(i, j int) bool { return x[i].Pos() < x[j].Pos() }