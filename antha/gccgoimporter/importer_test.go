// antha-tools/antha/gccgoimporter/importer_test.go: Part of the Antha language
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
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/antha-lang/antha-tools/antha/types"
)

type importerTest struct {
	pkgpath, name, want, wantval string
}

func runImporterTest(t *testing.T, imp types.Importer, test *importerTest) {
	pkg, err := imp(make(map[string]*types.Package), test.pkgpath)
	if err != nil {
		t.Error(err)
		return
	}

	obj := pkg.Scope().Lookup(test.name)
	if obj == nil {
		t.Errorf("%s: object not found", test.name)
		return
	}

	got := types.ObjectString(pkg, obj)
	if got != test.want {
		t.Errorf("%s: got %q; want %q", test.name, got, test.want)
	}

	if test.wantval != "" {
		gotval := obj.(*types.Const).Val().String()
		if gotval != test.wantval {
			t.Errorf("%s: got val %q; want val %q", test.name, gotval, test.wantval)
		}
	}
}

var importerTests = [...]importerTest{
	{pkgpath: "pointer", name: "Int8Ptr", want: "type Int8Ptr *int8"},
	{pkgpath: "complexnums", name: "NN", want: "const NN untyped complex", wantval: "(-1/1 + -1/1i)"},
	{pkgpath: "complexnums", name: "NP", want: "const NP untyped complex", wantval: "(-1/1 + 1/1i)"},
	{pkgpath: "complexnums", name: "PN", want: "const PN untyped complex", wantval: "(1/1 + -1/1i)"},
	{pkgpath: "complexnums", name: "PP", want: "const PP untyped complex", wantval: "(1/1 + 1/1i)"},
}

func TestGoxImporter(t *testing.T) {
	imp := GetImporter([]string{"testdata"})

	for _, test := range importerTests {
		runImporterTest(t, imp, &test)
	}
}

func TestObjImporter(t *testing.T) {
	// This test relies on gccgo being around, which it most likely will be if we
	// were compiled with gccgo.
	if runtime.Compiler != "gccgo" {
		t.Skip("This test needs gccgo")
		return
	}

	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	imp := GetImporter([]string{tmpdir})

	for _, test := range importerTests {
		gofile := filepath.Join("testdata", test.pkgpath+".go")
		ofile := filepath.Join(tmpdir, test.pkgpath+".o")

		cmd := exec.Command("gccgo", "-c", "-o", ofile, gofile)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Logf("%s", out)
			t.Fatalf("gccgo %s failed: %s", gofile, err)
		}

		runImporterTest(t, imp, &test)

		if err := os.Remove(ofile); err != nil {
			t.Fatal(err)
		}
	}

	if err = os.Remove(tmpdir); err != nil {
		t.Fatal(err)
	}
}