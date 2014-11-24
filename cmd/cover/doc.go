// antha-tools/cmd/cover/doc.go: Part of the Antha language
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
Cover is a program for analyzing the coverage profiles generated by
'go test -coverprofile=cover.out'.

Cover is also used by 'go test -cover' to rewrite the source code with
annotations to track which parts of each function are executed.
It operates on one Go source file at a time, computing approximate
basic block information by studying the source. It is thus more portable
than binary-rewriting coverage tools, but also a little less capable.
For instance, it does not probe inside && and || expressions, and can
be mildly confused by single statements with multiple function literals.

For usage information, please see:
	go help testflag
	go tool cover -help
*/
package main