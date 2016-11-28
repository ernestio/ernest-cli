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

func (m *Manager) getDatacenterByName(token string, name string) (d model.Datacenter, err error) {
	datacenters, err := m.ListDatacenters(token)
	if err != nil {
		return d, err
	}

	for _, d := range datacenters {
		if name == d.Name {
			return d, nil
		}
	}
	return d, errors.New("Datanceter does not exist")
}

func (m *Manager) getGroupByName(token string, name string) (g model.Group, err error) {
	groups, err := m.ListGroups(token)
	if err != nil {
		return g, err
	}

	for _, g := range groups {
		if name == g.Name {
			return g, nil
		}
	}
	return g, errors.New("Group does not exist")
}

// DeleteGroup : Deletes a group through ernest api
func (m *Manager) DeleteGroup(token string, group string) error {
	g, err := m.getGroupByName(token, group)
	if err != nil {
		return errors.New("Group '" + group + "' does not exist, please specify a different group name")
	}
	id := strconv.Itoa(g.ID)

	_, resp, err := m.doRequest("/api/groups/"+id, "DELETE", []byte(""), token, "")
	if resp.StatusCode == 403 {
		return errors.New("You're not allowed to perform this action, please log in")
	}
	if resp.StatusCode == 409 {
		return errors.New("Group '" + group + "' already exists, please specify a different group name")
	}

	return err
}

// CreateGroup : Creates a group through ernest api
func (m *Manager) CreateGroup(token string, group string) error {
	payload := []byte(`{"name": "` + group + `"}`)
	_, resp, err := m.doRequest("/api/groups/", "POST", payload, token, "")
	if resp.StatusCode == 403 {
		return errors.New("You're not allowed to perform this action, please log in")
	}
	if resp.StatusCode == 409 {
		return errors.New("Group '" + group + "' already exists, please specify a different group name")
	}
	if resp.StatusCode == 400 {
		return errors.New("Group '" + group + "' already exists, please specify a different group name")
	}

	return err
}

// ListGroups ...
func (m *Manager) ListGroups(token string) (groups []model.Group, err error) {
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
func (m *Manager) GroupAddUser(token string, user string, group string) (err error) {
	var u model.User
	if u, err = m.GetUserByUsername(token, user); err != nil {
		return errors.New("User '" + user + "' does not exist")
	}

	g, err := m.getGroupByName(token, group)
	if err != nil {
		return errors.New("Group '" + group + "' does not exist")
	}

	if g.ID == u.GroupID && g.ID != 0 {
		return errors.New("User '" + user + "' already belongs to '" + group + "' group")
	}

	payload := []byte(`{"username": "` + user + `", "group": "` + group + `"}`)
	_, resp, err := m.doRequest("/api/groups/"+group+"/users/", "POST", payload, token, "")
	if resp.StatusCode == 403 {
		return errors.New("You're not allowed to perform this action, please log in")
	}

	return err
}

// GroupRemoveUser ...
func (m *Manager) GroupRemoveUser(token string, user string, group string) (err error) {
	var u model.User
	if u, err = m.GetUserByUsername(token, user); err != nil {
		return errors.New("User '" + user + "' does not exist")
	}

	g, err := m.getGroupByName(token, group)
	if err != nil {
		return errors.New("Group '" + group + "' does not exist")
	}

	if g.ID != u.GroupID {
		return errors.New("User 'tmp_user' does not belong to 'tmp_group' group")
	}
	userid := strconv.Itoa(u.ID)
	groupid := strconv.Itoa(g.ID)

	_, resp, err := m.doRequest("/api/groups/"+groupid+"/users/"+userid, "DELETE", []byte(""), token, "")
	if resp.StatusCode == 403 {
		return errors.New("You're not allowed to perform this action, please log in")
	}

	return err
}

// GroupAddDatacenter ...
func (m *Manager) GroupAddDatacenter(token string, datacenter string, group string) (err error) {
	var d model.Datacenter
	if d, err = m.getDatacenterByName(token, datacenter); err != nil {
		return errors.New("Datacenter '" + datacenter + "' does not exist")
	}

	g, err := m.getGroupByName(token, group)
	if err != nil {
		return errors.New("Group '" + group + "' does not exist")
	}

	if g.ID == d.GroupID && g.ID != 0 {
		return errors.New("Datacenter '" + datacenter + "' already belongs to '" + group + "' group")
	}

	datacenterid := strconv.Itoa(d.ID)
	groupid := strconv.Itoa(g.ID)
	payload := []byte(`{"datacenterid": "` + datacenterid + `", "groupid": "` + groupid + `"}`)
	_, resp, err := m.doRequest("/api/groups/"+groupid+"/datacenters/", "POST", payload, token, "")
	if resp.StatusCode == 403 {
		return errors.New("You're not allowed to perform this action, please log in")
	}

	return err
}

// GroupRemoveDatacenter ...
func (m *Manager) GroupRemoveDatacenter(token string, datacenter string, group string) (err error) {
	var d model.Datacenter
	if d, err = m.getDatacenterByName(token, datacenter); err != nil {
		return errors.New("Datacenter '" + datacenter + "' does not exist")
	}

	g, err := m.getGroupByName(token, group)
	if err != nil {
		return errors.New("Group '" + group + "' does not exist")
	}

	if g.ID != d.GroupID {
		return errors.New("User 'tmp_user' does not belong to 'tmp_group' group")
	}
	datacenterid := strconv.Itoa(d.ID)
	groupid := strconv.Itoa(g.ID)

	_, resp, err := m.doRequest("/api/groups/"+groupid+"/datacenters/"+datacenterid, "DELETE", []byte(""), token, "")
	if resp.StatusCode == 403 {
		return errors.New("You're not allowed to perform this action, please log in")
	}

	return err
}
