package rootgrp

import (
	"context"
	"hosting-resources-service/cmd/server/rest/gen"
)

type Handlers struct {
	prefix string
}

func New(prefix string) *Handlers {
	return &Handlers{prefix: prefix}
}

func (h *Handlers) GetRoot(ctx context.Context, request gen.GetRootRequestObject) (gen.GetRootResponseObject, error) {
	return gen.GetRoot200ApplicationHalPlusJSONResponse(toRoot(h.prefix)), nil
}
