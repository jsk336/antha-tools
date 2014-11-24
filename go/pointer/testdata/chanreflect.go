// antha-tools/go/pointer/testdata/chanreflect.go: Part of the Antha language
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

// Test of channels with reflection.

var a, b int

func chanreflect1() {
	ch := make(chan *int, 0) // @line cr1make
	crv := reflect.ValueOf(ch)
	crv.Send(reflect.ValueOf(&a))
	print(crv.Interface())             // @types chan *int
	print(crv.Interface().(chan *int)) // @pointsto makechan@cr1make:12
	print(<-ch)                        // @pointsto main.a
}

func chanreflect1i() {
	// Exercises reflect.Value conversions to/from interfaces:
	// a different code path than for concrete types.
	ch := make(chan interface{}, 0)
	reflect.ValueOf(ch).Send(reflect.ValueOf(&a))
	v := <-ch
	print(v)        // @types *int
	print(v.(*int)) // @pointsto main.a
}

func chanreflect2() {
	ch := make(chan *int, 0)
	ch <- &b
	crv := reflect.ValueOf(ch)
	r, _ := crv.Recv()
	print(r.Interface())        // @types *int
	print(r.Interface().(*int)) // @pointsto main.b
}

func chanOfRecv() {
	// MakeChan(<-chan) is a no-op.
	t := reflect.ChanOf(reflect.RecvDir, reflect.TypeOf(&a))
	print(reflect.Zero(t).Interface())                      // @types <-chan *int
	print(reflect.MakeChan(t, 0).Interface().(<-chan *int)) // @pointsto
	print(reflect.MakeChan(t, 0).Interface().(chan *int))   // @pointsto
}

func chanOfSend() {
	// MakeChan(chan<-) is a no-op.
	t := reflect.ChanOf(reflect.SendDir, reflect.TypeOf(&a))
	print(reflect.Zero(t).Interface())                      // @types chan<- *int
	print(reflect.MakeChan(t, 0).Interface().(chan<- *int)) // @pointsto
	print(reflect.MakeChan(t, 0).Interface().(chan *int))   // @pointsto
}

func chanOfBoth() {
	t := reflect.ChanOf(reflect.BothDir, reflect.TypeOf(&a))
	print(reflect.Zero(t).Interface()) // @types chan *int
	ch := reflect.MakeChan(t, 0)
	print(ch.Interface().(chan *int)) // @pointsto <alloc in reflect.MakeChan>
	ch.Send(reflect.ValueOf(&b))
	ch.Interface().(chan *int) <- &a
	r, _ := ch.Recv()
	print(r.Interface().(*int))         // @pointsto main.a | main.b
	print(<-ch.Interface().(chan *int)) // @pointsto main.a | main.b
}

var unknownDir reflect.ChanDir // not a constant

func chanOfUnknown() {
	// Unknown channel direction: assume all three.
	// MakeChan only works on the bi-di channel type.
	t := reflect.ChanOf(unknownDir, reflect.TypeOf(&a))
	print(reflect.Zero(t).Interface())        // @types <-chan *int | chan<- *int | chan *int
	print(reflect.MakeChan(t, 0).Interface()) // @types chan *int
}

func main() {
	chanreflect1()
	chanreflect1i()
	chanreflect2()
	chanOfRecv()
	chanOfSend()
	chanOfBoth()
	chanOfUnknown()
}