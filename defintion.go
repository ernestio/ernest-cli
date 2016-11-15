/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"io/ioutil"

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
func (d *Definition) LoadFileImports() {
	d.data = LoadMapSlice(d.data)
}

// LoadFile : determines if the encountered string is
func LoadFile(path string) string {
	if len(path) < 3 || path[:2] != "@{" {
		return path
	}

	trimmedPath := path[2:]
	trimmedPath = trimmedPath[:len(trimmedPath)-1]

	payload, err := ioutil.ReadFile(trimmedPath)
	if err != nil {
		panic(err)
	}

	return string(payload)
}

// LoadMapSlice : loads all values into a slice
func LoadMapSlice(s yaml.MapSlice) yaml.MapSlice {
	for i, item := range s {
		switch v := item.Value.(type) {
		case string:
			s[i].Value = LoadFile(v)
		case yaml.MapSlice:
			s[i].Value = LoadMapSlice(v)
		case []interface{}:
			s[i].Value = LoadSlice(v)
		}
	}
	return s
}

// LoadSlice : loads all values into a slice
func LoadSlice(s []interface{}) []interface{} {
	for i, selector := range s {
		switch v := selector.(type) {
		case string:
			s[i] = LoadFile(v)
		case []interface{}:
			s[i] = LoadSlice(v)
		case yaml.MapSlice:
			s[i] = LoadMapSlice(v)
		}
	}
	return s
}
