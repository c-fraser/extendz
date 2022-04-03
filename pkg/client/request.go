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

// LoginRequest -> https://developer.paywithextend.com/#tocS_LoginRequest.
type LoginRequest struct {
	Email    string `json:"testEmail"`
	Password string `json:"testPassword"`
}

// LogoutRequest -> https://developer.paywithextend.com/#tocS_LogoutRequest.
type LogoutRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// RefreshTokenLoginRequest -> https://developer.paywithextend.com/#tocS_RefreshTokenLoginRequest.
type RefreshTokenLoginRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// ForgotPasswordRequest -> https://developer.paywithextend.com/#tocS_ForgotPasswordRequest.
type ForgotPasswordRequest struct {
	Email string `json:"testEmail"`
}

// VirtualCardPageableRequest -> https://developer.paywithextend.com/#tocS_VirtualCardPageableRequest.
type VirtualCardPageableRequest struct {
	Count              int      `json:"count"`
	Page               int      `json:"page"`
	SortField          string   `json:"sortField"`
	SortDirection      string   `json:"sortDirection"`
	Cardholder         string   `json:"cardholder"`
	Recipient          string   `json:"recipient"`
	CardholderOrViewer string   `json:"cardholderOrViewer"`
	CreditCardID       string   `json:"creditCardId"`
	Status             string   `json:"status"`
	Statuses           []string `json:"statuses"`
	Issued             bool     `json:"issued"`
	PendingRequest     bool     `json:"pendingRequest"`
	Search             string   `json:"search"`
	WithPermission     string   `json:"withPermission"`
}
