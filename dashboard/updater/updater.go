// antha-tools/dashboard/updater/updater.go: Part of the Antha language
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
	"bytes"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

var (
	builder   = flag.String("builder", "", "builder name")
	key       = flag.String("key", "", "builder key")
	gopath    = flag.String("gopath", "", "path to antha repo")
	dashboard = flag.String("dashboard", "build.golang.org", "Go Dashboard Host")
	batch     = flag.Int("batch", 100, "upload batch size")
)

// Do not benchmark beyond this commit.
// There is little sense in benchmarking till first commit,
// and the benchmark won't build anyway.
const Go1Commit = "0051c7442fed" // test/bench/shootout: update timing.log to Go 1.

// HgLog represents a single Mercurial revision.
type HgLog struct {
	Hash   string
	Branch string
	Files  string
}

func main() {
	flag.Parse()
	logs := hgLog()
	var hashes []string
	ngo1 := 0
	for i := range logs {
		if strings.HasPrefix(logs[i].Hash, Go1Commit) {
			break
		}
		if needsBenchmarking(&logs[i]) {
			hashes = append(hashes, logs[i].Hash)
		}
		ngo1++
	}
	fmt.Printf("found %v commits, %v after Go1, %v need benchmarking\n", len(logs), ngo1, len(hashes))
	for i := 0; i < len(hashes); i += *batch {
		j := i + *batch
		if j > len(hashes) {
			j = len(hashes)
		}
		fmt.Printf("sending %v-%v... ", i, j)
		res := postCommits(hashes[i:j])
		fmt.Printf("%s\n", res)
	}
}

func hgLog() []HgLog {
	var out bytes.Buffer
	cmd := exec.Command("hg", "log", "--encoding=utf-8", "--template", xmlLogTemplate)
	cmd.Dir = *gopath
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("failed to execute 'hg log': %v\n", err)
		os.Exit(1)
	}
	var top struct{ Log []HgLog }
	err = xml.Unmarshal([]byte("<Top>"+out.String()+"</Top>"), &top)
	if err != nil {
		fmt.Printf("failed to parse log: %v\n", err)
		os.Exit(1)
	}
	return top.Log
}

func needsBenchmarking(log *HgLog) bool {
	if log.Branch != "" {
		return false
	}
	for _, f := range strings.Split(log.Files, " ") {
		if (strings.HasPrefix(f, "include") || strings.HasPrefix(f, "src")) &&
			!strings.HasSuffix(f, "_test.go") && !strings.Contains(f, "testdata") {
			return true
		}
	}
	return false
}

func postCommits(hashes []string) string {
	args := url.Values{"builder": {*builder}, "key": {*key}}
	cmd := fmt.Sprintf("http://%v/updatebenchmark?%v", *dashboard, args.Encode())
	b, err := json.Marshal(hashes)
	if err != nil {
		return fmt.Sprintf("failed to encode request: %v\n", err)
	}
	r, err := http.Post(cmd, "text/json", bytes.NewReader(b))
	if err != nil {
		return fmt.Sprintf("failed to send http request: %v\n", err)
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return fmt.Sprintf("http request failed: %v\n", r.Status)
	}
	resp, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Sprintf("failed to read http response: %v\n", err)
	}
	return string(resp)
}

const xmlLogTemplate = `
        <Log>
        <Hash>{node|escape}</Hash>
        <Branch>{branches}</Branch>
        <Files>{files}</Files>
        </Log>
`