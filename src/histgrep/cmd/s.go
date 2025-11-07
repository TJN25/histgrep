/*
s.g
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/TJN25/histgrep/hsdata"
	"github.com/TJN25/histgrep/utils"
	// log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
	sCmd.Flags().BoolP("case-sensitive", "c", false, "Use case sensitive search")
	sCmd.Flags().BoolP("no-color", "f", false, "Do not include colors in output")
	sCmd.Flags().BoolP("pager", "p", false, "Display output in pager (Bubble Tea)")
	sCmd.Flags().BoolP("numbered", "", false, "Include line numbers in output")
	sCmd.Flags().StringP("exclude", "x", "SKIPEXCLUDE", "Exclude specific terms from output")
	sCmd.PersistentFlags().CountP("verbose", "v", "Level of verbosity (0-5) default (0)")
}

func sRun(cmd *cobra.Command, args []string) {
	data := hsdata.HsData{Terms: args}
	verbosity, _ := cmd.PersistentFlags().GetCount("verbose")
	utils.SetVerbosity(verbosity)
	config := sGetArgs(cmd, &data)
	utils.Log.Infof("\n    Running search with: \n    files: %v -> %v\n    Terms: %v\n    Format: %v\n", data.InputFile, data.OutputFile, data.Terms, data.FormatData)
	utils.Log.Tracef("Formatting input: %+v\n", data)
	utils.Log.Debugf("HsData.Input_file: %s\n", data.InputFile)
	utils.Log.Debugf("HsData.Files: %d files\n", len(data.Files))
	utils.Log.Debugf("HsData.Output_file: %s\n", data.OutputFile)
	utils.Log.Debugf("HsData.Terms: %v\n", data.Terms)
	utils.Log.Debugf("HsData.ExcludeTerms: %v\n", data.ExcludeTerms)
	utils.Log.Debugf("HsData.FormatData: %+v\n", data.FormatData)
	utils.Log.Debugf("HsData.Name: %s\n", data.Name)
	utils.Log.Debugf("HsData.NoColor: %t\n", data.NoColor)
	utils.Log.Debugf("HsData.UsePager: %t\n", data.UsePager)
	utils.Log.Debugf("HsData.IncludeNumbers: %t\n", data.IncludeNumbers)
	utils.Log.Debugf("HsData.CaseSensitive: %t\n", data.CaseSensitive)
	DoFormatting(&data)
	RunLoopFile(&data, config)
	// SaveHistory(&data)
}

func sGetArgs(cmd *cobra.Command, data *hsdata.HsData) *utils.Config {
	data.InputFile, _ = cmd.Flags().GetString("input")
	config := DoConfigFile(data)
	data.OutputFile, _ = cmd.Flags().GetString("output")
	data.Name, _ = cmd.Flags().GetString("name")
	if cmd.Flags().Changed("no-color") {
		data.NoColor, _ = cmd.Flags().GetBool("no-color")
	}
	if cmd.Flags().Changed("pager") {
		data.UsePager, _ = cmd.Flags().GetBool("pager")
	}
	if cmd.Flags().Changed("case-sensitive") {
		data.CaseSensitive, _ = cmd.Flags().GetBool("case-sensitive")
	}
	data.IncludeNumbers, _ = cmd.Flags().GetBool("numbered")
	exclude, _ := cmd.Flags().GetString("exclude")
	if exclude == "SKIPEXCLUDE" {
		data.ExcludeTerms = []string{}
	} else {
		data.ExcludeTerms = strings.Split(exclude, " ")
	}
	if data.Name == "-" {
		data.FormatData = UseDefaults(data, config)
		utils.Log.Debugf("%+v\n", data.FormatData)
	} else {
		file, err := utils.GetDataPath("formats.json")
		if err != nil {
			utils.Log.Println("Please create the config directory ($XDG_CONFIG_HOME/histgrep/ or $HOME/.histgrep/)")
			os.Exit(1)
		}
		formatMap := hsdata.FormatMap{}
		utils.FetchFormatting(file, &formatMap)
		format_data, ok := formatMap[data.Name]
		if ok {
			data.FormatData = format_data
		} else {
			utils.ErrorExit(fmt.Sprintf("Format not found: %v", data.Name))
		}
		data.FormatData = formatMap[data.Name]
		utils.Log.Tracef("%+v\n", formatMap)
	}
	utils.Log.Tracef("Args data: %+v\n", data)
	return config
}

/*
TODO: Process the current separator, skip_by, and skip_dir in some way that
makes sense and allows the later functions to work with the new format.
Update the other functions to use this format.
add a feature that allows each search term to have a set of exclude terms.
*/

func DoFormatting(data *hsdata.HsData) hsdata.FormattingData {
	format_data := hsdata.FormattingData{}
	return format_data
}

