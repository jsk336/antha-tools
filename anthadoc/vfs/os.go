// antha-tools/anthadoc/vfs/os.go: Part of the Antha language
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


package vfs

import (
	"fmt"
	"io/ioutil"
	"os"
	pathpkg "path"
	"path/filepath"
)

// OS returns an implementation of FileSystem reading from the
// tree rooted at root.  Recording a root is convenient everywhere
// but necessary on Windows, because the slash-separated path
// passed to Open has no way to specify a drive letter.  Using a root
// lets code refer to OS(`c:\`), OS(`d:\`) and so on.
func OS(root string) FileSystem {
	return osFS(root)
}

type osFS string

func (root osFS) String() string { return "os(" + string(root) + ")" }

func (root osFS) resolve(path string) string {
	// Clean the path so that it cannot possibly begin with ../.
	// If it did, the result of filepath.Join would be outside the
	// tree rooted at root.  We probably won't ever see a path
	// with .. in it, but be safe anyway.
	path = pathpkg.Clean("/" + path)

	return filepath.Join(string(root), path)
}

func (root osFS) Open(path string) (ReadSeekCloser, error) {
	f, err := os.Open(root.resolve(path))
	if err != nil {
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return nil, fmt.Errorf("Open: %s is a directory", path)
	}
	return f, nil
}

func (root osFS) Lstat(path string) (os.FileInfo, error) {
	return os.Lstat(root.resolve(path))
}

func (root osFS) Stat(path string) (os.FileInfo, error) {
	return os.Stat(root.resolve(path))
}

func (root osFS) ReadDir(path string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(root.resolve(path)) // is sorted
}