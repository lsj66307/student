package benchmark

import (
	"testing"
	"student-management-system/internal/domain"
)

// BenchmarkStudentCreation 基准测试学生创建性能
func BenchmarkStudentCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		student := &domain.Student{
			ID:     i,
			Name:   "测试学生",
			Age:    20,
			Gender: "男",
			Email:  "test@example.com",
			Phone:  "13800138000",
		}
		_ = student // 避免编译器优化
	}
}

// BenchmarkTeacherCreation 基准测试教师创建性能
func BenchmarkTeacherCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		teacher := &domain.Teacher{
			ID:         i,
			Name:       "测试教师",
			Age:        35,
			Gender:     "女",
			Email:      "teacher@example.com",
			Phone:      "13900139000",
			Subject:    "数学",
			Department: "数学系",
		}
		_ = teacher // 避免编译器优化
	}
}

// BenchmarkGradeCreation 基准测试成绩创建性能
func BenchmarkGradeCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mathScore := 95.5
		grade := &domain.Grade{
			ID:        i,
			StudentID: 1,
			MathScore: &mathScore,
		}
		_ = grade // 避免编译器优化
	}
}

// BenchmarkSliceOperations 基准测试切片操作性能
func BenchmarkSliceOperations(b *testing.B) {
	b.Run("AppendStudents", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var students []*domain.Student
			for j := 0; j < 100; j++ {
				student := &domain.Student{
					ID:     j,
					Name:   "学生",
					Age:    20,
					Gender: "男",
					Email:  "test@example.com",
					Phone:  "13800138000",
				}
				students = append(students, student)
			}
			_ = students
		}
	})

	b.Run("PreallocatedStudents", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			students := make([]*domain.Student, 0, 100)
			for j := 0; j < 100; j++ {
				student := &domain.Student{
					ID:     j,
					Name:   "学生",
					Age:    20,
					Gender: "男",
					Email:  "test@example.com",
					Phone:  "13800138000",
				}
				students = append(students, student)
			}
			_ = students
		}
	})
}

// BenchmarkMapOperations 基准测试映射操作性能
func BenchmarkMapOperations(b *testing.B) {
	b.Run("StudentMap", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			studentMap := make(map[int]*domain.Student)
			for j := 0; j < 100; j++ {
				student := &domain.Student{
					ID:     j,
					Name:   "学生",
					Age:    20,
					Gender: "男",
					Email:  "test@example.com",
					Phone:  "13800138000",
				}
				studentMap[j] = student
			}
			_ = studentMap
		}
	})

	b.Run("PreallocatedStudentMap", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			studentMap := make(map[int]*domain.Student, 100)
			for j := 0; j < 100; j++ {
				student := &domain.Student{
					ID:     j,
					Name:   "学生",
					Age:    20,
					Gender: "男",
					Email:  "test@example.com",
					Phone:  "13800138000",
				}
				studentMap[j] = student
			}
			_ = studentMap
		}
	})
}