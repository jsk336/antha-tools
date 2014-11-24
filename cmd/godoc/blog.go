// antha-tools/cmd/godoc/blog.go: Part of the Antha language
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
	"fmt"
	"github.com/antha-lang/antha/build"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/antha-lang/antha-tools/blog"
	"github.com/antha-lang/antha-tools/anthadoc/redirect"
)

const (
	blogRepo = "code.google.com/p/go.blog"
	blogURL  = "http://blog.golang.org/"
	blogPath = "/blog/"
)

var (
	blogServer   http.Handler // set by blogInit
	blogInitOnce sync.Once
	playEnabled  bool
)

func init() {
	// Initialize blog only when first accessed.
	http.HandleFunc(blogPath, func(w http.ResponseWriter, r *http.Request) {
		blogInitOnce.Do(blogInit)
		blogServer.ServeHTTP(w, r)
	})
}

func blogInit() {
	// Binary distributions will include the blog content in "/blog".
	root := filepath.Join(runtime.GOROOT(), "blog")

	// Prefer content from go.blog repository if present.
	if pkg, err := build.Import(blogRepo, "", build.FindOnly); err == nil {
		root = pkg.Dir
	}

	// If content is not available fall back to redirect.
	if fi, err := os.Stat(root); err != nil || !fi.IsDir() {
		fmt.Fprintf(os.Stderr, "Blog content not available locally. "+
			"To install, run \n\tgo get %v\n", blogRepo)
		blogServer = http.HandlerFunc(blogRedirectHandler)
		return
	}

	s, err := blog.NewServer(blog.Config{
		BaseURL:      blogPath,
		BasePath:     strings.TrimSuffix(blogPath, "/"),
		ContentPath:  filepath.Join(root, "content"),
		TemplatePath: filepath.Join(root, "template"),
		HomeArticles: 5,
		PlayEnabled:  playEnabled,
	})
	if err != nil {
		log.Fatal(err)
	}
	blogServer = s
}

func blogRedirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == blogPath {
		http.Redirect(w, r, blogURL, http.StatusFound)
		return
	}
	blogPrefixHandler.ServeHTTP(w, r)
}

var blogPrefixHandler = redirect.PrefixHandler(blogPath, blogURL)