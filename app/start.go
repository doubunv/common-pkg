package app

import (
	"github.com/zeromicro/go-zero/rest"
	"gitlab.888bbm.com/go-package/common-pkg/appMiddleware"
)

type SMOption func(s *ServerMiddleware)

func WithWhiteHeaderPathSMOption(whiteHeader map[string]int) SMOption {
	return func(s *ServerMiddleware) {
		s.whiteHeader = whiteHeader
	}
}

func WithDebugOption() SMOption {
	return func(s *ServerMiddleware) {
		s.isDebug = true
	}
}

func WithTestOption() SMOption {
	return func(s *ServerMiddleware) {
		s.isTest = true
	}
}

func WithCheckTokenHandleSMOption(fun appMiddleware.CheckRequestTokenFunc) SMOption {
	return func(s *ServerMiddleware) {
		s.checkTokenHandle = fun
	}
}

type ServerMiddleware struct {
	whiteHeader      map[string]int
	checkTokenHandle appMiddleware.CheckRequestTokenFunc

	Server *rest.Server

	isDebug bool
	isTest  bool
}

func NewServerMiddleware(s *rest.Server, opt ...SMOption) *ServerMiddleware {
	res := &ServerMiddleware{
		Server: s,
	}

	for _, item := range opt {
		item(res)
	}

	return res
}

func (s *ServerMiddleware) ApiUseMiddleware() {
	s.Server.Use(appMiddleware.NewCorsMiddleware().Handle)
	s.useApiHeaderMiddleware()
	s.mustUserAgentMiddleware()
}

func (s *ServerMiddleware) useApiHeaderMiddleware() {
	var apiHeaderOption = []appMiddleware.ApiHeadOption{
		appMiddleware.CloseVerifyOption(s.whiteHeader),
	}
	if s.isDebug {
		apiHeaderOption = append(apiHeaderOption, appMiddleware.WithDebugOption())
	}
	s.Server.Use(appMiddleware.NewApiHeaderMiddleware(
		apiHeaderOption...,
	).Handle)
}

func (s *ServerMiddleware) mustUserAgentMiddleware() {
	if s.checkTokenHandle == nil {
		panic("must use CheckTokenHandleSMOption.")
	}

	s.Server.Use(appMiddleware.NewUserAgentMiddleware(
		s.whiteHeader,
		appMiddleware.WithCheckOption(s.checkTokenHandle),
	).Handle)
}
