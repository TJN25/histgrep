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

type WriteFn func(*HsLine)

