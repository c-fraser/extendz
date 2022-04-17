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
	"net/url"
	"strconv"
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
	response, err := c.signIn()
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

// empty is used to denote the absence of a request payload.
var empty *any

// refreshToken renews the Client.token automatically according to the Client.validity duration.
func (c *Client) refreshToken(token string) {
	for {
		time.Sleep(c.validity)
		select {
		case <-c.closed:
			return
		default:
			response, err := c.renewAuth(token)
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
func (c *Client) signIn() (*LoginSignUpResponse, error) {
	return do[LoginRequest, LoginSignUpResponse](
		c.client,
		http.MethodPost,
		c.server+"/signin",
		unauthenticated,
		&LoginRequest{Email: c.email, Password: c.password})
}

// renewAuth -> https://developer.paywithextend.com/#renew-auth.
func (c *Client) renewAuth(token string) (*LoginSignUpResponse, error) {
	return do[RefreshTokenLoginRequest, LoginSignUpResponse](
		c.client,
		http.MethodPost,
		c.server+"/renewauth",
		unauthenticated,
		&RefreshTokenLoginRequest{RefreshToken: token})
}

// signOut -> https://developer.paywithextend.com/#sign-out.
func (c *Client) signOut() error {
	_, err := do[LogoutRequest, any](
		c.client,
		http.MethodDelete,
		c.server+"/signout",
		c.token(),
		&LogoutRequest{})
	return err
}

// ForgotPassword -> https://developer.paywithextend.com/#forgot-password.
func (c *Client) ForgotPassword(email string) (*Response, error) {
	return do[ForgotPasswordRequest, Response](
		c.client,
		http.MethodPost,
		c.server+"/forgot",
		c.token(),
		&ForgotPasswordRequest{Email: email})
}

// GetUserVirtualCards -> https://developer.paywithextend.com/#get-user-virtual-cards.
func (c *Client) GetUserVirtualCards(request *VirtualCardPageableRequest) (*VirtualCardsResponse, error) {
	return do[VirtualCardPageableRequest, VirtualCardsResponse](
		c.client,
		http.MethodGet,
		c.server+"/virtualcards",
		c.token(),
		request)
}

// GetVirtualCard -> https://developer.paywithextend.com/#get-virtual-card.
func (c *Client) GetVirtualCard(id string) (*VirtualCardResponse, error) {
	return do[any, VirtualCardResponse](
		c.client,
		http.MethodGet,
		c.server+"/virtualcards/"+id,
		c.token(),
		empty)
}

// GetVirtualCardTransactions -> https://developer.paywithextend.com/#get-virtual-card-transactions.
func (c *Client) GetVirtualCardTransactions(id string, count int, before, after, status string) (*TransactionsResponse, error) {
	v := url.Values{}
	if count > 0 && count <= 500 {
		v.Add("count", strconv.Itoa(count))
	}
	if before != "" {
		v.Add("before", before)
	}
	if after != "" {
		v.Add("after", after)
	}
	if status != "" {
		v.Add("status", status)
	}
	u := c.server + "/virtualcards/" + id + "/transactions"
	if len(v) > 0 {
		u += "?" + v.Encode()
	}
	return do[any, TransactionsResponse](c.client, http.MethodGet, u, c.token(), empty)
}

// CreateVirtualCard -> https://developer.paywithextend.com/#create-virtual-card.
func (c *Client) CreateVirtualCard(request *CreateVirtualCardRequest) (*VirtualCardResponse, error) {
	return do[CreateVirtualCardRequest, VirtualCardResponse](
		c.client,
		http.MethodPost,
		c.server+"/virtualcards",
		c.token(),
		request)
}

// UpdateVirtualCard -> https://developer.paywithextend.com/#update-virtual-card.
func (c *Client) UpdateVirtualCard(id string, request *UpdateVirtualCardRequest) (*VirtualCardResponse, error) {
	return do[UpdateVirtualCardRequest, VirtualCardResponse](
		c.client,
		http.MethodPut,
		c.server+"/virtualcards/"+id,
		c.token(),
		request)
}

// CancelVirtualCard -> https://developer.paywithextend.com/#cancel-virtual-card.
func (c *Client) CancelVirtualCard(id string) (*VirtualCardResponse, error) {
	return do[any, VirtualCardResponse](
		c.client,
		http.MethodPut,
		c.server+"/virtualcards/"+id+"/cancel",
		c.token(),
		empty)
}

// RejectVirtualCard -> https://developer.paywithextend.com/#reject-virtual-card.
func (c *Client) RejectVirtualCard(id string) (*VirtualCardResponse, error) {
	return do[any, VirtualCardResponse](
		c.client,
		http.MethodPut,
		c.server+"/virtualcards/"+id+"/reject",
		c.token(),
		empty)
}

// do an HTTP request with the method, url, token, and body, using the client.
func do[rq any, rs any](client *http.Client, method, url, token string, in *rq) (*rs, error) {
	var reader io.Reader
	if in != nil {
		data, err := json.Marshal(in)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(data)
	}
	request, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/vnd.paywithextend.v2021-03-12+json")
	if token != unauthenticated {
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var out rs
	if len(data) > 0 {
		err = json.Unmarshal(data, &out)
		if err != nil {
			return nil, err
		}
	}
	return &out, err
}
