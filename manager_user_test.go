/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestListUsers(t *testing.T) {
	convey.Convey("Given I get all users", t, func() {
		server := mockRequest("/api/users/", "GET", 200, `[]`)
		m := Manager{URL: server.URL}
		_, err := m.ListUsers("token")
		convey.Convey("Then It does not fail", func() {
			convey.So(err, convey.ShouldBeNil)
		})
	})
}

func TestCreateUser(t *testing.T) {
	t.Skip()
	convey.Convey("Given I create a client", t, func() {
		server := mockRequest("/api/groups/", "POST", 200, `{}`)
		m := Manager{URL: server.URL}
		err := m.CreateUser("name", "email", "user", "password", "adminuser", "adminpassword")
		convey.Convey("Then It does not fail", func() {
			convey.So(err, convey.ShouldBeNil)
		})
	})
}
