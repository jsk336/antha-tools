// antha-tools/cmd/vet/testdata/deadcode.go: Part of the Antha language
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


// +build ignore

// This file contains tests for the dead code checker.

package testdata

type T int

var x interface{}
var c chan int

func external() int // ok

func _() int {
}

func _() int {
	print(1)
}

func _() int {
	print(1)
	return 2
	println() // ERROR "unreachable code"
}

func _() int {
L:
	print(1)
	goto L
	println() // ERROR "unreachable code"
}

func _() int {
	print(1)
	panic(2)
	println() // ERROR "unreachable code"
}

// but only builtin panic
func _() int {
	var panic = func(int) {}
	print(1)
	panic(2)
	println() // ok
}

func _() int {
	{
		print(1)
		return 2
		println() // ERROR "unreachable code"
	}
	println() // ok
}

func _() int {
	{
		print(1)
		return 2
	}
	println() // ERROR "unreachable code"
}

func _() int {
L:
	{
		print(1)
		goto L
		println() // ERROR "unreachable code"
	}
	println() // ok
}

func _() int {
L:
	{
		print(1)
		goto L
	}
	println() // ERROR "unreachable code"
}

func _() int {
	print(1)
	{
		panic(2)
	}
}

func _() int {
	print(1)
	{
		panic(2)
		println() // ERROR "unreachable code"
	}
}

func _() int {
	print(1)
	{
		panic(2)
	}
	println() // ERROR "unreachable code"
}

func _() int {
	print(1)
	return 2
	{ // ERROR "unreachable code"
	}
}

func _() int {
L:
	print(1)
	goto L
	{ // ERROR "unreachable code"
	}
}

func _() int {
	print(1)
	panic(2)
	{ // ERROR "unreachable code"
	}
}

func _() int {
	{
		print(1)
		return 2
		{ // ERROR "unreachable code"
		}
	}
}

func _() int {
L:
	{
		print(1)
		goto L
		{ // ERROR "unreachable code"
		}
	}
}

func _() int {
	print(1)
	{
		panic(2)
		{ // ERROR "unreachable code"
		}
	}
}

func _() int {
	{
		print(1)
		return 2
	}
	{ // ERROR "unreachable code"
	}
}

func _() int {
L:
	{
		print(1)
		goto L
	}
	{ // ERROR "unreachable code"
	}
}

func _() int {
	print(1)
	{
		panic(2)
	}
	{ // ERROR "unreachable code"
	}
}

func _() int {
	print(1)
	if x == nil {
		panic(2)
	} else {
		panic(3)
	}
	println() // ERROR "unreachable code"
}

func _() int {
L:
	print(1)
	if x == nil {
		panic(2)
	} else {
		goto L
	}
	println() // ERROR "unreachable code"
}

func _() int {
L:
	print(1)
	if x == nil {
		panic(2)
	} else if x == 1 {
		return 0
	} else if x != 2 {
		panic(3)
	} else {
		goto L
	}
	println() // ERROR "unreachable code"
}

// if-else chain missing final else is not okay, even if the
// conditions cover every possible case.

func _() int {
	print(1)
	if x == nil {
		panic(2)
	} else if x != nil {
		panic(3)
	}
	println() // ok
}

func _() int {
	print(1)
	if x == nil {
		panic(2)
	}
	println() // ok
}

func _() int {
L:
	print(1)
	if x == nil {
		panic(2)
	} else if x == 1 {
		return 0
	} else if x != 1 {
		panic(3)
	}
	println() // ok
}

func _() int {
	print(1)
	for {
	}
	println() // ERROR "unreachable code"
}

func _() int {
	for {
		for {
			break
		}
	}
	println() // ERROR "unreachable code"
}

func _() int {
	for {
		for {
			break
			println() // ERROR "unreachable code"
		}
	}
}

func _() int {
	for {
		for {
			continue
			println() // ERROR "unreachable code"
		}
	}
}

func _() int {
	for {
	L:
		for {
			break L
		}
	}
	println() // ERROR "unreachable code"
}

func _() int {
	print(1)
	for {
		break
	}
	println() // ok
}

