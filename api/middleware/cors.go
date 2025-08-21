// middleware/cors.go
package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORS 跨域中间件
// 与Python FastAPI对比：
// # from fastapi.middleware.cors import CORSMiddleware
// #
// # app.add_middleware(
// #     CORSMiddleware,
// #     allow_origins=["*"],
// #     allow_credentials=True,
// #     allow_methods=["*"],
// #     allow_headers=["*"],
// # )

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许所有源
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		// 允许携带凭证
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 允许的请求头
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")

		// 允许的HTTP方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // 无内容响应
			return
		}

		c.Next()
	}
}
