// antha-tools/antha/ssa/interp/external_plan9.go: Part of the Antha language
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


package interp

import "syscall"

func ext۰syscall۰Close(fr *frame, args []value) value {
	panic("syscall.Close not yet implemented")
}
func ext۰syscall۰Fstat(fr *frame, args []value) value {
	panic("syscall.Fstat not yet implemented")
}
func ext۰syscall۰Kill(fr *frame, args []value) value {
	panic("syscall.Kill not yet implemented")
}
func ext۰syscall۰Lstat(fr *frame, args []value) value {
	panic("syscall.Lstat not yet implemented")
}
func ext۰syscall۰Open(fr *frame, args []value) value {
	panic("syscall.Open not yet implemented")
}
func ext۰syscall۰ParseDirent(fr *frame, args []value) value {
	panic("syscall.ParseDirent not yet implemented")
}
func ext۰syscall۰Read(fr *frame, args []value) value {
	panic("syscall.Read not yet implemented")
}
func ext۰syscall۰ReadDirent(fr *frame, args []value) value {
	panic("syscall.ReadDirent not yet implemented")
}
func ext۰syscall۰Stat(fr *frame, args []value) value {
	panic("syscall.Stat not yet implemented")
}
func ext۰syscall۰Write(fr *frame, args []value) value {
	// func Write(fd int, p []byte) (n int, err error)
	n, err := write(args[0].(int), valueToBytes(args[1]))
	return tuple{n, wrapError(err)}
}
func ext۰syscall۰RawSyscall(fr *frame, args []value) value {
	return tuple{^uintptr(0), uintptr(0), uintptr(0)}
}

func syswrite(fd int, b []byte) (int, error) {
	return syscall.Write(fd, b)
}