func _() int {
	for {
		for {
		}
		break // ERROR "unreachable code"
	}
	println() // ok
}

func _() int {
L:
	for {
		for {
			break L
		}
	}
	println() // ok
}

func _() int {
	print(1)
	for x == nil {
	}
	println() // ok
}

func _() int {
	for x == nil {
		for {
			break
		}
	}
	println() // ok
}

func _() int {
	for x == nil {
	L:
		for {
			break L
		}
	}
	println() // ok
}

func _() int {
	print(1)
	for true {
	}
	println() // ok
}

func _() int {
	for true {
		for {
			break
		}
	}
	println() // ok
}

func _() int {
	for true {
	L:
		for {
			break L
		}
	}
	println() // ok
}

func _() int {
	print(1)
	select {}
	println() // ERROR "unreachable code"
}

func _() int {
	print(1)
	select {
	case <-c:
		print(2)
		panic("abc")
		println() // ERROR "unreachable code"
	}
}

func _() int {
	print(1)
	select {
	case <-c:
		print(2)
		panic("abc")
	}
	println() // ERROR "unreachable code"
}

func _() int {
	print(1)
	select {
	case <-c:
		print(2)
		for {
		}
		println() // ERROR "unreachable code"
	}
}

func _() int {
	print(1)
	select {
	case <-c:
		print(2)
		for {
		}
	}
	println() // ERROR "unreachable code"
}

func _() int {
L:
	print(1)
	select {
	case <-c:
		print(2)
		panic("abc")
		println() // ERROR "unreachable code"
	case c <- 1:
		print(2)
		goto L
		println() // ERROR "unreachable code"
	}
}

func _() int {
L:
	print(1)
	select {
	case <-c:
		print(2)
		panic("abc")
	case c <- 1:
		print(2)
		goto L
	}
	println() // ERROR "unreachable code"
}

func _() int {
	print(1)
	select {
	case <-c:
		print(2)
		panic("abc")
		println() // ERROR "unreachable code"
	default:
		select {}
		println() // ERROR "unreachable code"
	}
}

func _() int {
	print(1)
	select {
	case <-c:
		print(2)
		panic("abc")
	default:
		select {}
	}
	println() // ERROR "unreachable code"
}

func _() int {
	print(1)
	select {
	case <-c:
		print(2)
	}
	println() // ok
}

func _() int {
L:
	print(1)
	select {
	case <-c:
		print(2)
		panic("abc")
		goto L // ERROR "unreachable code"
	case c <- 1:
		print(2)
	}
	println() // ok
}

func _() int {
	print(1)
	select {
	case <-c:
		print(2)
		panic("abc")
	default:
		print(2)
	}
	println() // ok
}

func _() int {
	print(1)
	select {
	default:
		break
	}
	println() // ok
}

func _() int {
	print(1)
	select {
	case <-c:
		print(2)
		panic("abc")
		break // ERROR "unreachable code"
	}
	println() // ok
}

func _() int {
	print(1)
L:
	select {
	case <-c:
		print(2)
		for {
			break L
		}
	}
	println() // ok
}

func _() int {
	print(1)
L:
	select {
	case <-c:
		print(2)
		panic("abc")
	case c <- 1:
		print(2)
		break L
	}
	println() // ok
}

func _() int {
	print(1)
	select {
	case <-c:
		print(1)
		panic("abc")
	default:
		select {}
		break // ERROR "unreachable code"
	}
	println() // ok
}

func _() int {
	print(1)
	switch x {
	case 1:
		print(2)
		panic(3)
		println() // ERROR "unreachable code"
	default:
		return 4
		println() // ERROR "unreachable code"
	}
}

func _() int {
	print(1)
	switch x {
	case 1:
		print(2)
		panic(3)
	default:
		return 4
	}
	println() // ERROR "unreachable code"
}

func _() int {
	print(1)
	switch x {
	default:
		return 4
		println() // ERROR "unreachable code"
	case 1:
		print(2)
		panic(3)
		println() // ERROR "unreachable code"
	}
}

func _() int {
	print(1)
	switch x {
	default:
		return 4
	case 1:
		print(2)
		panic(3)
	}
	println() // ERROR "unreachable code"
}

