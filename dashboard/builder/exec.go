// antha-tools/dashboard/builder/exec.go: Part of the Antha language
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
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

// run is a simple wrapper for exec.Run/Close
func run(d time.Duration, envv []string, dir string, argv ...string) error {
	if *verbose {
		log.Println("run", argv)
	}
	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Dir = dir
	cmd.Env = envv
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	err := timeout(d, cmd.Wait)
	if _, ok := err.(timeoutError); ok {
		cmd.Process.Kill()
	}
	return err
}

// runLog runs a process and returns the combined stdout/stderr. It returns
// process combined stdout and stderr output, exit status and error. The
// error returned is nil, if process is started successfully, even if exit
// status is not successful.
func runLog(timeout time.Duration, envv []string, dir string, argv ...string) (string, bool, error) {
	var b bytes.Buffer
	ok, err := runOutput(timeout, envv, &b, dir, argv...)
	return b.String(), ok, err
}

// runOutput runs a process and directs any output to the supplied writer.
// It returns exit status and error. The error returned is nil, if process
// is started successfully, even if exit status is not successful.
func runOutput(d time.Duration, envv []string, out io.Writer, dir string, argv ...string) (bool, error) {
	if *verbose {
		log.Println("runOutput", argv)
	}

	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Dir = dir
	cmd.Env = envv
	cmd.Stdout = out
	cmd.Stderr = out

	startErr := cmd.Start()
	if startErr != nil {
		return false, startErr
	}

	if err := timeout(d, cmd.Wait); err != nil {
		if _, ok := err.(timeoutError); ok {
			cmd.Process.Kill()
		}
		return false, err
	}
	return true, nil
}

// timeout runs f and returns its error value, or if the function does not
// complete before the provided duration it returns a timeout error.
func timeout(d time.Duration, f func() error) error {
	errc := make(chan error, 1)
	go func() {
		errc <- f()
	}()
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return timeoutError(d)
	case err := <-errc:
		return err
	}
}

type timeoutError time.Duration

func (e timeoutError) Error() string {
	return fmt.Sprintf("timed out after %v", time.Duration(e))
}