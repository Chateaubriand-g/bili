package proxy

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	//"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	zipkin "github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	b3 "github.com/openzipkin/zipkin-go/propagation/b3"
)

var proxyCache = sync.Map{}

func ReverseProxy(cli *api.Client, serviceName string, tracer *zipkin.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		serv, err := PickService(cli, serviceName)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"msg": "service unavailable", "error": err.Error()})
			return
		}

		target := fmt.Sprintf("http://%s:%d", serv.Address, serv.Port)
		proxy, exists := proxyCache.Load(target)
		if !exists {
			url, err := url.Parse(target)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "invalid service address"})
				return
			}
			newProxy := newZipkinReverseProxy(url, tracer)
			// newProxy.ModifyResponse = func(resp *http.Response) error {
			// 	if resp.StatusCode >= 500 {
			// 		body, err := io.ReadAll(resp.Body)
			// 		if err != nil {
			// 			return err
			// 		}
			// 		newBody := []byte(fmt.Sprintf(`{"msg":"server internal error: %s"}`, string(body)))
			// 		resp.Body = io.NopCloser(bytes.NewBuffer(newBody))
			// 		resp.ContentLength = int64(len(newBody))
			// 		resp.Header.Set("Content-Type", "application/json")
			// 	}
			// 	return nil
			// }
			proxyCache.Store(target, newProxy)
			proxy = newProxy
		}

		//ioriginalPath := c.Request.URL.Path
		//trimPath := strings.TrimPrefix(originalPath, "/api")
		//if trimPath == "" {
		//	trimPath = "/"
		//}
		//c.Request.URL.Path = trimPath

		log.Printf("[TRACE] %s - Final path for forwarding: %s", 123123, c.Request.URL.Path)

		proxy.(*httputil.ReverseProxy).ServeHTTP(c.Writer, c.Request)
	}
}

// newZipkinReverseProxy 创建带 Zipkin 链路追踪功能的反向代理
func newZipkinReverseProxy(target *url.URL, tracer *zipkin.Tracer) *httputil.ReverseProxy {
	// 创建带追踪能力的 Zipkin HTTP 客户端（自动注入 B3 头）
	// client, err := zipkinhttp.NewClient(tracer, zipkinhttp.ClientTrace(true))
	// if err != nil {
	// 	log.Fatalf("failed to create zipkin http client: %v", err)
	// }

	proxy := httputil.NewSingleHostReverseProxy(target)
	originalDirector := proxy.Director

	// 请求转发前的预处理
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		// 从请求上下文中获取父 span（如果 Gin 已经有 zipkin 中间件）
		parentSpan := zipkin.SpanFromContext(req.Context())

		// 为本次代理请求创建新的客户端 span
		var span zipkin.Span
		if parentSpan != nil {
			span = tracer.StartSpan(
				"proxy_request",
				zipkin.Parent(parentSpan.Context()),
				zipkin.Kind(model.Client),
			)
		} else {
			span = tracer.StartSpan("proxy_request", zipkin.Kind(model.Client))
		}

		zipkin.TagHTTPMethod.Set(span, req.Method)
		zipkin.TagHTTPUrl.Set(span, req.URL.String())
		zipkin.TagHTTPRoute.Set(span, req.URL.Path)
		const TagPeerHost zipkin.Tag = "peer.host"
		zipkin.Tag(TagPeerHost).Set(span, target.Host)

		// 将 span 放入请求上下文
		newreq := req.WithContext(zipkin.NewContext(req.Context(), span))

		// 注入B3追踪头（上游→代理→下游服务）
		if err := b3.InjectHTTP(newreq)(span.Context()); err != nil {
			span.Tag("error", fmt.Sprintf("注入B3头失败: %v", err))
			log.Printf("[TRACE] 注入B3头失败: %v", err)
		}

		req.Header = newreq.Header
		*req = *newreq
	}

	// ModifyResponse: 处理响应与 span 结束逻辑
	proxy.ModifyResponse = func(resp *http.Response) error {
		// 从响应中获取span并结束
		clientSpan := zipkin.SpanFromContext(resp.Request.Context())
		if clientSpan == nil {
			return nil
		}
		defer clientSpan.Finish()

		statusCode := resp.StatusCode
		clientSpan.Tag("http.status_code", fmt.Sprintf("%d", statusCode))
		zipkin.TagHTTPResponseSize.Set(clientSpan, strconv.FormatInt(resp.ContentLength, 10))

		// 处理下游错误，包装响应体
		if statusCode >= 500 {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				clientSpan.Tag("error", "读取下游错误响应失败: "+err.Error())
				return err
			}
			defer resp.Body.Close()

			newBody := []byte(fmt.Sprintf(`{"msg":"server internal error: %s"}`, string(body)))
			resp.Body = io.NopCloser(bytes.NewBuffer(newBody))
			resp.ContentLength = int64(len(newBody))
			resp.Header.Set("Content-Type", "application/json")

			clientSpan.Tag("error", fmt.Sprintf("下游服务错误: %d, %s", statusCode, string(body)))
		} else if statusCode >= 400 {
			clientSpan.Tag("error", fmt.Sprintf("下游服务错误: %d", statusCode))
		}

		return nil
	}

	// ErrorHandler: 捕获网络错误并结束 span
	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		rw.Header().Set("Content-Type", "application/json;charset=utf-8")
		rw.WriteHeader(http.StatusBadGateway)
		_, _ = rw.Write([]byte(fmt.Sprintf(`{"code":502,"msg":"代理转发失败","err":"%v"}`, err)))
	}

	return proxy
}
