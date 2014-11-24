// antha-tools/godoc/vfs/gatefs/gatefs.go: Part of the Antha language
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


// Package gatefs provides an implementation of the FileSystem
// interface that wraps another FileSystem and limits its concurrency.
package gatefs

import (
	"fmt"
	"os"

	"github.com/antha-lang/antha-tools/anthadoc/vfs"
)

// New returns a new FileSystem that delegates to fs.
// If gateCh is non-nil and buffered, it's used as a gate
// to limit concurrency on calls to fs.
func New(fs vfs.FileSystem, gateCh chan bool) vfs.FileSystem {
	if cap(gateCh) == 0 {
		return fs
	}
	return gatefs{fs, gate(gateCh)}
}

type gate chan bool

func (g gate) enter() { g <- true }
func (g gate) leave() { <-g }

type gatefs struct {
	fs vfs.FileSystem
	gate
}

func (fs gatefs) String() string {
	return fmt.Sprintf("gated(%s, %d)", fs.fs.String(), cap(fs.gate))
}

func (fs gatefs) Open(p string) (vfs.ReadSeekCloser, error) {
	fs.enter()
	defer fs.leave()
	rsc, err := fs.fs.Open(p)
	if err != nil {
		return nil, err
	}
	return gatef{rsc, fs.gate}, nil
}

func (fs gatefs) Lstat(p string) (os.FileInfo, error) {
	fs.enter()
	defer fs.leave()
	return fs.fs.Lstat(p)
}

func (fs gatefs) Stat(p string) (os.FileInfo, error) {
	fs.enter()
	defer fs.leave()
	return fs.fs.Stat(p)
}

func (fs gatefs) ReadDir(p string) ([]os.FileInfo, error) {
	fs.enter()
	defer fs.leave()
	return fs.fs.ReadDir(p)
}

type gatef struct {
	rsc vfs.ReadSeekCloser
	gate
}

func (f gatef) Read(p []byte) (n int, err error) {
	f.enter()
	defer f.leave()
	return f.rsc.Read(p)
}

func (f gatef) Seek(offset int64, whence int) (ret int64, err error) {
	f.enter()
	defer f.leave()
	return f.rsc.Seek(offset, whence)
}

func (f gatef) Close() error {
	f.enter()
	defer f.leave()
	return f.rsc.Close()
}