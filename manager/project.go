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

// CreateVcloudProject : Creates a VCloud project
func (m *Manager) CreateVcloudProject(token string, name string, rtype string, user string, password string, url string, network string, vseURL string) (string, error) {
	payload := []byte(`{"name": "` + name + `", "type":"` + rtype + `", "credentials":{"region": "", "username":"` + user + `", "password":"` + password + `", "external_network":"` + network + `", "vcloud_url":"` + url + `", "vse_url":"` + vseURL + `"}}`)
	body, res, err := m.doRequest("/api/projects/", "POST", payload, token, "")
	if err != nil {
		if res == nil {
			return "", ErrConnectionRefused
		}
		if res.StatusCode == 409 {
			return "Project '" + name + "' already exists, please specify a different name", err
		}
		return body, err
	}
	return body, err
}

// CreateAWSProject : Creates an AWS project
func (m *Manager) CreateAWSProject(token string, name string, rtype string, region string, awsAccessKeyID string, awsSecretAccessKey string) (string, error) {
	payload := []byte(`{"name": "` + name + `", "type":"` + rtype + `", "credentials":{"region":"` + region + `", "username":"` + name + `", "aws_access_key_id":"` + awsAccessKeyID + `", "aws_secret_access_key":"` + awsSecretAccessKey + `"}}`)
	body, res, err := m.doRequest("/api/projects/", "POST", payload, token, "")
	if err != nil {
		if res == nil {
			return "", ErrConnectionRefused
		}
		if res.StatusCode == 409 {
			return "Project '" + name + "' already exists, please specify a different name", err
		}
		return body, err
	}
	return body, err
}

// CreateAzureProject : Creates an Azure project
func (m *Manager) CreateAzureProject(token, name, rtype, region, subscriptionID, clientID, clientSecret, tenantID, environment string) (string, error) {
	payload := []byte(`{"name": "` + name + `", "type":"` + rtype + `", "credentials": {"region":"` + region + `", "username":"` + name + `", "azure_subscription_id":"` + subscriptionID + `", "azure_client_id":"` + clientID + `", "azure_client_secret": "` + clientSecret + `", "azure_tenant_id": "` + tenantID + `", "azure_environment": "` + environment + `"}}`)
	body, res, err := m.doRequest("/api/projects/", "POST", payload, token, "")
	if err != nil {
		if res == nil {
			return "", ErrConnectionRefused
		}
		if res.StatusCode == 409 {
			return "Project '" + name + "' already exists, please specify a different name", err
		}
		return body, err
	}
	return body, err
}

// ListProjects : Lists all projects on your account
func (m *Manager) ListProjects(token string) (projects []model.Project, err error) {
	body, res, err := m.doRequest("/api/projects/", "GET", []byte(""), token, "")
	if err != nil {
		if res == nil {
			return nil, ErrConnectionRefused
		}
		return nil, err
	}
	err = json.Unmarshal([]byte(body), &projects)
	if err != nil {
		return nil, err
	}
	return projects, err
}

// DeleteProject : Deletes an existing project by its name
func (m *Manager) DeleteProject(token string, name string) (err error) {
	g, err := m.getProjectByName(token, name)
	if err != nil {
		return errors.New("Project '" + name + "' does not exist, please specify a different project name")
	}
	id := strconv.Itoa(g.ID)

	body, res, err := m.doRequest("/api/projects/"+id, "DELETE", []byte(""), token, "")
	if err != nil {
		if res == nil {
			return ErrConnectionRefused
		}
		if res.StatusCode == 400 {
			return errors.New(body)
		}

		return err
	}
	return nil
}

// UpdateVCloudProject : updates vcloud project details
func (m *Manager) UpdateVCloudProject(token, name, user, password string) (err error) {
	g, err := m.getProjectByName(token, name)
	if err != nil {
		return errors.New("Project '" + name + "' does not exist, please specify a different project name")
	}
	id := strconv.Itoa(g.ID)

	payload := []byte(`{"username":"` + user + `", "password":"` + password + `"}`)
	body, res, err := m.doRequest("/api/projects/"+id, "PUT", payload, token, "")
	if err != nil {
		if res == nil {
			return ErrConnectionRefused
		}
		if res.StatusCode == 400 {
			return errors.New(body)
		}

		return err
	}

	return nil
}

// UpdateAWSProject : updates awsproject details
func (m *Manager) UpdateAWSProject(token, name, awsAccessKeyID, awsSecretAccessKey string) (err error) {
	g, err := m.getProjectByName(token, name)
	if err != nil {
		return errors.New("Project '" + name + "' does not exist, please specify a different project name")
	}
	id := strconv.Itoa(g.ID)

	payload := []byte(`{"aws_access_key_id":"` + awsAccessKeyID + `", "aws_secret_access_key":"` + awsSecretAccessKey + `"}`)
	body, res, err := m.doRequest("/api/projects/"+id, "PUT", payload, token, "")
	if err != nil {
		if res == nil {
			return ErrConnectionRefused
		}
		if res.StatusCode == 400 {
			return errors.New(body)
		}

		return err
	}

	return nil
}

// UpdateAzureProject : updates awsproject details
func (m *Manager) UpdateAzureProject(token, name, subscriptionID, clientID, clientSecret, tenantID, environment string) (err error) {
	g, err := m.getProjectByName(token, name)
	if err != nil {
		return errors.New("Project '" + name + "' does not exist, please specify a different project name")
	}
	id := strconv.Itoa(g.ID)

	payload := []byte(`{"azure_subscription_id":"` + subscriptionID + `", "azure_client_id":"` + clientID + `", "azure_client_secret": "` + clientSecret + `", "azure_tenant_id": "` + tenantID + `", "azure_environment": "` + environment + `"}`)
	body, res, err := m.doRequest("/api/projects/"+id, "PUT", payload, token, "")
	if err != nil {
		if res == nil {
			return ErrConnectionRefused
		}
		if res.StatusCode == 400 {
			return errors.New(body)
		}

		return err
	}

	return nil
}

func (m *Manager) getProjectByName(token string, name string) (project model.Project, err error) {
	body, res, err := m.doRequest("/api/projects/"+name, "GET", []byte(""), token, "")
	if err != nil {
		if res == nil {
			return project, ErrConnectionRefused
		}
		return
	}
	err = json.Unmarshal([]byte(body), &project)
	return
}

// InfoProject : updates awsproject details
func (m *Manager) InfoProject(token, name string) (p model.Project, err error) {
	return m.getProjectByName(token, name)
}
