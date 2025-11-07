package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Chateaubriand-g/bili/auth_service/config"
	"github.com/gin-gonic/gin"

	zipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	"github.com/openzipkin/zipkin-go/reporter"
	repoterhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

func InitZipkin(cfg *config.Config) (*zipkin.Tracer, reporter.Reporter, error) {
	// http上报器，指定zipkin服务器地址
	reporter := repoterhttp.NewReporter(cfg.Zipkin.URL)

	// 创建本地服务端点，标识当前服务+地址，地址为空则自动获取
	endpoint, err := zipkin.NewEndpoint(cfg.Zipkin.ServiceName, "")
	if err != nil {
		reporter.Close()
		return nil, nil, fmt.Errorf("zipkin newendpoint failed: %w", err)
	}

	// 创建采样器，指定采样频率，0.1表示10%
	sampler, err := zipkin.NewCountingSampler(cfg.Zipkin.SampleRate)
	if err != nil {
		reporter.Close()
		return nil, nil, fmt.Errorf("zipkin new sampler failed: %w", err)
	}

	tracer, err := zipkin.NewTracer(
		reporter,
		zipkin.WithLocalEndpoint(endpoint),
		zipkin.WithSampler(sampler),
		//zipkin.WithTraceID128Bit(true) 开启128位traceID，默认64位
	)
	if err != nil {
		reporter.Close()
		return nil, nil, fmt.Errorf("zipkin new tracer failed: %w", err)
	}

	return tracer, reporter, nil
}

func ZipkinMiddleware(tracer *zipkin.Tracer) gin.HandlerFunc {
	serverMiddleware := zipkinhttp.NewServerMiddleware(
		tracer,
		zipkinhttp.TagResponseSize(true),
	)

	return func(c *gin.Context) {
		handler := serverMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Request = r
			start := time.Now()

			c.Next()

			status := c.Writer.Status()
			if span := zipkin.SpanFromContext(r.Context()); span != nil {
				zipkin.TagHTTPStatusCode.Set(span, http.StatusText(status))
				if status >= 400 {
					zipkin.TagError.Set(span, "true")
				}
				span.Annotate(time.Now(), "reuqest_finished")
			}

			_ = start
		}))

		handler.ServeHTTP(c.Writer, c.Request)
	}
}

func CloseZipkin(reporter reporter.Reporter) {
	if reporter != nil {
		reporter.Close()
	}
}
