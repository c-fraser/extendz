# extendz

[![GoDoc](https://godoc.org/github.com/c-fraser/extendz?status.svg)](https://godoc.org/github.com/c-fraser/extendz)
[![Release](https://img.shields.io/github/v/release/c-fraser/extendz?logo=github&sort=semver)](https://github.com/c-fraser/extendz/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/c-fraser/extendz)](https://goreportcard.com/report/github.com/c-fraser/extendz)
[![Apache License 2.0](https://img.shields.io/badge/License-Apache2-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)

Utilities for interacting with the [Extend API](https://developer.paywithextend.com/#extend-api).

## CLI

The [CLI application](cmd/cli) enables [client operations](#client) to be executed via the command
line.

## Client

The [client package](pkg/client) contains a REST client for
the [Extend API](https://developer.paywithextend.com/#extend-api).

### Operations

- [X] Authentication
    - [X] Sign In
    - [X] Sign Out
    - [X] Renew Auth
    - [X] Forgot Password
- [ ] Virtual Cards
    - [X] Get User Virtual Cards
    - [X] Get Virtual Card
    - [ ] Get Virtual Card History
    - [X] Get Virtual Card Transactions
    - [X] Create Virtual Card
    - [X] Update Virtual Card
    - [ ] Update Virtual Card Customer Support Code
    - [X] Cancel Virtual Card
    - [ ] Cancel Virtual Card Update Request
    - [X] Reject Virtual Card
    - [ ] Get Virtual Card Permissions
    - [ ] Bulk Virtual Card Push
    - [ ] Get Bulk Push XLSX Template
    - [ ] Get Bulk Virtual Card Upload Statuses
    - [ ] Simulate Transaction for Test Cards
- [ ] Credit Cards
    - [ ] Get Credit Card
    - [ ] Update Credit Card
    - [ ] Delete Credit Card
    - [ ] Get User Credit Cards
    - [ ] Get Credit Card Permissions
    - [ ] Begin Verify Credit Cardholder Process
    - [ ] Update Credit Card Status
    - [ ] Verify Credit Cardholder
- [ ] Events
    - [ ] Get Event
    - [ ] Event List
- [ ] Users
    - [ ] Create User
    - [ ] Set User Avatar
    - [ ] Get Expensify Links
    - [ ] Create Expensify Link
    - [ ] Remove Expensify Links from User
    - [ ] Delete User
    - [ ] Get User
    - [ ] List Users
    - [ ] Update User
    - [ ] Login to Quickbooks
    - [ ] Revoke Quickbooks Token
    - [ ] Resend Email Verification Code
    - [ ] Verify Email
- [ ] Metrics
    - [ ] Get Spend Metrics
- [ ] Organizations
    - [ ] Get Organizations
    - [ ] Create Organization
    - [ ] Get Organization
    - [ ] Update Organization
    - [ ] Delete Organization
    - [ ] Get Organization Members
    - [ ] Invite a list of users by email address to an Organization
    - [ ] Get Invites
    - [ ] Get Organization Permissions
- [ ] Attachments
    - [ ] Get User Attachments
    - [ ] Upload Attachment
    - [ ] Get Attachment
    - [ ] Delete Attachment
    - [ ] Update Attachment
    - [ ] Get Attachment Permissions
- [ ] References
    - [ ] Get Reference Fields for Credit Card
    - [ ] Update Reference Fields for Credit Card
- [ ] Reports
    - [ ] Get Report
    - [ ] Generate Transaction Report
    - [ ] Get Transactions Report
    - [ ] Generate Virtual Card Report
- [ ] Statistics
    - [ ] Get Statistics
- [ ] Subscriptions
    - [ ] Get Webhook Attempts for Subscription
    - [ ] Subscription List
    - [ ] Create a Subscription
    - [ ] Get Subscription
    - [ ] Update Subscription
    - [ ] Delete Subscription
    - [ ] Regenerate secret for a Subscription
- [ ] Transactions
    - [ ] Get Transaction
    - [ ] Update Transaction
    - [ ] Get Transaction Permissions
