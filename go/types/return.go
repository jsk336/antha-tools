// antha-tools/go/types/return.go: Part of the Antha language
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


// This file implements isTerminating.

package types

import (
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/token"
)

// isTerminating reports if s is a terminating statement.
// If s is labeled, label is the label name; otherwise s
// is "".
func (check *checker) isTerminating(s ast.Stmt, label string) bool {
	switch s := s.(type) {
	default:
		unreachable()

	case *ast.BadStmt, *ast.DeclStmt, *ast.EmptyStmt, *ast.SendStmt,
		*ast.IncDecStmt, *ast.AssignStmt, *ast.GoStmt, *ast.DeferStmt,
		*ast.RangeStmt:
		// no chance

	case *ast.LabeledStmt:
		return check.isTerminating(s.Stmt, s.Label.Name)

	case *ast.ExprStmt:
		// the predeclared (possibly parenthesized) panic() function is terminating
		if call, _ := unparen(s.X).(*ast.CallExpr); call != nil {
			if id, _ := call.Fun.(*ast.Ident); id != nil {
				if obj := check.scope.LookupParent(id.Name); obj != nil {
					if b, _ := obj.(*Builtin); b != nil && b.id == _Panic {
						return true
					}
				}
			}
		}

	case *ast.ReturnStmt:
		return true

	case *ast.BranchStmt:
		if s.Tok == token.GOTO || s.Tok == token.FALLTHROUGH {
			return true
		}

	case *ast.BlockStmt:
		return check.isTerminatingList(s.List, "")

	case *ast.IfStmt:
		if s.Else != nil &&
			check.isTerminating(s.Body, "") &&
			check.isTerminating(s.Else, "") {
			return true
		}

	case *ast.SwitchStmt:
		return check.isTerminatingSwitch(s.Body, label)

	case *ast.TypeSwitchStmt:
		return check.isTerminatingSwitch(s.Body, label)

	case *ast.SelectStmt:
		for _, s := range s.Body.List {
			cc := s.(*ast.CommClause)
			if !check.isTerminatingList(cc.Body, "") || hasBreakList(cc.Body, label, true) {
				return false
			}

		}
		return true

	case *ast.ForStmt:
		if s.Cond == nil && !hasBreak(s.Body, label, true) {
			return true
		}
	}

	return false
}

func (check *checker) isTerminatingList(list []ast.Stmt, label string) bool {
	n := len(list)
	return n > 0 && check.isTerminating(list[n-1], label)
}

func (check *checker) isTerminatingSwitch(body *ast.BlockStmt, label string) bool {
	hasDefault := false
	for _, s := range body.List {
		cc := s.(*ast.CaseClause)
		if cc.List == nil {
			hasDefault = true
		}
		if !check.isTerminatingList(cc.Body, "") || hasBreakList(cc.Body, label, true) {
			return false
		}
	}
	return hasDefault
}

// TODO(gri) For nested breakable statements, the current implementation of hasBreak
//	     will traverse the same subtree repeatedly, once for each label. Replace
//           with a single-pass label/break matching phase.

// hasBreak reports if s is or contains a break statement
// referring to the label-ed statement or implicit-ly the
// closest outer breakable statement.
func hasBreak(s ast.Stmt, label string, implicit bool) bool {
	switch s := s.(type) {
	default:
		unreachable()

	case *ast.BadStmt, *ast.DeclStmt, *ast.EmptyStmt, *ast.ExprStmt,
		*ast.SendStmt, *ast.IncDecStmt, *ast.AssignStmt, *ast.GoStmt,
		*ast.DeferStmt, *ast.ReturnStmt:
		// no chance

	case *ast.LabeledStmt:
		return hasBreak(s.Stmt, label, implicit)

	case *ast.BranchStmt:
		if s.Tok == token.BREAK {
			if s.Label == nil {
				return implicit
			}
			if s.Label.Name == label {
				return true
			}
		}

	case *ast.BlockStmt:
		return hasBreakList(s.List, label, implicit)

	case *ast.IfStmt:
		if hasBreak(s.Body, label, implicit) ||
			s.Else != nil && hasBreak(s.Else, label, implicit) {
			return true
		}

	case *ast.CaseClause:
		return hasBreakList(s.Body, label, implicit)

	case *ast.SwitchStmt:
		if label != "" && hasBreak(s.Body, label, false) {
			return true
		}

	case *ast.TypeSwitchStmt:
		if label != "" && hasBreak(s.Body, label, false) {
			return true
		}

	case *ast.CommClause:
		return hasBreakList(s.Body, label, implicit)

	case *ast.SelectStmt:
		if label != "" && hasBreak(s.Body, label, false) {
			return true
		}

	case *ast.ForStmt:
		if label != "" && hasBreak(s.Body, label, false) {
			return true
		}

	case *ast.RangeStmt:
		if label != "" && hasBreak(s.Body, label, false) {
			return true
		}
	}

	return false
}

func hasBreakList(list []ast.Stmt, label string, implicit bool) bool {
	for _, s := range list {
		if hasBreak(s, label, implicit) {
			return true
		}
	}
	return false
}