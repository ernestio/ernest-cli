/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"unicode"

	"github.com/fatih/color"
	"github.com/urfave/cli"

	"gopkg.in/yaml.v2"
)

// setup ...
func setup(c *cli.Context) (*Manager, *Config) {
	config := getConfig()
	if config == nil {
		config = &Config{}
		if c.Command.Name != "target" {
			color.Red("Environment not configured, please use target command")
		}
	}
	m := Manager{URL: config.URL}
	return &m, config
}

func getServiceUUID(output []byte) (string, error) {
	var service struct {
		ID string `json:"id"`
	}
	err := json.Unmarshal(output, &service)
	if err != nil {
		return "", err
	}
	return service.ID, nil
}

func buildServiceUUID(payload []byte) (string, error) {
	var definition struct {
		Name       string
		Datacenter string
	}
	err := yaml.Unmarshal(payload, &definition)
	if err != nil {
		return "", err
	}
	return definition.Name + "-" + definition.Datacenter, nil
}

// askForConfirmation uses Scanln to parse user input. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user. Typically, you should use fmt to print out a question
// before calling askForConfirmation. E.g. fmt.Println("WARNING: Are you sure? (yes/no)")
func askForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return askForConfirmation()
	}
}

// posString returns the first index of element in slice.
// If slice does not contain element, returns -1.
func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

// containsString returns true iff slice contains element
func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}

// generate random string
func randString(n int) string {
	g := big.NewInt(0)
	max := big.NewInt(130)
	bs := make([]byte, n)

	for i := range bs {
		g, _ = rand.Int(rand.Reader, max)
		r := rune(g.Int64())
		for !unicode.IsNumber(r) && !unicode.IsLetter(r) {
			g, _ = rand.Int(rand.Reader, max)
			r = rune(g.Int64())
		}
		bs[i] = byte(g.Int64())
	}
	return string(bs)
}

type DatacenterTemplate struct {
	URL      string `yaml:"vcloud-url"`
	Network  string `yaml:"public-network"`
	Org      string `yaml:"org"`
	Password string `yaml:"password"`
	User     string `yaml:"user"`
	Token    string `yaml:"aws_secret_access_key"`
	Secret   string `yaml:"aws_access_key_id "`
	Region   string `yaml:"region"`
	Fake     bool   `yaml:"fake"`
}

func getDatacenterTemplate(template string, t *DatacenterTemplate) (err error) {
	payload, err := ioutil.ReadFile(template)
	if err != nil {
		return errors.New("Template file '" + template + "' not found")
	}
	if yaml.Unmarshal(payload, &t) != nil {
		return errors.New("Template file '" + template + "' is not valid yaml file")
	}
	return err
}
