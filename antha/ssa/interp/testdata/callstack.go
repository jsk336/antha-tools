// antha-tools/antha/ssa/interp/testdata/callstack.go: Part of the Antha language
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

package main

import (
	"fmt"
	"path"
	"runtime"
	"strings"
)

var stack string

func f() {
	pc := make([]uintptr, 6)
	pc = pc[:runtime.Callers(1, pc)]
	for _, f := range pc {
		Func := runtime.FuncForPC(f)
		name := Func.Name()
		if strings.Contains(name, "$") || strings.Contains(name, ".func") {
			name = "func" // anon funcs vary across toolchains
		}
		file, line := Func.FileLine(0)
		stack += fmt.Sprintf("%s at %s:%d\n", name, path.Base(file), line)
	}
}

func g() { f() }
func h() { g() }
func i() { func() { h() }() }

// Hack: the 'func' and the call to Caller are on the same line,
// to paper over differences between toolchains.
// (The interpreter's location info isn't yet complete.)
func runtimeCaller0() (uintptr, string, int, bool) { return runtime.Caller(0) }

func main() {
	i()
	if stack != `main.f at callstack.go:12
main.g at callstack.go:26
main.h at callstack.go:27
func at callstack.go:28
main.i at callstack.go:28
main.main at callstack.go:35
` {
		panic("unexpected stack: " + stack)
	}

	pc, file, line, _ := runtimeCaller0()
	got := fmt.Sprintf("%s @ %s:%d", runtime.FuncForPC(pc).Name(), path.Base(file), line)
	if got != "main.runtimeCaller0 @ callstack.go:33" {
		panic("runtime.Caller: " + got)
	}
}