func GetFormat(curr string, names *[]string, separators *[]string, positions *[]hsdata.FormatPosition) {
	utils.Log.Tracef("%v: %v\n", utils.CallerName(0), curr)
	words := strings.Split(curr, "}")
	var name_sep []string
	for _, word := range words {
		pos := hsdata.FormatPosition{}
		utils.Log.Debugf("Word: %v\n", word)
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
	utils.Log.Debugln("GetFormat: Finished.")
}

func MultiplePositions(name_sep *[]string) {
	current_separator, skip_by, skip_dir := SkipSeperators((*name_sep)[0])
	utils.Log.Printf("current_separator: %v skip_by: %v skip_dir: %v\n", current_separator, skip_by, skip_dir)
	utils.Log.Printf("New formatting used - name: %v\n", name_sep)
	os.Exit(0)
}

func NameAndSeparator(name_sep []string, pos *hsdata.FormatPosition) (string, string) {
	var separator, name string
	if len(name_sep) < 2 {
		utils.Log.Debugf("Short section - Name: N/A, Separator: %v\n", name_sep[0])
		separator = name_sep[0]
		pos.Separator = name_sep[0]
	} else {
		utils.Log.Debugf("Name: %v, Separator: %v\n", name_sep[1], name_sep[0])
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
		utils.Log.Debugf("Name and color: %v\n", names)
		name = names[0]
		pos.Name = name
		pos.Color, pos.ColorMap = GetColor(names[1])
	} else {
		utils.Log.Debugf("Name: %v\n", name)
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
	utils.Log.Debugf("%v: %v\n", utils.CallerName(0), curr)
	words := strings.Split(curr, "}")
	var name_sep []string
	var c_sep, c_name string
	for _, word := range words {
		utils.Log.Debugf("Word: %v\n", word)
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
				utils.Log.Debugf("Name: %v, Separator: %v\n", c_name, c_sep)
				*format_data = append(*format_data, hsdata.FormatPosition{Separator: c_sep, Name: c_name})
			} else {
				utils.Log.Debugf("Name: %v, Separator: %v\n", name_sep[0], name_sep[1])
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
	utils.Log.Debugln("GetFormatPos: Finished.")
}

func RunLoopFile(data *hsdata.HsData, config *utils.Config) {
	var err error
	line := hsdata.HsLine{}
	if data.UsePager {
		formatted_lines, err := utils.LoopFile(data, utils.SaveLine, line)
		if err != nil {
			utils.Log.Panicf("Running LoopFile failed: %v\n", err)
		}
		utils.ViewFileWithPager(formatted_lines, data, line, config)
	} else if data.OutputFile == "stdout" {
		_, err = utils.LoopFile(data, utils.PrintLine, line)
	} else {
		f, err := os.Create(data.OutputFile)
		if err != nil {
			utils.Log.Panicf("Running LoopFile failed: %v\n", err)
		}
		// remember to close the file
		defer f.Close()

		line.F = f
		_, err = utils.LoopFile(data, utils.WriteLine, line)
	}
	if err != nil {
		utils.Log.Panicf("Running LoopFile failed: %v\n", err)
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

func UseDefaults(data *hsdata.HsData, config *utils.Config) hsdata.FormattingData {
	config_file, err := utils.GetDataPath("formats.json")
	if err != nil {
		utils.Log.Println("Please create the config directory ($XDG_CONFIG_HOME/histgrep/ or $HOME/.histgrep/)")
		os.Exit(1)
	}
	utils.Log.Infof("Using config file %v\n", config_file)
	formatMap := hsdata.FormatMap{}
	utils.FetchFormatting(config_file, &formatMap)
	return formatMap.Get(config.Search.DefaultName) // This can fail and I should return an error instead.
}

func DoConfigFile(data *hsdata.HsData) *utils.Config {
	utils.Log.Debugf("Starting config file processing\n")

	var config *utils.Config
	file, err := utils.GetDataPath("histgrep.toml")
	if err != nil {
		utils.Log.Debugf("No config file found, returning nil config: %v\n", err)
		return config
	}
	utils.Log.Debugf("Found config file: %s\n", file)

	config, err = utils.LoadConfig(file)
	if err != nil {
		utils.Log.Debugf("Failed to load config, returning nil: %v\n", err)
		return config
	}
	utils.Log.Debugf("Config loaded successfully\n")

	data.CaseSensitive = config.Search.CaseSensitive
	data.UsePager = config.Display.PagerEnabled
	data.NoColor = !config.Display.ColorEnabled
	utils.Log.Tracef("Data updated - CaseSensitive: %t, UsePager: %t, NoColor: %t\n",
		data.CaseSensitive, data.UsePager, data.NoColor)

	utils.Log.Debugf("Current Input_file: %s\n", data.InputFile)
	if data.InputFile == "stdin" {
		utils.Log.Debugf("Processing stdin input\n")
		stat, _ := os.Stdin.Stat()
		utils.Log.Tracef("Stdin mode: %v, is char device: %t\n", stat.Mode(), (stat.Mode()&os.ModeCharDevice) != 0)
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			utils.Log.Debugf("Stdin is terminal, getting default log files\n")
			// Get matching log files
			logFiles, err := utils.GetMatchingLogFiles(config)
			if err != nil {
				utils.Log.Printf("Error getting log files: %v\n", err)
				os.Exit(1)
			}
			utils.Log.Tracef("Found %d log files: %v\n", len(logFiles), logFiles)
			data.Files = logFiles
			data.InputFile = "default_files"
		} else {
			utils.Log.Debugf("Stdin is not terminal, checking if empty\n")
			// Read a bit from stdin to check if it's empty
			buf := make([]byte, 1)
			_, err := os.Stdin.Read(buf)
			if err == io.EOF {
				utils.Log.Debugf("Stdin is empty, getting default log files\n")
				// Get matching log files
				logFiles, err := utils.GetMatchingLogFiles(config)
				if err != nil {
					utils.Log.Printf("Error getting log files: %v\n", err)
					os.Exit(1)
				}
				utils.Log.Tracef("Found %d log files: %v\n", len(logFiles), logFiles)
				data.Files = logFiles
				data.InputFile = "default_files"
			} else {
				utils.Log.Debugf("Stdin has data, will use stdin input\n")
			}
		}
	} else {
		utils.Log.Debugf("Not using stdin, Input_file: %s\n", data.InputFile)
	}

	utils.Log.Debugf("Final data.Input_file: %s, Files count: %d\n", data.InputFile, len(data.Files))
	return config
}

const PERIOD byte = 46
