/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package manager

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/ernestio/ernest-cli/helper"
	"github.com/fatih/color"
)

// ********************* Login *******************

// Login does a login action against the api
func (m *Manager) Login(username, password, verificationCode string) (token string, err error) {
	var t Token

	url := m.URL + "/auth"
	body := []byte(`{"username": "` + username + `", "password": "` + password + `", "verification_code": "` + verificationCode + `"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Ernest/"+m.Version)

	resp, err := m.client().Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		color.Red(err.Error())
	}

	if resp.StatusCode != 200 {
		e := helper.ResponseMessage(body)
		return "", errors.New(e.Message)
	}

	err = json.Unmarshal(body, &t)
	if err != nil {
		color.Red(err.Error())
	}

	token = t.Token

	return token, nil
}
