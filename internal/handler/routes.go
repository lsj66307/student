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
	scoreRepo := repository.NewScoreRepository(repository.DB)

	// 创建服务实例
	authService := service.NewAuthService(cfg, adminRepo)
	studentService := service.NewStudentService()
	teacherService := service.NewTeacherService()
	subjectService := service.NewSubjectService()
	scoreService := service.NewScoreService(scoreRepo)
	adminService := service.NewAdminService(adminRepo, loggerInstance)

	// 创建处理器实例
	authHandler := NewAuthHandler(authService)
	studentHandler := NewStudentHandler(studentService, customValidator)
	teacherHandler := NewTeacherHandler(teacherService)
	subjectHandler := NewSubjectHandler(subjectService, customValidator)
	scoreHandler := NewScoreHandler(scoreService)
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
			protected.POST("/auth/logout", authHandler.Logout)        // 用户登出

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

			// 科目相关路由（需要认证）
			subjects := protected.Group("/subjects")
			{
				subjects.POST("", subjectHandler.CreateSubject)       // 创建科目
				subjects.GET("", subjectHandler.GetSubjects)          // 获取科目列表
				subjects.GET("/:id", subjectHandler.GetSubject)       // 获取单个科目
				subjects.PUT("/:id", subjectHandler.UpdateSubject)    // 更新科目
				subjects.DELETE("/:id", subjectHandler.DeleteSubject) // 删除科目
			}

			// 成绩相关路由（需要认证）
			scores := protected.Group("/scores")
			{
				scores.POST("", scoreHandler.CreateScore)       // 创建成绩
				scores.GET("", scoreHandler.GetScores)          // 获取成绩列表
				scores.GET("/:id", scoreHandler.GetScore)       // 获取单个成绩
				scores.PUT("/:id", scoreHandler.UpdateScore)    // 更新成绩
				scores.DELETE("/:id", scoreHandler.DeleteScore) // 删除成绩
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
