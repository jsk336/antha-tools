// antha-tools/oracle/callgraph.go: Part of the Antha language
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
	"sort"

	"github.com/antha-lang/antha-tools/antha/callgraph"
	"github.com/antha-lang/antha-tools/antha/ssa"
	"github.com/antha-lang/antha-tools/antha/types"
	"github.com/antha-lang/antha-tools/oracle/serial"
)

// doCallgraph displays the entire callgraph of the current program,
// or if a query -pos was provided, the query package.
func doCallgraph(o *Oracle, qpos *QueryPos) (queryResult, error) {
	buildSSA(o)

	// Run the pointer analysis and build the callgraph.
	o.ptaConfig.BuildCallGraph = true
	cg := ptrAnalysis(o).CallGraph
	cg.DeleteSyntheticNodes()

	var qpkg *types.Package
	var roots []*callgraph.Node
	if qpos == nil {
		// No -pos provided: show complete callgraph.
		roots = append(roots, cg.Root)

	} else {
		// A query -pos was provided: restrict result to
		// functions belonging to the query package.
		qpkg = qpos.info.Pkg
		isQueryPkg := func(fn *ssa.Function) bool {
			return fn.Pkg != nil && fn.Pkg.Object == qpkg
		}

		// First compute the nodes to keep and remove.
		var nodes, remove []*callgraph.Node
		for fn, cgn := range cg.Nodes {
			if isQueryPkg(fn) {
				nodes = append(nodes, cgn)
			} else {
				remove = append(remove, cgn)
			}
		}

		// Compact the Node.ID sequence of the remaining
		// nodes, preserving the original order.
		sort.Sort(nodesByID(nodes))
		for i, cgn := range nodes {
			cgn.ID = i
		}

		// Compute the set of roots:
		// in-package nodes with out-of-package callers.
		// For determinism, roots are ordered by original Node.ID.
		for _, cgn := range nodes {
			for _, e := range cgn.In {
				if !isQueryPkg(e.Caller.Func) {
					roots = append(roots, cgn)
					break
				}
			}
		}

		// Finally, discard all out-of-package nodes.
		for _, cgn := range remove {
			cg.DeleteNode(cgn)
		}
	}

	return &callgraphResult{qpkg, cg.Nodes, roots}, nil
}

type callgraphResult struct {
	qpkg  *types.Package
	nodes map[*ssa.Function]*callgraph.Node
	roots []*callgraph.Node
}

func (r *callgraphResult) display(printf printfFunc) {
	descr := "the entire program"
	if r.qpkg != nil {
		descr = fmt.Sprintf("package %s", r.qpkg.Path())
	}

	printf(nil, `
Below is a call graph of %s.
The numbered nodes form a spanning tree.
Non-numbered nodes indicate back- or cross-edges to the node whose
 number follows in parentheses.
`, descr)

	printed := make(map[*callgraph.Node]int)
	var print func(caller *callgraph.Node, indent int)
	print = func(caller *callgraph.Node, indent int) {
		if num, ok := printed[caller]; !ok {
			num = len(printed)
			printed[caller] = num

			// Sort the children into name order for deterministic* output.
			// (*mostly: anon funcs' names are not globally unique.)
			var funcs funcsByName
			for callee := range callgraph.CalleesOf(caller) {
				funcs = append(funcs, callee.Func)
			}
			sort.Sort(funcs)

			printf(caller.Func, "%d\t%*s%s", num, 4*indent, "", caller.Func.RelString(r.qpkg))
			for _, callee := range funcs {
				print(r.nodes[callee], indent+1)
			}
		} else {
			printf(caller.Func, "\t%*s%s (%d)", 4*indent, "", caller.Func.RelString(r.qpkg), num)
		}
	}
	for _, root := range r.roots {
		print(root, 0)
	}
}

type nodesByID []*callgraph.Node

func (s nodesByID) Len() int           { return len(s) }
func (s nodesByID) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s nodesByID) Less(i, j int) bool { return s[i].ID < s[j].ID }

type funcsByName []*ssa.Function

func (s funcsByName) Len() int           { return len(s) }
func (s funcsByName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s funcsByName) Less(i, j int) bool { return s[i].String() < s[j].String() }

func (r *callgraphResult) toSerial(res *serial.Result, fset *token.FileSet) {
	cg := make([]serial.CallGraph, len(r.nodes))
	for _, n := range r.nodes {
		j := &cg[n.ID]
		fn := n.Func
		j.Name = fn.String()
		j.Pos = fset.Position(fn.Pos()).String()
		for callee := range callgraph.CalleesOf(n) {
			j.Children = append(j.Children, callee.ID)
		}
		sort.Ints(j.Children)
	}
	res.Callgraph = cg
}