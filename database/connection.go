package database

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
	// 创建学生表
	studentQuery := `
	CREATE TABLE IF NOT EXISTS students (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		age INTEGER NOT NULL CHECK (age > 0 AND age < 150),
		gender VARCHAR(10) NOT NULL CHECK (gender IN ('男', '女')),
		email VARCHAR(255) UNIQUE NOT NULL,
		phone VARCHAR(20) NOT NULL,
		major VARCHAR(100) NOT NULL,
		grade VARCHAR(50) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := DB.Exec(studentQuery)
	if err != nil {
		return fmt.Errorf("failed to create students table: %v", err)
	}

	// 创建老师表
	teacherQuery := `
	CREATE TABLE IF NOT EXISTS teachers (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		age INTEGER NOT NULL CHECK (age >= 22 AND age <= 70),
		gender VARCHAR(10) NOT NULL CHECK (gender IN ('男', '女')),
		email VARCHAR(255) UNIQUE NOT NULL,
		phone VARCHAR(20) NOT NULL,
		subject VARCHAR(100) NOT NULL,
		title VARCHAR(50) NOT NULL,
		department VARCHAR(100) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = DB.Exec(teacherQuery)
	if err != nil {
		return fmt.Errorf("failed to create teachers table: %v", err)
	}

	// 创建更新时间触发器函数
	triggerFuncSQL := `
	CREATE OR REPLACE FUNCTION update_updated_at_column()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = CURRENT_TIMESTAMP;
		RETURN NEW;
	END;
	$$ language 'plpgsql';
	`

	_, err = DB.Exec(triggerFuncSQL)
	if err != nil {
		return fmt.Errorf("failed to create trigger function: %v", err)
	}

	// 创建学生表触发器
	studentTriggerSQL := `
	DROP TRIGGER IF EXISTS update_students_updated_at ON students;
	CREATE TRIGGER update_students_updated_at
		BEFORE UPDATE ON students
		FOR EACH ROW
		EXECUTE FUNCTION update_updated_at_column();
	`

	_, err = DB.Exec(studentTriggerSQL)
	if err != nil {
		return fmt.Errorf("failed to create student trigger: %v", err)
	}

	// 创建老师表触发器
	teacherTriggerSQL := `
	DROP TRIGGER IF EXISTS update_teachers_updated_at ON teachers;
	CREATE TRIGGER update_teachers_updated_at
		BEFORE UPDATE ON teachers
		FOR EACH ROW
		EXECUTE FUNCTION update_updated_at_column();
	`

	_, err = DB.Exec(teacherTriggerSQL)
	if err != nil {
		return fmt.Errorf("failed to create teacher trigger: %v", err)
	}

	log.Println("Database tables created successfully")
	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
