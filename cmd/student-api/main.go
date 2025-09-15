package main

import (
	"log"
	"student-management-system/internal/config"
	handler "student-management-system/internal/handler"
	repo "student-management-system/internal/repository"
)

func main() {
	// 加载配置
	log.Println("正在加载配置...")
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库连接
	log.Println("正在初始化数据库连接...")
	err = repo.InitDB()
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer repo.CloseDB()

	// 初始化Redis连接
	log.Println("正在初始化Redis连接...")
	err = repo.InitRedis(cfg)
	if err != nil {
		log.Fatalf("Redis连接失败: %v", err)
	}
	defer repo.CloseRedis()

	// 创建数据库表
	log.Println("正在创建数据库表...")
	err = repo.CreateTables()
	if err != nil {
		log.Fatalf("创建数据库表失败: %v", err)
	}

	// 设置路由
	log.Println("正在设置路由...")
	router := handler.SetupRoutes(cfg)

	// 启动服务器
	log.Println("学生管理系统启动中...")
	log.Println("服务器运行在: http://localhost:3060")
	log.Println("API文档: http://localhost:3060/")
	log.Println("健康检查: http://localhost:3060/health")
	log.Println("按 Ctrl+C 停止服务器")

	// 在端口3060启动服务器
	if err := router.Run(":3060"); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
