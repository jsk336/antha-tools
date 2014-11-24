// antha-tools/go/types/typeutil/ui.go: Part of the Antha language
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


package typeutil

// This file defines utilities for user interfaces that display types.

import "github.com/antha-lang/antha-tools/antha/types"

// IntuitiveMethodSet returns the intuitive method set of a type, T.
//
// The result contains MethodSet(T) and additionally, if T is a
// concrete type, methods belonging to *T if there is no identically
// named method on T itself.  This corresponds to user intuition about
// method sets; this function is intended only for user interfaces.
//
// The order of the result is as for types.MethodSet(T).
//
func IntuitiveMethodSet(T types.Type, msets *types.MethodSetCache) []*types.Selection {
	var result []*types.Selection
	mset := msets.MethodSet(T)
	if _, ok := T.Underlying().(*types.Interface); ok {
		for i, n := 0, mset.Len(); i < n; i++ {
			result = append(result, mset.At(i))
		}
	} else {
		pmset := msets.MethodSet(types.NewPointer(T))
		for i, n := 0, pmset.Len(); i < n; i++ {
			meth := pmset.At(i)
			if m := mset.Lookup(meth.Obj().Pkg(), meth.Obj().Name()); m != nil {
				meth = m
			}
			result = append(result, meth)
		}
	}
	return result
}