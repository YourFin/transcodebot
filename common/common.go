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
func (system OS) ToString() string {
	return string(system)
}
func (system Arch) ToString() string {
	return string(system)
}

func (system SystemType) ToString() string {
	return system.OS.ToString() + "-" + system.Arch.ToString()
}

// Settings to pass to ffmpeg to use for transcoding
// Most of these settings line up with something in the ffmeg documentation, and
type TranscodeSettings struct {
	//Optionally short-circuit all of the built in logic and specify the command manually, and ignore all other options.
	//Pseudocode of what happens if this option is not zero-length:
	//for line in RawffmpegOptions:
	//	Sys.Execute("ffmpeg -nostdin -i $inputfile $line")
	RawffmpegOptions []string

	//If false, error if a stream is found that cannot be converted
	HandleUnuseableStreams bool
	//If true and HandleUnuseableStreams is true, throw out streams that can't be converted to a target type
	//If false, pass through the stream to the output unmolested
	TossUnuseableStreams bool
	//The containing file type, i.e. webm, mkv, mp4, mov
	ContainerType string

	//Video options:

	//Video Codec to use
	VideoCodec string
	//Subpixel format
	PixFormat string
	//Speed parameter for the primary (second in two pass) encoding pass
	PrimaryPassSpeed uint8
	//Whether to do two-pass encoding
	TwoPass bool
	//Speed parameter for the log (first and possibly omitted) pass
	PreliminaryPassSpeed uint8

	//Audio codec to use
	AudioCodec string

	//Subtitle codec to use
	SubtitleCodec string
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
	SettingsDir string
)


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
