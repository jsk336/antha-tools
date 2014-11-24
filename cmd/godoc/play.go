// antha-tools/cmd/godoc/play.go: Part of the Antha language
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
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/format"
	"net/http"

	// This package registers "/compile" and "/share" handlers
	// that redirect to the golang.org playground.
	_ "github.com/antha-lang/antha-tools/playground"
)

func init() {
	http.HandleFunc("/fmt", fmtHandler)
}

type fmtResponse struct {
	Body  string
	Error string
}

// fmtHandler takes a Go program in its "body" form value, formats it with
// standard gofmt formatting, and writes a fmtResponse as a JSON object.
func fmtHandler(w http.ResponseWriter, r *http.Request) {
	resp := new(fmtResponse)
	body, err := format.Source([]byte(r.FormValue("body")))
	if err != nil {
		resp.Error = err.Error()
	} else {
		resp.Body = string(body)
	}
	json.NewEncoder(w).Encode(resp)
}

// disabledHandler serves a 501 "Not Implemented" response.
func disabledHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprint(w, "This functionality is not available via local godoc.")
}