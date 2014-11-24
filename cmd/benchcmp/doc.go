// antha-tools/cmd/benchcmp/doc.go: Part of the Antha language
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

The benchcmp command displays performance changes between benchmarks.

Benchcmp parses the output of two 'go test' benchmark runs,
correlates the results per benchmark, and displays the deltas.

To measure the performance impact of a change, use 'go test'
to run benchmarks before and after the change:

	go test -run=NONE -bench=. ./... > old.txt
	# make changes
	go test -run=NONE -bench=. ./... > new.txt

Then feed the benchmark results to benchcmp:

	benchcmp old.txt new.txt

Benchcmp will summarize and display the performance changes,
in a format like this:

	$ benchcmp old.txt new.txt
	benchmark           old ns/op     new ns/op     delta
	BenchmarkConcat     523           68.6          -86.88%

	benchmark           old allocs     new allocs     delta
	BenchmarkConcat     3              1              -66.67%

	benchmark           old bytes     new bytes     delta
	BenchmarkConcat     80            48            -40.00%

*/
package main