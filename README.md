# 学生管理系统

一个基于 Go 语言开发的学生管理系统，提供学生信息的增删改查功能。

## 功能特性

- ✅ 创建学生信息
- ✅ 查看学生列表（支持分页）
- ✅ 查看单个学生详情
- ✅ 更新学生信息
- ✅ 删除学生信息
- ✅ RESTful API 设计
- ✅ SQLite 数据库存储
- ✅ CORS 跨域支持

## 技术栈

- **语言**: Go 1.25.0
- **Web 框架**: Gin
- **数据库**: SQLite
- **ORM**: 原生 SQL

## 项目结构

```
student-management-system/
├── main.go                 # 程序入口
├── go.mod                  # Go模块文件
├── student_management.db   # SQLite数据库文件
├── database/
│   └── connection.go       # 数据库连接和表创建
├── models/
│   ├── student.go          # 学生数据模型
│   └── student_service.go  # 学生业务逻辑
└── handlers/
    ├── student_handler.go  # HTTP处理函数
    └── routes.go           # 路由配置
```

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 启动服务器

```bash
go run main.go
```

服务器将在 `http://localhost:3060` 启动

### 3. 测试 API

访问 `http://localhost:3060` 查看 API 文档
访问 `http://localhost:3060/health` 进行健康检查

## API 接口文档

### 基础信息

- **Base URL**: `http://localhost:3060`
- **Content-Type**: `application/json`

### 接口列表

#### 1. 创建学生

```http
POST /api/students
```

**请求体**:

```json
{
  "name": "张三",
  "age": 20,
  "gender": "男",
  "email": "zhangsan@example.com",
  "phone": "13800138000",
  "major": "计算机科学",
  "grade": "大二"
}
```

**响应**:

```json
{
  "code": 201,
  "message": "学生创建成功",
  "data": {
    "id": 1,
    "name": "张三",
    "age": 20,
    "gender": "男",
    "email": "zhangsan@example.com",
    "phone": "13800138000",
    "major": "计算机科学",
    "grade": "大二",
    "created_at": "2025-09-12T07:07:59Z",
    "updated_at": "2025-09-12T07:07:59Z"
  }
}
```

#### 2. 获取学生列表

```http
GET /api/students?page=1&size=10
```

**查询参数**:

- `page`: 页码（默认：1）
- `size`: 每页数量（默认：10，最大：100）

**响应**:

```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "id": 1,
      "name": "张三",
      "age": 20,
      "gender": "男",
      "email": "zhangsan@example.com",
      "phone": "13800138000",
      "major": "计算机科学",
      "grade": "大二",
      "created_at": "2025-09-12T07:07:59Z",
      "updated_at": "2025-09-12T07:07:59Z"
    }
  ],
  "total": 1,
  "page": 1,
  "size": 10
}
```

#### 3. 获取单个学生

```http
GET /api/students/{id}
```

**响应**:

```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": 1,
    "name": "张三",
    "age": 20,
    "gender": "男",
    "email": "zhangsan@example.com",
    "phone": "13800138000",
    "major": "计算机科学",
    "grade": "大二",
    "created_at": "2025-09-12T07:07:59Z",
    "updated_at": "2025-09-12T07:07:59Z"
  }
}
```

#### 4. 更新学生信息

```http
PUT /api/students/{id}
```

**请求体**:

```json
{
  "name": "李四",
  "age": 21,
  "major": "软件工程"
}
```

#### 5. 删除学生

```http
DELETE /api/students/{id}
```

**响应**:

```json
{
  "code": 200,
  "message": "删除成功"
}
```

## 测试示例

### 使用 curl 测试

```bash
# 创建学生
curl -X POST http://localhost:3060/api/students \
  -H "Content-Type: application/json" \
  -d '{
    "name": "张三",
    "age": 20,
    "gender": "男",
    "email": "zhangsan@example.com",
    "phone": "13800138000",
    "major": "计算机科学",
    "grade": "大二"
  }'

# 获取学生列表
curl -X GET http://localhost:3060/api/students

# 获取单个学生
curl -X GET http://localhost:3060/api/students/1

# 更新学生信息
curl -X PUT http://localhost:3060/api/students/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "李四", "age": 21}'

# 删除学生
curl -X DELETE http://localhost:3060/api/students/1
```

## 数据库设计

### students 表结构

| 字段名     | 类型     | 说明     | 约束         |
| ---------- | -------- | -------- | ------------ |
| id         | INTEGER  | 学生 ID  | 主键，自增   |
| name       | TEXT     | 姓名     | 非空         |
| age        | INTEGER  | 年龄     | 非空，1-150  |
| gender     | TEXT     | 性别     | 非空，男/女  |
| email      | TEXT     | 邮箱     | 非空，唯一   |
| phone      | TEXT     | 电话     | 非空         |
| major      | TEXT     | 专业     | 非空         |
| grade      | TEXT     | 年级     | 非空         |
| created_at | DATETIME | 创建时间 | 默认当前时间 |
| updated_at | DATETIME | 更新时间 | 自动更新     |

## 错误码说明

| 错误码 | 说明           |
| ------ | -------------- |
| 200    | 请求成功       |
| 201    | 创建成功       |
| 400    | 请求参数错误   |
| 404    | 资源不存在     |
| 500    | 服务器内部错误 |

## 开发说明

### 环境要求

- Go 1.25.0 或更高版本
- SQLite3

### 部署说明

1. 编译项目：`go build -o student-management-system`
2. 运行：`./student-management-system`
3. 服务器将在端口 3060 启动

### 注意事项

- 项目使用 SQLite 数据库，数据文件为 `student_management.db`
- 首次运行会自动创建数据库表
- 支持 CORS 跨域请求
- 生产环境建议设置 `GIN_MODE=release`

## 许可证

MIT License
