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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"
)

// Client makes RESTful calls to https://developer.paywithextend.com/#extend-api endpoints.
//
// Authentication (token retrieval and renewal) is automatically managed by the Client via the given
// username and testPassword.
//
// Close should be invoked upon exit to release Client resources.
type Client struct {
	// server is the URL of the Extend API.
	server string
	// email is the email address to use to signIn.
	email string
	// password corresponds to the email.
	password string
	// client is the http.Client to use to make HTTP requests.
	client *http.Client
	// aToken stores the current authentication token.
	aToken *atomic.Value
	// validity is the duration the token is valid for.
	validity time.Duration
	// closed is the channel used to stop the refreshToken goroutine.
	closed chan struct{}
}

// NewClient initializes and returns (a reference to) a Client.
func NewClient(server, email, password string) (*Client, error) {
	c := &Client{
		server:   server,
		email:    email,
		password: password,
		client:   &http.Client{Timeout: 10 * time.Second},
		aToken:   &atomic.Value{},
		validity: 10 * time.Minute,
		closed:   make(chan struct{}, 1),
	}
	var response LoginSignUpResponse
	err := c.signIn(&response)
	if err != nil {
		return nil, err
	}
	c.aToken.Store(response.Token)
	go c.refreshToken(response.RefreshToken)
	return c, nil
}

// token returns the Client.aToken value stored in the atomic.Value.
func (c *Client) token() string {
	switch val := c.aToken.Load().(type) {
	case string:
		return val
	default:
		return unauthenticated
	}
}

// unauthenticated represents an anonymous request (lack of auth token).
const unauthenticated = ""

// refreshToken renews the Client.token automatically according to the Client.validity duration.
func (c *Client) refreshToken(token string) {
	for {
		time.Sleep(c.validity)
		select {
		case <-c.closed:
			return
		default:
			var response LoginSignUpResponse
			err := c.renewAuth(token, &response)
			if err != nil {
				continue
			}
			c.aToken.Swap(response.Token)
			token = response.RefreshToken
		}
	}
}

// Close the Client.
func (c *Client) Close() {
	_ = c.signOut()
	c.closed <- struct{}{}
}

// signIn -> https://developer.paywithextend.com/#sign-in.
func (c *Client) signIn(response *LoginSignUpResponse) error {
	return do(
		c.client,
		http.MethodPost,
		c.server+"/signin",
		unauthenticated,
		LoginRequest{Email: c.email, Password: c.password},
		response)
}

// renewAuth -> https://developer.paywithextend.com/#renew-auth.
func (c *Client) renewAuth(token string, response *LoginSignUpResponse) error {
	return do(
		c.client,
		http.MethodPost,
		c.server+"/renewauth",
		unauthenticated,
		RefreshTokenLoginRequest{RefreshToken: token},
		response)
}

// signOut -> https://developer.paywithextend.com/#sign-out.
func (c *Client) signOut() error {
	response := make(map[string]any)
	return do(
		c.client,
		http.MethodDelete,
		c.server+"/signout",
		c.token(),
		LogoutRequest{},
		&response)
}

// ForgotPassword -> https://developer.paywithextend.com/#forgot-password.
func (c *Client) ForgotPassword(email string, response *Response) error {
	return do(
		c.client,
		http.MethodPost,
		c.server+"/forgot",
		c.token(),
		ForgotPasswordRequest{Email: email},
		response)
}

// GetUserVirtualCards -> https://developer.paywithextend.com/#get-user-virtual-cards.
func (c Client) GetUserVirtualCards() error {
	// TODO
	return nil
}

// do an HTTP request with the method, url, token, and body, using the client.
func do[In any, Out any](client *http.Client, method, url, token string, in In, out *Out) error {
	data, err := json.Marshal(in)
	if err != nil {
		return err
	}
	request, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/vnd.paywithextend.v2021-03-12+json")
	if token != unauthenticated {
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	data, err = io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if len(data) > 0 {
		err = json.Unmarshal(data, out)
		if err != nil {
			return err
		}
	}
	return nil
}
