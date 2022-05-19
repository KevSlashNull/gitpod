// Copyright (c) 2022 Gitpod GmbH. All rights reserved.
// Licensed under the GNU Affero General Public License (AGPL).
// See License-AGPL.txt in the project root for license information.

package cmd

import (
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "idea-cli",
}

func Execute() {
	waitUntilBackendPluginIsAccessible()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func waitUntilBackendPluginIsAccessible() {
	for {
		resp, httpError := http.Get("http://localhost:63342/api/gitpod/cli?op=ping")
		if httpError == nil && resp.StatusCode == http.StatusOK {
			break
		}
		time.Sleep(1000)
	}
}

func init() {}
