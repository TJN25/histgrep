/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	// "io"
	"github.com/TJN25/histgrep/hsdata"
	"github.com/TJN25/histgrep/utils"
	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"
	"encoding/json"
)

const (
	ADD = iota
	CHANGE = iota
	DELETE = iota
	ADD_DEFULT = iota
	CHANGE_DEFULT = iota
	DELETE_DEFULT = iota
)


// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Modify configuration settings",
	Long: `Modify configuration settings for the histgrep command. Add, remove, 
or update formatting options and default values.`,
	Run: configRun,
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().StringP("add-format", "a", "", "Add a new format options with -i [input] -o [output] -n [name]")
	configCmd.Flags().StringP("change-format", "c", "", "Add a new format options with -i [input] -o [output] -n [name]")
	configCmd.Flags().StringP("delete-format", "d", "", "Delete a format with -n [name]")
	configCmd.Flags().StringP("add-default", "A", "", "Add a new default with -i [input] -o [output] -n [name]")
    configCmd.Flags().StringP("change-default", "C", "", "Add a new default with -i [input] -o [output] -n [name]")
    configCmd.Flags().StringP("delete-default", "D", "", "Delete a default with -n [name]")

	configCmd.Flags().StringP("input", "i", "-", "Input format")
    configCmd.Flags().StringP("output", "o", "-", "Output format")
    configCmd.Flags().StringP("name", "n", "-", "Format name")

	configCmd.PersistentFlags().CountP("verbose", "v", "Level of verbosity (0-5) default (0)")
}

func configRun(cmd *cobra.Command, args []string) {
	data := hsdata.ConfigData{}
	configGetArgs(cmd, &data)
	utils.SetVerbosity(data.Verbosity)

	configCallFunction(&data)
}

func configGetArgs(cmd *cobra.Command, data *hsdata.ConfigData) {

	data.Input, _ = cmd.Flags().GetString("input")
	data.Output, _ = cmd.Flags().GetString("output")

	add_format, _ := cmd.Flags().GetString("add-format")
    change_format, _ := cmd.Flags().GetString("change-format")
    delete_format, _ := cmd.Flags().GetString("delete-format")

	add_default_format, _ := cmd.Flags().GetString("add-default")
	change_default_format, _ := cmd.Flags().GetString("change-default")
	delete_default_format, _ := cmd.Flags().GetString("delete-default")

	if add_format != "" {
        data.Name = add_format
		data.Action = ADD
	}else if change_format != "" {
        data.Name = change_format
        data.Action = CHANGE
    }else if add_default_format != "" {
        data.Name = add_default_format
        data.Action = ADD_DEFULT
    }else if change_default_format != "" {
        data.Name = change_default_format
        data.Action = CHANGE_DEFULT
    }else if delete_format != "" {
        data.Name = delete_format
        data.Action = DELETE
    }else if delete_default_format != "" {
        data.Name = delete_default_format
        data.Action = DELETE_DEFULT
    }

	data.Verbosity, _ = cmd.PersistentFlags().GetCount("verbose")
}

func configCallFunction(data *hsdata.ConfigData) {

	switch data.Action {
	case ADD:
		name := utils.GetDataPath("formats.json")
		data.Path = name
		log.Info(data.Path)
		add(data)
	case CHANGE:
		name := utils.GetDataPath("formats.json")
		data.Path = name
		log.Info(data.Path)
		change(data)
	case DELETE:
		name := utils.GetDataPath("formats.json")
		data.Path = name
		log.Info(data.Path)
		delete_config(data)
	case ADD_DEFULT:
		name := utils.GetDataPath("defaults.json")
		data.Path = name
		log.Info(data.Path)
		add_default(data)
	case CHANGE_DEFULT:
		name := utils.GetDataPath("defaults.json")
		data.Path = name
		log.Info(data.Path)
		change_default(data)
	case DELETE_DEFULT:
		name := utils.GetDataPath("defaults.json")
		data.Path = name
		log.Info(data.Path)
		delete_default(data)
	default:
		log.Fatal("No action provided")
}
}

func add(data *hsdata.ConfigData) {

	jsonFile, err := os.ReadFile(data.Path)
	if err != nil {
		log.Fatal(fmt.Sprintf("Cannot or find formats.json\n%v", err))
	}
	configMap := hsdata.ConfigMap{}
	utils.FetchFormatting(data.Path, &configMap)
	json.Unmarshal(jsonFile, &configMap)
	log.Info(configMap)

	_, ok := configMap[data.Name]
	if ok {
		utils.ErrorExit(fmt.Sprintf("Name %v already exists\n", data.Name))
	}

	save := hsdata.ConfigSave{
		Input: data.Input, 
		Output: data.Output, 
	}

	configMap[data.Name] = save
	b, _ := json.MarshalIndent(configMap, "", " ")
	log.Info(configMap)
	os.WriteFile(data.Path, b, os.ModePerm)
}

func change(data *hsdata.ConfigData) {

	jsonFile, err := os.ReadFile(data.Path)
	if err != nil {
		log.Fatal(fmt.Sprintf("Cannot find formats.json\n%v", err))
	}
	configMap := hsdata.ConfigMap{}
	utils.FetchFormatting(data.Path, &configMap)
	json.Unmarshal(jsonFile, &configMap)
	log.Info(configMap)

	_, ok := configMap[data.Name]
	if !ok {
		log.Warn(fmt.Sprintf("%v not in formats.json\n", data.Name))
		add(data)
		return
	}

	save := hsdata.ConfigSave{
        Input: data.Input,
        Output: data.Output,
    }
    configMap.Update(data.Name, save)
	b, _ := json.MarshalIndent(configMap, "", " ")
	log.Info(configMap)
	os.WriteFile(data.Path, b, os.ModePerm)
}

func delete_config(data *hsdata.ConfigData) {
	jsonFile, err := os.ReadFile(data.Path)
	if err != nil {
		log.Fatal(fmt.Sprintf("Cannot find formats.json\n%v", err))
	}
	configMap := hsdata.ConfigMap{}
	json.Unmarshal(jsonFile, &configMap)
	log.Info(configMap)

    err = configMap.Delete_config(data.Name)
    if err != nil {
		log.Fatal(fmt.Sprintf("Cannot find %v in formats.json\n%v", data.Name, err))
    }
	b, _ := json.MarshalIndent(configMap, "", " ")
	log.Info(configMap)
	os.WriteFile(data.Path, b, os.ModePerm)
}

func add_default(data *hsdata.ConfigData) {
	save := hsdata.DefaultsData{
		Name: data.Name, 
	}
	b, _ := json.MarshalIndent(save, "", " ")
	log.Info(save)
	os.WriteFile(data.Path, b, os.ModePerm)
}

func change_default(data *hsdata.ConfigData) {

}

func delete_default(data *hsdata.ConfigData) {

}

