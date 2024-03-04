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

package list

import (
	"fmt"
	"os"

	"github/yasun1/myquota/pkg/quota"

	"github.com/spf13/cobra"
)

var args struct {
	username string
}

var Cmd = &cobra.Command{
	Use:   "list <skuIDs>",
	Short: "List the quota cost under the account",
	Long: "List the quota cost in the organization that the account is belonged to. " +
		"If no skuIDs are specified, will list all the quota of the organization.",
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
		quota.FPrintQuotaCost(orgID)
	} else {
		skuMap := quota.AllSkus()

		var specifiedSKus []quota.Sku
		for _, skuName := range argv {
			if _, existed := skuMap[skuName]; !existed {
				panic(fmt.Errorf("[E] The sku '%s' is invalid\n", skuName))
			}

			specifiedSKus = append(specifiedSKus, skuMap[skuName])
		}

		quota.FPrintUsageForSkus(orgID, specifiedSKus...)
	}
}
