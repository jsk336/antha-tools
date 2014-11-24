// antha-tools/go/ssa/lvalue.go: Part of the Antha language
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


package ssa

// lvalues are the union of addressable expressions and map-index
// expressions.

import (
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/token"

	"github.com/antha-lang/antha-tools/antha/types"
)

// An lvalue represents an assignable location that may appear on the
// left-hand side of an assignment.  This is a generalization of a
// pointer to permit updates to elements of maps.
//
type lvalue interface {
	store(fn *Function, v Value) // stores v into the location
	load(fn *Function) Value     // loads the contents of the location
	address(fn *Function) Value  // address of the location
	typ() types.Type             // returns the type of the location
}

// An address is an lvalue represented by a true pointer.
type address struct {
	addr    Value
	starPos token.Pos // source position, if from explicit *addr
	expr    ast.Expr  // source syntax [debug mode]
}

func (a *address) load(fn *Function) Value {
	load := emitLoad(fn, a.addr)
	load.pos = a.starPos
	return load
}

func (a *address) store(fn *Function, v Value) {
	store := emitStore(fn, a.addr, v)
	store.pos = a.starPos
	if a.expr != nil {
		// store.Val is v, converted for assignability.
		emitDebugRef(fn, a.expr, store.Val, false)
	}
}

func (a *address) address(fn *Function) Value {
	if a.expr != nil {
		emitDebugRef(fn, a.expr, a.addr, true)
	}
	return a.addr
}

func (a *address) typ() types.Type {
	return deref(a.addr.Type())
}

// An element is an lvalue represented by m[k], the location of an
// element of a map or string.  These locations are not addressable
// since pointers cannot be formed from them, but they do support
// load(), and in the case of maps, store().
//
type element struct {
	m, k Value      // map or string
	t    types.Type // map element type or string byte type
	pos  token.Pos  // source position of colon ({k:v}) or lbrack (m[k]=v)
}

func (e *element) load(fn *Function) Value {
	l := &Lookup{
		X:     e.m,
		Index: e.k,
	}
	l.setPos(e.pos)
	l.setType(e.t)
	return fn.emit(l)
}

func (e *element) store(fn *Function, v Value) {
	up := &MapUpdate{
		Map:   e.m,
		Key:   e.k,
		Value: emitConv(fn, v, e.t),
	}
	up.pos = e.pos
	fn.emit(up)
}

func (e *element) address(fn *Function) Value {
	panic("map/string elements are not addressable")
}

func (e *element) typ() types.Type {
	return e.t
}

// A blank is a dummy variable whose name is "_".
// It is not reified: loads are illegal and stores are ignored.
//
type blank struct{}

func (bl blank) load(fn *Function) Value {
	panic("blank.load is illegal")
}

func (bl blank) store(fn *Function, v Value) {
	// no-op
}

func (bl blank) address(fn *Function) Value {
	panic("blank var is not addressable")
}

func (bl blank) typ() types.Type {
	// This should be the type of the blank Ident; the typechecker
	// doesn't provide this yet, but fortunately, we don't need it
	// yet either.
	panic("blank.typ is unimplemented")
}