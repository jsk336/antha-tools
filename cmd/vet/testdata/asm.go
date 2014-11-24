// antha-tools/cmd/vet/testdata/asm.go: Part of the Antha language
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

// This file contains declarations to test the assembly in test_asm.s.

package testdata

func arg1(x int8, y uint8)
func arg2(x int16, y uint16)
func arg4(x int32, y uint32)
func arg8(x int64, y uint64)
func argint(x int, y uint)
func argptr(x *byte, y *byte, c chan int, m map[int]int, f func())
func argstring(x, y string)
func argslice(x, y []string)
func argiface(x interface{}, y interface {
	m()
})
func returnint() int
func returnbyte(x int) byte
func returnnamed(x byte) (r1 int, r2 int16, r3 string, r4 byte)

func noprof(x int)
func dupok(x int)
func nosplit(x int)
func rodata(x int)
func noptr(x int)
func wrapper(x int)