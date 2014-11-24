// antha-tools/antha/gccgoimporter/importer.go: Part of the Antha language
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


// Package gccgoimporter implements Import for gccgo-generated object files.
package gccgoimporter

import (
	"debug/elf"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/antha-lang/antha-tools/antha/types"
)

// Locate the file from which to read export data.
// This is intended to replicate the logic in gofrontend.
func findExportFile(searchpaths []string, pkgpath string) (string, error) {
	for _, spath := range searchpaths {
		pkgfullpath := filepath.Join(spath, pkgpath)
		pkgdir, name := filepath.Split(pkgfullpath)

		for _, filepath := range [...]string{
			pkgfullpath,
			pkgfullpath + ".gox",
			pkgdir + "lib" + name + ".so",
			pkgdir + "lib" + name + ".a",
			pkgfullpath + ".o",
		} {
			fi, err := os.Stat(filepath)
			if err == nil && !fi.IsDir() {
				return filepath, nil
			}
		}
	}

	return "", fmt.Errorf("%s: could not find export data (tried %s)", pkgpath, strings.Join(searchpaths, ":"))
}

// Opens the export data file at the given path. If this is an ELF file,
// searches for and opens the .go_export section.
// This is intended to replicate the logic in gofrontend, although it doesn't handle archive files yet.
func openExportFile(fpath string) (reader io.ReadSeeker, closer io.Closer, err error) {
	f, err := os.Open(fpath)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			f.Close()
		}
	}()
	closer = f

	var magic [4]byte
	_, err = f.ReadAt(magic[:], 0)
	if err != nil {
		return
	}

	if string(magic[:]) == "v1;\n" {
		// Raw export data.
		reader = f
		return
	}

	ef, err := elf.NewFile(f)
	if err != nil {
		return
	}

	sec := ef.Section(".go_export")
	if sec == nil {
		err = fmt.Errorf("%s: .go_export section not found", fpath)
		return
	}

	reader = sec.Open()
	return
}

func GetImporter(searchpaths []string) types.Importer {
	return func(imports map[string]*types.Package, pkgpath string) (pkg *types.Package, err error) {
		if pkgpath == "unsafe" {
			return types.Unsafe, nil
		}

		fpath, err := findExportFile(searchpaths, pkgpath)
		if err != nil {
			return
		}

		reader, closer, err := openExportFile(fpath)
		if err != nil {
			return
		}
		defer closer.Close()

		var p parser
		p.init(fpath, reader, imports)
		pkg = p.parsePackage()
		return
	}
}