// antha-tools/cmd/goimports/doc.go: Part of the Antha language
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

Command goimports updates your Go import lines,
adding missing ones and removing unreferenced ones.

     $ antha get antha-tools/cmd/goimports

It's a drop-in replacement for your editor's gofmt-on-save hook.
It has the same command-line interface as gofmt and formats
your code in the same way.

For emacs, make sure you have the latest (Go 1.2+) go-mode.el:
   https://go.googlecode.com/hg/misc/emacs/go-mode.el
Then in your .emacs file:
   (setq gofmt-command "goimports")
   (add-to-list 'load-path "/home/you/goroot/misc/emacs/")
   (require 'go-mode-load)
   (add-hook 'before-save-hook 'gofmt-before-save)

For vim, set "gofmt_command" to "goimports":
    https://code.google.com/p/antha/source/detail?r=39c724dd7f252
    https://code.google.com/p/antha/source/browse#hg%2Fmisc%2Fvim
    etc

For GoSublime, follow the steps described here:
    http://michaelwhatcott.com/gosublime-goimports/

For other editors, you probably know what to do.

Happy hacking!

*/
package main