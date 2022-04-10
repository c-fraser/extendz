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

// Response -> https://developer.paywithextend.com/#tocS_SimpleResponse.
type Response struct {
	Msg string `json:"msg"`
}

// LoginSignUpResponse -> https://developer.paywithextend.com/#tocS_LoginSignUpResponse.
type LoginSignUpResponse struct {
	User         User   `json:"user"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

// VirtualCardsResponse -> https://developer.paywithextend.com/#tocS_VirtualCardsResponse.
type VirtualCardsResponse struct {
	Pagination   Pagination    `json:"pagination"`
	VirtualCards []VirtualCard `json:"virtualCards"`
}

// VirtualCardResponse -> https://developer.paywithextend.com/#tocS_VirtualCardResponse.
type VirtualCardResponse struct {
	VirtualCard VirtualCard `json:"virtualCard"`
}

// TransactionsResponse -> https://developer.paywithextend.com/#tocS_TransactionsResponse.
type TransactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
}
