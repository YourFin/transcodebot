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
	"path/filepath"
	"fmt"

	"github.com/yourfin/transcodebot/common"
)

type BuildSettings struct {
	OutputLocation string
	OutputPrefix string
	NoCompress bool
}

func Build(settings BuildSettings) error {
	common.Println(settings.OutputLocation)
	common.Println(settings.OutputPrefix)
	common.Println(settings.NoCompress)

	//get the dir we were called from so we can come back
	calledPath, err := os.Getwd()
	if err != nil {
		common.PrintError("getting working directory err:", err)
	}
	calledPath, err = filepath.Abs(calledPath)
	if err != nil {
		common.PrintError("absolute path err:", err)
	}

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

	return nil
}