func _() int {
	print(1)
	switch x {
	case 1:
		print(2)
		fallthrough
	default:
		return 4
		println() // ERROR "unreachable code"
	}
}

func _() int {
	print(1)
	switch x {
	case 1:
		print(2)
		fallthrough
	default:
		return 4
	}
	println() // ERROR "unreachable code"
}

func _() int {
	print(1)
	switch {
	}
	println() // ok
}

func _() int {
	print(1)
	switch x {
	case 1:
		print(2)
		panic(3)
	case 2:
		return 4
	}
	println() // ok
}

func _() int {
	print(1)
	switch x {
	case 2:
		return 4
	case 1:
		print(2)
		panic(3)
	}
	println() // ok
}

func _() int {
	print(1)
	switch x {
	case 1:
		print(2)
		fallthrough
	case 2:
		return 4
	}
	println() // ok
}

func _() int {
	print(1)
	switch x {
	case 1:
		print(2)
		panic(3)
	}
	println() // ok
}

func _() int {
	print(1)
L:
	switch x {
	case 1:
		print(2)
		panic(3)
		break L // ERROR "unreachable code"
	default:
		return 4
	}
	println() // ok
}

func _() int {
	print(1)
	switch x {
	default:
		return 4
		break // ERROR "unreachable code"
	case 1:
		print(2)
		panic(3)
	}
	println() // ok
}

func _() int {
	print(1)
L:
	switch x {
	case 1:
		print(2)
		for {
			break L
		}
	default:
		return 4
	}
	println() // ok
}

func _() int {
	print(1)
	switch x.(type) {
	case int:
		print(2)
		panic(3)
		println() // ERROR "unreachable code"
	default:
		return 4
		println() // ERROR "unreachable code"
	}
}

func _() int {
	print(1)
	switch x.(type) {
	case int:
		print(2)
		panic(3)
	default:
		return 4
	}
	println() // ERROR "unreachable code"
}

func _() int {
	print(1)
	switch x.(type) {
	default:
		return 4
		println() // ERROR "unreachable code"
	case int:
		print(2)
		panic(3)
		println() // ERROR "unreachable code"
	}
}

func _() int {
	print(1)
	switch x.(type) {
	default:
		return 4
	case int:
		print(2)
		panic(3)
	}
	println() // ERROR "unreachable code"
}

func _() int {
	print(1)
	switch x.(type) {
	case int:
		print(2)
		fallthrough
	default:
		return 4
		println() // ERROR "unreachable code"
	}
}

func _() int {
	print(1)
	switch x.(type) {
	case int:
		print(2)
		fallthrough
	default:
		return 4
	}
	println() // ERROR "unreachable code"
}

func _() int {
	print(1)
	switch {
	}
	println() // ok
}

func _() int {
	print(1)
	switch x.(type) {
	case int:
		print(2)
		panic(3)
	case float64:
		return 4
	}
	println() // ok
}

func _() int {
	print(1)
	switch x.(type) {
	case float64:
		return 4
	case int:
		print(2)
		panic(3)
	}
	println() // ok
}

func _() int {
	print(1)
	switch x.(type) {
	case int:
		print(2)
		fallthrough
	case float64:
		return 4
	}
	println() // ok
}

func _() int {
	print(1)
	switch x.(type) {
	case int:
		print(2)
		panic(3)
	}
	println() // ok
}

func _() int {
	print(1)
L:
	switch x.(type) {
	case int:
		print(2)
		panic(3)
		break L // ERROR "unreachable code"
	default:
		return 4
	}
	println() // ok
}

func _() int {
	print(1)
	switch x.(type) {
	default:
		return 4
		break // ERROR "unreachable code"
	case int:
		print(2)
		panic(3)
	}
	println() // ok
}

func _() int {
	print(1)
L:
	switch x.(type) {
	case int:
		print(2)
		for {
			break L
		}
	default:
		return 4
	}
	println() // ok
}

// again, but without the leading print(1).
// testing that everything works when the terminating statement is first.

func _() int {
	println() // ok
}

func _() int {
	return 2
	println() // ERROR "unreachable code"
}

func _() int {
L:
	goto L
	println() // ERROR "unreachable code"
}

