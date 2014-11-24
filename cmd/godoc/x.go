// antha-tools/cmd/godoc/x.go: Part of the Antha language
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


// This file contains the handlers that serve go-import redirects for Go
// sub-repositories. It specifies the mapping from import paths like
// "golang.org/x/tools" to the actual repository locations.

package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

const xPrefix = "/x/"

var xMap = map[string]string{
	"benchmarks": "https://code.google.com/p/go.benchmarks",
	"blog":       "https://code.google.com/p/go.blog",
	"codereview": "https://code.google.com/p/go.codereview",
	"crypto":     "https://code.google.com/p/go.crypto",
	"exp":        "https://code.google.com/p/go.exp",
	"image":      "https://code.google.com/p/go.image",
	"mobile":     "https://code.google.com/p/go.mobile",
	"net":        "https://code.google.com/p/go.net",
	"sys":        "https://code.google.com/p/go.sys",
	"talks":      "https://code.google.com/p/go.talks",
	"text":       "https://code.google.com/p/go.text",
	"tools":      "https://code.google.com/p/go.tools",
}

func init() {
	http.HandleFunc(xPrefix, xHandler)
}

func xHandler(w http.ResponseWriter, r *http.Request) {
	head, tail := strings.TrimPrefix(r.URL.Path, xPrefix), ""
	if i := strings.Index(head, "/"); i != -1 {
		head, tail = head[:i], head[i:]
	}
	repo, ok := xMap[head]
	if !ok {
		http.NotFound(w, r)
		return
	}
	data := struct {
		Prefix, Head, Tail, Repo string
	}{xPrefix, head, tail, repo}
	if err := xTemplate.Execute(w, data); err != nil {
		log.Println("xHandler:", err)
	}
}

var xTemplate = template.Must(template.New("x").Parse(`<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="go-import" content="golang.org{{.Prefix}}{{.Head}} hg {{.Repo}}">
<meta http-equiv="refresh" content="0; url=https://godoc.org/golang.org{{.Prefix}}{{.Head}}{{.Tail}}">
</head>
<body>
Nothing to see here; <a href="https://godoc.org/golang.org{{.Prefix}}{{.Head}}{{.Tail}}">move along</a>.
</body>
</html>
`))