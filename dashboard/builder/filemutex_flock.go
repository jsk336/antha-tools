// antha-tools/dashboard/builder/filemutex_flock.go: Part of the Antha language
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


// +build darwin dragonfly freebsd linux netbsd openbsd

package main

import (
	"sync"
	"syscall"
)

// FileMutex is similar to sync.RWMutex, but also synchronizes across processes.
// This implementation is based on flock syscall.
type FileMutex struct {
	mu sync.RWMutex
	fd int
}

func MakeFileMutex(filename string) *FileMutex {
	if filename == "" {
		return &FileMutex{fd: -1}
	}
	fd, err := syscall.Open(filename, syscall.O_CREAT|syscall.O_RDONLY, mkdirPerm)
	if err != nil {
		panic(err)
	}
	return &FileMutex{fd: fd}
}

func (m *FileMutex) Lock() {
	m.mu.Lock()
	if m.fd != -1 {
		if err := syscall.Flock(m.fd, syscall.LOCK_EX); err != nil {
			panic(err)
		}
	}
}

func (m *FileMutex) Unlock() {
	if m.fd != -1 {
		if err := syscall.Flock(m.fd, syscall.LOCK_UN); err != nil {
			panic(err)
		}
	}
	m.mu.Unlock()
}

func (m *FileMutex) RLock() {
	m.mu.RLock()
	if m.fd != -1 {
		if err := syscall.Flock(m.fd, syscall.LOCK_SH); err != nil {
			panic(err)
		}
	}
}

func (m *FileMutex) RUnlock() {
	if m.fd != -1 {
		if err := syscall.Flock(m.fd, syscall.LOCK_UN); err != nil {
			panic(err)
		}
	}
	m.mu.RUnlock()
}