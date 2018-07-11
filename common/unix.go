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

// +build darwin dragonfly freebsd js,wasm linux nacl netbsd openbsd solaris

package common

import (
	"os"
	"fmt"

	homedir "github.com/mitchellh/go-homedir"
)

const BuildType = "unix"

func IsSuperUser() bool {
	if superuserForced {
		return forceSuperuser
	}
	return os.Getuid() == 0
}

//Returns full path to the default settings directory location
// ~/.local/share/transcodebot on unix's if normal user
func GetDefaultSettingsDir() string {
	if IsSuperUser() {
		return "/etc/transcodebot/"
	} else {
		home, err := homedir.Expand("~/.local/share/transcodebot/")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return home
	}
}
