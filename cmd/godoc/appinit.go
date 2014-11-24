// antha-tools/cmd/godoc/appinit.go: Part of the Antha language
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

package main

// This file replaces main.go when running godoc under app-engine.
// See README.godoc-app for details.

import (
	"archive/zip"
	"log"
	"path"

	"github.com/antha-lang/antha-tools/anthadoc"
	"github.com/antha-lang/antha-tools/anthadoc/static"
	"github.com/antha-lang/antha-tools/anthadoc/vfs"
	"github.com/antha-lang/antha-tools/anthadoc/vfs/mapfs"
	"github.com/antha-lang/antha-tools/anthadoc/vfs/zipfs"
)

func init() {
	playEnabled = true

	log.Println("initializing godoc ...")
	log.Printf(".zip file   = %s", zipFilename)
	log.Printf(".zip GOROOT = %s", zipGoroot)
	log.Printf("index files = %s", indexFilenames)

	goroot := path.Join("/", zipGoroot) // fsHttp paths are relative to '/'

	// read .zip file and set up file systems
	const zipfile = zipFilename
	rc, err := zip.OpenReader(zipfile)
	if err != nil {
		log.Fatalf("%s: %s\n", zipfile, err)
	}
	// rc is never closed (app running forever)
	fs.Bind("/", zipfs.New(rc, zipFilename), goroot, vfs.BindReplace)
	fs.Bind("/lib/godoc", mapfs.New(static.Files), "/", vfs.BindReplace)

	corpus := godoc.NewCorpus(fs)
	corpus.Verbose = false
	corpus.MaxResults = 10000 // matches flag default in main.go
	corpus.IndexEnabled = true
	corpus.IndexFiles = indexFilenames
	if err := corpus.Init(); err != nil {
		log.Fatal(err)
	}
	go corpus.RunIndexer()

	pres = godoc.NewPresentation(corpus)
	pres.TabWidth = 8
	pres.ShowPlayground = true
	pres.ShowExamples = true
	pres.DeclLinks = true

	readTemplates(pres, true)
	registerHandlers(pres)

	log.Println("godoc initialization complete")
}