// Copyright (c) 2022 Gitpod GmbH. All rights reserved.
// Licensed under the GNU Affero General Public License (AGPL).
// See License-AGPL.txt in the project root for license information.

package server

import (
	"fmt"
	"net/url"

	"github.com/gitpod-io/gitpod/common-go/baseserver"
	"github.com/gitpod-io/gitpod/public-api-server/pkg/apiv1"
	"github.com/gitpod-io/gitpod/public-api-server/pkg/proxy"
	"github.com/gitpod-io/gitpod/public-api/config"
	v1 "github.com/gitpod-io/gitpod/public-api/v1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func Start(logger *logrus.Entry, cfg *config.Configuration) error {
	gitpodAPI, err := url.Parse(cfg.GitpodServiceURL)
	if err != nil {
		logger.WithError(err).Fatal("Failed to parse Gitpod API URL.")
	}

	registry := prometheus.NewRegistry()

	srv, err := baseserver.New("public_api_server",
		baseserver.WithLogger(logger),
		baseserver.WithDebugPort(cfg.PProfPort),
		baseserver.WithGRPC(&baseserver.ServerConfiguration{Address: fmt.Sprintf(":%d", cfg.GRPCPort)}),
		baseserver.WithMetricsRegistry(registry),
	)
	if err != nil {
		return fmt.Errorf("failed to initialize public api server: %w", err)
	}

	if registerErr := register(srv, gitpodAPI, registry); registerErr != nil {
		return fmt.Errorf("failed to register services: %w", registerErr)
	}

	if listenErr := srv.ListenAndServe(); listenErr != nil {
		return fmt.Errorf("failed to serve public api server: %w", err)
	}

	return nil
}

func register(srv *baseserver.Server, serverAPIURL *url.URL, registry *prometheus.Registry) error {
	proxy.RegisterMetrics(registry)

	connPool := &proxy.NoConnectionPool{ServerAPI: serverAPIURL}

	v1.RegisterWorkspacesServiceServer(srv.GRPC(), apiv1.NewWorkspaceService(connPool))
	v1.RegisterPrebuildsServiceServer(srv.GRPC(), v1.UnimplementedPrebuildsServiceServer{})

	return nil
}
