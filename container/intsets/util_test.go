// antha-tools/container/intsets/util_test.go: Part of the Antha language
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

import "testing"

func TestNLZ(t *testing.T) {
	// Test the platform-specific edge case.
	// NB: v must be a var (not const) so that the word() conversion is dynamic.
	// Otherwise the compiler will report an error.
	v := uint64(0x0000801000000000)
	n := nlz(word(v))
	want := 32 // (on 32-bit)
	if bitsPerWord == 64 {
		want = 16
	}
	if n != want {
		t.Errorf("%d-bit nlz(%d) = %d, want %d", bitsPerWord, v, n, want)
	}
}

// Backdoor for testing.
func (s *Sparse) Check() error { return s.check() }