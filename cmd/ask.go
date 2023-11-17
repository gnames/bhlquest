/*
Copyright © 2023 Dmitry Mozzherin <dmozzherin@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/gnames/gnfmt"
	"github.com/spf13/cobra"
)

// askCmd represents the ask command
var askCmd = &cobra.Command{
	Use:   "ask",
	Short: "Allows to ask questions about BHL content.",
	Long:  `Allows to ask questions about BHL content.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			slog.Error("No question is given.")
			os.Exit(1)
		}

		flags := []flagFunc{scoreThresholdFlag, maxResultsFlag}
		for _, v := range flags {
			v(cmd)
		}

		bq := bhlquestFactory()
		q := args[0]
		answer, err := bq.Ask(q)
		if err != nil {
			slog.Error("No reply", "error", err)
			os.Exit(1)
		}
		enc := gnfmt.GNjson{}
		out, _ := enc.Encode(answer)
		fmt.Println(string(out))
	},
}

func init() {
	rootCmd.AddCommand(askCmd)

	askCmd.Flags().IntP(
		"max-results", "m", 0,
		"defines the maximum number of results",
	)
	askCmd.Flags().Float64P(
		"score-threshold", "s", 0,
		"results with lower score will be ignored",
	)
}
