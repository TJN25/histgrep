package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/TJN25/histgrep/hsdata"
	// log "github.com/sirupsen/logrus"
)

type searchTerm struct {
	term        string
	conditional string
}

func LoopFile(hs_dat *hsdata.HsData, write_fn hsdata.WriteFn, current_line hsdata.HsLine) ([]string, error) {
	Log.Tracef("%v: Loop file: %v\n", CallerName(0), hs_dat)
	Log.Debugf("Input: %v, Output: %v, Color: %v, Excludes: %v\n", hs_dat.FormatData.Input, hs_dat.FormatData.Output, hs_dat.FormatData.Color, hs_dat.FormatData.Excludes)

	hasConditionalTerms := false
	searchTerms := make([]searchTerm, 0)

	for _, term := range hs_dat.Terms {
		if strings.HasPrefix(term, "^") {
			hasConditionalTerms = true
			searchTerms = append(searchTerms, searchTerm{term: strings.TrimPrefix(term, "^"), conditional: "StartsWith"})
		} else if strings.HasSuffix(term, "$") {
			hasConditionalTerms = true
			searchTerms = append(searchTerms, searchTerm{term: strings.TrimSuffix(term, "$"), conditional: "EndsWith"})
		} else {
			searchTerms = append(searchTerms, searchTerm{term: term, conditional: "Contains"})
		}
	}
	_, ok := hs_dat.Reader.(*BufferedInput)
	if !ok {
		var err error
		hs_dat.Reader, err = GetScanner(hs_dat)
		if err != nil {
			var formatted_lines []string
			Log.Fprintf(os.Stderr, "Scanner error: %v\n", err)
			return formatted_lines, err
		}
	}

	lines_remaining := true
	match_found := false
	line_count := 0
	for lines_remaining {
		var line string
		var err error
		switch r := hs_dat.Reader.(type) {
		case *bufio.Reader:
			line, err = r.ReadString('\n')
			line, _ = strings.CutSuffix(line, "\n")
		case *BufferedInput:
			line, err = r.ReadLine()
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			Log.Printf("Read: %s, error: %v\n", line, err)
			continue
		}
		do_write := true
		if len(hs_dat.ExcludeTerms) > 0 {
			for _, term := range hs_dat.ExcludeTerms {
				var search_line string
				var search_term string
				if !hs_dat.CaseSensitive {
					search_line = strings.ToLower(line)
					search_term = strings.ToLower(term)
				} else {
					search_line = line
					search_term = term
				}

				if strings.Contains(search_line, search_term) {
					do_write = false
					break
				}
			}
		}
		if len(hs_dat.Terms) > 0 {
			if strings.Contains(line, "histgrep") {
				continue
			}
			for _, term := range hs_dat.Terms {
				var search_line string
				var search_term string
				if !hs_dat.CaseSensitive {
					search_line = strings.ToLower(line)
					search_term = strings.ToLower(term)
				} else {
					search_line = line
					search_term = term
				}

				if strings.Contains(search_line, search_term) {
					current_line.Line = line
				} else {
					do_write = false
					break
				}
			}
		} else {
			current_line.Line = line
		}
		if do_write {
			line_count++
			match_found = true
			if (hs_dat.FormatData.Output["keys"])[0] != "BLANK" {
				words_map := getInputNames(current_line.Line, &hs_dat.FormatData)
				Log.Tracef("%+v\n", words_map)
				current_line.Line = FormatLine(&words_map, &hs_dat.FormatData, hs_dat.NoColor)
				if hs_dat.IncludeNumbers {
					number_str := strconv.Itoa(line_count)
					padding := 4 - len(number_str)
					if padding > 0 {
						number_str = strings.Repeat(" ", padding) + number_str
					}
					current_line.Line = number_str + "| " + current_line.Line
				}
				Log.Debugf("%s\n", current_line)
			}
			if current_line.Line != "" {
				write_fn(&current_line)
			}
		}
	}
	if !match_found {
		return []string{"No matches found for the given terms"}, nil
	}
	return current_line.OutLines, nil
}

func GetScanner(hs_dat *hsdata.HsData) (interface{}, error) {
	if hs_dat.Input_file == "stdin" {
		var lines []string
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		if len(lines) == 0 {
			// If stdin is empty, use default files
			return getBufferedInputFromFiles(hs_dat.Files)
		}
		return &BufferedInput{content: lines}, nil
	} else if hs_dat.Input_file == "default_files" {
		return getBufferedInputFromFiles(hs_dat.Files)
	} else {
		file, err := os.Open(hs_dat.Input_file)
		if err != nil {
			return nil, err
		}
		return bufio.NewReader(file), nil
	}
}

