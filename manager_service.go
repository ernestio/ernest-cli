/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"runtime"

	"github.com/fatih/color"
)

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
	body, resp, err := m.doRequest("/api/services/"+serviceName, "GET", []byte(""), token, "")
	if err != nil {
		if resp.StatusCode == 403 {
			return service, errors.New("You don't have permissions to perform this action")
		}
		if resp.StatusCode == 404 {
			return service, errors.New("Specified service not found")
		}
		return service, err
	}
	if body == "null" {
		return service, errors.New("Unexpected endpoint response : " + string(body))
	}
	err = json.Unmarshal([]byte(body), &service)
	if err != nil {
		return service, err
	}
	return service, nil
}

// ServiceBuildStatus ...
func (m *Manager) ServiceBuildStatus(token string, serviceName string, serviceID string) (service Service, err error) {
	body, resp, err := m.doRequest("/api/services/"+serviceName+"/builds/"+serviceID, "GET", []byte(""), token, "")
	if err != nil {
		if resp.StatusCode == 403 {
			return service, errors.New("You don't have permissions to perform this action")
		}
		if resp.StatusCode == 404 {
			return service, errors.New("Specified build not found")
		}
		return service, err
	}
	if body == "null" {
		return service, errors.New("Unexpected endpoint response : " + string(body))
	}
	err = json.Unmarshal([]byte(body), &service)
	if err != nil {
		return service, err
	}
	return service, nil
}

// ResetService ...
func (m *Manager) ResetService(name string, token string) error {
	s, err := m.ServiceStatus(token, name)
	if err != nil {
		return err
	}
	if s.Status != "in_progress" {
		return errors.New("The service '" + name + "' can't be resetted as is on status '" + s.Status + "'")
	}
	_, _, err = m.doRequest("/api/services/"+name+"/reset/", "POST", nil, token, "application/yaml")
	return err
}

// Destroy : Destroys an existing service
func (m *Manager) Destroy(token string, name string, monit bool) error {
	body, resp, err := m.doRequest("/api/services/"+name, "DELETE", nil, token, "application/yaml")
	if err != nil {
		if resp.StatusCode == 404 {
			return errors.New("Specified service name does not exist")
		}
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

// Apply : Applies a yaml to create / update a new service
func (m *Manager) Apply(token string, path string, monit bool) (string, error) {
	var d Definition

	payload, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.New("You should specify a valid template path or store an ernest.yml on the current folder")
	}

	err = d.Load(payload)
	if err != nil {
		return "", errors.New("Could not process definition yaml")
	}

	// Load any imported files
	if err := d.LoadFileImports(); err != nil {
		return "", err
	}

	payload, err = d.Save()

	if err != nil {
		return "", errors.New("Could not finalize definition yaml")
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
