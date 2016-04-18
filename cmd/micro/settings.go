package main

import (
	"encoding/json"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// The options that the user can set
var settings Settings

// All the possible settings
var possibleSettings = []string{"colorscheme", "tabsize", "autoindent", "syntax"}

// The Settings struct contains the settings for micro
type Settings struct {
	Colorscheme string `json:"colorscheme"`
	TabSize     int    `json:"tabsize"`
	AutoIndent  bool   `json:"autoindent"`
	Syntax      bool   `json:"syntax"`
}

// InitSettings initializes the options map and sets all options to their default values
func InitSettings() {
	home, err := homedir.Dir()
	if err != nil {
		TermMessage("Error finding your home directory\nCan't load settings file")
		return
	}

	filename := home + "/.micro/settings.json"
	if _, e := os.Stat(filename); e == nil {
		input, err := ioutil.ReadFile(filename)
		if err != nil {
			TermMessage("Error reading settings.json file: " + err.Error())
			return
		}

		json.Unmarshal(input, &settings)
	} else {
		settings = DefaultSettings()
		err := WriteSettings(filename)
		if err != nil {
			TermMessage("Error writing settings.json file: " + err.Error())
		}

		os.Mkdir(home+"/.micro", 755)
	}
}

// WriteSettings writes the settings to the specified filename as JSON
func WriteSettings(filename string) error {
	var err error
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	if _, e := os.Stat(home + "/.micro"); e == nil {
		txt, _ := json.MarshalIndent(settings, "", "    ")
		err = ioutil.WriteFile(filename, txt, 0644)
	}
	return err
}

// DefaultSettings returns the default settings for micro
func DefaultSettings() Settings {
	return Settings{
		Colorscheme: "default",
		TabSize:     4,
		AutoIndent:  true,
		Syntax:      true,
	}
}

// SetOption prompts the user to set an option and checks that the response is valid
func SetOption(view *View, args []string) {
	home, err := homedir.Dir()
	if err != nil {
		messenger.Error("Error finding your home directory\nCan't load settings file")
	}

	filename := home + "/.micro/settings.json"
	if len(args) == 2 {
		option := strings.TrimSpace(args[0])
		value := strings.TrimSpace(args[1])

		if Contains(possibleSettings, option) {
			if option == "tabsize" {
				tsize, err := strconv.Atoi(value)
				if err != nil {
					messenger.Error("Invalid value for " + option)
					return
				}
				settings.TabSize = tsize
			} else if option == "colorscheme" {
				settings.Colorscheme = value
				LoadSyntaxFiles()
				view.buf.UpdateRules()
			} else if option == "syntax" {
				if value == "on" {
					settings.Syntax = true
				} else if value == "off" {
					settings.Syntax = false
				} else {
					messenger.Error("Invalid value for " + option)
					return
				}
				LoadSyntaxFiles()
				view.buf.UpdateRules()
			}
			err := WriteSettings(filename)
			if err != nil {
				messenger.Error("Error writing to settings.json: " + err.Error())
				return
			}
		} else {
			messenger.Error("Option " + option + " does not exist")
		}
	} else {
		messenger.Error("Invalid option, please use option value")
	}
}
