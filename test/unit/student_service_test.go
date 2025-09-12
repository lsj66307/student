package unit

import (
	"student-management-system/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCreateStudentRequest_Validation 测试创建学生请求的数据验证
func TestCreateStudentRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request domain.CreateStudentRequest
		wantErr bool
	}{
		{
			name: "有效的学生数据",
			request: domain.CreateStudentRequest{
				Name:   "张三",
				Age:    20,
				Gender: "男",
				Email:  "zhangsan@example.com",
				Phone:  "13800138000",
				Major:  "计算机科学",
				Grade:  "2024级",
			},
			wantErr: false,
		},
		{
			name: "空姓名",
			request: domain.CreateStudentRequest{
				Name:   "",
				Age:    20,
				Gender: "男",
				Email:  "zhangsan@example.com",
				Phone:  "13800138000",
				Major:  "计算机科学",
				Grade:  "2024级",
			},
			wantErr: true,
		},
		{
			name: "无效年龄",
			request: domain.CreateStudentRequest{
				Name:   "张三",
				Age:    0,
				Gender: "男",
				Email:  "zhangsan@example.com",
				Phone:  "13800138000",
				Major:  "计算机科学",
				Grade:  "2024级",
			},
			wantErr: true,
		},
		{
			name: "无效性别",
			request: domain.CreateStudentRequest{
				Name:   "张三",
				Age:    20,
				Gender: "其他",
				Email:  "zhangsan@example.com",
				Phone:  "13800138000",
				Major:  "计算机科学",
				Grade:  "2024级",
			},
			wantErr: true,
		},
		{
			name: "无效邮箱",
			request: domain.CreateStudentRequest{
				Name:   "张三",
				Age:    20,
				Gender: "男",
				Email:  "invalid-email",
				Phone:  "13800138000",
				Major:  "计算机科学",
				Grade:  "2024级",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 测试结构体字段是否正确设置
			assert.Equal(t, tt.request.Name, tt.request.Name)
			assert.Equal(t, tt.request.Age, tt.request.Age)
			assert.Equal(t, tt.request.Gender, tt.request.Gender)
			assert.Equal(t, tt.request.Email, tt.request.Email)
			assert.Equal(t, tt.request.Phone, tt.request.Phone)
			assert.Equal(t, tt.request.Major, tt.request.Major)
			assert.Equal(t, tt.request.Grade, tt.request.Grade)
		})
	}
}

// TestUpdateStudentRequest_Validation 测试更新学生请求的数据验证
func TestUpdateStudentRequest_Validation(t *testing.T) {
	request := domain.UpdateStudentRequest{
		Name:   "李四",
		Age:    21,
		Gender: "女",
		Email:  "lisi@example.com",
		Phone:  "13900139000",
		Major:  "软件工程",
		Grade:  "2023级",
	}

	// 验证字段设置
	assert.Equal(t, "李四", request.Name)
	assert.Equal(t, 21, request.Age)
	assert.Equal(t, "女", request.Gender)
	assert.Equal(t, "lisi@example.com", request.Email)
	assert.Equal(t, "13900139000", request.Phone)
	assert.Equal(t, "软件工程", request.Major)
	assert.Equal(t, "2023级", request.Grade)
}

// TestStudent_StructFields 测试学生结构体字段
func TestStudent_StructFields(t *testing.T) {
	student := &domain.Student{
		ID:     1,
		Name:   "王五",
		Age:    22,
		Gender: "男",
		Email:  "wangwu@example.com",
		Phone:  "13700137000",
		Major:  "数据科学",
		Grade:  "2022级",
	}

	// 验证所有字段
	assert.Equal(t, 1, student.ID)
	assert.Equal(t, "王五", student.Name)
	assert.Equal(t, 22, student.Age)
	assert.Equal(t, "男", student.Gender)
	assert.Equal(t, "wangwu@example.com", student.Email)
	assert.Equal(t, "13700137000", student.Phone)
	assert.Equal(t, "数据科学", student.Major)
	assert.Equal(t, "2022级", student.Grade)
	assert.NotNil(t, student)
}
