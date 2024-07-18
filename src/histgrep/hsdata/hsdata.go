package hsdata

import (
	"os"
    "fmt"
    "errors"
)

type HsLine struct {
	Line string
	F *os.File
}

type HsData struct {
	Input_file string
	Output_file string
	Terms []string
	LineFormat string
	OutputFormat string
    Name string
}

type MapFormat map[string]string

type FormattingData struct {
	LineMap *MapFormat
	Names []string
    Separators []string
	Fnames []string
    Fseparators []string
    Positions []FormatPosition
    Fpositions []FormatPosition
}

type FormatArray []FormatPosition

type FormatPosition struct {
	Name string
	Separator string
	Range int
	Direction int
    Color string
    ColorMap map[string]string
}

type ConfigData struct {
    Input string
    Output string
    Name string
	Action int
	Verbosity int
	Path string
}

type ConfigMap map[string]ConfigSave

func (cp *ConfigMap) Get(name string) ConfigSave { 
    return (*cp)[name]
}

func (cp *ConfigMap) Add(name string, cs ConfigSave) { 
    (*cp)[name] = cs
}

func (cp *ConfigMap) Update(name string, cs ConfigSave) { 
    (*cp)[name] = (*cp)[name].Update(cs.Input, cs.Output)
}

func (cp *ConfigMap) Delete_config(name string) error { 

	_, ok := (*cp)[name]
	if ok {
		delete(*cp, name)
	} else {
        return errors.New(fmt.Sprintf("Key error: %v not in ConfigMap", name))
    }

    return nil

}

type ConfigSave struct {
    Input string
    Output string
}

func (cs ConfigSave) Add(input string, output string) ConfigSave { 
	cs.Input = input
	cs.Output = output
    return cs
}

func (cs ConfigSave) Update(input string, output string) ConfigSave { 
    if input != "-" {
        cs.Input = input
    }
    if output != "-" {
        cs.Output = output
    }
    return cs
}

type InfoData struct {
	Name string
	Names_only bool
	Defaults_only bool
	Formats_only bool
	Verbosity int
}

type DefaultsData struct {
    Name string
}

type HistoryArray struct {
    Calls []HsData
}

func (ha *HistoryArray) Add(hs HsData) { 
    ha.Calls = append(ha.Calls, hs)
}

func (ha *HistoryArray) Print(print_type string, count int) {
    fmt.Println("--- History ---\n ")
    if print_type == "head" {
        for i, v := range ha.Calls {
            if i > count {
                return
            }
            PrintHistoryLine(&v)
        } 
    }else if print_type == "tail" {
        l := len(ha.Calls)
        for i, v := range ha.Calls {
            if i < (l - count) {
                continue
            }
            PrintHistoryLine(&v)
        } 
    } else {
        for _, v := range ha.Calls {
            PrintHistoryLine(&v)
        }
    }
}

func PrintHistoryLine(hs *HsData) {
    fmt.Printf("Files: %s%v%s -> %s%v%s\n    Terms: %s%v%s\n    Formats (%v): %s\"%v\"%s -> %s\"%v\"%s\n", ColorBlue, hs.Input_file, ColorNone, ColorBlue, hs.Output_file, ColorNone, ColorGreen, hs.Terms, ColorNone, hs.Name, ColorRed, hs.LineFormat, ColorNone, ColorRed, hs.OutputFormat, ColorNone)

}

type WriteFn func(*HsLine)

const ColorRed = "\033[0;31m"
const ColorGreen = "\033[0;32m"
const ColorBlue = "\033[0;34m"
const ColorNone = "\033[0m"
const ColorGrey = "\033[1;30m"
