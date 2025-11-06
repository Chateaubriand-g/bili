package proxy

import (
	"bytes"
	"log"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	//"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
)

var proxyCache = sync.Map{}

func ReverseProxy(cli *api.Client, serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		serv, err := PickService(cli, serviceName)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"msg": "service unavailable","error":err.Error()})
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
			newProxy := httputil.NewSingleHostReverseProxy(url)
			newProxy.ModifyResponse = func(resp *http.Response) error {
				if resp.StatusCode >= 500 {
					body, err := io.ReadAll(resp.Body)
					if err != nil {
						return err
					}
					newBody := []byte(fmt.Sprintf(`{"msg":"server internal error: %s"}`, string(body)))
					resp.Body = io.NopCloser(bytes.NewBuffer(newBody))
					resp.ContentLength = int64(len(newBody))
					resp.Header.Set("Content-Type", "application/json")
				}
				return nil
			}
			proxyCache.Store(target, newProxy)
			proxy = newProxy
		}

		//ioriginalPath := c.Request.URL.Path
		//trimPath := strings.TrimPrefix(originalPath, "/api")
		//if trimPath == "" {
		//	trimPath = "/"
		//}
		//c.Request.URL.Path = trimPath
		
		log.Printf("[TRACE] %s - Final path for forwarding: %s",123123,c.Request.URL.Path)

		proxy.(*httputil.ReverseProxy).ServeHTTP(c.Writer, c.Request)
	}
}
