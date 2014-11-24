// antha-tools/go/pointer/testdata/channels.go: Part of the Antha language
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

func incr(x int) int { return x + 1 }

func decr(x int) int { return x - 1 }

var unknown bool // defeat dead-code elimination

func chan1() {
	chA := make(chan func(int) int, 0) // @line c1makeA
	chB := make(chan func(int) int, 0) // @line c1makeB
	chA <- incr
	chB <- decr
	chB <- func(int) int { return 1 }

	print(chA)   // @pointsto makechan@c1makeA:13
	print(<-chA) // @pointsto main.incr

	print(chB)   // @pointsto makechan@c1makeB:13
	print(<-chB) // @pointsto main.decr | chan1$1
}

func chan2() {
	chA := make(chan func(int) int, 0) // @line c2makeA
	chB := make(chan func(int) int, 0) // @line c2makeB
	chA <- incr
	chB <- decr
	chB <- func(int) int { return 1 }

	// Channels flow together.
	// Labelsets remain distinct but elements are merged.
	chAB := chA
	if unknown {
		chAB = chB
	}

	print(chA)   // @pointsto makechan@c2makeA:13
	print(<-chA) // @pointsto main.incr

	print(chB)   // @pointsto makechan@c2makeB:13
	print(<-chB) // @pointsto main.decr | chan2$1

	print(chAB)   // @pointsto makechan@c2makeA:13 | makechan@c2makeB:13
	print(<-chAB) // @pointsto main.incr | main.decr | chan2$1

	(<-chA)(3)
}

// @calls main.chan2 -> main.incr

func chan3() {
	chA := make(chan func(int) int, 0) // @line c3makeA
	chB := make(chan func(int) int, 0) // @line c3makeB
	chA <- incr
	chB <- decr
	chB <- func(int) int { return 1 }
	print(chA)   // @pointsto makechan@c3makeA:13
	print(<-chA) // @pointsto main.incr
	print(chB)   // @pointsto makechan@c3makeB:13
	print(<-chB) // @pointsto main.decr | chan3$1

	(<-chA)(3)
}

// @calls main.chan3 -> main.incr

func chan4() {
	chA := make(chan func(int) int, 0) // @line c4makeA
	chB := make(chan func(int) int, 0) // @line c4makeB

	select {
	case chA <- incr:
	case chB <- decr:
	case a := <-chA:
		print(a) // @pointsto main.incr
	case b := <-chB:
		print(b) // @pointsto main.decr
	default:
		print(chA) // @pointsto makechan@c4makeA:13
		print(chB) // @pointsto makechan@c4makeB:13
	}

	for k := range chA {
		print(k) // @pointsto main.incr
	}
	// Exercise constraint generation (regtest for a crash).
	for _ = range chA {
	}
}

// Multi-word channel value in select with multiple receive cases.
// (Regtest for a crash.)
func chan5() {
	type T struct {
		x *int
		y interface{}
	}
	ch := make(chan T)
	ch <- T{new(int), incr} // @line ch5new
	select {
	case a := <-ch:
		print(a.x) // @pointsto new@ch5new:13
		print(a.y) // @types func(x int) int
	case b := <-ch:
		print(b.x) // @pointsto new@ch5new:13
		print(b.y) // @types func(x int) int
	}
}

func main() {
	chan1()
	chan2()
	chan3()
	chan4()
	chan5()
}