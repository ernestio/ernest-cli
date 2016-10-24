/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"errors"
	"strconv"
)

// Datacenter : Representation of a datacenter
type Datacenter struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Type            string `json:"type"`
	GroupID         int    `json:"group_id"`
	GroupName       string `json:"group_name"`
	Region          string `json:"region"`
	VCloudURL       string `json:"vcloud_url"`
	VseURL          string `json:"vse_url"`
	ExternalNetwork string `json:"external_network"`
	Username        string `json:"username"`
}

// CreateVcloudDatacenter : Creates a VCloud datacenter
func (m *Manager) CreateVcloudDatacenter(token string, name string, rtype string, user string, password string, url string, network string, vseURL string) (string, error) {
	payload := []byte(`{"name": "` + name + `", "type":"` + rtype + `", "region": "", "username":"` + user + `", "password":"` + password + `", "external_network":"` + network + `", "vcloud_url":"` + url + `", "vse_url":"` + vseURL + `"}`)
	body, res, err := m.doRequest("/api/datacenters/", "POST", payload, token, "")
	if err != nil {
		if res.StatusCode == 409 {
			return "Datacenter '" + name + "' already exists, please specify a different name", err
		} else {
			return body, err
		}
	}
	return body, err
}

// CreateAWSDatacenter : Creates an AWS datacenter
func (m *Manager) CreateAWSDatacenter(token string, name string, rtype string, region string, awstoken string, awssecret string) (string, error) {
	payload := []byte(`{"name": "` + name + `", "type":"` + rtype + `", "region":"` + region + `", "username":"` + name + `", "token":"` + awstoken + `", "secret":"` + awssecret + `"}`)
	body, res, err := m.doRequest("/api/datacenters/", "POST", payload, token, "")
	if err != nil {
		if res.StatusCode == 409 {
			return "Datacenter '" + name + "' already exists, please specify a different name", err
		} else {
			return body, err
		}
	}
	return body, err
}

// ListDatacenters : Lists all datacenters on your account
func (m *Manager) ListDatacenters(token string) (datacenters []Datacenter, err error) {
	body, _, err := m.doRequest("/api/datacenters/", "GET", []byte(""), token, "")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(body), &datacenters)
	if err != nil {
		return nil, err
	}
	return datacenters, err
}

// DeleteDatacenter : Deletes an existing datacenter by its name
func (m *Manager) DeleteDatacenter(token string, name string) (err error) {
	g, err := m.getDatacenterByName(token, name)
	if err != nil {
		return errors.New("Datacenter '" + name + "' does not exist, please specify a different datacenter name")
	}
	id := strconv.Itoa(g.ID)

	body, res, err := m.doRequest("/api/datacenters/"+id, "DELETE", []byte(""), token, "")
	if err != nil {
		if res.StatusCode == 400 {
			return errors.New(body)
		}

		return err
	}
	return nil
}
