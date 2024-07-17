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

func LoopFile(hs_dat *hsdata.HsData, write_fn hsdata.WriteFn, current_line hsdata.HsLine, format_data *hsdata.FormattingData) error {
	log.Info(fmt.Sprintf("%v: Loop file: %v", CallerName(0), hs_dat))
	log.Info(fmt.Sprintf("Names: %v, Separators: %v, fn: %v, fs: %v\n", format_data.Names, format_data.Separators, format_data.Fnames, format_data.Fseparators))
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
		// line := reader.Text()
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
			if (format_data.Names)[0] != "BLANK" {
				words_map := getInputNames(current_line.Line, &format_data.Names, &format_data.Separators)
				log.Debug(words_map)
				current_line.Line = FormatLine(&words_map, &format_data.Fnames, &format_data.Fseparators)
				log.Debug(current_line)
			}
			write_fn(&current_line)
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

func FormatLine(terms *MapFormat, f_names *[]string, f_separators *[]string) string {
	log.Debug(fmt.Sprintf("Terms: %v, Names: %v, Separators: %v", terms, f_names, f_separators))
	var line string = ""
	sep_len := len(*f_separators)
	for i, term := range *f_names {
		line += (*terms)[term]
		if i < sep_len {
			line += (*f_separators)[i]
		}

	}
	return line
}

func getInputNames(line string, names *[]string, separators *[]string) MapFormat {
	log.Debug(fmt.Sprintf("%v: Line: %v, Names: %v, Separators: %v", CallerName(0), line, names, separators))
	var words = make(MapFormat)
	var curr []string
	var remainder string = line
	var idx int = 0
	var separator_name string
	for i, separator := range *separators {
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
		if i < len(*names) {
			words[(*names)[i]] = curr[0]
		}
		if len(curr) < 2 {
			idx = len(*names)
			break
		}
		remainder = curr[1]
		idx = i + 1
		if remainder == "" {
			break
		}
	}
	if idx < len(*names) {
		words[(*names)[idx]] = remainder
	}
	return words
}


