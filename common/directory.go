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
	"fmt"
	"path/filepath"
	"io/ioutil"
	"errors"

	"github.com/songmu/prompter"
)

var (
	settingsDir string
	settingsDirSet bool = false
)

//Sets internal settings directory. Will panic if called more than once.
func SetSettingsDir(settingsDirIn string) {
	if settingsDirSet {
		panic("common: attempted to set settings dir when alread set")
	}
	settingsDirSet = true
	//Need to do this to avoid shadowing settingsDir with
	// :=
	var err error
	settingsDir, err = filepath.Abs(settingsDirIn)
	if err != nil {
		PrintError("Set settings dir err:", err)
	}
}

//Safely get the settings directory
func SettingsDir() string {
	if ! settingsDirSet {
		panic("Attempted to get settings dir when not set")
	}
	return settingsDir
}

//Safely write data bytes to a file inside the settings directory.
func SettingsWriteFile(data []byte, relPath ...string) error {
	fullPath := SettingsDir() + string(filepath.Separator) + filepath.Join(relPath...)
	err := CowardlyCreateDir(parentDir)
	if err != nil {
		return err
	}
	parentInfo, err := os.Stat(parentDir)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fullPath, data, parentInfo.Mode())
}

//Will ask for user confirmation for creating folders outside SettingsDir, or to create SettingsDir.
//Otherwise will error out
func CowardlyCreateDir(dirname string) error {
	//TODO: clean this mess up
	existing, nonexistent, info, err := findExistingParentDir(dirname)

	//I'm not happy with this, but every solution I come up with
	//seems grosser
	prompt := func(toPrompt string, dirToMake string) error {
		if prompter.YN(
			fmt.Sprintf(
				toPrompt,
				dirToMake,
				existing,
			),
			false, // Default to no for no tty
		) {
			return nil
		} else {
			return errors.New("Cowardly refusing to create output directory: " + dirname)
		}
	}

	var promptErr error
	if err != nil {
		return err
	} else if nonexistent != "" {
		if pathStartsWith(dirname, SettingsDir()) {
			if !pathStartsWith(existing, SettingsDir()) {
				promptErr = prompt("Settings directory %s does not exist, but %s does.\nCreate intermediate folders?", SettingsDir())
			} else {
				promptErr = nil
			}
		} else {
			err = prompt("%s does not exist, but %s does.\nCreate intermediate folders?", dirname)
		}
	} else {
		return err
	}
	if promptErr != nil {
		return errors.New(fmt.Sprintf("create dir: '%s'; parent did not exist", dirname))
	} else {
		return os.MkdirAll(dirname, info.Mode())
	}
}

//Returns true if the first part of checkee is targetStart. Just returns false if it can't parse the data, which probably needs to be replaced with error handling
func pathStartsWith(checkee, targetStart string) bool {
	//Sanitize input
	checkee, err := filepath.Abs(checkee)
	if err != nil {
		return false
	}

	targetStart, err = filepath.Abs(targetStart)
	if err != nil {
		return false
	}

	targetLen := len(targetStart)
	return targetLen <= len(checkee) && targetStart == checkee[:targetLen]
}

func settingsDirExists() bool {
	existing, _, _, err := findExistingParentDir(SettingsDir())
	return err != nil && existing == SettingsDir()
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
			Println("fileinfo:", fileinfo)
			Println("err:", err)
			Println("dirname:", dirname)
			panic("Reached a place that was thought to be unreachable. Contact the maintainer of transcodebot with the above three lines")
		}
	}
}
