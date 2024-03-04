/*
Copyright (c) 2018 Red Hat, Inc.

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

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github/yasun1/myquota/cmd/myquota/assign"
	"github/yasun1/myquota/cmd/myquota/list"
	"github/yasun1/myquota/cmd/myquota/remove"
	"github/yasun1/myquota/pkg/flags"

	_ "github.com/golang/glog"
	"github.com/spf13/cobra"
	// "github.com/spf13/pflag"
)

var root = &cobra.Command{
	Use: "myquota",
	Long: "Command line tool for manage ocm resource quotas." +
		" The default stage is stage ocm, setting OCM_ENV=prod will change to prod ocm.",
}

func init() {
	// Send logs to the standard error stream by default:
	// err := flag.Set("logtostderr", "true")
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Can't set default error stream: %v\n", err)
	// 	os.Exit(1)
	// }

	// Register the options that are managed by the 'flag' package, so that they will also be parsed
	// by the 'pflag' package:
	// pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	// Add the command line flags:
	fs := root.PersistentFlags()
	flags.AddDebugFlag(fs)

	// Set log title
	log.SetPrefix("[quota] ")
	log.SetFlags(log.LstdFlags | log.LUTC)

	// Register the subcommands:
	root.AddCommand(assign.Cmd)
	root.AddCommand(remove.Cmd)
	root.AddCommand(list.Cmd)
}

func main() {
	// This is needed to make `glog` believe that the flags have already been parsed, otherwise
	// every log messages is prefixed by an error message stating the the flags haven't been
	// parsed.
	err := flag.CommandLine.Parse([]string{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't parse empty command line to satisfy 'glog': %v\n", err)
		os.Exit(1)
	}

	// Execute the root command:
	root.SetArgs(os.Args[1:])
	err = root.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute root command: %v\n", err)
		os.Exit(1)
	}
}
