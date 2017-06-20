/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package manager

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/fatih/color"

	yaml "gopkg.in/yaml.v2"
)

// Manager manages all api communications
type Manager struct {
	URL string `json:"url"`
}

// Token holds the JWT token that is received when authenticating
type Token struct {
	Token string `json:"token"`
}

// Session ...
type Session struct {
	UserID  int    `json:"id"`
	GroupID int    `json:"group_id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"admin"`
}

var CONNECTIONREFUSED = errors.New("Connection refused")

func (m *Manager) client() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	return client
}

func (m *Manager) doRequest(url, method string, payload []byte, token string, contentType string) (string, *http.Response, error) {
	url = m.URL + url
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	resp, err := m.client().Do(req)

	if err != nil {
		return err.Error(), resp, err
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		color.Red(err.Error())
	}
	body := string(responseBody)

	if resp.StatusCode != 200 {
		return string(body), resp, errors.New(resp.Status)
	}
	return string(body), resp, nil
}

func (m *Manager) createClient(token string, name string) (string, error) {
	payload := []byte(`{"name":"` + name + `"}`)
	body, resp, err := m.doRequest("/api/groups/", "POST", payload, token, "")
	if err != nil {
		if resp == nil {
			return "", CONNECTIONREFUSED
		}
		return body, err
	}

	color.Green("SUCCESS: Group " + name + " created")

	var group struct {
		ID int `json:"id"`
	}
	err = json.Unmarshal([]byte(body), &group)
	if err != nil {
		return "", errors.New("ERROR: Couldn't read response from server")
	}
	return strconv.Itoa(group.ID), nil
}

// GetSession ..
func (m *Manager) GetSession(token string) (session Session, err error) {
	body, resp, err := m.doRequest("/api/session/", "GET", nil, token, "application/yaml")
	if err != nil {
		if resp == nil {
			return session, CONNECTIONREFUSED
		}
		return session, err
	}
	err = json.Unmarshal([]byte(body), &session)
	if err != nil {
		return session, err
	}
	return session, nil
}

// GetUUID ...
func (m *Manager) GetUUID(token string, payload []byte) string {
	id, err := buildServiceUUID(payload)
	if err != nil {
		log.Fatal(err)
	}
	body, _, _ := m.doRequest("/api/services/uuid/", "POST", []byte(`{"id":"`+id+`"}`), token, "")
	var dat map[string]interface{}
	err = json.Unmarshal([]byte(body), &dat)
	if err != nil {
		return ""
	}

	if str, ok := dat["uuid"].(string); ok {
		return str
	}
	return ""
}

func buildServiceUUID(payload []byte) (string, error) {
	var definition struct {
		Name       string
		Datacenter string
	}
	if err := yaml.Unmarshal(payload, &definition); err != nil {
		return "", err
	}
	return definition.Name + "-" + definition.Datacenter, nil
}
