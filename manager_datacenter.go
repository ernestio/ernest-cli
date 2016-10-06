/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
)

// Datacenter ...
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

// CreateVcloudDatacenter ...
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

// CreateAWSDatacenter ...
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

// ListDatacenters ...
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
