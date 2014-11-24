// antha-tools/go/pointer/print.go: Part of the Antha language
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


package pointer

import "fmt"

func (c *addrConstraint) String() string {
	return fmt.Sprintf("addr n%d <- {&n%d}", c.dst, c.src)
}

func (c *copyConstraint) String() string {
	return fmt.Sprintf("copy n%d <- n%d", c.dst, c.src)
}

func (c *loadConstraint) String() string {
	return fmt.Sprintf("load n%d <- n%d[%d]", c.dst, c.src, c.offset)
}

func (c *storeConstraint) String() string {
	return fmt.Sprintf("store n%d[%d] <- n%d", c.dst, c.offset, c.src)
}

func (c *offsetAddrConstraint) String() string {
	return fmt.Sprintf("offsetAddr n%d <- n%d.#%d", c.dst, c.src, c.offset)
}

func (c *typeFilterConstraint) String() string {
	return fmt.Sprintf("typeFilter n%d <- n%d.(%s)", c.dst, c.src, c.typ)
}

func (c *untagConstraint) String() string {
	return fmt.Sprintf("untag n%d <- n%d.(%s)", c.dst, c.src, c.typ)
}

func (c *invokeConstraint) String() string {
	return fmt.Sprintf("invoke n%d.%s(n%d ...)", c.iface, c.method.Name(), c.params+1)
}

func (n nodeid) String() string {
	return fmt.Sprintf("n%d", n)
}