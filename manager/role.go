/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package manager

import (
	"encoding/json"
	"errors"
)

type roleObj struct {
	ID       string `json:"resource_id"`
	User     string `json:"user_id"`
	Role     string `json:"role"`
	Resource string `json:"resource_type"`
}

type respErr struct {
	Message string `json:"message"`
}

// SetRole : ...
func (m *Manager) SetRole(token, user, project, env, role string) (body string, err error) {
	return m.roleRequest(token, "POST", user, project, env, role)
}

// UnsetRole : ...
func (m *Manager) UnsetRole(token, user, project, env, role string) (body string, err error) {
	return m.roleRequest(token, "DELETE", user, project, env, role)
}

func (m *Manager) roleRequest(token, verb, user, project, env, role string) (body string, err error) {
	rType := "projects"
	rID := project
	if env != "" {
		rType = "environments"
		rID = project + "-" + env
	}
	r := roleObj{
		ID:       rID,
		User:     user,
		Resource: rType,
		Role:     role,
	}

	req, err := json.Marshal(r)
	if err != nil {
		return body, errors.New("Invalid input")
	}

	body, resp, err := m.doRequest("/api/roles/", verb, req, token, "")
	if err != nil {
		if resp.StatusCode == 403 {
			return body, errors.New("You're not allowed to perform this action, please contact the resource owner")
		}
		var rerr respErr
		err := json.Unmarshal([]byte(body), &rerr)
		if err == nil {
			body = rerr.Message
		}
	}

	return
}
