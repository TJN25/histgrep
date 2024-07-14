/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// rgxCmd represents the rgx command
var rgxCmd = &cobra.Command{
	Use:   "rgx",
	Short: "Search with regex",
	Long: `Search files with a regular expression.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("rgx called")
	},
}

func init() {
	rootCmd.AddCommand(rgxCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rgxCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rgxCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
