/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/TJN25/histgrep/hsdata"
	"github.com/TJN25/histgrep/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var sCmd = &cobra.Command{
	Use:   "s",
	Short: "Search files",
	Long:  `Search files with a series of search terms.`,
	Run:   sRun,
}

func init() {
	rootCmd.AddCommand(sCmd)
	sCmd.Flags().StringP("input", "i", "stdin", "Input file (leave blank for stdin)")
	sCmd.Flags().StringP("output", "o", "stdout", "Output file (leave blank for stdout)")
	sCmd.Flags().StringP("name", "n", "-", "Name of saved format (add with histgrep add-format -n [name] -i [input] -o [output])")
	sCmd.Flags().BoolP("common-cmds", "c", false, "Keep common commands such as cd, ls and cat")
	sCmd.Flags().BoolP("no-color", "f", false, "Do not include colors in output")
	sCmd.Flags().BoolP("pager", "p", false, "Display output in pager (Bubble Tea)")
	sCmd.PersistentFlags().CountP("verbose", "v", "Level of verbosity (0-5) default (0)")
}

func sRun(cmd *cobra.Command, args []string) {
	data := hsdata.HsData{Terms: args}
	verbosity, _ := cmd.PersistentFlags().GetCount("verbose")
	utils.SetVerbosity(verbosity)
	sGetArgs(cmd, &data)
	log.Info(fmt.Sprintf("\n    Running search with: \n    files: %v -> %v\n    Terms: %v\n    Format: %v\n", data.Input_file, data.Output_file, data.Terms, data.FormatData))
	log.Debug(fmt.Sprintf("Formatting input: %v", data))
	DoFormatting(&data)
	RunLoopFile(&data)
	// SaveHistory(&data)
}

func sGetArgs(cmd *cobra.Command, data *hsdata.HsData) {
	data.Input_file, _ = cmd.Flags().GetString("input")
	data.Output_file, _ = cmd.Flags().GetString("output")
	data.Name, _ = cmd.Flags().GetString("name")
	data.KeepCommonCmds, _ = cmd.Flags().GetBool("common-cmds")
	data.NoColor, _ = cmd.Flags().GetBool("no-color")
	data.UsePager, _ = cmd.Flags().GetBool("pager")
	if data.Name == "-" {
		data.FormatData = UseDefaults(data)
		log.Debug(data.FormatData)
	} else {
		file := utils.GetDataPath("formats.json")
		formatMap := hsdata.FormatMap{}
		utils.FetchFormatting(file, &formatMap)
		format_data, ok := formatMap[data.Name]
		if ok {
			data.FormatData = format_data
		} else {
			utils.ErrorExit(fmt.Sprintf("Format not found: %v", data.Name))
		}
		data.FormatData = formatMap[data.Name]
		log.Debug(formatMap)
	}
	log.Info(fmt.Sprintf("Args data: %v", data))
}

/*
TODO: Process the current separator, skip_by, and skip_dir in some way that
TODO: makes sense and allows the later functions to work with the new format.
TODO: Update the other functions to use this format.
*/

/*
TODO: add a feature that allows each search term to have a set of exclude terms.
*/

func DoFormatting(data *hsdata.HsData) hsdata.FormattingData {
	format_data := hsdata.FormattingData{}
	return format_data
}

func GetFormat(curr string, names *[]string, separators *[]string, positions *[]hsdata.FormatPosition) {
	log.Debug(fmt.Sprintf("%v: %v", utils.CallerName(0), curr))
	words := strings.Split(curr, "}")
	var name_sep []string
	for _, word := range words {
		pos := hsdata.FormatPosition{}
		log.Debug(fmt.Sprintf("Word: %v", word))
		name_sep = strings.Split(word, "{")
		if strings.Contains(name_sep[0], "...") {
			MultiplePositions(&name_sep)
		}
		curr_name, separator := NameAndSeparator(name_sep, &pos)
		if separator != "" {
			*separators = append(*separators, separator)
		}
		if curr_name != "" {
			*names = append(*names, curr_name)
		}
		*positions = append(*positions, pos)
	}
	log.Debug("GetFormat: Finished.")
}

func MultiplePositions(name_sep *[]string) {
	current_separator, skip_by, skip_dir := SkipSeperators((*name_sep)[0])
	fmt.Printf("current_separator: %v skip_by: %v skip_dir: %v\n", current_separator, skip_by, skip_dir)
	fmt.Printf("New formatting used - name: %v\n", name_sep)
	os.Exit(0)
}

