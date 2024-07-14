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

	configCmd.Flags().BoolP("add-format", "a", false, "Add a new format options with -i [input] -o [output] -n [name]")
	configCmd.Flags().BoolP("change-format", "c", false, "Add a new format options with -i [input] -o [output] -n [name]")
	configCmd.Flags().BoolP("delete-format", "d", false, "Delete a format with -n [name]")
	configCmd.Flags().BoolP("add-default", "A", false, "Add a new default with -i [input] -o [output] -n [name]")
    configCmd.Flags().BoolP("change-default", "C", false, "Add a new default with -i [input] -o [output] -n [name]")
    configCmd.Flags().BoolP("delete-default", "D", false, "Delete a default with -n [name]")

	configCmd.Flags().StringP("input", "i", "-", "Input format")
    configCmd.Flags().StringP("output", "o", "-", "Output format")
    configCmd.Flags().StringP("name", "n", "-", "Format name")

	configCmd.PersistentFlags().CountP("verbose", "v", "Level of verbosity (0-5) default (0)")
}

func configRun(cmd *cobra.Command, args []string) {
	data := hsdata.ConfigData{}
	configGetArgs(cmd, &data)
	setVerbosity(&data)

	configCallFunction(&data)
}

func configGetArgs(cmd *cobra.Command, data *hsdata.ConfigData) {

	data.Input, _ = cmd.Flags().GetString("input")
	data.Output, _ = cmd.Flags().GetString("output")
	data.Name, _ = cmd.Flags().GetString("name")

	add, _ := cmd.Flags().GetBool("add-format")
    change, _ := cmd.Flags().GetBool("change-format")
    delete, _ := cmd.Flags().GetBool("delete-format")

	add_default, _ := cmd.Flags().GetBool("add-default")
	change_default, _ := cmd.Flags().GetBool("change-default")
	delete_default, _ := cmd.Flags().GetBool("delete-default")

	if add {
		data.Action = ADD
	}else if change {
        data.Action = CHANGE
    }else if add_default {
        data.Action = ADD_DEFULT
    }else if change_default {
        data.Action = CHANGE_DEFULT
    }else if delete {
        data.Action = DELETE
    }else if delete_default {
        data.Action = DELETE_DEFULT
    }

	data.Verbosity, _ = cmd.PersistentFlags().GetCount("verbose")
}

func configCallFunction(data *hsdata.ConfigData) {
	name, err := utils.GetDataPath("histgrep/formats.json")
	if err != nil {
			fmt.Printf("Searched for %v/histgrep/formats.json and %v/.histgrep/formats.json\n", utils.XDG_CONFIG_HOME, utils.HOME_PATH)
			fmt.Println("Please create the config directory ($XDG_CONFIG_HOME/histgrep/formats.json or $HOME/.histgrep/formats.json)")
			os.Exit(1)
    }
	data.Path = name

	log.Info(data.Path)

	switch data.Action {
	case ADD:
		add(data)
	case CHANGE:
		change(data)
	case DELETE:
		delete_config(data)
	case ADD_DEFULT:
		add_default(data)
	case CHANGE_DEFULT:
		change_default(data)
	case DELETE_DEFULT:
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
	fmt.Println(configMap)
	os.WriteFile(data.Path, b, os.ModePerm)
}

func change(data *hsdata.ConfigData) {

}

func delete_config(data *hsdata.ConfigData) {
	jsonFile, err := os.ReadFile(data.Path)
	if err != nil {
		log.Fatal(fmt.Sprintf("Cannot or find formats.json\n%v", err))
	}
	configMap := hsdata.ConfigMap{}
	json.Unmarshal(jsonFile, &configMap)
	log.Info(configMap)

	_, ok := configMap[data.Name]
	if ok {
		delete(configMap, data.Name)
	} else {
		utils.ErrorExit(fmt.Sprintf("Name %v does not exist\n", data.Name))
	}
	b, _ := json.MarshalIndent(configMap, "", " ")
	log.Info(configMap)
	os.WriteFile(data.Path, b, os.ModePerm)
}

func add_default(data *hsdata.ConfigData) {

}

func change_default(data *hsdata.ConfigData) {

}

func delete_default(data *hsdata.ConfigData) {

}

func setVerbosity(data *hsdata.ConfigData) {

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

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
}
