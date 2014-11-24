// antha-tools/oracle/definition.go: Part of the Antha language
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

	"github.com/antha-lang/antha-tools/antha/types"
	"github.com/antha-lang/antha-tools/oracle/serial"
)

// definition reports the location of the definition of an identifier.
//
// TODO(adonovan): opt: for intra-file references, the parser's
// resolution might be enough; we should start with that.
//
func definition(o *Oracle, qpos *QueryPos) (queryResult, error) {
	id, _ := qpos.path[0].(*ast.Ident)
	if id == nil {
		return nil, fmt.Errorf("no identifier here")
	}

	obj := qpos.info.ObjectOf(id)
	if obj == nil {
		// Happens for y in "switch y := x.(type)", but I think that's all.
		return nil, fmt.Errorf("no object for identifier")
	}

	return &definitionResult{qpos, obj}, nil
}

type definitionResult struct {
	qpos *QueryPos
	obj  types.Object // object it denotes
}

func (r *definitionResult) display(printf printfFunc) {
	printf(r.obj, "defined here as %s", r.qpos.ObjectString(r.obj))
}

func (r *definitionResult) toSerial(res *serial.Result, fset *token.FileSet) {
	definition := &serial.Definition{
		Desc: r.obj.String(),
	}
	if pos := r.obj.Pos(); pos != token.NoPos { // Package objects have no Pos()
		definition.ObjPos = fset.Position(pos).String()
	}
	res.Definition = definition
}