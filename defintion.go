/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"errors"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// Definition ...
type Definition struct {
	// data map[interface{}]interface{}
	data yaml.MapSlice
}

// Load the yaml
func (d *Definition) Load(data []byte) error {
	return yaml.Unmarshal(data, &d.data)
}

// Save the definition as a byte slice
func (d *Definition) Save() ([]byte, error) {
	return yaml.Marshal(d.data)
}

// LoadFileImports : loads any referenced files and maps them to the import definition
func (d *Definition) LoadFileImports() error {
	var err error
	d.data, err = LoadMapSlice(d.data)
	return err
}

// LoadFile : determines if the encountered string is
func LoadFile(path string) (string, error) {
	if len(path) < 3 || path[:2] != "@{" {
		return path, nil
	}

	trimmedPath := path[2:]
	trimmedPath = trimmedPath[:len(trimmedPath)-1]

	if _, err := os.Stat(trimmedPath); os.IsNotExist(err) {
		return "", errors.New("Can't access referenced file " + trimmedPath)
	}

	payload, err := ioutil.ReadFile(trimmedPath)
	if err != nil {
		return "", errors.New("Can't access referenced file " + trimmedPath)
	}

	return string(payload), nil
}

// LoadMapSlice : loads all values into a slice
func LoadMapSlice(s yaml.MapSlice) (yaml.MapSlice, error) {
	var err error
	for i, item := range s {
		switch v := item.Value.(type) {
		case string:
			if s[i].Value, err = LoadFile(v); err != nil {
				return s, err
			}
		case yaml.MapSlice:
			if s[i].Value, err = LoadMapSlice(v); err != nil {
				return s, err
			}
		case []interface{}:
			if s[i].Value, err = LoadSlice(v); err != nil {
				return s, err
			}
		}
	}
	return s, err
}

// LoadSlice : loads all values into a slice
func LoadSlice(s []interface{}) ([]interface{}, error) {
	var err error
	for i, selector := range s {
		switch v := selector.(type) {
		case string:
			if s[i], err = LoadFile(v); err != nil {
				return s, err
			}
		case []interface{}:
			if s[i], err = LoadSlice(v); err != nil {
				return s, err
			}
		case yaml.MapSlice:
			if s[i], err = LoadMapSlice(v); err != nil {
				return s, err
			}
		}
	}
	return s, err
}
