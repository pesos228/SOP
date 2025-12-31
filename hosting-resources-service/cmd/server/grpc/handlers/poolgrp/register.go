package poolgrp

import (
	"hosting-resources-service/cmd/server/grpc/handlers/poolgrp/gen"
	"hosting-resources-service/internal/pool"

	"google.golang.org/grpc"
)

func Register(serv *grpc.Server, poolBus *pool.Business) {
	apiImpl := New(poolBus)
	gen.RegisterResourcesServer(serv, apiImpl)
}
