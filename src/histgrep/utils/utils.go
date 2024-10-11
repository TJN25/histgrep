package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TJN25/histgrep/hsdata"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
)

var XDG_CONFIG_HOME string = os.Getenv("XDG_CONFIG_HOME")
var HOME_PATH string = os.Getenv("HOME")
var HISTGREP_CONFIG_PATH string = os.Getenv("HISTGREP_CONFIG_PATH")

func Btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func CallerName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip + 1)
	if !ok {
		return ""
	}
	f := runtime.FuncForPC(pc)
	if f == nil {
		return ""
	}
	return f.Name()
}

func GetDataPath(file string) (string, error) {
	log.Info(fmt.Sprintf("Checking for %v", file))

	_, err := os.Stat(fmt.Sprintf("%v/%v", HISTGREP_CONFIG_PATH, file))
	if err == nil {
		log.Debug(fmt.Sprintf("File exists: %v/%v", HISTGREP_CONFIG_PATH, file))
		return fmt.Sprintf("%v/%v", HISTGREP_CONFIG_PATH, file), nil
	} else {
		log.Warn(fmt.Sprintf("File missing: %v/%v", HISTGREP_CONFIG_PATH, file))
	}

	_, err = os.Stat(fmt.Sprintf("%v/histgrep/%v", XDG_CONFIG_HOME, file))
	if err == nil {
		return fmt.Sprintf("%v/histgrep/%v", XDG_CONFIG_HOME, file), nil
	} else {
		log.Warn(fmt.Sprintf("File missing: %v/histgrep/%v", XDG_CONFIG_HOME, file))
	}

	_, err = os.Stat(fmt.Sprintf("%v/.histgrep/%v", HOME_PATH, file))
	if err == nil {
		return fmt.Sprintf("%v/.histgrep/%v", HOME_PATH, file), nil
	} else {
		log.Warn(fmt.Sprintf("File missing: %v/.histgrep/%v", HOME_PATH, file))
		fmt.Printf("Searched for %v/%v, %v, and %v/.%v\n", XDG_CONFIG_HOME, file, file, HOME_PATH, file)
		return "", errors.New("File not found")
	}
}

func ErrorExit(msg string) {
	log.Error(msg)
	fmt.Fprintln(os.Stderr, "exiting")
	os.Exit(1)
}

func FetchFormatting(file string, fm *hsdata.FormatMap) {
	jsonFile, err := os.ReadFile(file)
	if err != nil {
		ErrorExit(fmt.Sprintf("Cannot find %v\n%v", file, err))
	}
	json.Unmarshal(jsonFile, fm)
	log.Info(fmt.Sprintf("FetchFormatting: %v, from %v", fm, file))
}

func FetchDefaults(file string, df *hsdata.DefaultsData) {
	jsonFile, err := os.ReadFile(file)
	if err != nil {
		ErrorExit(fmt.Sprintf("Cannot find %v\n%v", file, err))
	}
	json.Unmarshal(jsonFile, df)
	log.Info(fmt.Sprintf("FetchDefaults: %v, from %v", df, file))
}

func SetVerbosity(verbosity int) {

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	switch verbosity {
	case 0:
		log.SetLevel(log.ErrorLevel)
	case 1:
		log.SetLevel(log.WarnLevel)
	case 2:
		log.SetLevel(log.InfoLevel)
	case 3:
		log.SetLevel(log.DebugLevel)
	case 4:
		log.SetLevel(log.TraceLevel)
	default:
		log.SetLevel(log.TraceLevel)
	}
}
