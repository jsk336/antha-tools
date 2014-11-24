// antha-tools/antha/ssa/ssautil/switch_test.go: Part of the Antha language
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


package ssautil_test

import (
	"github.com/antha-lang/antha/parser"
	"strings"
	"testing"

	"github.com/antha-lang/antha-tools/antha/loader"
	"github.com/antha-lang/antha-tools/antha/ssa"
	"github.com/antha-lang/antha-tools/antha/ssa/ssautil"
)

func TestSwitches(t *testing.T) {
	conf := loader.Config{ParserMode: parser.ParseComments}
	f, err := conf.ParseFile("testdata/switches.go", nil)
	if err != nil {
		t.Error(err)
		return
	}

	conf.CreateFromFiles("main", f)
	iprog, err := conf.Load()
	if err != nil {
		t.Error(err)
		return
	}

	prog := ssa.Create(iprog, 0)
	mainPkg := prog.Package(iprog.Created[0].Pkg)
	mainPkg.Build()

	for _, mem := range mainPkg.Members {
		if fn, ok := mem.(*ssa.Function); ok {
			if fn.Synthetic != "" {
				continue // e.g. init()
			}
			// Each (multi-line) "switch" comment within
			// this function must match the printed form
			// of a ConstSwitch.
			var wantSwitches []string
			for _, c := range f.Comments {
				if fn.Syntax().Pos() <= c.Pos() && c.Pos() < fn.Syntax().End() {
					text := strings.TrimSpace(c.Text())
					if strings.HasPrefix(text, "switch ") {
						wantSwitches = append(wantSwitches, text)
					}
				}
			}

			switches := ssautil.Switches(fn)
			if len(switches) != len(wantSwitches) {
				t.Errorf("in %s, found %d switches, want %d", fn, len(switches), len(wantSwitches))
			}
			for i, sw := range switches {
				got := sw.String()
				if i >= len(wantSwitches) {
					continue
				}
				want := wantSwitches[i]
				if got != want {
					t.Errorf("in %s, found switch %d: got <<%s>>, want <<%s>>", fn, i, got, want)
				}
			}
		}
	}
}