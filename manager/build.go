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

// ListBuilds ...
func (m *Manager) ListBuilds(project, env, token string) (builds []model.Build, err error) {
	body, resp, err := m.doRequest("/api/projects/"+project+"/envs/"+env+"/builds/", "GET", []byte(""), token, "")
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

// BuildStatus ...
func (m *Manager) BuildStatus(token, project, env, index string) (build model.Build, err error) {
	buildID, err := m.BuildIDFromIndex(token, project, env, index)
	if err != nil {
		return build, err
	}

	return m.BuildStatusByID(token, project, env, buildID)
}

// BuildIDFromIndex ...
func (m *Manager) BuildIDFromIndex(token, project, env, index string) (string, error) {
	builds, _ := m.ListBuilds(project, env, token)
	num, _ := strconv.Atoi(index)
	if num < 1 || num > len(builds) {
		return "", errors.New("Invalid build ID")
	}
	num = len(builds) - num
	return builds[num].ID, nil
}

// BuildStatusByID ...
func (m *Manager) BuildStatusByID(token, project, env, buildID string) (build model.Build, err error) {
	body, resp, err := m.doRequest("/api/projects/"+project+"/envs/"+env+"/builds/"+buildID, "GET", []byte(""), token, "")
	if err != nil {
		if resp == nil {
			return build, ErrConnectionRefused
		}
		if resp.StatusCode == 403 {
			return build, errors.New("You don't have permissions to perform this action")
		}
		if resp.StatusCode == 404 {
			return build, errors.New("Specified build not found")
		}
		return build, err
	}
	if body == "null" {
		return build, errors.New("Unexpected endpoint response : " + string(body))
	}
	err = json.Unmarshal([]byte(body), &build)
	return build, err
}

// BuildDefinitionByID ...
func (m *Manager) BuildDefinitionByID(token, project, env, buildID string) ([]byte, error) {
	body, resp, err := m.doRequest("/api/projects/"+project+"/envs/"+env+"/builds/"+buildID+"/definition/", "GET", []byte(""), token, "")
	if err != nil {
		if resp == nil {
			return nil, ErrConnectionRefused
		}
		if resp.StatusCode == 403 {
			return nil, errors.New("You don't have permissions to perform this action")
		}
		if resp.StatusCode == 404 {
			return nil, errors.New("Specified build not found")
		}
		return nil, err
	}
	if body == "null" {
		return nil, errors.New("Unexpected endpoint response : " + string(body))
	}

	return []byte(body), err
}

// LatestBuildDefinition ...
func (m *Manager) LatestBuildDefinition(token, project, env string) ([]byte, error) {
	id, err := m.LatestBuildID(token, project, env)
	if err != nil {
		return nil, err
	}

	return m.BuildDefinitionByID(token, project, env, id)
}

// BuildDefinitionFromIndex ...
func (m *Manager) BuildDefinitionFromIndex(token, project, env, index string) ([]byte, error) {
	id, err := m.BuildIDFromIndex(token, project, env, index)
	if err != nil {
		return nil, err
	}

	return m.BuildDefinitionByID(token, project, env, id)
}

// LatestBuildID ...
func (m *Manager) LatestBuildID(token, project, env string) (string, error) {
	builds, err := m.ListBuilds(project, env, token)
	if err != nil {
		return "", err
	}

	if len(builds) < 1 {
		return "", errors.New("Specified build not found")
	}

	return builds[0].ID, nil
}

// LatestBuildStatus ...
func (m *Manager) LatestBuildStatus(token, project, env string) (build model.Build, err error) {
	id, err := m.LatestBuildID(token, project, env)
	if err != nil {
		return build, err
	}

	return m.BuildStatusByID(token, project, env, id)
}

// Apply : Applies a yaml to create / update a new env
func (m *Manager) Apply(token, path string, credentials map[string]interface{}, monit, dry bool) (string, error) {
	var d model.Definition

	payload, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.New("You should specify a valid template path or store an ernest.yml on the current folder")
	}

	err = d.Load(payload)
	if err != nil {
		return "", errors.New("Could not process definition yaml")
	}

	_, err = m.EnvStatus(token, d.Project, d.Name)
	if err != nil {
		m.CreateEnv(token, d.Name, d.Project, credentials)
	}

	// Load any imported files
	if err := d.LoadFileImports(); err != nil {
		return "", err
	}

	return m.ApplyEnv(d, token, credentials, monit, dry)
}

// Import : Imports an existing env
func (m *Manager) Import(token string, name string, project string, filters []string) (streamID string, err error) {
	a := model.Action{
		Type: "import",
	}

	a.Options.Filters = filters

	data, err := json.Marshal(a)
	if err != nil {
		return "", err
	}

	body, resp, rerr := m.doRequest("/api/projects/"+project+"/envs/"+name+"/actions/", "POST", data, token, "application/yaml")
	if resp == nil {
		return "", ErrConnectionRefused
	}
	if rerr != nil {
		return "", rerr
	}

	err = json.Unmarshal([]byte(body), &a)
	if err != nil {
		return "", errors.New(body)
	}

	return a.ResourceID, helper.Monitorize(m.URL, "/events", token, a.ResourceID)
}

// ApplyEnv : Applies a yaml to create / update a new env
func (m *Manager) ApplyEnv(d model.Definition, token string, credentials map[string]interface{}, monit, dry bool) (string, error) {
	payload, err := d.Save()

	if err != nil {
		return "", errors.New("Could not finalize definition yaml")
	}

	if dry {
		return m.dryApply(token, payload, d)
	}

	var response struct {
		ID      string `json:"id,omitempty"`
		Name    string `json:"name,omitempty"`
		Project string `json:"project,omitempty"`
		Message string `json:"message,omitempty"`
	}

	body, resp, rerr := m.doRequest("/api/projects/"+d.Project+"/envs/"+d.Name+"/builds/", "POST", payload, token, "application/yaml")
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

	if monit {
		err = helper.Monitorize(m.URL, "/events", token, response.ID)
		if err != nil {
			return response.ID, err
		}

		fmt.Println("================\nPlatform Details\n================\n ")
		var build model.Build

		build, err = m.BuildStatusByID(token, d.Project, d.Name, response.ID)
		if err != nil {
			return response.ID, err
		}

		view.PrintEnvInfo(&build)
	}

	return response.ID, nil
}

func (m *Manager) dryApply(token string, payload []byte, d model.Definition) (string, error) {
	var body string
	body, resp, err := m.doRequest("/api/projects/"+d.Project+"/envs/"+d.Name+"/builds/?dry=true", "POST", payload, token, "application/yaml")
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