func NameAndSeparator(name_sep []string, pos *hsdata.FormatPosition) (string, string) {
	var separator, name string
	if len(name_sep) < 2 {
		log.Debug(fmt.Sprintf("Short section - Name: N/A, Separator: %v", name_sep[0]))
		separator = name_sep[0]
		pos.Separator = name_sep[0]
	} else {
		log.Debug(fmt.Sprintf("Name: %v, Separator: %v", name_sep[1], name_sep[0]))
		if name_sep[0] == "" {
			name = GetName(name_sep[1], pos)
			return name, ""
		}
		pos.Separator = name_sep[0]
		separator = name_sep[0]
		name = GetName(name_sep[1], pos)
	}
	return name, separator
}

func GetName(name string, pos *hsdata.FormatPosition) string {
	if strings.Contains(name, ":") {
		names := strings.Split(name, ":")
		log.Debug(fmt.Sprintf("Name and color: %v", names))
		name = names[0]
		pos.Name = name
		pos.Color, pos.ColorMap = GetColor(names[1])
	} else {
		log.Debug(fmt.Sprintf("Name: %v", name))
		pos.Name = name
		pos.Color, pos.ColorMap = GetColor("none")
	}
	return name
}

func GetColor(word string) (string, map[string]string) {
	var colors = map[string]string{}
	if strings.Contains(word, ";") {
		words := strings.Split(word, ";")
		default_found := false

		for _, curr := range words {
			if strings.Contains(curr, "=") {
				items := strings.Split(curr, "=")
				colors[items[0]] = items[1]
			} else {
				colors["default"] = curr
				word = curr
				default_found = true
			}
		}
		if !default_found {
			colors["default"] = "none"
			word = "none"
		}
	} else {
		colors["default"] = word
	}
	return word, colors
}

func GetFormatPositons(curr string, format_data *[]hsdata.FormatPosition) {
	log.Debug(fmt.Sprintf("%v: %v", utils.CallerName(0), curr))
	words := strings.Split(curr, "}")
	var name_sep []string
	var c_sep, c_name string
	for _, word := range words {
		log.Debug(fmt.Sprintf("Word: %v", word))
		name_sep = strings.Split(word, "{")
		if strings.Contains(name_sep[0], "...") {
			c_sep, skip_by, skip_dir := SkipSeperators(name_sep[0])
			if len(name_sep) < 2 {
				c_name = "N/A"
			} else {
				c_name = name_sep[1]
			}
			*format_data = append(*format_data, hsdata.FormatPosition{Separator: c_sep, Name: c_name, Direction: skip_dir, Range: skip_by})
		} else {
			if len(name_sep) < 2 {
				if len(name_sep) == 0 {
					continue
				}
				c_sep = name_sep[0]
				c_name = "N/A"
				log.Debug(fmt.Sprintf("Name: %v, Separator: %v", c_name, c_sep))
				*format_data = append(*format_data, hsdata.FormatPosition{Separator: c_sep, Name: c_name})
			} else {
				log.Debug(fmt.Sprintf("Name: %v, Separator: %v", name_sep[0], name_sep[1]))
				if name_sep[0] != "" {
					c_sep = name_sep[0]
					c_name = name_sep[1]
					*format_data = append(*format_data, hsdata.FormatPosition{Separator: c_sep, Name: c_name})
				} else {
					c_name = name_sep[1]
					*format_data = append(*format_data, hsdata.FormatPosition{Name: c_name})
					continue
				}
			}
		}
	}
	log.Debug("GetFormatPos: Finished.")
}

func RunLoopFile(data *hsdata.HsData) {
	var err error
	line := hsdata.HsLine{}
	if data.UsePager {
		formatted_lines, err := utils.LoopFile(data, utils.SaveLine, line)
		if err != nil {
			log.Panic(err)
		}
		utils.ViewFileWithPager(formatted_lines, data, line)
	} else if data.Output_file == "stdout" {
		_, err = utils.LoopFile(data, utils.PrintLine, line)
	} else {
		f, err := os.Create(data.Output_file)
		if err != nil {
			log.Panic(err)
		}
		// remember to close the file
		defer f.Close()

		line.F = f
		_, err = utils.LoopFile(data, utils.WriteLine, line)
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
		} else if (int(separator[i]-48) > 0) && (int(separator[i]-48) < 10) {
			skip_by *= 10
			skip_by += int(separator[i] - 48)
		}
	}
	return current_separator, skip_by, skip_dir
}

func UseDefaults(data *hsdata.HsData) hsdata.FormattingData {
	file := utils.GetDataPath("defaults.json")
	config_file := utils.GetDataPath("formats.json")
	log.Info(fmt.Sprintf("Using defaults file %v", file))
	log.Info(fmt.Sprintf("Using config file %v", config_file))
	formatMap := hsdata.FormatMap{}
	defaults := hsdata.DefaultsData{}
	utils.FetchFormatting(config_file, &formatMap)
	utils.FetchDefaults(file, &defaults)
	return formatMap.Get(defaults.Name)
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
	os.WriteFile(file, b, os.ModePerm)
}

const PERIOD byte = 46
