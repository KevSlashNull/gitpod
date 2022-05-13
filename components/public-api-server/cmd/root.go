// Copyright (c) 2022 Gitpod GmbH. All rights reserved.
// Licensed under the GNU Affero General Public License (AGPL).
// See License-AGPL.txt in the project root for license information.

package cmd

import (
	"context"
	"os"

	"github.com/gitpod-io/gitpod/common-go/log"
	"github.com/spf13/cobra"
)

var (
	// ServiceName is the name we use for tracing/logging
	ServiceName = "public-api-server"
	// Version of this service - set during build
	Version = ""
)

var rootOpts struct {
	CfgFile string
	JsonLog bool
	Verbose bool
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   ServiceName,
	Short: "Serves public API services",
}

func Execute() {
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		log.WithError(err).Error("Failed to execute command.")
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&rootOpts.CfgFile, "config", "", "config file")
	rootCmd.PersistentFlags().BoolVarP(&rootOpts.JsonLog, "json-log", "j", true, "produce JSON log output on verbose level")
	rootCmd.PersistentFlags().BoolVarP(&rootOpts.Verbose, "verbose", "v", false, "Enable verbose JSON logging")
}

func getConfig() *config.Configuration {
	ctnt, err := os.ReadFile(rootOpts.CfgFile)
	if err != nil {
		log.WithError(err).Fatal("cannot read configuration. Maybe missing --config?")
	}

	var cfg config.Configuration
	dec := json.NewDecoder(bytes.NewReader(ctnt))
	dec.DisallowUnknownFields()
	err = dec.Decode(&cfg)
	if err != nil {
		log.WithError(err).Fatal("cannot decode configuration. Maybe missing --config?")
	}

	return &cfg
}
