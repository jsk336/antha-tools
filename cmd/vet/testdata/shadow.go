// antha-tools/cmd/vet/testdata/shadow.go: Part of the Antha language
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


// This file contains tests for the shadowed variable checker.
// Some of these errors are caught by the compiler (shadowed return parameters for example)
// but are nonetheless useful tests.

package testdata

import "os"

func ShadowRead(f *os.File, buf []byte) (err error) {
	var x int
	if f != nil {
		err := 3 // OK - different type.
		_ = err
	}
	if f != nil {
		_, err := f.Read(buf) // ERROR "declaration of err shadows declaration at testdata/shadow.go:13"
		if err != nil {
			return err
		}
		i := 3 // OK
		_ = i
	}
	if f != nil {
		var _, err = f.Read(buf) // ERROR "declaration of err shadows declaration at testdata/shadow.go:13"
		if err != nil {
			return err
		}
	}
	for i := 0; i < 10; i++ {
		i := i // OK: obviously intentional idiomatic redeclaration
		go func() {
			println(i)
		}()
	}
	var shadowTemp interface{}
	switch shadowTemp := shadowTemp.(type) { // OK: obviously intentional idiomatic redeclaration
	case int:
		println("OK")
		_ = shadowTemp
	}
	if shadowTemp := shadowTemp; true { // OK: obviously intentional idiomatic redeclaration
		var f *os.File // OK because f is not mentioned later in the function.
		// The declaration of x is a shadow because x is mentioned below.
		var x int // ERROR "declaration of x shadows declaration at testdata/shadow.go:14"
		_, _, _ = x, f, shadowTemp
	}
	// Use a couple of variables to trigger shadowing errors.
	_, _ = err, x
	return
}