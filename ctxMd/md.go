package ctxMd

import (
	"context"
	"google.golang.org/grpc/metadata"
)

func SetMdCtxFromOut(ctx context.Context, key string, value string) metadata.MD {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	md.Set(key, value)

	return md
}
