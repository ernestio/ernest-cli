/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package manager

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ernestio/api-gateway/models"
	"github.com/ernestio/ernest-cli/helper"
	"github.com/ernestio/ernest-cli/model"
	"github.com/ernestio/ernest-cli/view"
)

// ListEnvs ...
func (m *Manager) ListEnvs(token string) (envs []model.Env, err error) {
	body, resp, err := m.doRequest("/api/envs/", "GET", []byte(""), token, "")
	if err != nil {
		if resp == nil {
			return nil, ErrConnectionRefused
		}
		return nil, err
	}
	err = json.Unmarshal([]byte(body), &envs)
	if err != nil {
		return nil, err
	}
	return envs, err
}

// EnvStatus ...
func (m *Manager) EnvStatus(token, project, env string) (environment model.Env, err error) {
	body, resp, err := m.doRequest("/api/projects/"+project+"/envs/"+env, "GET", []byte(""), token, "")
	if err != nil {
		if resp == nil {
			return environment, ErrConnectionRefused
		}
		if resp.StatusCode == 403 {
			return environment, errors.New("You don't have permissions to perform this action")
		}
		if resp.StatusCode == 404 {
			return environment, errors.New("Specified environment name does not exist")
		}
		return environment, err
	}
	if body == "null" {
		return environment, errors.New("Unexpected endpoint response : " + string(body))
	}
	err = json.Unmarshal([]byte(body), &environment)

	return environment, err
}

// ResetEnv ...
func (m *Manager) ResetEnv(project, env, token string) error {
	e, err := m.EnvStatus(token, project, env)
	if err != nil {
		return err
	}
	if e.Status != "in_progress" {
		return errors.New("The environment '" + project + " / " + env + "' cannot be reset as its status is '" + e.Status + "'")
	}
	_, resp, err := m.doRequest("/api/projects/"+project+"/envs/"+env+"/actions/reset/", "POST", nil, token, "application/yaml")
	if err != nil {
		if resp == nil {
			return ErrConnectionRefused
		}
	}
	return err
}

// RevertEnv reverts a env to a previous known state using a build ID
func (m *Manager) RevertEnv(project, env, buildID, token string, dry bool) (string, error) {
	// get requested manifest
	b, err := m.BuildStatus(token, project, env, buildID)
	if err != nil {
		return "", err
	}
	payload := []byte(b.Definition)

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
		return m.dryApply(token, payload, d)
	}

	var response struct {
		ID      string `json:"id,omitempty"`
		Name    string `json:"name,omitempty"`
		Message string `json:"message,omitempty"`
	}

	body, resp, rerr := m.doRequest("/api/projects/"+d.Project+"/envs/", "POST", payload, token, "application/yaml")
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
	var build model.Build

	build, err = m.BuildStatus(token, project, env, response.ID)
	if err != nil {
		return response.ID, err
	}

	view.PrintEnvInfo(&build)

	return response.ID, nil
}

// Destroy : Destroys an existing env
func (m *Manager) Destroy(token, project, env string, monit bool) error {
	s, err := m.EnvStatus(token, project, env)
	if err != nil {
		return err
	}
	if s.Status == "in_progress" {
		return errors.New("The environment " + env + " cannot be destroyed as it is currently '" + s.Status + "'")
	}

	body, resp, err := m.doRequest("/api/projects/"+project+"/envs/"+env, "DELETE", nil, token, "application/yaml")
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

// ForceDestroy : Destroys an existing env by forcing it
func (m *Manager) ForceDestroy(token, project, env string) error {
	_, resp, err := m.doRequest("/api/projects/"+project+"/envs/"+env+"/actions/force/", "DELETE", nil, token, "application/yaml")
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

// UpdateEnv : Updates credentials on a specific environment
func (m *Manager) UpdateEnv(token, name, project string, credentials map[string]interface{}) error {
	e := models.Env{
		Name:        name,
		Credentials: credentials,
	}

	payload, err := json.Marshal(e)
	if err != nil {
		return err
	}

	_, resp, rerr := m.doRequest("/api/projects/"+project+"/envs/"+name, "PUT", payload, token, "application/json")
	if resp == nil {
		return ErrConnectionRefused
	}

	return rerr
}

// CreateEnv : Creates a new empty environmnet
func (m *Manager) CreateEnv(token, name, project string, credentials map[string]interface{}) error {
	e := models.Env{
		Name:        name,
		Credentials: credentials,
	}

	payload, err := json.Marshal(e)
	if err != nil {
		return err
	}

	_, resp, rerr := m.doRequest("/api/projects/"+project+"/envs/", "POST", payload, token, "application/json")
	if resp == nil {
		return ErrConnectionRefused
	}

	return rerr
}

// Import : Imports an existing env
func (m *Manager) Import(token string, name string, project string, filters []string) (streamID string, err error) {
	var s struct {
		Name          string   `json:"name"`
		Project       string   `json:"project"`
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

	body, resp, rerr := m.doRequest("/api/projects/"+project+"/envs/"+name+"/actions/import/", "POST", payload, token, "application/yaml")
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
