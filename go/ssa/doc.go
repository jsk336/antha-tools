// antha-tools/go/ssa/doc.go: Part of the Antha language
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


// Package ssa defines a representation of the elements of Go programs
// (packages, types, functions, variables and constants) using a
// static single-assignment (SSA) form intermediate representation
// (IR) for the bodies of functions.
//
// THIS INTERFACE IS EXPERIMENTAL AND IS LIKELY TO CHANGE.
//
// For an introduction to SSA form, see
// http://en.wikipedia.org/wiki/Static_single_assignment_form.
// This page provides a broader reading list:
// http://www.dcs.gla.ac.uk/~jsinger/ssa.html.
//
// The level of abstraction of the SSA form is intentionally close to
// the source language to facilitate construction of source analysis
// tools.  It is not intended for machine code generation.
//
// All looping, branching and switching constructs are replaced with
// unstructured control flow.  Higher-level control flow constructs
// such as multi-way branch can be reconstructed as needed; see
// ssautil.Switches() for an example.
//
// To construct an SSA-form program, call ssa.Create on a
// loader.Program, a set of type-checked packages created from
// parsed Go source files.  The resulting ssa.Program contains all the
// packages and their members, but SSA code is not created for
// function bodies until a subsequent call to (*Package).Build.
//
// The builder initially builds a naive SSA form in which all local
// variables are addresses of stack locations with explicit loads and
// stores.  Registerisation of eligible locals and φ-node insertion
// using dominance and dataflow are then performed as a second pass
// called "lifting" to improve the accuracy and performance of
// subsequent analyses; this pass can be skipped by setting the
// NaiveForm builder flag.
//
// The primary interfaces of this package are:
//
//    - Member: a named member of a Go package.
//    - Value: an expression that yields a value.
//    - Instruction: a statement that consumes values and performs computation.
//
// A computation that yields a result implements both the Value and
// Instruction interfaces.  The following table shows for each
// concrete type which of these interfaces it implements.
//
//                      Value?          Instruction?    Member?
//   *Alloc             ✔               ✔
//   *BinOp             ✔               ✔
//   *Builtin           ✔
//   *Call              ✔               ✔
//   *Capture           ✔
//   *ChangeInterface   ✔               ✔
//   *ChangeType        ✔               ✔
//   *Const             ✔
//   *Convert           ✔               ✔
//   *DebugRef                          ✔
//   *Defer                             ✔
//   *Extract           ✔               ✔
//   *Field             ✔               ✔
//   *FieldAddr         ✔               ✔
//   *Function          ✔                               ✔ (func)
//   *Global            ✔                               ✔ (var)
//   *Go                                ✔
//   *If                                ✔
//   *Index             ✔               ✔
//   *IndexAddr         ✔               ✔
//   *Jump                              ✔
//   *Lookup            ✔               ✔
//   *MakeChan          ✔               ✔
//   *MakeClosure       ✔               ✔
//   *MakeInterface     ✔               ✔
//   *MakeMap           ✔               ✔
//   *MakeSlice         ✔               ✔
//   *MapUpdate                         ✔
//   *NamedConst                                        ✔ (const)
//   *Next              ✔               ✔
//   *Panic                             ✔
//   *Parameter         ✔
//   *Phi               ✔               ✔
//   *Range             ✔               ✔
//   *Return                            ✔
//   *RunDefers                         ✔
//   *Select            ✔               ✔
//   *Send                              ✔
//   *Slice             ✔               ✔
//   *Store                             ✔
//   *Type                                              ✔ (type)
//   *TypeAssert        ✔               ✔
//   *UnOp              ✔               ✔
//
// Other key types in this package include: Program, Package, Function
// and BasicBlock.
//
// The program representation constructed by this package is fully
// resolved internally, i.e. it does not rely on the names of Values,
// Packages, Functions, Types or BasicBlocks for the correct
// interpretation of the program.  Only the identities of objects and
// the topology of the SSA and type graphs are semantically
// significant.  (There is one exception: Ids, used to identify field
// and method names, contain strings.)  Avoidance of name-based
// operations simplifies the implementation of subsequent passes and
// can make them very efficient.  Many objects are nonetheless named
// to aid in debugging, but it is not essential that the names be
// either accurate or unambiguous.  The public API exposes a number of
// name-based maps for client convenience.
//
// The ssa/ssautil package provides various utilities that depend only
// on the public API of this package.
//
// TODO(adonovan): Consider the exceptional control-flow implications
// of defer and recover().
//
// TODO(adonovan): write a how-to document for all the various cases
// of trying to determine corresponding elements across the four
// domains of source locations, ast.Nodes, types.Objects,
// ssa.Values/Instructions.
//
package ssa