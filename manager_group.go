/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
)

// Group ...
type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// DeleteGroup ...
func (m *Manager) DeleteGroup(token string, group string) error {
	_, _, err := m.doRequest("/api/groups/"+group, "DELETE", []byte(""), token, "")
	if err != nil {
		fmt.Println(err)
		return err
	}
	color.Green("SUCCESS: Group " + group + " deleted")
	return nil
}

// CreateGroup ...
func (m *Manager) CreateGroup(token string, group string) error {
	payload := []byte(`{"name": "` + group + `"}`)
	_, _, err := m.doRequest("/api/groups/", "POST", payload, token, "")
	if err != nil {
		fmt.Println(err)
		return err
	}
	color.Green("SUCCESS: Group " + group + " created")
	return nil
}

// ListGroups ...
func (m *Manager) ListGroups(token string) (groups []Group, err error) {
	body, _, err := m.doRequest("/api/groups/", "GET", []byte(""), token, "")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(body), &groups)
	if err != nil {
		return nil, err
	}
	return groups, err
}

// GroupAddUser : relates a group and a user
func (m *Manager) GroupAddUser(token string, user string, group string) error {
	payload := []byte(`{"username": "` + user + `", "group": "` + group + `"}`)
	_, _, err := m.doRequest("/api/groups/"+group+"/users/", "POST", payload, token, "")
	if err != nil {
		fmt.Println(err)
		return err
	}
	color.Green("SUCCESS: Added user " + user + " to group " + group)
	return nil
}

// GroupRemoveUser ...
func (m *Manager) GroupRemoveUser(token string, userid string, groupid string) error {
	_, _, err := m.doRequest("/api/groups/"+groupid+"/users/"+userid, "DELETE", []byte(""), token, "")
	if err != nil {
		fmt.Println(err)
		return err
	}
	color.Green("SUCCESS: Removed user " + userid + " from group " + groupid)
	return nil
}

// GroupAddDatacenter ...
func (m *Manager) GroupAddDatacenter(token string, datacenterid string, groupid string) error {
	payload := []byte(`{"datacenterid": "` + datacenterid + `", "groupid": "` + groupid + `"}`)
	_, _, err := m.doRequest("/api/groups/"+groupid+"/datacenters/", "POST", payload, token, "")
	if err != nil {
		fmt.Println(err)
		return err
	}
	color.Green("SUCCESS: Added datacenter " + datacenterid + " to group " + groupid)
	return nil
}

// GroupRemoveDatacenter ...
func (m *Manager) GroupRemoveDatacenter(token string, datacenterid string, groupid string) error {
	_, _, err := m.doRequest("/api/groups/"+groupid+"/datacenters/"+datacenterid, "DELETE", []byte(""), token, "")
	if err != nil {
		fmt.Println(err)
		return err
	}
	color.Green("SUCCESS: Removed datacenter " + datacenterid + " from group " + groupid)
	return nil
}