func getBufferedInputFromFiles(files []string) (*BufferedInput, error) {
	var allLines []string
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("error reading file %s: %v", file, err)
		}
		lines := strings.Split(string(content), "\n")
		allLines = append(allLines, lines...)
	}
	return &BufferedInput{content: allLines}, nil
}

func WriteLine(line *hsdata.HsLine) {
	_, err := line.F.WriteString(line.Line + "\n")
	if err != nil {
		Log.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func PrintLine(line *hsdata.HsLine) {
	Log.Fprintf(os.Stdout, "%s\n", line.Line)
}

func SaveLine(line *hsdata.HsLine) {
	line.OutLines = append(line.OutLines, line.Line)
}

type MapFormat map[string]string

func FormatLine(terms *MapFormat, format_data *hsdata.FormattingData, no_color bool) string {
	f_keys := (*format_data).Output["keys"]
	f_separators := (*format_data).Output["separators"]
	f_colors := format_data.Color
	f_excludes := (*format_data).Excludes
	Log.Debugf("Terms: %v, Names: %v, Separators: %v\n", terms, f_keys, f_separators)
	var line string = ""
	sep_len := len(f_separators)
	for i, term := range f_keys {
		excludes, ok := f_excludes[term]
		if ok {
			starts_with, ok := excludes["starts_with"]
			if ok {
				for _, exclude := range starts_with {
					if strings.HasPrefix((*terms)[term], exclude) {
						return ""
					}
				}
			}
			contains, ok := excludes["contains"]
			if ok {
				for _, exclude := range contains {
					if strings.Contains((*terms)[term], exclude) {
						return ""
					}
				}
			}
			ends_with, ok := excludes["ends_with"]
			if ok {
				for _, exclude := range ends_with {
					if strings.HasSuffix((*terms)[term], exclude) {
						return ""
					}
				}
			}
		}
		color_map, ok := f_colors[term]
		if ok {
			color := color_map["default"]
			for key, try_color := range color_map {
				if strings.Contains((*terms)[term], key) {
					color = try_color
					break
				}
			}
			if !no_color {
				line += InsertColor(color)
			}
		} else {
			color := "white"
			if !no_color {
				line += InsertColor(color)
			}
		}
		line += (*terms)[term]

		if !no_color {
			line += hsdata.ColorNone
		}
		if i < sep_len {
			color_map, ok := f_colors["SEPARATOR"]
			if ok {
				if !no_color {
					line += InsertColor(color_map["default"])
				}
			} else {
				if !no_color {
					line += hsdata.ColorNone
				}
			}
			line += f_separators[i]
			if !no_color {
				line += hsdata.ColorNone
			}
		}

	}
	return line
}

func InsertColor(color string) string {
	if color == "red" {
		return hsdata.ColorRed
	} else if color == "green" {
		return hsdata.ColorGreen
	} else if color == "blue" {
		return hsdata.ColorBlue
	} else if color == "grey" {
		return hsdata.ColorGrey
	}
	return hsdata.ColorNone
}

func getInputNames(line string, format_data *hsdata.FormattingData) MapFormat {
	Log.Debugf("%v: Line: %v, Keys: %v, Separators: %v\n", CallerName(0), line, (*format_data).Input["keys"], (*format_data).Input["separators"])
	keys := (*format_data).Input["keys"]
	separators := (*format_data).Input["separators"]
	var words = make(MapFormat)
	var curr []string
	var remainder string = line
	var idx int = 0
	var separator_name string
	for i, separator := range separators {
		switch separator {
		case " ":
			separator_name = "SPACE"
		case "":
			separator_name = "N/A"
		case "  ":
			separator_name = "SPACE:SPACE"
		default:
			separator_name = separator
		}
		Log.Debugf("Idx: %v, Separator: %v\n", i, separator_name)
		curr = strings.SplitN(remainder, separator, 2)
		Log.Debugf("Length %v\n", len(curr))
		if len(curr) > 1 {
			Log.Debugf("Name: %v, Remainder: %v\n", curr[0], curr[1])
		} else {
			Log.Debugf("Name: %v, Remainder: N/A\n", curr[0])
		}
		if separator == "" {
			continue
		}
		if i < len(keys) {
			words[keys[i]] = curr[0]
		}
		if len(curr) < 2 {
			idx = len(keys)
			break
		}
		remainder = curr[1]
		idx = i + 1
		if remainder == "" {
			break
		}
	}
	if idx < len(keys) {
		words[keys[idx]] = remainder
	}
	return words
}

type BufferedInput struct {
	content []string
	index   int
}

func (bi *BufferedInput) ReadLine() (string, error) {
	if bi.index >= len(bi.content) {
		return "", io.EOF
	}
	line := bi.content[bi.index]
	bi.index++
	return line, nil
}

func (bi *BufferedInput) Reset() {
	bi.index = 0
}
