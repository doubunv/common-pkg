package headInfo

import (
	"context"
	"github.com/doubunv/common-pkg/consts"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
	"strconv"
	"strings"
)

func GetTrace(ctx context.Context) string {
	return trace.SpanContextFromContext(ctx).TraceID().String()
}

func GetTokenUid(ctx context.Context) int64 {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return 0
	}

	list := md.Get(consts.TokenUid)
	parseInt, err := strconv.ParseInt(strings.Join(list, ""), 10, 64)
	if err != nil {
		return 0
	}

	return parseInt
}

func GetTokenUidRole(ctx context.Context) string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return ""
	}
	res := strings.Join(md.Get(consts.TokenUidRole), "")
	return res
}

func GetJwtToken(ctx context.Context) string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return ""
	}
	res := strings.Join(md.Get(consts.HeaderToken), "")
	return res
}

func GetClientIp(ctx context.Context) string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return ""
	}
	res := strings.Join(md.Get(consts.ClientIp), "")
	return res
}
func GetUserAgent(ctx context.Context) string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return ""
	}
	res := strings.Join(md.Get(consts.UserAgent), "")
	return res
}

func GetVersion(ctx context.Context) string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return ""
	}
	res := strings.Join(md.Get(consts.Version), "")
	return res
}

func GetSource(ctx context.Context) string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return ""
	}
	res := strings.Join(md.Get(consts.Source), "")
	return res
}

func SetTokenUid(ctx context.Context, value string) context.Context {
	md, _ := metadata.FromOutgoingContext(ctx)
	md = md.Copy()
	md.Set(consts.TokenUid, value)
	return metadata.NewOutgoingContext(ctx, md)
}

//func SetTokenUid(ctx context.Context, value string) context.Context {
//	md := ctxMd.SetMdCtxFromOut(ctx, consts.TokenUid, value)
//	ctx = metadata.NewOutgoingContext(ctx, md)
//	return ctx
//}

func GetTrance(ctx context.Context) string {
	return trace.SpanContextFromContext(ctx).TraceID().String()
}

func GetContentLanguage(ctx context.Context) string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return ""
	}
	res := strings.Join(md.Get(consts.ContentLanguage), "")
	return res
}

func GetBusinessCode(ctx context.Context) string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return ""
	}
	res := strings.Join(md.Get(consts.BusinessCode), "")
	return res
}

func SetBusinessCode(ctx context.Context, value string) context.Context {
	md, _ := metadata.FromOutgoingContext(ctx)
	md = md.Copy()
	md.Set(consts.BusinessCode, value)
	return metadata.NewOutgoingContext(ctx, md)
}

//func SetBusinessCode(ctx context.Context, value string) context.Context {
//	md := ctxMd.SetMdCtxFromOut(ctx, consts.BusinessCode, value)
//	ctx = metadata.NewOutgoingContext(ctx, md)
//	return ctx
//}

func GetBusiness(ctx context.Context) string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return ""
	}
	res := strings.Join(md.Get(consts.Business), "")
	return res
}

func GetOriginHostUrl(ctx context.Context) string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return ""
	}
	res := strings.Join(md.Get(consts.OriginUrl), "")

	res = strings.ReplaceAll(res, "http://", "")
	res = strings.ReplaceAll(res, "https://", "")

	return res
}

func GetTerminal(ctx context.Context) string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return ""
	}
	res := strings.Join(md.Get(consts.Source), "")
	if res == consts.IOS || res == consts.Android {
		return consts.APP
	}
	return res
}

func GetTimezone(ctx context.Context) string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return ""
	}
	res := strings.Join(md.Get(consts.Timezone), "")
	return res
}

func GetDeviceId(ctx context.Context) string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return ""
	}
	res := strings.Join(md.Get(consts.DeviceId), "")
	return res
}

func GetReqPath(ctx context.Context) string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return ""
	}
	res := strings.Join(md.Get(consts.ReqPath), "")
	return res
}
