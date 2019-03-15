// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime_provider

import (
	"google.golang.org/grpc"

	providerclient "openpitrix.io/openpitrix/pkg/client/runtime_provider"
	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	runtimeprovider "openpitrix.io/openpitrix/pkg/service/runtime_provider"
)

type Server struct {
	runtimeprovider.Server
}

func Serve(cfg *config.Config) {
	pi.SetGlobal(cfg)
	err := providerclient.RegisterRuntimeProvider(Provider, ProviderConfig)
	if err != nil {
		logger.Critical(nil, "failed to register provider config: %+v", err)
	}
	s := Server{}
	manager.NewGrpcServer("runtime-provider-kubernetes", constants.RuntimeProviderManagerPort).
		ShowErrorCause(cfg.Grpc.ShowErrorCause).
		Serve(func(server *grpc.Server) {
			pb.RegisterRuntimeProviderManagerServer(server, &s)
		})
}
