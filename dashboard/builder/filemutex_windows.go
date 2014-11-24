// antha-tools/dashboard/builder/filemutex_windows.go: Part of the Antha language
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


package main

import (
	"sync"
	"syscall"
	"unsafe"
)

var (
	modkernel32      = syscall.NewLazyDLL("kernel32.dll")
	procLockFileEx   = modkernel32.NewProc("LockFileEx")
	procUnlockFileEx = modkernel32.NewProc("UnlockFileEx")
)

const (
	INVALID_FILE_HANDLE     = ^syscall.Handle(0)
	LOCKFILE_EXCLUSIVE_LOCK = 2
)

func lockFileEx(h syscall.Handle, flags, reserved, locklow, lockhigh uint32, ol *syscall.Overlapped) (err error) {
	r1, _, e1 := syscall.Syscall6(procLockFileEx.Addr(), 6, uintptr(h), uintptr(flags), uintptr(reserved), uintptr(locklow), uintptr(lockhigh), uintptr(unsafe.Pointer(ol)))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func unlockFileEx(h syscall.Handle, reserved, locklow, lockhigh uint32, ol *syscall.Overlapped) (err error) {
	r1, _, e1 := syscall.Syscall6(procUnlockFileEx.Addr(), 5, uintptr(h), uintptr(reserved), uintptr(locklow), uintptr(lockhigh), uintptr(unsafe.Pointer(ol)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

// FileMutex is similar to sync.RWMutex, but also synchronizes across processes.
// This implementation is based on flock syscall.
type FileMutex struct {
	mu sync.RWMutex
	fd syscall.Handle
}

func MakeFileMutex(filename string) *FileMutex {
	if filename == "" {
		return &FileMutex{fd: INVALID_FILE_HANDLE}
	}
	fd, err := syscall.CreateFile(&(syscall.StringToUTF16(filename)[0]), syscall.GENERIC_READ|syscall.GENERIC_WRITE,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE, nil, syscall.OPEN_ALWAYS, syscall.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		panic(err)
	}
	return &FileMutex{fd: fd}
}

func (m *FileMutex) Lock() {
	m.mu.Lock()
	if m.fd != INVALID_FILE_HANDLE {
		var ol syscall.Overlapped
		if err := lockFileEx(m.fd, LOCKFILE_EXCLUSIVE_LOCK, 0, 1, 0, &ol); err != nil {
			panic(err)
		}
	}
}

func (m *FileMutex) Unlock() {
	if m.fd != INVALID_FILE_HANDLE {
		var ol syscall.Overlapped
		if err := unlockFileEx(m.fd, 0, 1, 0, &ol); err != nil {
			panic(err)
		}
	}
	m.mu.Unlock()
}

func (m *FileMutex) RLock() {
	m.mu.RLock()
	if m.fd != INVALID_FILE_HANDLE {
		var ol syscall.Overlapped
		if err := lockFileEx(m.fd, 0, 0, 1, 0, &ol); err != nil {
			panic(err)
		}
	}
}

func (m *FileMutex) RUnlock() {
	if m.fd != INVALID_FILE_HANDLE {
		var ol syscall.Overlapped
		if err := unlockFileEx(m.fd, 0, 1, 0, &ol); err != nil {
			panic(err)
		}
	}
	m.mu.RUnlock()
}