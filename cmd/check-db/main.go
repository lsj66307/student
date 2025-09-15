package main

import (
	"database/sql"
	"fmt"
	"log"

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
		log.Fatal("数据库连接失败:", err)
	}
	defer db.Close()

	// 测试连接
	err = db.Ping()
	if err != nil {
		log.Fatal("数据库连接测试失败:", err)
	}
	fmt.Println("数据库连接成功!")

	// 检查学生表数据
	fmt.Println("\n=== 检查学生表数据 ===")
	rows, err := db.Query("SELECT student_id, name, email, major FROM students ORDER BY student_id")
	if err != nil {
		log.Fatal("查询学生数据失败:", err)
	}
	defer rows.Close()

	studentCount := 0
	fmt.Println("学生ID\t\t姓名\t\t邮箱\t\t\t专业")
	fmt.Println("------------------------------------------------------------")

	for rows.Next() {
		var studentID, name, email, major string
		err := rows.Scan(&studentID, &name, &email, &major)
		if err != nil {
			log.Fatal("扫描学生数据失败:", err)
		}
		fmt.Printf("%s\t%s\t\t%s\t%s\n", studentID, name, email, major)
		studentCount++
	}

	fmt.Printf("\n学生总数: %d\n", studentCount)

	// 检查教师表数据
	fmt.Println("\n=== 检查教师表数据 ===")
	rows2, err := db.Query("SELECT id, name, email, subject, department FROM teachers ORDER BY id")
	if err != nil {
		log.Fatal("查询教师数据失败:", err)
	}
	defer rows2.Close()

	teacherCount := 0
	fmt.Println("教师ID\t姓名\t\t邮箱\t\t\t\t科目\t\t部门")
	fmt.Println("------------------------------------------------------------------------")

	for rows2.Next() {
		var teacherID int
		var name, email, subject, department string
		err := rows2.Scan(&teacherID, &name, &email, &subject, &department)
		if err != nil {
			log.Fatal("扫描教师数据失败:", err)
		}
		fmt.Printf("%d\t%s\t\t%s\t\t%s\t\t%s\n", teacherID, name, email, subject, department)
		teacherCount++
	}

	fmt.Printf("\n教师总数: %d\n", teacherCount)

	// 检查表结构
	fmt.Println("\n=== 检查学生表结构 ===")
	rows3, err := db.Query(`
		SELECT column_name, data_type, is_nullable 
		FROM information_schema.columns 
		WHERE table_name = 'students' 
		ORDER BY ordinal_position
	`)
	if err != nil {
		log.Fatal("查询学生表结构失败:", err)
	}
	defer rows3.Close()

	fmt.Println("列名\t\t\t数据类型\t\t可空")
	fmt.Println("--------------------------------------------")
	for rows3.Next() {
		var columnName, dataType, isNullable string
		err := rows3.Scan(&columnName, &dataType, &isNullable)
		if err != nil {
			log.Fatal("扫描表结构失败:", err)
		}
		fmt.Printf("%s\t\t%s\t\t%s\n", columnName, dataType, isNullable)
	}

	fmt.Println("\n=== 检查教师表结构 ===")
	rows4, err := db.Query(`
		SELECT column_name, data_type, is_nullable 
		FROM information_schema.columns 
		WHERE table_name = 'teachers' 
		ORDER BY ordinal_position
	`)
	if err != nil {
		log.Fatal("查询教师表结构失败:", err)
	}
	defer rows4.Close()

	fmt.Println("列名\t\t\t数据类型\t\t可空")
	fmt.Println("--------------------------------------------")
	for rows4.Next() {
		var columnName, dataType, isNullable string
		err := rows4.Scan(&columnName, &dataType, &isNullable)
		if err != nil {
			log.Fatal("扫描表结构失败:", err)
		}
		fmt.Printf("%s\t\t%s\t\t%s\n", columnName, dataType, isNullable)
	}

	fmt.Println("\n数据库检查完成!")
}
