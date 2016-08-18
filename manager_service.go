/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"runtime"
	"time"

	"github.com/fatih/color"
)

// Service ...
type Service struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Datacenter int       `json:"datacenter_id"`
	Version    time.Time `json:"version"`
	Status     string    `json:"status"`
	Definition string    `json:"definition"`
	Result     string    `json:"result"`
	Endpoint   string    `json:"endpoint"`
}

// ListServices ...
func (m *Manager) ListServices(token string) (services []Service, err error) {
	body, _, err := m.doRequest("/api/services/", "GET", []byte(""), token, "")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(body), &services)
	if err != nil {
		return nil, err
	}
	return services, err
}

// ListBuilds ...
func (m *Manager) ListBuilds(name string, token string) (builds []Service, err error) {
	body, _, err := m.doRequest("/api/services/"+name+"/builds/", "GET", []byte(""), token, "")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(body), &builds)
	if err != nil {
		return nil, err
	}
	return builds, err
}

// ServiceStatus ...
func (m *Manager) ServiceStatus(token string, serviceName string) (service Service, err error) {
	body, _, err := m.doRequest("/api/services/"+serviceName, "GET", []byte(""), token, "")
	if err != nil {
		return service, err
	}
	err = json.Unmarshal([]byte(body), &service)
	if err != nil {
		return service, err
	}
	return service, nil
}

// ServiceBuildStatus ...
func (m *Manager) ServiceBuildStatus(token string, serviceName string, serviceID string) (service Service, err error) {
	body, _, err := m.doRequest("/api/services/"+serviceName+"/builds/"+serviceID, "GET", []byte(""), token, "")
	if err != nil {
		return service, err
	}
	err = json.Unmarshal([]byte(body), &service)
	if err != nil {
		return service, err
	}
	return service, nil
}

// ResetService ...
func (m *Manager) ResetService(name string, token string) error {
	_, _, err := m.doRequest("/api/services/"+name+"/reset/", "POST", nil, token, "application/yaml")
	return err
}

// Destroy ...
func (m *Manager) Destroy(token string, name string, monit bool) error {
	body, _, err := m.doRequest("/api/services/"+name, "DELETE", nil, token, "application/yaml")
	if err != nil {
		return err
	}

	var res map[string]interface{}
	err = json.Unmarshal([]byte(body), &res)
	if err != nil {
		return err
	}

	if monit == true {
		if str, ok := res["stream_id"].(string); ok {
			Monitorize(m.URL, token, str)
			runtime.Goexit()
		}
	}

	return nil
}

// Apply ...
func (m *Manager) Apply(token string, path string, monit bool) (string, error) {
	payload, err := ioutil.ReadFile(path)
	if err != nil {
		color.Red(err.Error())
		return "", nil
	}

	color.Green("Environment creation requested")
	println("Ernest will show you all output from your requested service creation")
	println("You can cancel at any moment with Ctrl+C, even the service is still being created, you won't have any output")

	streamID := m.GetUUID(token, payload)
	if streamID == "" {
		color.Red("Please log in")
		return "", nil
	}

	if monit == true {
		go Monitorize(m.URL, token, streamID)
	} else {
		println("Additionally you can trace your service on ernest monitor tool with id: " + streamID)
	}

	if body, _, err := m.doRequest("/api/services/", "POST", payload, token, "application/yaml"); err != nil {
		return "", errors.New(body)
	}
	if monit == true {
		runtime.Goexit()
	}
	return streamID, nil
}