func _() int {
	panic(2)
	println() // ERROR "unreachable code"
}

// but only builtin panic
func _() int {
	var panic = func(int) {}
	panic(2)
	println() // ok
}

func _() int {
	{
		return 2
		println() // ERROR "unreachable code"
	}
}

func _() int {
	{
		return 2
	}
	println() // ERROR "unreachable code"
}

func _() int {
L:
	{
		goto L
		println() // ERROR "unreachable code"
	}
}

func _() int {
L:
	{
		goto L
	}
	println() // ERROR "unreachable code"
}

func _() int {
	{
		panic(2)
		println() // ERROR "unreachable code"
	}
}

func _() int {
	{
		panic(2)
	}
	println() // ERROR "unreachable code"
}

func _() int {
	return 2
	{ // ERROR "unreachable code"
	}
	println() // ok
}

func _() int {
L:
	goto L
	{ // ERROR "unreachable code"
	}
	println() // ok
}

func _() int {
	panic(2)
	{ // ERROR "unreachable code"
	}
	println() // ok
}

func _() int {
	{
		return 2
		{ // ERROR "unreachable code"
		}
	}
	println() // ok
}

func _() int {
L:
	{
		goto L
		{ // ERROR "unreachable code"
		}
	}
	println() // ok
}

func _() int {
	{
		panic(2)
		{ // ERROR "unreachable code"
		}
	}
	println() // ok
}

func _() int {
	{
		return 2
	}
	{ // ERROR "unreachable code"
	}
	println() // ok
}

func _() int {
L:
	{
		goto L
	}
	{ // ERROR "unreachable code"
	}
	println() // ok
}

func _() int {
	{
		panic(2)
	}
	{ // ERROR "unreachable code"
	}
	println() // ok
}

// again, with func literals

var _ = func() int {
}

var _ = func() int {
	print(1)
}

var _ = func() int {
	print(1)
	return 2
	println() // ERROR "unreachable code"
}

var _ = func() int {
L:
	print(1)
	goto L
	println() // ERROR "unreachable code"
}

var _ = func() int {
	print(1)
	panic(2)
	println() // ERROR "unreachable code"
}

// but only builtin panic
var _ = func() int {
	var panic = func(int) {}
	print(1)
	panic(2)
	println() // ok
}

var _ = func() int {
	{
		print(1)
		return 2
		println() // ERROR "unreachable code"
	}
	println() // ok
}

var _ = func() int {
	{
		print(1)
		return 2
	}
	println() // ERROR "unreachable code"
}

var _ = func() int {
L:
	{
		print(1)
		goto L
		println() // ERROR "unreachable code"
	}
	println() // ok
}

var _ = func() int {
L:
	{
		print(1)
		goto L
	}
	println() // ERROR "unreachable code"
}

var _ = func() int {
	print(1)
	{
		panic(2)
	}
}

var _ = func() int {
	print(1)
	{
		panic(2)
		println() // ERROR "unreachable code"
	}
}

var _ = func() int {
	print(1)
	{
		panic(2)
	}
	println() // ERROR "unreachable code"
}

var _ = func() int {
	print(1)
	return 2
	{ // ERROR "unreachable code"
	}
}

var _ = func() int {
L:
	print(1)
	goto L
	{ // ERROR "unreachable code"
	}
}

var _ = func() int {
	print(1)
	panic(2)
	{ // ERROR "unreachable code"
	}
}

var _ = func() int {
	{
		print(1)
		return 2
		{ // ERROR "unreachable code"
		}
	}
}

var _ = func() int {
L:
	{
		print(1)
		goto L
		{ // ERROR "unreachable code"
		}
	}
}

var _ = func() int {
	print(1)
	{
		panic(2)
		{ // ERROR "unreachable code"
		}
	}
}

var _ = func() int {
	{
		print(1)
		return 2
	}
	{ // ERROR "unreachable code"
	}
}

var _ = func() int {
L:
	{
		print(1)
		goto L
	}
	{ // ERROR "unreachable code"
	}
}

var _ = func() int {
	print(1)
	{
		panic(2)
	}
	{ // ERROR "unreachable code"
	}
}

