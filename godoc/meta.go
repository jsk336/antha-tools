// antha-tools/godoc/meta.go: Part of the Antha language
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


package godoc

import (
	"bytes"
	"encoding/json"
	"log"
	pathpkg "path"
	"strings"
	"time"

	"github.com/antha-lang/antha-tools/anthadoc/vfs"
)

var (
	doctype   = []byte("<!DOCTYPE ")
	jsonStart = []byte("<!--{")
	jsonEnd   = []byte("}-->")
)

// ----------------------------------------------------------------------------
// Documentation Metadata

// TODO(adg): why are some exported and some aren't? -brad
type Metadata struct {
	Title    string
	Subtitle string
	Template bool   // execute as template
	Path     string // canonical path for this page
	filePath string // filesystem path relative to goroot
}

func (m *Metadata) FilePath() string { return m.filePath }

// extractMetadata extracts the Metadata from a byte slice.
// It returns the Metadata value and the remaining data.
// If no metadata is present the original byte slice is returned.
//
func extractMetadata(b []byte) (meta Metadata, tail []byte, err error) {
	tail = b
	if !bytes.HasPrefix(b, jsonStart) {
		return
	}
	end := bytes.Index(b, jsonEnd)
	if end < 0 {
		return
	}
	b = b[len(jsonStart)-1 : end+1] // drop leading <!-- and include trailing }
	if err = json.Unmarshal(b, &meta); err != nil {
		return
	}
	tail = tail[end+len(jsonEnd):]
	return
}

// UpdateMetadata scans $GOROOT/doc for HTML files, reads their metadata,
// and updates the DocMetadata map.
func (c *Corpus) updateMetadata() {
	metadata := make(map[string]*Metadata)
	var scan func(string) // scan is recursive
	scan = func(dir string) {
		fis, err := c.fs.ReadDir(dir)
		if err != nil {
			log.Println("updateMetadata:", err)
			return
		}
		for _, fi := range fis {
			name := pathpkg.Join(dir, fi.Name())
			if fi.IsDir() {
				scan(name) // recurse
				continue
			}
			if !strings.HasSuffix(name, ".html") {
				continue
			}
			// Extract metadata from the file.
			b, err := vfs.ReadFile(c.fs, name)
			if err != nil {
				log.Printf("updateMetadata %s: %v", name, err)
				continue
			}
			meta, _, err := extractMetadata(b)
			if err != nil {
				log.Printf("updateMetadata: %s: %v", name, err)
				continue
			}
			// Store relative filesystem path in Metadata.
			meta.filePath = name
			if meta.Path == "" {
				// If no Path, canonical path is actual path.
				meta.Path = meta.filePath
			}
			// Store under both paths.
			metadata[meta.Path] = &meta
			metadata[meta.filePath] = &meta
		}
	}
	scan("/doc")
	c.docMetadata.Set(metadata)
}

// MetadataFor returns the *Metadata for a given relative path or nil if none
// exists.
//
func (c *Corpus) MetadataFor(relpath string) *Metadata {
	if m, _ := c.docMetadata.Get(); m != nil {
		meta := m.(map[string]*Metadata)
		// If metadata for this relpath exists, return it.
		if p := meta[relpath]; p != nil {
			return p
		}
		// Try with or without trailing slash.
		if strings.HasSuffix(relpath, "/") {
			relpath = relpath[:len(relpath)-1]
		} else {
			relpath = relpath + "/"
		}
		return meta[relpath]
	}
	return nil
}

// refreshMetadata sends a signal to update DocMetadata. If a refresh is in
// progress the metadata will be refreshed again afterward.
//
func (c *Corpus) refreshMetadata() {
	select {
	case c.refreshMetadataSignal <- true:
	default:
	}
}

// RefreshMetadataLoop runs forever, updating DocMetadata when the underlying
// file system changes. It should be launched in a goroutine.
func (c *Corpus) refreshMetadataLoop() {
	for {
		<-c.refreshMetadataSignal
		c.updateMetadata()
		time.Sleep(10 * time.Second) // at most once every 10 seconds
	}
}