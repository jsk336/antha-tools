// antha-tools/playground/common.go: Part of the Antha language
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


// Package playground registers HTTP handlers at "/compile" and "/share" that
// proxy requests to the golang.org playground service.
// This package may be used unaltered on App Engine.
package playground

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "http://play.golang.org"

func init() {
	http.HandleFunc("/compile", bounce)
	http.HandleFunc("/share", bounce)
}

func bounce(w http.ResponseWriter, r *http.Request) {
	b := new(bytes.Buffer)
	if err := passThru(b, r); err != nil {
		http.Error(w, "Server error.", http.StatusInternalServerError)
		report(r, err)
		return
	}
	io.Copy(w, b)
}

func passThru(w io.Writer, req *http.Request) error {
	defer req.Body.Close()
	url := baseURL + req.URL.Path
	r, err := client(req).Post(url, req.Header.Get("Content-type"), req.Body)
	if err != nil {
		return fmt.Errorf("making POST request: %v", err)
	}
	defer r.Body.Close()
	if _, err := io.Copy(w, r.Body); err != nil {
		return fmt.Errorf("copying response Body: %v", err)
	}
	return nil
}