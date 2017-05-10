/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package manager

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/ernestio/ernest-cli/helper"
	"github.com/ernestio/ernest-cli/model"
)

// ListUsers ...
func (m *Manager) ListUsers(token string) (users []model.User, err error) {
	body, resp, err := m.doRequest("/api/users/", "GET", []byte(""), token, "")
	if err != nil {
		if resp == nil {
			return nil, errors.New("Connection refused")
		}
		if resp.StatusCode == 400 {
			return users, errors.New("You're not allowed to perform this action, please log in")
		}
		if resp.StatusCode == 404 {
			return users, errors.New("Couldn't found any users")
		}
		return nil, err
	}
	err = json.Unmarshal([]byte(body), &users)
	if err != nil {
		return nil, err
	}
	return users, err
}

// GetUserByUsername : Gets a user by name
func (m *Manager) GetUserByUsername(token string, name string) (user model.User, err error) {
	users, err := m.ListUsers(token)
	for _, u := range users {
		if u.Username == name {
			return u, nil
		}
	}
	return user, errors.New("User not found")
}

// GetUser ...
func (m *Manager) GetUser(token string, userid string) (user model.User, err error) {
	body, resp, err := m.doRequest("/api/users/"+userid, "GET", nil, token, "application/yaml")
	if err != nil {
		if resp == nil {
			return user, errors.New("Connection refused")
		}
		return user, err
	}
	err = json.Unmarshal([]byte(body), &user)
	if err != nil {
		return user, err
	}
	return user, nil
}

// CreateUser ...
func (m *Manager) CreateUser(token string, name string, email string, user string, password string) error {
	payload := []byte(`{"group_id": 0, "username": "` + user + `", "email": "` + email + `", "password": "` + password + `"}`)
	body, resp, err := m.doRequest("/api/users/", "POST", payload, token, "")
	if err != nil {
		if resp == nil {
			return errors.New("Connection refused")
		}
		if resp.StatusCode != 200 {
			e := helper.ResponseMessage([]byte(body))
			if strings.Contains(e.Message, "invalid jwt") {
				return errors.New("You're not allowed to perform this action, please log in")
			}
			return errors.New(e.Message)
		}
		return err
	}
	return nil
}

// ChangePassword ...
func (m *Manager) ChangePassword(token string, userid int, username string, usergroup int, oldpassword string, newpassword string) error {
	payload := []byte(`{"id":` + strconv.Itoa(userid) + `, "username": "` + username + `", "group_id": ` + strconv.Itoa(usergroup) + `, "password": "` + newpassword + `", "oldpassword": "` + oldpassword + `"}`)
	body, resp, err := m.doRequest("/api/users/"+strconv.Itoa(userid), "PUT", payload, token, "application/yaml")
	if err != nil {
		if resp.StatusCode != 200 {
			e := helper.ResponseMessage([]byte(body))
			return errors.New(e.Message)
		}
		return err
	}
	return nil
}

// ChangePasswordByAdmin ...
func (m *Manager) ChangePasswordByAdmin(token string, userid int, username string, usergroup int, newpassword string) error {
	payload := []byte(`{"id":` + strconv.Itoa(userid) + `, "username": "` + username + `", "group_id": ` + strconv.Itoa(usergroup) + `, "password": "` + newpassword + `"}`)
	body, resp, err := m.doRequest("/api/users/"+strconv.Itoa(userid), "PUT", payload, token, "application/yaml")
	if err != nil {
		if resp.StatusCode != 200 {
			e := helper.ResponseMessage([]byte(body))
			return errors.New(e.Message)
		}
		return err
	}
	return nil
}
