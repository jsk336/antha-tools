// antha-tools/cmd/vet/buildtag.go: Part of the Antha language
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
	"bytes"
	"fmt"
	"os"
	"strings"
	"unicode"
)

var (
	nl         = []byte("\n")
	slashSlash = []byte("//")
	plusBuild  = []byte("+build")
)

// checkBuildTag checks that build tags are in the correct location and well-formed.
func checkBuildTag(name string, data []byte) {
	if !vet("buildtags") {
		return
	}
	lines := bytes.SplitAfter(data, nl)

	// Determine cutpoint where +build comments are no longer valid.
	// They are valid in leading // comments in the file followed by
	// a blank line.
	var cutoff int
	for i, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			cutoff = i
			continue
		}
		if bytes.HasPrefix(line, slashSlash) {
			continue
		}
		break
	}

	for i, line := range lines {
		line = bytes.TrimSpace(line)
		if !bytes.HasPrefix(line, slashSlash) {
			continue
		}
		text := bytes.TrimSpace(line[2:])
		if bytes.HasPrefix(text, plusBuild) {
			fields := bytes.Fields(text)
			if !bytes.Equal(fields[0], plusBuild) {
				// Comment is something like +buildasdf not +build.
				fmt.Fprintf(os.Stderr, "%s:%d: possible malformed +build comment\n", name, i+1)
				continue
			}
			if i >= cutoff {
				fmt.Fprintf(os.Stderr, "%s:%d: +build comment must appear before package clause and be followed by a blank line\n", name, i+1)
				setExit(1)
				continue
			}
			// Check arguments.
		Args:
			for _, arg := range fields[1:] {
				for _, elem := range strings.Split(string(arg), ",") {
					if strings.HasPrefix(elem, "!!") {
						fmt.Fprintf(os.Stderr, "%s:%d: invalid double negative in build constraint: %s\n", name, i+1, arg)
						setExit(1)
						break Args
					}
					if strings.HasPrefix(elem, "!") {
						elem = elem[1:]
					}
					for _, c := range elem {
						if !unicode.IsLetter(c) && !unicode.IsDigit(c) && c != '_' && c != '.' {
							fmt.Fprintf(os.Stderr, "%s:%d: invalid non-alphanumeric build constraint: %s\n", name, i+1, arg)
							setExit(1)
							break Args
						}
					}
				}
			}
			continue
		}
		// Comment with +build but not at beginning.
		if bytes.Contains(line, plusBuild) && i < cutoff {
			fmt.Fprintf(os.Stderr, "%s:%d: possible malformed +build comment\n", name, i+1)
			continue
		}
	}
}