// antha-tools/cmd/eg/eg.go: Part of the Antha language
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

// The eg command performs example-based refactoring.
package main

import (
	"flag"
	"fmt"
	"github.com/antha-lang/antha/parser"
	"github.com/antha-lang/antha/printer"
	"github.com/antha-lang/antha/token"
	"os"
	"path/filepath"

	"github.com/antha-lang/antha-tools/antha/loader"
	"github.com/antha-lang/antha-tools/refactor/eg"
)

var (
	helpFlag       = flag.Bool("help", false, "show detailed help message")
	templateFlag   = flag.String("t", "", "template.go file specifying the refactoring")
	transitiveFlag = flag.Bool("transitive", false, "apply refactoring to all dependencies too")
	writeFlag      = flag.Bool("w", false, "rewrite input files in place (by default, the results are printed to standard output)")
	verboseFlag    = flag.Bool("v", false, "show verbose matcher diagnostics")
)

const usage = `eg: an example-based refactoring tool.

Usage: eg -t template.go [-w] [-transitive] <args>...
-t template.go	specifies the template file (use -help to see explanation)
-w          	causes files to be re-written in place.
-transitive 	causes all dependencies to be refactored too.
` + loader.FromArgsUsage

func main() {
	if err := doMain(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s.\n", filepath.Base(os.Args[0]), err)
		os.Exit(1)
	}
}

func doMain() error {
	flag.Parse()
	args := flag.Args()

	if *helpFlag {
		fmt.Fprint(os.Stderr, eg.Help)
		os.Exit(2)
	}

	if *templateFlag == "" {
		return fmt.Errorf("no -t template.go file specified")
	}

	conf := loader.Config{
		Fset:          token.NewFileSet(),
		ParserMode:    parser.ParseComments,
		SourceImports: true,
	}

	// The first Created package is the template.
	if err := conf.CreateFromFilenames("template", *templateFlag); err != nil {
		return err //  e.g. "foo.go:1: syntax error"
	}

	if len(args) == 0 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}

	if _, err := conf.FromArgs(args, true); err != nil {
		return err
	}

	// Load, parse and type-check the whole program.
	iprog, err := conf.Load()
	if err != nil {
		return err
	}

	// Analyze the template.
	template := iprog.Created[0]
	xform, err := eg.NewTransformer(iprog.Fset, template, *verboseFlag)
	if err != nil {
		return err
	}

	// Apply it to the input packages.
	var pkgs []*loader.PackageInfo
	if *transitiveFlag {
		for _, info := range iprog.AllPackages {
			pkgs = append(pkgs, info)
		}
	} else {
		pkgs = iprog.InitialPackages()
	}
	var hadErrors bool
	for _, pkg := range pkgs {
		if pkg == template {
			continue
		}
		for _, file := range pkg.Files {
			n := xform.Transform(&pkg.Info, pkg.Pkg, file)
			if n == 0 {
				continue
			}
			filename := iprog.Fset.File(file.Pos()).Name()
			fmt.Fprintf(os.Stderr, "=== %s (%d matches):\n", filename, n)
			if *writeFlag {
				if err := eg.WriteAST(iprog.Fset, filename, file); err != nil {
					fmt.Fprintf(os.Stderr, "Error: %s\n", err)
					hadErrors = true
				}
			} else {
				printer.Fprint(os.Stdout, iprog.Fset, file)
			}
		}
	}
	if hadErrors {
		os.Exit(1)
	}
	return nil
}