// antha-tools/anthadoc/server_test.go: Part of the Antha language
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

package anthadoc

import (
	"errors"
	"expvar"
	"net/http"
	"net/http/httptest"
	"testing"
	"text/template"
)

var (
	// NOTE: with no plain-text in the template, template.Execute will not
	// return an error when http.ResponseWriter.Write does return an error.
	tmpl = template.Must(template.New("test").Parse("{{.Foo}}"))
)

type withFoo struct {
	Foo int
}

type withoutFoo struct {
}

type errResponseWriter struct {
}

func (*errResponseWriter) Header() http.Header {
	return http.Header{}
}

func (*errResponseWriter) WriteHeader(int) {
}

func (*errResponseWriter) Write(p []byte) (int, error) {
	return 0, errors.New("error")
}

func TestApplyTemplateToResponseWriter(t *testing.T) {
	for _, tc := range []struct {
		desc    string
		rw      http.ResponseWriter
		data    interface{}
		expVars int
	}{
		{
			desc:    "no error",
			rw:      &httptest.ResponseRecorder{},
			data:    &withFoo{},
			expVars: 0,
		},
		{
			desc:    "template error",
			rw:      &httptest.ResponseRecorder{},
			data:    &withoutFoo{},
			expVars: 0,
		},
		{
			desc:    "ResponseWriter error",
			rw:      &errResponseWriter{},
			data:    &withFoo{},
			expVars: 1,
		},
	} {
		httpErrors.Init()
		applyTemplateToResponseWriter(tc.rw, tmpl, tc.data)
		gotVars := 0
		httpErrors.Do(func(expvar.KeyValue) {
			gotVars++
		})
		if gotVars != tc.expVars {
			t.Errorf("applyTemplateToResponseWriter(%q): got %d vars, want %d", tc.desc, gotVars, tc.expVars)
		}
	}
}