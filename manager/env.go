/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/ernestio/ernest-cli/helper"
	"github.com/ernestio/ernest-cli/model"
	"github.com/ernestio/ernest-cli/view"
)

// ListEnvs ...
func (m *Manager) ListEnvs(token string) (services []model.Service, err error) {
	body, resp, err := m.doRequest("/api/environments/", "GET", []byte(""), token, "")
	if err != nil {
		if resp == nil {
			return nil, ErrConnectionRefused
		}
		return nil, err
	}
	err = json.Unmarshal([]byte(body), &services)
	if err != nil {
		return nil, err
	}
	return services, err
}

// ListBuilds ...
func (m *Manager) ListBuilds(name string, token string) (builds []model.Service, err error) {
	body, resp, err := m.doRequest("/api/environments/"+name+"/builds/", "GET", []byte(""), token, "")
	if err != nil {
		if resp == nil {
			return nil, ErrConnectionRefused
		}
		return nil, err
	}
	err = json.Unmarshal([]byte(body), &builds)
	if err != nil {
		return nil, err
	}
	return builds, err
}

// EnvStatus ...
func (m *Manager) EnvStatus(token string, serviceName string) (service model.Service, err error) {
	body, resp, err := m.doRequest("/api/environments/"+serviceName, "GET", []byte(""), token, "")
	if err != nil {
		if resp == nil {
			return service, ErrConnectionRefused
		}
		if resp.StatusCode == 403 {
			return service, errors.New("You don't have permissions to perform this action")
		}
		if resp.StatusCode == 404 {
			return service, errors.New("Specified environment name does not exist")
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

// EnvBuildStatus ...
func (m *Manager) EnvBuildStatus(token string, serviceName string, serviceID string) (service model.Service, err error) {
	builds, _ := m.ListBuilds(serviceName, token)
	num, _ := strconv.Atoi(serviceID)
	if num < 1 || num > len(builds) {
		return service, errors.New("Invalid build ID")
	}
	num = len(builds) - num
	serviceID = builds[num].ID

	body, resp, err := m.doRequest("/api/environments/"+serviceName+"/builds/"+serviceID, "GET", []byte(""), token, "")
	if err != nil {
		if resp == nil {
			return service, ErrConnectionRefused
		}
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

// ResetEnv ...
func (m *Manager) ResetEnv(name string, token string) error {
	s, err := m.EnvStatus(token, name)
	if err != nil {
		return err
	}
	if s.Status != "in_progress" {
		return errors.New("The environment '" + name + "' cannot be reset as its status is '" + s.Status + "'")
	}
	_, resp, err := m.doRequest("/api/environments/"+name+"/reset/", "POST", nil, token, "application/yaml")
	if err != nil {
		if resp == nil {
			return ErrConnectionRefused
		}
	}
	return err
}

// RevertEnv reverts a service to a previous known state using a build ID
func (m *Manager) RevertEnv(name string, buildID string, token string, dry bool) (string, error) {
	// get requested manifest
	s, err := m.EnvBuildStatus(token, name, buildID)
	if err != nil {
		return "", err
	}
	payload := []byte(s.Definition)

	// apply requested manifest
	var d model.Definition

	err = d.Load(payload)
	if err != nil {
		return "", errors.New("Could not process definition yaml")
	}

	payload, err = d.Save()
	if err != nil {
		return "", errors.New("Could not finalize definition yaml")
	}

	if dry {
		return m.dryApply(token, payload)
	}

	var response struct {
		ID      string `json:"id,omitempty"`
		Name    string `json:"name,omitempty"`
		Message string `json:"message,omitempty"`
	}

	body, resp, rerr := m.doRequest("/api/environments/", "POST", payload, token, "application/yaml")
	if resp == nil {
		return "", ErrConnectionRefused
	}

	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		return "", errors.New(body)
	}

	if rerr != nil {
		return "", errors.New(response.Message)
	}

	err = helper.Monitorize(m.URL, "/events", token, response.ID)
	if err != nil {
		return "", err
	}

	fmt.Println("================\nPlatform Details\n================\n ")
	var srv model.Service

	srv, err = m.EnvStatus(token, name)
	if err != nil {
		return response.ID, err
	}

	view.PrintEnvInfo(&srv)

	return response.ID, nil
}

// Destroy : Destroys an existing service
func (m *Manager) Destroy(token string, name string, monit bool) error {
	s, err := m.EnvStatus(token, name)
	if err != nil {
		return err
	}
	if s.Status == "in_progress" {
		return errors.New("The service " + name + " cannot be destroyed as it is currently '" + s.Status + "'")
	}

	body, resp, err := m.doRequest("/api/environments/"+name, "DELETE", nil, token, "application/yaml")
	if err != nil {
		if resp == nil {
			return ErrConnectionRefused
		}
		if resp.StatusCode == 404 {
			return errors.New("Specified environment name does not exist")
		}
		return err
	}

	var res map[string]interface{}
	err = json.Unmarshal([]byte(body), &res)
	if err != nil {
		return err
	}

	if id, ok := res["id"].(string); ok {
		err = helper.Monitorize(m.URL, "/events", token, id)
		if err != nil {
			return err
		}
	} else {
		return errors.New("could not read response")
	}

	return nil
}

// ForceDestroy : Destroys an existing service by forcing it
func (m *Manager) ForceDestroy(token, name string) error {
	_, resp, err := m.doRequest("/api/environments/"+name+"/force/", "DELETE", nil, token, "application/yaml")
	if err != nil {
		if resp == nil {
			return ErrConnectionRefused
		}
		if resp.StatusCode == 404 {
			return errors.New("Specified environment name does not exist")
		}
		return err
	}

	return nil
}

// Apply : Applies a yaml to create / update a new service
func (m *Manager) Apply(token string, path string, monit, dry bool) (string, error) {
	var d model.Definition

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

	if dry {
		return m.dryApply(token, payload)
	}

	var response struct {
		ID      string `json:"id,omitempty"`
		Name    string `json:"name,omitempty"`
		Message string `json:"message,omitempty"`
	}

	body, resp, rerr := m.doRequest("/api/environments/", "POST", payload, token, "application/yaml")
	if resp == nil {
		return "", ErrConnectionRefused
	}

	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		return "", errors.New(body)
	}

	if rerr != nil {
		return "", errors.New(response.Message)
	}

	err = helper.Monitorize(m.URL, "/events", token, response.ID)
	if err != nil {
		return response.ID, err
	}

	fmt.Println("================\nPlatform Details\n================\n ")
	var srv model.Service

	srv, err = m.EnvStatus(token, response.Name)
	if err != nil {
		return response.ID, err
	}

	view.PrintEnvInfo(&srv)

	return response.ID, nil
}

// Import : Imports an existing service
func (m *Manager) Import(token string, name string, project string, filters []string) (streamID string, err error) {
	var s struct {
		Name          string   `json:"name"`
		Project       string   `json:"datacenter"`
		ImportFilters []string `json:"import_filters,omitempty"`
	}
	s.Name = name
	s.Project = project
	s.ImportFilters = filters
	payload, err := json.Marshal(s)
	if err != nil {
		return "", errors.New("Invalid name or project")
	}

	var response struct {
		ID      string `json:"id,omitempty"`
		Name    string `json:"name,omitempty"`
		Message string `json:"message,omitempty"`
	}

	body, resp, rerr := m.doRequest("/api/environments/import/", "POST", payload, token, "application/yaml")
	if resp == nil {
		return "", ErrConnectionRefused
	}

	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		return "", errors.New(body)
	}

	if rerr != nil {
		return "", errors.New(response.Message)
	}

	err = helper.Monitorize(m.URL, "/events", token, response.ID)

	return response.ID, err
}

func (m *Manager) dryApply(token string, payload []byte) (string, error) {
	var body string
	body, resp, err := m.doRequest("/api/environments/?dry=true", "POST", payload, token, "application/yaml")
	if err != nil {
		if resp == nil {
			return "", ErrConnectionRefused
		}
		var internalError struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal([]byte(body), &internalError); err != nil {
			return "", errors.New(body)
		}
		return "", errors.New(internalError.Message)
	}
	view.EnvDry(body)
	return "", nil
}
