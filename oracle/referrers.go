// antha-tools/oracle/referrers.go: Part of the Antha language
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


package oracle

import (
	"fmt"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/token"
	"sort"

	"github.com/antha-lang/antha-tools/antha/types"
	"github.com/antha-lang/antha-tools/oracle/serial"
)

// Referrers reports all identifiers that resolve to the same object
// as the queried identifier, within any package in the analysis scope.
//
func referrers(o *Oracle, qpos *QueryPos) (queryResult, error) {
	id, _ := qpos.path[0].(*ast.Ident)
	if id == nil {
		return nil, fmt.Errorf("no identifier here")
	}

	obj := qpos.info.ObjectOf(id)
	if obj == nil {
		// Happens for y in "switch y := x.(type)", but I think that's all.
		return nil, fmt.Errorf("no object for identifier")
	}

	// Iterate over all antha/types' Uses facts for the entire program.
	var refs []*ast.Ident
	for _, info := range o.typeInfo {
		for id2, obj2 := range info.Uses {
			if sameObj(obj, obj2) {
				refs = append(refs, id2)
			}
		}
	}
	sort.Sort(byNamePos(refs))

	return &referrersResult{
		query: id,
		obj:   obj,
		refs:  refs,
	}, nil
}

// same reports whether x and y are identical, or both are PkgNames
// referring to the same Package.
//
func sameObj(x, y types.Object) bool {
	if x == y {
		return true
	}
	if _, ok := x.(*types.PkgName); ok {
		if _, ok := y.(*types.PkgName); ok {
			return x.Pkg() == y.Pkg()
		}
	}
	return false
}

// -------- utils --------

type byNamePos []*ast.Ident

func (p byNamePos) Len() int           { return len(p) }
func (p byNamePos) Less(i, j int) bool { return p[i].NamePos < p[j].NamePos }
func (p byNamePos) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type referrersResult struct {
	query *ast.Ident   // identifier of query
	obj   types.Object // object it denotes
	refs  []*ast.Ident // set of all other references to it
}

func (r *referrersResult) display(printf printfFunc) {
	if r.query.Pos() != r.obj.Pos() {
		printf(r.query, "reference to %s", r.obj.Name())
	}
	// TODO(adonovan): pretty-print object using same logic as
	// (*describeValueResult).display.
	printf(r.obj, "defined here as %s", r.obj)
	for _, ref := range r.refs {
		if r.query != ref {
			printf(ref, "referenced here")
		}
	}
}

// TODO(adonovan): encode extent, not just Pos info, in Serial form.

func (r *referrersResult) toSerial(res *serial.Result, fset *token.FileSet) {
	referrers := &serial.Referrers{
		Pos:  fset.Position(r.query.Pos()).String(),
		Desc: r.obj.String(),
	}
	if pos := r.obj.Pos(); pos != token.NoPos { // Package objects have no Pos()
		referrers.ObjPos = fset.Position(pos).String()
	}
	for _, ref := range r.refs {
		referrers.Refs = append(referrers.Refs, fset.Position(ref.NamePos).String())
	}
	res.Referrers = referrers
}