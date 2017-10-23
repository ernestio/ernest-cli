/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package manager

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ernestio/ernest-cli/helper"
	"github.com/ernestio/ernest-cli/model"
	"github.com/ernestio/ernest-cli/view"
	"github.com/fatih/color"
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
	req := []byte(`{"type": "reset"}`)
	_, resp, err := m.doRequest("/api/projects/"+project+"/envs/"+env+"/actions/", "POST", req, token, "application/json")
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
	payload, err := m.BuildDefinitionFromIndex(token, project, env, buildID)
	if err != nil {
		return "", err
	}

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

	return m.ApplyEnv(d, token, nil, true, false)
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
	e := model.Env{
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

	switch resp.StatusCode {
	case 404:
		return errors.New("Specified environment does not exist")
	case 403:
		return errors.New("You don't have permissions to perform this action, please login as a resource owner")
	case 401:
		return errors.New("Invalid session, please log in")
	}

	return rerr
}

// CreateEnv : Creates a new empty environmnet
func (m *Manager) CreateEnv(token, name, project string, credentials map[string]interface{}) error {
	e := model.Env{
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

	switch resp.StatusCode {
	case 404:
		return errors.New("Specified project does not exist")
	case 403:
		return errors.New("You don't have permissions to perform this action, please login as a resource owner")
	}

	return rerr
}

// SyncEnv : Sync's an environmnet
func (m *Manager) SyncEnv(token, name, project string) error {
	a := model.Action{
		Type: "sync",
	}

	data, err := json.Marshal(a)
	if err != nil {
		return err
	}

	body, resp, rerr := m.doRequest("/api/projects/"+project+"/envs/"+name+"/actions/", "POST", data, token, "application/json")
	if resp == nil {
		return ErrConnectionRefused
	}
	if rerr != nil {
		return rerr
	}

	err = json.Unmarshal([]byte(body), &a)
	if err != nil {
		return errors.New(body)
	}

	stream := helper.OpenStream(m.URL, "/events", token, a.ResourceID)

	for {
		var m map[string]interface{}

		msg, ok := <-stream
		if !ok {
			return errors.New("could not monitor sync progress")
		}

		if msg.Data == nil {
			continue
		}

		err = json.Unmarshal(bytes.Trim(msg.Data, "\x00"), &m)
		if err != nil {
			return err
		}

		if m["_subject"].(string) == "build.import.done" {
			break
		}
	}

	time.Sleep(time.Second * 1)

	b, err := m.BuildStatusByID(token, project, name, a.ResourceID)
	if err != nil {
		return err
	}

	switch b.Status {
	case "done":
		color.Green("No changes detected")
	case "awaiting_resolution":
		color.Red("Changes detected")
		fmt.Println("")

		body, err := m.BuildMappingChanges(token, project, name, a.ResourceID)
		if err != nil {
			return err
		}

		view.SyncChanges(body)
	case "errored":
		color.Red("Sync failed!")
	}

	return nil
}

// ResolveEnv : Resolves an issue with an environment
func (m *Manager) ResolveEnv(token, name, project, resolution string) error {
	a := model.Action{
		Type: "resolve",
	}

	a.Options.Resolution = resolution

	data, err := json.Marshal(a)
	if err != nil {
		return err
	}

	body, resp, rerr := m.doRequest("/api/projects/"+project+"/envs/"+name+"/actions/", "POST", data, token, "application/json")
	if resp == nil {
		return ErrConnectionRefused
	}
	if rerr != nil {
		return rerr
	}

	err = json.Unmarshal([]byte(body), &a)
	if err != nil {
		return errors.New(body)
	}

	if a.Status == "done" {
		color.Green("Sync successfully resolved!")
		return nil
	}

	return helper.Monitorize(m.URL, "/events", token, a.ResourceID)
}
