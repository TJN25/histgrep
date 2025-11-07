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

func LoopFile(hsDat *hsdata.HsData, write_fn hsdata.WriteFn, currentLine hsdata.HsLine) ([]string, error) {
	Log.Tracef("%v: Loop file: %v\n", CallerName(0), hsDat)
	Log.Debugf("Input: %v, Output: %v, Color: %v, Excludes: %v\n", hsDat.FormatData.Input, hsDat.FormatData.Output, hsDat.FormatData.Color, hsDat.FormatData.Excludes)

	hasConditionalTerms := false
	searchTerms := make([]searchTerm, 0)

	for _, term := range hsDat.Terms {
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
	_, ok := hsDat.Reader.(*BufferedInput)
	if !ok {
		var err error
		hsDat.Reader, err = GetScanner(hsDat)
		if err != nil {
			var formatted_lines []string
			Log.Fprintf(os.Stderr, "Scanner error: %v\n", err)
			return formatted_lines, err
		}
	}

	lines_remaining := true
	matchFound := false
	lineCount := 0
	for lines_remaining {
		var line string
		var err error
		switch r := hsDat.Reader.(type) {
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
		if len(hsDat.ExcludeTerms) > 0 {
			for _, term := range hsDat.ExcludeTerms {
				var searchLine string
				var searchTerm string
				if !hsDat.CaseSensitive {
					searchLine = strings.ToLower(line)
					searchTerm = strings.ToLower(term)
				} else {
					searchLine = line
					searchTerm = term
				}

				if strings.Contains(searchLine, searchTerm) {
					do_write = false
					break
				}
			}
		}
		currentLine.Line = line // The full line is assigned here. `do_write` will control if it's used for output.

		if len(searchTerms) > 0 {
			if strings.Contains(line, "histgrep") {
				continue
			}

			// First pass: Basic contains check for all terms.
			for _, term := range searchTerms {
				searchLine := line
				searchTerm := term.term
				if !hsDat.CaseSensitive {
					searchLine = strings.ToLower(line)
					searchTerm = strings.ToLower(term.term)
				}

				if !strings.Contains(searchLine, searchTerm) {
					do_write = false
					break
				}
			}

			// Second pass: Conditional checks for startsWith/endsWith for matching lines
			if do_write && hasConditionalTerms {
				wordsMap := getInputNames(line, &hsDat.FormatData)

				for _, term := range searchTerms {
					if term.conditional == "StartsWith" || term.conditional == "EndsWith" {
						conditionalMatchFound := false
						for _, currentSegment := range wordsMap {
							searchSegment := currentSegment
							searchTerm := term.term
							if !hsDat.CaseSensitive {
								searchSegment = strings.ToLower(currentSegment)
								searchTerm = strings.ToLower(term.term)
							}

							if term.conditional == "StartsWith" && strings.HasPrefix(searchSegment, searchTerm) {
								conditionalMatchFound = true
								break
							}
							if term.conditional == "EndsWith" && strings.HasSuffix(searchSegment, searchTerm) {
								conditionalMatchFound = true
								break
							}
						}
						if !conditionalMatchFound {
							do_write = false
							break
						}
					}
				}
			}
		} else {
			currentLine.Line = line
		}
		if do_write {
			lineCount++
			matchFound = true
			if (hsDat.FormatData.Output["keys"])[0] != "BLANK" {
				wordsMap := getInputNames(currentLine.Line, &hsDat.FormatData)
				Log.Tracef("%+v\n", wordsMap)
				currentLine.Line = FormatLine(&wordsMap, &hsDat.FormatData, hsDat.NoColor)
				if hsDat.IncludeNumbers {
					numberStr := strconv.Itoa(lineCount)
					padding := 4 - len(numberStr)
					if padding > 0 {
						numberStr = strings.Repeat(" ", padding) + numberStr
					}
					currentLine.Line = numberStr + "| " + currentLine.Line
				}
				Log.Debugf("%s\n", currentLine)
			}
			if currentLine.Line != "" {
				write_fn(&currentLine)
			}
		}
	}
	if !matchFound {
		return []string{"No matches found for the given terms"}, nil
	}
	return currentLine.OutLines, nil
}

func GetScanner(hsDat *hsdata.HsData) (interface{}, error) {
	if hsDat.InputFile == "stdin" {
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
			return getBufferedInputFromFiles(hsDat.Files)
		}
		return &BufferedInput{content: lines}, nil
	} else if hsDat.InputFile == "default_files" {
		return getBufferedInputFromFiles(hsDat.Files)
	} else {
		file, err := os.Open(hsDat.InputFile)
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
