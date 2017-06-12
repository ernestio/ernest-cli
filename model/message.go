/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package model

// Message represents an incomming websocket message
type Message struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
	Level   string `json:"level"`
}

type ServiceNew struct {
	ID      string `json:"id"`
	Subject string `json:"_subject"`
	//  Components []Component `json:"components"`
	Changes []ComponentNew `json:"changes"`
}

type ComponentNew struct {
	//  Name     string `json:"name"`
	//  Provider string `json:"provider"`
	//  State    string `json:"_state"`
	ID       string `json:"_component_id"`
	Subject  string `json:"_subject"`
	Type     string `json:"_component"`
	State    string `json:"_state"`
	Action   string `json:"_action"`
	Provider string `json:"_provider"`
	Name     string `json:"name"`
	Error    string `json:"error,omitempty"`
	Service  string `json:"service,omitempty"`
}
