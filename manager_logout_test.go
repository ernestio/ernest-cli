/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestSuccessLogout(t *testing.T) {
	convey.Convey("Given I do a success logout", t, func() {
		server := mockRequest("/session/", "DELETE", 200, ``)
		m := Manager{URL: server.URL}
		err := m.Logout("foo")
		convey.Convey("Then It does not fail", func() {
			convey.So(err, convey.ShouldBeNil)
		})
	})
}
