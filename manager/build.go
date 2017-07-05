/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package manager

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/ernestio/ernest-cli/model"
)

// ListBuilds : GET /api/services/<service>/builds
func (m *Manager) ListBuilds(name string, token string) (builds []model.Service, err error) {
	body, resp, err := m.doRequest("/api/services/"+name+"/builds/", "GET", []byte(""), token, "")
	if err != nil {
		if resp == nil {
			return nil, CONNECTIONREFUSED
		}
		return nil, err
	}
	err = json.Unmarshal([]byte(body), &builds)
	if err != nil {
		return nil, err
	}
	return builds, err
}

// ServiceBuildStatus : GET /api/services/<service>/builds/<build_id>
func (m *Manager) ServiceBuildStatus(token string, serviceName string, serviceID string) (service model.Service, err error) {
	builds, _ := m.ListBuilds(serviceName, token)
	num, _ := strconv.Atoi(serviceID)
	if num < 1 || num > len(builds) {
		return service, errors.New("Invalid build ID")
	}
	num = len(builds) - num
	serviceID = builds[num].ID

	body, resp, err := m.doRequest("/api/services/"+serviceName+"/builds/"+serviceID, "GET", []byte(""), token, "")
	if err != nil {
		if resp == nil {
			return service, CONNECTIONREFUSED
		} else if resp.StatusCode == 403 {
			return service, errors.New("You don't have permissions to perform this action")
		} else if resp.StatusCode == 404 {
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

// DelBuild : DELETE /api/services/<service>/builds/<build_id>
func (m *Manager) DelBuild(token, service, buildID string) error {
	body, resp, err := m.doRequest("/api/services/"+service+"/builds/"+buildID, "DELETE", []byte(""), token, "")
	if err != nil {
		if resp == nil {
			return CONNECTIONREFUSED
		}
		if resp.StatusCode == 403 {
			return errors.New("You don't have permissions to perform this action")
		} else if resp.StatusCode == 404 {
			return errors.New("Specified build not found")
		}
		return err
	}

	if body == "null" {
		return errors.New("Unexpected endpoint response : " + string(body))
	}

	return nil
}