var _ = func() int {
	print(1)
	if x == nil {
		panic(2)
	} else {
		panic(3)
	}
	println() // ERROR "unreachable code"
}

var _ = func() int {
L:
	print(1)
	if x == nil {
		panic(2)
	} else {
		goto L
	}
	println() // ERROR "unreachable code"
}

var _ = func() int {
L:
	print(1)
	if x == nil {
		panic(2)
	} else if x == 1 {
		return 0
	} else if x != 2 {
		panic(3)
	} else {
		goto L
	}
	println() // ERROR "unreachable code"
}

// if-else chain missing final else is not okay, even if the
// conditions cover every possible case.

var _ = func() int {
	print(1)
	if x == nil {
		panic(2)
	} else if x != nil {
		panic(3)
	}
	println() // ok
}

var _ = func() int {
	print(1)
	if x == nil {
		panic(2)
	}
	println() // ok
}

var _ = func() int {
L:
	print(1)
	if x == nil {
		panic(2)
	} else if x == 1 {
		return 0
	} else if x != 1 {
		panic(3)
	}
	println() // ok
}

var _ = func() int {
	print(1)
	for {
	}
	println() // ERROR "unreachable code"
}

var _ = func() int {
	for {
		for {
			break
		}
	}
	println() // ERROR "unreachable code"
}

var _ = func() int {
	for {
		for {
			break
			println() // ERROR "unreachable code"
		}
	}
}

var _ = func() int {
	for {
		for {
			continue
			println() // ERROR "unreachable code"
		}
	}
}

var _ = func() int {
	for {
	L:
		for {
			break L
		}
	}
	println() // ERROR "unreachable code"
}

var _ = func() int {
	print(1)
	for {
		break
	}
	println() // ok
}

var _ = func() int {
	for {
		for {
		}
		break // ERROR "unreachable code"
	}
	println() // ok
}

var _ = func() int {
L:
	for {
		for {
			break L
		}
	}
	println() // ok
}

var _ = func() int {
	print(1)
	for x == nil {
	}
	println() // ok
}

var _ = func() int {
	for x == nil {
		for {
			break
		}
	}
	println() // ok
}

var _ = func() int {
	for x == nil {
	L:
		for {
			break L
		}
	}
	println() // ok
}

var _ = func() int {
	print(1)
	for true {
	}
	println() // ok
}

var _ = func() int {
	for true {
		for {
			break
		}
	}
	println() // ok
}

var _ = func() int {
	for true {
	L:
		for {
			break L
		}
	}
	println() // ok
}

var _ = func() int {
	print(1)
	select {}
	println() // ERROR "unreachable code"
}

var _ = func() int {
	print(1)
	select {
	case <-c:
		print(2)
		panic("abc")
		println() // ERROR "unreachable code"
	}
}

var _ = func() int {
	print(1)
	select {
	case <-c:
		print(2)
		panic("abc")
	}
	println() // ERROR "unreachable code"
}

var _ = func() int {
	print(1)
	select {
	case <-c:
		print(2)
		for {
		}
		println() // ERROR "unreachable code"
	}
}

var _ = func() int {
	print(1)
	select {
	case <-c:
		print(2)
		for {
		}
	}
	println() // ERROR "unreachable code"
}

var _ = func() int {
L:
	print(1)
	select {
	case <-c:
		print(2)
		panic("abc")
		println() // ERROR "unreachable code"
	case c <- 1:
		print(2)
		goto L
		println() // ERROR "unreachable code"
	}
}

var _ = func() int {
L:
	print(1)
	select {
	case <-c:
		print(2)
		panic("abc")
	case c <- 1:
		print(2)
		goto L
	}
	println() // ERROR "unreachable code"
}

var _ = func() int {
	print(1)
	select {
	case <-c:
		print(2)
		panic("abc")
		println() // ERROR "unreachable code"
	default:
		select {}
		println() // ERROR "unreachable code"
	}
}

var _ = func() int {
	print(1)
	select {
	case <-c:
		print(2)
		panic("abc")
	default:
		select {}
	}
	println() // ERROR "unreachable code"
}

