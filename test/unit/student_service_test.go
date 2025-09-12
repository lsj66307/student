package unit

import (
	"testing"
	"student-management-system/internal/domain"
	"student-management-system/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStudentRepository 模拟学生仓储
type MockStudentRepository struct {
	mock.Mock
}

func (m *MockStudentRepository) Create(student *domain.Student) error {
	args := m.Called(student)
	return args.Error(0)
}

func (m *MockStudentRepository) GetByID(id int) (*domain.Student, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Student), args.Error(1)
}

func (m *MockStudentRepository) GetAll() ([]*domain.Student, error) {
	args := m.Called()
	return args.Get(0).([]*domain.Student), args.Error(1)
}

func (m *MockStudentRepository) Update(student *domain.Student) error {
	args := m.Called(student)
	return args.Error(0)
}

func (m *MockStudentRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// TestStudentService_CreateStudent 测试创建学生
func TestStudentService_CreateStudent(t *testing.T) {
	// 准备测试数据
	req := &domain.CreateStudentRequest{
		Name:    "张三",
		Age:     20,
		Gender:  "男",
		Email:   "zhangsan@example.com",
		Phone:   "13800138000",
		Address: "北京市朝阳区",
	}

	// 创建模拟仓储
	mockRepo := new(MockStudentRepository)
	mockRepo.On("Create", mock.AnythingOfType("*domain.Student")).Return(nil)

	// 创建服务实例
	studentService := service.NewStudentService(mockRepo)

	// 执行测试
	student, err := studentService.CreateStudent(req)

	// 断言结果
	assert.NoError(t, err)
	assert.NotNil(t, student)
	assert.Equal(t, req.Name, student.Name)
	assert.Equal(t, req.Age, student.Age)
	assert.Equal(t, req.Gender, student.Gender)
	assert.Equal(t, req.Email, student.Email)
	assert.Equal(t, req.Phone, student.Phone)
	assert.Equal(t, req.Address, student.Address)

	// 验证模拟调用
	mockRepo.AssertExpectations(t)
}

// TestStudentService_GetStudentByID 测试根据ID获取学生
func TestStudentService_GetStudentByID(t *testing.T) {
	// 准备测试数据
	expectedStudent := &domain.Student{
		ID:      1,
		Name:    "张三",
		Age:     20,
		Gender:  "男",
		Email:   "zhangsan@example.com",
		Phone:   "13800138000",
		Address: "北京市朝阳区",
	}

	// 创建模拟仓储
	mockRepo := new(MockStudentRepository)
	mockRepo.On("GetByID", 1).Return(expectedStudent, nil)

	// 创建服务实例
	studentService := service.NewStudentService(mockRepo)

	// 执行测试
	student, err := studentService.GetStudentByID(1)

	// 断言结果
	assert.NoError(t, err)
	assert.NotNil(t, student)
	assert.Equal(t, expectedStudent.ID, student.ID)
	assert.Equal(t, expectedStudent.Name, student.Name)
	assert.Equal(t, expectedStudent.Age, student.Age)

	// 验证模拟调用
	mockRepo.AssertExpectations(t)
}

// TestStudentService_GetAllStudents 测试获取所有学生
func TestStudentService_GetAllStudents(t *testing.T) {
	// 准备测试数据
	expectedStudents := []*domain.Student{
		{
			ID:      1,
			Name:    "张三",
			Age:     20,
			Gender:  "男",
			Email:   "zhangsan@example.com",
			Phone:   "13800138000",
			Address: "北京市朝阳区",
		},
		{
			ID:      2,
			Name:    "李四",
			Age:     21,
			Gender:  "女",
			Email:   "lisi@example.com",
			Phone:   "13900139000",
			Address: "上海市浦东新区",
		},
	}

	// 创建模拟仓储
	mockRepo := new(MockStudentRepository)
	mockRepo.On("GetAll").Return(expectedStudents, nil)

	// 创建服务实例
	studentService := service.NewStudentService(mockRepo)

	// 执行测试
	students, err := studentService.GetAllStudents()

	// 断言结果
	assert.NoError(t, err)
	assert.NotNil(t, students)
	assert.Len(t, students, 2)
	assert.Equal(t, expectedStudents[0].Name, students[0].Name)
	assert.Equal(t, expectedStudents[1].Name, students[1].Name)

	// 验证模拟调用
	mockRepo.AssertExpectations(t)
}

// TestStudentService_UpdateStudent 测试更新学生
func TestStudentService_UpdateStudent(t *testing.T) {
	// 准备测试数据
	existingStudent := &domain.Student{
		ID:      1,
		Name:    "张三",
		Age:     20,
		Gender:  "男",
		Email:   "zhangsan@example.com",
		Phone:   "13800138000",
		Address: "北京市朝阳区",
	}

	updateReq := &domain.UpdateStudentRequest{
		Name:    "张三丰",
		Age:     21,
		Gender:  "男",
		Email:   "zhangsanfeng@example.com",
		Phone:   "13800138001",
		Address: "北京市海淀区",
	}

	// 创建模拟仓储
	mockRepo := new(MockStudentRepository)
	mockRepo.On("GetByID", 1).Return(existingStudent, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.Student")).Return(nil)

	// 创建服务实例
	studentService := service.NewStudentService(mockRepo)

	// 执行测试
	student, err := studentService.UpdateStudent(1, updateReq)

	// 断言结果
	assert.NoError(t, err)
	assert.NotNil(t, student)
	assert.Equal(t, updateReq.Name, student.Name)
	assert.Equal(t, updateReq.Age, student.Age)
	assert.Equal(t, updateReq.Email, student.Email)

	// 验证模拟调用
	mockRepo.AssertExpectations(t)
}