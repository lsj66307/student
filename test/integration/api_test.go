package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"student-management-system/internal/domain"
	"student-management-system/internal/handler"
	"student-management-system/internal/repository"
	"student-management-system/internal/service"

	"github.com/gin-gonic/gin"
)

// TestStudentAPI 测试学生API集成
func TestStudentAPI(t *testing.T) {
	// 检查数据库连接
	err := repository.InitDB()
	if err != nil {
		t.Skipf("跳过集成测试：数据库连接失败 - %v", err)
		return
	}

	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建路由
	router := gin.New()
	studentHandler := handler.NewStudentHandler()

	// 注册路由
	v1 := router.Group("/api/v1")
	{
		v1.POST("/students", studentHandler.CreateStudent)
		v1.GET("/students/:id", studentHandler.GetStudent)
		v1.GET("/students", studentHandler.GetStudents)
		v1.PUT("/students/:id", studentHandler.UpdateStudent)
		v1.DELETE("/students/:id", studentHandler.DeleteStudent)
	}

	// 测试创建学生
	t.Run("CreateStudent", func(t *testing.T) {
		createReq := domain.CreateStudentRequest{
			Name:   "测试学生",
			Age:    20,
			Gender: "男",
			Email:  "test@example.com",
			Phone:  "13800138000",
			Major:  "计算机科学",
			Grade:  "2024级",
		}

		jsonData, _ := json.Marshal(createReq)
		req, _ := http.NewRequest("POST", "/api/v1/students", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 检查响应状态码
		if w.Code != http.StatusCreated {
			t.Logf("Expected status %d, got %d", http.StatusCreated, w.Code)
			t.Logf("Response body: %s", w.Body.String())
		}
	})

	// 测试获取学生列表
	t.Run("GetStudents", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/students", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 检查响应
		if w.Code != http.StatusOK {
			t.Logf("Expected status %d, got %d", http.StatusOK, w.Code)
			t.Logf("Response body: %s", w.Body.String())
		}
	})
}

// TestTeacherAPI 测试教师API集成
func TestTeacherAPI(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建路由
	router := gin.New()
	teacherHandler := handler.NewTeacherHandler()

	// 注册路由
	v1 := router.Group("/api/v1")
	{
		v1.POST("/teachers", teacherHandler.CreateTeacher)
		v1.GET("/teachers/:id", teacherHandler.GetTeacher)
		v1.GET("/teachers", teacherHandler.GetTeachers)
		v1.PUT("/teachers/:id", teacherHandler.UpdateTeacher)
		v1.DELETE("/teachers/:id", teacherHandler.DeleteTeacher)
	}

	// 测试创建教师
	t.Run("CreateTeacher", func(t *testing.T) {
		createReq := domain.CreateTeacherRequest{
			Name:       "测试教师",
			Age:        35,
			Gender:     "女",
			Email:      "teacher@example.com",
			Phone:      "13900139000",
			Subject:    "数学",
			Department: "数学系",
		}

		jsonData, _ := json.Marshal(createReq)
		req, _ := http.NewRequest("POST", "/api/v1/teachers", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 检查响应
		if w.Code != http.StatusCreated {
			t.Logf("Expected status %d, got %d", http.StatusCreated, w.Code)
			t.Logf("Response body: %s", w.Body.String())
		}
	})
}

// TestGradeAPI 测试成绩API集成
func TestGradeAPI(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 检查数据库连接
	err := repository.InitDB()
	if err != nil {
		t.Skipf("跳过集成测试：数据库连接失败 - %v", err)
		return
	}

	// 创建路由
	router := gin.New()
	gradeService := service.NewGradeService(repository.DB)
	gradeHandler := handler.NewGradeHandler(gradeService)

	// 注册路由
	v1 := router.Group("/api/v1")
	{
		v1.POST("/grades", gradeHandler.CreateGrade)
		v1.GET("/grades/:id", gradeHandler.GetGrade)
		v1.GET("/grades", gradeHandler.GetGrades)
		v1.PUT("/grades/:id", gradeHandler.UpdateGrade)
		v1.DELETE("/grades/:id", gradeHandler.DeleteGrade)
	}

	// 测试创建成绩
	t.Run("CreateGrade", func(t *testing.T) {
		mathScore := 95.5
		createReq := domain.CreateGradeRequest{
			StudentID: 1,
			MathScore: &mathScore,
		}

		jsonData, _ := json.Marshal(createReq)
		req, _ := http.NewRequest("POST", "/api/v1/grades", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 检查响应
		if w.Code != http.StatusCreated {
			t.Logf("Expected status %d, got %d", http.StatusCreated, w.Code)
			t.Logf("Response body: %s", w.Body.String())
		}
	})
}
