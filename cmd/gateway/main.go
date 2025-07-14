package main

import (
	"fmt"
	"github.com/teakingwang/ginmicro/config"
	"github.com/teakingwang/ginmicro/pkg/logger"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := logger.Init(true); err != nil {
		panic("logger init failed: " + err.Error())
	}

	if err := config.LoadConfigFromConsul("config/user"); err != nil {
		logger.Warn("load from consul failed: %v", err)
		if err := config.LoadConfig(); err != nil {
			panic(fmt.Errorf("failed to load config: %v", err))
		}
	}

	if err := config.LoadConfigFromConsul("config/order"); err != nil {
		logger.Warn("load from consul failed: %v, falling back to local config", err)
		if err := config.LoadConfig(); err != nil {
			panic(fmt.Errorf("failed to load config: %v", err))
		}
	}

	r := gin.Default()

	// 添加 JWT 验证中间件，跳过登录/注册接口
	//r.Use(middleware.JWTGinMiddleware())

	// 转发到 user-service (监听 50051)/
	userTarget := "http://user:" + config.Config.Server.User.HTTPPort
	r.Any("/v1/user/*proxyPath", ReverseProxy(userTarget))

	// 转发到 order-service (监听 50052)
	orderTarget := "http://order:" + config.Config.Server.Order.HTTPPort
	r.Any("/v1/order/*proxyPath", ReverseProxy(orderTarget))

	log.Println("🚪 Gateway listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start gateway: %v", err)
	}
}

func ReverseProxy(target string) gin.HandlerFunc {
	targetURL, err := url.Parse(target)
	if err != nil {
		panic("Invalid proxy target: " + target)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// 替换 Director 保留原始路径
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.Host = targetURL.Host
		// 保留原始请求路径，例如 /v1/order/xxx
		req.URL.Path = req.URL.Path
		// 保留查询参数
		req.URL.RawQuery = req.URL.RawQuery
	}

	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// 保证路径拼接时不会重复斜杠
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
