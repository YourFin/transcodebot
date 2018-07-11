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

package build

import (
	"os"
	"os/exec"
	"path/filepath"
	"errors"
	"fmt"

	"github.com/songmu/prompter"

	"github.com/yourfin/transcodebot/common"
)

//Settings for building the clients
type BuildSettings struct {
	//The folder where the output will be placed, must be passed in as an absolute path.
	//Default is is $cmd.SettingsDir/build
	OutputLocation string

	//Prefix for build output files.
	//Will be followed by target arch and a file extension if applicable
	//Default is transcode-client
	OutputPrefix string
	//Whether or not to compress the files
	//if this variable is true, then the output binaries will not be zipped
	//Default false
	NoCompress bool

	//List of system os/arch combinations to target
	Targets []common.SystemType
}

//Builds client binaries according to the passed in settings
func Build(settings BuildSettings) error {
	//get the dir we were called from so we can come back
	calledPath, err := os.Getwd()
	if err != nil {
		common.PrintError("getting working directory err:", err)
	}
	calledPath, err = filepath.Abs(calledPath)
	if err != nil {
		common.PrintError("absolute path err:", err)
	}

	//go back to the original working directory after the build
	defer func() {
		err = os.Chdir(calledPath)
		if err != nil {
			panic(fmt.Sprintf("change back to original working dir err: %s", err))
		}
	}()

	err = os.Chdir(filepath.Join(
		os.Getenv("GOPATH"),
		"src",
		"github.com",
		"yourfin",
		"transcodebot",
		"client"))
	if err != nil {
		common.PrintError("Moving to build dir err:", err, "\nAre you sure your GOPATH environment variable is set?")
	}

	// Handle (non)existence of build directory
	existing, nonexistent, info, err := findExistingParentDir(settings.OutputLocation)
	if err != nil {
		return err
	} else if nonexistent != "" {
		if prompter.YN(
			fmt.Sprintf(
				"%s does not exist, but %s does.\nCreate intermediate folders?",
				settings.OutputLocation,
				existing,
			),
			// Default to no for no tty
			false) {
				err = os.MkdirAll(settings.OutputLocation, info.Mode())
				if err != nil {
					common.PrintError("Creating output directory err:", err)
				}
			} else {
				common.PrintError("Cowardly refusing to create output directory: " + settings.OutputLocation)
			}
	}

	//Compile
	doneChan := make(chan int)
	for ii, target := range settings.Targets {
		builtName := filepath.Join(settings.OutputLocation, settings.OutputPrefix + target.ToString())
		if target.OS == common.Windows {
			builtName = builtName + ".exe"
		}
		command := exec.Command("go", "build", "-a", "-o", builtName)
		//Duplicate entries are removed automatically on execution
		command.Env = append(
			os.Environ(),
			"CGO=0",
			"GOARCH=" + target.Arch.ToString(),
			"GOOS=" + target.OS.ToString(),
		)
		//Compile them
		go func(index int, target common.SystemType) {
			//go build doesn't use stdout
			stderr, err := command.CombinedOutput()
			if len(stderr) != 0 {
				common.PrintError("Compile error building", target.ToString(), ":", string(stderr[:]))
			} else if err != nil {
				common.PrintError("Compile error building", target.ToString(), ":", err)
			}
			doneChan <- index
		}(ii, target)
	}
	for finishedCompiles := 0; finishedCompiles < len(settings.Targets); finishedCompiles++ {
		doneNumber := <- doneChan
		common.PrintVerbose(settings.Targets[doneNumber].ToString(), "compile finished")
	}
	return nil
}

//Works upwards on the path until it finds an existing dir
//Will probably break on windows with non-existent drives
func findExistingParentDir(dirname string) (existing string, nonexistent string, fileinfo os.FileInfo, err error) {
	nonexistent = ""
	err = nil
	for {
		fileinfo, err = os.Stat(dirname)
		if err != nil && os.IsNotExist(err) {
			nonexistent = filepath.Join(filepath.Base(dirname), nonexistent)
			dirname = filepath.Dir(dirname)
		} else if err != nil {
			panic("Unhandled error finding parent dir (build.findExistingParentDir):" + err.Error())
		}	else if !fileinfo.IsDir() {
			return dirname, "", nil, errors.New("non-directory file exists as subset of filepath:" + dirname)
		} else if fileinfo.IsDir() {
			existing = dirname
			return
		} else {
			fmt.Println("fileinfo:", fileinfo)
			fmt.Println("err:", err)
			fmt.Println("dirname:", dirname)
			panic("Reached a place that was thought to be unreachable. Contact the maintainer of transcodebot with the above three lines")
		}
	}
}
