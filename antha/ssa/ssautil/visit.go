// antha-tools/antha/ssa/ssautil/visit.go: Part of the Antha language
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


package ssautil

import "github.com/antha-lang/antha-tools/antha/ssa"

// This file defines utilities for visiting the SSA representation of
// a Program.
//
// TODO(adonovan): test coverage.

// AllFunctions finds and returns the set of functions potentially
// needed by program prog, as determined by a simple linker-style
// reachability algorithm starting from the members and method-sets of
// each package.  The result may include anonymous functions and
// synthetic wrappers.
//
// Precondition: all packages are built.
//
func AllFunctions(prog *ssa.Program) map[*ssa.Function]bool {
	visit := visitor{
		prog: prog,
		seen: make(map[*ssa.Function]bool),
	}
	visit.program()
	return visit.seen
}

type visitor struct {
	prog *ssa.Program
	seen map[*ssa.Function]bool
}

func (visit *visitor) program() {
	for _, pkg := range visit.prog.AllPackages() {
		for _, mem := range pkg.Members {
			if fn, ok := mem.(*ssa.Function); ok {
				visit.function(fn)
			}
		}
	}
	for _, T := range visit.prog.TypesWithMethodSets() {
		mset := visit.prog.MethodSets.MethodSet(T)
		for i, n := 0, mset.Len(); i < n; i++ {
			visit.function(visit.prog.Method(mset.At(i)))
		}
	}
}

func (visit *visitor) function(fn *ssa.Function) {
	if !visit.seen[fn] {
		visit.seen[fn] = true
		var buf [10]*ssa.Value // avoid alloc in common case
		for _, b := range fn.Blocks {
			for _, instr := range b.Instrs {
				for _, op := range instr.Operands(buf[:0]) {
					if fn, ok := (*op).(*ssa.Function); ok {
						visit.function(fn)
					}
				}
			}
		}
	}
}