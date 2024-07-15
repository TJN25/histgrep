/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/TJN25/histgrep/utils"
	"github.com/TJN25/histgrep/hsdata"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "View the formats and defaults.",
	Long: `View the formats and defaults.`,
	Run: infoRun,
}

func init() {
	rootCmd.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	infoCmd.PersistentFlags().CountP("verbose", "v", "Level of verbosity (0-5) default (0)")
}

func infoRun(cmd *cobra.Command, args []string) {

	verbosity, _ := cmd.PersistentFlags().GetCount("verbose")
	switch verbosity {
		case 0:
			log.SetLevel(log.FatalLevel)
		case 1:
			log.SetLevel(log.ErrorLevel)
		case 2:
			log.SetLevel(log.WarnLevel)
		case 3:
			log.SetLevel(log.InfoLevel)
		case 4:
			log.SetLevel(log.DebugLevel)
		case 5:
			log.SetLevel(log.TraceLevel)
		default:
			log.SetLevel(log.TraceLevel)
	}
	file, err := utils.GetDataPath("formats.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Info(fmt.Sprintf("Using config file %v", file))
	configMap := hsdata.ConfigMap{}
	utils.FetchFormatting(file, &configMap)
	for k, v := range configMap {
		fmt.Printf("Name: %v \n    input: \"%v\" \n    output: \"%v\"\n",k, v.Input, v.Output)
	}
	}
