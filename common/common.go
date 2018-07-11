// Copyright © 2018 Patrick Nuckolls <nuckollsp at gmail>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package common

import (
	"os"
	"log"
	"fmt"
)

//Operating system name type
//Should be limited to those that Go can target
type OS string
//Operating system system architecture
//Should be limited to those that Go can target
type Arch string

//Basic system information
type SystemType struct {
	OS OS
	Arch Arch
}

// Any additional architectures/OS's need to be added here
const (
	Linux OS = "linux"
	Windows OS = "windows"
	OSx OS = "darwin"
	Amd64 Arch = "amd64"
	I386 Arch = "386"
)

var (
	forceSuperuser bool
	superuserForced bool
)

func (system OS) ToString() string {
	return string(system)
}
func (system Arch) ToString() string {
	return string(system)
}

func (system SystemType) ToString() string {
	return system.OS.ToString() + "-" + system.Arch.ToString()
}

func IsInteractive() bool {
	return true
}

func ForceSuperuser(value bool) {
	if superuserForced {
		PrintError("Superuser forced twice")
	}
	forceSuperuser = value
	superuserForced = true
}

func PrintError(in ...interface{}) {
	if IsSuperUser() {
		log.Fatal(in...)
	} else {
		fmt.Println(in...)
		os.Exit(1)
	}
}

func Println(in ...interface{}) {
	if IsSuperUser() {
		log.Println(in...)
	} else {
		fmt.Println(in...)
	}
}

// Print only if appropriate for verbosity
func PrintVerbose(in ...interface{}) {
	if ! IsSuperUser() {
		fmt.Println(in...)
	}
}
