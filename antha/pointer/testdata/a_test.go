// antha-tools/antha/pointer/testdata/a_test.go: Part of the Antha language
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

package a

// This test exercises the synthesis of testmain packages for tests.
// The test framework doesn't directly let us perform negative
// assertions (i.e. that TestingQuux isn't called, or that its
// parameter's PTS is empty) so this test is rather roundabout.

import "testing"

func log(f func(*testing.T)) {
	// The PTS of f is the set of called tests.  TestingQuux is not present.
	print(f) // @pointsto main.Test | main.TestFoo
}

func Test(t *testing.T) {
	// Don't assert @pointsto(t) since its label contains a fragile line number.
	log(Test)
}

func TestFoo(t *testing.T) {
	// Don't assert @pointsto(t) since its label contains a fragile line number.
	log(TestFoo)
}

func TestingQuux(t *testing.T) {
	// We can't assert @pointsto(t) since this is dead code.
	log(TestingQuux)
}

func BenchmarkFoo(b *testing.B) {
}

func ExampleBar() {
}

// Excludes TestingQuux.
// @calls testing.tRunner -> main.Test
// @calls testing.tRunner -> main.TestFoo
// @calls testing.runExample -> main.ExampleBar
// @calls (*testing.B).runN -> main.BenchmarkFoo