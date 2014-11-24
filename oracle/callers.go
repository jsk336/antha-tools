// antha-tools/oracle/callers.go: Part of the Antha language
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
	"github.com/antha-lang/antha/token"

	"github.com/antha-lang/antha-tools/antha/callgraph"
	"github.com/antha-lang/antha-tools/antha/ssa"
	"github.com/antha-lang/antha-tools/oracle/serial"
)

// Callers reports the possible callers of the function
// immediately enclosing the specified source location.
//
func callers(o *Oracle, qpos *QueryPos) (queryResult, error) {
	pkg := o.prog.Package(qpos.info.Pkg)
	if pkg == nil {
		return nil, fmt.Errorf("no SSA package")
	}
	if !ssa.HasEnclosingFunction(pkg, qpos.path) {
		return nil, fmt.Errorf("this position is not inside a function")
	}

	buildSSA(o)

	target := ssa.EnclosingFunction(pkg, qpos.path)
	if target == nil {
		return nil, fmt.Errorf("no SSA function built for this location (dead code?)")
	}

	// Run the pointer analysis, recording each
	// call found to originate from target.
	o.ptaConfig.BuildCallGraph = true
	cg := ptrAnalysis(o).CallGraph
	cg.DeleteSyntheticNodes()
	edges := cg.CreateNode(target).In
	// TODO(adonovan): sort + dedup calls to ensure test determinism.

	return &callersResult{
		target:    target,
		callgraph: cg,
		edges:     edges,
	}, nil
}

type callersResult struct {
	target    *ssa.Function
	callgraph *callgraph.Graph
	edges     []*callgraph.Edge
}

func (r *callersResult) display(printf printfFunc) {
	root := r.callgraph.Root
	if r.edges == nil {
		printf(r.target, "%s is not reachable in this program.", r.target)
	} else {
		printf(r.target, "%s is called from these %d sites:", r.target, len(r.edges))
		for _, edge := range r.edges {
			if edge.Caller == root {
				printf(r.target, "the root of the call graph")
			} else {
				printf(edge.Site, "\t%s from %s", edge.Site.Common().Description(), edge.Caller.Func)
			}
		}
	}
}

func (r *callersResult) toSerial(res *serial.Result, fset *token.FileSet) {
	root := r.callgraph.Root
	var callers []serial.Caller
	for _, edge := range r.edges {
		var c serial.Caller
		c.Caller = edge.Caller.Func.String()
		if edge.Caller == root {
			c.Desc = "synthetic call"
		} else {
			c.Pos = fset.Position(edge.Site.Pos()).String()
			c.Desc = edge.Site.Common().Description()
		}
		callers = append(callers, c)
	}
	res.Callers = callers
}