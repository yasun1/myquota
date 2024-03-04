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

package remove

import (
	"fmt"
	"os"

	"github/yasun1/myquota/pkg/quota"

	"github.com/spf13/cobra"
)

var args struct {
	username string
	qtype    string
	force    bool
}

var Cmd = &cobra.Command{
	Use:   "remove <skuID>",
	Short: "Remove the 'Manual' resource quota under the account",
	Long:  "Remove the 'Manual' resource quota from the organization that the account is belonged to.",
	Run:   run,
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
	fs.BoolVarP(
		&args.force,
		"force",
		"f",
		false,
		"If the force is true, will ignore checking the consumed quota and forcely remove the quota from the organization.",
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
		panic(fmt.Errorf("[E] The input sku '%s' is invalid\n", skuName))
	}
	sku := skuMap[skuName]
	sku.Type = args.qtype

	// Remove the quota
	quota.RemoveQuota(orgID, sku, args.force)
}
