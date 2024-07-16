/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
    "os"
    log "github.com/sirupsen/logrus"
    "github.com/TJN25/histgrep/utils"
    "github.com/TJN25/histgrep/hsdata"
    "encoding/json"

	"github.com/spf13/cobra"
)

// histCmd represents the hist command
var histCmd = &cobra.Command{
	Use:   "hist",
	Short: "Print history",
	Long: `Print history of commands and search terms that have been used.`,
	Run: histRun,
}

func init() {
	rootCmd.AddCommand(histCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// histCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// histCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func histRun(cmd *cobra.Command, args []string) {
    file := utils.GetDataPath("history.json")
    log.Info(fmt.Sprintf("Saving history to: %v", file))
	jsonFile, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(fmt.Sprintf("Cannot find history.json\n%v", err))
	}
    history := hsdata.HistoryArray{}
    json.Unmarshal(jsonFile, &history)
    history.Print("All", 0)
}
