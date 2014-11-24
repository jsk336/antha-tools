// antha-tools/dashboard/app/build/init.go: Part of the Antha language
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
	"fmt"
	"net/http"

	"appengine"
	"appengine/datastore"

	"cache"
	"key"
)

func initHandler(w http.ResponseWriter, r *http.Request) {
	d := dashboardForRequest(r)
	c := d.Context(appengine.NewContext(r))
	defer cache.Tick(c)
	for _, p := range d.Packages {
		err := datastore.Get(c, p.Key(c), new(Package))
		if _, ok := err.(*datastore.ErrFieldMismatch); ok {
			// Some fields have been removed, so it's okay to ignore this error.
			err = nil
		}
		if err == nil {
			continue
		} else if err != datastore.ErrNoSuchEntity {
			logErr(w, r, err)
			return
		}
		if _, err := datastore.Put(c, p.Key(c), p); err != nil {
			logErr(w, r, err)
			return
		}
	}

	// Create secret key.
	key.Secret(c)

	fmt.Fprint(w, "OK")
}