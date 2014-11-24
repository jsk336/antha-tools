// antha-tools/oracle/callstack.go: Part of the Antha language
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

// Callstack displays an arbitrary path from a root of the callgraph
// to the function at the current position.
//
// The information may be misleading in a context-insensitive
// analysis. e.g. the call path X->Y->Z might be infeasible if Y never
// calls Z when it is called from X.  TODO(adonovan): think about UI.
//
// TODO(adonovan): permit user to specify a starting point other than
// the analysis root.
//
func callstack(o *Oracle, qpos *QueryPos) (queryResult, error) {
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

	// Run the pointer analysis and build the complete call graph.
	o.ptaConfig.BuildCallGraph = true
	cg := ptrAnalysis(o).CallGraph
	cg.DeleteSyntheticNodes()

	// Search for an arbitrary path from a root to the target function.
	isEnd := func(n *callgraph.Node) bool { return n.Func == target }
	callpath := callgraph.PathSearch(cg.Root, isEnd)
	if callpath != nil {
		callpath = callpath[1:] // remove synthetic edge from <root>
	}

	return &callstackResult{
		qpos:     qpos,
		target:   target,
		callpath: callpath,
	}, nil
}

type callstackResult struct {
	qpos     *QueryPos
	target   *ssa.Function
	callpath []*callgraph.Edge
}

func (r *callstackResult) display(printf printfFunc) {
	if r.callpath != nil {
		printf(r.qpos, "Found a call path from root to %s", r.target)
		printf(r.target, "%s", r.target)
		for i := len(r.callpath) - 1; i >= 0; i-- {
			edge := r.callpath[i]
			printf(edge.Site, "%s from %s", edge.Site.Common().Description(), edge.Caller.Func)
		}
	} else {
		printf(r.target, "%s is unreachable in this analysis scope", r.target)
	}
}

func (r *callstackResult) toSerial(res *serial.Result, fset *token.FileSet) {
	var callers []serial.Caller
	for i := len(r.callpath) - 1; i >= 0; i-- { // (innermost first)
		edge := r.callpath[i]
		callers = append(callers, serial.Caller{
			Pos:    fset.Position(edge.Site.Pos()).String(),
			Caller: edge.Caller.Func.String(),
			Desc:   edge.Site.Common().Description(),
		})
	}
	res.Callstack = &serial.CallStack{
		Pos:     fset.Position(r.target.Pos()).String(),
		Target:  r.target.String(),
		Callers: callers,
	}
}