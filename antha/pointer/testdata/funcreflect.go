// antha-tools/antha/pointer/testdata/funcreflect.go: Part of the Antha language
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

var zero, a, b int
var false2 bool

func f(p *int, q hasF) *int {
	print(p)      // @pointsto main.a
	print(q)      // @types *T
	print(q.(*T)) // @pointsto new@newT1:22
	return &b
}

func g(p *bool) (*int, *bool, hasF) {
	return &b, p, new(T) // @line newT2
}

func reflectValueCall() {
	rvf := reflect.ValueOf(f)
	res := rvf.Call([]reflect.Value{
		// argument order is not significant:
		reflect.ValueOf(new(T)), // @line newT1
		reflect.ValueOf(&a),
	})
	print(res[0].Interface())        // @types *int
	print(res[0].Interface().(*int)) // @pointsto main.b
}

// @calls main.reflectValueCall -> main.f

func reflectValueCallIndirect() {
	rvf := reflect.ValueOf(g)
	call := rvf.Call // kids, don't try this at home

	// Indirect call uses shared contour.
	//
	// Also notice that argument position doesn't matter, and args
	// of inappropriate type (e.g. 'a') are ignored.
	res := call([]reflect.Value{
		reflect.ValueOf(&a),
		reflect.ValueOf(&false2),
	})
	res0 := res[0].Interface()
	print(res0)         // @types *int | *bool | *T
	print(res0.(*int))  // @pointsto main.b
	print(res0.(*bool)) // @pointsto main.false2
	print(res0.(hasF))  // @types *T
	print(res0.(*T))    // @pointsto new@newT2:19
}

// @calls main.reflectValueCallIndirect -> bound$(reflect.Value).Call
// @calls bound$(reflect.Value).Call -> main.g

func reflectTypeInOut() {
	var f func(float64, bool) (string, int)
	print(reflect.Zero(reflect.TypeOf(f).In(0)).Interface())    // @types float64
	print(reflect.Zero(reflect.TypeOf(f).In(1)).Interface())    // @types bool
	print(reflect.Zero(reflect.TypeOf(f).In(-1)).Interface())   // @types float64 | bool
	print(reflect.Zero(reflect.TypeOf(f).In(zero)).Interface()) // @types float64 | bool

	print(reflect.Zero(reflect.TypeOf(f).Out(0)).Interface()) // @types string
	print(reflect.Zero(reflect.TypeOf(f).Out(1)).Interface()) // @types int
	print(reflect.Zero(reflect.TypeOf(f).Out(2)).Interface()) // @types

	print(reflect.Zero(reflect.TypeOf(3).Out(0)).Interface()) // @types
}

type hasF interface {
	F()
}

type T struct{}

func (T) F()    {}
func (T) g(int) {}

type U struct{}

func (U) F(int)    {}
func (U) g(string) {}

type I interface {
	f()
}

var nonconst string

func reflectTypeMethodByName() {
	TU := reflect.TypeOf([]interface{}{T{}, U{}}[0])
	print(reflect.Zero(TU)) // @types T | U

	F, _ := TU.MethodByName("F")
	print(reflect.Zero(F.Type)) // @types func(T) | func(U, int)
	print(F.Func)               // @pointsto (main.T).F | (main.U).F

	g, _ := TU.MethodByName("g")
	print(reflect.Zero(g.Type)) // @types func(T, int) | func(U, string)
	print(g.Func)               // @pointsto (main.T).g | (main.U).g

	// Non-literal method names are treated less precisely.
	U := reflect.TypeOf(U{})
	X, _ := U.MethodByName(nonconst)
	print(reflect.Zero(X.Type)) // @types func(U, int) | func(U, string)
	print(X.Func)               // @pointsto (main.U).F | (main.U).g

	// Interface methods.
	rThasF := reflect.TypeOf(new(hasF)).Elem()
	print(reflect.Zero(rThasF)) // @types hasF
	F2, _ := rThasF.MethodByName("F")
	print(reflect.Zero(F2.Type)) // @types func()
	print(F2.Func)               // @pointsto (main.hasF).F

}

func reflectTypeMethod() {
	m := reflect.TypeOf(T{}).Method(0)
	print(reflect.Zero(m.Type)) // @types func(T) | func(T, int)
	print(m.Func)               // @pointsto (main.T).F | (main.T).g
}

func main() {
	reflectValueCall()
	reflectValueCallIndirect()
	reflectTypeInOut()
	reflectTypeMethodByName()
	reflectTypeMethod()
}