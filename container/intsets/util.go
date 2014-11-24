// antha-tools/container/intsets/util.go: Part of the Antha language
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


package intsets

var a [1 << 8]byte

func init() {
	for i := range a {
		var n byte
		for x := i; x != 0; x >>= 1 {
			if x&1 != 0 {
				n++
			}
		}
		a[i] = n
	}
}

// popcount returns the population count (number of set bits) of x.
func popcount(x word) int {
	return int(a[byte(x>>(0*8))] +
		a[byte(x>>(1*8))] +
		a[byte(x>>(2*8))] +
		a[byte(x>>(3*8))] +
		a[byte(x>>(4*8))] +
		a[byte(x>>(5*8))] +
		a[byte(x>>(6*8))] +
		a[byte(x>>(7*8))])
}

// nlz returns the number of leading zeros of x.
// From Hacker's Delight, fig 5.11.
func nlz(x word) int {
	x |= (x >> 1)
	x |= (x >> 2)
	x |= (x >> 4)
	x |= (x >> 8)
	x |= (x >> 16)
	x |= (x >> 32)
	return popcount(^x)
}

// ntz returns the number of trailing zeros of x.
// From Hacker's Delight, fig 5.13.
func ntz(x word) int {
	if x == 0 {
		return bitsPerWord
	}
	n := 1
	if bitsPerWord == 64 {
		if (x & 0xffffffff) == 0 {
			n = n + 32
			x = x >> 32
		}
	}
	if (x & 0x0000ffff) == 0 {
		n = n + 16
		x = x >> 16
	}
	if (x & 0x000000ff) == 0 {
		n = n + 8
		x = x >> 8
	}
	if (x & 0x0000000f) == 0 {
		n = n + 4
		x = x >> 4
	}
	if (x & 0x00000003) == 0 {
		n = n + 2
		x = x >> 2
	}
	return n - int(x&1)
}