// antha-tools/present/image.go: Part of the Antha language
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


package present

import (
	"fmt"
	"strings"
)

func init() {
	Register("image", parseImage)
}

type Image struct {
	URL    string
	Width  int
	Height int
}

func (i Image) TemplateName() string { return "image" }

func parseImage(ctx *Context, fileName string, lineno int, text string) (Elem, error) {
	args := strings.Fields(text)
	img := Image{URL: args[1]}
	a, err := parseArgs(fileName, lineno, args[2:])
	if err != nil {
		return nil, err
	}
	switch len(a) {
	case 0:
		// no size parameters
	case 2:
		// If a parameter is empty (underscore) or invalid
		// leave the field set to zero. The "image" action
		// template will then omit that img tag attribute and
		// the browser will calculate the value to preserve
		// the aspect ratio.
		if v, ok := a[0].(int); ok {
			img.Height = v
		}
		if v, ok := a[1].(int); ok {
			img.Width = v
		}
	default:
		return nil, fmt.Errorf("incorrect image invocation: %q", text)
	}
	return img, nil
}