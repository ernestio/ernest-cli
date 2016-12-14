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

// CreateVcloudDatacenter : Creates a VCloud datacenter
func (m *Manager) CreateVcloudDatacenter(token string, name string, rtype string, user string, password string, url string, network string, vseURL string) (string, error) {
	payload := []byte(`{"name": "` + name + `", "type":"` + rtype + `", "region": "", "username":"` + user + `", "password":"` + password + `", "external_network":"` + network + `", "vcloud_url":"` + url + `", "vse_url":"` + vseURL + `"}`)
	body, res, err := m.doRequest("/api/datacenters/", "POST", payload, token, "")
	if err != nil {
		if res.StatusCode == 409 {
			return "Datacenter '" + name + "' already exists, please specify a different name", err
		}
		return body, err
	}
	return body, err
}

// CreateAWSDatacenter : Creates an AWS datacenter
func (m *Manager) CreateAWSDatacenter(token string, name string, rtype string, region string, awsAccessKeyID string, awsSecretAccessKey string) (string, error) {
	payload := []byte(`{"name": "` + name + `", "type":"` + rtype + `", "region":"` + region + `", "username":"` + name + `", "aws_access_key_id":"` + awsAccessKeyID + `", "aws_secret_access_key":"` + awsSecretAccessKey + `"}`)
	body, res, err := m.doRequest("/api/datacenters/", "POST", payload, token, "")
	if err != nil {
		if res.StatusCode == 409 {
			return "Datacenter '" + name + "' already exists, please specify a different name", err
		}
		return body, err
	}
	return body, err
}

// ListDatacenters : Lists all datacenters on your account
func (m *Manager) ListDatacenters(token string) (datacenters []model.Datacenter, err error) {
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

// UpdateVCloudDatacenter : updates vcloud datacenter details
func (m *Manager) UpdateVCloudDatacenter(token, name, user, password string) (err error) {
	g, err := m.getDatacenterByName(token, name)
	if err != nil {
		return errors.New("Datacenter '" + name + "' does not exist, please specify a different datacenter name")
	}
	id := strconv.Itoa(g.ID)

	payload := []byte(`{"username":"` + user + `", "password":"` + password + `"}`)
	body, res, err := m.doRequest("/api/datacenters/"+id, "PUT", payload, token, "")
	if err != nil {
		if res.StatusCode == 400 {
			return errors.New(body)
		}

		return err
	}

	return nil
}

// UpdateAWSDatacenter : updates awsdatacenter details
func (m *Manager) UpdateAWSDatacenter(token, name, awsAccessKeyID, awsSecretAccessKey string) (err error) {
	g, err := m.getDatacenterByName(token, name)
	if err != nil {
		return errors.New("Datacenter '" + name + "' does not exist, please specify a different datacenter name")
	}
	id := strconv.Itoa(g.ID)

	payload := []byte(`{"aws_access_key_id":"` + awsAccessKeyID + `", "aws_secret_access_key":"` + awsSecretAccessKey + `"}`)
	body, res, err := m.doRequest("/api/datacenters/"+id, "PUT", payload, token, "")
	if err != nil {
		if res.StatusCode == 400 {
			return errors.New(body)
		}

		return err
	}

	return nil
}
