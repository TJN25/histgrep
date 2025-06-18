package utils

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/TJN25/histgrep/hsdata"

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
	Log.Debugf("Checking for %v\n", file)
	searchPaths := []string{
		filepath.Join(HISTGREP_CONFIG_PATH, file),
		filepath.Join(XDG_CONFIG_HOME, "histgrep", file),
		filepath.Join(HOME_PATH, ".histgrep", file),
	}

	for i, path := range searchPaths {
		Log.Debugf("Checking path %d: %s\n", i+1, path)

		if _, err := os.Stat(path); err == nil {
			Log.Infof("Found config file: %s\n", path)
			return path, nil
		}
	}

	Log.Errorf("Config file '%s' not found in any of these locations:\n", file)
	for i, path := range searchPaths {
		Log.Errorf("  %d. %s\n", i+1, path)
	}

	return "", fmt.Errorf("config file '%s' not found", file)
}

func ErrorExit(msg string) {
	Log.Error(msg)
	fmt.Fprintln(os.Stderr, "exiting")
	os.Exit(1)
}

func FetchFormatting(file string, fm *hsdata.FormatMap) {
	jsonFile, err := os.ReadFile(file)
	if err != nil {
		ErrorExit(fmt.Sprintf("Cannot find %v\n%v", file, err))
	}
	json.Unmarshal(jsonFile, fm)
	Log.Infof("FetchFormatting: %v, from %v\n", fm, file)
}

func SetVerbosity(verbosity int) {

	InitializeLogger(verbosity)
	// log.SetFormatter(&log.TextFormatter{
	// 	FullTimestamp: true,
	// })
	//
	// switch verbosity {
	// case 0:
	// 	log.SetLevel(log.ErrorLevel)
	// case 1:
	// 	log.SetLevel(log.WarnLevel)
	// case 2:
	// 	log.SetLevel(log.InfoLevel)
	// case 3:
	// 	log.SetLevel(log.DebugLevel)
	// case 4:
	// 	log.SetLevel(log.TraceLevel)
	// default:
	// 	log.SetLevel(log.TraceLevel)
	// }
}
