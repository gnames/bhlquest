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
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gnames/bhlquest/pkg/config"
	"github.com/gnames/gnsys"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:embed bhlquest.yaml
var configText string

var (
	cfgFile string
	opts    []config.Option
)

// fConfig purpose is to achieve automatic import of data from the
// configuration file, if it exists.
type configData struct {
	BHLDir     string
	LlmUtilURL string
	DbHost     string
	DbUser     string
	DbPass     string
	DbBHLQuest string
	DbBHLNames string
	PortREST   int
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bhlquest",
	Short: "adds AI capabilities to Biodiversity Heritage Library.",
	Long:  `Adds AI capabilities to Biodiversity Heritage Library.`,
	Run: func(cmd *cobra.Command, _ []string) {
		versionFlag(cmd)
		flags := []flagFunc{debugFlag, noConfirmFlag}
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
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "changes log to debug level")
	rootCmd.PersistentFlags().BoolP("yes-to-confirmations", "y", false,
		"skip all confirmation dialogs")
	rootCmd.Flags().BoolP("version", "V", false, "show version and build timestamp")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var configDir string
	var err error
	configFile := "bhlquest"
	home, err := os.UserHomeDir()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	configDir = filepath.Join(home, ".config")

	// Search config in home directory with name ".gnmatcher" (without extension).
	viper.AddConfigPath(configDir)
	viper.SetConfigName(configFile)

	_ = viper.BindEnv("BHLDir", "BHLQ_BHL_DIR")
	_ = viper.BindEnv("LlmUtilURL", "BHLQ_LLM_UTIL_URL")
	_ = viper.BindEnv("DbHost", "BHLQ_DB_HOST")
	_ = viper.BindEnv("DbUser", "BHLQ_DB_USER")
	_ = viper.BindEnv("DbPass", "BHLQ_DB_PASS")
	_ = viper.BindEnv("DbBHLQuest", "BHLQ_DB_BHL_QUEST")
	_ = viper.BindEnv("DbBHLNames", "BHLQ_DB_BHL_NAMES")
	_ = viper.BindEnv("Port", "BHLQ_PORT")
	_ = viper.BindEnv("InitClasses", "BHLQ_INIT_CLASSES")
	_ = viper.BindEnv("ScoreThreshold", "BHLQ_SCORE_THRESHOLD")
	_ = viper.BindEnv("MaxResultsNum", "BHLQ_MAX_RESULTS_NUM")

	viper.AutomaticEnv() // read in environment variables that match

	configPath := filepath.Join(configDir, fmt.Sprintf("%s.yaml", configFile))
	touchConfigFile(configPath)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		msg := fmt.Sprintf("Using config file %s.", viper.ConfigFileUsed())
		slog.Info(msg)
	}

	getOpts()
}

// getOpts imports data from the configuration file. Some of the settings can
// be overriden by command line flags.
func getOpts() {
	cfgCli := &configData{}
	err := viper.Unmarshal(cfgCli)
	if err != nil {
		msg := fmt.Sprintf("Cannot deserialize config file: %s", err)
		slog.Error(msg)
		os.Exit(1)
	}

	if cfgCli.BHLDir != "" {
		opts = append(opts, config.OptBHLDir(cfgCli.BHLDir))
	}
	if cfgCli.LlmUtilURL != "" {
		opts = append(opts, config.OptLlmUtilURL(cfgCli.LlmUtilURL))
	}
	if cfgCli.DbHost != "" {
		opts = append(opts, config.OptDbHost(cfgCli.DbHost))
	}
	if cfgCli.DbUser != "" {
		opts = append(opts, config.OptDbUser(cfgCli.DbUser))
	}
	if cfgCli.DbPass != "" {
		opts = append(opts, config.OptDbPass(cfgCli.DbPass))
	}
	if cfgCli.DbBHLQuest != "" {
		opts = append(opts, config.OptDbBHLQuest(cfgCli.DbBHLQuest))
	}
	if cfgCli.DbBHLNames != "" {
		opts = append(opts, config.OptDbBHLNames(cfgCli.DbBHLNames))
	}
	if cfgCli.PortREST != 0 {
		opts = append(opts, config.OptPort(cfgCli.PortREST))
	}
}

// touchConfigFile checks if config file exists, and if not, it gets created.
func touchConfigFile(configPath string) error {
	exists, err := gnsys.FileExists(configPath)
	if exists || err != nil {
		return err
	}

	msg := fmt.Sprintf("Creating config file '%s'", configPath)
	slog.Info(msg)
	createConfig(configPath)
	return nil
}

// createConfig creates config file.
func createConfig(path string) {
	err := gnsys.MakeDir(filepath.Dir(path))
	if err != nil {
		msg := fmt.Sprintf("Cannot create dir %s: %s", path, err)
		slog.Error(msg)
		os.Exit(1)
	}

	err = os.WriteFile(path, []byte(configText), 0644)
	if err != nil {
		msg := fmt.Sprintf("Cannot write to file %s: %s", path, err)
		slog.Error(msg)
		os.Exit(1)
	}
}
