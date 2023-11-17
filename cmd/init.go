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
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Makes database, embeds BHL items and saves them.",
	Long:  `Makes database, embeds BHL items and saves them.`,
	Run: func(cmd *cobra.Command, args []string) {

		flags := []flagFunc{classesFlag, noConfirmFlag, rebuildDbFlag}
		for _, v := range flags {
			v(cmd)
		}

		bq := bhlquestFactory()
		err := bq.Init()
		if err != nil {
			slog.Error("Init failed", "error", err)
			os.Exit(1)
		}

		slog.Info("Init finished")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP(
		"classes", "c", "",
		`limit intake to certain classes, e.g. "-t 'Aves,Mammalia'".`,
	)
	initCmd.Flags().BoolP(
		"rebuild-db", "r", false,
		"deletes database and starts it from scratch",
	)
}
