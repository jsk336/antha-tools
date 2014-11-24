// antha-tools/go/gccgoimporter/gccgoinstallation.go: Part of the Antha language
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


package gccgoimporter

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/antha-lang/antha-tools/antha/types"
)

// Information about a specific installation of gccgo.
type GccgoInstallation struct {
	// Version of gcc (e.g. 4.8.0).
	GccVersion string

	// Target triple (e.g. x86_64-unknown-linux-gnu).
	TargetTriple string

	// Built-in library paths used by this installation.
	LibPaths []string
}

// Ask the driver at the given path for information for this GccgoInstallation.
func (inst *GccgoInstallation) InitFromDriver(gccgoPath string) (err error) {
	cmd := exec.Command(gccgoPath, "-###", "-S", "-x", "go", "-")
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return
	}

	err = cmd.Start()
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "Target: "):
			inst.TargetTriple = line[8:]

		case line[0] == ' ':
			args := strings.Fields(line)
			for _, arg := range args[1:] {
				if strings.HasPrefix(arg, "-L") {
					inst.LibPaths = append(inst.LibPaths, arg[2:])
				}
			}
		}
	}

	stdout, err := exec.Command(gccgoPath, "-dumpversion").Output()
	if err != nil {
		return
	}
	inst.GccVersion = strings.TrimSpace(string(stdout))

	return
}

// Return the list of export search paths for this GccgoInstallation.
func (inst *GccgoInstallation) SearchPaths() (paths []string) {
	for _, lpath := range inst.LibPaths {
		spath := filepath.Join(lpath, "go", inst.GccVersion)
		fi, err := os.Stat(spath)
		if err != nil || !fi.IsDir() {
			continue
		}
		paths = append(paths, spath)

		spath = filepath.Join(spath, inst.TargetTriple)
		fi, err = os.Stat(spath)
		if err != nil || !fi.IsDir() {
			continue
		}
		paths = append(paths, spath)
	}

	paths = append(paths, inst.LibPaths...)

	return
}

// Return an importer that searches incpaths followed by the gcc installation's
// built-in search paths and the current directory.
func (inst *GccgoInstallation) GetImporter(incpaths []string) types.Importer {
	return GetImporter(append(append(incpaths, inst.SearchPaths()...), "."))
}