var _ = func() int {
	print(1)
	select {
	case <-c:
		print(2)
	}
	println() // ok
}

var _ = func() int {
L:
	print(1)
	select {
	case <-c:
		print(2)
		panic("abc")
		goto L // ERROR "unreachable code"
	case c <- 1:
		print(2)
	}
	println() // ok
}

var _ = func() int {
	print(1)
	select {
	case <-c:
		print(2)
		panic("abc")
	default:
		print(2)
	}
	println() // ok
}

var _ = func() int {
	print(1)
	select {
	default:
		break
	}
	println() // ok
}

var _ = func() int {
	print(1)
	select {
	case <-c:
		print(2)
		panic("abc")
		break // ERROR "unreachable code"
	}
	println() // ok
}

var _ = func() int {
	print(1)
L:
	select {
	case <-c:
		print(2)
		for {
			break L
		}
	}
	println() // ok
}

var _ = func() int {
	print(1)
L:
	select {
	case <-c:
		print(2)
		panic("abc")
	case c <- 1:
		print(2)
		break L
	}
	println() // ok
}

var _ = func() int {
	print(1)
	select {
	case <-c:
		print(1)
		panic("abc")
	default:
		select {}
		break // ERROR "unreachable code"
	}
	println() // ok
}

var _ = func() int {
	print(1)
	switch x {
	case 1:
		print(2)
		panic(3)
		println() // ERROR "unreachable code"
	default:
		return 4
		println() // ERROR "unreachable code"
	}
}

var _ = func() int {
	print(1)
	switch x {
	case 1:
		print(2)
		panic(3)
	default:
		return 4
	}
	println() // ERROR "unreachable code"
}

var _ = func() int {
	print(1)
	switch x {
	default:
		return 4
		println() // ERROR "unreachable code"
	case 1:
		print(2)
		panic(3)
		println() // ERROR "unreachable code"
	}
}

var _ = func() int {
	print(1)
	switch x {
	default:
		return 4
	case 1:
		print(2)
		panic(3)
	}
	println() // ERROR "unreachable code"
}

var _ = func() int {
	print(1)
	switch x {
	case 1:
		print(2)
		fallthrough
	default:
		return 4
		println() // ERROR "unreachable code"
	}
}

var _ = func() int {
	print(1)
	switch x {
	case 1:
		print(2)
		fallthrough
	default:
		return 4
	}
	println() // ERROR "unreachable code"
}

var _ = func() int {
	print(1)
	switch {
	}
	println() // ok
}

var _ = func() int {
	print(1)
	switch x {
	case 1:
		print(2)
		panic(3)
	case 2:
		return 4
	}
	println() // ok
}

var _ = func() int {
	print(1)
	switch x {
	case 2:
		return 4
	case 1:
		print(2)
		panic(3)
	}
	println() // ok
}

var _ = func() int {
	print(1)
	switch x {
	case 1:
		print(2)
		fallthrough
	case 2:
		return 4
	}
	println() // ok
}

var _ = func() int {
	print(1)
	switch x {
	case 1:
		print(2)
		panic(3)
	}
	println() // ok
}

var _ = func() int {
	print(1)
L:
	switch x {
	case 1:
		print(2)
		panic(3)
		break L // ERROR "unreachable code"
	default:
		return 4
	}
	println() // ok
}

var _ = func() int {
	print(1)
	switch x {
	default:
		return 4
		break // ERROR "unreachable code"
	case 1:
		print(2)
		panic(3)
	}
	println() // ok
}

var _ = func() int {
	print(1)
L:
	switch x {
	case 1:
		print(2)
		for {
			break L
		}
	default:
		return 4
	}
	println() // ok
}

var _ = func() int {
	print(1)
	switch x.(type) {
	case int:
		print(2)
		panic(3)
		println() // ERROR "unreachable code"
	default:
		return 4
		println() // ERROR "unreachable code"
	}
}

var _ = func() int {
	print(1)
	switch x.(type) {
	case int:
		print(2)
		panic(3)
	default:
		return 4
	}
	println() // ERROR "unreachable code"
}

