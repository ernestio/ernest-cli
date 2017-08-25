/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package model

import (
	"errors"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// Definition ...
type Definition struct {
	// data map[interface{}]interface{}
	data    yaml.MapSlice
	Name    string
	Project string
}

// NewDefinition : creates a new definition from a name and project
func NewDefinition(name, project string) (d Definition) {
	d.data = append(d.data, yaml.MapItem{
		Key:   "name",
		Value: name,
	})
	d.data = append(d.data, yaml.MapItem{
		Key:   "project",
		Value: project,
	})
	d.Name = name
	d.Project = project

	return
}

// Load the yaml
func (d *Definition) Load(data []byte) (err error) {
	err = yaml.Unmarshal(data, &d.data)
	for _, item := range d.data {
		if item.Key == "name" {
			d.Name = item.Value.(string)
		}
		if item.Key == "project" {
			d.Project = item.Value.(string)
		}
	}

	return
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

// AttachMap : will attach the contents of a map to the end of the definition
func (d *Definition) AttachMap(key string, m map[string]string) {
	x := yaml.MapSlice{}
	for k, v := range m {
		x = append(x, yaml.MapItem{
			Key:   k,
			Value: v,
		})
	}
	d.data = append(d.data, yaml.MapItem{
		Key:   key,
		Value: x,
	})
}

// AttachFile : will attach the contents of a file to the end of the definition
func (d *Definition) AttachFile(key, path string) (err error) {
	var str string
	path = "@{" + path + "}"
	if str, err = LoadFile(path); err == nil {
		d.data = append(d.data, yaml.MapItem{
			Key:   key,
			Value: str,
		})
	}
	return
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
