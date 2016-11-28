/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package manager

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/fatih/color"
)

// ********************* Login *******************

// Login does a login action against the api
func (m *Manager) Login(username string, password string) (token string, err error) {
	var t Token

	f := url.Values{}
	f.Add("username", username)
	f.Add("password", password)

	url := m.URL + "/auth"
	req, err := http.NewRequest("POST", url, strings.NewReader(f.Encode()))
	req.Form = f
	req.PostForm = f
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := m.client().Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New("The keypair user / password does not match any user on the database, please try again")
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		color.Red(err.Error())
	}

	err = json.Unmarshal(responseBody, &t)
	if err != nil {
		color.Red(err.Error())
	}

	token = t.Token

	return token, nil
}
