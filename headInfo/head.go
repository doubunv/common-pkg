package headInfo

import (
	"context"
	"encoding/json"
	"gitlab.888bbm.com/go-package/common-pkg/consts"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
	"net"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type Head struct {
	AuthorizationJwt string `json:"authorization_jwt"` // 用户token
	Version          string `json:"version"`           // APP版本
	Source           string `json:"source"`            // 来源渠道	* Android * Ios * Pc
	ClientIp         string `json:"client_ip"`         // 客户端IP
	Trace            string `json:"trace"`             // 链路路由
	TokenUid         string `json:"token_uid"`         // 用户ID
	ReqPath          string `json:"req_path"`          // 请求path
	Business         string `json:"business"`
	BusinessCode     string `json:"business_code"`
	ContentLanguage  string `json:"content_language"`
}

func GetHead(r *http.Request) *Head {
	header := r.Header
	return &Head{
		AuthorizationJwt: strings.Trim(header.Get(consts.HeaderToken), " "),
		Version:          strings.Trim(header.Get("Version"), " "),
		Source:           strings.Trim(header.Get("Source"), " "),
		ClientIp:         getClientIP(r),
		TokenUid:         strings.Trim(header.Get("TokenUid"), " "),
		Trace:            trace.SpanContextFromContext(r.Context()).TraceID().String(),
		ReqPath:          r.URL.Path,
		Business:         strings.Trim(header.Get(consts.Business), " "),
		BusinessCode:     strings.Trim(header.Get(consts.BusinessCode), " "),
		ContentLanguage:  strings.Trim(header.Get(consts.ContentLanguage), " "),
	}
}

func (h *Head) Verify() error {
	return nil
}

func (h *Head) String() string {
	data, _ := json.Marshal(h)
	return string(data)
}

func ContextHeadInLog(ctx context.Context, h *Head) context.Context {
	ctxNew := logx.ContextWithFields(ctx,
		logx.Field(consts.Token, h.AuthorizationJwt),
		logx.Field(consts.Version, h.Version),
		logx.Field(consts.Source, h.Source),
		logx.Field(consts.ClientIp, h.ClientIp),
		logx.Field(consts.Trace, h.Trace),
		logx.Field(consts.TokenUid, h.TokenUid),
		logx.Field(consts.ReqPath, h.ReqPath),
		logx.Field(consts.Business, h.Business),
		logx.Field(consts.BusinessCode, h.BusinessCode),
		logx.Field(consts.ContentLanguage, h.BusinessCode),
	)
	return ctxNew
}

func getClientIP(r *http.Request) string {
	ip := r.Header.Get("x_forwarded_realip")
	if ip == "" {
		ip = r.Header.Get("X-Real-Ip")
	}
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	return ip
}

func GetFullHead(r *http.Request) map[string][]string {
	headers := make(map[string][]string)

	for k, v := range r.Header {
		headers[k] = v
	}

	return headers
}

func HeadInMetadata(ctx context.Context, h Head) context.Context {
	md := metadata.Pairs(
		consts.Token, h.AuthorizationJwt,
		consts.TokenUid, h.TokenUid,
		consts.ClientIp, h.ClientIp,
		consts.ReqPath, h.ReqPath,
		consts.Version, h.Version,
		consts.Source, h.Source,
		consts.Business, h.Business,
		consts.BusinessCode, h.BusinessCode,
		consts.ContentLanguage, h.ContentLanguage,
	)

	ctxNew := metadata.NewOutgoingContext(ctx, md)
	return ctxNew
}
