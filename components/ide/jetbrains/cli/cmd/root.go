// Copyright (c) 2022 Gitpod GmbH. All rights reserved.
// Licensed under the GNU Affero General Public License (AGPL).
// See License-AGPL.txt in the project root for license information.

package cmd

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "idea-cli",
}

func Execute() {
	waitUntilBackendPluginIsReady()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func waitUntilBackendPluginIsReady() {
	for {
		// We expect the backend plugin API to return "400 Bad Request" as the operation ("op" query parameter) is undefined.
		resp, httpError := http.Get("http://localhost:63342/api/gitpod/cli")
		if httpError != nil {
			log.Fatal(httpError)
		}
		if resp.StatusCode == http.StatusBadRequest {
			break
		}
		time.Sleep(1000)
	}
}

func init() {}
