/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"
    "encoding/json"
	"github.com/TJN25/histgrep/hsdata"
	"github.com/TJN25/histgrep/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// sCmd represents the s command
var sCmd = &cobra.Command{
	Use:   "s",
	Short: "Search files",
	Long: `Search files with a series of search terms.`,
	Run: sRun,
}

func init() {
	rootCmd.AddCommand(sCmd)

	sCmd.Flags().StringP("input", "i", "stdin", "Input file (leave blank for stdin)")
	sCmd.Flags().StringP("line-format", "l", "-", "Format of the line e.g. \"{foo}:{bar}:{string}\"")
	sCmd.Flags().StringP("output-format", "f", "-", "Format for the output e.g. \"{foo}\\t{bar}\\t{string}\"")
	sCmd.Flags().StringP("output", "o", "stdout", "Output file (leave blank for stdout)")
	sCmd.Flags().StringP("name", "n", "-", "Name of saved format (add with histgrep add-format -n [name] -i [input] -o [output])")
	sCmd.PersistentFlags().CountP("verbose", "v", "Level of verbosity (0-5) default (0)")
}


func sRun(cmd *cobra.Command, args []string) {
	data := hsdata.HsData{Terms: args}
    verbosity, _ := cmd.PersistentFlags().GetCount("verbose")
	utils.SetVerbosity(verbosity)
	sGetArgs(cmd, &data)
    log.Info(fmt.Sprintf("\n    Running search with: \n    files: %v -> %v\n    Terms: %v\n    Input Format: %v\n    Output Format: %v\n", data.Input_file, data.Output_file, data.Terms, data.LineFormat, data.OutputFormat))
	log.Debug(fmt.Sprintf("Formatting input: %v", data))
    format_data := DoFormatting(&data)
    RunLoopFile(&data, &format_data)
    SaveHistory(&data)

}

func sGetArgs(cmd *cobra.Command, data *hsdata.HsData) {

	data.Input_file, _ = cmd.Flags().GetString("input")
	data.Output_file, _ = cmd.Flags().GetString("output")

	data.Name, _ = cmd.Flags().GetString("name")

	if data.Name == "-" {
		data.LineFormat, _ = cmd.Flags().GetString("line-format")
		data.OutputFormat, _ = cmd.Flags().GetString("output-format")
        if data.LineFormat == "-" || data.OutputFormat == "-" {
            UseDefaults(data)
        }
    } else {
		file := utils.GetDataPath("formats.json")
		configMap := hsdata.ConfigMap{}
		utils.FetchFormatting(file, &configMap)
		format_c, _:= configMap[data.Name]
		if format_c.Input == "" {
			data.LineFormat, _ = cmd.Flags().GetString("line-format")
			data.OutputFormat, _ = cmd.Flags().GetString("output-format")
			if data.LineFormat == "-" || data.OutputFormat == "-" {
				utils.ErrorExit(fmt.Sprintf("Cannot find %v in %v", data.Name, file))
			}
			config_data := hsdata.ConfigData{
				Input: data.LineFormat,
				Output: data.OutputFormat,
				Name: data.Name,
				Path: file,
			}
			add(&config_data)
		} else {
			data.LineFormat = format_c.Input
			data.OutputFormat = format_c.Output
		}

		log.Info(fmt.Sprintf("Fetech formatting: %v for %v", format_c, data.Name))
	}
	log.Info(fmt.Sprintf("Args data: %v", data))
}

/*
TODO: Process the current separator, skip_by, and skip_dir in some way that
TODO: makes sense and allows the later functions to work with the new format.
TODO: Update the other functions to use this format.
*/

func DoFormatting(data *hsdata.HsData) hsdata.FormattingData {
    var names, separators, f_names, f_separators []string
	if data.LineFormat == "-" || data.OutputFormat == "-" {
		names = append(names, "BLANK")
		separators = append(names, "BLANK")
		f_names = append(names, "BLANK")
		f_separators = append(names, "BLANK")
	} else {
		formatInput(data.LineFormat, &names, &separators)
		log.Debug(fmt.Sprintf("Formatting output: %v", data))
		formatInput(data.OutputFormat, &f_names, &f_separators)
	}
	format_data := hsdata.FormattingData{
		Names: &names,
		Separators: &separators,
		Fnames: &f_names,
		Fseparators: &f_separators,
	}
    return format_data
}

