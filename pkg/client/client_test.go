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
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

const (
	testEmail         = "_@gmail.com"
	testPassword      = "P4$sW0rD"
	testToken         = "abc123DEF456ghi789JKL012"
	testVirtualCardId = "vc_1234"
	testTransactionId = "txn_1234"
)

func TestSignIn(t *testing.T) {
	server := newTestServer(t, http.MethodPost, "/signin", "", testToken)
	defer server.Close()

	client := newTestClient(t, server)
	defer client.Close()

	response, err := client.signIn()
	if err != nil {
		t.Errorf("Failed to signin: %v", err)
	}
	if response.Token != testToken {
		t.Errorf("Unexpected repsonse token: %s", response.Token)
	}
}

func TestRenewAuth(t *testing.T) {
	server := newTestServer(
		t,
		http.MethodPost,
		"/renewauth",
		readTestdata(t, "refresh_token_request.json"),
		readTestdata(t, "login_signup_response.json"))
	defer server.Close()

	client := newTestClient(t, server)
	defer client.Close()

	response, err := client.renewAuth(testToken)
	if err != nil {
		t.Errorf("Failed to renew auth: %v", err)
	}
	if response.RefreshToken != testToken {
		t.Errorf("Unexpected repsonse token: %s", response.RefreshToken)
	}
}

func TestSignOut(t *testing.T) {
	server := newTestServer(t, http.MethodDelete, "/signout", "", "")
	defer server.Close()

	client := newTestClient(t, server)
	defer client.Close()

	err := client.signOut()
	if err != nil {
		t.Errorf("Failed to signout: %v", err)
	}
}

func TestForgotPassword(t *testing.T) {
	server := newTestServer(
		t,
		http.MethodPost,
		"/forgot",
		readTestdata(t, "forgot_password_request.json"),
		readTestdata(t, "response.json"))
	defer server.Close()

	client := newTestClient(t, server)
	defer client.Close()

	response, err := client.ForgotPassword(testEmail)
	if err != nil {
		t.Errorf("Failed to reset password: %v", err)
	}
	if response.Msg != "ok" {
		t.Errorf("Unexpected repsonse message: %s", response.Msg)
	}
}

func TestGetUserVirtualCards(t *testing.T) {
	server := newTestServer(
		t,
		http.MethodGet,
		"/virtualcards",
		readTestdata(t, "virtual_card_pageable_request.json"),
		readTestdata(t, "virtual_cards_response.json"))
	defer server.Close()

	client := newTestClient(t, server)
	defer client.Close()

	response, err := client.GetUserVirtualCards(&VirtualCardPageableRequest{
		Count:              0,
		Page:               0,
		SortField:          "string",
		SortDirection:      "string",
		Cardholder:         "string",
		Recipient:          "string",
		CardholderOrViewer: "string",
		CreditCardID:       "string",
		Status:             "string",
		Statuses:           []string{"string"},
		Issued:             true,
		PendingRequest:     true,
		Search:             "string",
		WithPermission:     "string",
	})
	if err != nil {
		t.Errorf("Failed to get user virtual cards: %v", err)
	}
	if len(response.VirtualCards) != 1 {
		t.Fatalf("Unexpected virtual cards response: %v", response.VirtualCards)
	}
	if vc := response.VirtualCards[0]; vc.ID != testVirtualCardId {
		t.Errorf("Unexpected virtual card ID: %v", vc)
	}
}

func TestGetVirtualCard(t *testing.T) {
	server := newTestServer(
		t,
		http.MethodGet,
		"/virtualcards/"+testVirtualCardId,
		"",
		readTestdata(t, "virtual_card_response.json"))
	defer server.Close()

	client := newTestClient(t, server)
	defer client.Close()

	response, err := client.GetVirtualCard(testVirtualCardId)
	if err != nil {
		t.Errorf("Failed to get virtual card: %v", err)
	}
	if response.VirtualCard.ID != testVirtualCardId {
		t.Errorf("Unexpected virtual card ID: %v", response.VirtualCard)
	}
}

func TestGetVirtualCardTransactions(t *testing.T) {
	server := newTestServer(
		t,
		http.MethodGet,
		"/virtualcards/"+testVirtualCardId+"/transactions?count=25&status=CLEARED",
		"",
		readTestdata(t, "transactions_response.json"))
	defer server.Close()

	client := newTestClient(t, server)
	defer client.Close()

	response, err := client.GetVirtualCardTransactions(testVirtualCardId, 25, "", "", "CLEARED")
	if err != nil {
		t.Errorf("Failed to get virtual card transactions: %v", err)
	}
	if tx := response.Transactions[0]; tx.ID != testTransactionId {
		t.Errorf("Unexpected transaction ID: %v", response.Transactions)
	}
}

