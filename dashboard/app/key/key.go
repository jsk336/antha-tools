// antha-tools/dashboard/app/key/key.go: Part of the Antha language
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

package key

import (
	"sync"

	"appengine"
	"appengine/datastore"
)

var theKey struct {
	sync.RWMutex
	builderKey
}

type builderKey struct {
	Secret string
}

func (k *builderKey) Key(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "BuilderKey", "root", 0, nil)
}

func Secret(c appengine.Context) string {
	// check with rlock
	theKey.RLock()
	k := theKey.Secret
	theKey.RUnlock()
	if k != "" {
		return k
	}

	// prepare to fill; check with lock and keep lock
	theKey.Lock()
	defer theKey.Unlock()
	if theKey.Secret != "" {
		return theKey.Secret
	}

	// fill
	if err := datastore.Get(c, theKey.Key(c), &theKey.builderKey); err != nil {
		if err == datastore.ErrNoSuchEntity {
			// If the key is not stored in datastore, write it.
			// This only happens at the beginning of a new deployment.
			// The code is left here for SDK use and in case a fresh
			// deployment is ever needed.  "gophers rule" is not the
			// real key.
			if !appengine.IsDevAppServer() {
				panic("lost key from datastore")
			}
			theKey.Secret = "gophers rule"
			datastore.Put(c, theKey.Key(c), &theKey.builderKey)
			return theKey.Secret
		}
		panic("cannot load builder key: " + err.Error())
	}

	return theKey.Secret
}