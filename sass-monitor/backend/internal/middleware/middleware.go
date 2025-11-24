package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter 限流器结构
type RateLimiter struct {
	clients map[string][]time.Time
	mutex   sync.RWMutex
}

// NewRateLimiter 创建新的限流器
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		clients: make(map[string][]time.Time),
	}
}

// IsAllowed 检查是否允许请求
func (rl *RateLimiter) IsAllowed(clientIP string, maxRequests int, window time.Duration) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()

	// 获取或创建客户端请求记录
	requests, exists := rl.clients[clientIP]
	if !exists {
		rl.clients[clientIP] = []time.Time{now}
		return true
	}

	// 清理过期的请求记录
	var validRequests []time.Time
	for _, reqTime := range requests {
		if now.Sub(reqTime) < window {
			validRequests = append(validRequests, reqTime)
		}
	}

	// 检查是否超过限制
	if len(validRequests) >= maxRequests {
		return false
	}

	// 添加当前请求
	validRequests = append(validRequests, now)
	rl.clients[clientIP] = validRequests

	return true
}

// 全局限流器实例
var globalRateLimiter = NewRateLimiter()

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// 检查是否超过限制（每分钟100个请求）
		if !globalRateLimiter.IsAllowed(clientIP, 100, time.Minute) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SecurityMiddleware 安全中间件
func SecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置安全头
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// 检查User-Agent，防止简单的爬虫
		userAgent := c.GetHeader("User-Agent")
		if userAgent == "" || strings.Contains(strings.ToLower(userAgent), "bot") {
			// 对于可疑的User-Agent进行记录，但不阻止
			// 在生产环境中可能需要更严格的检查
		}

		c.Next()
	}
}

// CORSMiddleware 自定义CORS中间件
func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 检查是否在允许的源列表中
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "43200") // 12 hours

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// LoggingMiddleware 日志中间件
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %s %s %d %s %s\n",
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
		)
	})
}

// ErrorHandlingMiddleware 错误处理中间件
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}