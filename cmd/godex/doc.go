// antha-tools/cmd/godex/doc.go: Part of the Antha language
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


// The godex command prints (dumps) exported information of packages
// or selected package objects.
//
// In contrast to godoc, godex extracts this information from compiled
// object files. Hence the exported data is truly what a compiler will
// see, at the cost of missing commentary.
//
// Usage: godex [flags] {path[.name]}
//
// Each argument must be a (possibly partial) package path, optionally
// followed by a dot and the name of a package object:
//
//	godex math
//	godex math.Sin
//	godex math.Sin fmt.Printf
//	godex antha/types
//
// godex automatically tries all possible package path prefixes if only a
// partial package path is given. For instance, for the path "github.com/antha-lang/antha/types",
// godex prepends "github.com/antha-lang/antha-tools".
//
// The prefixes are computed by searching the directories specified by
// the GOROOT and GOPATH environment variables (and by excluding the
// build OS- and architecture-specific directory names from the path).
// The search order is depth-first and alphabetic; for a partial path
// "foo", a package "a/foo" is found before "b/foo".
//
// Absolute and relative paths may be provided, which disable automatic
// prefix generation:
//
//	godex $GOROOT/pkg/darwin_amd64/sort
//	godex ./sort
//
// All but the last path element may contain dots; a dot in the last path
// element separates the package path from the package object name. If the
// last path element contains a dot, terminate the argument with another
// dot (indicating an empty object name). For instance, the path for a
// package foo.bar would be specified as in:
//
//	godex foo.bar.
//
// The flags are:
//
//	-s=""
//		only consider packages from src, where src is one of the supported compilers
//	-v=false
//		verbose mode
//
// The following sources (-s arguments) are supported:
//
//	gc
//		gc-generated object files
//	gccgo
//		gccgo-generated object files
//	gccgo-new
//		gccgo-generated object files using a condensed format (experimental)
//	source
//		(uncompiled) source code (not yet implemented)
//
// If no -s argument is provided, godex will try to find a matching source.
//
package main

// BUG(gri): support for -s=source is not yet implemented
// BUG(gri): gccgo-importing appears to have occasional problems stalling godex; try -s=gc as work-around