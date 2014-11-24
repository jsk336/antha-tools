// antha-tools/go/pointer/testdata/structreflect.go: Part of the Antha language
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

// +build ignore

package main

import "reflect"

type A struct {
	f *int
	g interface{}
	h bool
}

var dyn string

func reflectTypeFieldByName() {
	f, _ := reflect.TypeOf(A{}).FieldByName("f")
	print(f.Type) // @pointsto *int

	g, _ := reflect.TypeOf(A{}).FieldByName("g")
	print(g.Type)               // @pointsto interface{}
	print(reflect.Zero(g.Type)) // @pointsto <alloc in reflect.Zero>
	print(reflect.Zero(g.Type)) // @types interface{}

	print(reflect.Zero(g.Type).Interface()) // @pointsto
	print(reflect.Zero(g.Type).Interface()) // @types

	h, _ := reflect.TypeOf(A{}).FieldByName("h")
	print(h.Type) // @pointsto bool

	missing, _ := reflect.TypeOf(A{}).FieldByName("missing")
	print(missing.Type) // @pointsto

	dyn, _ := reflect.TypeOf(A{}).FieldByName(dyn)
	print(dyn.Type) // @pointsto *int | bool | interface{}
}

func reflectTypeField() {
	fld := reflect.TypeOf(A{}).Field(0)
	print(fld.Type) // @pointsto *int | bool | interface{}
}

func main() {
	reflectTypeFieldByName()
	reflectTypeField()
}