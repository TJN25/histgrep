package utils
import (
	"strings"
	"bufio"
    "io"
	"os"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/TJN25/histgrep/hsdata"
)

func LoopFile(hs_dat *hsdata.HsData, write_fn hsdata.WriteFn, current_line hsdata.HsLine) error {
	log.Info(fmt.Sprintf("%v: Loop file: %v", CallerName(0), hs_dat))
	log.Info(fmt.Sprintf("Input: %v, Output: %v, Color: %v, Excludes: %v\n", hs_dat.FormatData.Input, hs_dat.FormatData.Output, hs_dat.FormatData.Color, hs_dat.FormatData.Excludes))
    fmt.Println("")
	reader, err := GetScanner(hs_dat.Input_file)
	var do_write bool = true
	if err != nil {
		fmt.Fprintln(os.Stderr, "Scanner error", err)
		return err
	}

    lines_remaining := true
    for lines_remaining {
        line, err := reader.ReadString('\n')
        line, _ = strings.CutSuffix(line, "\n")
        if err != nil {
            if err == io.EOF {
                break
            }
            fmt.Printf("Read: %s, error: %v\n", line, err)
            continue
        }
		do_write = true

		if strings.Contains(line, "histgrep") {
			continue
		}
		for _, term := range hs_dat.Terms {
			if strings.Contains(line, term) {
				current_line.Line = line
			}else {
				do_write = false
				break 
			}
		}
		if do_write {
			if (hs_dat.FormatData.Output["keys"])[0] != "BLANK" {
				words_map := getInputNames(current_line.Line, &hs_dat.FormatData)
				log.Debug(words_map)
				current_line.Line = FormatLine(&words_map, &hs_dat.FormatData)
				log.Debug(current_line)
			}
            if current_line.Line != "" {
                write_fn(&current_line)
            }
		}
	}
	return nil
}

func GetScanner(input_file string) (*bufio.Reader, error) {

	var reader *bufio.Reader
	if input_file == "stdin" {
		reader = bufio.NewReader(os.Stdin) // read file by line
	} else {
		file, err := os.Open(input_file) //open file
		if err != nil {
			return reader, err
		}
		reader = bufio.NewReader(file) // read file by line
	}

	return reader, nil;
}

func WriteLine(line *hsdata.HsLine) {
	_, err := line.F.WriteString(line.Line + "\n")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func PrintLine(line *hsdata.HsLine) {
    fmt.Fprintf(os.Stdout, "%s\n", line.Line)
}

type MapFormat map[string]string

func FormatLine(terms *MapFormat, format_data *hsdata.FormattingData) string {
    f_keys := (*format_data).Output["keys"]
    f_separators := (*format_data).Output["separators"]
    f_colors := format_data.Color
    f_excludes := (*format_data).Excludes
    log.Debug(fmt.Sprintf("Terms: %v, Names: %v, Separators: %v", terms, f_keys, f_separators))
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
            line += InsertColor(color)
        } else {
            color := "white"
            line += InsertColor(color)
        }
        line += (*terms)[term]
        line += hsdata.ColorNone
        if i < sep_len {
            color_map, ok := f_colors["SEPARATOR"]
            if ok {
                line += InsertColor(color_map["default"])
            } else {
                line += hsdata.ColorNone
            }
            line += f_separators[i]
            line += hsdata.ColorNone
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
	log.Debug(fmt.Sprintf("%v: Line: %v, Keys: %v, Separators: %v", CallerName(0), line, (*format_data).Input["keys"], (*format_data).Input["separators"]))
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
		log.Debug(fmt.Sprintf("Idx: %v, Separator: %v", i, separator_name))
		curr = strings.SplitN(remainder, separator, 2)
		log.Debug(fmt.Sprintf("Length %v", len(curr)))
		if len(curr) > 1 {
			log.Debug(fmt.Sprintf("Name: %v, Remainder: %v", curr[0], curr[1]))
		} else {
			log.Debug(fmt.Sprintf("Name: %v, Remainder: N/A", curr[0]))
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