var _ = func() int {
	print(1)
	switch x.(type) {
	default:
		return 4
		println() // ERROR "unreachable code"
	case int:
		print(2)
		panic(3)
		println() // ERROR "unreachable code"
	}
}

var _ = func() int {
	print(1)
	switch x.(type) {
	default:
		return 4
	case int:
		print(2)
		panic(3)
	}
	println() // ERROR "unreachable code"
}

var _ = func() int {
	print(1)
	switch x.(type) {
	case int:
		print(2)
		fallthrough
	default:
		return 4
		println() // ERROR "unreachable code"
	}
}

var _ = func() int {
	print(1)
	switch x.(type) {
	case int:
		print(2)
		fallthrough
	default:
		return 4
	}
	println() // ERROR "unreachable code"
}

var _ = func() int {
	print(1)
	switch {
	}
	println() // ok
}

var _ = func() int {
	print(1)
	switch x.(type) {
	case int:
		print(2)
		panic(3)
	case float64:
		return 4
	}
	println() // ok
}

var _ = func() int {
	print(1)
	switch x.(type) {
	case float64:
		return 4
	case int:
		print(2)
		panic(3)
	}
	println() // ok
}

var _ = func() int {
	print(1)
	switch x.(type) {
	case int:
		print(2)
		fallthrough
	case float64:
		return 4
	}
	println() // ok
}

var _ = func() int {
	print(1)
	switch x.(type) {
	case int:
		print(2)
		panic(3)
	}
	println() // ok
}

var _ = func() int {
	print(1)
L:
	switch x.(type) {
	case int:
		print(2)
		panic(3)
		break L // ERROR "unreachable code"
	default:
		return 4
	}
	println() // ok
}

var _ = func() int {
	print(1)
	switch x.(type) {
	default:
		return 4
		break // ERROR "unreachable code"
	case int:
		print(2)
		panic(3)
	}
	println() // ok
}

var _ = func() int {
	print(1)
L:
	switch x.(type) {
	case int:
		print(2)
		for {
			break L
		}
	default:
		return 4
	}
	println() // ok
}

// again, but without the leading print(1).
// testing that everything works when the terminating statement is first.

var _ = func() int {
	println() // ok
}

var _ = func() int {
	return 2
	println() // ERROR "unreachable code"
}

var _ = func() int {
L:
	goto L
	println() // ERROR "unreachable code"
}

var _ = func() int {
	panic(2)
	println() // ERROR "unreachable code"
}

// but only builtin panic
var _ = func() int {
	var panic = func(int) {}
	panic(2)
	println() // ok
}

var _ = func() int {
	{
		return 2
		println() // ERROR "unreachable code"
	}
}

var _ = func() int {
	{
		return 2
	}
	println() // ERROR "unreachable code"
}

var _ = func() int {
L:
	{
		goto L
		println() // ERROR "unreachable code"
	}
}

var _ = func() int {
L:
	{
		goto L
	}
	println() // ERROR "unreachable code"
}

var _ = func() int {
	{
		panic(2)
		println() // ERROR "unreachable code"
	}
}

var _ = func() int {
	{
		panic(2)
	}
	println() // ERROR "unreachable code"
}

var _ = func() int {
	return 2
	{ // ERROR "unreachable code"
	}
	println() // ok
}

var _ = func() int {
L:
	goto L
	{ // ERROR "unreachable code"
	}
	println() // ok
}

var _ = func() int {
	panic(2)
	{ // ERROR "unreachable code"
	}
	println() // ok
}

var _ = func() int {
	{
		return 2
		{ // ERROR "unreachable code"
		}
	}
	println() // ok
}

var _ = func() int {
L:
	{
		goto L
		{ // ERROR "unreachable code"
		}
	}
	println() // ok
}

var _ = func() int {
	{
		panic(2)
		{ // ERROR "unreachable code"
		}
	}
	println() // ok
}

var _ = func() int {
	{
		return 2
	}
	{ // ERROR "unreachable code"
	}
	println() // ok
}

var _ = func() int {
L:
	{
		goto L
	}
	{ // ERROR "unreachable code"
	}
	println() // ok
}

var _ = func() int {
	{
		panic(2)
	}
	{ // ERROR "unreachable code"
	}
	println() // ok
}