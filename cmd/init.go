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

	"github.com/gnames/bhlquest/internal/io/bhlnio"
	"github.com/gnames/bhlquest/internal/io/embedio"
	"github.com/gnames/bhlquest/internal/io/llmutilio"
	"github.com/gnames/bhlquest/internal/io/storageio"
	bhlquest "github.com/gnames/bhlquest/pkg"
	"github.com/gnames/bhlquest/pkg/config"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Makes database, embeds BHL items and saves them.",
	Long:  `Makes database, embeds BHL items and saves them.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New(opts...)

		bn, err := bhlnio.New(cfg)
		if err != nil {
			msg := fmt.Sprintf("Init of BHLNames failed: %s", err)
			slog.Error(msg)
			os.Exit(1)
		}

		stg := storageio.New(cfg)

		llm, err := llmutilio.New(cfg)
		if err != nil {
			msg := fmt.Sprintf("Init of BHLNames failed: %s", err)
			slog.Error(msg)
			os.Exit(1)
		}
		emb, err := embedio.New(cfg, stg, llm)
		if err != nil {
			msg := fmt.Sprintf("Init of BHLNames failed: %s", err)
			slog.Error(msg)
			os.Exit(1)
		}

		cp := bhlquest.Components{
			BHLNames: bn,
			Embed:    emb,
		}

		bq := bhlquest.New(cfg, cp)
		err = bq.Init()
		if err != nil {
			msg := fmt.Sprintf("Init failed: %s", err)
			slog.Error(msg)
			os.Exit(1)
		}

		slog.Info("Init finished")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
