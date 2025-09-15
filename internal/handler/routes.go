package handler

import (
	"student-management-system/internal/config"
	"student-management-system/internal/repository"
	"student-management-system/internal/service"
	"student-management-system/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置所有路由
func SetupRoutes(cfg *config.Config) *gin.Engine {
	// 创建Gin引擎
	router := gin.Default()

	// 添加CORS中间件
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 创建服务和处理器
	studentHandler := NewStudentHandler()
	teacherHandler := NewTeacherHandler()

	// 创建成绩服务和处理器
	gradeService := service.NewGradeService(repository.DB)
	gradeHandler := NewGradeHandler(gradeService)

	// 创建认证服务和处理器
	authService := service.NewAuthService(cfg)
	authHandler := NewAuthHandler(authService)

	// API路由组
	api := router.Group("/api/v1")
	{
		// 认证相关路由（无需认证）
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)           // 管理员登录
			auth.POST("/validate", authHandler.ValidateToken) // 验证token
		}

		// 需要认证的路由组
		protected := api.Group("")
		protected.Use(middleware.JWTAuth()) // 应用JWT认证中间件
		{
			// 认证用户信息路由
			protected.GET("/auth/profile", authHandler.GetProfile)     // 获取当前管理员信息
			protected.POST("/auth/refresh", authHandler.RefreshToken) // 刷新token

			// 学生相关路由（需要认证）
			students := protected.Group("/students")
			{
				students.POST("", studentHandler.CreateStudent)       // 创建学生
				students.GET("", studentHandler.GetStudents)          // 获取学生列表
				students.GET("/:id", studentHandler.GetStudent)       // 获取单个学生
				students.PUT("/:id", studentHandler.UpdateStudent)    // 更新学生
				students.DELETE("/:id", studentHandler.DeleteStudent) // 删除学生
			}

			// 老师相关路由（需要认证）
			teachers := protected.Group("/teachers")
			{
				teachers.POST("", teacherHandler.CreateTeacher)       // 创建老师
				teachers.GET("", teacherHandler.GetTeachers)          // 获取老师列表
				teachers.GET("/:id", teacherHandler.GetTeacher)       // 获取单个老师
				teachers.PUT("/:id", teacherHandler.UpdateTeacher)    // 更新老师
				teachers.DELETE("/:id", teacherHandler.DeleteTeacher) // 删除老师
			}

			// 成绩相关路由（需要认证）
			grades := protected.Group("/grades")
			{
				grades.POST("", gradeHandler.CreateGrade)       // 创建成绩
				grades.GET("", gradeHandler.GetGrades)          // 获取成绩列表
				grades.GET("/:id", gradeHandler.GetGrade)       // 获取单个成绩
				grades.PUT("/:id", gradeHandler.UpdateGrade)    // 更新成绩
				grades.DELETE("/:id", gradeHandler.DeleteGrade) // 删除成绩
			}
		}
	}

	// 健康检查路由
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "学生管理系统运行正常",
		})
	})

	// 根路径
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "欢迎使用学生管理系统API",
			"version": "1.0.0",
			"endpoints": gin.H{
				"health":         "GET /health",
				"login":          "POST /api/v1/auth/login",
				"profile":        "GET /api/v1/auth/profile (需要认证)",
				"create_student": "POST /api/v1/students (需要认证)",
				"get_students":   "GET /api/v1/students (需要认证)",
				"get_student":    "GET /api/v1/students/{id} (需要认证)",
				"update_student": "PUT /api/v1/students/{id} (需要认证)",
				"delete_student": "DELETE /api/v1/students/{id} (需要认证)",
				"create_teacher": "POST /api/v1/teachers (需要认证)",
				"get_teachers":   "GET /api/v1/teachers (需要认证)",
				"create_grade":   "POST /api/v1/grades (需要认证)",
				"get_grades":     "GET /api/v1/grades (需要认证)",
			},
		})
	})

	return router
}
