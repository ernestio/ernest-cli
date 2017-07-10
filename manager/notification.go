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

// CreateNotification : Creates a notification
func (m *Manager) CreateNotification(token string, name string, ntype string, config string) (string, error) {
	mPayload := make(map[string]string)
	mPayload["name"] = name
	mPayload["type"] = ntype
	mPayload["config"] = config
	payload, err := json.Marshal(mPayload)
	if err != nil {
		return "Internal error processing your input", err
	}

	// payload := []byte(`{"name": "` + name + `", "type":"` + ntype + `", "` + config + `"}`)
	body, res, err := m.doRequest("/api/notifications/", "POST", payload, token, "")
	if err != nil {
		if res == nil {
			return "", ErrConnectionRefused
		}
		if res.StatusCode == 409 {
			return "Notification '" + name + "' already exists, please specify a different name", err
		}
		if res.StatusCode == 403 {
			return body, errors.New(body)
		}
		return body, err
	}
	return body, err
}

// ListNotifications : Lists all notifications on your account
func (m *Manager) ListNotifications(token string) (notifications []model.Notification, err error) {
	body, res, err := m.doRequest("/api/notifications/", "GET", []byte(""), token, "")
	if err != nil {
		if res == nil {
			return nil, ErrConnectionRefused
		}
		if res.StatusCode == 403 {
			return nil, errors.New(body)
		}
		return nil, err
	}
	err = json.Unmarshal([]byte(body), &notifications)
	if err != nil {
		return nil, err
	}
	return notifications, err
}

// DeleteNotification : Deletes an existing notification by its name
func (m *Manager) DeleteNotification(token string, name string) (err error) {
	g, err := m.getNotificationByName(token, name)
	if err != nil {
		return errors.New("Notification '" + name + "' does not exist, please specify a different notification name")
	}
	id := strconv.Itoa(g.ID)

	body, res, err := m.doRequest("/api/notifications/"+id, "DELETE", []byte(""), token, "")
	if err != nil {
		if res == nil {
			return ErrConnectionRefused
		}
		if res.StatusCode == 400 {
			return errors.New(body)
		}
		if res.StatusCode == 403 {
			return errors.New(body)
		}

		return err
	}
	return nil
}

// UpdateNotification : updates notification details
func (m *Manager) UpdateNotification(token, name, config string) (err error) {
	g, err := m.getNotificationByName(token, name)
	if err != nil {
		return errors.New("Notification '" + name + "' does not exist, please specify a different notification name")
	}
	id := strconv.Itoa(g.ID)

	mPayload := make(map[string]string)
	mPayload["name"] = name
	mPayload["config"] = config
	payload, err := json.Marshal(mPayload)
	if err != nil {
		return errors.New("Internal error processing your input")
	}

	body, res, err := m.doRequest("/api/notifications/"+id, "PUT", payload, token, "")
	if err != nil {
		if res == nil {
			return ErrConnectionRefused
		}
		if res.StatusCode == 400 {
			return errors.New(body)
		}
		if res.StatusCode == 403 {
			return errors.New(body)
		}

		return err
	}

	return nil
}

// AddServiceToNotification : updates notification details
func (m *Manager) AddServiceToNotification(token, service, name string, delete bool) (err error) {
	g, err := m.getNotificationByName(token, name)
	if err != nil {
		return errors.New("Notification '" + name + "' does not exist, please specify a different notification name")
	}
	id := strconv.Itoa(g.ID)

	mPayload := make(map[string]string)
	mPayload["name"] = name
	mPayload["service"] = service
	payload, err := json.Marshal(mPayload)
	if err != nil {
		return errors.New("Internal error processing your input")
	}

	method := "POST"
	if delete {
		method = "DELETE"
	}
	body, res, err := m.doRequest("/api/notifications/"+id+"/"+service, method, payload, token, "")
	if err != nil {
		if res == nil {
			return ErrConnectionRefused
		}
		if res.StatusCode == 400 {
			return errors.New(body)
		}
		if res.StatusCode == 403 {
			return errors.New(body)
		}

		return err
	}

	return nil
}

func (m *Manager) getNotificationByName(token string, name string) (d model.Notification, err error) {
	notifications, err := m.ListNotifications(token)
	if err != nil {
		return d, err
	}

	for _, d := range notifications {
		if name == d.Name {
			return d, nil
		}
	}
	return d, errors.New("Notify does not exist")
}
