// antha-tools/dashboard/app/cache/cache.go: Part of the Antha language
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

package cache

import (
	"fmt"
	"net/http"
	"time"

	"appengine"
	"appengine/memcache"
)

// TimeKey specifies the memcache entity that keeps the logical datastore time.
var TimeKey = "cachetime"

const (
	nocache = "nocache"
	expiry  = 600 // 10 minutes
)

func newTime() uint64 { return uint64(time.Now().Unix()) << 32 }

// Now returns the current logical datastore time to use for cache lookups.
func Now(c appengine.Context) uint64 {
	t, err := memcache.Increment(c, TimeKey, 0, newTime())
	if err != nil {
		c.Errorf("cache.Now: %v", err)
		return 0
	}
	return t
}

// Tick sets the current logical datastore time to a never-before-used time
// and returns that time. It should be called to invalidate the cache.
func Tick(c appengine.Context) uint64 {
	t, err := memcache.Increment(c, TimeKey, 1, newTime())
	if err != nil {
		c.Errorf("cache.Tick: %v", err)
		return 0
	}
	return t
}

// Get fetches data for name at time now from memcache and unmarshals it into
// value. It reports whether it found the cache record and logs any errors to
// the admin console.
func Get(r *http.Request, now uint64, name string, value interface{}) bool {
	if now == 0 || r.FormValue(nocache) != "" {
		return false
	}
	c := appengine.NewContext(r)
	key := fmt.Sprintf("%s.%d", name, now)
	_, err := memcache.JSON.Get(c, key, value)
	if err == nil {
		c.Debugf("cache hit %q", key)
		return true
	}
	c.Debugf("cache miss %q", key)
	if err != memcache.ErrCacheMiss {
		c.Errorf("get cache %q: %v", key, err)
	}
	return false
}

// Set puts value into memcache under name at time now.
// It logs any errors to the admin console.
func Set(r *http.Request, now uint64, name string, value interface{}) {
	if now == 0 || r.FormValue(nocache) != "" {
		return
	}
	c := appengine.NewContext(r)
	key := fmt.Sprintf("%s.%d", name, now)
	err := memcache.JSON.Set(c, &memcache.Item{
		Key:        key,
		Object:     value,
		Expiration: expiry,
	})
	if err != nil {
		c.Errorf("set cache %q: %v", key, err)
	}
}