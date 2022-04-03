// Copyright 2022 c-fraser
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testEmail    = "_@gmail.com"
	testPassword = "P4$sW0rD"
	testToken    = "abc123DEF456ghi789JKL012"
)

func TestSignIn(t *testing.T) {
	server := newTestServer(t, http.MethodPost, "/signin", testToken)
	defer server.Close()

	client := newTestClient(t, server)
	defer client.Close()

	var response LoginSignUpResponse
	err := client.signIn(&response)
	if err != nil {
		t.Errorf("Failed to signin: %v", err)
	}
	if response.Token != testToken {
		t.Errorf("Unexpected repsonse token: %s", response.Token)
	}
}

func TestRenewAuth(t *testing.T) {
	server := newTestServer(t, http.MethodPost, "/renewauth", readTestdata(t, "login_signup_response.json"))
	defer server.Close()

	client := newTestClient(t, server)
	defer client.Close()

	var response LoginSignUpResponse
	err := client.renewAuth(testToken, &response)
	if err != nil {
		t.Errorf("Failed to renew auth: %v", err)
	}
	if response.RefreshToken != testToken {
		t.Errorf("Unexpected repsonse token: %s", response.RefreshToken)
	}
}

func TestSignOut(t *testing.T) {
	server := newTestServer(t, http.MethodDelete, "/signout", "")
	defer server.Close()

	client := newTestClient(t, server)
	defer client.Close()

	err := client.signOut()
	if err != nil {
		t.Errorf("Failed to signout: %v", err)
	}
}

func TestForgotPassword(t *testing.T) {
	server := newTestServer(t, http.MethodPost, "/forgot", readTestdata(t, "response.json"))
	defer server.Close()

	client := newTestClient(t, server)
	defer client.Close()

	var response Response
	err := client.ForgotPassword(testEmail, &response)
	if err != nil {
		t.Errorf("Failed to reset password: %v", err)
	}
	if response.Msg != "ok" {
		t.Errorf("Unexpected repsonse message: %s", response.Msg)
	}
}

func newTestClient(t *testing.T, s *httptest.Server) *Client {
	c, err := NewClient(s.URL, testEmail, testPassword)
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}
	return c
}

func newTestServer(t *testing.T, method, path, data string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if content := r.Header.Get("Content-Type"); content != "application/json" {
			t.Errorf("Unexpected 'Content-Type' header: %s", content)
		}
		if accept := r.Header.Get("Accept"); accept != "application/vnd.paywithextend.v2021-03-12+json" {
			t.Errorf("Unexpected 'Accept' header: %s", accept)
		}
		if r.URL.Path == "/signin" {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write(
				[]byte(
					fmt.Sprintf(
						`{"user": {},"token": "%s","refreshToken": "%s"}`,
						testToken,
						testToken)))
			if err != nil {
				t.Errorf("Failed to signin: %v", err)
			}
			return
		}
		if r.URL.Path == "/signout" {
			w.WriteHeader(http.StatusOK)
			return
		}
		if auth := r.Header.Get("Authorization"); r.URL.Path != "/renewauth" && auth != fmt.Sprintf("Bearer %s", testToken) {
			t.Errorf("Unexpected 'Authorization' header: %s", auth)
		}
		if r.Method != method {
			t.Errorf("Unexpected HTTP method: %s", r.Method)
		}
		if r.URL.Path != path {
			t.Errorf("Unexpected path %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(data))
		if err != nil {
			t.Errorf("Failed to write response data: %v", err)
		}
	}))
}

func readTestdata(t *testing.T, filename string) string {
	data, err := ioutil.ReadFile("testdata/" + filename)
	if err != nil {
		t.Fatalf("Failed to read testdata file: %s", filename)
	}
	return string(data)
}
