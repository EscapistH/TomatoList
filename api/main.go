// main.go
package main

import (
	"log"

	"TomatoList/controllers"
	"TomatoList/database"
	"TomatoList/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	database.InitDatabase()
	db := database.GetDB()

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode) // 生产环境使用ReleaseMode

	// 创建Gin路由器
	router := gin.Default()

	// 中间件
	router.Use(middleware.CORS())                 // CORS中间件
	router.Use(middleware.DatabaseMiddleware(db)) // 数据库中间件

	// API路由
	api := router.Group("/api")
	{
		// 认证路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
		}

		// 需要认证的路由
		authorized := api.Group("/")
		authorized.Use(middleware.JWTAuth()) // JWT认证中间件
		{
			// 任务路由
			authorized.GET("/tasks", controllers.GetTasks)
			authorized.GET("/tasks/:id", controllers.GetTask)
			authorized.POST("/tasks", controllers.CreateTask)
			authorized.PUT("/tasks/:id", controllers.UpdateTask)
			authorized.DELETE("/tasks/:id", controllers.DeleteTask)

			// 番茄钟路由
			authorized.POST("/pomodoros", controllers.StartPomodoro)
			authorized.POST("/pomodoros/:id/complete", controllers.CompletePomodoro)
			authorized.GET("/pomodoros", controllers.GetPomodoros)
			authorized.GET("/pomodoros/stats", controllers.GetPomodoroStats)
		}
	}

	// 启动服务器
	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
