package handlers

import (
	"github.com/gin-gonic/gin"
	"student-management-system/database"
	"student-management-system/models"
)

// SetupRoutes 设置所有路由
func SetupRoutes() *gin.Engine {
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

	// 创建处理器
	studentHandler := NewStudentHandler()
	teacherHandler := NewTeacherHandler()
	
	// 创建成绩服务和处理器
	gradeService := models.NewGradeService(database.DB)
	gradeHandler := NewGradeHandler(gradeService)

	// API路由组
	api := router.Group("/api/v1")
	{
		// 学生相关路由
		students := api.Group("/students")
		{
			students.POST("", studentHandler.CreateStudent)       // 创建学生
			students.GET("", studentHandler.GetStudents)          // 获取学生列表
			students.GET("/:id", studentHandler.GetStudent)       // 获取单个学生
			students.PUT("/:id", studentHandler.UpdateStudent)    // 更新学生
			students.DELETE("/:id", studentHandler.DeleteStudent) // 删除学生
		}

		// 老师相关路由
		teachers := api.Group("/teachers")
		{
			teachers.POST("", teacherHandler.CreateTeacher)       // 创建老师
			teachers.GET("", teacherHandler.GetTeachers)          // 获取老师列表
			teachers.GET("/:id", teacherHandler.GetTeacher)       // 获取单个老师
			teachers.PUT("/:id", teacherHandler.UpdateTeacher)    // 更新老师
			teachers.DELETE("/:id", teacherHandler.DeleteTeacher) // 删除老师
		}

		// 成绩相关路由
		grades := api.Group("/grades")
		{
			grades.POST("", gradeHandler.CreateGrade)       // 创建成绩
			grades.GET("", gradeHandler.GetGrades)          // 获取成绩列表
			grades.GET("/:id", gradeHandler.GetGrade)       // 获取单个成绩
			grades.PUT("/:id", gradeHandler.UpdateGrade)    // 更新成绩
			grades.DELETE("/:id", gradeHandler.DeleteGrade) // 删除成绩
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
				"create_student": "POST /api/students",
				"get_students":   "GET /api/students",
				"get_student":    "GET /api/students/{id}",
				"update_student": "PUT /api/students/{id}",
				"delete_student": "DELETE /api/students/{id}",
			},
		})
	})

	return router
}
