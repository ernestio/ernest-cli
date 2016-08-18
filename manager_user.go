/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"strconv"

	"github.com/fatih/color"
)

// User ...
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	GroupID  int    `json:"group_id"`
	IsAdmin  bool   `json:"admin"`
}

// ListUsers ...
func (m *Manager) ListUsers(token string) (users []User, err error) {
	body, _, err := m.doRequest("/api/users/", "GET", []byte(""), token, "")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(body), &users)
	if err != nil {
		return nil, err
	}
	return users, err
}

// GetUser ...
func (m *Manager) GetUser(token string, userid string) (user User, err error) {
	res, _, err := m.doRequest("/api/users/"+userid, "GET", nil, token, "application/yaml")
	if err != nil {
		return user, err
	}
	err = json.Unmarshal([]byte(res), &user)
	if err != nil {
		return user, err
	}
	return user, nil
}

// CreateUser ...
func (m *Manager) CreateUser(token string, name string, email string, user string, password string) error {
	payload := []byte(`{"group_id": 0, "username": "` + user + `", "email": "` + email + `", "password": "` + password + `"}`)
	_, _, err := m.doRequest("/api/users/", "POST", payload, token, "")
	if err != nil {
		return err
	}
	color.Green("SUCCESS: User " + user + " created")
	return nil
}

// ChangePassword ...
func (m *Manager) ChangePassword(token string, userid int, username string, usergroup int, oldpassword string, newpassword string) error {
	payload := []byte(`{"id":` + strconv.Itoa(userid) + `, "username": "` + username + `", "group_id": ` + strconv.Itoa(usergroup) + `, "password": "` + newpassword + `", "oldpassword": "` + oldpassword + `"}`)
	_, _, err := m.doRequest("/api/users/"+strconv.Itoa(userid), "PUT", payload, token, "application/yaml")
	if err != nil {
		return err
	}
	return nil
}

// ChangePasswordByAdmin ...
func (m *Manager) ChangePasswordByAdmin(token string, userid int, username string, usergroup int, newpassword string) error {
	payload := []byte(`{"id":` + strconv.Itoa(userid) + `, "username": "` + username + `", "group_id": ` + strconv.Itoa(usergroup) + `, "password": "` + newpassword + `"}`)
	_, _, err := m.doRequest("/api/users/"+string(userid), "PUT", payload, token, "application/yaml")
	if err != nil {
		return err
	}
	return nil
}
