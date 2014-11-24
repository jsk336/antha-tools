// antha-tools/antha/types/token_test.go: Part of the Antha language
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


// This file checks invariants of token.Token ordering that we rely on
// since package antha/token doesn't provide any guarantees at the moment.

package types

import (
	"github.com/antha-lang/antha/token"
	"testing"
)

var assignOps = map[token.Token]token.Token{
	token.ADD_ASSIGN:     token.ADD,
	token.SUB_ASSIGN:     token.SUB,
	token.MUL_ASSIGN:     token.MUL,
	token.QUO_ASSIGN:     token.QUO,
	token.REM_ASSIGN:     token.REM,
	token.AND_ASSIGN:     token.AND,
	token.OR_ASSIGN:      token.OR,
	token.XOR_ASSIGN:     token.XOR,
	token.SHL_ASSIGN:     token.SHL,
	token.SHR_ASSIGN:     token.SHR,
	token.AND_NOT_ASSIGN: token.AND_NOT,
}

func TestZeroTok(t *testing.T) {
	// zero value for token.Token must be token.ILLEGAL
	var zero token.Token
	if token.ILLEGAL != zero {
		t.Errorf("%s == %d; want 0", token.ILLEGAL, zero)
	}
}

func TestAssignOp(t *testing.T) {
	// there are fewer than 256 tokens
	for i := 0; i < 256; i++ {
		tok := token.Token(i)
		got := assignOp(tok)
		want := assignOps[tok]
		if got != want {
			t.Errorf("for assignOp(%s): got %s; want %s", tok, got, want)
		}
	}
}