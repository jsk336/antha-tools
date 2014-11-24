// antha-tools/cmd/gotype/doc.go: Part of the Antha language
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


/*
The gotype command does syntactic and semantic analysis of Go files
and packages like the front-end of a Go compiler. Errors are reported
if the analysis fails; otherwise gotype is quiet (unless -v is set).

Without a list of paths, gotype reads from standard input, which
must provide a single Go source file defining a complete package.

If a single path is specified that is a directory, gotype checks
the Go files in that directory; they must all belong to the same
package.

Otherwise, each path must be the filename of Go file belonging to
the same package.

Usage:
	gotype [flags] [path...]

The flags are:
	-a
		use all (incl. _test.go) files when processing a directory
	-e
		report all errors (not just the first 10)
	-v
		verbose mode
	-gccgo
		use gccimporter instead of gcimporter

Debugging flags:
	-seq
		parse sequentially, rather than in parallel
	-ast
		print AST (forces -seq)
	-trace
		print parse trace (forces -seq)
	-comments
		parse comments (ignored unless -ast or -trace is provided)

Examples:

To check the files a.go, b.go, and c.go:

	gotype a.go b.go c.go

To check an entire package in the directory dir and print the processed files:

	gotype -v dir

To check an entire package including tests in the local directory:

	gotype -a .

To verify the output of a pipe:

	echo "package foo" | gotype

*/
package main