func TestCreateVirtualCard(t *testing.T) {
	server := newTestServer(
		t,
		http.MethodPost,
		"/virtualcards",
		"",
		readTestdata(t, "virtual_card_response.json"))
	defer server.Close()

	client := newTestClient(t, server)
	defer client.Close()

	var request CreateVirtualCardRequest
	err := json.Unmarshal([]byte(readTestdata(t, "create_virtual_card_request.json")), &request)
	if err != nil {
		t.Errorf("Failed to initialize request: %v", err)
	}
	response, err := client.CreateVirtualCard(&request)
	if err != nil {
		t.Errorf("Failed to create virtual card: %v", err)
	}
	if response.VirtualCard.ID != testVirtualCardId {
		t.Errorf("Unexpected virtual card ID: %v", response.VirtualCard)
	}
}

func TestUpdateVirtualCard(t *testing.T) {
	server := newTestServer(
		t,
		http.MethodPut,
		"/virtualcards/"+testVirtualCardId,
		"",
		readTestdata(t, "virtual_card_response.json"))
	defer server.Close()

	client := newTestClient(t, server)
	defer client.Close()

	var request UpdateVirtualCardRequest
	err := json.Unmarshal([]byte(readTestdata(t, "update_virtual_card_request.json")), &request)
	if err != nil {
		t.Errorf("Failed to initialize request: %v", err)
	}
	response, err := client.UpdateVirtualCard(testVirtualCardId, &request)
	if err != nil {
		t.Errorf("Failed to update virtual card: %v", err)
	}
	if response.VirtualCard.ID != testVirtualCardId {
		t.Errorf("Unexpected virtual card ID: %v", response.VirtualCard)
	}
}

func TestCancelVirtualCard(t *testing.T) {
	server := newTestServer(
		t,
		http.MethodPut,
		"/virtualcards/"+testVirtualCardId+"/cancel",
		"",
		readTestdata(t, "virtual_card_response.json"))
	defer server.Close()

	client := newTestClient(t, server)
	defer client.Close()

	response, err := client.CancelVirtualCard(testVirtualCardId)
	if err != nil {
		t.Errorf("Failed to cancel virtual card: %v", err)
	}
	if response.VirtualCard.ID != testVirtualCardId {
		t.Errorf("Unexpected virtual card ID: %v", response.VirtualCard)
	}
}

func TestRejectVirtualCard(t *testing.T) {
	server := newTestServer(
		t,
		http.MethodPut,
		"/virtualcards/"+testVirtualCardId+"/reject",
		"",
		readTestdata(t, "virtual_card_response.json"))
	defer server.Close()

	client := newTestClient(t, server)
	defer client.Close()

	response, err := client.RejectVirtualCard(testVirtualCardId)
	if err != nil {
		t.Errorf("Failed to reject virtual card: %v", err)
	}
	if response.VirtualCard.ID != testVirtualCardId {
		t.Errorf("Unexpected virtual card ID: %v", response.VirtualCard)
	}
}

func newTestClient(t *testing.T, s *httptest.Server) *Client {
	c, err := NewClient(s.URL, testEmail, testPassword)
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}
	return c
}

func newTestServer(t *testing.T, method, path, data, response string) *httptest.Server {
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
		if query := "?" + r.URL.RawQuery; r.URL.RawQuery != "" {
			if r.URL.Path+query != path {
				t.Errorf("Unexpected path %s", r.URL.Path+query)
			}
		} else {
			if r.URL.Path != path {
				t.Errorf("Unexpected path %s", r.URL.Path)
			}
		}
		if data != "" {
			all, err := io.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				t.Errorf("Failed to read request data: %v", err)
			}
			var actual interface{}
			var expected interface{}
			err = json.Unmarshal(all, &actual)
			err = json.Unmarshal([]byte(data), &expected)
			if err != nil {
				t.Errorf("Failed to unmarshal request data: %v", err)
			}
			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("Unexpected request data, expected: %v, actual: %v", expected, actual)
			}
		}
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(response))
		if err != nil {
			t.Errorf("Failed to write response: %v", err)
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
