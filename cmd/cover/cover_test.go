// antha-tools/cmd/cover/cover_test.go: Part of the Antha language
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
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

const (
	// Data directory, also the package directory for the test.
	testdata = "testdata"

	// Binaries we compile.
	testcover = "./testcover.exe"
)

var (
	// Files we use.
	testMain    = filepath.Join(testdata, "main.go")
	testTest    = filepath.Join(testdata, "test.go")
	coverInput  = filepath.Join(testdata, "test_line.go")
	coverOutput = filepath.Join(testdata, "test_cover.go")
)

var debug = false // Keeps the rewritten files around if set.

// Run this shell script, but do it in Go so it can be run by "go test".
//
//	replace the word LINE with the line number < testdata/test.go > testdata/test_line.go
// 	go build -o ./testcover
// 	./testcover -mode=count -var=CoverTest -o ./testdata/test_cover.go testdata/test_line.go
//	go run ./testdata/main.go ./testdata/test.go
//
func TestCover(t *testing.T) {
	// Read in the test file (testTest) and write it, with LINEs specified, to coverInput.
	file, err := ioutil.ReadFile(testTest)
	if err != nil {
		t.Fatal(err)
	}
	lines := bytes.Split(file, []byte("\n"))
	for i, line := range lines {
		lines[i] = bytes.Replace(line, []byte("LINE"), []byte(fmt.Sprint(i+1)), -1)
	}
	err = ioutil.WriteFile(coverInput, bytes.Join(lines, []byte("\n")), 0666)

	// defer removal of test_line.go
	if !debug {
		defer os.Remove(coverInput)
	}

	// antha build -o testcover
	cmd := exec.Command("go", "build", "-o", testcover)
	run(cmd, t)

	// defer removal of testcover
	defer os.Remove(testcover)

	// ./testcover -mode=count -var=coverTest -o ./testdata/test_cover.go testdata/test_line.go
	cmd = exec.Command(testcover, "-mode=count", "-var=coverTest", "-o", coverOutput, coverInput)
	run(cmd, t)

	// defer removal of ./testdata/test_cover.go
	if !debug {
		defer os.Remove(coverOutput)
	}

	// antha run ./testdata/main.go ./testdata/test.go
	cmd = exec.Command("go", "run", testMain, coverOutput)
	run(cmd, t)
}

func run(c *exec.Cmd, t *testing.T) {
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		t.Fatal(err)
	}
}