func formatInput(line string, names *[]string, separators *[]string) {
	log.Debug(fmt.Sprintf("%v: %v", utils.CallerName(0), line))
    words := strings.Split(line,"}")
    var name_sep []string
    for _, word := range words {
		log.Debug(fmt.Sprintf("Word: %v", word))
		name_sep = strings.Split(word,"{")
		if strings.Contains(name_sep[0], "...") {
			current_separator, skip_by, skip_dir := SkipSeperators(name_sep[0])
			fmt.Printf("current_separator: %v skip_by: %v skip_dir: %v\n", current_separator, skip_by, skip_dir)
			fmt.Printf("New formatting used - name: %v\n", name_sep)
			os.Exit(0)
		}
		if len(name_sep) < 2 {
			log.Debug(fmt.Sprintf("Short section - Name: N/A, Separator: %v", name_sep[0]))
			*separators = append(*separators, name_sep[0])
		}else {
			log.Debug(fmt.Sprintf("Name: %v, Separator: %v", name_sep[0], name_sep[1]))
			if name_sep[0] == "" {
				*names = append(*names, name_sep[1])
				continue
			}
			*separators = append(*separators, name_sep[0])
			*names = append(*names, name_sep[1])
		}
		log.Debug("formatInput: Finished.")
    }
}

func RunLoopFile(data *hsdata.HsData, format_data *hsdata.FormattingData) {
	var err error
	line := hsdata.HsLine{}
	if data.Output_file == "stdout" {
		err = utils.LoopFile(data, utils.PrintLine, line, format_data)
	} else {
		f, err := os.Create(data.Output_file)
		if err != nil {
			log.Panic(err)
		}
	// remember to close the file
	defer f.Close()

		line.F = f
		err = utils.LoopFile(data, utils.WriteLine, line, format_data)
	}
	if err != nil {
		log.Panic(err)
	}
}

func SkipSeperators(separator string) (string, int, int) {
	write_separator := true
	current_separator := ""
	skip_by, skip_dir := 0, 1
	for i := 0; i < len(separator); i++ {
		if write_separator && separator[i] != PERIOD {
			current_separator += string(separator[i])
		} else if separator[i] == PERIOD {
			write_separator = false
		} else if separator[i] == '-' {
			skip_dir = -1
		} else if (int(separator[i] - 48) > 0) && (int(separator[i] - 48) < 10) {
			skip_by *= 10
			skip_by += int(separator[i] - 48)
		}
	}
	return current_separator, skip_by, skip_dir
}

func UseDefaults(data *hsdata.HsData) {
	file := utils.GetDataPath("defaults.json")
	config_file := utils.GetDataPath("formats.json")
	log.Info(fmt.Sprintf("Using defaults file %v", file))
	log.Info(fmt.Sprintf("Using config file %v", config_file))
	configMap := hsdata.ConfigMap{}
	defaults := hsdata.DefaultsData{}
	utils.FetchFormatting(config_file, &configMap)
	utils.FetchDefaults(file, &defaults)
    defaultsConfig := configMap.Get(defaults.Name)
    data.Name = defaults.Name
    data.LineFormat = defaultsConfig.Input
    data.OutputFormat = defaultsConfig.Output
}

func SaveHistory(data *hsdata.HsData) {
    file := utils.GetDataPath("history.json")
    log.Info(fmt.Sprintf("Saving history to: %v", file))
	jsonFile, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(fmt.Sprintf("Cannot find history.json\n%v", err))
	}
    history := hsdata.HistoryArray{}
	json.Unmarshal(jsonFile, &history)
    history.Add(*data)
	log.Info(history)
	b, _ := json.Marshal(history)
    log.Debug(fmt.Sprintf("%#v", b))
    os.WriteFile(file, b, os.ModePerm)

}
const PERIOD byte = 46
