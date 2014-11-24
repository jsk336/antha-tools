// antha-tools/dashboard/app/build/dash.go: Part of the Antha language
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


// +build appengine

package build

import (
	"net/http"
	"strings"

	"appengine"
)

// Dashboard describes a unique build dashboard.
type Dashboard struct {
	Name     string     // This dashboard's name and namespace
	RelPath  string     // The relative url path
	Packages []*Package // The project's packages to build
}

// dashboardForRequest returns the appropriate dashboard for a given URL path.
func dashboardForRequest(r *http.Request) *Dashboard {
	if strings.HasPrefix(r.URL.Path, gccgoDash.RelPath) {
		return gccgoDash
	}
	return goDash
}

// Context returns a namespaced context for this dashboard, or panics if it
// fails to create a new context.
func (d *Dashboard) Context(c appengine.Context) appengine.Context {
	// No namespace needed for the original Go dashboard.
	if d.Name == "Go" {
		return c
	}
	n, err := appengine.Namespace(c, d.Name)
	if err != nil {
		panic(err)
	}
	return n
}

// the currently known dashboards.
var dashboards = []*Dashboard{goDash, gccgoDash}

// goDash is the dashboard for the main antha repository.
var goDash = &Dashboard{
	Name:     "Go",
	RelPath:  "/",
	Packages: goPackages,
}

// goPackages is a list of all of the packages built by the main antha repository.
var goPackages = []*Package{
	{
		Kind: "go",
		Name: "Go",
	},
	{
		Kind: "subrepo",
		Name: "go.blog",
		Path: "code.google.com/p/go.blog",
	},
	{
		Kind: "subrepo",
		Name: "go.codereview",
		Path: "code.google.com/p/go.codereview",
	},
	{
		Kind: "subrepo",
		Name: "go.crypto",
		Path: "code.google.com/p/go.crypto",
	},
	{
		Kind: "subrepo",
		Name: "go.exp",
		Path: "code.google.com/p/go.exp",
	},
	{
		Kind: "subrepo",
		Name: "go.image",
		Path: "code.google.com/p/go.image",
	},
	{
		Kind: "subrepo",
		Name: "go.net",
		Path: "code.google.com/p/go.net",
	},
	{
		Kind: "subrepo",
		Name: "go.talks",
		Path: "code.google.com/p/go.talks",
	},
	{
		Kind: "subrepo",
		Name: "go.tools",
		Path: "github.com/antha-lang/antha-tools",
	},
}

// gccgoDash is the dashboard for gccgo.
var gccgoDash = &Dashboard{
	Name:    "Gccgo",
	RelPath: "/gccgo/",
	Packages: []*Package{
		{
			Kind: "gccgo",
			Name: "Gccgo",
		},
	},
}