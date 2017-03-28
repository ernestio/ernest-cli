/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package model

// Message represents an incomming websocket message
type Message struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
	Level   string `json:"level"`
	User    string `json:"user"`
}
