// antha-tools/cmd/godoc/remotesearch.go: Part of the Antha language
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


// +build !appengine

package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func handleRemoteSearch() {
	// Command-line queries.
	for i := 0; i < flag.NArg(); i++ {
		res, err := remoteSearch(flag.Arg(i))
		if err != nil {
			log.Fatalf("remoteSearch: %s", err)
		}
		io.Copy(os.Stdout, res.Body)
	}
	return
}

// remoteSearchURL returns the search URL for a given query as needed by
// remoteSearch. If html is set, an html result is requested; otherwise
// the result is in textual form.
// Adjust this function as necessary if modeNames or FormValue parameters
// change.
func remoteSearchURL(query string, html bool) string {
	s := "/search?m=text&q="
	if html {
		s = "/search?q="
	}
	return s + url.QueryEscape(query)
}

func remoteSearch(query string) (res *http.Response, err error) {
	// list of addresses to try
	var addrs []string
	if *serverAddr != "" {
		// explicit server address - only try this one
		addrs = []string{*serverAddr}
	} else {
		addrs = []string{
			defaultAddr,
			"golang.org",
		}
	}

	// remote search
	search := remoteSearchURL(query, *html)
	for _, addr := range addrs {
		url := "http://" + addr + search
		res, err = http.Get(url)
		if err == nil && res.StatusCode == http.StatusOK {
			break
		}
	}

	if err == nil && res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
	}

	return
}