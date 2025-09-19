package handler

import (
	"student-management-system/internal/config"
	"student-management-system/internal/repository"
	"student-management-system/internal/service"
	"student-management-system/pkg/logger"
	"student-management-system/pkg/middleware"
	"student-management-system/pkg/validator"

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

	// 创建验证器实例
	customValidator := validator.NewValidator()

	// 获取logger实例
	loggerInstance := logger.GetLogger().Logger

	// 创建Repository实例
	adminRepo := repository.NewAdminRepository(repository.DB, loggerInstance)

	// 创建服务实例
	authService := service.NewAuthService(cfg, adminRepo)
	studentService := service.NewStudentService()
	teacherService := service.NewTeacherService()
	gradeService := service.NewGradeService(repository.DB)
	adminService := service.NewAdminService(adminRepo, loggerInstance)

	// 创建处理器实例
	authHandler := NewAuthHandler(authService)
	studentHandler := NewStudentHandler(studentService, customValidator)
	teacherHandler := NewTeacherHandler(teacherService)
	gradeHandler := NewGradeHandler(gradeService)
	adminHandler := NewAdminHandler(adminService, loggerInstance)

	// API路由组
	api := router.Group("/api/v1")
	{
		// 认证相关路由（无需认证）
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)            // 管理员登录
			auth.POST("/validate", authHandler.ValidateToken) // 验证token
		}

		// 需要认证的路由组
		protected := api.Group("")
		protected.Use(middleware.JWTAuth()) // 应用JWT认证中间件
		{
			// 认证用户信息路由
			protected.GET("/auth/profile", authHandler.GetProfile)    // 获取当前管理员信息
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

			// 管理员相关路由（需要认证）
			admins := protected.Group("/admins")
			{
				admins.POST("", adminHandler.CreateAdmin)       // 创建管理员
				admins.GET("", adminHandler.ListAdmins)         // 获取管理员列表
				admins.GET("/:id", adminHandler.GetAdmin)       // 获取单个管理员
				admins.PUT("/:id", adminHandler.UpdateAdmin)    // 更新管理员
				admins.DELETE("/:id", adminHandler.DeleteAdmin) // 删除管理员
			}
		}
	}

	return router
}
