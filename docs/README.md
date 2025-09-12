# 学生管理系统文档

## 项目概述

学生管理系统是一个基于 Go 语言开发的 RESTful API 服务，用于管理学生、教师和成绩信息。

## 目录结构

```
docs/
├── README.md              # 项目文档首页
├── api/                   # API 文档
│   ├── README.md         # API 文档说明
│   ├── students.md       # 学生相关 API
│   ├── teachers.md       # 教师相关 API
│   └── grades.md         # 成绩相关 API
├── deployment/           # 部署文档
│   ├── docker.md        # Docker 部署指南
│   ├── kubernetes.md    # Kubernetes 部署指南
│   └── production.md    # 生产环境部署
├── development/          # 开发文档
│   ├── setup.md         # 开发环境搭建
│   ├── coding-standards.md # 编码规范
│   └── testing.md       # 测试指南
└── architecture/         # 架构文档
    ├── overview.md      # 系统架构概览
    ├── database.md     # 数据库设计
    └── security.md     # 安全设计
```

## 快速开始

### 环境要求

- Go 1.21+
- PostgreSQL 12+
- Redis 6+

### 安装和运行

1. 克隆项目
```bash
git clone <repository-url>
cd student-management-system
```

2. 安装依赖
```bash
go mod download
```

3. 配置环境
```bash
cp configs/config.yaml.example configs/config.yaml
# 编辑配置文件
```

4. 运行应用
```bash
make run
# 或者
go run cmd/student-api/main.go
```

### 使用 Docker

```bash
# 构建并运行
docker-compose up -d

# 仅运行应用（开发模式）
docker-compose --profile development up -d

# 生产模式
docker-compose --profile production up -d
```

## 主要功能

- **学生管理**: 创建、查询、更新、删除学生信息
- **教师管理**: 创建、查询、更新、删除教师信息
- **成绩管理**: 创建、查询、更新、删除成绩信息
- **数据验证**: 完整的输入数据验证
- **错误处理**: 统一的错误处理和响应格式
- **日志记录**: 结构化日志记录
- **配置管理**: 灵活的配置管理系统

## API 文档

详细的 API 文档请参考 [API 文档](./api/README.md)。

## 开发指南

- [开发环境搭建](./development/setup.md)
- [编码规范](./development/coding-standards.md)
- [测试指南](./development/testing.md)

## 部署指南

- [Docker 部署](./deployment/docker.md)
- [Kubernetes 部署](./deployment/kubernetes.md)
- [生产环境部署](./deployment/production.md)

## 架构设计

- [系统架构概览](./architecture/overview.md)
- [数据库设计](./architecture/database.md)
- [安全设计](./architecture/security.md)

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](../LICENSE) 文件了解详情。

## 联系方式

- 项目维护者: [Your Name]
- 邮箱: [your.email@example.com]
- 项目链接: [https://github.com/yourusername/student-management-system]