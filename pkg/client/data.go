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

// User -> https://developer.paywithextend.com/#tocS_User.
type User struct {
	ID                string            `json:"id"`
	FirstName         string            `json:"firstName"`
	LastName          string            `json:"lastName"`
	Email             string            `json:"testEmail"`
	Phone             string            `json:"phone"`
	PhoneIsoCountry   string            `json:"phoneIsoCountry"`
	AvatarType        string            `json:"avatarType"`
	AvatarURL         string            `json:"avatarUrl"`
	CreatedAt         string            `json:"createdAt"`
	UpdatedAt         string            `json:"updatedAt"`
	Currency          string            `json:"currency"`
	Locale            string            `json:"locale"`
	Timezone          string            `json:"timezone"`
	Verified          bool              `json:"verified"`
	HasExpensifyLink  bool              `json:"hasExpensifyLink"`
	QuickbooksTokenID string            `json:"quickbooksTokenId"`
	EmployeeID        string            `json:"employeeId"`
	IssuerSanctions   []IssuerSanctions `json:"issuerSanctions"`
	OrganizationID    string            `json:"organizationId"`
	OrganizationRole  string            `json:"organizationRole"`
}

// IssuerSanctions -> https://developer.paywithextend.com/#tocS_IssuerSanctions.
type IssuerSanctions struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// Pagination -> https://developer.paywithextend.com/#tocS_PaginationView.
type Pagination struct {
	Page          int `json:"page"`
	PageItemCount int `json:"pageItemCount"`
	TotalItems    int `json:"totalItems"`
	NumberOfPages int `json:"numberOfPages"`
}

// VirtualCard -> https://developer.paywithextend.com/#tocS_VirtualCard.
type VirtualCard struct {
	ID                    string              `json:"id"`
	Status                string              `json:"status"`
	RecipientID           string              `json:"recipientId"`
	Recipient             User                `json:"recipient"`
	CardholderID          string              `json:"cardholderId"`
	Cardholder            User                `json:"cardholder"`
	CardImage             CardImage           `json:"cardImage"`
	DisplayName           string              `json:"displayName"`
	Expires               string              `json:"expires"`
	Currency              string              `json:"currency"`
	LimitCents            int                 `json:"limitCents"`
	BalanceCents          int                 `json:"balanceCents"`
	SpentCents            int                 `json:"spentCents"`
	LifetimeSpentCents    int                 `json:"lifetimeSpentCents"`
	AwaitingBudget        bool                `json:"awaitingBudget"`
	Last4                 string              `json:"last4"`
	NumberFormat          string              `json:"numberFormat"`
	ValidFrom             string              `json:"validFrom"`
	ValidTo               string              `json:"validTo"`
	InactiveSince         string              `json:"inactiveSince"`
	Timezone              string              `json:"timezone"`
	CreditCardID          string              `json:"creditCardId"`
	Recurs                bool                `json:"recurs"`
	Recurrence            Recurrence          `json:"recurrence"`
	Pending               VirtualCardRevision `json:"pending"`
	Notes                 string              `json:"notes"`
	CreatedAt             string              `json:"createdAt"`
	UpdatedAt             string              `json:"updatedAt"`
	Address               Address             `json:"address"`
	Direct                bool                `json:"direct"`
	Features              VirtualCardFeature  `json:"features"`
	ActiveUntil           string              `json:"activeUntil"`
	MinTransactionCents   int                 `json:"minTransactionCents"`
	MaxTransactionCents   int                 `json:"maxTransactionCents"`
	MaxTransactionCount   int                 `json:"maxTransactionCount"`
	TokenReferenceIds     string              `json:"tokenReferenceIds"`
	Network               string              `json:"network"`
	CompanyName           string              `json:"companyName"`
	CreditCardDisplayName string              `json:"creditCardDisplayName"`
	Issuer                string              `json:"issuer"`
	ValidMccRanges        []MccRange          `json:"validMccRanges"`
}

// CardImage -> https://developer.paywithextend.com/#tocS_CardImage.
type CardImage struct {
	ID                  string            `json:"id"`
	ContentType         string            `json:"contentType"`
	Urls                map[string]string `json:"urls"`
	TextColorRGBA       string            `json:"textColorRGBA"`
	HasTextShadow       bool              `json:"hasTextShadow"`
	ShadowTextColorRGBA string            `json:"shadowTextColorRGBA"`
}

// Recurrence -> https://developer.paywithextend.com/#tocS_Recurrence.
type Recurrence struct {
	ID               string `json:"id"`
	BalanceCents     int    `json:"balanceCents"`
	Period           string `json:"period"`
	Interval         int    `json:"interval"`
	Terminator       string `json:"terminator"`
	Count            int    `json:"count"`
	Until            string `json:"until"`
	ByWeekDay        int    `json:"byWeekDay"`
	ByMonthDay       int    `json:"byMonthDay"`
	ByYearDay        int    `json:"byYearDay"`
	PrevRecurrenceAt string `json:"prevRecurrenceAt"`
	NextRecurrenceAt string `json:"nextRecurrenceAt"`
	CurrentCount     int    `json:"currentCount"`
	RemainingCount   int    `json:"remainingCount"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
}

// VirtualCardRevision -> https://developer.paywithextend.com/#tocS_VirtualCardRevision.
type VirtualCardRevision struct {
	BalanceCents       int            `json:"balanceCents"`
	ValidFrom          string         `json:"validFrom"`
	ValidTo            string         `json:"validTo"`
	Recurs             bool           `json:"recurs"`
	ActiveUntil        string         `json:"activeUntil"`
	Currency           string         `json:"currency"`
	Recurrence         Recurrence     `json:"recurrence"`
	ReceiptAttachments map[string]any `json:"receiptAttachments"`
}

// Address -> https://developer.paywithextend.com/#tocS_Address.
type Address struct {
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	City     string `json:"city"`
	Province string `json:"province"`
	Postal   string `json:"postal"`
	Country  string `json:"country"`
}

// VirtualCardFeature -> https://developer.paywithextend.com/#tocS_VirtualCardFeature.
type VirtualCardFeature struct {
	Recurrence       bool   `json:"recurrence"`
	CustomAddress    bool   `json:"customAddress"`
	CustomMin        bool   `json:"customMin"`
	CustomMax        bool   `json:"customMax"`
	WalletsEnabled   string `json:"walletsEnabled"`
	MccControl       bool   `json:"mccControl"`
	QboReportEnabled bool   `json:"qboReportEnabled"`
}

// MccRange -> https://developer.paywithextend.com/#tocS_MccRange.
type MccRange struct {
	Lowest  int `json:"lowest"`
	Highest int `json:"highest"`
}
