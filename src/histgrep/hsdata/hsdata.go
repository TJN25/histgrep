package hsdata

import (
	"os"
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
	Verbosity int
}

type MapFormat map[string]string

type FormattingData struct {
	LineMap *MapFormat
	Names *[]string
    Separators *[]string
	Fnames *[]string
    Fseparators *[]string
}

type FormatArray []FormatPosition

type FormatPosition struct {
	Name string
	Separator string
	Range int
	Direction int
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

type ConfigSave struct {
    Input string
    Output string
}

type WriteFn func(*HsLine)


