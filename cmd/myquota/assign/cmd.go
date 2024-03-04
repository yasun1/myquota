/*
Copyright (c) 2019 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package assign

import (
	"fmt"
	"os"

	"github/yasun1/myquota/pkg/quota"

	"github.com/spf13/cobra"
)

var args struct {
	username string
	qtype    string
	number   int
	force    bool
}

var Cmd = &cobra.Command{
	Use:   "assign <skuID>",
	Short: "Assign the resource quota to the account",
	Long: "Assign the  resource quota to the account. " +
		"If the resource quota does not exist, create the resource quota for it; " +
		"If the resource quota exists, update the resource quota to the specified value.",
	Run: run,
}

func init() {
	fs := Cmd.Flags()
	fs.StringVarP(
		&args.username,
		"username",
		"u",
		"",
		"The username of the account.",
	)
	fs.StringVarP(
		&args.qtype,
		"qtype",
		"t",
		"Manual",
		"The type of the quota.",
	)
	fs.IntVarP(
		&args.number,
		"number",
		"n",
		0,
		"The number is the applied sku account.",
	)
}

func run(cmd *cobra.Command, argv []string) {
	if args.username == "" {
		fmt.Fprintf(os.Stderr, "[E] The option '--username' is mandatory.\n\n")
		os.Exit(1)
	}

	orgID, err := quota.GetOrgID(args.username)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(0)
	}

	if len(argv) == 0 {
		fmt.Fprintf(os.Stderr, "[E] The sku id is required.\n\n")
		os.Exit(1)
	}

	skuName := argv[0]
	skuMap := quota.AllSkus()
	if _, existed := skuMap[skuName]; !existed {
		panic(fmt.Errorf("[E] The input sku '%s' is invalid", skuName))
	}
	sku := skuMap[skuName]
	sku.Allowed = args.number
	sku.Type = args.qtype

	// Assign quota
	_, err = quota.AssignQuota(orgID, sku)
	if err != nil {
		panic(err)
	}

	// Print the usage of the just assigned quota
	quota.FPrintUsageForSkus(orgID, sku)
}
