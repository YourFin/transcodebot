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
	"github.com/yourfin/transcodebot/build"

	"github.com/spf13/cobra"
	"github.com/yourfin/transcodebot/server/transcode"
	"github.com/yourfin/transcodebot/common"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		appender, err := build.MakeAppender(args[0])
		if err != nil {
			common.PrintError("MakeAppender err: ", err)
		}
		err = appender.AppendFile(args[1])
		if err != nil {
			common.PrintError("append file err: ", err)
		}
		err = appender.Close()
		if err != nil {
			common.PrintError("close appender err: ", err)
		}
	},
}
var watchSettings transcode.WatchSettings

func init() {
	rootCmd.AddCommand(watchCmd)

	watchSettings = transcode.WatchSettings{}

	regexHelp := `regex that input files must match`
	watchCmd.PersistentFlags().StringVarP(&watchSettings.Regex, "regex", "x", `\.(mp4|mov|mpeg|webm|mkv|avi|mts|wmv)$`, regexHelp)

	watchCmd.PersistentFlags().BoolVarP(&watchSettings.Recursive, "recursive", "r", false, "search recursivly for files to transcode")

	//Defined in ./transcode.go
	transcodeServerSettings = addCommonOptions(watchCmd)
}
