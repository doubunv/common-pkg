package core

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/opensearch-project/opensearch-go"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"io"
	"math/rand"
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
		*opensearch.Client
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

// InitESClient 初始化 OpenSearch 客户端
func InitESClient(conf *Config) (*Es, error) {
	tlsConfig := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true, // 建议生产环境关闭跳过
	}

	transport := &http.Transport{
		TLSClientConfig:     tlsConfig,
		MaxIdleConns:        500,
		MaxIdleConnsPerHost: 500,
		MaxConnsPerHost:     1000,
		IdleConnTimeout:     60 * time.Second,
		DisableKeepAlives:   false,
		TLSHandshakeTimeout: 10 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	headers := BuildAuthHeader(conf)
	cfg := opensearch.Config{
		Addresses:     conf.Addresses,
		Transport:     transport,
		Header:        headers.Clone(), // 确保线程安全
		RetryOnStatus: []int{502, 503, 504},
		MaxRetries:    3,
		RetryBackoff: func(attempt int) time.Duration {
			base := 100 * time.Millisecond
			jitter := time.Duration(rand.Intn(100)) * time.Millisecond
			return (1<<attempt)*base + jitter
		},
	}

	client, err := opensearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating OpenSearch client failed: %w", err)
	}

	res, err := client.Info()
	if err != nil {
		return nil, fmt.Errorf("OpenSearch ping failed: %w", err)
	}
	defer func() {
		defer res.Body.Close()
	}()

	io.Copy(io.Discard, res.Body)
	if res.IsError() {
		return nil, fmt.Errorf("OpenSearch returned error: %s", res.String())
	}

	return &Es{Client: client}, nil
}

func MustNewEs(conf *Config) *Es {
	es, err := InitESClient(conf)
	if err != nil {
		panic(err)
	}

	return es
}
