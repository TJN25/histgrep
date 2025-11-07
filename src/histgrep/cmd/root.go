/*
Copyright © 2024 dani
*/
package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/TJN25/histgrep/hsdata"
	"github.com/TJN25/histgrep/utils"
	"github.com/spf13/cobra"
)

const VERSION = "0.3.0"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "histgrep",
	Short: "HistGrep - Enhanced Command History Search",
	Long:  `HistGrep is a terminal-based command-line tool for searching through history files or other logs.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: rootRun,
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.histgrep.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("version", "V", false, "Print version information")
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		versionFlag, _ := cmd.Flags().GetBool("version")
		if versionFlag {
			fmt.Printf("HistGrep version %s\n", VERSION)
			os.Exit(0)
		}
	}

}

func rootRun(cmd *cobra.Command, args []string) {
	data := hsdata.HsData{Terms: []string{}}
	verbosity, _ := cmd.PersistentFlags().GetCount("verbose")
	utils.SetVerbosity(verbosity)
	config := rootGetArgs(&data)
	log.Info(fmt.Sprintf("\n    Running search with: \n    files: %v -> %v\n    Terms: %v\n    Format: %v\n", data.InputFile, data.OutputFile, data.Terms, data.FormatData))
	log.Debug(fmt.Sprintf("Formatting input: %v", data))
	DoFormatting(&data)
	RunLoopFile(&data, config)
}

func rootGetArgs(data *hsdata.HsData) *utils.Config {
	data.InputFile = "stdin"
	config := DoConfigFile(data)
	data.OutputFile = "stdout"
	data.IncludeNumbers = false
	data.ExcludeTerms = []string{}
	data.UsePager = true
	data.FormatData = UseDefaults(data, config)
	log.Debug(data.FormatData)
	log.Info(fmt.Sprintf("Args data: %v", data))
	return config
}
