package helpers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"student-management-system/internal/domain"

	"github.com/gin-gonic/gin"
)

// SetupTestRouter 设置测试路由
func SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

// CreateTestStudent 创建测试学生数据
func CreateTestStudent(id int) *domain.Student {
	return &domain.Student{
		ID:     id,
		Name:   "测试学生",
		Age:    20,
		Gender: "男",
		Email:  "test@example.com",
		Phone:  "13800138000",
	}
}

// CreateTestTeacher 创建测试教师数据
func CreateTestTeacher(id int) *domain.Teacher {
	return &domain.Teacher{
		ID:         id,
		Name:       "测试教师",
		Age:        35,
		Gender:     "女",
		Email:      "teacher@example.com",
		Phone:      "13900139000",
		Subject:    "数学",
		Department: "数学系",
	}
}

// CreateTestGrade 创建测试成绩数据
func CreateTestGrade(id, studentID int) *domain.Grade {
	mathScore := 95.5
	return &domain.Grade{
		ID:        id,
		StudentID: studentID,
		MathScore: &mathScore,
	}
}

// CreateTestCreateStudentRequest 创建测试学生创建请求
func CreateTestCreateStudentRequest() *domain.CreateStudentRequest {
	return &domain.CreateStudentRequest{
		Name:   "新学生",
		Age:    18,
		Gender: "女",
		Email:  "newstudent@example.com",
		Phone:  "13700137000",
	}
}

// CreateTestCreateTeacherRequest 创建测试教师创建请求
func CreateTestCreateTeacherRequest() *domain.CreateTeacherRequest {
	return &domain.CreateTeacherRequest{
		Name:       "新教师",
		Age:        30,
		Gender:     "男",
		Email:      "newteacher@example.com",
		Phone:      "13600136000",
		Subject:    "英语",
		Department: "外语系",
	}
}

// CreateTestCreateGradeRequest 创建测试成绩创建请求
func CreateTestCreateGradeRequest(studentID int) *domain.CreateGradeRequest {
	mathScore := 88.0
	englishScore := 92.5
	return &domain.CreateGradeRequest{
		StudentID:    studentID,
		MathScore:    &mathScore,
		EnglishScore: &englishScore,
	}
}

// MakeJSONRequest 创建JSON请求
func MakeJSONRequest(method, url string, body interface{}) (*http.Request, error) {
	var bodyReader *strings.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = strings.NewReader(string(jsonData))
	} else {
		bodyReader = strings.NewReader("")
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// PerformRequest 执行HTTP请求
func PerformRequest(router *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// AssertStatusCode 断言状态码
func AssertStatusCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected status code %d, got %d", expected, actual)
	}
}

// AssertJSONResponse 断言JSON响应
func AssertJSONResponse(t *testing.T, w *httptest.ResponseRecorder, expected interface{}) {
	var actual interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &actual); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
		return
	}

	expectedJSON, _ := json.Marshal(expected)
	actualJSON, _ := json.Marshal(actual)

	if string(expectedJSON) != string(actualJSON) {
		t.Errorf("Expected JSON %s, got %s", expectedJSON, actualJSON)
	}
}

// LogResponse 记录响应内容（用于调试）
func LogResponse(t *testing.T, w *httptest.ResponseRecorder) {
	t.Logf("Response Status: %d", w.Code)
	t.Logf("Response Headers: %v", w.Header())
	t.Logf("Response Body: %s", w.Body.String())
}

// CreateTestStudents 批量创建测试学生
func CreateTestStudents(count int) []*domain.Student {
	students := make([]*domain.Student, count)
	for i := 0; i < count; i++ {
		students[i] = CreateTestStudent(i + 1)
	}
	return students
}

// CreateTestTeachers 批量创建测试教师
func CreateTestTeachers(count int) []*domain.Teacher {
	teachers := make([]*domain.Teacher, count)
	for i := 0; i < count; i++ {
		teachers[i] = CreateTestTeacher(i + 1)
	}
	return teachers
}

// CreateTestGrades 批量创建测试成绩
func CreateTestGrades(count, studentID int) []*domain.Grade {
	grades := make([]*domain.Grade, count)
	for i := 0; i < count; i++ {
		grades[i] = CreateTestGrade(i+1, studentID)
	}
	return grades
}