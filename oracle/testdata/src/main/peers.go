// antha-tools/oracle/testdata/src/main/peers.go: Part of the Antha language
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

package peers

// Tests of channel 'peers' query.
// See go.tools/oracle/oracle_test.go for explanation.
// See peers.golden for expected query results.

var a2 int

func main() {
	chA := make(chan *int)
	a1 := 1
	chA <- &a1

	chA2 := make(chan *int, 2)
	if a2 == 0 {
		chA = chA2
	}

	chB := make(chan *int)
	b := 3
	chB <- &b

	<-chA  // @pointsto pointsto-chA "chA"
	<-chA2 // @pointsto pointsto-chA2 "chA2"
	<-chB  // @pointsto pointsto-chB "chB"

	select {
	case rA := <-chA: // @peers peer-recv-chA "<-"
		_ = rA // @pointsto pointsto-rA "rA"
	case rB := <-chB: // @peers peer-recv-chB "<-"
		_ = rB // @pointsto pointsto-rB "rB"

	case <-chA: // @peers peer-recv-chA' "<-"

	case chA2 <- &a2: // @peers peer-send-chA' "<-"
	}

	for _ = range chA {
	}
}