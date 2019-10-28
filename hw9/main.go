package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

// Settings - parameters to start program
type Settings struct {
	LogFile      string `json:"logfile"`
	VerboseLevel string `json:"level"`
}

func readSettings() (Settings, error) {
	settings := Settings{}
	settingsFile := flag.String("settings", "settings.json", "path to settings file")
	flag.Parse()
	data, err := ioutil.ReadFile(*settingsFile)
	if err == nil {
		err = json.Unmarshal(data, &settings)
	}
	return settings, err
}

func main() {
	var err error
	var settings Settings
	if settings, err = readSettings(); err != nil {
		log.Fatalf("Application can't start, error: %v", err)
	}
	fmt.Println(settings)

}
