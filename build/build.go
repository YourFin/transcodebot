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
	"fmt"
	"net"
	"time"
	//"crypto/x509"

	cert "github.com/yourfin/transcodebot/certificate"
	"github.com/yourfin/transcodebot/common"
)

//Settings for building the clients
type BuildSettings struct {
	//Prefix for build output files.
	//Will be followed by target arch and a file extension if applicable
	//Default is transcode-client
	OutputPrefix string

	//Whether or not to compress the files
	//if this variable is true, then the output binaries will not be zipped
	//Default false
	NoCompress bool

	//Force a new server certificate to be generated
	//Invalidates all previous clients
	ForceNewCert bool

	//Valid IP's for the main server
	ServerIPs []net.IP

	//List of system os/arch combinations to target
	Targets []common.SystemType
}
const build_extention = "clients"

//Builds client binaries according to the passed in settings
func Build(settings BuildSettings) error {
	buildDir := common.SettingsDir(build_extention)

	if settings.ForceNewCert { //or no cert exists
		cert.GenRootCert(settings.ServerIPs)
	}
	rootCert := cert.ReadCert("root")
	rootKey := cert.ReadRsaKey("root")

	//get the dir we were called from so we can come back
	calledPath, err := os.Getwd()
	if err != nil {
		common.PrintError("getting working directory err: ", err)
	}
	calledPath, err = filepath.Abs(calledPath)
	if err != nil {
		common.PrintError("absolute path err: ", err)
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
		common.PrintError("Moving to build dir err: ", err, "\nAre you sure your GOPATH environment variable is set?")
	}

	common.CowardlyCreateDir(buildDir)

	//Compile
	common.Println("Building...")
	doneChan := make(chan int)
	for ii, target := range settings.Targets {
		//Generate new client certificate
		privKey, cert := cert.GenClientCert(target.ToString() + time.Now().String(), rootCert, rootKey)

		builtName := filepath.Join(buildDir, settings.OutputPrefix + target.ToString())
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
		//Note that range variables are shared between
		//loops but others are not, hence the passing by
		//value
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
	common.PrintVerbose("All complies finished. Binaries at:", buildDir)
	return nil
}
