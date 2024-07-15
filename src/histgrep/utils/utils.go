package utils

import (
    "runtime"
    "os"
    "fmt"
    "encoding/json"
    "github.com/TJN25/histgrep/hsdata"
    log "github.com/sirupsen/logrus"
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
    }
    return "", err
}

func ErrorExit(msg string) {
		log.Error(msg)
		fmt.Fprintln(os.Stderr, "exiting")
		os.Exit(1)
}

func FetchFormatting(file string, configMap *hsdata.ConfigMap) *hsdata.ConfigMap {

	jsonFile, err := os.ReadFile(file)
	if err != nil {
        ErrorExit(fmt.Sprintf("Cannot find %v\n%v", file, err))
	}
	json.Unmarshal(jsonFile, configMap)
	log.Info(configMap)
    return configMap

}
