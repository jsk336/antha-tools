// antha-tools/go/ssa/example_test.go: Part of the Antha language
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


package ssa_test

import (
	"fmt"
	"os"

	"github.com/antha-lang/antha-tools/antha/loader"
	"github.com/antha-lang/antha-tools/antha/ssa"
)

// This program demonstrates how to run the SSA builder on a "Hello,
// World!" program and shows the printed representation of packages,
// functions and instructions.
//
// Within the function listing, the name of each BasicBlock such as
// ".0.entry" is printed left-aligned, followed by the block's
// Instructions.
//
// For each instruction that defines an SSA virtual register
// (i.e. implements Value), the type of that value is shown in the
// right column.
//
// Build and run the ssadump.go program if you want a standalone tool
// with similar functionality. It is located at
// antha-tools/cmd/ssadump.
//
func Example() {
	const hello = `
package main

import "fmt"

const message = "Hello, World!"

func main() {
	fmt.Println(message)
}
`
	var conf loader.Config

	// Parse the input file.
	file, err := conf.ParseFile("hello.go", hello)
	if err != nil {
		fmt.Print(err) // parse error
		return
	}

	// Create single-file main package.
	conf.CreateFromFiles("main", file)

	// Load the main package and its dependencies.
	iprog, err := conf.Load()
	if err != nil {
		fmt.Print(err) // type error in some package
		return
	}

	// Create SSA-form program representation.
	prog := ssa.Create(iprog, ssa.SanityCheckFunctions)
	mainPkg := prog.Package(iprog.Created[0].Pkg)

	// Print out the package.
	mainPkg.WriteTo(os.Stdout)

	// Build SSA code for bodies of functions in mainPkg.
	mainPkg.Build()

	// Print out the package-level functions.
	mainPkg.Func("init").WriteTo(os.Stdout)
	mainPkg.Func("main").WriteTo(os.Stdout)

	// Output:
	//
	// package main:
	//   func  init       func()
	//   var   init$guard bool
	//   func  main       func()
	//   const message    message = "Hello, World!":untyped string
	//
	// # Name: main.init
	// # Package: main
	// # Synthetic: package initializer
	// func init():
	// 0:                                                                entry P:0 S:2
	// 	t0 = *init$guard                                                   bool
	// 	if t0 goto 2 else 1
	// 1:                                                           init.start P:1 S:1
	// 	*init$guard = true:bool
	// 	t1 = fmt.init()                                                      ()
	// 	jump 2
	// 2:                                                            init.done P:2 S:0
	// 	return
	//
	// # Name: main.main
	// # Package: main
	// # Location: hello.go:8:6
	// func main():
	// 0:                                                                entry P:0 S:0
	// 	t0 = new [1]interface{} (varargs)                       *[1]interface{}
	// 	t1 = &t0[0:int]                                            *interface{}
	// 	t2 = make interface{} <- string ("Hello, World!":string)    interface{}
	// 	*t1 = t2
	// 	t3 = slice t0[:]                                          []interface{}
	// 	t4 = fmt.Println(t3...)                              (n int, err error)
	// 	return
}