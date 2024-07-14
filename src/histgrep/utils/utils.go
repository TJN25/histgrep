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
	_, err := os.Stat(fmt.Sprintf("%v/%v", XDG_CONFIG_HOME, file)) 
	if err != nil {
	_, err := os.Stat(fmt.Sprintf("%v/.%v", HOME_PATH, file)) 
		if err != nil {
            return "", err
		}
			return fmt.Sprintf("%v/.%v", HOME_PATH, file), nil
    } else {
		return fmt.Sprintf("%v/%v", XDG_CONFIG_HOME, file), nil
	}
}

func ErrorExit(msg string) {
		log.Error(msg)
		fmt.Fprintln(os.Stderr, "exiting")
		os.Exit(1)
}

func FetchFormatting(file string, name string, configMap *hsdata.ConfigMap) *hsdata.ConfigMap {

	jsonFile, err := os.ReadFile(file)
	if err != nil {
        ErrorExit(fmt.Sprintf("Cannot find %v\n%v", file, err))
	}
	json.Unmarshal(jsonFile, configMap)
	log.Info(configMap)
    return configMap

}
