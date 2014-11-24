// antha-tools/go/pointer/example_test.go: Part of the Antha language
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


package pointer_test

import (
	"fmt"
	"sort"

	"github.com/antha-lang/antha-tools/antha/callgraph"
	"github.com/antha-lang/antha-tools/antha/loader"
	"github.com/antha-lang/antha-tools/antha/pointer"
	"github.com/antha-lang/antha-tools/antha/ssa"
)

// This program demonstrates how to use the pointer analysis to
// obtain a conservative call-graph of a Go program.
// It also shows how to compute the points-to set of a variable,
// in this case, (C).f's ch parameter.
//
func Example() {
	const myprog = `
package main

import "fmt"

type I interface {
	f(map[string]int)
}

type C struct{}

func (C) f(m map[string]int) {
	fmt.Println("C.f()")
}

func main() {
	var i I = C{}
	x := map[string]int{"one":1}
	i.f(x) // dynamic method call
}
`
	// Construct a loader.
	conf := loader.Config{SourceImports: true}

	// Parse the input file.
	file, err := conf.ParseFile("myprog.go", myprog)
	if err != nil {
		fmt.Print(err) // parse error
		return
	}

	// Create single-file main package and import its dependencies.
	conf.CreateFromFiles("main", file)

	iprog, err := conf.Load()
	if err != nil {
		fmt.Print(err) // type error in some package
		return
	}

	// Create SSA-form program representation.
	prog := ssa.Create(iprog, 0)
	mainPkg := prog.Package(iprog.Created[0].Pkg)

	// Build SSA code for bodies of all functions in the whole program.
	prog.BuildAll()

	// Configure the pointer analysis to build a call-graph.
	config := &pointer.Config{
		Mains:          []*ssa.Package{mainPkg},
		BuildCallGraph: true,
	}

	// Query points-to set of (C).f's parameter m, a map.
	C := mainPkg.Type("C").Type()
	Cfm := prog.LookupMethod(C, mainPkg.Object, "f").Params[1]
	config.AddQuery(Cfm)

	// Run the pointer analysis.
	result, err := pointer.Analyze(config)
	if err != nil {
		panic(err) // internal error in pointer analysis
	}

	// Find edges originating from the main package.
	// By converting to strings, we de-duplicate nodes
	// representing the same function due to context sensitivity.
	var edges []string
	callgraph.GraphVisitEdges(result.CallGraph, func(edge *callgraph.Edge) error {
		caller := edge.Caller.Func
		if caller.Pkg == mainPkg {
			edges = append(edges, fmt.Sprint(caller, " --> ", edge.Callee.Func))
		}
		return nil
	})

	// Print the edges in sorted order.
	sort.Strings(edges)
	for _, edge := range edges {
		fmt.Println(edge)
	}
	fmt.Println()

	// Print the labels of (C).f(m)'s points-to set.
	fmt.Println("m may point to:")
	var labels []string
	for _, l := range result.Queries[Cfm].PointsTo().Labels() {
		label := fmt.Sprintf("  %s: %s", prog.Fset.Position(l.Pos()), l)
		labels = append(labels, label)
	}
	sort.Strings(labels)
	for _, label := range labels {
		fmt.Println(label)
	}

	// Output:
	// (main.C).f --> fmt.Println
	// main.init --> fmt.init
	// main.main --> (main.C).f
	//
	// m may point to:
	//   myprog.go:18:21: makemap
}