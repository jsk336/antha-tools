// antha-tools/anthadoc/static/bake.go: Part of the Antha language
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


// +build ignore

// Command bake takes a list of file names and writes a Go source file to
// standard output that declares a map of string constants containing the input files.
//
// For example, the command
// 	bake foo.html bar.txt
// produces a source file in package main that declares the variable bakedFiles
// that is a map with keys "foo.html" and "bar.txt" that contain the contents
// of foo.html and bar.txt.
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"unicode/utf8"
)

func main() {
	if err := bake(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func bake(files []string) error {
	w := bufio.NewWriter(os.Stdout)
	fmt.Fprintf(w, "%v\n\npackage static\n\n", warning)
	fmt.Fprintf(w, "var Files = map[string]string{\n")
	for _, fn := range files {
		b, err := ioutil.ReadFile(fn)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "\t%q: ", fn)
		if utf8.Valid(b) {
			fmt.Fprintf(w, "`%s`", sanitize(b))
		} else {
			fmt.Fprintf(w, "%q", b)
		}
		fmt.Fprintln(w, ",\n")
	}
	fmt.Fprintln(w, "}")
	return w.Flush()
}

// sanitize prepares a valid UTF-8 string as a raw string constant.
func sanitize(b []byte) []byte {
	// Replace ` with `+"`"+`
	b = bytes.Replace(b, []byte("`"), []byte("`+\"`\"+`"), -1)

	// Replace BOM with `+"\xEF\xBB\xBF"+`
	// (A BOM is valid UTF-8 but not permitted in Go source files.
	// I wouldn't bother handling this, but for some insane reason
	// jquery.js has a BOM somewhere in the middle.)
	return bytes.Replace(b, []byte("\xEF\xBB\xBF"), []byte("`+\"\\xEF\\xBB\\xBF\"+`"), -1)
}

const warning = "// DO NOT EDIT ** This file was generated with the bake tool ** DO NOT EDIT //"