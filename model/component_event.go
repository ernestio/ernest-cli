/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package model

// ComponentEvent represents an event of type component from a SSE stream
type ComponentEvent struct {
	ID         string           `json:"_component_id"`
	Subject    string           `json:"_subject"`
	Type       string           `json:"_component"`
	State      string           `json:"_state"`
	Action     string           `json:"_action"`
	Provider   string           `json:"_provider"`
	Name       string           `json:"name"`
	Error      string           `json:"error,omitempty"`
	Service    string           `json:"service,omitempty"`
	Components []ComponentEvent `json:"components,omitempty"`
}
