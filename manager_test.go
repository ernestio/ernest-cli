/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
)

func mockRequest(route string, method string, status int, output string) *httptest.Server {
	r := mux.NewRouter()
	r.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		s := output
		if s == "" {
			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)
			s = buf.String()
		}
		w.WriteHeader(status)
		w.Header().Set("X-Auth-Token", "")
		fmt.Fprint(w, s)
	}).Methods(method)

	return httptest.NewServer(r)
}
