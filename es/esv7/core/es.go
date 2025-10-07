package core

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	es7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type (
	Config struct {
		Addresses  []string
		Username   string
		Password   string
		MaxRetries int
	}

	Es struct {
		*es7.Client
	}

	// esTransport is a transport for elasticsearch client
	esTransport struct{}
)

func (t *esTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	var (
		ctx        = req.Context()
		span       oteltrace.Span
		startTime  = time.Now()
		propagator = otel.GetTextMapPropagator()
		indexName  = strings.Split(req.URL.RequestURI(), "/")[1]
		tracer     = trace.TracerFromContext(ctx)
	)

	ctx, span = tracer.Start(ctx,
		req.URL.Path,
		oteltrace.WithSpanKind(oteltrace.SpanKindClient),
		oteltrace.WithAttributes(semconv.HTTPClientAttributesFromHTTPRequest(req)...),
	)
	defer func() {
		metricClientReqDur.Observe(time.Since(startTime).Milliseconds(), indexName)
		metricClientReqErrTotal.Inc(indexName, strconv.FormatBool(err != nil))

		span.End()
	}()

	req = req.WithContext(ctx)
	propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err = http.DefaultTransport.RoundTrip(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	span.SetAttributes(semconv.DBSQLTableKey.String(indexName))
	span.SetAttributes(semconv.HTTPAttributesFromHTTPStatusCode(resp.StatusCode)...)
	span.SetStatus(semconv.SpanStatusFromHTTPStatusCodeAndSpanKind(resp.StatusCode, oteltrace.SpanKindClient))

	return
}

// BuildAuthHeader 构建认证头
func BuildAuthHeader(conf *Config) http.Header {
	headers := http.Header{}
	authValue := "Basic " + basicAuth(conf.Username, conf.Password)
	headers.Set("Authorization", authValue)
	return headers
}

// basicAuth 生成 Basic Auth 头
func basicAuth(username, password string) string {
	cred := fmt.Sprintf("%s:%s", username, password)
	return base64.StdEncoding.EncodeToString([]byte(cred))
}

func NewEs(conf *Config) (*Es, error) {
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		// 生产环境建议关闭
		InsecureSkipVerify: true,
	}
	c := es7.Config{
		Addresses:  conf.Addresses,
		Username:   conf.Username,
		Password:   conf.Password,
		MaxRetries: conf.MaxRetries,
		//Transport:  &esTransport{},
		Header: BuildAuthHeader(conf),
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
			// 高并发时可调大
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			MaxConnsPerHost:     100,
			IdleConnTimeout:     90 * time.Second,
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}

	client, err := es7.NewClient(c)
	if err != nil {
		return nil, err
	}

	return &Es{
		Client: client,
	}, nil
}

func MustNewEs(conf *Config) *Es {
	es, err := NewEs(conf)
	if err != nil {
		panic(err)
	}

	return es
}
