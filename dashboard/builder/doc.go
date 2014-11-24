// antha-tools/dashboard/builder/doc.go: Part of the Antha language
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

Go Builder is a continuous build client for the Go project.
It integrates with the Go Dashboard AppEngine application.

Go Builder is intended to run continuously as a background process.

It periodically pulls updates from the Go Mercurial repository.

When a newer revision is found, Go Builder creates a clone of the repository,
runs all.bash, and reports build success or failure to the Go Dashboard.

For a release revision (a change description that matches "release.YYYY-MM-DD"),
Go Builder will create a tar.gz archive of the GOROOT and deliver it to the
Go Google Code project's downloads section.

Usage:

  gobuilder goos-goarch...

  Several goos-goarch combinations can be provided, and the builder will
  build them in serial.

Optional flags:

  -dashboard="godashboard.appspot.com": Go Dashboard Host
    The location of the Go Dashboard application to which Go Builder will
    report its results.

  -release: Build and deliver binary release archive

  -rev=N: Build revision N and exit

  -cmd="./all.bash": Build command (specify absolute or relative to go/src)

  -v: Verbose logging

  -external: External package builder mode (will not report Go build
     state to dashboard or issue releases)

The key file should be located at $HOME/.gobuildkey or, for a builder-specific
key, $HOME/.gobuildkey-$BUILDER (eg, $HOME/.gobuildkey-linux-amd64).

The build key file is a text file of the format:

  godashboard-key
  googlecode-username
  googlecode-password

If the Google Code credentials are not provided the archival step
will be skipped.

*/
package main