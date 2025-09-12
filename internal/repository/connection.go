package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Config 数据库配置
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// DB 全局数据库连接
var DB *sql.DB

// InitDB 初始化数据库连接
func InitDB() error {
	// 使用提供的连接信息
	host := "192.168.31.114"
	port := 5432
	user := "postgres"
	password := "mm152002"
	dbname := "postgres"

	var err error
	var defaultDB *sql.DB

	log.Printf("正在连接到 PostgreSQL: %s:%d", host, port)

	// 先连接到默认的postgres数据库
	defaultConnStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, "postgres")

	defaultDB, err = sql.Open("postgres", defaultConnStr)
	if err != nil {
		return fmt.Errorf("打开数据库连接失败: %v", err)
	}

	// 测试连接
	err = defaultDB.Ping()
	if err != nil {
		return fmt.Errorf("数据库连接测试失败: %v", err)
	}

	log.Printf("成功连接到PostgreSQL")

	// 创建student_management数据库（如果不存在）
	_, err = defaultDB.Exec("CREATE DATABASE " + dbname)
	if err != nil {
		// 数据库可能已存在，忽略错误
		log.Println("数据库可能已存在:", err)
	}
	defaultDB.Close()

	// 连接到student_management数据库
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("连接目标数据库失败: %v", err)
	}

	// 测试连接
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("目标数据库ping失败: %v", err)
	}

	log.Println("成功连接到student_management数据库")
	return nil
}

// CreateTables 创建数据库表
func CreateTables() error {
	log.Println("正在创建数据库表...")

	// 创建学生表
	studentsTable := `
	CREATE TABLE IF NOT EXISTS students (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		age INTEGER,
		gender VARCHAR(10),
		phone VARCHAR(20),
		email VARCHAR(100),
		address TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := DB.Exec(studentsTable)
	if err != nil {
		return fmt.Errorf("failed to create students table: %v", err)
	}

	// 创建教师表
	teachersTable := `
	CREATE TABLE IF NOT EXISTS teachers (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		subject VARCHAR(50) NOT NULL,
		email VARCHAR(100),
		phone VARCHAR(20),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = DB.Exec(teachersTable)
	if err != nil {
		return fmt.Errorf("failed to create teachers table: %v", err)
	}

	// 创建成绩表 - 简化为5个科目
	gradesTable := `
	DROP TABLE IF EXISTS grades CASCADE;
	CREATE TABLE grades (
		id SERIAL PRIMARY KEY,
		student_id INTEGER NOT NULL REFERENCES students(id) ON DELETE CASCADE,
		chinese_score DECIMAL(5,2) CHECK (chinese_score >= 0 AND chinese_score <= 100),
		math_score DECIMAL(5,2) CHECK (math_score >= 0 AND math_score <= 100),
		english_score DECIMAL(5,2) CHECK (english_score >= 0 AND english_score <= 100),
		sports_score DECIMAL(5,2) CHECK (sports_score >= 0 AND sports_score <= 100),
		music_score DECIMAL(5,2) CHECK (music_score >= 0 AND music_score <= 100),
		chinese_teacher_id INTEGER REFERENCES teachers(id),
		math_teacher_id INTEGER REFERENCES teachers(id),
		english_teacher_id INTEGER REFERENCES teachers(id),
		sports_teacher_id INTEGER REFERENCES teachers(id),
		music_teacher_id INTEGER REFERENCES teachers(id),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(student_id)
	);
	`

	_, err = DB.Exec(gradesTable)
	if err != nil {
		return fmt.Errorf("failed to create grades table: %v", err)
	}

	// 创建更新时间触发器函数
	updateFunction := `
	CREATE OR REPLACE FUNCTION update_updated_at_column()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = CURRENT_TIMESTAMP;
		RETURN NEW;
	END;
	$$ language 'plpgsql';
	`

	_, err = DB.Exec(updateFunction)
	if err != nil {
		return fmt.Errorf("failed to create update function: %v", err)
	}

	// 为学生表创建触发器
	studentTrigger := `
	DROP TRIGGER IF EXISTS update_students_updated_at ON students;
	CREATE TRIGGER update_students_updated_at
		BEFORE UPDATE ON students
		FOR EACH ROW
		EXECUTE FUNCTION update_updated_at_column();
	`

	_, err = DB.Exec(studentTrigger)
	if err != nil {
		return fmt.Errorf("failed to create students trigger: %v", err)
	}

	// 为教师表创建触发器
	teacherTrigger := `
	DROP TRIGGER IF EXISTS update_teachers_updated_at ON teachers;
	CREATE TRIGGER update_teachers_updated_at
		BEFORE UPDATE ON teachers
		FOR EACH ROW
		EXECUTE FUNCTION update_updated_at_column();
	`

	_, err = DB.Exec(teacherTrigger)
	if err != nil {
		return fmt.Errorf("failed to create teachers trigger: %v", err)
	}

	// 为成绩表创建触发器
	gradeTrigger := `
	DROP TRIGGER IF EXISTS update_grades_updated_at ON grades;
	CREATE TRIGGER update_grades_updated_at
		BEFORE UPDATE ON grades
		FOR EACH ROW
		EXECUTE FUNCTION update_updated_at_column();
	`

	_, err = DB.Exec(gradeTrigger)
	if err != nil {
		return fmt.Errorf("failed to create grades trigger: %v", err)
	}

	log.Println("数据库表创建成功")
	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
