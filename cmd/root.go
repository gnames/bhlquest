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
	_ "embed"
	"os"

	"github.com/gnames/bhlquest/pkg/config"
	"github.com/spf13/cobra"
)

//go:embed bhlquest.yaml
var configText string

var (
	cfgFile string
	opts    []config.Option
)

// fConfig purpose is to achieve automatic import of data from the
// configuration file, if it exists.
type fConfig struct {
	BHLDir     string
	LlmUtilURL string
	DbHost     string
	DbUser     string
	DbPass     string
	DbName     string
	PortREST   int
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bhlquest",
	Short: "adds AI capabilities to Biodiversity Heritage Library.",
	Long:  `Adds AI capabilities to Biodiversity Heritage Library.`,
	Run: func(cmd *cobra.Command, _ []string) {
		versionFlag(cmd)
		flags := []flagFunc{debugFlag}
		for _, v := range flags {
			v(cmd)
		}
		_ = cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "changes log to debug level")
	rootCmd.Flags().BoolP("version", "V", false, "show version and build timestamp")
}
