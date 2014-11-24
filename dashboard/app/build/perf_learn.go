// antha-tools/dashboard/app/build/perf_learn.go: Part of the Antha language
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
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"sort"

	"appengine"
	"appengine/datastore"
)

func init() {
	http.HandleFunc("/perflearn", perfLearnHandler)
}

const (
	learnPercentile       = 0.95
	learnSignalMultiplier = 1.2
	learnMinSignal        = 0.5
)

func perfLearnHandler(w http.ResponseWriter, r *http.Request) {
	d := dashboardForRequest(r)
	c := d.Context(appengine.NewContext(r))

	pc, err := GetPerfConfig(c, r)
	if err != nil {
		logErr(w, r, err)
		return
	}

	p, err := GetPackage(c, "")
	if err != nil {
		logErr(w, r, err)
		return
	}

	update := r.FormValue("update") != ""
	noise := make(map[string]string)

	data := &perfLearnData{}

	commits, err := GetCommits(c, 0, p.NextNum)
	if err != nil {
		logErr(w, r, err)
		return
	}

	for _, builder := range pc.BuildersForBenchmark("") {
		for _, benchmark := range pc.BenchmarksForBuilder(builder) {
			for _, metric := range pc.MetricsForBenchmark(benchmark) {
				for _, procs := range pc.ProcList(builder) {
					values, err := GetPerfMetricsForCommits(c, builder, fmt.Sprintf("%v-%v", benchmark, procs), metric, 0, p.NextNum)
					if err != nil {
						logErr(w, r, err)
						return
					}
					var dd []float64
					last := uint64(0)
					for i, v := range values {
						if v == 0 {
							if commits[i].NeedsBenchmarking {
								last = 0
							}
							continue
						}
						if last != 0 {
							v1 := v
							if v1 < last {
								v1, last = last, v1
							}
							diff := float64(v1)/float64(last)*100 - 100
							dd = append(dd, diff)
						}
						last = v
					}
					if len(dd) == 0 {
						continue
					}
					sort.Float64s(dd)

					baseIdx := int(float64(len(dd)) * learnPercentile)
					baseVal := dd[baseIdx]
					signalVal := baseVal * learnSignalMultiplier
					if signalVal < learnMinSignal {
						signalVal = learnMinSignal
					}
					signalIdx := -1
					noiseNum := 0
					signalNum := 0

					var diffs []*perfLearnDiff
					for i, d := range dd {
						if d > 3*signalVal {
							d = 3 * signalVal
						}
						diffs = append(diffs, &perfLearnDiff{Num: i, Val: d})
						if signalIdx == -1 && d >= signalVal {
							signalIdx = i
						}
						if d < signalVal {
							noiseNum++
						} else {
							signalNum++
						}
					}
					diffs[baseIdx].Hint = "95%"
					if signalIdx != -1 {
						diffs[signalIdx].Hint = "signal"
					}
					diffs = diffs[len(diffs)*4/5:]
					name := fmt.Sprintf("%v/%v-%v/%v", builder, benchmark, procs, metric)
					data.Entries = append(data.Entries, &perfLearnEntry{len(data.Entries), name, baseVal, noiseNum, signalVal, signalNum, diffs})

					if len(dd) >= 100 || r.FormValue("force") != "" {
						nname := fmt.Sprintf("%v|%v-%v", builder, benchmark, procs)
						n := noise[nname] + fmt.Sprintf("|%v=%.2f", metric, signalVal)
						noise[nname] = n
					}
				}
			}
		}
	}

	if update {
		var noiseLevels []string
		for k, v := range noise {
			noiseLevels = append(noiseLevels, k+v)
		}
		tx := func(c appengine.Context) error {
			pc, err := GetPerfConfig(c, r)
			if err != nil {
				return err
			}
			pc.NoiseLevels = noiseLevels
			if _, err := datastore.Put(c, PerfConfigKey(c), pc); err != nil {
				return fmt.Errorf("putting PerfConfig: %v", err)
			}
			return nil
		}
		if err := datastore.RunInTransaction(c, tx, nil); err != nil {
			logErr(w, r, err)
			return
		}
	}

	var buf bytes.Buffer
	if err := perfLearnTemplate.Execute(&buf, data); err != nil {
		logErr(w, r, err)
		return
	}

	buf.WriteTo(w)
}

var perfLearnTemplate = template.Must(
	template.New("perf_learn.html").Funcs(tmplFuncs).ParseFiles("build/perf_learn.html"),
)

type perfLearnData struct {
	Entries []*perfLearnEntry
}

type perfLearnEntry struct {
	Num       int
	Name      string
	BaseVal   float64
	NoiseNum  int
	SignalVal float64
	SignalNum int
	Diffs     []*perfLearnDiff
}

type perfLearnDiff struct {
	Num  int
	Val  float64
	Hint string
}