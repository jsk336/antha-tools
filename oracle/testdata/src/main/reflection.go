// antha-tools/oracle/testdata/src/main/reflection.go: Part of the Antha language
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

package reflection

// This is a test of 'pointsto', but we split it into a separate file
// so that pointsto.go doesn't have to import "reflect" each time.

import "reflect"

var a int
var b bool

func main() {
	m := make(map[*int]*bool)
	m[&a] = &b

	mrv := reflect.ValueOf(m)
	if a > 0 {
		mrv = reflect.ValueOf(&b)
	}
	if a > 0 {
		mrv = reflect.ValueOf(&a)
	}

	_ = mrv                  // @pointsto mrv "mrv"
	p1 := mrv.Interface()    // @pointsto p1 "p1"
	p2 := mrv.MapKeys()      // @pointsto p2 "p2"
	p3 := p2[0]              // @pointsto p3 "p3"
	p4 := reflect.TypeOf(p1) // @pointsto p4 "p4"

	_, _, _, _ = p1, p2, p3, p4
}