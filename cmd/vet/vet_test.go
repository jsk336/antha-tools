// antha-tools/cmd/vet/vet_test.go: Part of the Antha language
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


package main_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

const (
	dataDir = "testdata"
	binary  = "testvet"
)

// Run this shell script, but do it in Go so it can be run by "go test".
// 	go build -o testvet
// 	$(GOROOT)/test/errchk ./testvet -shadow -printfuncs='Warn:1,Warnf:1' testdata/*.go testdata/*.s
// 	rm testvet
//
func TestVet(t *testing.T) {
	// Plan 9 and Windows systems can't be guaranteed to have Perl and so can't run errchk.
	switch runtime.GOOS {
	case "plan9", "windows":
		t.Skip("skipping test; no Perl on %q", runtime.GOOS)
	}

	// antha build
	cmd := exec.Command("go", "build", "-o", binary)
	run(cmd, t)

	// defer removal of vet
	defer os.Remove(binary)

	// errchk ./testvet
	gos, err := filepath.Glob(filepath.Join(dataDir, "*.go"))
	if err != nil {
		t.Fatal(err)
	}
	asms, err := filepath.Glob(filepath.Join(dataDir, "*.s"))
	if err != nil {
		t.Fatal(err)
	}
	files := append(gos, asms...)
	errchk := filepath.Join(runtime.GOROOT(), "test", "errchk")
	flags := []string{
		"./" + binary,
		"-printfuncs=Warn:1,Warnf:1",
		"-test", // TODO: Delete once -shadow is part of -all.
	}
	cmd = exec.Command(errchk, append(flags, files...)...)
	if !run(cmd, t) {
		t.Fatal("vet command failed")
	}
}

func run(c *exec.Cmd, t *testing.T) bool {
	output, err := c.CombinedOutput()
	os.Stderr.Write(output)
	if err != nil {
		t.Fatal(err)
	}
	// Errchk delights by not returning non-zero status if it finds errors, so we look at the output.
	// It prints "BUG" if there is a failure.
	if !c.ProcessState.Success() {
		return false
	}
	return !bytes.Contains(output, []byte("BUG"))
}