// antha-tools/anthadoc/snippet.go: Part of the Antha language
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


// This file contains the infrastructure to create a code
// snippet for search results.
//
// Note: At the moment, this only creates HTML snippets.

package godoc

import (
	"bytes"
	"fmt"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/token"
)

type Snippet struct {
	Line int
	Text string // HTML-escaped
}

func (p *Presentation) newSnippet(fset *token.FileSet, decl ast.Decl, id *ast.Ident) *Snippet {
	// TODO instead of pretty-printing the node, should use the original source instead
	var buf1 bytes.Buffer
	p.writeNode(&buf1, fset, decl)
	// wrap text with <pre> tag
	var buf2 bytes.Buffer
	buf2.WriteString("<pre>")
	FormatText(&buf2, buf1.Bytes(), -1, true, id.Name, nil)
	buf2.WriteString("</pre>")
	return &Snippet{fset.Position(id.Pos()).Line, buf2.String()}
}

func findSpec(list []ast.Spec, id *ast.Ident) ast.Spec {
	for _, spec := range list {
		switch s := spec.(type) {
		case *ast.ImportSpec:
			if s.Name == id {
				return s
			}
		case *ast.ValueSpec:
			for _, n := range s.Names {
				if n == id {
					return s
				}
			}
		case *ast.TypeSpec:
			if s.Name == id {
				return s
			}
		}
	}
	return nil
}

func (p *Presentation) genSnippet(fset *token.FileSet, d *ast.GenDecl, id *ast.Ident) *Snippet {
	s := findSpec(d.Specs, id)
	if s == nil {
		return nil //  declaration doesn't contain id - exit gracefully
	}

	// only use the spec containing the id for the snippet
	dd := &ast.GenDecl{
		Doc:    d.Doc,
		TokPos: d.Pos(),
		Tok:    d.Tok,
		Lparen: d.Lparen,
		Specs:  []ast.Spec{s},
		Rparen: d.Rparen,
	}

	return p.newSnippet(fset, dd, id)
}

func (p *Presentation) funcSnippet(fset *token.FileSet, d *ast.FuncDecl, id *ast.Ident) *Snippet {
	if d.Name != id {
		return nil //  declaration doesn't contain id - exit gracefully
	}

	// only use the function signature for the snippet
	dd := &ast.FuncDecl{
		Doc:  d.Doc,
		Recv: d.Recv,
		Name: d.Name,
		Type: d.Type,
	}

	return p.newSnippet(fset, dd, id)
}

// NewSnippet creates a text snippet from a declaration decl containing an
// identifier id. Parts of the declaration not containing the identifier
// may be removed for a more compact snippet.
func NewSnippet(fset *token.FileSet, decl ast.Decl, id *ast.Ident) *Snippet {
	// TODO(bradfitz, adg): remove this function.  But it's used by indexer, which
	// doesn't have a *Presentation, and NewSnippet needs a TabWidth.
	var p Presentation
	p.TabWidth = 4
	return p.NewSnippet(fset, decl, id)
}

// NewSnippet creates a text snippet from a declaration decl containing an
// identifier id. Parts of the declaration not containing the identifier
// may be removed for a more compact snippet.
func (p *Presentation) NewSnippet(fset *token.FileSet, decl ast.Decl, id *ast.Ident) *Snippet {
	var s *Snippet
	switch d := decl.(type) {
	case *ast.GenDecl:
		s = p.genSnippet(fset, d, id)
	case *ast.FuncDecl:
		s = p.funcSnippet(fset, d, id)
	}

	// handle failure gracefully
	if s == nil {
		var buf bytes.Buffer
		fmt.Fprintf(&buf, `<span class="alert">could not generate a snippet for <span class="highlight">%s</span></span>`, id.Name)
		s = &Snippet{fset.Position(id.Pos()).Line, buf.String()}
	}
	return s
}