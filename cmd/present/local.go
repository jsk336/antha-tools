// antha-tools/cmd/present/local.go: Part of the Antha language
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
	"flag"
	"fmt"
	"github.com/antha-lang/antha/build"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"

	"github.com/antha-lang/antha-tools/playground/socket"
	"github.com/antha-lang/antha-tools/present"
)

const basePkg = "github.com/antha-lang/antha-tools/cmd/present"

var basePath string

func main() {
	httpAddr := flag.String("http", "127.0.0.1:3999", "HTTP service address (e.g., '127.0.0.1:3999')")
	originHost := flag.String("orighost", "", "host component of web origin URL (e.g., 'localhost')")
	flag.StringVar(&basePath, "base", "", "base path for slide template and static resources")
	flag.BoolVar(&present.PlayEnabled, "play", true, "enable playground (permit execution of arbitrary user code)")
	nativeClient := flag.Bool("nacl", false, "use Native Client environment playground (prevents non-Go code execution)")
	flag.Parse()

	if basePath == "" {
		p, err := build.Default.Import(basePkg, "", build.FindOnly)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't find gopresent files: %v\n", err)
			fmt.Fprintf(os.Stderr, basePathMessage, basePkg)
			os.Exit(1)
		}
		basePath = p.Dir
	}

	ln, err := net.Listen("tcp", *httpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	_, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		log.Fatal(err)
	}
	origin := &url.URL{Scheme: "http"}
	if *originHost != "" {
		origin.Host = net.JoinHostPort(*originHost, port)
	} else if ln.Addr().(*net.TCPAddr).IP.IsUnspecified() {
		name, _ := os.Hostname()
		origin.Host = net.JoinHostPort(name, port)
	} else {
		reqHost, reqPort, err := net.SplitHostPort(*httpAddr)
		if err != nil {
			log.Fatal(err)
		}
		if reqPort == "0" {
			origin.Host = net.JoinHostPort(reqHost, port)
		} else {
			origin.Host = *httpAddr
		}
	}

	if present.PlayEnabled {
		if *nativeClient {
			socket.RunScripts = false
			socket.Environ = func() []string {
				if runtime.GOARCH == "amd64" {
					return environ("GOOS=nacl", "GOARCH=amd64p32")
				}
				return environ("GOOS=nacl")
			}
		}
		playScript(basePath, "SocketTransport")
		http.Handle("/socket", socket.NewHandler(origin))
	}
	http.Handle("/static/", http.FileServer(http.Dir(basePath)))

	if !ln.Addr().(*net.TCPAddr).IP.IsLoopback() &&
		present.PlayEnabled && !*nativeClient {
		log.Print(localhostWarning)
	}

	log.Printf("Open your web browser and visit %s", origin.String())
	log.Fatal(http.Serve(ln, nil))
}

func playable(c present.Code) bool {
	return present.PlayEnabled && c.Play
}

func environ(vars ...string) []string {
	env := os.Environ()
	for _, r := range vars {
		k := strings.SplitAfter(r, "=")[0]
		var found bool
		for i, v := range env {
			if strings.HasPrefix(v, k) {
				env[i] = r
				found = true
			}
		}
		if !found {
			env = append(env, r)
		}
	}
	return env
}

const basePathMessage = `
By default, gopresent locates the slide template files and associated
static content by looking for a %q package
in your Go workspaces (GOPATH).

You may use the -base flag to specify an alternate location.
`

const localhostWarning = `
WARNING!  WARNING!  WARNING!

The present server appears to be listening on an address that is not localhost.
Anyone with access to this address and port will have access to this machine as
the user running present.

To avoid this message, listen on localhost or run with -play=false.

If you don't understand this message, hit Control-C to terminate this process.

WARNING!  WARNING!  WARNING!
`