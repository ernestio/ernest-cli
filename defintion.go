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
	data map[interface{}]interface{}
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
	d.data = LoadMap(d.data)
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

// LoadSlice : loads all values into a slice
func LoadSlice(s []interface{}) []interface{} {
	for i, selector := range s {
		switch v := selector.(type) {
		case string:
			s[i] = LoadFile(v)
		case []interface{}:
			s[i] = LoadSlice(v)
		case map[interface{}]interface{}:
			s[i] = LoadMap(v)
		}
	}
	return s
}

// LoadMap : loads all values into a map
func LoadMap(m map[interface{}]interface{}) map[interface{}]interface{} {
	for field, selector := range m {
		switch v := selector.(type) {
		case string:
			m[field] = LoadFile(v)
		case []interface{}:
			m[field] = LoadSlice(v)
		case map[interface{}]interface{}:
			m[field] = LoadMap(v)
		}
	}

	return m
}
