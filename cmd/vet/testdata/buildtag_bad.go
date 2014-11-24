// antha-tools/cmd/vet/testdata/buildtag_bad.go: Part of the Antha language
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

// This file contains misplaced or malformed build constraints.
// The Go tool will skip it, because the constraints are invalid.
// It serves only to test the tag checker during make test.

// Mention +build // ERROR "possible malformed \+build comment"

// +build !!bang // ERROR "invalid double negative in build constraint"
// +build @#$ // ERROR "invalid non-alphanumeric build constraint"

// +build toolate // ERROR "build comment must appear before package clause and be followed by a blank line"
package bad

// This is package 'bad' rather than 'main' so the erroneous build
// tag doesn't end up looking like a package doc for the vet command
// when examined by godoc.