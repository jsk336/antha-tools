// antha-tools/cmd/vet/testdata/composite.go: Part of the Antha language
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


// This file contains tests for the untagged struct literal checker.

// This file contains the test for untagged struct literals.

package testdata

import (
	"flag"
	"github.com/antha-lang/antha/scanner"
)

var Okay1 = []string{
	"Name",
	"Usage",
	"DefValue",
}

var Okay2 = map[string]bool{
	"Name":     true,
	"Usage":    true,
	"DefValue": true,
}

var Okay3 = struct {
	X string
	Y string
	Z string
}{
	"Name",
	"Usage",
	"DefValue",
}

type MyStruct struct {
	X string
	Y string
	Z string
}

var Okay4 = MyStruct{
	"Name",
	"Usage",
	"DefValue",
}

// Testing is awkward because we need to reference things from a separate package
// to trigger the warnings.

var BadStructLiteralUsedInTests = flag.Flag{ // ERROR "unkeyed fields"
	"Name",
	"Usage",
	nil, // Value
	"DefValue",
}

// Used to test the check for slices and arrays: If that test is disabled and
// vet is run with --compositewhitelist=false, this line triggers an error.
// Clumsy but sufficient.
var scannerErrorListTest = scanner.ErrorList{nil, nil}