// antha-tools/anthadoc/vfs/httpfs/httpfs.go: Part of the Antha language
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


// Package httpfs implements http.FileSystem using a anthadoc vfs.FileSystem.
package httpfs

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/antha-lang/antha-tools/anthadoc/vfs"
)

func New(fs vfs.FileSystem) http.FileSystem {
	return &httpFS{fs}
}

type httpFS struct {
	fs vfs.FileSystem
}

func (h *httpFS) Open(name string) (http.File, error) {
	fi, err := h.fs.Stat(name)
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return &httpDir{h.fs, name, nil}, nil
	}
	f, err := h.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return &httpFile{h.fs, f, name}, nil
}

// httpDir implements http.File for a directory in a FileSystem.
type httpDir struct {
	fs      vfs.FileSystem
	name    string
	pending []os.FileInfo
}

func (h *httpDir) Close() error               { return nil }
func (h *httpDir) Stat() (os.FileInfo, error) { return h.fs.Stat(h.name) }
func (h *httpDir) Read([]byte) (int, error) {
	return 0, fmt.Errorf("cannot Read from directory %s", h.name)
}

func (h *httpDir) Seek(offset int64, whence int) (int64, error) {
	if offset == 0 && whence == 0 {
		h.pending = nil
		return 0, nil
	}
	return 0, fmt.Errorf("unsupported Seek in directory %s", h.name)
}

func (h *httpDir) Readdir(count int) ([]os.FileInfo, error) {
	if h.pending == nil {
		d, err := h.fs.ReadDir(h.name)
		if err != nil {
			return nil, err
		}
		if d == nil {
			d = []os.FileInfo{} // not nil
		}
		h.pending = d
	}

	if len(h.pending) == 0 && count > 0 {
		return nil, io.EOF
	}
	if count <= 0 || count > len(h.pending) {
		count = len(h.pending)
	}
	d := h.pending[:count]
	h.pending = h.pending[count:]
	return d, nil
}

// httpFile implements http.File for a file (not directory) in a FileSystem.
type httpFile struct {
	fs vfs.FileSystem
	vfs.ReadSeekCloser
	name string
}

func (h *httpFile) Stat() (os.FileInfo, error) { return h.fs.Stat(h.name) }
func (h *httpFile) Readdir(int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("cannot Readdir from file %s", h.name)
}