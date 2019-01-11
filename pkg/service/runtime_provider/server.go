// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime_provider

import (
	"google.golang.org/grpc"

	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	runtimeprovider "openpitrix.io/openpitrix/pkg/service/runtime_provider"
)

type Server struct {
}

func Serve(cfg *config.Config) {
	pi.SetGlobal(cfg)
	err := runtimeprovider.RegisterRuntimeProvider(Provider, ProviderConfig)
	if err != nil {
		logger.Critical(nil, "failed to register provider config: %+v", err)
	}
	s := Server{}
	manager.NewGrpcServer("runtime-provider-qingcloud", constants.RuntimeProviderManagerPort).
		ShowErrorCause(cfg.Grpc.ShowErrorCause).
		Serve(func(server *grpc.Server) {
			pb.RegisterRuntimeProviderManagerServer(server, &s)
		})
}
