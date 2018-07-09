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

package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"


	"github.com/yourfin/transcodebot/common"
	"github.com/yourfin/transcodebot/build"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "build client binaries",
	Long: `Build client binaries for target platforms`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if len(args) != 0 {
			// TODO: Figure out how to call parent help function here
			common.PrintError("build does not take any arguments")
		}

		if buildSettings.OutputLocation == "" {
			if common.IsSuperUser() {
				if common.Os == "windows" {
					buildSettings.OutputLocation = filepath.Join(os.Getenv("ProgramData"), "transcodebot", "build")
				} else {
					// unix
					buildSettings.OutputLocation = "/var/transcodebot/build"
				}
			} else {
				// Not superuser
				buildSettings.OutputLocation = filepath.Join(SettingsDir, "build", "")
			}
		}

		err = build.Build(buildSettings)
		if err != nil {
			common.PrintError("build err:", err)
		}
	},
}

var buildSettings build.BuildSettings

func init() {
	rootCmd.AddCommand(buildCmd)

	// Configuration flags
	buildCmd.PersistentFlags().StringVar(&buildSettings.OutputLocation, "output-location", "", `The folder to place output binaries
(default $settings/clients)`)
	buildCmd.PersistentFlags().StringVar(&buildSettings.OutputPrefix, "output-prefix", "trancode-client", "The start of the binary names")
	buildCmd.PersistentFlags().BoolVarP(&buildSettings.NoCompress, "no-compress", "Z", false, "Don't zip binaries")
}
