// antha-tools/anthadoc/analysis/json.go: Part of the Antha language
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


package analysis

// This file defines types used by client-side JavaScript.

type anchorJSON struct {
	Text string // HTML
	Href string // URL
}

type commOpJSON struct {
	Op anchorJSON
	Fn string
}

// JavaScript's onClickComm() expects a commJSON.
type commJSON struct {
	Ops []commOpJSON
}

// Indicates one of these forms of fact about a type T:
// T "is implemented by <ByKind> type <Other>"  (ByKind != "", e.g. "array")
// T "implements <Other>"                       (ByKind == "")
type implFactJSON struct {
	ByKind string `json:",omitempty"`
	Other  anchorJSON
}

// Implements facts are grouped by form, for ease of reading.
type implGroupJSON struct {
	Descr string
	Facts []implFactJSON
}

// JavaScript's onClickIdent() expects a TypeInfoJSON.
type TypeInfoJSON struct {
	Name        string // type name
	Size, Align int64
	Methods     []anchorJSON
	ImplGroups  []implGroupJSON
}

// JavaScript's onClickCallees() expects a calleesJSON.
type calleesJSON struct {
	Descr   string
	Callees []anchorJSON // markup for called function
}

type callerJSON struct {
	Func  string
	Sites []anchorJSON
}

// JavaScript's onClickCallers() expects a callersJSON.
type callersJSON struct {
	Callee  string
	Callers []callerJSON
}

// JavaScript's cgAddChild requires a global array of PCGNodeJSON
// called CALLGRAPH, representing the intra-package call graph.
// The first element is special and represents "all external callers".
type PCGNodeJSON struct {
	Func    anchorJSON
	Callees []int // indices within CALLGRAPH of nodes called by this one
}