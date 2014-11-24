// antha-tools/refactor/eg/testdata/A1.go: Part of the Antha language
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

package A1

import (
	. "fmt"
	myfmt "fmt"
	"os"
	"strings"
)

func example(n int) {
	x := "foo" + strings.Repeat("\t", n)
	// Match, despite named import.
	myfmt.Errorf("%s", x)

	// Match, despite dot import.
	Errorf("%s", x)

	// Match: multiple matches in same function are possible.
	myfmt.Errorf("%s", x)

	// No match: wildcarded operand has the wrong type.
	myfmt.Errorf("%s", 3)

	// No match: function operand doesn't match.
	myfmt.Printf("%s", x)

	// No match again, dot import.
	Printf("%s", x)

	// Match.
	myfmt.Fprint(os.Stderr, myfmt.Errorf("%s", x+"foo"))

	// No match: though this literally matches the template,
	// fmt doesn't resolve to a package here.
	var fmt struct{ Errorf func(string, string) }
	fmt.Errorf("%s", x)

	// Recursive matching:

	// Match: both matches are well-typed, so both succeed.
	myfmt.Errorf("%s", myfmt.Errorf("%s", x+"foo").Error())

	// Outer match succeeds, inner doesn't: 3 has wrong type.
	myfmt.Errorf("%s", myfmt.Errorf("%s", 3).Error())

	// Inner match succeeds, outer doesn't: the inner replacement
	// has the wrong type (error not string).
	myfmt.Errorf("%s", myfmt.Errorf("%s", x+"foo"))
}