/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package model

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// Config is the configuration struct for this app
type Config struct {
	URL          string `json:"url"`
	Token        string `json:"token"`
	User         string `json:"user"`
	Password     string `json:"-"`
	UserID       string `json:"userid"`
	Verification string `json:"verification_code"`
}

// GetConfig : Get config defined on the .ernest file
func GetConfig() *Config {
	source := getConfigPath()
	c := Config{}
	file, err := os.Open(source)
	if err != nil {
		return nil
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		log.Println("Config file is invalid")
		log.Panic("error:", err)
	}
	c.URL = strings.TrimSuffix(c.URL, "/")
	return &c
}

// Get the config path to use, default is .ernest on the same
// folder, but you can override it with a same named file on
// your home
func getConfigPath() string {
	dir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	return dir + "/.ernest"
}

// SaveConfig ...
func SaveConfig(c *Config) error {
	body, err := json.Marshal(&c)
	if err != nil {
		return errors.New("Can't save config file")
	}
	err = ioutil.WriteFile(getConfigPath(), body, 0600)
	if err != nil {
		return errors.New("Can't save config file")
	}

	return nil
}
