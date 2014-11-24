// antha-tools/cmd/vet/testdata/nilfunc.go: Part of the Antha language
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


package testdata

func F() {}

type T struct {
	F func()
}

func (T) M() {}

var Fv = F

func Comparison() {
	var t T
	var fn func()
	if fn == nil || Fv == nil || t.F == nil {
		// no error; these func vars or fields may be nil
	}
	if F == nil { // ERROR "comparison of function F == nil is always false"
		panic("can't happen")
	}
	if t.M == nil { // ERROR "comparison of function M == nil is always false"
		panic("can't happen")
	}
	if F != nil { // ERROR "comparison of function F != nil is always true"
		if t.M != nil { // ERROR "comparison of function M != nil is always true"
			return
		}
	}
	panic("can't happen")
}