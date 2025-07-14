package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/teakingwang/ginmicro/pkg/auth"
)

// 路径白名单（跳过 JWT 验证）
var noAuthPaths = map[string]bool{
	"/v1/user/login":  true,
	"/v1/user/signup": true,
}

func JWTGinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if noAuthPaths[c.FullPath()] {
			c.Next()
			return
		}

		token := extractTokenFromHeader(c.Request)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}

		claims, err := auth.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("user", claims)
		c.Next()
	}
}

func extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	return parts[1]
}
