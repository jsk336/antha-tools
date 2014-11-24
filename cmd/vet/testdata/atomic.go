// antha-tools/cmd/vet/testdata/atomic.go: Part of the Antha language
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


// This file contains tests for the atomic checker.

package testdata

import (
	"sync/atomic"
)

type Counter uint64

func AtomicTests() {
	x := uint64(1)
	x = atomic.AddUint64(&x, 1)        // ERROR "direct assignment to atomic value"
	_, x = 10, atomic.AddUint64(&x, 1) // ERROR "direct assignment to atomic value"
	x, _ = atomic.AddUint64(&x, 1), 10 // ERROR "direct assignment to atomic value"

	y := &x
	*y = atomic.AddUint64(y, 1) // ERROR "direct assignment to atomic value"

	var su struct{ Counter uint64 }
	su.Counter = atomic.AddUint64(&su.Counter, 1) // ERROR "direct assignment to atomic value"
	z1 := atomic.AddUint64(&su.Counter, 1)
	_ = z1 // Avoid err "z declared and not used"

	var sp struct{ Counter *uint64 }
	*sp.Counter = atomic.AddUint64(sp.Counter, 1) // ERROR "direct assignment to atomic value"
	z2 := atomic.AddUint64(sp.Counter, 1)
	_ = z2 // Avoid err "z declared and not used"

	au := []uint64{10, 20}
	au[0] = atomic.AddUint64(&au[0], 1) // ERROR "direct assignment to atomic value"
	au[1] = atomic.AddUint64(&au[0], 1)

	ap := []*uint64{&au[0], &au[1]}
	*ap[0] = atomic.AddUint64(ap[0], 1) // ERROR "direct assignment to atomic value"
	*ap[1] = atomic.AddUint64(ap[0], 1)

	x = atomic.AddUint64() // Used to make vet crash; now silently ignored.
}