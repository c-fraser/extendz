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

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/c-fraser/extendz"
	extend "github.com/c-fraser/extendz/pkg/client"
	"github.com/hokaccha/go-prettyjson"
	"github.com/urfave/cli/v2"
)

const (
	extendApiBaseUrl = "https://api.paywithextend.com/"
	emailEnv         = "EXTEND_EMAIL"
	passwordEnv      = "EXTEND_PASSWORD"
)

// main is the entry point into the extendz CLI application.
func main() {
	email := os.Getenv(emailEnv)
	password := os.Getenv(passwordEnv)
	if email == "" || password == "" {
		log.Fatalf("The '%s' and '%s' environment variables must be set", emailEnv, passwordEnv)
	}
	client, err := extend.NewClient(extendApiBaseUrl, email, password)
	if err != nil {
		log.Fatalf("Failed to initialize Extend API client: %v", err)
	}

	app := cli.NewApp()
	app.Name = "extendz"
	app.Usage = "A tool for interacting with the Extend API"
	app.Version = extendz.VERSION
	app.Commands = cli.Commands{
		&cli.Command{
			Name:  "get-user-virtual-cards",
			Usage: "Get the virtual cards for a user",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "request",
					Aliases:  []string{"r"},
					Usage:    "the https://developer.paywithextend.com/#tocS_VirtualCardPageableRequest JSON",
					Required: false,
				},
			},
			Action: func(c *cli.Context) error {
				s := c.String("request")
				var request extend.VirtualCardPageableRequest
				if s != "" {
					err := json.Unmarshal([]byte(s), &request)
					if err != nil {
						return err
					}
				}
				response, err := client.GetUserVirtualCards(&request)
				if err != nil {
					return err
				}
				return printResponse(response)
			},
		},
		&cli.Command{
			Name:  "get-virtual-card",
			Usage: "Get a virtual card",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "id",
					Aliases:  []string{"i"},
					Usage:    "the virtual card ID",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				id := c.String("id")
				response, err := client.GetVirtualCard(id)
				if err != nil {
					return err
				}
				return printResponse(response)
			},
		},
		&cli.Command{
			Name:  "get-virtual-card-transactions",
			Usage: "Get the transactions for a virtual card",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "id",
					Aliases:  []string{"i"},
					Usage:    "the virtual card ID",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "count",
					Aliases:  []string{"c"},
					Usage:    "the number of transactions to get",
					Required: false,
				},
				&cli.StringFlag{
					Name:     "before",
					Aliases:  []string{"b"},
					Usage:    "get transactions before timestamp",
					Required: false,
				},
				&cli.StringFlag{
					Name:     "after",
					Aliases:  []string{"a"},
					Usage:    "get transactions after timestamp",
					Required: false,
				},
				&cli.StringFlag{
					Name:     "status",
					Aliases:  []string{"s"},
					Usage:    "the comma-delimited list of transaction statuses to get",
					Required: false,
				},
			},
			Action: func(c *cli.Context) error {
				id := c.String("id")
				count := c.Int("count")
				before := c.String("before")
				after := c.String("after")
				status := c.String("status")
				response, err := client.GetVirtualCardTransactions(id, count, before, after, status)
				if err != nil {
					return err
				}
				return printResponse(response)
			},
		},
		&cli.Command{
			Name:  "create-virtual-card",
			Usage: "Create a virtual card",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "request",
					Aliases:  []string{"r"},
					Usage:    "the https://developer.paywithextend.com/#tocS_CreateVirtualCardRequest JSON",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				s := c.String("request")
				var request extend.CreateVirtualCardRequest
				err := json.Unmarshal([]byte(s), &request)
				if err != nil {
					return err
				}
				response, err := client.CreateVirtualCard(&request)
				if err != nil {
					return err
				}
				return printResponse(response)
			},
		},
		&cli.Command{
			Name:  "update-virtual-card",
			Usage: "Update a virtual card",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "id",
					Aliases:  []string{"i"},
					Usage:    "the virtual card ID",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "request",
					Aliases:  []string{"r"},
					Usage:    "the https://developer.paywithextend.com/#tocS_UpdateVirtualCardRequest JSON",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				id := c.String("id")
				s := c.String("request")
				var request extend.UpdateVirtualCardRequest
				err := json.Unmarshal([]byte(s), &request)
				if err != nil {
					return err
				}
				response, err := client.UpdateVirtualCard(id, &request)
				if err != nil {
					return err
				}
				return printResponse(response)
			},
		},
		&cli.Command{
			Name:  "cancel-virtual-card",
			Usage: "Cancel a virtual card",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "id",
					Aliases:  []string{"i"},
					Usage:    "the virtual card ID",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				id := c.String("id")
				response, err := client.CancelVirtualCard(id)
				if err != nil {
					return err
				}
				return printResponse(response)
			},
		},
		&cli.Command{
			Name:  "reject-virtual-card",
			Usage: "Reject a virtual card",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "id",
					Aliases:  []string{"i"},
					Usage:    "the virtual card ID",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				id := c.String("id")
				response, err := client.RejectVirtualCard(id)
				if err != nil {
					return err
				}
				return printResponse(response)
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// printResponse prints the response as (pretty) JSON.
func printResponse(response any) error {
	b, err := prettyjson.Marshal(response)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}
