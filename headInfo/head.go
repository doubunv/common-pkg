package headInfo

import (
	"context"
	"encoding/json"
	"github.com/doubunv/common-pkg/consts"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
	"net"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

const PLATFORM_BUSINESS_CODE = "99999999"

type Head struct {
	AuthorizationJwt string `json:"authorization_jwt"` // 用户token
	Version          string `json:"version"`           // APP版本
	Source           string `json:"source"`            // 来源渠道
	ClientIp         string `json:"client_ip"`         // 客户端IP
	Trace            string `json:"trace"`             // 链路路由
	TokenUid         string `json:"token_uid"`         // 用户ID
	ReqPath          string `json:"req_path"`          // 请求path
	Business         string `json:"business"`
	BusinessCode     string `json:"business_code"`
	ContentLanguage  string `json:"content_language"`
	TokenUidRole     string `json:"token_uid_role"`
	ReqOrigin        string `json:"req_origin"` // 请求地址
	Timezone         string `json:"timezone"`
	UserAgent        string `json:"user_agent"`
	DeviceId         string `json:"device_id"`
}

func GetHead(r *http.Request) *Head {
	header := r.Header
	return &Head{
		AuthorizationJwt: strings.Trim(header.Get(consts.HeaderToken), " "),
		Version:          strings.Trim(header.Get(consts.Version), " "),
		Source:           strings.Trim(header.Get(consts.Source), " "),
		ClientIp:         getClientIP(r),
		Trace:            getTraceId(r),
		ReqPath:          r.URL.Path,
		Business:         strings.Trim(header.Get(consts.Business), " "),
		BusinessCode:     strings.Trim(header.Get(consts.BusinessCode), " "),
		ContentLanguage:  strings.Trim(header.Get(consts.ContentLanguage), " "),
		TokenUid:         "",
		TokenUidRole:     "",
		ReqOrigin:        strings.Trim(header.Get("Origin"), " "),
		Timezone:         strings.Trim(header.Get(consts.Timezone), " "),
		UserAgent:        strings.Trim(header.Get("user-agent"), ""),
		DeviceId:         strings.Trim(header.Get(consts.ContentLanguage), ""),
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
		logx.Field(consts.HeaderToken, h.AuthorizationJwt),
		logx.Field(consts.Version, h.Version),
		logx.Field(consts.Source, h.Source),
		logx.Field(consts.ClientIp, h.ClientIp),
		logx.Field(consts.Trace, h.Trace),
		logx.Field(consts.TokenUid, h.TokenUid),
		logx.Field(consts.ReqPath, h.ReqPath),
		logx.Field(consts.Business, h.Business),
		logx.Field(consts.BusinessCode, h.BusinessCode),
		logx.Field(consts.ContentLanguage, h.ContentLanguage),
		logx.Field(consts.TokenUidRole, h.TokenUidRole),
		logx.Field(consts.OriginUrl, h.ReqOrigin),
		logx.Field(consts.Timezone, h.Timezone),
		logx.Field(consts.UserAgent, h.UserAgent),
		logx.Field(consts.DeviceId, h.DeviceId),
	)
	return ctxNew
}

func getTraceId(r *http.Request) string {
	traceId := r.Header.Get("X-Request-ID")
	if traceId == "" {
		traceId = strings.Trim(r.Header.Get(consts.RequestId), " ")
	}
	if traceId == "" {
		traceId = trace.SpanContextFromContext(r.Context()).TraceID().String()
	}
	return traceId
}

func getClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		ip = strings.Split(ip, ",")[0]
	}
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
		consts.HeaderToken, h.AuthorizationJwt,
		consts.TokenUid, h.TokenUid,
		consts.ClientIp, h.ClientIp,
		consts.ReqPath, h.ReqPath,
		consts.Version, h.Version,
		consts.Source, h.Source,
		consts.Business, h.Business,
		consts.BusinessCode, h.BusinessCode,
		consts.ContentLanguage, h.ContentLanguage,
		consts.TokenUidRole, h.TokenUidRole,
		consts.OriginUrl, h.ReqOrigin,
		consts.Timezone, h.Timezone,
		consts.UserAgent, h.UserAgent,
		consts.Trace, h.Trace,
		consts.DeviceId, h.DeviceId,
	)

	ctxNew := metadata.NewOutgoingContext(ctx, md)
	return ctxNew
}
