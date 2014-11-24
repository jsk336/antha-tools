// antha-tools/oracle/testdata/src/main/implements-json.go: Part of the Antha language
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

package main

// Tests of 'implements' query, -output=json.
// See go.tools/oracle/oracle_test.go for explanation.
// See implements.golden for expected query results.

func main() {
}

type E interface{} // @implements E "E"

type F interface { // @implements F "F"
	f()
}

type FG interface { // @implements FG "FG"
	f()
	g() []int // @implements slice "..int"
}

type C int // @implements C "C"
type D struct{}

func (c *C) f() {} // @implements starC ".C"
func (d D) f()  {} // @implements D "D"

func (d *D) g() []int { return nil } // @implements starD ".D"