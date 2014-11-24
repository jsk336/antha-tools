// antha-tools/cmd/present/doc.go: Part of the Antha language
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


/*
Present displays slide presentations and articles. It runs a web server that
presents slide and article files from the current directory.

It may be run as a stand-alone command or an App Engine app.
Instructions for deployment to App Engine are in the README of the
antha-tools repository.

Usage of present:
  -base="": base path for slide template and static resources
  -http="127.0.0.1:3999": host:port to listen on

Input files are named foo.extension, where "extension" defines the format of
the generated output. The supported formats are:
	.slide        // HTML5 slide presentation
	.article      // article format, such as a blog post

The present file format is documented by the present package:
http://godoc.org/antha-tools/present
*/
package main