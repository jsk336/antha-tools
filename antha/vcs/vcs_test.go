// antha-tools/antha/vcs/vcs_test.go: Part of the Antha language
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


package vcs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// Test that RepoRootForImportPath creates the correct RepoRoot for a given importPath.
// TODO(cmang): Add tests for SVN and BZR.
func TestRepoRootForImportPath(t *testing.T) {
	tests := []struct {
		path string
		want *RepoRoot
	}{
		{
			"code.google.com/p/go",
			&RepoRoot{
				VCS:  vcsHg,
				Repo: "https://code.google.com/p/go",
			},
		},
		{
			"code.google.com/r/go",
			&RepoRoot{
				VCS:  vcsHg,
				Repo: "https://code.google.com/r/go",
			},
		},
		{
			"github.com/golang/groupcache",
			&RepoRoot{
				VCS:  vcsGit,
				Repo: "https://github.com/golang/groupcache",
			},
		},
	}

	for _, test := range tests {
		got, err := RepoRootForImportPath(test.path, false)
		if err != nil {
			t.Errorf("RepoRootForImport(%q): %v", test.path, err)
			continue
		}
		want := test.want
		if got.VCS.Name != want.VCS.Name || got.Repo != want.Repo {
			t.Errorf("RepoRootForImport(%q) = VCS(%s) Repo(%s), want VCS(%s) Repo(%s)", test.path, got.VCS, got.Repo, want.VCS, want.Repo)
		}
	}
}

// Test that FromDir correctly inspects a given directory and returns the right VCS.
func TestFromDir(t *testing.T) {
	type testStruct struct {
		path string
		want *Cmd
	}

	tests := make([]testStruct, len(vcsList))
	tempDir, err := ioutil.TempDir("", "vcstest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	for i, vcs := range vcsList {
		tests[i] = testStruct{
			filepath.Join(tempDir, vcs.Name, "."+vcs.Cmd),
			vcs,
		}
	}

	for _, test := range tests {
		os.MkdirAll(test.path, 0755)
		got, _, _ := FromDir(test.path, tempDir)
		if got.Name != test.want.Name {
			t.Errorf("FromDir(%q, %q) = %s, want %s", test.path, tempDir, got, test.want)
		}
		os.RemoveAll(test.path)
	}
}