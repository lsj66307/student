package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	// 数据库连接配置
	host := "192.168.31.114"
	port := 5432
	user := "postgres"
	password := "mm152002"
	dbname := "student_management"

	// 构建连接字符串
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// 连接数据库
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}
	defer db.Close()

	// 测试连接
	err = db.Ping()
	if err != nil {
		log.Fatal("数据库连接测试失败:", err)
	}

	log.Println("数据库连接成功")

	// 删除现有表并重新创建
	sqlCommands := []string{
		"DROP TABLE IF EXISTS teachers CASCADE",
		"DROP TABLE IF EXISTS students CASCADE",
		`CREATE TABLE IF NOT EXISTS students (
			id SERIAL PRIMARY KEY,
			student_id VARCHAR(20) UNIQUE,
			name VARCHAR(100) NOT NULL,
			age INTEGER,
			gender VARCHAR(10),
			phone VARCHAR(20),
			email VARCHAR(100),
			address TEXT,
			major VARCHAR(100),
			enrollment_date DATE,
			graduation_date DATE,
			status VARCHAR(20) DEFAULT 'active',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS teachers (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			age INTEGER,
			gender VARCHAR(10),
			email VARCHAR(100),
			phone VARCHAR(20),
			subject VARCHAR(50) NOT NULL,
			title VARCHAR(50),
			department VARCHAR(100),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	// 执行SQL命令
	for i, cmd := range sqlCommands {
		preview := cmd
		if len(cmd) > 50 {
			preview = cmd[:50] + "..."
		}
		log.Printf("执行SQL命令 %d: %s", i+1, preview)
		_, err := db.Exec(cmd)
		if err != nil {
			log.Printf("执行SQL命令失败: %v", err)
			os.Exit(1)
		}
		log.Printf("SQL命令 %d 执行成功", i+1)
	}

	log.Println("数据库表结构修复完成!")

	// 验证表结构
	log.Println("验证teachers表结构...")
	rows, err := db.Query("SELECT column_name, data_type FROM information_schema.columns WHERE table_name = 'teachers' ORDER BY ordinal_position")
	if err != nil {
		log.Printf("查询teachers表结构失败: %v", err)
	} else {
		log.Println("teachers表结构:")
		for rows.Next() {
			var columnName, dataType string
			rows.Scan(&columnName, &dataType)
			log.Printf("  %s: %s", columnName, dataType)
		}
		rows.Close()
	}

	log.Println("验证students表结构...")
	rows, err = db.Query("SELECT column_name, data_type FROM information_schema.columns WHERE table_name = 'students' ORDER BY ordinal_position")
	if err != nil {
		log.Printf("查询students表结构失败: %v", err)
	} else {
		log.Println("students表结构:")
		for rows.Next() {
			var columnName, dataType string
			rows.Scan(&columnName, &dataType)
			log.Printf("  %s: %s", columnName, dataType)
		}
		rows.Close()
	}
}