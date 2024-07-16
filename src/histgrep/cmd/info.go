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
	infoCmd.Flags().BoolP("names-only", "N", false, "Show only the names of the formats and defaults")
	infoCmd.Flags().BoolP("defaults-only", "D", false, "Show only the defaults")
	infoCmd.Flags().BoolP("formats-only", "F", false, "Show only the defaults")
	infoCmd.Flags().StringP("name", "n", "-", "Print info about a specific format")
	infoCmd.PersistentFlags().CountP("verbose", "v", "Level of verbosity (0-4) default (0)")
}

func infoRun(cmd *cobra.Command, args []string) {

	data := hsdata.InfoData{}
	infoGetArgs(cmd, &data)

	switch data.Verbosity {
		case 0:
			log.SetLevel(log.ErrorLevel)
		case 1:
			log.SetLevel(log.WarnLevel)
		case 2:
			log.SetLevel(log.InfoLevel)
		case 3:
			log.SetLevel(log.DebugLevel)
		case 4:
			log.SetLevel(log.TraceLevel)
		default:
			log.SetLevel(log.TraceLevel)
	}
	log.Info(fmt.Sprintf("Args data: %v", data))

	if !data.Formats_only {
		DoDefaults(&data)
	}
	if !data.Defaults_only {
		DoFormats(&data)
	}
}

func DoFormats(data *hsdata.InfoData) {
	file := utils.GetDataPath("formats.json")
	log.Info(fmt.Sprintf("Using config file %v", file))
	configMap := hsdata.ConfigMap{}
	utils.FetchFormatting(file, &configMap)
	if data.Name != "-" {
		fmt.Println("\n--- Format ---\n ")
		PrintOneFormat(configMap, data.Name)
		return
	} 
	fmt.Println("\n--- Formats ---\n ")
	PrintFormats(configMap, data.Names_only)
}

func PrintFormats(configMap hsdata.ConfigMap, names_only bool) {
	for k, v := range configMap {
		if names_only {
			fmt.Println(k)
		} else {
			fmt.Printf("Name: %v \n    input: \"%v\" \n    output: \"%v\"\n",k, v.Input, v.Output)
		}
	}
}

func PrintOneFormat(configMap hsdata.ConfigMap, name string) {
	v, _ := configMap[name]
	fmt.Printf("Name: %v \n    input: \"%v\" \n    output: \"%v\"\n",name, v.Input, v.Output)
}

func DoDefaults(data *hsdata.InfoData) {
	file := utils.GetDataPath("defaults.json")
	config_file := utils.GetDataPath("formats.json")
	log.Info(fmt.Sprintf("Using defaults file %v", file))
	log.Info(fmt.Sprintf("Using config file %v", config_file))
	fmt.Println("\n--- Defaults ---\n ")
	configMap := hsdata.ConfigMap{}
	defaults := hsdata.DefaultsData{}
	utils.FetchFormatting(config_file, &configMap)
	utils.FetchDefaults(file, &defaults)
    defaultsConfig := configMap.Get(defaults.Name)

	PrintDefaults(defaultsConfig, defaults.Name, data.Names_only)
}

func PrintDefaults(cs hsdata.ConfigSave, name string, names_only bool) {
    if names_only {
        fmt.Println(name)
    } else {
        fmt.Printf("Name: %v \n    input: \"%v\" \n    output: \"%v\"\n",name, cs.Input, cs.Output)
    }
}

func infoGetArgs(cmd *cobra.Command, data *hsdata.InfoData) {
	data.Names_only, _ = cmd.Flags().GetBool("names-only")
	data.Defaults_only, _ = cmd.Flags().GetBool("defaults-only")
	data.Formats_only, _ = cmd.Flags().GetBool("formats-only")
	data.Name, _ = cmd.Flags().GetString("name")
	data.Verbosity, _ = cmd.PersistentFlags().GetCount("verbose")
}
